package responses

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuccessResponse sends a standardized success response
func SuccessResponse(ctx *gin.Context, statusCode int, data interface{}) {
	ctx.JSON(statusCode, gin.H{
		"success": true,
		"data":    data,
	})
}

// ErrorResponse sends a standardized error response
func ErrorResponse(ctx *gin.Context, statusCode int, message string) {
	ctx.JSON(statusCode, gin.H{
		"success": false,
		"error": gin.H{
			"code":    statusCode,
			"message": message,
		},
	})
}

// PaginatedResponse sends a standardized paginated response
func PaginatedResponse(ctx *gin.Context, statusCode int, data interface{}, total int64, page, limit int) {
	ctx.JSON(statusCode, gin.H{
		"success": true,
		"data":    data,
		"meta": gin.H{
			"total":     total,
			"page":      page,
			"limit":     limit,
			"totalPage": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// JSONResponse is a generic JSON response writer
func JSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// EmptyResponse writes an empty response with status code
func EmptyResponse(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}
