# Luxbiss API Documentation

Welcome to the Luxbiss API reference. This document provides exact request and response structures for all endpoints.

## 📌 Global Constants

- **Base URL:** `http://localhost:8080/api/v1`
- **Content-Type:** `application/json`
- **Authorization:** `Bearer <JWT_ACCESS_TOKEN>`

---

## 🏥 Health Module (`/health`)

### 1. Simple Health Check
Checks if the server is up and running.

- **Method:** `GET`
- **Path:** `/health`
- **Auth Required:** No
- **Response (200 OK):**
```json
{
  "success": true,
  "message": "Server is healthy",
  "data": {
    "status": "ok"
  },
  "request_id": "8b9cad0e-1f20..."
}
```

---

### 2. Readiness Check
Checks if the server is ready to handle requests.

- **Method:** `GET`
- **Path:** `/health/ready`
- **Auth Required:** No
- **Response (200 OK):**
```json
{
  "success": true,
  "message": "Server is ready",
  "data": {
    "status": "ready"
  },
  "request_id": "8b9cad0e-1f20..."
}
```

---

## 🔐 Authentication Module (`/auth`)

### 1. Register a New Account
Creates a new user and returns authentication tokens.

- **Method:** `POST`
- **Path:** `/auth/register`
- **Request Body:**
```json
{
  "name": "Jane Smith",
  "email": "jane@example.com",
  "password": "StrongPassword123!"
}
```
| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `name` | string | Yes | 2-100 characters. |
| `email` | string | Yes | Valid email format. |
| `password` | string | Yes | Min 8 chars, 1 uppercase, 1 lowercase, 1 number, 1 symbol. |

- **Response (201 Created):**
```json
{
  "success": true,
  "message": "Registration successful",
  "data": {
    "access_token": "eyJhbG...",
    "refresh_token": "eyJhbG...",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Jane Smith",
      "email": "jane@example.com",
      "role": "user",
      "is_active": true,
      "created_at": "2024-02-27T15:04:05Z",
      "updated_at": "2024-02-27T15:04:05Z"
    }
  },
  "request_id": "8b9cad0e-1f20..."
}
```

---

### 2. Login
Authenticates an existing user.

- **Method:** `POST`
- **Path:** `/auth/login`
- **Request Body:**
```json
{
  "email": "jane@example.com",
  "password": "StrongPassword123!"
}
```

- **Response (200 OK):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbG...",
    "refresh_token": "eyJhbG...",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Jane Smith",
      "email": "jane@example.com",
      "role": "user",
      "is_active": true,
      "created_at": "2024-02-27T15:04:05Z",
      "updated_at": "2024-02-27T15:04:05Z"
    }
  }
}
```

---

### 3. Refresh Access Token
Exchange a refresh token for a new set of tokens.

- **Method:** `POST`
- **Path:** `/auth/refresh`
- **Request Body:**
```json
{
  "refresh_token": "eyJhbG..."
}
```

- **Response (200 OK):**
```json
{
  "success": true,
  "message": "Token refreshed successfully",
  "data": {
    "access_token": "eyJhbG...",
    "refresh_token": "eyJhbG...",
    "user": { ... }
  }
}
```

---

### 4. Forgot Password
Triggers an OTP email for password recovery.

- **Method:** `POST`
- **Path:** `/auth/forgot-password`
- **Request Body:**
```json
{
  "email": "jane@example.com"
}
```

- **Response (200 OK):**
```json
{
  "success": true,
  "message": "If an account exists with this email, you will receive an OTP",
  "request_id": "..."
}
```

---

### 5. Reset Password
Changes password using the OTP received via email.

- **Method:** `POST`
- **Path:** `/auth/reset-password`
- **Request Body:**
```json
{
  "email": "jane@example.com",
  "otp": "123456",
  "password": "NewStrongPassword456!"
}
```

- **Response (200 OK):**
```json
{
  "success": true,
  "message": "Password reset successful"
}
```

---

## 👤 User Module (`/users`)
*Requires Authentication Header: `Authorization: Bearer <token>`*

### 1. Get My Profile
Returns details of the currently authenticated user.

- **Method:** `GET`
- **Path:** `/users/me`
- **Response (200 OK):**
```json
{
  "success": true,
  "message": "User retrieved successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Jane Smith",
    "email": "jane@example.com",
    "role": "user",
    "is_active": true,
    "created_at": "2024-02-27T15:04:05Z",
    "updated_at": "2024-02-27T15:04:05Z"
  }
}
```

---

### 2. List All Users (Admin Restricted)
- **Method:** `GET`
- **Path:** `/users?page=1&per_page=20`
- **Response (200 OK):**
```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Jane Smith",
      "email": "jane@example.com",
      "role": "user",
      "is_active": true,
      "created_at": "2024-02-27T15:04:05Z",
      "updated_at": "2024-02-27T15:04:05Z"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 125
  }
}
```

---

### 3. Update User Profile
- **Method:** `PUT`
- **Path:** `/users/:id`
- **Request Body:**
```json
{
  "name": "Jane Doe Updated",
  "email": "jane.updated@example.com",
  "is_active": true
}
```
*Note: All fields are optional.*

- **Response (200 OK):**
```json
{
  "success": true,
  "message": "User updated successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Jane Doe Updated",
    "email": "jane.updated@example.com",
    "role": "user",
    "is_active": true,
    "created_at": "2024-02-27T15:04:05Z",
    "updated_at": "2024-02-27T15:30:00Z"
  }
}
```

---

## 🧪 Common Error Scenarios

### Validation Error (400 Bad Request)
Returned when request body format is incorrect or fails validation rules.

```json
{
  "success": false,
  "message": "Validation failed",
  "errors": [
    {
      "field": "password",
      "message": "must be at least 8 characters"
    },
    {
      "field": "email",
      "message": "must be a valid email address"
    }
  ],
  "request_id": "..."
}
```

### Unauthorized Error (401 Unauthorized)
Returned when the token is missing, invalid, or expired.

```json
{
  "success": false,
  "message": "Invalid or expired token",
  "request_id": "..."
}
```

### Forbidden Error (403 Forbidden)
Returned when a user tries to access a route they don't have the role for (e.g., non-admin accessing User List).

```json
{
  "success": false,
  "message": "You do not have permission to perform this action",
  "request_id": "..."
}
```
