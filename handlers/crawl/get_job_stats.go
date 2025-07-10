package crawl

import (
	"sykell-challenge/backend/utils/crawl"

	"github.com/gin-gonic/gin"
)

// HandleGetJobStats returns crawl job statistics
func HandleGetJobStats(g *gin.Context) {
	stats, err := crawl.GetJobStats()
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	g.JSON(200, gin.H{"stats": stats})
}
