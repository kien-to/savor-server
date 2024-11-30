package main

import (
	"log"

	"savor-backend-2/config"
	_ "savor-backend-2/docs" // This will be auto-generated
	"savor-backend-2/handlers"
	"savor-backend-2/middleware"

	"github.com/gin-gonic/gin"
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

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Initialize Firebase
	app, err := config.InitializeFirebase()
	if err != nil {
		log.Fatalf("Error initializing Firebase: %v\n", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Auth routes
	auth := r.Group("/auth")
	{
		auth.POST("/signup", handlers.SignUp(app))
		auth.POST("/google", handlers.GoogleAuth(app))
		auth.POST("/facebook", handlers.FacebookAuth(app))
		auth.POST("/login", handlers.Login(app))
	}

	// Protected routes example
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(app))
	{
		protected.GET("/profile", handlers.GetProfile)
	}

	r.Run(":8080")
}
