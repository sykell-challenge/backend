package helpers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Common HTTP response helpers that can be used across all handlers

// SendBadRequestError sends a 400 error response
func SendBadRequestError(c *gin.Context, message string) {
	c.JSON(400, gin.H{"error": message})
}

// SendNotFoundError sends a 404 error response
func SendNotFoundError(c *gin.Context, message string) {
	c.JSON(404, gin.H{"error": message})
}

// SendConflictError sends a 409 error response
func SendConflictError(c *gin.Context, message string) {
	c.JSON(409, gin.H{"error": message})
}

// SendInternalError sends a 500 error response
func SendInternalError(c *gin.Context, message string) {
	c.JSON(500, gin.H{"error": message})
}

// SendSuccessResponse sends a 200 success response
func SendSuccessResponse(c *gin.Context, data gin.H) {
	c.JSON(200, data)
}

// SendCreatedResponse sends a 201 created response
func SendCreatedResponse(c *gin.Context, data gin.H) {
	c.JSON(201, data)
}

// ParseIDParam extracts and validates ID from URL parameter
func ParseIDParam(c *gin.Context, paramName string) (uint, bool) {
	idStr := c.Param(paramName)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		SendBadRequestError(c, "Invalid ID format")
		return 0, false
	}
	return uint(id), true
}

// ParseLimitQuery parses limit query parameter with default and max values
func ParseLimitQuery(c *gin.Context, defaultLimit, maxLimit int) int {
	limit := defaultLimit
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= maxLimit {
			limit = parsedLimit
		}
	}
	return limit
}

// ParsePaginationParams parses page and limit query parameters
func ParsePaginationParams(c *gin.Context) (page, limit int) {
	page = 1
	if pageStr := c.Query("page"); pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	limit = 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	return page, limit
}

// RequireQueryParam validates that a query parameter is present
func RequireQueryParam(c *gin.Context, paramName string) (string, bool) {
	value := c.Query(paramName)
	if value == "" {
		SendBadRequestError(c, paramName+" query parameter is required")
		return "", false
	}
	return value, true
}

// HandleDBError handles common database errors
func HandleDBError(c *gin.Context, err error, notFoundMessage string) bool {
	if err == gorm.ErrRecordNotFound {
		SendNotFoundError(c, notFoundMessage)
		return true
	}
	if err != nil {
		SendInternalError(c, err.Error())
		return true
	}
	return false
}

// ValidateJSONBinding validates JSON binding and sends error if invalid
func ValidateJSONBinding(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		SendBadRequestError(c, err.Error())
		return false
	}
	return true
}

// ValidateQueryBinding validates query parameter binding and sends error if invalid
func ValidateQueryBinding(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindQuery(obj); err != nil {
		SendBadRequestError(c, err.Error())
		return false
	}
	return true
}
