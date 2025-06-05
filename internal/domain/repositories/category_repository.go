package repositories

import (
	"gorm.io/gorm"

	"github.com/Mutonya/Savanah/internal/domain/models"
)

type CategoryRepository interface {
	Create(category *models.Category) error
	GetByID(id uint) (*models.Category, error)
	GetAll() ([]models.Category, error)
	Update(category *models.Category) error
	Delete(id uint) error
	GetProducts(categoryID uint, page, limit int) ([]models.Product, int64, error)
	GetAveragePrice(categoryID uint) (float64, error)
	GetSubcategories(parentID uint) ([]models.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) GetByID(id uint) (*models.Category, error) {
	var category models.Category
	if err := r.db.Preload("Children").Preload("Products").
		First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetAll() ([]models.Category, error) {
	var categories []models.Category
	if err := r.db.Where("parent_id IS NULL").
		Preload("Children").
		Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}

func (r *categoryRepository) GetProducts(categoryID uint, page, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var count int64

	offset := (page - 1) * limit

	// Get all subcategory IDs
	var subcategoryIDs []uint
	if err := r.db.Model(&models.Category{}).
		Where("parent_id = ?", categoryID).
		Pluck("id", &subcategoryIDs).Error; err != nil {
		return nil, 0, err
	}

	// Include the parent category ID
	allCategoryIDs := append(subcategoryIDs, categoryID)

	query := r.db.Preload("Category").
		Where("category_id IN ?", allCategoryIDs)

	if err := query.
		Offset(offset).
		Limit(limit).
		Find(&products).
		Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return products, count, nil
}

func (r *categoryRepository) GetAveragePrice(categoryID uint) (float64, error) {
	// Get all subcategory IDs
	var subcategoryIDs []uint
	if err := r.db.Model(&models.Category{}).
		Where("parent_id = ?", categoryID).
		Pluck("id", &subcategoryIDs).Error; err != nil {
		return 0, err
	}

	// Include the parent category ID
	allCategoryIDs := append(subcategoryIDs, categoryID)

	var avgPrice float64
	if err := r.db.Model(&models.Product{}).
		Where("category_id IN ?", allCategoryIDs).
		Select("AVG(price)").
		Scan(&avgPrice).Error; err != nil {
		return 0, err
	}

	return avgPrice, nil
}

func (r *categoryRepository) GetSubcategories(parentID uint) ([]models.Category, error) {
	var categories []models.Category
	if err := r.db.Where("parent_id = ?", parentID).
		Preload("Children").
		Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
