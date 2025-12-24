package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/helwyza-f/ayocuci-backend/internal/config"
	"github.com/helwyza-f/ayocuci-backend/internal/database"
	"github.com/helwyza-f/ayocuci-backend/internal/module/audit_log"
	"github.com/helwyza-f/ayocuci-backend/internal/module/auth"
	"github.com/helwyza-f/ayocuci-backend/internal/module/client"
	"github.com/helwyza-f/ayocuci-backend/internal/module/customer"
	"github.com/helwyza-f/ayocuci-backend/internal/module/employee"
	"github.com/helwyza-f/ayocuci-backend/internal/module/expense"
	"github.com/helwyza-f/ayocuci-backend/internal/module/laundry_service"
	"github.com/helwyza-f/ayocuci-backend/internal/module/order"
	"github.com/helwyza-f/ayocuci-backend/internal/module/outlet"
	"github.com/helwyza-f/ayocuci-backend/internal/routes"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using default env")
	}

	dbPath := config.GetEnv("DB_URL", "ayocuci.db")
	database.Connect(dbPath)
	database.DB.AutoMigrate(&auth.User{}, &client.Client{}, &outlet.Outlet{}, &laundry_service.Service{}, &customer.Customer{}, &employee.UserOutlet{}, &order.Order{},
	&order.OrderItem{}, &expense.Expense{}, &audit_log.AuditLog{})

	r := gin.Default()
	r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"}, // Izinkan semua (buat dev/test oke banget)
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))
	routes.Register(r, database.DB)

	port := config.GetEnv("APP_PORT", "8080")
	r.Run(":" + port)
}
