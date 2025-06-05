package repositories

import (
	"context"
	"gorm.io/gorm"

	"github.com/Mutonya/Savanah/internal/domain/models"
)

type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	GetByID(ctx context.Context, id uint) (*models.Product, error)
	GetAll(ctx context.Context, page, limit int) ([]models.Product, int64, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id uint) error
	GetByCategory(ctx context.Context, categoryID uint, page, limit int) ([]models.Product, int64, error)
	GetAveragePrice(ctx context.Context, categoryID uint) (float64, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) GetByID(ctx context.Context, id uint) (*models.Product, error) {
	var product models.Product
	if err := r.db.WithContext(ctx).Preload("Category").First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetAll(ctx context.Context, page, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var count int64

	offset := (page - 1) * limit

	tx := r.db.WithContext(ctx).Model(&models.Product{}).Preload("Category").
		Offset(offset).
		Limit(limit)

	if err := tx.Find(&products).Error; err != nil {
		return nil, 0, err
	}
	if err := tx.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return products, count, nil
}

func (r *productRepository) Update(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *productRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Product{}, id).Error
}

func (r *productRepository) GetByCategory(ctx context.Context, categoryID uint, page, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var count int64

	offset := (page - 1) * limit

	tx := r.db.WithContext(ctx).Preload("Category").
		Where("category_id = ?", categoryID).
		Offset(offset).
		Limit(limit)

	if err := tx.Find(&products).Error; err != nil {
		return nil, 0, err
	}
	if err := tx.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return products, count, nil
}

func (r *productRepository) GetAveragePrice(ctx context.Context, categoryID uint) (float64, error) {
	var avgPrice float64
	if err := r.db.WithContext(ctx).Model(&models.Product{}).
		Where("category_id = ?", categoryID).
		Select("AVG(price)").
		Scan(&avgPrice).Error; err != nil {
		return 0, err
	}
	return avgPrice, nil
}
