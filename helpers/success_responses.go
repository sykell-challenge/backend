package helpers

import (
	"github.com/gin-gonic/gin"
)

func SendSuccessResponse(c *gin.Context, data gin.H) {
	c.JSON(200, data)
}

func SendCreatedResponse(c *gin.Context, data gin.H) {
	c.JSON(201, data)
}
