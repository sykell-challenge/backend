package main

import (
	"os"
	"strings"
	"sykell-challenge/backend/auth"
	"sykell-challenge/backend/db"
	crawlHandler "sykell-challenge/backend/handlers/crawl"
	"sykell-challenge/backend/handlers/url"
	"sykell-challenge/backend/handlers/user"
	"sykell-challenge/backend/services/socket"
	"sykell-challenge/backend/utils/crawl"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	database := db.GetDB()
	db.MigrateAll()

	// Initialize crawl manager with database connection
	crawl.InitializeCrawlManager(database)

	// Initialize handlers
	urlHandler := url.NewURLHandler(database)
	userHandler := user.NewUserHandler(database)

	router := gin.Default()

	// Configure CORS with environment-based origins
	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}

	// Set allowed origins based on environment
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		panic("ALLOWED_ORIGINS environment variable is required!")
	} else {
		corsConfig.AllowOrigins = strings.Split(allowedOrigins, ",")
	}

	router.Use(cors.New(corsConfig))

	// Public user routes (no authentication required)
	router.POST("/users", userHandler.CreateUser)
	router.POST("/users/login", userHandler.LoginUser)

	// Protected routes (require JWT authentication)
	protected := router.Group("/")
	protected.Use(auth.JWTMiddleware())

	// URL routes (protected)
	protected.GET("/urls", urlHandler.GetURLs)
	protected.GET("/urls/search", urlHandler.SearchURLByString)
	protected.GET("/urls/search/fuzzy", urlHandler.FuzzySearchURLs)
	protected.GET("/urls/stats", urlHandler.GetURLStats)
	protected.GET("/urls/:id", urlHandler.GetURLByID)
	protected.GET("/urls/:id/links", urlHandler.GetURLLinks)
	protected.GET("/urls/:id/links/internal", urlHandler.GetURLInternalLinks)
	protected.GET("/urls/:id/links/external", urlHandler.GetURLExternalLinks)
	protected.GET("/urls/:id/links/broken", urlHandler.GetURLBrokenLinks)
	protected.POST("/urls", urlHandler.CreateURL)
	protected.PUT("/urls/:id", urlHandler.UpdateURL)
	protected.PATCH("/urls/:id/status", urlHandler.UpdateURLStatus)
	protected.DELETE("/urls/:id", urlHandler.DeleteURL)

	// User routes (protected)
	protected.GET("/users/:id", userHandler.GetUserByID)
	protected.PUT("/users/:id", userHandler.UpdateUser)
	protected.DELETE("/users/:id", userHandler.DeleteUser)

	// Crawl routes (protected)
	protected.POST("/crawl", crawlHandler.HandleCrawlURL)
	protected.GET("/crawl/:jobId", crawlHandler.HandleGetCrawlStatus)
	protected.GET("/crawl/:jobId/url", crawlHandler.HandleGetURLByJobID)
	protected.DELETE("/crawl/:jobId", crawlHandler.HandleCancelCrawl)
	protected.GET("/crawl", crawlHandler.HandleListActiveJobs)
	protected.GET("/crawl-history", crawlHandler.HandleGetJobHistory)
	protected.GET("/crawl-stats", crawlHandler.HandleGetJobStats)
	protected.GET("/crawl/by-url/:urlId", crawlHandler.HandleGetJobsByURL)

	server := socket.InitSocketServer()

	router.Any("/socket.io/*any", gin.WrapH(server.ServeHandler(nil)))

	router.Run("0.0.0.0:8080")
}
