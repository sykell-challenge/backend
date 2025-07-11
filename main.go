package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"sykell-challenge/backend/auth"
	"sykell-challenge/backend/db"
	"sykell-challenge/backend/handlers/crawl"
	"sykell-challenge/backend/handlers/url"
	"sykell-challenge/backend/handlers/user"
	"sykell-challenge/backend/services/socket"
	"sykell-challenge/backend/services/taskq"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	database := db.GetDB()
	db.MigrateAll()

	// Initialize task queue for background crawling
	taskq.InitTaskQueue()

	// Initialize handlers
	urlHandler := url.NewURLHandler(database)
	userHandler := user.NewUserHandler(database)
	crawlHandler := crawl.NewCrawlHandler(database)

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
	protected.DELETE("/crawl/:jobId", crawlHandler.HandleCancelCrawl)
	protected.GET("/crawl-history", crawlHandler.HandleGetAllCrawlJobs)

	server := socket.InitSocketServer()

	router.Any("/socket.io/*any", gin.WrapH(server.ServeHandler(nil)))

	// Create HTTP server
	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Println("Server starting on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Shutdown task queue first
	taskq.ShutdownTaskQueue()

	// Shutdown HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
