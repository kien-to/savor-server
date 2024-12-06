package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"savor-server/models"

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

type PhoneAuthInput struct {
	PhoneNumber string `json:"phone_number" binding:"required" example:"+1234567890"`
	Code        string `json:"code" binding:"required" example:"123456"`
}

type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required" example:"user@example.com"`
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

		// Generate custom token
		token, err := client.CustomToken(context.Background(), user.UID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user_id": user.UID, "token": token})
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
			log.Printf("Failed to bind JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		client, err := app.Auth(context.Background())
		if err != nil {
			log.Printf("Failed to initialize Firebase auth client: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize auth client"})
			return
		}

		// Add debug logging for token verification
		log.Printf("Attempting to verify Google ID token length: %d", len(input.IdToken))
		token, err := client.VerifyIDToken(context.Background(), input.IdToken)
		if err != nil {
			log.Printf("Failed to verify Google ID token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid token: %v", err)})
			return
		}

		customToken, err := client.CustomToken(context.Background(), token.UID)
		if err != nil {
			log.Printf("Failed to generate custom token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user_id": token.UID, "token": customToken})
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

		// After verifying ID token, generate custom token
		customToken, err := client.CustomToken(context.Background(), token.UID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user_id": token.UID, "token": customToken})
	}
}

// @Summary     Authenticate with Phone
// @Description Authenticate user using phone number and verification code
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       input body PhoneAuthInput true "Phone Auth Input"
// @Success     200 {object} map[string]string "Successfully authenticated"
// @Failure     400 {object} map[string]string "Invalid input"
// @Failure     401 {object} map[string]string "Invalid code"
// @Router      /auth/phone [post]
func PhoneAuth(app *firebase.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input PhoneAuthInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		client, err := app.Auth(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize auth client"})
			return
		}

		token, err := client.VerifyIDToken(context.Background(), input.Code)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid verification code"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user_id": token.UID,
			"token":   input.Code,
		})
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

		// Generate custom token
		token, err := client.CustomToken(context.Background(), user.UID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, models.AuthResponse{UserID: user.UID, Token: token})
	}
}

// @Summary     Forgot Password
// @Description Send password reset email to user
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       input body ForgotPasswordInput true "Email Address"
// @Success     200 {object} map[string]string "Reset email sent successfully"
// @Failure     400 {object} map[string]string "Invalid input"
// @Failure     500 {object} map[string]string "Server error"
// @Router      /auth/forgot-password [post]
func ForgotPassword(app *firebase.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input ForgotPasswordInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		apiKey := "AIzaSyBcF4jUafzDMU7oAjWBlNLJARr282r9Duo"
		// os.Getenv("FIREBASE_API_KEY")
		if apiKey == "" {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Firebase API key not configured",
			})
			return
		}

		url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:sendOobCode?key=%s", apiKey)

		payload := map[string]interface{}{
			"requestType": "PASSWORD_RESET",
			"email":       input.Email,
		}

		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Failed to process request",
			})
			return
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Failed to send reset email",
			})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			var firebaseError struct {
				Error struct {
					Message string `json:"message"`
				} `json:"error"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&firebaseError); err != nil {
				c.JSON(http.StatusBadRequest, models.ErrorResponse{
					Error: "Failed to process Firebase response",
				})
				return
			}
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: firebaseError.Error.Message,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Password reset email sent successfully",
		})
	}
}
