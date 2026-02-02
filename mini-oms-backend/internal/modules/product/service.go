package product

import (
	"errors"
	"mini-oms-backend/internal/models"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type ProductRequest struct {
	CategoryID  *uuid.UUID `json:"category_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	Stock       int        `json:"stock"`
	ImageURL    string     `json:"image_url"`
}

func (s *Service) GetAll() ([]models.Product, error) {
	return s.repo.FindAll()
}

func (s *Service) GetByID(id uuid.UUID) (*models.Product, error) {
	return s.repo.FindByID(id)
}

func (s *Service) Create(req *ProductRequest) (*models.Product, error) {
	if req.Name == "" || req.Price <= 0 || req.Stock < 0 {
		return nil, errors.New("invalid product data")
	}

	product := &models.Product{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
	}

	if err := s.repo.Create(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *Service) Update(id uuid.UUID, req *ProductRequest) (*models.Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("product not found")
	}

	product.CategoryID = req.CategoryID
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock
	product.ImageURL = req.ImageURL

	if err := s.repo.Update(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *Service) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
