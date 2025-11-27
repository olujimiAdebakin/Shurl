package controllers

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/olujimiAdebakin/Shurl/dtos"
	"github.com/olujimiAdebakin/Shurl/initializers"
	"github.com/olujimiAdebakin/Shurl/models"
	"golang.org/x/crypto/bcrypt"
)

type ContextUserStruct struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

// SignUpWithToken godoc
// @Summary Create a new user account
// @Description Create a new user account with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body dtos.CreateUserRequest true "Signup request"
// @Success 201 {object} dtos.SuccessResponse{data=dtos.LoginResponse}
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 409 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /users/signup [post]
func SignUpWithToken(c *gin.Context) {
	var req dtos.CreateUserRequest

	// Validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			Success: false,
			Error:   "Invalid input: " + err.Error(),
		})
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			Success: false,
			Error:   "Failed to hash the password",
		})
		return
	}

	// Create the user
	user := models.User{
		Email:    req.Email,
		Password: string(hash),
		Role:     models.RoleUser,
		Name:     req.Name,
	}

	if err := initializers.DB.Create(&user).Error; err != nil {
		// Check for duplicate email
		if strings.Contains(err.Error(), "duplicate") ||
		   strings.Contains(err.Error(), "unique") {
			c.JSON(http.StatusConflict, dtos.ErrorResponse{
				Success: false,
				Error:   "User with this email already exists",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			Success: false,
			Error:   "Failed to create user",
		})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			Success: false,
			Error:   "Failed to create token",
		})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, dtos.SuccessResponse{
		Success: true,
		Data: dtos.LoginResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Token: tokenString,
		},
	})
}


// LoginWithToken godoc
// @Summary Authenticate user and receive JWT token
// @Description Login with email and password to receive JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body dtos.LoginUserRequest true "Login request"
// @Success 200 {object} dtos.SuccessResponse{data=dtos.LoginResponse}
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 401 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /users/login [post]
func LoginWithToken(c *gin.Context) {
	var req dtos.LoginUserRequest

	// Validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			Success: false,
			Error:   "Invalid input: " + err.Error(),
		})
		return
	}

	// Look up user by email
	var user models.User
	result := initializers.DB.Where("email = ?", req.Email).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, dtos.ErrorResponse{
			Success: false,
			Error:   "Invalid email or password",
		})
		return
	}

	// Compare password with stored hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dtos.ErrorResponse{
			Success: false,
			Error:   "Invalid email or password",
		})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			Success: false,
			Error:   "Failed to create token",
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, dtos.SuccessResponse{
		Success: true,
		Data: dtos.LoginResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Token: tokenString,
		},
	})
}


// Validate godoc
// @Summary Verify JWT token and get user information
// @Description Validate JWT token and return authenticated user information
// @Tags Authentication
// @Security Bearer
// @Accept json
// @Produce json
// @Success 200 {object} dtos.SuccessResponse{data=ContextUserStruct}
// @Failure 401 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /users/validate [get]
// @SecurityScheme Bearer
func Validate(c *gin.Context) {
	// Retrieve user from context (set by RequireAuth middleware)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, dtos.ErrorResponse{
			Success: false,
			Error:   "User not found in context",
		})
		return
	}

	// Type assert to ContextUserStruct
	contextUser, ok := user.(ContextUserStruct)
	if !ok {
		c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			Success: false,
			Error:   "Invalid user data format",
		})
		return
	}

	// Return user information
	c.JSON(http.StatusOK, dtos.SuccessResponse{
		Success: true,
		Data:    contextUser,
	})
}
