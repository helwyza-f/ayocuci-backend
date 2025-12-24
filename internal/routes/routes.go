package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/helwyza-f/ayocuci-backend/internal/middleware"
	"github.com/helwyza-f/ayocuci-backend/internal/module/audit_log"
	"github.com/helwyza-f/ayocuci-backend/internal/module/auth"
	"github.com/helwyza-f/ayocuci-backend/internal/module/client"
	"github.com/helwyza-f/ayocuci-backend/internal/module/customer"
	"github.com/helwyza-f/ayocuci-backend/internal/module/employee"
	"github.com/helwyza-f/ayocuci-backend/internal/module/expense"
	"github.com/helwyza-f/ayocuci-backend/internal/module/laundry_service"
	"github.com/helwyza-f/ayocuci-backend/internal/module/order"
	"github.com/helwyza-f/ayocuci-backend/internal/module/outlet"
	"github.com/helwyza-f/ayocuci-backend/internal/module/report"
)

func Register(r *gin.Engine, db *gorm.DB) {
	// ===== AUTH =====
	userRepo := auth.NewRepository(db)
	clientRepo := client.NewRepository(db)

	authService := auth.NewService(
		db,
		userRepo,
		clientRepo,
	)
	authHandler := auth.NewHandler(authService)

	// ===== AUDIT LOG =====
	auditRepo := audit_log.NewRepository(db)
	auditService := audit_log.NewService(auditRepo)

	// ===== OUTLET =====
	outletRepo := outlet.NewRepository(db)
	outletService := outlet.NewService(outletRepo)
	outletHandler := outlet.NewHandler(outletService)

	// ===== LAUNDRY SERVICE =====
	serviceRepo := laundry_service.NewRepository(db)
	serviceLogic := laundry_service.NewService(serviceRepo)
	serviceHandler := laundry_service.NewHandler(serviceLogic)

	// ===== CUSTOMER =====
	customerRepo := customer.NewRepository(db)
	customerService := customer.NewService(customerRepo)
	customerHandler := customer.NewHandler(customerService)

	// ===== EMPLOYEE =====
	employeeRepo := employee.NewRepository(db)
	employeeService := employee.NewService(
		employeeRepo,
		userRepo,
		outletRepo,
	)
	employeeHandler := employee.NewHandler(employeeService)

	// ===== ORDER =====
	orderRepo := order.NewRepository(db)
	orderService := order.NewService(orderRepo, auditService)
	orderHandler := order.NewHandler(orderService)

	// ===== EXPENSE =====
	expenseRepo := expense.NewRepository(db)
	expenseService := expense.NewService(expenseRepo)
	expenseHandler := expense.NewHandler(expenseService)

	// ===== REPORT =====
	reportService := report.NewService(db)
	reportHandler := report.NewHandler(reportService)

	// ===== PUBLIC =====
	r.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)

	// =============================
	// AUTHENTICATED ROUTES
	// =============================
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())

	// ===== OUTLET (OWNER) =====
	api.POST("/outlets", outletHandler.Create)
	api.GET("/outlets", outletHandler.List)

	// ===== SELECT OUTLET (ACK ONLY) =====
	// frontend pakai ini setelah user pilih outlet
	api.POST("/select-outlet",
		middleware.OutletContext(db),
		func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message":   "outlet selected",
				"outlet_id": c.GetUint("outlet_id"),
			})
		},
	)

	// ===== EMPLOYEE =====
	
	api.GET("/me/outlets", 
	middleware.RoleMiddleware("owner", "staff"),
	employeeHandler.MyOutlets)
	// Di router
	api.PUT("/employees/:user_id/transfer", middleware.RoleMiddleware("owner"), employeeHandler.Transfer)
	api.POST("/employees",middleware.RoleMiddleware("owner"), employeeHandler.CreateUser)
	
	api.POST("/outlets/:outlet_id/employees",middleware.RoleMiddleware("owner"), employeeHandler.Assign)
	api.DELETE("/outlets/:id",middleware.RoleMiddleware("owner"), outletHandler.Delete)
	

	// =============================
	// ROUTES WITH OUTLET CONTEXT
	// =============================
	apiOutlet := api.Group("/")
	apiOutlet.Use(middleware.OutletContext(db))

	// ===== EMPLOYEE =====
	apiOutlet.GET("/employees", employeeHandler.ListEmployees)

	// ===== LAUNDRY SERVICES (OWNER) =====
	apiOutlet.POST("/services",
		middleware.RoleMiddleware("owner", "staff"),
		serviceHandler.Create,
	)
	apiOutlet.GET("/services",
		middleware.RoleMiddleware("owner", "staff"),
		serviceHandler.List,
	)
	apiOutlet.PUT("/services/:id",
		middleware.RoleMiddleware("owner", "staff"),
		serviceHandler.Update,
	)
	apiOutlet.DELETE("/services/:id",
		middleware.RoleMiddleware("owner"),
		serviceHandler.Delet,
	)

	// ===== CUSTOMERS (OWNER + KASIR) =====
	apiOutlet.POST("/customers",
		middleware.RoleMiddleware("owner", "kasir", "staff"),
		customerHandler.Create,
	)
	apiOutlet.GET("/customers",
		middleware.RoleMiddleware("owner", "kasir", "staff"),
		customerHandler.List,
	)
	apiOutlet.GET("/customers/:id",
		middleware.RoleMiddleware("owner", "kasir", "staff"),
		customerHandler.GetByID,
	)
	apiOutlet.PUT("/customers/:id",
		middleware.RoleMiddleware("owner", "kasir", "staff"),
		customerHandler.Update,
	)

	// ===== ORDERS =====
	apiOutlet.POST("/orders",
		middleware.RoleMiddleware("owner", "kasir", "staff"),
		orderHandler.Create,
	)
	apiOutlet.GET("/orders",
		middleware.RoleMiddleware("owner", "kasir", "kurir", "staff"),
		orderHandler.List,
	)
	apiOutlet.PUT("/orders/:id/status",
		middleware.RoleMiddleware("owner", "kasir", "kurir", "staff"),
		orderHandler.UpdateStatus,
	)

	// ===== EXPENSES (OWNER) =====
	apiOutlet.POST("/expenses",
		middleware.RoleMiddleware("owner", "staff"),
		expenseHandler.Create,
	)
	apiOutlet.GET("/expenses",
		middleware.RoleMiddleware("owner", "staff"),
		expenseHandler.List,
	)
	apiOutlet.PUT("/expenses/:id",
		middleware.RoleMiddleware("owner", "staff"),
		expenseHandler.Update,
	)
	apiOutlet.DELETE("/expenses/:id",
		middleware.RoleMiddleware("owner", "staff"),
		expenseHandler.Delete,
	)

	// ===== REPORT (OWNER) =====
	apiOutlet.GET("/reports/daily",
		middleware.RoleMiddleware("owner","staff"),
		reportHandler.Daily,
	)
	apiOutlet.GET("/reports/range",
		middleware.RoleMiddleware("owner","staff"),
		reportHandler.Range,
	)
}
