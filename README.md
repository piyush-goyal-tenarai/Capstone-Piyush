# Production-Ready Order Management System

A production-grade Order Management REST API built in Go with PostgreSQL, implementing layered architecture, resilience patterns, observability, security, Docker, CI/CD, and AWS deployment readiness.

## Tech Stack

- **Go 1.21** with Gin framework
- **PostgreSQL** for persistence
- **Docker** for containerization
- **GitHub Actions** for CI/CD
- **AWS** (EC2/ECS, ALB, RDS, CloudWatch, Secrets Manager)

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/orders | Create order |
| GET | /api/v1/orders/{id} | Get order by ID |
| GET | /api/v1/orders | List all orders |
| PUT | /api/v1/orders/{id}/status | Update order status |
| GET | /health/live | Liveness probe |
| GET | /health/ready | Readiness probe (checks DB) |
| GET | /metrics | Application metrics |

## Order Statuses

`CREATED` → `CONFIRMED` → `CANCELLED`

## Architecture

```
Client → Load Balancer → Gin Router → Middleware Chain → Controller → Service → Repository → PostgreSQL
                                          │
                                     Request ID
                                     Logging
                                     Recovery
                                     Metrics
                                     Secure Headers
```

## Quick Start

```bash
# 1. Create database
psql -U postgres -c "CREATE DATABASE orderdb;"

# 2. Copy env file
cp .env.example .env

# 3. Install dependencies
go mod tidy

# 4. Run
go run main.go

# 5. Test
go test -v -race ./tests/...
```

## Docker

```bash
docker build -t order-management .
docker run -p 8080:8080 -e DB_HOST=host.docker.internal -e DB_PASSWORD=postgres order-management
```

## SoapUI Testing

Import the endpoints into SoapUI:
- Base URL: `http://localhost:8080`
- Create Order: POST `/api/v1/orders` with JSON body `{"customerId":101,"productId":"P001","quantity":2}`
- List Orders: GET `/api/v1/orders`
- Get Order: GET `/api/v1/orders/1`
- Update Status: PUT `/api/v1/orders/1/status` with JSON body `{"status":"CONFIRMED"}`

## Project Structure

```
order-management/
├── main.go                    Entry point + graceful shutdown
├── go.mod                     Dependencies
├── Dockerfile                 Multi-stage Docker build
├── .env.example               Environment template
├── config/config.go           Environment-based configuration
├── model/order.go             Domain models + DTOs
├── controller/                HTTP handlers (request/response)
├── service/                   Business logic + product price fetch
├── repository/                PostgreSQL data access
├── middleware/                 Logging, recovery, request ID, security headers
├── resilience/                Retry, circuit breaker
├── metrics/                   Request count, latency, error tracking
├── database/                  PostgreSQL connection + pooling
├── tests/                     Unit tests + race condition tests
├── scripts/init.sql           Database schema + seed data
├── aws/deploy.sh              AWS deployment commands
└── .github/workflows/ci.yml   CI/CD pipeline
```
