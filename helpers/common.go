package helpers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SendBadRequestError(c *gin.Context, message string) {
	c.JSON(400, gin.H{"error": message})
}

func SendNotFoundError(c *gin.Context, message string) {
	c.JSON(404, gin.H{"error": message})
}

func SendConflictError(c *gin.Context, message string) {
	c.JSON(409, gin.H{"error": message})
}

func SendInternalError(c *gin.Context, message string) {
	c.JSON(500, gin.H{"error": message})
}

func SendSuccessResponse(c *gin.Context, data gin.H) {
	c.JSON(200, data)
}

func SendCreatedResponse(c *gin.Context, data gin.H) {
	c.JSON(201, data)
}

func ParseIDParam(c *gin.Context, paramName string) (uint, bool) {
	idStr := c.Param(paramName)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		SendBadRequestError(c, "Invalid ID format")
		return 0, false
	}
	return uint(id), true
}

func ParseLimitQuery(c *gin.Context, defaultLimit, maxLimit int) int {
	limit := defaultLimit
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= maxLimit {
			limit = parsedLimit
		}
	}
	return limit
}

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

func RequireQueryParam(c *gin.Context, paramName string) (string, bool) {
	value := c.Query(paramName)
	if value == "" {
		SendBadRequestError(c, paramName+" query parameter is required")
		return "", false
	}
	return value, true
}

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

func ValidateJSONBinding(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		SendBadRequestError(c, err.Error())
		return false
	}
	return true
}

func ValidateQueryBinding(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindQuery(obj); err != nil {
		SendBadRequestError(c, err.Error())
		return false
	}
	return true
}
