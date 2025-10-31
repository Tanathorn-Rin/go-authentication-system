package helpers

import (
	"authentication/config"
	"context"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// Claims represents the JWT token claims structure
// It includes user identification, email, role, and standard JWT claims
type Claims struct {
	UserID string `json:"user_id"` // Unique identifier for the user
	Email  string `json:"email"`   // User's email address
	Role   string `json:"role"`    // User's role (ADMIN or USER)

	jwt.StandardClaims // Standard JWT claims (expiration, issuer, etc.)
}

// jwtKey is the secret key used for signing and validating JWT tokens
var jwtKey []byte

// SetJWTKey sets the JWT secret key from a string
// This key is used to sign and validate all JWT tokens
func SetJWTKey(key string) {
	jwtKey = []byte(key)
}

// GetJWTKey returns the current JWT secret key
// Returns: The JWT key as a byte slice
func GetJWTKey() []byte {
	return []byte(jwtKey)
}

// ValidateToken parses and validates a JWT token string
// It checks the token signature and expiration, then returns the claims
// Parameters:
//   - tokenString: The JWT token string to validate
//
// Returns:
//   - *Claims: The parsed claims from the token if valid
//   - error: Any error encountered during validation
func ValidateToken(tokenString string) (*Claims, error) {
	// Use the dynamically set JWT key here
	secretKey := GetJWTKey() // This retrieves the key set in SetJWTKey

	// Parse the token with claims
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	// Check if the token is valid and return the claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GenerateTokens creates both an access token and a refresh token for a user
// The access token expires in 24 hours, refresh token in 7 days
// Parameters:
//   - email: User's email address
//   - userID: Unique user identifier
//   - userType: User's role (ADMIN or USER)
//
// Returns:
//   - string: Signed access token
//   - string: Signed refresh token
func GenerateTokens(email, userID, userType string) (string, string) {
	log.Printf("JWT Key %v Type: %T", jwtKey, jwtKey)

	// Token expiration times
	tokenExpiry := time.Now().Add(24 * time.Hour).Unix()            // Access token valid for 24 hours
	refreshTokenExpiry := time.Now().Add(7 * 24 * time.Hour).Unix() // Refresh token valid for 7 days

	// Create claims for access token with user information
	claims := &Claims{
		Email:  email,
		UserID: userID,
		Role:   userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpiry,
		},
	}

	// Create claims for refresh token (only expiration, no user data)
	refreshClaims := &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshTokenExpiry,
		},
	}

	// Generate and sign the access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedAcessToken, err := accessToken.SignedString(jwtKey)
	if err != nil {
		panic(err)
	}

	// Generate and sign the refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		panic(err)
	}

	return signedAcessToken, signedRefreshToken
}

// HashPassword hashes a plain text password using bcrypt
// Uses bcrypt's default cost factor for security
// Parameters:
//   - password: Plain text password to hash
//
// Returns:
//   - string: Hashed password
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

// UpdateAllTokens updates the user's access and refresh tokens in the database
// Also updates the updated_at timestamp
// Parameters:
//   - signedToken: New access token to store
//   - signedRefreshToken: New refresh token to store
//   - userID: ID of the user to update
//
// Returns:
//   - error: Any error encountered during the database update
func UpdateAllTokens(signedToken, signedRefreshToken, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	userCollection := config.OpenCollection("users")

	// Create an update object to set new tokens and timestamp
	updateObj := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "token", Value: signedToken},
			bson.E{Key: "refresh_token", Value: signedRefreshToken},
			bson.E{Key: "updated_at", Value: time.Now()},
		}},
	}

	// Create a filter to find the user by ID
	filter := bson.M{"user_id": userID}

	// Update the user document in MongoDB
	_, err := userCollection.UpdateOne(ctx, filter, updateObj)

	return err
}

// VerifyPassword compares a hashed password with a plain text password
// Uses bcrypt to securely compare passwords
// Parameters:
//   - foundPwd: The hashed password from the database
//   - pwd: The plain text password to verify
//
// Returns:
//   - bool: true if passwords match, false otherwise
//   - error: Any error from bcrypt comparison (nil if passwords match)
func VerifyPassword(foundPwd, pwd string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(foundPwd), []byte(pwd))

	return err == nil, err
}
