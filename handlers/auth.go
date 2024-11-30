package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"savor-backend-2/models"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

type SignUpInput struct {
	Email    string `json:"email" binding:"required" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

type SocialAuthInput struct {
	IdToken string `json:"id_token" binding:"required" example:"eyJhbGciOiJS..."`
}

type FacebookAuthInput struct {
	AccessToken string `json:"access_token" binding:"required" example:"EAAaYA6ZA..."`
}

// LoginInput represents the login request body
type LoginInput struct {
	Email    string `json:"email" binding:"required" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// @Summary     Sign up a new user
// @Description Register a new user with email and password
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       input body SignUpInput true "Sign Up Credentials"
// @Success     200 {object} map[string]string "Successfully created user"
// @Failure     400 {object} map[string]string "Invalid input"
// @Failure     500 {object} map[string]string "Server error"
// @Router      /auth/signup [post]
func SignUp(app *firebase.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input SignUpInput

		if err := c.ShouldBindJSON(&input); err != nil {
			fmt.Println(err)
			fmt.Println("something went wrong")
			log.Printf("Firebase error: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		client, err := app.Auth(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize auth client"})
			return
		}

		params := (&auth.UserToCreate{}).
			Email(input.Email).
			Password(input.Password)

		user, err := client.CreateUser(context.Background(), params)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user_id": user.UID})
	}
}

// @Summary     Authenticate with Google
// @Description Authenticate user using Google ID token
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       input body SocialAuthInput true "Google ID Token"
// @Success     200 {object} map[string]string "Successfully authenticated"
// @Failure     400 {object} map[string]string "Invalid input"
// @Failure     401 {object} map[string]string "Invalid token"
// @Failure     500 {object} map[string]string "Server error"
// @Router      /auth/google [post]
func GoogleAuth(app *firebase.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input SocialAuthInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		client, err := app.Auth(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize auth client"})
			return
		}

		token, err := client.VerifyIDToken(context.Background(), input.IdToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user_id": token.UID})
	}
}

// @Summary     Authenticate with Facebook
// @Description Authenticate user using Facebook access token
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       input body FacebookAuthInput true "Facebook Access Token"
// @Success     200 {object} map[string]string "Successfully authenticated"
// @Failure     400 {object} map[string]string "Invalid input"
// @Failure     401 {object} map[string]string "Invalid token"
// @Failure     500 {object} map[string]string "Server error"
// @Router      /auth/facebook [post]
func FacebookAuth(app *firebase.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input FacebookAuthInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		client, err := app.Auth(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize auth client"})
			return
		}

		token, err := client.VerifyIDToken(context.Background(), input.AccessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user_id": token.UID})
	}
}

// @Summary     Get user profile
// @Description Get authenticated user's profile
// @Tags        profile
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} map[string]string "User profile"
// @Failure     401 {object} map[string]string "Unauthorized"
// @Failure     500 {object} map[string]string "Server error"
// @Router      /api/profile [get]
func GetProfile(c *gin.Context) {
	userId := c.GetString("user_id")
	c.JSON(http.StatusOK, gin.H{"user_id": userId})
}

// @Summary     Login user
// @Description Authenticate user with email and password
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       input body LoginInput true "Login Credentials"
// @Success     200 {object} models.AuthResponse
// @Failure     400 {object} models.ErrorResponse
// @Failure     500 {object} models.ErrorResponse
// @Router      /auth/login [post]
func Login(app *firebase.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input LoginInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		client, err := app.Auth(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to initialize auth client"})
			return
		}

		user, err := client.GetUserByEmail(context.Background(), input.Email)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid credentials"})
			return
		}

		c.JSON(http.StatusOK, models.AuthResponse{UserID: user.UID})
	}
} 