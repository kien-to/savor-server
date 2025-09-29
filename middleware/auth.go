package middleware

import (
	"context"
	"fmt"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(client *auth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("DEBUG: AuthMiddleware called for %s %s\n", c.Request.Method, c.Request.URL.Path)

		authHeader := c.GetHeader("Authorization")
		fmt.Printf("DEBUG: Authorization header: %s\n", authHeader)

		if authHeader == "" {
			fmt.Println("ERROR: Authorization header is empty")
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header required"})
			return
		}

		idToken := strings.Replace(authHeader, "Bearer ", "", 1)
		fmt.Printf("DEBUG: Extracted ID token (first 20 chars): %s...\n", idToken[:min(20, len(idToken))])

		token, err := client.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			fmt.Printf("ERROR: Error verifying token: %v\n", err)
			fmt.Printf("DEBUG: Token (first 50 chars): %s...\n", idToken[:min(50, len(idToken))])
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}

		fmt.Printf("DEBUG: Token verification successful - UID: %s\n", token.UID)
		fmt.Printf("DEBUG: Token claims: %+v\n", token.Claims)

		c.Set("user_id", token.UID)
		fmt.Printf("DEBUG: Set user_id in context: %s\n", token.UID)
		c.Next()
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
