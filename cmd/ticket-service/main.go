package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"ticket-service/internal/config"
	"ticket-service/internal/database"
	"ticket-service/internal/handlers"
)

func main() {
	cfg := config.Load()

	// Инициализируем подключение к базе данных
	database.Init(cfg)

	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)

	router := gin.Default()

	handlers.RegisterTicketRoutes(router)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "ticket-service",
		})
	})

	addr := ":" + cfg.HTTPPort
	log.Printf("ticket-service: starting HTTP server on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("ticket-service: failed to start HTTP server: %v", err)
	}
}
