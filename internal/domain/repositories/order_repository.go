package repositories

import (
	"context"
	"gorm.io/gorm"

	"github.com/Mutonya/Savanah/internal/domain/models"
)

type OrderRepository interface {
	Create(ctx context.Context, order *models.Order) error
	GetByID(ctx context.Context, id uint) (*models.Order, error)
	GetByCustomerID(ctx context.Context, customerID uint, page, limit int) ([]models.Order, int64, error)
	Update(ctx context.Context, order *models.Order) error
	UpdateStatus(ctx context.Context, orderID uint, status models.OrderStatus) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *orderRepository) GetByID(ctx context.Context, id uint) (*models.Order, error) {
	var order models.Order
	if err := r.db.WithContext(ctx).Preload("Customer").
		Preload("OrderItems").
		Preload("OrderItems.Product").
		First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetByCustomerID(ctx context.Context, customerID uint, page, limit int) ([]models.Order, int64, error) {
	var orders []models.Order
	var count int64

	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).Preload("OrderItems").
		Preload("OrderItems.Product").
		Where("customer_id = ?", customerID).
		Offset(offset).
		Limit(limit).
		Find(&orders).
		Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return orders, count, nil
}

func (r *orderRepository) Update(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

func (r *orderRepository) UpdateStatus(ctx context.Context, orderID uint, status models.OrderStatus) error {
	return r.db.WithContext(ctx).Model(&models.Order{}).
		Where("id = ?", orderID).
		Update("status", status).Error
}
