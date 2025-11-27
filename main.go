package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/olujimiAdebakin/Shurl/controllers"
	_ "github.com/olujimiAdebakin/Shurl/docs"
	"github.com/olujimiAdebakin/Shurl/initializers"
	"github.com/olujimiAdebakin/Shurl/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Shurl API
// @version 1.0
// @description A modern URL shortener API built with Go and Gin
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@shurl.dev
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @basePath /api/v1
// @schemes http https
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description JWT token with 'Bearer ' prefix

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	fmt.Println("Welcome to Shurl API!!!!")

	// Initialize router
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	// Health check endpoint (define before wildcards to avoid conflicts)
	// @Summary Health Check
	// @Description Check if the API is running
	// @Tags System
	// @Produce json
	// @Success 200 {object} map[string]string "status: ok, message: API is running"
	// @Router /health [get]
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "API is running",
		})
	})

	// Swagger UI endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 group
	v1 := router.Group("/api/v1")

	// User routes
	users := v1.Group("/users")
	{
		// @Summary Sign Up
		// @Description Create a new user account
		// @Tags Authentication
		// @Accept json
		// @Produce json
		// @Param request body dtos.CreateUserRequest true "Signup credentials"
		// @Success 201 {object} map[string]interface{} "User created with token"
		// @Failure 400 {object} map[string]interface{} "Bad request"
		// @Failure 409 {object} map[string]interface{} "Email already exists"
		// @Router /users/signup [post]
		users.POST("/signup", controllers.SignUpWithToken)

		// @Summary Login
		// @Description Authenticate user and receive JWT token
		// @Tags Authentication
		// @Accept json
		// @Produce json
		// @Param request body dtos.LoginUserRequest true "Login credentials"
		// @Success 200 {object} dtos.LoginResponse "User authenticated with token"
		// @Failure 400 {object} map[string]interface{} "Bad request"
		// @Failure 401 {object} map[string]interface{} "Invalid credentials"
		// @Router /users/login [post]
		users.POST("/login", controllers.LoginWithToken)

		// @Summary Validate Token
		// @Description Verify JWT token and get user information
		// @Tags Authentication
		// @Security Bearer
		// @Accept json
		// @Produce json
		// @Success 200 {object} map[string]interface{} "User information from token"
		// @Failure 401 {object} map[string]interface{} "Invalid or missing token"
		// @Router /users/validate [get]
		users.GET("/validate", middleware.RequireAuthWithToken, controllers.Validate)
	}

	// Link routes
	links := v1.Group("/links")
	{
		// @Summary Get Link Info
		// @Description Retrieve link details and increment click count
		// @Tags Links
		// @Accept json
		// @Produce json
		// @Param shortCode path string true "Short code of the link"
		// @Success 200 {object} dtos.LinkResponse "Link information"
		// @Failure 404 {object} map[string]interface{} "Link not found"
		// @Router /links/{shortCode} [get]
		links.GET("/:shortCode", controllers.GetLink)

		// @Summary Create Link
		// @Description Create a new shortened URL
		// @Tags Links
		// @Security Bearer
		// @Accept json
		// @Produce json
		// @Param request body dtos.CreateLinkRequest true "Link creation details"
		// @Success 201 {object} dtos.LinkResponse "Link created"
		// @Failure 400 {object} map[string]interface{} "Bad request"
		// @Failure 401 {object} map[string]interface{} "Unauthorized"
		// @Failure 409 {object} map[string]interface{} "Short code already exists"
		// @Router /links [post]
		links.POST("", middleware.RequireAuthWithToken, controllers.CreateLink)

		// @Summary Update Link
		// @Description Modify an existing link (owner only)
		// @Tags Links
		// @Security Bearer
		// @Accept json
		// @Produce json
		// @Param shortCode path string true "Short code of the link"
		// @Param request body dtos.LinkUpdateRequest true "Link update details"
		// @Success 200 {object} dtos.LinkResponse "Link updated"
		// @Failure 400 {object} map[string]interface{} "Bad request"
		// @Failure 401 {object} map[string]interface{} "Unauthorized"
		// @Failure 403 {object} map[string]interface{} "Forbidden - not the owner"
		// @Failure 404 {object} map[string]interface{} "Link not found"
		// @Router /links/{shortCode} [patch]
		links.PATCH("/:shortCode", middleware.RequireAuthWithToken, controllers.UpdateLink)

		// @Summary Delete Link
		// @Description Delete a link (owner only)
		// @Tags Links
		// @Security Bearer
		// @Accept json
		// @Produce json
		// @Param shortCode path string true "Short code of the link"
		// @Success 200 {object} map[string]interface{} "Link deleted"
		// @Failure 401 {object} map[string]interface{} "Unauthorized"
		// @Failure 403 {object} map[string]interface{} "Forbidden - not the owner"
		// @Failure 404 {object} map[string]interface{} "Link not found"
		// @Router /links/{shortCode} [delete]
		links.DELETE("/:shortCode", middleware.RequireAuthWithToken, controllers.DeleteLink)

		// @Summary Get User Links
		// @Description Retrieve all links created by authenticated user
		// @Tags Links
		// @Security Bearer
		// @Accept json
		// @Produce json
		// @Success 200 {array} dtos.LinkResponse "User's links"
		// @Failure 401 {object} map[string]interface{} "Unauthorized"
		// @Router /links [get]
		links.GET("", middleware.RequireAuthWithToken, controllers.GetUserLinks)
	}

	// Redirect route - accessible at root level (e.g., localhost:8080/my-link)
	// IMPORTANT: This should be defined AFTER all other routes to avoid conflicts
	// @Summary Redirect to Link
	// @Description Redirect to original URL and increment clicks
	// @Tags Redirect
	// @Accept json
	// @Produce json
	// @Param shortCode path string true "Short code of the link"
	// @Success 301 "Redirect to original URL"
	// @Failure 404 {object} map[string]interface{} "Link not found"
	// @Router /{shortCode} [get]
	router.GET("/:shortCode", controllers.RedirectLink)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port) // listens on 0.0.0.0:PORT
}