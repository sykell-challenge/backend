package helpers

import (
	"github.com/gin-gonic/gin"
)

func RequireQueryParam(c *gin.Context, paramName string) (string, bool) {
	value := c.Query(paramName)
	if value == "" {
		SendBadRequestError(c, paramName+" query parameter is required")
		return "", false
	}
	return value, true
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
