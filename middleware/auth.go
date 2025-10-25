package middleware

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// CustomClaims represents the claims in our custom JWT tokens
type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func AuthMiddleware(client *auth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			fmt.Println("ERROR: Authorization header is empty")
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header required"})
			return
		}

		idToken := strings.Replace(authHeader, "Bearer ", "", 1)

		// Try Firebase token verification first
		firebaseToken, err := client.VerifyIDToken(context.Background(), idToken)
		if err == nil {
			// Firebase token verification successful
			c.Set("user_id", firebaseToken.UID)
			c.Next()
			return
		}

		// For Firebase custom tokens, we need to decode them to get the user ID
		// Firebase custom tokens are JWT tokens with a specific structure
		parts := strings.Split(idToken, ".")
		if len(parts) != 3 {
			fmt.Printf("ERROR: Invalid token format - not a valid JWT\n")
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token format"})
			return
		}

		// Decode the payload (second part of JWT)
		payload := parts[1]
		// Add padding if needed
		if len(payload)%4 != 0 {
			payload += strings.Repeat("=", 4-len(payload)%4)
		}

		// Decode base64 payload
		decoded, err := base64.URLEncoding.DecodeString(payload)
		if err != nil {
			fmt.Printf("ERROR: Failed to decode JWT payload: %v\n", err)
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token payload"})
			return
		}

		// Parse the JSON payload to extract user ID
		var claims map[string]interface{}
		if err := json.Unmarshal(decoded, &claims); err != nil {
			fmt.Printf("ERROR: Failed to parse JWT claims: %v\n", err)
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token claims"})
			return
		}

		// Extract user ID from claims
		userID, ok := claims["uid"].(string)
		if !ok {
			fmt.Printf("ERROR: No uid found in token claims: %v\n", claims)
			c.AbortWithStatusJSON(401, gin.H{"error": "No user ID in token"})
			return
		}

		c.Set("user_id", userID)
		c.Next()
		return
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
