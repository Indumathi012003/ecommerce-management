package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"order-management/model"

	"github.com/sony/gobreaker"
)

// Exported controlled errors for controllers and callers
var (
	ErrCircuitOpen        = errors.New("circuit breaker is OPEN: product service unavailable")
	ErrServiceUnavailable = errors.New("product service unavailable after max retries")
)

type ProductClient interface {
	GetProductPrice(ctx context.Context, productID int64) (float64, error)
}

type resilientProductClient struct {
	baseURL    string
	httpClient *http.Client
	cb         *gobreaker.CircuitBreaker
}

// Accepts optional baseURL: NewProductClient() or NewProductClient("http://...")
func NewProductClient(baseURL ...string) ProductClient {
	url := ""
	if len(baseURL) > 0 && baseURL[0] != "" {
		url = baseURL[0]
	} else {
		url = os.Getenv("PRODUCT_SERVICE_URL")
		if url == "" {
			url = "http://localhost:8082"
		}
	}

	cbSettings := gobreaker.Settings{
		Name:        "ProductServiceCB",
		MaxRequests: 2,
		Interval:    10 * time.Second,
		Timeout:     5 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 3
		},
	}

	return &resilientProductClient{
		baseURL: url,
		httpClient: &http.Client{
			Timeout: 1500 * time.Millisecond,
		},
		cb: gobreaker.NewCircuitBreaker(cbSettings),
	}
}

func (c *resilientProductClient) GetProductPrice(ctx context.Context, productID int64) (float64, error) {
	result, err := c.cb.Execute(func() (interface{}, error) {
		var lastErr error
		url := fmt.Sprintf("%s/products/%d", c.baseURL, productID)

		const maxRetries = 3
		for attempt := 1; attempt <= maxRetries; attempt++ {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				return nil, err
			}

			resp, err := c.httpClient.Do(req)
			if err == nil {
				defer resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					var p model.Product
					if err := json.NewDecoder(resp.Body).Decode(&p); err == nil {
						return p.Price, nil
					}
				}
				lastErr = fmt.Errorf("status %d", resp.StatusCode)
			} else {
				lastErr = err
			}

			// Exponential Backoff: 100ms, 200ms, 400ms
			backoff := time.Duration(100*(1<<(attempt-1))) * time.Millisecond
			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		return nil, fmt.Errorf("%w: %v", ErrServiceUnavailable, lastErr)
	})

	if err != nil {
		if errors.Is(err, gobreaker.ErrOpenState) || errors.Is(err, gobreaker.ErrTooManyRequests) {
			return 0, ErrCircuitOpen
		}
		return 0, err
	}

	return result.(float64), nil
}
