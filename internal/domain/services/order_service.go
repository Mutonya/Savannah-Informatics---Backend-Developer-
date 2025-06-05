package services

import (
	"context"
	"github.com/Mutonya/Savanah/internal/domain/models"
	"github.com/Mutonya/Savanah/internal/domain/repositories"
)

type OrderService interface {
	CreateOrder(ctx context.Context, customerID uint, req *OrderCreateRequest) (*models.Order, error)
	GetOrder(ctx context.Context, customerID, orderID uint) (*models.Order, error)
	GetOrders(ctx context.Context, customerID uint, page, limit int) ([]models.Order, int64, error)
	UpdateOrderStatus(ctx context.Context, orderID uint, status models.OrderStatus) (*models.Order, error)
}

type OrderCreateRequest struct {
	Items []OrderItemRequest `json:"items" binding:"required,min=1"`
}
type OrderStatusUpdateRequest struct {
	Status models.OrderStatus `json:"status" binding:"required,oneof=pending completed cancelled"`
}

type OrderItemRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

type orderService struct {
	orderRepo    repositories.OrderRepository
	productRepo  repositories.ProductRepository
	customerRepo repositories.CustomerRepository
	notifier     NotificationService
}

func NewOrderService(
	orderRepo repositories.OrderRepository,
	productRepo repositories.ProductRepository,
	customerRepo repositories.CustomerRepository,
	notifier NotificationService,
) OrderService {
	return &orderService{
		orderRepo:    orderRepo,
		productRepo:  productRepo,
		customerRepo: customerRepo,
		notifier:     notifier,
	}
}
func (s *orderService) CreateOrder(ctx context.Context, customerID uint, req *OrderCreateRequest) (*models.Order, error) {
	// Get customer
	_, err := s.customerRepo.GetByID(customerID)
	if err != nil {
		return nil, err
	}

	// Create order object
	order := &models.Order{
		CustomerID: customerID,
		Status:     models.OrderStatusPending,
	}

	var total float64
	var orderItems []models.OrderItem

	// Process each item
	for _, item := range req.Items {
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return nil, err
		}

		itemTotal := product.Price * float64(item.Quantity)
		total += itemTotal

		orderItems = append(orderItems, models.OrderItem{
			ProductID: product.ID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})
	}

	order.Total = total
	order.OrderItems = orderItems

	// Save order - after this, order won't have Customer loaded
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	// Reload the order with related data preloaded
	fullOrder, err := s.orderRepo.GetByID(ctx, order.ID)
	if err != nil {
		// fallback to original order if reload fails
		fullOrder = order
	}

	// Send notifications with fully loaded order
	if err := s.notifier.SendOrderConfirmation(fullOrder); err != nil {
		// Log error but do not fail the operation
	}

	return fullOrder, nil
}

func (s *orderService) GetOrder(ctx context.Context, customerID, orderID uint) (*models.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order.CustomerID != customerID {
		return nil, models.ErrOrderNotFound
	}

	return order, nil
}

func (s *orderService) GetOrders(ctx context.Context, customerID uint, page, limit int) ([]models.Order, int64, error) {
	return s.orderRepo.GetByCustomerID(ctx, customerID, page, limit)
}

func (s *orderService) UpdateOrderStatus(ctx context.Context, orderID uint, status models.OrderStatus) (*models.Order, error) {
	if err := s.orderRepo.UpdateStatus(ctx, orderID, status); err != nil {
		return nil, err
	}

	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// Send status update notification
	if err := s.notifier.SendStatusUpdate(order); err != nil {
		// Log error but don't fail the operation
	}

	return order, nil
}
