# API Examples & Requests

This document provides example requests and responses for the Mini OMS API.

## Base URL
```
http://localhost:8080/api
```

---

## 1. Authentication

### POST /api/auth/register
Register a new user account.

**Request:**
```json
POST /api/auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "SecurePass123!"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "user",
      "created_at": "2026-02-02T14:00:00Z",
      "updated_at": "2026-02-02T14:00:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNTUwZTg0MDAtZTI5Yi00MWQ0LWE3MTYtNDQ2NjU1NDQwMDAwIiwiZW1haWwiOiJqb2huQGV4YW1wbGUuY29tIiwicm9sZSI6InVzZXIiLCJleHAiOjE3MDcwNDMyMDB9.xyz...",
    "token_type": "Bearer",
    "expires_in": 86400
  }
}
```

### POST /api/auth/login
Login with existing credentials.

**Request:**
```json
POST /api/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "SecurePass123!"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "user",
      "created_at": "2026-02-02T14:00:00Z",
      "updated_at": "2026-02-02T14:00:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 86400
  }
}
```

---

## 2. Products

### GET /api/products
Get all products (public).

**Request:**
```
GET /api/products
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Products retrieved successfully",
  "data": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440001",
      "category_id": "770e8400-e29b-41d4-a716-446655440002",
      "name": "Laptop Gaming ASUS ROG",
      "description": "High performance gaming laptop",
      "price": 15000000,
      "stock": 5,
      "image_url": "https://example.com/laptop.jpg",
      "created_at": "2026-02-01T10:00:00Z",
      "updated_at": "2026-02-01T10:00:00Z"
    }
  ]
}
```

### POST /api/products
Create new product (admin only).

**Request:**
```json
POST /api/products
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "name": "Mouse Wireless Logitech",
  "description": "Wireless gaming mouse with RGB",
  "price": 500000,
  "stock": 25,
  "category_id": "770e8400-e29b-41d4-a716-446655440002",
  "image_url": "https://example.com/mouse.jpg"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Product created successfully",
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440003",
    "name": "Mouse Wireless Logitech",
    "description": "Wireless gaming mouse with RGB",
    "price": 500000,
    "stock": 25,
    "image_url": "https://example.com/mouse.jpg",
    "created_at": "2026-02-02T14:30:00Z",
    "updated_at": "2026-02-02T14:30:00Z"
  }
}
```

---

## 3. Orders

### POST /api/orders
Create new order (protected).

**Request:**
```json
POST /api/orders
Authorization: Bearer <user-token>
Content-Type: application/json

{
  "items": [
    {
      "product_id": "660e8400-e29b-41d4-a716-446655440001",
      "quantity": 1
    },
    {
      "product_id": "660e8400-e29b-41d4-a716-446655440003",
      "quantity": 2
    }
  ],
  "notes": "Tolong kirim dengan bubble wrap"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Order created successfully",
  "data": {
    "id": "880e8400-e29b-41d4-a716-446655440010",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "order_number": "ORD-20260202-1001",
    "total_amount": 16000000,
    "status": "created",
    "notes": "Tolong kirim dengan bubble wrap",
    "created_at": "2026-02-02T15:00:00Z",
    "updated_at": "2026-02-02T15:00:00Z",
    "items": [
      {
        "id": "990e8400-e29b-41d4-a716-446655440020",
        "order_id": "880e8400-e29b-41d4-a716-446655440010",
        "product_id": "660e8400-e29b-41d4-a716-446655440001",
        "product_name": "Laptop Gaming ASUS ROG",
        "product_price": 15000000,
        "quantity": 1,
        "subtotal": 15000000,
        "created_at": "2026-02-02T15:00:00Z"
      },
      {
        "id": "990e8400-e29b-41d4-a716-446655440021",
        "order_id": "880e8400-e29b-41d4-a716-446655440010",
        "product_id": "660e8400-e29b-41d4-a716-446655440003",
        "product_name": "Mouse Wireless Logitech",
        "product_price": 500000,
        "quantity": 2,
        "subtotal": 1000000,
        "created_at": "2026-02-02T15:00:00Z"
      }
    ]
  }
}
```

### GET /api/orders
Get all orders (user sees own, admin sees all).

**Request:**
```
GET /api/orders
Authorization: Bearer <user-token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Orders retrieved successfully",
  "data": [
    {
      "id": "880e8400-e29b-41d4-a716-446655440010",
      "user_id": "550e8400-e29b-41d4-a716-446655440000",
      "order_number": "ORD-20260202-1001",
      "total_amount": 16000000,
      "status": "created",
      "notes": "Tolong kirim dengan bubble wrap",
      "created_at": "2026-02-02T15:00:00Z",
      "updated_at": "2026-02-02T15:00:00Z"
    }
  ]
}
```

---

## 4. Payments

### POST /api/payments
Create payment for an order (protected).

**Request:**
```json
POST /api/payments
Authorization: Bearer <user-token>
Content-Type: application/json

{
  "order_id": "880e8400-e29b-41d4-a716-446655440010",
  "payment_method": "bank_transfer",
  "payment_proof_url": "https://example.com/uploads/proof-123.jpg",
  "notes": "Transfer dari BCA a/n John Doe"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Payment created successfully",
  "data": {
    "id": "aa0e8400-e29b-41d4-a716-446655440030",
    "order_id": "880e8400-e29b-41d4-a716-446655440010",
    "payment_number": "PAY-20260202-3001",
    "amount": 16000000,
    "payment_method": "bank_transfer",
    "status": "pending",
    "payment_proof_url": "https://example.com/uploads/proof-123.jpg",
    "verified_by": null,
    "verified_at": null,
    "notes": "Transfer dari BCA a/n John Doe",
    "created_at": "2026-02-02T15:10:00Z",
    "updated_at": "2026-02-02T15:10:00Z"
  }
}
```

### GET /api/payments/order/:orderId
Get payment by order ID (protected).

**Request:**
```
GET /api/payments/order/880e8400-e29b-41d4-a716-446655440010
Authorization: Bearer <user-token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Payment retrieved successfully",
  "data": {
    "id": "aa0e8400-e29b-41d4-a716-446655440030",
    "order_id": "880e8400-e29b-41d4-a716-446655440010",
    "payment_number": "PAY-20260202-3001",
    "amount": 16000000,
    "payment_method": "bank_transfer",
    "status": "pending",
    "payment_proof_url": "https://example.com/uploads/proof-123.jpg",
    "notes": "Transfer dari BCA a/n John Doe",
    "created_at": "2026-02-02T15:10:00Z",
    "updated_at": "2026-02-02T15:10:00Z"
  }
}
```

---

## Error Responses

### 400 Bad Request
```json
{
  "success": false,
  "message": "Invalid request body"
}
```

### 401 Unauthorized
```json
{
  "success": false,
  "message": "Invalid or expired token"
}
```

### 403 Forbidden
```json
{
  "success": false,
  "message": "Access forbidden: admin only"
}
```

### 404 Not Found
```json
{
  "success": false,
  "message": "Product not found"
}
```

---

## cURL Examples

### Register
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","password":"SecurePass123!"}'
```

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"SecurePass123!"}'
```

### Create Order
```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "items": [
      {"product_id":"660e8400-e29b-41d4-a716-446655440001","quantity":1}
    ],
    "notes": "Urgent delivery"
  }'
```

### Create Payment
```bash
curl -X POST http://localhost:8080/api/payments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "order_id":"880e8400-e29b-41d4-a716-446655440010",
    "payment_method":"bank_transfer",
    "payment_proof_url":"https://example.com/proof.jpg",
    "notes":"Paid via BCA"
  }'
```
