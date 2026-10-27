package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"order-management/config"
	"order-management/database"
	"order-management/middleware"
	"order-management/services/order/controller"
	gqlserver "order-management/services/order/graphql"
	"order-management/services/order/repository"
	"order-management/services/order/service"
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

	repo := repository.NewOrderRepository(db)
	svc := service.NewOrderService(repo)
	ctrl := controller.NewOrderController(svc)
	gs, err := gqlserver.NewServer(svc)
	if err != nil {
		log.Fatal("GraphQL schema failed:", err)
	}

	r := gin.New()
	r.Use(
		middleware.RequestID(),
		middleware.Recovery(),
		middleware.Logging(),
	)

	r.GET("/health/live", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"service": "order-service", "status": "ALIVE"})
	})

	r.GET("/health/ready", func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "DOWN", "error": err.Error()})
			return
		}
	
		if err := sqlDB.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "DOWN", "error": err.Error()})
			return
		}
	
		c.JSON(http.StatusOK, gin.H{"status": "READY"})
	})

	v1 := r.Group("/api/v1")
	{
		v1.POST("/orders", ctrl.Create)
		v1.GET("/orders", ctrl.GetAll)
		v1.GET("/orders/:id", ctrl.GetByID)
		v1.PUT("/orders/:id/status", ctrl.UpdateStatus)
	}

	r.POST("/orders", ctrl.Create)
	r.GET("/orders", ctrl.GetAll)
	r.GET("/orders/:id", ctrl.GetByID)
	r.PUT("/orders/:id/status", ctrl.UpdateStatus)

	h := gqlserver.NewHTTPHandler(gs)
	r.Any("/graphql", gin.WrapH(h))
	r.GET("/playground", gin.WrapH(h))

	log.Println("Order Service: http://localhost:" + cfg.OrderServicePort + "/api/v1/orders")
	log.Fatal(r.Run(":" + cfg.OrderServicePort))
}
