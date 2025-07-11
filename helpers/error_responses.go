package helpers

import (
	"github.com/gin-gonic/gin"
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
