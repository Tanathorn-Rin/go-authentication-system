package controllers

import (
	"authentication/config"
	"authentication/helpers"
	"authentication/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// validate is the validator instance used for struct validation
var validate = validator.New()

// userCollection is the MongoDB collection for user documents
var userCollection = config.OpenCollection("users")

// Signup handles user registration
// It validates input, checks for duplicate email/phone, hashes password,
// generates JWT tokens, and creates a new user in the database
// Returns: gin.HandlerFunc that processes the signup request
func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user models.User

		// Bind JSON request body to user struct
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate user input against struct validation tags
		if validationErr := validate.Struct(user); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Check if email or phone already exists in database
		count, err := userCollection.CountDocuments(ctx, bson.M{
			"$or": []bson.M{
				{"email": user.Email},
				{"phone": user.Phone},
			},
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Prevent duplicate registrations
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email or phone already exists"})
			return
		}

		// Generate additional user data (hash password, create ID, generate tokens)
		hashedPassword := helpers.HashPassword(*user.Password)
		user.Password = &hashedPassword
		user.Created_at = time.Now()
		user.Updated_at = time.Now()
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		accessToken, refreshToken := helpers.GenerateTokens(*user.Email, user.User_id, *user.Role)
		user.Token = &accessToken
		user.Refresh_token = &refreshToken

		_, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": insertErr.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
	}
}

// Login handles user authentication
// It verifies email and password, generates new JWT tokens,
// updates tokens in database, and returns user data with tokens
// Returns: gin.HandlerFunc that processes the login request
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		var user models.User
		var foundUser models.User

		// Bind JSON request body to user struct
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Find user by email in database
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
			return
		}

		// Verify the provided password matches the stored hashed password
		passwordIsValid, msg := helpers.VerifyPassword(*foundUser.Password, *user.Password)
		if !passwordIsValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			return
		}

		// Generate new access and refresh tokens for the user
		token, refreshToken := helpers.GenerateTokens(*foundUser.Email, foundUser.User_id, *foundUser.Role)
		// Update the tokens in the database
		helpers.UpdateAllTokens(token, refreshToken, foundUser.User_id)

		// Return user data along with tokens
		c.JSON(http.StatusOK, gin.H{
			"user":          foundUser,
			"token":         token,
			"refresh_token": refreshToken,
		})
	}
}

// GetUser retrieves a specific user by ID
// Regular users can only access their own profile, ADMIN can access any profile
// Requires JWT authentication
// Returns: gin.HandlerFunc that processes the get user request
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestedUserId := c.Param("id")

		// Get JWT claims from context (set by auth middleware)
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Type assertion to extract claims data
		tokenClaims, ok := claims.(*helpers.Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
			return
		}

		tokenUserId := tokenClaims.UserID
		userType := tokenClaims.Role

		// Authorization check: Users can only view their own profile, ADMIN can view all
		if userType != "ADMIN" && tokenUserId != requestedUserId {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		// Find user in database by user_id
		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": requestedUserId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// Return user data
		c.JSON(http.StatusOK, user)
	}
}

// GetUsers retrieves all users from the database
// This endpoint is restricted to ADMIN users only
// Requires JWT authentication
// Returns: gin.HandlerFunc that processes the get all users request
func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve JWT claims from context
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Type assertion to extract claims data
		tokenClaims, ok := claims.(*helpers.Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
			return
		}

		// Authorization check: Only ADMIN users can view all users
		if tokenClaims.Role != "ADMIN" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		// Retrieve all users from the database
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		// Find all users
		cursor, err := userCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(ctx)

		var users []models.User
		if err := cursor.All(ctx, &users); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return the list of users
		c.JSON(http.StatusOK, users)
	}
}
