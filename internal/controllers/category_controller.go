package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/Mutonya/Savanah/internal/domain/services"
	"github.com/Mutonya/Savanah/internal/utils/responses"
)

type CategoryController struct {
	categoryService services.CategoryService
}

func NewCategoryController(categoryService services.CategoryService) *CategoryController {
	return &CategoryController{categoryService: categoryService}
}

func (c *CategoryController) CreateCategory(ctx *gin.Context) {
	var req services.CategoryCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid category creation request")
		responses.ErrorResponse(ctx, http.StatusBadRequest, "invalid request payload")
		return
	}

	category, err := c.categoryService.CreateCategory(&req)
	if err != nil {
		log.Error().Err(err).Interface("request", req).Msg("Failed to create category")
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "failed to create category")
		return
	}

	log.Info().Uint("categoryID", category.ID).Msg("Category created successfully")
	responses.SuccessResponse(ctx, http.StatusCreated, category)
}

func (c *CategoryController) GetCategories(ctx *gin.Context) {
	categories, err := c.categoryService.GetCategories()
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch categories")
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "failed to fetch categories")
		return
	}

	log.Info().Int("count", len(categories)).Msg("Categories fetched successfully")
	responses.SuccessResponse(ctx, http.StatusOK, categories)
}

func (c *CategoryController) GetCategory(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Warn().Str("id", ctx.Param("id")).Msg("Invalid category ID format")
		responses.ErrorResponse(ctx, http.StatusBadRequest, "invalid category ID")
		return
	}

	category, err := c.categoryService.GetCategory(uint(id))
	if err != nil {
		log.Error().Err(err).Uint("categoryID", uint(id)).Msg("Failed to fetch category")
		responses.ErrorResponse(ctx, http.StatusNotFound, "category not found")
		return
	}

	log.Info().Uint("categoryID", uint(id)).Msg("Category fetched successfully")
	responses.SuccessResponse(ctx, http.StatusOK, category)
}

func (c *CategoryController) GetCategoryProducts(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Warn().Str("id", ctx.Param("id")).Msg("Invalid category ID format")
		responses.ErrorResponse(ctx, http.StatusBadRequest, "invalid category ID")
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	products, total, err := c.categoryService.GetCategoryProducts(uint(id), page, limit)
	if err != nil {
		log.Error().Err(err).Uint("categoryID", uint(id)).Msg("Failed to fetch category products")
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "failed to fetch category products")
		return
	}

	log.Info().Uint("categoryID", uint(id)).Int("count", len(products)).Msg("Category products fetched successfully")
	responses.PaginatedResponse(ctx, http.StatusOK, products, int64(total), page, limit)
}

func (c *CategoryController) GetAveragePrice(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Warn().Str("id", ctx.Param("id")).Msg("Invalid category ID format")
		responses.ErrorResponse(ctx, http.StatusBadRequest, "invalid category ID")
		return
	}

	avgPrice, err := c.categoryService.GetAveragePrice(uint(id))
	if err != nil {
		log.Error().Err(err).Uint("categoryID", uint(id)).Msg("Failed to calculate average price")
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "failed to calculate average price")
		return
	}

	log.Info().Uint("categoryID", uint(id)).Float64("averagePrice", avgPrice).Msg("Average price calculated successfully")
	responses.SuccessResponse(ctx, http.StatusOK, gin.H{"average_price": avgPrice})
}

func (c *CategoryController) UpdateCategory(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Warn().Str("id", ctx.Param("id")).Msg("Invalid category ID format")
		responses.ErrorResponse(ctx, http.StatusBadRequest, "invalid category ID")
		return
	}

	var req services.CategoryUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid category update request")
		responses.ErrorResponse(ctx, http.StatusBadRequest, "invalid request payload")
		return
	}

	category, err := c.categoryService.UpdateCategory(uint(id), &req)
	if err != nil {
		log.Error().Err(err).Uint("categoryID", uint(id)).Msg("Failed to update category")
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "failed to update category")
		return
	}

	log.Info().Uint("categoryID", uint(id)).Msg("Category updated successfully")
	responses.SuccessResponse(ctx, http.StatusOK, category)
}

func (c *CategoryController) DeleteCategory(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Warn().Str("id", ctx.Param("id")).Msg("Invalid category ID format")
		responses.ErrorResponse(ctx, http.StatusBadRequest, "invalid category ID")
		return
	}

	if err := c.categoryService.DeleteCategory(uint(id)); err != nil {
		log.Error().Err(err).Uint("categoryID", uint(id)).Msg("Failed to delete category")
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "failed to delete category")
		return
	}

	log.Info().Uint("categoryID", uint(id)).Msg("Category deleted successfully")
	ctx.Status(http.StatusNoContent)
}
