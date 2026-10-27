package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"order-management/config"
	"order-management/database"
	"order-management/middleware"
	"order-management/services/product/controller"
	"order-management/services/product/repository"
	"order-management/services/product/service"
)

func main() {
	cfg := config.Load()
	db, err := database.InitDB(&cfg)
	if err != nil {
		log.Fatal("database connection failed:", err)
	}
	
	sqlDB, err := db.DB()
	if err == nil {
		defer sqlDB.Close()
	}

	// Dependency Injection: Repository -> Service -> Controller
	repo := repository.NewProductRepository(db)
	svc := service.NewProductService(repo)
	ctrl := controller.NewProductController(svc)

	r := gin.New()
	r.Use(middleware.RequestID(), middleware.Recovery(), middleware.Logging())

	// Health check
	r.GET("/health/live", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"service": "product-service", "status": "ALIVE"})
	})

	// Product CRUD Endpoints (Accessible by ADMIN)
	r.POST("/products", ctrl.AddProduct)
	r.GET("/products", ctrl.ListProducts)
	r.GET("/products/:id", ctrl.GetProduct)
	r.PUT("/products/:id", ctrl.UpdateProduct)
	r.DELETE("/products/:id", ctrl.DeleteProduct)

	// Also support /api/v1/products prefix
	v1 := r.Group("/api/v1")
	{
		v1.POST("/products", ctrl.AddProduct)
		v1.GET("/products", ctrl.ListProducts)
		v1.GET("/products/:id", ctrl.GetProduct)
		v1.PUT("/products/:id", ctrl.UpdateProduct)
		v1.DELETE("/products/:id", ctrl.DeleteProduct)
	}

	log.Printf("Product Service running on port :%s\n", cfg.ProductServicePort)
	if err := r.Run(":" + cfg.ProductServicePort); err != nil {
		log.Fatalf("Product Service failed: %v", err)
	}
}
