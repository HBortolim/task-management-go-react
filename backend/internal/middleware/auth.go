package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// JwtMiddleware represents the JWT authentication middleware
type JwtMiddleware struct {
	jwtSecret string
}

// TokenClaims represents the JWT token claims
type TokenClaims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

// NewJwtMiddleware creates a new JWT middleware
func NewJwtMiddleware(jwtSecret string) *JwtMiddleware {
	return &JwtMiddleware{
		jwtSecret: jwtSecret,
	}
}

// GenerateToken generates a new JWT token
func (m *JwtMiddleware) GenerateToken(userID primitive.ObjectID, expiryHours int) (string, error) {
	// Set expiration time
	expirationTime := time.Now().Add(time.Duration(expiryHours) * time.Hour)

	// Create claims
	claims := &TokenClaims{
		UserID: userID.Hex(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.jwtSecret))

	return tokenString, err
}

// AuthRequired is a middleware to verify JWT token
func (m *JwtMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims := &TokenClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(m.jwtSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Convert user ID string to ObjectID
		userID, err := primitive.ObjectIDFromHex(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			c.Abort()
			return
		}

		// Set user ID in context
		c.Set("userId", userID)
		c.Next()
	}
}
