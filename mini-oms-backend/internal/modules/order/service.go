package order

import (
	"errors"
	"fmt"
	"math/rand"
	"mini-oms-backend/internal/models"
	"mini-oms-backend/internal/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
	db   *gorm.DB
}

func NewService(repo *Repository, db *gorm.DB) *Service {
	return &Service{
		repo: repo,
		db:   db,
	}
}

type OrderItemRequest struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
}

type CreateOrderRequest struct {
	Items []OrderItemRequest `json:"items"`
	Notes string             `json:"notes"`
}

func (s *Service) GetAll() ([]models.Order, error) {
	return s.repo.FindAll()
}

func (s *Service) GetByUserID(userID uuid.UUID) ([]models.Order, error) {
	return s.repo.FindByUserID(userID)
}

func (s *Service) GetByID(id uuid.UUID) (*models.Order, error) {
	return s.repo.FindByID(id)
}

func (s *Service) CreateOrder(userID uuid.UUID, req *CreateOrderRequest) (*models.Order, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("order must have at least one item")
	}

	// Start database transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var totalAmount float64
	var orderItems []models.OrderItem

	// Process each order item
	for _, item := range req.Items {
		// Get product with ROW LOCK to prevent race conditions
		// We use s.repo.FindProductByIDWithLock (we need to expose this via interface or access repo directly)
		// Since product repo is separate, we'll assume we can call it differently or query directly via tx
		var product models.Product
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&product, "id = ?", item.ProductID).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("product not found")
		}

		// Check stock
		if !product.IsInStock(item.Quantity) {
			tx.Rollback()
			return nil, errors.New("insufficient stock for product: " + product.Name)
		}

		// Calculate subtotal
		subtotal := product.Price * float64(item.Quantity)
		totalAmount += subtotal

		// Create order item with product snapshot
		orderItem := models.OrderItem{
			ProductID:    product.ID,
			ProductName:  product.Name,
			ProductPrice: product.Price,
			Quantity:     item.Quantity,
			Subtotal:     subtotal,
		}
		orderItems = append(orderItems, orderItem)

		// Reduce stock
		if err := product.ReduceStock(item.Quantity); err != nil {
			tx.Rollback()
			return nil, err
		}

		// Update product stock in transaction
		if err := tx.Save(product).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Generate Order Number: ORD-YYYYMMDD-HHMMSS-XXXX (Collision resistant)
	// We use local RNG to avoid global lock contention and seed it with time
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	orderNumber := fmt.Sprintf("ORD-%s-%04d", time.Now().Format("20060102-150405"), rng.Intn(10000))

	// Create order
	order := &models.Order{
		UserID:      userID,
		OrderNumber: orderNumber,
		TotalAmount: totalAmount,
		Status:      models.OrderStatusCreated,
		Notes:       req.Notes,
		OrderItems:  orderItems,
	}

	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Log Audit
	if err := utils.LogAudit(tx, userID, "ORDER_CREATED", "Order", order.ID, fmt.Sprintf("Order created with %d items", len(orderItems))); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Reload order with relations
	return s.repo.FindByID(order.ID)
}

func (s *Service) CancelOrder(orderID, userID uuid.UUID, role string) error {
	// Find order
	order, err := s.repo.FindByID(orderID)
	if err != nil {
		return errors.New("order not found")
	}

	// Permission check (only owner or admin can cancel)
	if role != "admin" && order.UserID != userID {
		return errors.New("unauthorized")
	}

	// Check status (can only cancel if created or processing)
	// If already paid (processing), usually needs refund logic, but for now we just cancel
	if order.Status == "completed" || order.Status == "canceled" {
		return errors.New("cannot cancel completed or already canceled order")
	}

	// Start transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Update Order Status
	order.Status = "canceled"
	if err := tx.Save(order).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 2. Restore Stock
	for _, item := range order.OrderItems {
		// Find product with LOCK
		var product models.Product
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&product, "id = ?", item.ProductID).Error; err != nil {
			// If product not found/deleted, we skip restocking and continue
			continue
		}

		product.IncreaseStock(item.Quantity)

		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Log Audit
	if err := utils.LogAudit(tx, userID, "ORDER_CANCELED", "Order", order.ID, "Order canceled by user/admin"); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *Service) GetStats() (map[string]interface{}, error) {
	var totalOrders int64
	var totalRevenue float64
	var pendingPayments int64

	// Total Orders
	s.db.Model(&models.Order{}).Count(&totalOrders)

	// Total Revenue (Sum of Orders with status processing, completed, paid, or success)
	// Adding 'paid' and 'success' for backward compatibility or if data manually seeded
	statuses := []string{
		models.OrderStatusProcessing,
		models.OrderStatusCompleted,
		"paid",
		"success",
	}
	s.db.Model(&models.Order{}).Where("status IN ?", statuses).Select("COALESCE(SUM(total_amount), 0)").Scan(&totalRevenue)

	// Pending Payments (Status Created + Payment Pending)
	s.db.Model(&models.Payment{}).Where("status = ?", "pending").Count(&pendingPayments)

	fmt.Printf("DEBUG STATS: Orders=%d, Revenue=%.2f, Pending=%d\n", totalOrders, totalRevenue, pendingPayments)

	return map[string]interface{}{
		"total_orders":     totalOrders,
		"total_revenue":    totalRevenue,
		"pending_payments": pendingPayments,
	}, nil
}
