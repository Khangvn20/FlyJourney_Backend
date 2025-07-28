# FlyJourney API Documentation

Tài liệu hướng dẫn sử dụng API cho ứng dụng đặt vé máy bay FlyJourney.

## Base URL
```
http://localhost:3000/api/v1
```

## Authentication
API sử dụng JWT (JSON Web Token) để xác thực. Token được gửi trong header:
```
Authorization: Bearer <jwt_token>
```

## Response Format
Tất cả API responses đều có format:
```json
{
  "status": true/false,
  "errorCode": "SUCCESS|ERROR_CODE",
  "errorMessage": "Message description",
  "data": {}
}
```

---

# User Management APIs

## 1. User Registration
Đăng ký tài khoản người dùng mới.

**Endpoint:** `POST /auth/register`  
**Authentication:** Không cần

### Request Body
```json
{
  "first_name": "Nguyen",
  "last_name": "Van A", 
  "email": "user@example.com",
  "password": "password123",
  "phone": "0123456789",
}
```

### Response Example
```json
{
  "status": true,
  "errorCode": "SUCCESS",
  "errorMessage": "User registered successfully",
  "data": {
    "user": {
      "user_id": 1,
      "first_name": "Nguyen",
      "last_name": "Van A",
      "email": "user@example.com",
      "phone": "0123456789",
      "created_at": "2025-07-29T10:30:00Z"
    }
  }
}
```

## 2. User Login
Đăng nhập và nhận JWT token.

**Endpoint:** `POST /auth/login`  
**Authentication:** Không cần

### Request Body
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

### Response Example
```json
{
  "status": true,
  "errorCode": "SUCCESS", 
  "errorMessage": "Login successful",
  "data": {
    "user": {
      "user_id": 1,
      "first_name": "Nguyen",
      "last_name": "Van A",
      "email": "user@example.com",
      "phone": "0123456789"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2025-07-30T10:30:00Z"
  }
}
```

## 3. User Logout
Đăng xuất và vô hiệu hóa token.

**Endpoint:** `POST /auth/logout`  
**Authentication:** Required

### Headers
```
Authorization: Bearer <jwt_token>
```

### Response Example
```json
{
  "status": true,
  "errorCode": "SUCCESS",
  "errorMessage": "Logout successful"
}
```

## 4. Get User Info
Lấy thông tin chi tiết của user hiện tại.

**Endpoint:** `GET /user/profile`  
**Authentication:** Required

### Headers
```
Authorization: Bearer <jwt_token>
```

### Response Example
```json
{
  "status": true,
  "errorCode": "SUCCESS",
  "errorMessage": "User profile retrieved successfully",
  "data": {
    "user_id": 1,
    "first_name": "Nguyen",
    "last_name": "Van A",
    "email": "user@example.com",
    "phone": "0123456789",
    "created_at": "2025-07-29T10:30:00Z",
    "updated_at": "2025-07-29T10:30:00Z"
  }
}
```

## 5. Update User Info
Cập nhật thông tin người dùng.

**Endpoint:** `PUT /user/profile`  
**Authentication:** Required

### Headers
```
Authorization: Bearer <jwt_token>
```

### Request Body
```json
{
  "first_name": "Nguyen",
  "last_name": "Van B",
  "phone": "0987654321",
  "address": "456 New Street",
}
```

### Response Example
```json
{
  "status": true,
  "errorCode": "SUCCESS",
  "errorMessage": "User profile updated successfully",
  "data": {
    "user_id": 1,
    "first_name": "Nguyen", 
    "last_name": "Van B",
    "email": "user@example.com",
    "phone": "0987654321",
    "address": "456 New Street",
    "date_of_birth": "1990-12-25T00:00:00Z",
    "updated_at": "2025-07-29T11:00:00Z"
  }
}
```

---

# Flight Search APIs

## 6. Search Flights for User
Tìm kiếm chuyến bay cho người dùng (chỉ hiển thị flights available).

**Endpoint:** `POST /user/flights/search`  
**Authentication:** Required

### Headers
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

### Request Body
```json
{
  "departure_airport_code": "HAN",
  "arrival_airport_code": "SGN",
  "departure_date": "01/08/2025",
  "flight_class": "economy",
  "passenger": {
    "adults": 2,
    "children": 1,
    "infants": 0
  },
  "airline_ids": [1, 2], //tùy chọn
  "max_stops": 1,//tùy chọn
  "page": 1,
  "limit": 10,
  "sort_by": "departure_time",
  "sort_order": "asc"
}
```

### Request Parameters
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| departure_airport_code | string | ✅ | Mã sân bay khởi hành (VD: "HAN") |
| arrival_airport_code | string | ✅ | Mã sân bay đến (VD: "SGN") |
| departure_date | string | ✅ | Ngày khởi hành (format: dd/mm/yyyy) |
| flight_class | string | ❌ | Hạng ghế: "economy", "premium_economy", "business", "first", "all" |
| passenger.adults | int | ✅ | Số người lớn (≥1) |
| passenger.children | int | ❌ | Số trẻ em (2-11 tuổi) |
| passenger.infants | int | ❌ | Số em bé (<2 tuổi) |
| airline_ids | array | ❌ | Danh sách ID hãng bay |
| max_stops | int | ❌ | Số điểm dừng tối đa |
| page | int | ❌ | Trang hiện tại (default: 1) |
| limit | int | ❌ | Số items/trang (default: 10) |
| sort_by | string | ❌ | Sắp xếp theo: "departure_time", "price", "duration", "stops" |
| sort_order | string | ❌ | Thứ tự: "asc", "desc" |

### Response Example
```json
{
  "status": true,
  "errorCode": "SUCCESS",
  "errorMessage": "Successfully searched flights for user",
  "data": {
    "search_results": [
      {
        "flight_id": 65,
        "flight_number": "BL190",
        "airline_id": 5,
        "airline_name": "Bamboo Airways",
        "departure_airport_code": "HAN",
        "arrival_airport_code": "SGN",
        "departure_airport": "Sân Bay Nội Bài",
        "arrival_airport": "Sân Bay Tân Sơn Nhất",
        "departure_time": "2025-08-01T09:30:00Z",
        "arrival_time": "2025-08-01T11:45:00Z",
        "duration": 135,
        "stops_count": 0,
        "distance": 1166,
        "flight_class": "economy",
        "total_seats": 150,
        "fare_class_details": {
          "fare_class_code": "Q",
          "cabin_class": "Economy Class",
          "refundable": false,
          "changeable": false,
          "baggage_kg": "1 kiện x 23kg (miễn phí)",
          "description": "Discounted Economy: Vé rẻ, không hoàn tiền, thay đổi hạn chế."
        },
        "pricing": {
          "base_prices": {
            "adult": 876621,
            "child": 876621,
            "infant": 0
          },
          "total_prices": {
            "adult": 1645621,
            "child": 1453371,
            "infant": 192250
          },
          "taxes": {
            "adult": 769000,
            "child": 576750,
            "infant": 192250
          },
          "grand_total": 4744992,
          "currency": "VND"
        },
        "tax_and_fees": 769000
      }
    ],
    "total_count": 15,
    "total_pages": 2,
    "page": 1,
    "limit": 10,
    "departure_airport": "HAN",
    "arrival_airport": "SGN",
    "departure_date": "01/08/2025",
    "flight_class": "economy",
    "passengers": {
      "adults": 2,
      "children": 1,
      "infants": 0
    },
    "sort_by": "departure_time",
    "sort_order": "asc"
  }
}
```

## 7. Get Flight Details for User
Lấy thông tin chi tiết chuyến bay cho người dùng.

**Endpoint:** `GET /user/flights/{flight_id}`  
**Authentication:** Required

### Headers
```
Authorization: Bearer <jwt_token>
```

### Path Parameters
| Parameter | Type | Description |
|-----------|------|-------------|
| flight_id | int | ID của chuyến bay |

### Response Example
```json
{
  "status": true,
  "errorCode": "SUCCESS",
  "errorMessage": "Successfully retrieved flight details for user",
  "data": {
    "flight_id": 65,
    "airline_id": 5,
    "flight_number": "BL190",
    "departure_airport": "Sân Bay Nội Bài",
    "arrival_airport": "Sân Bay Tân Sơn Nhất",
    "departure_time": "2025-08-01T09:30:00Z",
    "arrival_time": "2025-08-01T11:45:00Z",
    "duration_minutes": 135,
    "stops_count": 0,
    "tax_and_fees": 769000,
    "distance": 1166,
    "flight_classes": [
      {
        "flight_class_id": 101,
        "class": "economy",
        "base_price": 876621,
        "available_seats": 45,
        "base_price_child": 876621,
        "base_price_infant": 0,
        "fare_class_code": "Q",
        "fare_class_details": {
          "fare_class_code": "Q",
          "cabin_class": "Economy Class",
          "refundable": false,
          "changeable": false,
          "baggage_kg": "1 kiện x 23kg (miễn phí)",
          "description": "Discounted Economy: Vé rẻ, không hoàn tiền, thay đổi hạn chế."
        }
      },
      {
        "flight_class_id": 102,
        "class": "business",
        "base_price": 2500000,
        "available_seats": 12,
        "base_price_child": 2000000,
        "base_price_infant": 500000,
        "fare_class_code": "C",
        "fare_class_details": {
          "fare_class_code": "C",
          "cabin_class": "Business Class",
          "refundable": true,
          "changeable": true,
          "baggage_kg": "2 kiện x 32kg (miễn phí)",
          "description": "Business Class: Vé linh hoạt, dịch vụ premium."
        }
      }
    ]
  }
}
```

---

# Error Codes

| Error Code | Description |
|------------|-------------|
| SUCCESS | Thành công |
| INVALID_REQUEST | Request không hợp lệ |
| UNAUTHORIZED | Chưa đăng nhập |
| FORBIDDEN | Không có quyền truy cập |
| NOT_FOUND | Không tìm thấy dữ liệu |
| INTERNAL_ERROR | Lỗi server |

---

# Date Format

Tất cả date inputs sử dụng format: **dd/mm/yyyy**
- ✅ Đúng: "01/08/2025", "25/12/2024"
- ❌ Sai: "2025-08-01", "08/01/2025"

---

# Pricing Logic

## Tax Calculation
- **Adult**: 100% tax
- **Child** (2-11 tuổi): 75% tax của adult
- **Infant** (<2 tuổi): 25% tax của adult

## Total Price Calculation
```
Total Price = Base Price + Tax
Grand Total = (Adult Total × Adults) + (Child Total × Children) + (Infant Total × Infants)
```

## Example
Với 2 adults, 1 child:
- Adult total: 1,645,621 VND × 2 = 3,291,242 VND
- Child total: 1,453,371 VND × 1 = 1,453,371 VND  
- **Grand Total: 4,744,613 VND**

---

# Status Codes

## HTTP Status Codes
- **200** - Success
- **400** - Bad Request
- **401** - Unauthorized
- **403** - Forbidden
- **404** - Not Found
- **500** - Internal Server Error

## Flight Status
- **scheduled** - Đã lên lịch
- **boarding** - Đang lên máy bay
- **departed** - Đã khởi hành
- **arrived** - Đã đến
- **delayed** - Bị hoãn
- **cancelled** - Bị hủy

---

# Rate Limiting

API có giới hạn:
- **100 requests/minute** per IP
- **1000 requests/hour** per user

---

# Examples

## Curl Examples

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Search Flights
```bash
curl -X POST http://localhost:8080/api/v1/user/flights/search \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <jwt_token>" \
  -d '{
    "departure_airport_code": "HAN",
    "arrival_airport_code": "SGN",
    "departure_date": "01/08/2025",
    "flight_class": "economy",
    "passenger": {
      "adults": 2,
      "children": 0,
      "infants": 0
    }
  }'
```

### Get Flight Details
```bash
curl -X GET http://localhost:8080/api/v1/user/flights/65 \
  -H "Authorization: Bearer <jwt_token>"
```

### Update User Profile
```bash
curl -X PUT http://localhost:8080/api/v1/user/profile \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <jwt_token>" \
  -d '{
    "first_name": "Nguyen",
    "last_name": "Van B", 
    "phone": "0987654321",
    "address": "456 New Street"
  }'
```

### Logout
```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer <jwt_token>"
```

---

# Common Error Responses

## Validation Error
```json
{
  "status": false,
  "errorCode": "INVALID_REQUEST",
  "errorMessage": "Key: 'FlightSearchRequest.DepartureDate' Error:Field validation for 'DepartureDate' failed on the 'required' tag"
}
```

## Authentication Error
```json
{
  "status": false,
  "errorCode": "UNAUTHORIZED", 
  "errorMessage": "Invalid or expired token"
}
```

## Not Found Error
```json
{
  "status": false,
  "errorCode": "NOT_FOUND",
  "errorMessage": "Flight not found"
}
```

## Internal Server Error
```json
{
  "status": false,
  "errorCode": "INTERNAL_ERROR",
  "errorMessage": "Internal server error occurred"
}
```

---

# Security Notes

1. **JWT Token**: Expires sau 24h, cần refresh hoặc login lại
2. **Password**: Minimum 8 characters, phải có chữ hoa, chữ thường, số
3. **Rate Limiting**: Tránh spam requests
4. **HTTPS**: Production phải sử dụng HTTPS
5. **Validation**: Tất cả inputs đều được validate

---

# Testing with Postman

Import collection với các endpoints trên vào Postman:

1. Set base URL: `http://localhost:8080/api/v1`
2. Login để lấy JWT token
3. Set Authorization header cho các protected endpoints
4. Test các scenarios khác nhau

**Environment Variables:**
- `base_url`: `http://localhost:8080/api/v1`
- `jwt_token`: `<token_from_login_response>`

---

*Lưu ý: Tài liệu này được cập nhật theo phiên bản API hiện tại. Kiểm tra changelog để biết các thay đổi mới nhất.*
