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

// db: Holds the database connection
//
// Constructor: Creates new repository instance with dependency injection
type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *models.Category) error {
	//WithContext(ctx) propagate ctx ensures the DB operation can be cancelled or timed out from upstream
	return r.db.Create(category).Error
}

// Preload("Children"): Eager loads subcategories
//
// Preload("Products"): Eager loads products
//
// Returns full category hierarchy in one query
func (r *categoryRepository) GetByID(id uint) (*models.Category, error) {
	var category models.Category
	if err := r.db.Preload("Children").Preload("Products").
		First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

//Retrieves only root categories (no parent)

//Eager loads immediate children

// Forms hierarchical structure
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

//Gets products from category + all subcategories

//Returns:

//Products slice (for current page)

//Total count (for pagination UI)

// Error (if any)
func (r *categoryRepository) GetProducts(categoryID uint, page, limit int) ([]models.Product, int64, error) {
	//Gets products from category + all subcategories
	var products []models.Product
	var count int64
	//Products slice (for current page)
	offset := (page - 1) * limit

	// Get all subcategory IDs
	var subcategoryIDs []uint
	if err := r.db.Model(&models.Category{}).
		Where("parent_id = ?", categoryID).
		Pluck("id", &subcategoryIDs).Error; err != nil { //Pluck queries a single column from a model, returning in the slice dest
		return nil, 0, err
	}

	allCategoryIDs := append(subcategoryIDs, categoryID)

	query := r.db.
		Preload("Category").
		Where("category_id IN ?", allCategoryIDs)
	//Order("id ASC") // or "id DESC" for descending order
	//Total count (for pagination UI)
	if err := query.
		Count(&count).
		Offset(offset).
		Limit(limit).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, count, nil
}

/*
		Calculates average price across category hierarchy

	 	# Uses SQL AVG() function for efficiency

	 	Includes all subcategories
*/
func (r *categoryRepository) GetAveragePrice(categoryID uint) (float64, error) {
	//Calculates average price across category hierarchy
	var subcategoryIDs []uint
	if err := r.db.Model(&models.Category{}).
		Where("parent_id = ?", categoryID).
		Pluck("id", &subcategoryIDs).Error; err != nil {
		return 0, err
	}
	//Includes all subcategories

	allCategoryIDs := append(subcategoryIDs, categoryID)
	//Uses SQL AVG() function for efficiency
	var avgPrice float64
	if err := r.db.Model(&models.Product{}).
		Where("category_id IN ?", allCategoryIDs).
		Select("AVG(price)").
		Scan(&avgPrice).Error; err != nil {
		return 0, err
	}

	return avgPrice, nil
}

/*
Fetches direct children of a category

Preloads grandchildren (recursive structure)

Forms hierarchical tree
*/
func (r *categoryRepository) GetSubcategories(parentID uint) ([]models.Category, error) {
	var categories []models.Category
	if err := r.db.Where("parent_id = ?", parentID).
		Preload("Children"). //Loads relationships in single query
		Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
