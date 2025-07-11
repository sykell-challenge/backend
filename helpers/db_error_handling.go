package helpers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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
