# Customer, Products, and Orders Management System (Go Version)

## Overview

This is a Go-based backend service for managing customers, products with hierarchical categories, and orders. The system includes authentication via OpenID Connect, REST APIs for product management and order processing, and integrations with Africa's Talking SMS gateway for customer notifications.

## Features

- **Product Management**:
    - Hierarchical category system (unlimited depth)
    - Product upload and categorization
    - Average price calculation by category

- **Order Processing**:
    - Order creation and management
    - SMS notifications to customers via Africa's Talking
    - Email notifications to administrators

- **Authentication**:
    - OpenID Connect for customer authentication

- **Database**:
    - PostgreSQL with proper schema design
    - Database migrations

- **Testing & Deployment**:
    - Unit tests with coverage checking
    - CI/CD pipeline
    - Containerized deployment (Docker)

## Prerequisites

- Docker
- Docker Compose
- Go 1.20+
- Africa's Talking API credentials (for SMS functionality)
- SMTP server credentials (for email notifications)

## Getting Started

1. **Clone the repository**:
   ```bash
   git clone https://github.com/Mutonya/Savannah-Informatics---Backend-Developer-.git
   cd docker-compose up --build
   ```

2. **Set up environment variables**:
   Create a `.env` file in the project root with the following variables:
   ```
   DB_HOST=db
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=savannah
   SSL_MODE=disable
   
   APP_PORT=8080
   
   AFRICAS_TALKING_API_KEY=your-api-key
   AFRICAS_TALKING_USERNAME=your-username
   CURRENCY=Ksh
   SMSSENDERID=
   
   SMTP_HOST=your-smtp-server
   SMTP_PORT=587
   SMTP_USER=your-email
   SMTP_PASSWORD=your-password

   OAUTH_CLIENT_ID=
   OAUTH_CLIENT_SECRET=
   OAUTH_REDIRECT_URL=http://localhost:8080/auth/callback
   OAUTH_PROVIDER_URL=https://accounts.google.com
   
   ```

3. **Build and start the containers**:
   ```bash
   docker-compose up --build
   ```

4. **Run database migrations**:
   ```bash
   docker-compose exec app ./migrate -path ./migrations -database "postgres://postgres:postgres@db:5432/savannah?sslmode=disable" up
   ```

5. **Access the application**:
   The API will be available at `http://localhost:8080/api/v1/`

## API Endpoints

### Authentication
- `GET /auth/login` - Initiate OAuth2/OIDC login flow
- `GET /auth/callback` - OAuth2 callback handler

### API v1 (Authenticated)
All API v1 routes require valid JWT authentication.

#### Customer
- `GET /api/v1/profile` - Get current user profile

#### Products
- `POST /api/v1/products` - Create new product
- `GET /api/v1/products` - List all products
- `GET /api/v1/products/:id` - Get product details
- `PUT /api/v1/products/:id` - Update product
- `DELETE /api/v1/products/:id` - Delete product

#### Categories
- `POST /api/v1/categories` - Create new category
- `GET /api/v1/categories` - List all categories
- `GET /api/v1/categories/:id` - Get category details
- `GET /api/v1/categories/:id/products` - Get products in category
- `GET /api/v1/categories/:id/average-price` - Get average price for category
- `PUT /api/v1/categories/:id` - Update category
- `DELETE /api/v1/categories/:id` - Delete category

#### Orders
- `POST /api/v1/orders` - Create new order
- `GET /api/v1/orders` - List user's orders
- `GET /api/v1/orders/:id` - Get order details
- `PUT /api/v1/orders/:id/status` - Update order status

## Authentication Flow

1. Client accesses `/auth/login`
2. Server redirects to OIDC provider
3. User authenticates with provider
4. Provider redirects to `/auth/callback`
5. Server exchanges code for tokens
6. Server issues JWT to client
7. Client uses JWT for API v1 requests
## Running Tests

```bash
docker-compose exec app go test -v ./...
```

To check test coverage:
```bash
docker-compose exec app go test -coverprofile=coverage.out ./...
docker-compose exec app go tool cover -html=coverage.out
```

## Project Structure

```
.
├── cmd/              # Main application entry points
├── internal/
│   ├── auth/         # Authentication handlers
│   ├── config/       # Configuration loading
│   ├── db/           # Database connection and migrations
│   ├── handlers/     # HTTP request handlers
│   ├── models/       # Database models
│   ├── services/     # Business logic
│   └── utils/        # Utility functions
├── migrations/       # Database migrations
├── pkg/              # Reusable packages
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── go.sum
```

## Deployment

The application is configured for deployment with Docker. For Kubernetes deployment:

1. Build the Docker image:
   ```bash
   docker build -t savannah-go-backend .
   ```


## CI/CD

The project includes GitHub Actions configuration for:
- Automated testing on push
- Coverage reporting
- Docker image building
- Deployment to staging/production (configured via repository secrets)


```

## License

This project is proprietary software developed for Savannah Informatics.
