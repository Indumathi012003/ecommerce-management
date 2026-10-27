package main

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var secret = []byte("lab2-secret-key-change-me")

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func jwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Bearer token"})
			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrSignatureInvalid
			}
			return secret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		// Forward claims to downstream services
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			role, _ := claims["role"].(string)
			c.Request.Header.Set("X-User-Role", role)
			c.Set("role", role)

			if cid, ok := claims["customer_id"].(float64); ok && cid > 0 {
				c.Request.Header.Set("X-Customer-ID", strconv.Itoa(int(cid)))
			}
		}

		c.Next()
	}
}

// adminOnlyMiddleware ensures ONLY ADMIN can access protected endpoints like /products
func adminOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != "ADMIN" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "forbidden: product service is only accessible by admin",
			})
			return
		}
		c.Next()
	}
}

func proxyTo(targetURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := http.NewRequestWithContext(
			c.Request.Context(),
			c.Request.Method,
			targetURL+c.Request.URL.RequestURI(),
			c.Request.Body,
		)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		// Copy incoming headers
		for k, v := range c.Request.Header {
			req.Header[k] = v
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "downstream service unavailable"})
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
	}
}

func main() {
	_ = godotenv.Load()

	orderServiceURL := os.Getenv("ORDER_SERVICE_URL")
	if orderServiceURL == "" {
		orderServiceURL = "http://localhost:8081"
	}

	productServiceURL := os.Getenv("PRODUCT_SERVICE_URL")
	if productServiceURL == "" {
		productServiceURL = "http://localhost:8082"
	}

	adminPassword := os.Getenv("ADMIN_PASSWORD")

	customerPassword := os.Getenv("CUSTOMER_PASSWORD")

	r := gin.Default()

	// Public Gateway Health: Port/Process check
	r.GET("/health", func(c *gin.Context) {
	    c.JSON(http.StatusOK, gin.H{"service": "gateway", "status": "UP"})
	})

	// Database / Downstream Readiness: Proxied to Order Service (:8081)
	r.GET("/health/ready", proxyTo(orderServiceURL))

	// Login Endpoint using env passwords
	r.POST("/login", func(c *gin.Context) {
		var req loginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}

		var role string
		var customerID int

		if req.Username == "admin" && req.Password == adminPassword {
			role = "ADMIN"
			customerID = 0
		} else if strings.HasPrefix(req.Username, "customer") && req.Password == customerPassword {
			role = "CUSTOMER"
			idStr := strings.TrimPrefix(strings.TrimPrefix(req.Username, "customer-"), "customer")
			customerID, _ = strconv.Atoi(idStr)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":         req.Username,
			"role":        role,
			"customer_id": customerID,
			"exp":         time.Now().Add(15 * time.Minute).Unix(),
		})

		tokenStr, err := token.SignedString(secret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": tokenStr, "role": role, "customer_id": customerID})
	})

	// Protected Routes (Proxied to downstream services)
	ordersProxy := proxyTo(orderServiceURL)
	productsProxy := proxyTo(productServiceURL)

	// Order Service endpoints (accessible by logged-in users)
	orders := r.Group("/api/v1/orders", jwtMiddleware())
	orders.Any("", ordersProxy)
	orders.Any("/*path", ordersProxy)

	// Product Service endpoints (STRICTLY ADMIN ONLY)
	products := r.Group("/products", jwtMiddleware(), adminOnlyMiddleware())
	products.Any("", productsProxy)
	products.Any("/*path", productsProxy)

	productsV1 := r.Group("/api/v1/products", jwtMiddleware(), adminOnlyMiddleware())
	productsV1.Any("", productsProxy)
	productsV1.Any("/*path", productsProxy)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
