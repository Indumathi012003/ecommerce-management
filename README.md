# Order & Product Management System (Microservices)

## Architecture Overview
The project is organized cleanly with separate layered services:
- **`gateway/`**: Single entrypoint (`:8080`) providing JWT Authentication, reverse proxying, and role-based access control.
- **`services/order/`**: Full layered Order Service (`:8081`): Controller, Service, Repository, Resilient Product Client, and GraphQL.
- **`services/product/`**: Full layered Product Service (`:8082`): Controller, Service, Repository. Accessible strictly by **ADMIN**!
- **`shared/`**: `config/`, `database/`, `middleware/`, and `model/`.

## Product Service Layered Features (ADMIN ONLY)
1. **Add Product**: `POST /products` or `POST /api/v1/products`
2. **View Product(s)**: `GET /products/:id` and `GET /products`
3. **Update Product**: `PUT /products/:id`
4. **Delete Product**: `DELETE /products/:id`

*Security Rule*: Gateway blocks non-admin users with `403 Forbidden` when attempting to access `/products`!

## Order Service Features
- `POST /api/v1/orders` (Creates order, fetches live price via resilient client)
- `GET /api/v1/orders/:id`
- `GET /api/v1/orders`
- `PUT /api/v1/orders/:id/status`

## Resilience in Product Client
- Request Timeout (1500ms)
- Finite Retries (max 3)
- Exponential Backoff (`100ms * 2^(attempt-1)`)
- Circuit Breaker (`gobreaker` trips after 3 consecutive failures to `OPEN`)
- Context Cancellation

## Run Locally
```bash
# 1. Start Product Service
go run ./services/product/main.go

# 2. Start Order Service
go run ./services/order/main.go

# 3. Start Gateway
go run ./gateway/main.go
```

## Run Tests
```bash
go test -v -race ./tests/...
```
