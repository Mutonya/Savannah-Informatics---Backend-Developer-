package services

import (
	"context"

	"github.com/Mutonya/Savanah/internal/domain/models"
	"github.com/Mutonya/Savanah/internal/domain/repositories"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req *ProductCreateRequest) (*models.Product, error)
	GetProduct(ctx context.Context, id uint) (*models.Product, error)
	GetProducts(ctx context.Context, page, limit int) ([]models.Product, int64, error)
	UpdateProduct(ctx context.Context, id uint, req *ProductUpdateRequest) (*models.Product, error)
	DeleteProduct(ctx context.Context, id uint) error
}

type ProductCreateRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	SKU         string  `json:"sku" binding:"required"`
	CategoryID  uint    `json:"category_id" binding:"required"`
}

type ProductUpdateRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"gt=0"`
	SKU         string  `json:"sku"`
	CategoryID  uint    `json:"category_id"`
}

type productService struct {
	productRepo repositories.ProductRepository
}

func NewProductService(productRepo repositories.ProductRepository) ProductService {
	return &productService{productRepo: productRepo}
}

func (s *productService) CreateProduct(ctx context.Context, req *ProductCreateRequest) (*models.Product, error) {
	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		SKU:         req.SKU,
		CategoryID:  req.CategoryID,
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) GetProduct(ctx context.Context, id uint) (*models.Product, error) {
	return s.productRepo.GetByID(ctx, id)
}

func (s *productService) GetProducts(ctx context.Context, page, limit int) ([]models.Product, int64, error) {
	return s.productRepo.GetAll(ctx, page, limit)
}

func (s *productService) UpdateProduct(ctx context.Context, id uint, req *ProductUpdateRequest) (*models.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.SKU != "" {
		product.SKU = req.SKU
	}
	if req.CategoryID > 0 {
		product.CategoryID = req.CategoryID
	}

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) DeleteProduct(ctx context.Context, id uint) error {
	return s.productRepo.Delete(ctx, id)
}
