package main

import (
	// "context"
	"context"
	"log"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/joho/godotenv"

	"savor-server/config"
	_ "savor-server/docs" // This will be auto-generated
	"savor-server/handlers"
	"savor-server/middleware"
	"savor-server/services"

	"savor-server/db" // Add this import

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Auth Service API
// @version         1.0
// @description     A Firebase authentication service with social login support.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// http://localhost:8080/swagger/index.html#/auth/post_auth_signup

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load .env file before any other setup
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Set up logging
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Initialize Firebase
	app, err := config.InitializeFirebase()
	if err != nil {
		log.Fatalf("Error initializing Firebase: %v\n", err)
	}

	// Create Firebase Auth client
	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("Error creating Firebase Auth client: %v\n", err)
	}

	// Initialize database connection
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		// Fallback to local database for development
		// connStr = "postgresql://postgres:PttZOUsZxyCRGDIlgKecSGThePtJbNYB@crossover.proxy.rlwy.net:31363/railway"
		connStr = "postgres://savor_user:your_password@localhost:5432/savor?sslmode=disable"
		log.Printf("Using local database connection: %s", connStr)
	} else {
		log.Printf("Using Railway database connection")
		// Don't log the full connection string for security
		log.Printf("DATABASE_URL is set and will be used")
	}

	log.Printf("Attempting to connect to database...")
	database, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Test the connection
	log.Printf("Testing database connection...")
	err = database.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Printf("Database connection successful!")

	// Assign to your global db variable
	db.DB = database

	// Initialize Stripe
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")
	if stripeKey == "" {
		log.Fatal("STRIPE_SECRET_KEY is required")
	}
	config.InitializeStripe(stripeKey)

	// Initialize Google Maps
	services.InitializeGoogleMaps()

	// Initialize Gin router with appropriate mode
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = "debug" // Default to debug for local development
	}
	gin.SetMode(ginMode)
	r := gin.Default()
	r.Use(gin.Logger(), gin.Recovery())

	// Initialize session store
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("savor_session", store))

	// Add CORS middleware
	allowedOrigins := []string{"http://localhost:3000"}
	if frontendURL := os.Getenv("FRONTEND_URL"); frontendURL != "" {
		allowedOrigins = append(allowedOrigins, frontendURL)
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check endpoint for Railway
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "healthy",
			"database": "connected",
		})
	})

	// Auth routes
	auth := r.Group("/auth")
	{
		auth.POST("/signup", handlers.SignUp(app))
		auth.POST("/google", handlers.GoogleAuth(app))
		auth.POST("/facebook", handlers.FacebookAuth(app))
		auth.POST("/login", handlers.Login(app))
		auth.POST("/phone", handlers.PhoneAuth(app))
		auth.POST("/forgot-password", handlers.ForgotPassword(app))
	}

	// Protected routes example
	protected := r.Group("/api/settings")
	protected.Use(middleware.AuthMiddleware(authClient))
	{
		protected.GET("/profile", handlers.GetProfile)
	}

	// Home routes
	homeGroup := r.Group("/api/home")
	{
		homeGroup.GET("", handlers.GetHomePageData)
		homeGroup.GET("/search", handlers.SearchStores)
		homeGroup.POST("/stores/:id/save", handlers.SaveStore)
		homeGroup.POST("/stores/:id/unsave", handlers.UnsaveStore)
		homeGroup.GET("/stores/favorites", handlers.GetFavorites)
	}

	storesGroup := r.Group("/api/stores")
	{
		storesGroup.GET("/:id", handlers.GetStoreDetail)
		storesGroup.POST("/:id/toggle-save", middleware.AuthMiddleware(authClient), handlers.ToggleSaveStore)
		storesGroup.GET("/favorites", middleware.AuthMiddleware(authClient), handlers.GetFavorites)
	}

	// Maps routes
	mapsGroup := r.Group("/api/maps")
	{
		mapsGroup.GET("/distance", handlers.CalculateDistance)
		mapsGroup.GET("/directions", handlers.GetDirections)
		mapsGroup.GET("/stores/:storeId", handlers.GetStoreWithDistance)
	}

	paymentGroup := r.Group("/api/payment")
	{
		paymentGroup.POST("/create-intent", middleware.AuthMiddleware(authClient), handlers.CreateReservation)
		paymentGroup.POST("/confirm", middleware.AuthMiddleware(authClient), handlers.ConfirmReservation)
		paymentGroup.POST("/confirm-pay-at-store", middleware.AuthMiddleware(authClient), handlers.ConfirmPayAtStore)
	}

	reservationsGroup := r.Group("/api/reservations")
	{
		// reservationsGroup.GET("", handlers.GetReservations)
		reservationsGroup.GET("", middleware.AuthMiddleware(authClient), handlers.GetUserReservations)
		reservationsGroup.GET("/demo", handlers.GetDemoReservations)
		reservationsGroup.GET("/session", handlers.GetSessionReservations)
		reservationsGroup.GET("/guest", handlers.GetGuestReservations)
		reservationsGroup.POST("/guest", handlers.CreateGuestReservation)
		reservationsGroup.DELETE("/:id", handlers.DeleteReservation)
	}

	storeManagementGroup := r.Group("/api/store-management")
	storeManagementGroup.Use(middleware.AuthMiddleware(authClient))
	{
		storeManagementGroup.POST("/create", handlers.CreateStore)
		storeManagementGroup.GET("/my-store", handlers.GetMyStore)
		storeManagementGroup.PUT("/update", handlers.UpdateStore)
		storeManagementGroup.POST("/toggle-selling", handlers.ToggleStoreSelling)
		storeManagementGroup.POST("/bag-details", handlers.UpdateBagDetails)
		storeManagementGroup.POST("/pickup-schedule", handlers.UpdatePickupSchedule)
	}

	// Start server with port from environment variable (Railway) or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s", port)
	r.Run(":" + port)
}
