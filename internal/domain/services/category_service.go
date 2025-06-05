package services

import (
	"github.com/Mutonya/Savanah/internal/domain/models"
	"github.com/Mutonya/Savanah/internal/domain/repositories"
)

type CategoryService interface {
	CreateCategory(req *CategoryCreateRequest) (*models.Category, error)
	GetCategory(id uint) (*models.Category, error)
	GetCategories() ([]models.Category, error)
	UpdateCategory(id uint, req *CategoryUpdateRequest) (*models.Category, error)
	DeleteCategory(id uint) error
	GetCategoryProducts(categoryID uint, page, limit int) ([]models.Product, int64, error)
	GetAveragePrice(categoryID uint) (float64, error)
}

type CategoryCreateRequest struct {
	Name     string `json:"name" binding:"required"`
	ParentID *uint  `json:"parent_id"`
}

type CategoryUpdateRequest struct {
	Name     string `json:"name"`
	ParentID *uint  `json:"parent_id"`
}

type categoryService struct {
	categoryRepo repositories.CategoryRepository
}

func NewCategoryService(categoryRepo repositories.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo: categoryRepo}
}

func (s *categoryService) CreateCategory(req *CategoryCreateRequest) (*models.Category, error) {
	category := &models.Category{
		Name:     req.Name,
		ParentID: req.ParentID,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) GetCategory(id uint) (*models.Category, error) {
	return s.categoryRepo.GetByID(id)
}

func (s *categoryService) GetCategories() ([]models.Category, error) {
	return s.categoryRepo.GetAll()
}

func (s *categoryService) UpdateCategory(id uint, req *CategoryUpdateRequest) (*models.Category, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.ParentID != nil {
		category.ParentID = req.ParentID
	}

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) DeleteCategory(id uint) error {
	return s.categoryRepo.Delete(id)
}

func (s *categoryService) GetCategoryProducts(categoryID uint, page, limit int) ([]models.Product, int64, error) {
	return s.categoryRepo.GetProducts(categoryID, page, limit)
}

func (s *categoryService) GetAveragePrice(categoryID uint) (float64, error) {
	return s.categoryRepo.GetAveragePrice(categoryID)
}
