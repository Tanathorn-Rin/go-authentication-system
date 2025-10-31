# Authentication System with Go

A robust REST API authentication system built with Go, Gin framework, MongoDB, and JWT tokens. This project provides user registration, login, and role-based access control (RBAC) features.

## 🚀 Features

- ✅ User Registration & Login
- ✅ JWT Authentication (Access & Refresh Tokens)
- ✅ Password Hashing with bcrypt
- ✅ Role-Based Access Control (ADMIN/USER)
- ✅ MongoDB Integration
- ✅ Input Validation
- ✅ Protected Routes with Middleware

## 📋 Prerequisites

Before running this project, make sure you have the following installed:

- [Go](https://golang.org/dl/) (version 1.20 or higher)
- [MongoDB](https://www.mongodb.com/try/download/community) or [Docker](https://www.docker.com/products/docker-desktop)
- [Git](https://git-scm.com/downloads)

## 🛠️ Installation

### 1. Clone the Repository

```bash
git clone https://github.com/Tanathorn-Rin/go-authentication-system.git
cd go-authentication-system
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Set Up MongoDB

#### Option A: Using Docker (Recommended for Quick Start)

```bash
docker run -d -p 27017:27017 --name mongodb-auth mongo:latest
```

#### Option B: Install MongoDB Locally (macOS)

```bash
brew tap mongodb/brew
brew install mongodb-community
brew services start mongodb-community
```

#### Option C: Using MongoDB Atlas (Cloud)

1. Create a free account at [MongoDB Atlas](https://www.mongodb.com/cloud/atlas)
2. Create a cluster and get your connection string
3. Update the connection string in `config/database.go`

### 4. Run the Application

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## 📚 API Endpoints

### Public Routes

#### 1. User Signup
```http
POST /signup
Content-Type: application/json

{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john@example.com",
  "password": "securepassword123",
  "phone": "1234567890",
  "role": "USER"
}
```

**Response:**
```json
{
  "message": "User created successfully"
}
```

#### 2. User Login
```http
POST /login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "securepassword123"
}
```

**Response:**
```json
{
  "user": { ... },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Protected Routes (Requires Authentication)

#### 3. Get All Users (ADMIN only)
```http
GET /users
Authorization: Bearer <your_jwt_token>
```

**Response:**
```json
[
  {
    "user_id": "...",
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "role": "USER",
    ...
  }
]
```

#### 4. Get User by ID
```http
GET /users/:id
Authorization: Bearer <your_jwt_token>
```

**Note:** Regular users can only access their own profile. ADMIN can access any user profile.

**Response:**
```json
{
  "user_id": "...",
  "first_name": "John",
  "last_name": "Doe",
  "email": "john@example.com",
  "role": "USER",
  ...
}
```

## 🔐 Authentication

This API uses JWT (JSON Web Tokens) for authentication. After logging in, include the token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

### Token Expiration
- **Access Token:** 24 hours
- **Refresh Token:** 7 days

## 🏗️ Project Structure

```
Authentication-system/
├── config/
│   ├── auth-key.go        # JWT key generation
│   └── database.go        # MongoDB connection
├── controllers/
│   └── userControllers.go # User-related handlers
├── helpers/
│   └── token.go           # JWT & password utilities
├── middleware/
│   └── auth.go            # JWT authentication middleware
├── models/
│   └── user.go            # User model/schema
├── routes/
│   └── routes.go          # Route definitions
├── main.go                # Application entry point
├── go.mod                 # Go module dependencies
└── go.sum                 # Dependency checksums
```

## 🔧 Configuration

### Database
Update MongoDB connection string in `config/database.go`:
```go
clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
```

### Database Name
The default database name is `usersdb`. Change it in `config/database.go` if needed:
```go
return Client.Database("usersdb").Collection(collectionName)
```

### Server Port
The server runs on port `8080` by default. Change it in `main.go`:
```go
port := "8080"
```

## 📦 Dependencies

- [gin-gonic/gin](https://github.com/gin-gonic/gin) - HTTP web framework
- [go-playground/validator](https://github.com/go-playground/validator) - Input validation
- [golang-jwt/jwt](https://github.com/golang-jwt/jwt) - JWT implementation
- [mongo-driver](https://github.com/mongodb/mongo-go-driver) - MongoDB driver
- [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) - Password hashing

## 🧪 Testing the API

### Using cURL

**Signup:**
```bash
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Jane",
    "last_name": "Smith",
    "email": "jane@example.com",
    "password": "password123",
    "phone": "9876543210",
    "role": "USER"
  }'
```

**Login:**
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "jane@example.com",
    "password": "password123"
  }'
```

**Get Users (with token):**
```bash
curl -X GET http://localhost:8080/users \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### Using Postman

1. Import the endpoints listed above
2. For protected routes, add the token in the Authorization tab:
   - Type: Bearer Token
   - Token: Paste your JWT token

## 🛡️ Security Features

- **Password Hashing:** Passwords are hashed using bcrypt before storage
- **JWT Tokens:** Secure token-based authentication
- **Role-Based Access Control:** Different permissions for ADMIN and USER roles
- **Input Validation:** All inputs are validated before processing
- **Secure Headers:** Authorization headers for protected routes

## 🐛 Troubleshooting

### MongoDB Connection Error
```
server selection error: server selection timeout
```
**Solution:** Make sure MongoDB is running:
```bash
# If using Docker
docker start mongodb-auth

# If using Homebrew
brew services start mongodb-community
```

### Port Already in Use
```
bind: address already in use
```
**Solution:** Change the port in `main.go` or kill the process using port 8080:
```bash
lsof -ti:8080 | xargs kill -9
```

### Module Errors
```
no required module provides package
```
**Solution:** Run:
```bash
go mod tidy
go mod download
```

## 📝 License

This project is open source and available under the [MIT License](LICENSE).

## 👤 Author

**Tanathorn Rin**
- GitHub: [@Tanathorn-Rin](https://github.com/Tanathorn-Rin)

## 🤝 Contributing

Contributions, issues, and feature requests are welcome! Feel free to check the [issues page](https://github.com/Tanathorn-Rin/go-authentication-system/issues).

## ⭐ Show your support

Give a ⭐️ if this project helped you!

---

**Note:** This is a learning project. For production use, consider adding:
- Environment variables for sensitive data
- HTTPS/TLS support
- Rate limiting
- Email verification
- Password reset functionality
- Refresh token rotation
- Better error handling
- Unit tests
- API documentation (Swagger/OpenAPI)
