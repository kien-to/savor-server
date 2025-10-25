package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"savor-server/db"
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
		fmt.Println("SignUp: Starting signup request from %s\n", c.ClientIP())

		var input SignUpInput

		if err := c.ShouldBindJSON(&input); err != nil {
			fmt.Println("SignUp: JSON binding error: %v", err)
			fmt.Println("SignUp: Request body binding failed - invalid JSON format or missing required fields")
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request format: %v", err.Error())})
			return
		}

		fmt.Println("SignUp: Successfully parsed request - Email: %s, Password length: %d", input.Email, len(input.Password))

		// Validate input
		if input.Email == "" {
			fmt.Println("SignUp: Validation failed - empty email")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
			return
		}

		if input.Password == "" {
			fmt.Println("SignUp: Validation failed - empty password")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password is required"})
			return
		}

		if len(input.Password) < 6 {
			fmt.Println("SignUp: Validation failed - password too short: %d characters", len(input.Password))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters long"})
			return
		}

		fmt.Println("SignUp: Input validation passed, initializing Firebase Auth client")

		client, err := app.Auth(context.Background())
		if err != nil {
			fmt.Println("SignUp: Failed to initialize Firebase Auth client: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize auth client"})
			return
		}

		fmt.Println("SignUp: Firebase Auth client initialized, creating user")

		params := (&auth.UserToCreate{}).
			Email(input.Email).
			Password(input.Password)

		user, err := client.CreateUser(context.Background(), params)
		if err != nil {
			fmt.Printf("SignUp: Firebase CreateUser failed: %v", err)
			fmt.Printf("SignUp: Error type: %T", err)

			// Provide more specific error messages
			errorMsg := err.Error()
			if strings.Contains(errorMsg, "email-already-exists") {
				errorMsg = "Email address is already registered"
			} else if strings.Contains(errorMsg, "invalid-email") {
				errorMsg = "Invalid email address format"
			} else if strings.Contains(errorMsg, "weak-password") {
				errorMsg = "Password is too weak"
			}

			c.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
			return
		}

		// Generate custom token
		token, err := client.CustomToken(context.Background(), user.UID)
		if err != nil {
			fmt.Printf("SignUp: Failed to generate custom token for user %s: %v", user.UID, err)
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
	log.Printf("GetProfile: Fetching profile for user_id: %s", userId)

	// Ensure table exists
	if _, err := db.DB.Exec(`
        CREATE TABLE IF NOT EXISTS user_profiles (
            user_id TEXT PRIMARY KEY,
            name TEXT,
            phone TEXT,
            email TEXT,
            created_at TIMESTAMPTZ DEFAULT NOW(),
            updated_at TIMESTAMPTZ DEFAULT NOW()
        )`); err != nil {
		log.Printf("ERROR: Failed ensuring user_profiles table: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var name, phone, email string
	err := db.DB.QueryRow(`
        SELECT COALESCE(name,''), COALESCE(phone,''), COALESCE(email,'')
        FROM user_profiles WHERE user_id = $1
    `, userId).Scan(&name, &phone, &email)

	if err != nil {
		log.Printf("GetProfile: No profile found for user_id %s (error: %v), returning empty profile", userId, err)
		// Return minimal profile if not found
		c.JSON(http.StatusOK, gin.H{
			"user_id": userId,
			"name":    "",
			"phone":   "",
			"email":   "",
		})
		return
	}

	log.Printf("GetProfile: Found profile for user_id %s - name: '%s', phone: '%s', email: '%s'", userId, name, phone, email)

	c.JSON(http.StatusOK, gin.H{
		"user_id": userId,
		"name":    name,
		"phone":   phone,
		"email":   email,
	})
}

// UpdateProfileInput holds editable profile fields
type UpdateProfileInput struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

// UpdateProfile stores profile fields (session-backed for now)
func UpdateProfile(c *gin.Context) {
	userId := c.GetString("user_id")
	if userId == "" {
		log.Printf("UpdateProfile: Unauthorized - no user_id in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	log.Printf("UpdateProfile: Starting update for user_id: %s", userId)

	var input UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("UpdateProfile: Invalid input JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	log.Printf("UpdateProfile: Received input - Name: '%s', Phone: '%s', Email: '%s'", input.Name, input.Phone, input.Email)

	// Ensure table exists
	if _, err := db.DB.Exec(`
        CREATE TABLE IF NOT EXISTS user_profiles (
            user_id TEXT PRIMARY KEY,
            name TEXT,
            phone TEXT,
            email TEXT,
            created_at TIMESTAMPTZ DEFAULT NOW(),
            updated_at TIMESTAMPTZ DEFAULT NOW()
        )`); err != nil {
		log.Printf("ERROR: Failed ensuring user_profiles table: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	name := strings.TrimSpace(input.Name)
	phone := strings.TrimSpace(input.Phone)
	email := strings.TrimSpace(input.Email)

	log.Printf("UpdateProfile: Upserting to database - user_id: %s, name: '%s', phone: '%s', email: '%s'", userId, name, phone, email)

	// Upsert profile
	if _, err := db.DB.Exec(`
        INSERT INTO user_profiles (user_id, name, phone, email)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (user_id) DO UPDATE SET
            name = EXCLUDED.name,
            phone = EXCLUDED.phone,
            email = EXCLUDED.email,
            updated_at = NOW()
    `, userId, name, phone, email); err != nil {
		log.Printf("ERROR: Upserting user profile failed for user_id %s: %v", userId, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save profile"})
		return
	}

	log.Printf("UpdateProfile: Successfully saved profile for user_id: %s", userId)

	c.JSON(http.StatusOK, gin.H{
		"user_id": userId,
		"name":    name,
		"phone":   phone,
		"email":   email,
	})
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
		log.Printf("Login: Starting login request from %s", c.ClientIP())

		var input LoginInput
		if err := c.ShouldBindJSON(&input); err != nil {
			log.Printf("Login: JSON binding error: %v", err)
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: fmt.Sprintf("Invalid request format: %v", err.Error())})
			return
		}

		log.Printf("Login: Successfully parsed request - Email: %s, Password length: %d", input.Email, len(input.Password))

		// Verify email and password using Firebase Auth REST API
		apiKey := os.Getenv("FIREBASE_API_KEY")
		if apiKey == "" {
			log.Printf("Login: FIREBASE_API_KEY environment variable not set")
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Firebase API key not configured"})
			return
		}

		log.Printf("Login: Firebase API key found, proceeding with authentication")

		// Use Firebase Auth REST API to verify email/password
		url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=%s", apiKey)

		payload := map[string]interface{}{
			"email":             input.Email,
			"password":          input.Password,
			"returnSecureToken": true,
		}

		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Login: Failed to marshal request payload: %v", err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to process request"})
			return
		}

		log.Printf("Login: Making request to Firebase Auth REST API")
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
		if err != nil {
			log.Printf("Login: HTTP request to Firebase failed: %v", err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to authenticate"})
			return
		}
		defer resp.Body.Close()

		log.Printf("Login: Firebase Auth API response status: %d", resp.StatusCode)

		if resp.StatusCode != http.StatusOK {
			var firebaseError struct {
				Error struct {
					Message string `json:"message"`
				} `json:"error"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&firebaseError); err != nil {
				log.Printf("Login: Failed to decode Firebase error response: %v", err)
				c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid credentials"})
				return
			}

			log.Printf("Login: Firebase authentication failed: %s", firebaseError.Error.Message)

			// Map Firebase error messages to user-friendly messages
			errorMessage := "Invalid credentials"
			switch firebaseError.Error.Message {
			case "EMAIL_NOT_FOUND":
				errorMessage = "Invalid credentials"
			case "INVALID_PASSWORD":
				errorMessage = "Invalid credentials"
			case "USER_DISABLED":
				errorMessage = "Account has been disabled"
			default:
				errorMessage = "Invalid credentials"
			}

			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: errorMessage})
			return
		}

		// Parse the successful response to get the user ID
		var authResponse struct {
			LocalID string `json:"localId"`
			IDToken string `json:"idToken"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
			log.Printf("Login: Failed to decode Firebase success response: %v", err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to process authentication response"})
			return
		}

		log.Printf("Login: Firebase authentication successful for user: %s", authResponse.LocalID)

		// Generate custom token using the verified user ID
		client, err := app.Auth(context.Background())
		if err != nil {
			log.Printf("Login: Failed to initialize Firebase Auth client: %v", err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to initialize auth client"})
			return
		}

		customToken, err := client.CustomToken(context.Background(), authResponse.LocalID)
		if err != nil {
			log.Printf("Login: Failed to generate custom token for user %s: %v", authResponse.LocalID, err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to generate token"})
			return
		}

		log.Printf("Login: Custom token generated successfully for user %s", authResponse.LocalID)
		c.JSON(http.StatusOK, models.AuthResponse{UserID: authResponse.LocalID, Token: customToken})
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

		apiKey := os.Getenv("FIREBASE_API_KEY")
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
