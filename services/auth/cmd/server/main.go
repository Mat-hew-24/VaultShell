package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"syncpad/services/auth/internal/handler"
	"syncpad/services/auth/internal/repository"
	"syncpad/services/auth/internal/service"
)

func main() {
	db, err := sql.Open("postgres", "postgres://syncuser:syncpass@localhost:5432/syncpad?sslmode=disable") //PGSQL CONNECTION STRING
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Database is not reachable:", err) //ERROR HAPS CODE GET STUCK IG
	}

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	// THIS IS THE GIN ROUTER
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to SyncPad Auth Service!")
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "auth",
			"status":  "running",
		})
	})

	r.GET("/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "This route will return a list of users (not implemented yet).",
		})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "auth"})
	})

	port := ":8081"
	log.Printf("Auth service starting on port %s", port)
	log.Fatal(r.Run(port))
}
