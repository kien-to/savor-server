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
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			fmt.Println("Authorization header is empty")
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header required"})
			return
		}

		idToken := strings.Replace(authHeader, "Bearer ", "", 1)
		token, err := client.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			fmt.Println("Error verifying token:", err)
			fmt.Println("Token:", idToken)
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("user_id", token.UID)
		c.Next()
	}
}
