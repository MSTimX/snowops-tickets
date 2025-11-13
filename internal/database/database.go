package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"ticket-service/internal/config"
	"ticket-service/internal/models"
)

// DB хранит глобальное подключение к базе данных.
var DB *gorm.DB

// Init открывает соединение с PostgreSQL и сохраняет его в глобальной переменной DB.
func Init(cfg *config.Config) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSSLMode,
		cfg.DBTimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("ticket-service: failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&models.Ticket{}, &models.TicketAssignment{}); err != nil {
		log.Fatalf("ticket-service: failed to run migrations: %v", err)
	}

	DB = db

	log.Println("ticket-service: database connection established")
}

// GetDB возвращает активное подключение к базе данных.
func GetDB() *gorm.DB {
	if DB == nil {
		log.Fatal("ticket-service: database is not initialized")
	}
	return DB
}
