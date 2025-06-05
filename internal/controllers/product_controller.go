package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/Mutonya/Savanah/internal/domain/services"
	"github.com/Mutonya/Savanah/internal/utils/responses"
)

type ProductController struct {
	productService services.ProductService
}

func NewProductController(productService services.ProductService) *ProductController {
	return &ProductController{productService: productService}
}

// @Summary Create a new product
// @Description Create a new product with categories
// @Tags products
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param product body services.ProductCreateRequest true "Product data"
// @Success 201 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/products [post]
func (c *ProductController) CreateProduct(ctx *gin.Context) {
	var req services.ProductCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid product creation request")
		responses.ErrorResponse(ctx, http.StatusBadRequest, "invalid request payload")
		return
	}

	product, err := c.productService.CreateProduct(ctx, &req)

	if err != nil {
		log.Error().Err(err).Interface("request", req).Msg("Failed to create product")
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "failed to create product")
		return
	}

	log.Info().Uint("productID", product.ID).Msg("Product created successfully")
	responses.SuccessResponse(ctx, http.StatusCreated, product)
}

// @Summary Get all products
// @Description Get a list of all products with pagination
// @Tags products
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} responses.PaginatedResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/products [get]
func (c *ProductController) GetProducts(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	products, total, err := c.productService.GetProducts(ctx, page, limit)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch products")
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "failed to fetch products")
		return
	}

	log.Info().Int("count", len(products)).Msg("Products fetched successfully")
	responses.PaginatedResponse(ctx, http.StatusOK, products, total, page, limit)
}

// @Summary Get a single product
// @Description Get details of a specific product
// @Tags products
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "Product ID"
// @Success 200 {object} responses.SuccessResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/products/{id} [get]
func (c *ProductController) GetProduct(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Warn().Str("id", ctx.Param("id")).Msg("Invalid product ID format")
		responses.ErrorResponse(ctx, http.StatusBadRequest, "invalid product ID")
		return
	}

	product, err := c.productService.GetProduct(ctx, uint(id))
	if err != nil {
		log.Error().Err(err).Uint("productID", uint(id)).Msg("Failed to fetch product")
		responses.ErrorResponse(ctx, http.StatusNotFound, "product not found")
		return
	}

	log.Info().Uint("productID", uint(id)).Msg("Product fetched successfully")
	responses.SuccessResponse(ctx, http.StatusOK, product)
}

// @Summary Update a product
// @Description Update an existing product's details
// @Tags products
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "Product ID"
// @Param product body services.ProductUpdateRequest true "Product data"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/products/{id} [put]
func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Warn().Str("id", ctx.Param("id")).Msg("Invalid product ID format")
		responses.ErrorResponse(ctx, http.StatusBadRequest, "invalid product ID")
		return
	}

	var req services.ProductUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid product update request")
		responses.ErrorResponse(ctx, http.StatusBadRequest, "invalid request payload")
		return
	}

	product, err := c.productService.UpdateProduct(ctx, uint(id), &req)
	if err != nil {
		log.Error().Err(err).Uint("productID", uint(id)).Msg("Failed to update product")
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "failed to update product")
		return
	}

	log.Info().Uint("productID", uint(id)).Msg("Product updated successfully")
	responses.SuccessResponse(ctx, http.StatusOK, product)
}

// @Summary Delete a product
// @Description Delete a product from the system
// @Tags products
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "Product ID"
// @Success 204
// @Failure 400 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/products/{id} [delete]
func (c *ProductController) DeleteProduct(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Warn().Str("id", ctx.Param("id")).Msg("Invalid product ID format")
		responses.ErrorResponse(ctx, http.StatusBadRequest, "invalid product ID")
		return
	}

	if err := c.productService.DeleteProduct(ctx, uint(id)); err != nil {
		log.Error().Err(err).Uint("productID", uint(id)).Msg("Failed to delete product")
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "failed to delete product")
		return
	}

	log.Info().Uint("productID", uint(id)).Msg("Product deleted successfully")
	ctx.Status(http.StatusNoContent)
}
