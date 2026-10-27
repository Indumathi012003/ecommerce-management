package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"order-management/model"
	"order-management/services/product/service"
)

type ProductController struct {
	service service.ProductService
}

func NewProductController(s service.ProductService) *ProductController {
	return &ProductController{service: s}
}

// requireAdmin ensures only ADMIN role can perform product management actions
func requireAdmin(ctx *gin.Context) bool {
	role := ctx.GetHeader("X-User-Role")
	// If call is coming from Gateway with authenticated user, check role
	if role != "" && role != "ADMIN" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden: product service is only accessible by admin"})
		return false
	}
	return true
}

// Add Product: POST /products
func (c *ProductController) AddProduct(ctx *gin.Context) {
	if !requireAdmin(ctx) {
		return
	}

	var req model.CreateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := c.service.AddProduct(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, product)
}

// View Single Product: GET /products/:id
func (c *ProductController) GetProduct(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	product, err := c.service.GetProduct(ctx.Request.Context(), id)
	if service.IsNotFound(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, product)
}

// View All Products: GET /products
func (c *ProductController) ListProducts(ctx *gin.Context) {
	if !requireAdmin(ctx) {
		return
	}

	products, err := c.service.ListProducts(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, products)
}

// Update Product: PUT /products/:id
func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	if !requireAdmin(ctx) {
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	var req model.UpdateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := c.service.UpdateProduct(ctx.Request.Context(), id, req)
	if service.IsNotFound(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, product)
}

// Delete Product: DELETE /products/:id
func (c *ProductController) DeleteProduct(ctx *gin.Context) {
	if !requireAdmin(ctx) {
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	err = c.service.DeleteProduct(ctx.Request.Context(), id)
	if service.IsNotFound(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "product deleted successfully"})
}
