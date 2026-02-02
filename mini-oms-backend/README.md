# Mini Order Management System - Backend

Backend API untuk Mini Order Management System menggunakan Golang Echo Framework, PostgreSQL, dan GORM.

## Tech Stack

- **Framework**: Echo v4
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT (golang-jwt/jwt)
- **Password Hashing**: bcrypt

## Project Structure

```
mini-oms-backend/
├── cmd/api/              # Application entry point
├── internal/
│   ├── config/           # Configuration management
│   ├── db/               # Database connection
│   ├── middlewares/      # JWT, RBAC middlewares
│   ├── models/           # GORM models
│   ├── modules/          # Business modules
│   │   ├── auth/         # Authentication (register, login)
│   │   ├── category/     # Product categories
│   │   ├── product/      # Product management
│   │   ├── order/        # Order management
│   │   └── payment/      # Payment simulation
│   └── utils/            # Helper functions
└── migrations/           # SQL migration files
```

## Setup

1. **Clone repository**
```bash
git clone <repository-url>
cd mini-oms-backend
```

2. **Install dependencies**
```bash
go mod download
```

3. **Setup PostgreSQL database**
```bash
createdb mini_oms
```

4. **Configure environment**
```bash
cp .env.example .env
# Edit .env dengan konfigurasi database Anda
```

5. **Run migrations** (opsional, menggunakan GORM auto-migrate)
```bash
go run cmd/api/main.go
```

6. **Run application**
```bash
go run cmd/api/main.go
```

Server akan berjalan di `http://localhost:8080`

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register user baru
- `POST /api/auth/login` - Login dan dapatkan JWT token

### Products (Public)
- `GET /api/products` - List semua products
- `GET /api/products/:id` - Detail product

### Products (Admin Only)
- `POST /api/products` - Create product
- `PUT /api/products/:id` - Update product
- `DELETE /api/products/:id` - Delete product

### Orders (Protected)
- `GET /api/orders` - List orders (user: own orders, admin: all)
- `GET /api/orders/:id` - Detail order
- `POST /api/orders` - Create order

### Payments (Protected)
- `POST /api/payments` - Create payment
- `GET /api/payments/order/:orderId` - Get payment by order

## Configuration Choices

### 1. Configuration Management: `os.Getenv`
**Alasan**: 
- Lebih sederhana dan tidak memerlukan dependency tambahan
- Cocok untuk project kecil-menengah
- Standard library Go sudah cukup untuk read environment variables
- Lebih ringan dibanding Viper

Jika project berkembang lebih besar dan memerlukan multiple config sources (file, env, remote), baru pertimbangkan Viper.

### 2. Database Migration: GORM AutoMigrate
**Alasan**:
- Lebih cepat untuk development
- Otomatis sinkronisasi model dengan schema
- Cocok untuk MVP dan prototyping
- Tidak perlu maintain migration files terpisah

Untuk production, disarankan menggunakan migration tool seperti `golang-migrate` atau `goose` untuk kontrol yang lebih baik.

### 3. Primary Key: UUID
**Alasan**:
- Lebih secure (tidak bisa ditebak)
- Distributed-friendly
- Menghindari information leakage dari sequential ID

## Example Requests

### Register
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "SecurePass123!"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "SecurePass123!"
  }'
```

### Create Order (Protected)
```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "items": [
      {"product_id": "uuid-here", "quantity": 2},
      {"product_id": "uuid-here", "quantity": 1}
    ],
    "notes": "Tolong kirim pagi hari"
  }'
```

## License

MIT
