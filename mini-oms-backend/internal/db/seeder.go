package db

import (
	"fmt"
	"log"
	"math/rand"
	"mini-oms-backend/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	seedUsers(db)
	seedProducts(db)
	fixOrderNumbers(db) // Fix data lama
}

func seedUsers(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Where("role = ?", "admin").Count(&count)
	if count > 0 {
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := models.User{
		Name:     "Admin Tofuu",
		Email:    "admin@whaleestore.com",
		Password: string(hashedPassword),
		Role:     "admin",
	}
	db.Create(&admin)
	log.Println("Seeded Admin User: admin@whaleestore.com / admin123")
}

func seedProducts(db *gorm.DB) {
	var count int64
	db.Model(&models.Product{}).Count(&count)
	if count > 0 {
		return
	}

	products := []models.Product{
		{Name: "Whale Plushie Huge", Price: 250000, Stock: 50, Description: "Super soft huge whale plushie"},
		{Name: "Ocean Blue Hoodie", Price: 350000, Stock: 20, Description: "Stylish hoodie with whale logo"},
		{Name: "Marine Tumbler", Price: 120000, Stock: 100, Description: "Keep your drinks cold"},
	}

	for _, p := range products {
		db.Create(&p)
	}
	log.Println("Seeded Products")
}

func fixOrderNumbers(db *gorm.DB) {
	var orders []models.Order
	// Find orders with empty order number
	if err := db.Where("order_number = ? OR order_number IS NULL", "").Find(&orders).Error; err != nil {
		return
	}

	for _, order := range orders {
		newNumber := fmt.Sprintf("ORD-FIX-%s-%d", order.CreatedAt.Format("20060102"), rand.Intn(10000))
		db.Model(&order).Update("order_number", newNumber)
		log.Printf("Fixed OrderNumber for %s: %s\n", order.ID, newNumber)
	}
}
