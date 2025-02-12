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
	connStr := "postgres://savor_user:your_password@localhost:5432/savor?sslmode=disable"
	database, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Test the connection
	err = database.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Assign to your global db variable
	db.DB = database

	// Initialize Stripe
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")
	if stripeKey == "" {
		log.Fatal("STRIPE_SECRET_KEY is required")
	}
	config.InitializeStripe(stripeKey)

	// Initialize Gin router with debug mode
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.Use(gin.Logger(), gin.Recovery())

	// Initialize session store
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("savor_session", store))

	// Add CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
	// protected.Use(middleware.AuthMiddleware(authClient))
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
		// homeGroup.GET("/stores/favorites", handlers.GetFavorites)
	}

	storesGroup := r.Group("/api/stores")
	{
		storesGroup.GET("/:id", handlers.GetStoreDetail)
		storesGroup.POST("/:id/toggle-save", middleware.AuthMiddleware(authClient), handlers.ToggleSaveStore)
		storesGroup.GET("/favorites", middleware.AuthMiddleware(authClient), handlers.GetFavorites)
	}

	paymentGroup := r.Group("/api/payment")
	{
		paymentGroup.POST("/create-intent", middleware.AuthMiddleware(authClient), handlers.CreateReservation)
		paymentGroup.POST("/confirm", middleware.AuthMiddleware(authClient), handlers.ConfirmReservation)
		paymentGroup.POST("/confirm-pay-at-store", middleware.AuthMiddleware(authClient), handlers.ConfirmPayAtStore)
	}

	reservationsGroup := r.Group("/api/reservations")
	{
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

	r.Run(":8080")
}
