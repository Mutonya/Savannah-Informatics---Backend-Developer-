package controllers

import (
	"github.com/Mutonya/Savanah/internal/domain/services"
	"github.com/Mutonya/Savanah/internal/utils/responses"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

type OrderController struct {
	orderService        services.OrderService
	notificationService services.NotificationService
}

func NewOrderController(orderService services.OrderService, notificationService services.NotificationService) *OrderController {
	return &OrderController{
		orderService:        orderService,
		notificationService: notificationService,
	}
}

// @Summary Create a new order
// @Description Create a new order for products
// @Tags orders
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param order body services.OrderCreateRequest true "Order data"
// @Success 201 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/orders [post]
func (c *OrderController) CreateOrder(ctx *gin.Context) {
	customerID, _ := ctx.Get("customerID")

	var req services.OrderCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid order creation request")
		responses.ErrorResponse(ctx, http.StatusBadRequest, "invalid request payload")
		return
	}

	order, err := c.orderService.CreateOrder(ctx, customerID.(uint), &req)
	if err != nil {
		log.Error().Err(err).Uint("customerID", customerID.(uint)).Msg("Failed to create order")
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "failed to create order")
		return
	}

	// Send notifications
	if err := c.notificationService.SendOrderConfirmation(order); err != nil {
		log.Error().Err(err).Uint("orderID", order.ID).Msg("Failed to send order notifications")
		// Continue despite notification failure
	}

	log.Info().Uint("orderID", order.ID).Uint("customerID", customerID.(uint)).Msg("Order created successfully")
	responses.SuccessResponse(ctx, http.StatusCreated, order)
}

// @Summary Get all orders
// @Description Get a list of all orders with pagination
// @Tags orders
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} responses.PaginatedResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/orders [get]
func (c *OrderController) GetOrders(ctx *gin.Context) {
	customerID, _ := ctx.Get("customerID")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	orders, total, err := c.orderService.GetOrders(ctx, customerID.(uint), page, limit)
	if err != nil {
		log.Error().Err(err).Uint("customerID", customerID.(uint)).Msg("Failed to fetch orders")
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "failed to fetch orders")
		return
	}

	log.Info().Uint("customerID", customerID.(uint)).Int("count", len(orders)).Msg("Orders fetched successfully")
	responses.PaginatedResponse(ctx, http.StatusOK, orders, int64(total), page, limit)
}

// @Summary Get order details
// @Description Get details of a specific order
// @Tags orders
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "Order ID"
// @Success 200 {object} responses.SuccessResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/orders/{id} [get]
func (c *OrderController) GetOrder(ctx *gin.Context) {
	customerID, _ := ctx.Get("customerID")
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Warn().Str("id", ctx.Param("id")).Msg("Invalid order ID format")
		responses.ErrorResponse(ctx, http.StatusBadRequest, "invalid order ID")
		return
	}

	order, err := c.orderService.GetOrder(ctx, customerID.(uint), uint(id))
	if err != nil {
		log.Error().Err(err).Uint("orderID", uint(id)).Msg("Failed to fetch order")
		responses.ErrorResponse(ctx, http.StatusNotFound, "order not found")
		return
	}

	log.Info().Uint("orderID", uint(id)).Msg("Order fetched successfully")
	responses.SuccessResponse(ctx, http.StatusOK, order)
}

// @Summary Update order status
// @Description Update the status of an existing order
// @Tags orders
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "Order ID"
// @Param status body services.OrderStatusUpdateRequest true "Status data"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/orders/{id}/status [put]
func (c *OrderController) UpdateOrderStatus(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Warn().Str("id", ctx.Param("id")).Msg("Invalid order ID format")
		responses.ErrorResponse(ctx, http.StatusBadRequest, "invalid order ID")
		return
	}

	var req services.OrderStatusUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid status update request")
		responses.ErrorResponse(ctx, http.StatusBadRequest, "invalid request payload")
		return
	}

	order, err := c.orderService.UpdateOrderStatus(ctx, uint(id), req.Status)
	if err != nil {
		log.Error().Err(err).Uint("orderID", uint(id)).Msg("Failed to update order status")
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "failed to update order status")
		return
	}

	// Send status update notification
	if err := c.notificationService.SendStatusUpdate(order); err != nil {
		log.Error().Err(err).Uint("orderID", order.ID).Msg("Failed to send status notification")
		// Continue despite notification failure
	}

	log.Info().Uint("orderID", uint(id)).Str("status", string(req.Status)).Msg("Order status updated successfully")
	responses.SuccessResponse(ctx, http.StatusOK, order)
}
