package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"order-management/model"
	"order-management/services/order/client"
	"order-management/services/order/service"
)

type OrderController struct {
	service service.OrderService
}

func NewOrderController(s service.OrderService) *OrderController {
	return &OrderController{service: s}
}

func (c *OrderController) Create(ctx *gin.Context) {
	var req model.CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := c.service.CreateOrder(ctx.Request.Context(), req)
	if err != nil {
		if errors.Is(err, client.ErrCircuitOpen) || errors.Is(err, client.ErrServiceUnavailable) {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{
				"error":   "Product service is temporarily unavailable",
				"details": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, order)
}

func (c *OrderController) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	order, err := c.service.GetOrder(ctx.Request.Context(), id)
	if service.IsNotFound(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	role := ctx.GetHeader("X-User-Role")
	customerIDStr := ctx.GetHeader("X-Customer-ID")

	// ADMIN can access any order
	if role == "ADMIN" {
		ctx.JSON(http.StatusOK, order)
		return
	}

	// CUSTOMER can only access their own order
	customerID, _ := strconv.Atoi(customerIDStr)
	if order.CustomerID != customerID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden: you cannot access another customer's order"})
		return
	}

	ctx.JSON(http.StatusOK, order)
}

func (c *OrderController) GetAll(ctx *gin.Context) {
	role := ctx.GetHeader("X-User-Role")
	customerIDStr := ctx.GetHeader("X-Customer-ID")

	orders, err := c.service.ListOrders(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ADMIN can access all orders
	if role == "ADMIN" {
		ctx.JSON(http.StatusOK, orders)
		return
	}

	// CUSTOMER: Filter out other customers' orders
	customerID, _ := strconv.Atoi(customerIDStr)
	var customerOrders []*model.Order
	for _, o := range orders {
		if o.CustomerID == customerID {
			customerOrders = append(customerOrders, o)
		}
	}
	ctx.JSON(http.StatusOK, customerOrders)
}

func (c *OrderController) UpdateStatus(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	var req model.UpdateStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := c.service.UpdateOrderStatus(ctx.Request.Context(), id, req.Status)
	if service.IsNotFound(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, order)
}
