package order

import (
	"errors"
	"mini-oms-backend/internal/models"

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
		// Get product
		product, err := s.repo.FindProductByID(item.ProductID)
		if err != nil {
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

	// Create order
	order := &models.Order{
		UserID:      userID,
		TotalAmount: totalAmount,
		Status:      "created",
		Notes:       req.Notes,
		OrderItems:  orderItems,
	}

	if err := tx.Create(order).Error; err != nil {
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
