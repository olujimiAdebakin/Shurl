package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/olujimiAdebakin/Shurl/controllers"
	"github.com/olujimiAdebakin/Shurl/dtos"
	"github.com/olujimiAdebakin/Shurl/initializers"
	"github.com/olujimiAdebakin/Shurl/models"
)

func RequireAuthWithToken(c *gin.Context) {
	// Get token from Authorization header
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, dtos.ErrorResponse{
			Success: false,
			Error:   "Unauthorized - No token provided",
		})
		c.Abort()
		return
	}

	// Extract token from "Bearer <token>" format
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	if tokenString == authHeader {
		c.JSON(http.StatusUnauthorized, dtos.ErrorResponse{
			Success: false,
			Error:   "Unauthorized - Invalid token format. Use 'Bearer <token>'",
		})
		c.Abort()
		return
	}

	// Decode/Validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, dtos.ErrorResponse{
			Success: false,
			Error:   "Unauthorized - Invalid token",
		})
		c.Abort()
		return
	}

	// Extract claims and validate
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.JSON(http.StatusUnauthorized, dtos.ErrorResponse{
				Success: false,
				Error:   "Unauthorized - Token expired",
			})
			c.Abort()
			return
		}

		// Find user with token sub
		var user models.User
		initializers.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.JSON(http.StatusUnauthorized, dtos.ErrorResponse{
				Success: false,
				Error:   "Unauthorized - User not found",
			})
			c.Abort()
			return
		}

		// Convert user model to context struct
		contextUser := controllers.ContextUserStruct{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
			Role:  user.Role,
		}

		// Attach user to context
		c.Set("user", contextUser)

		// Continue to next handler
		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, dtos.ErrorResponse{
			Success: false,
			Error:   "Unauthorized - Invalid token claims",
		})
		c.Abort()
		return
	}
}

func RequireAuthWithCookie(c *gin.Context) {
	// Get cookie from request
	tokenString, err := c.Cookie("Authorization")

	if err != nil {
		c.JSON(http.StatusUnauthorized, dtos.ErrorResponse{
			Success: false,
			Error:   "Unauthorized - No token provided",
		})
		c.Abort()
		return
	}

	// Decode/Validate token (same as above)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, dtos.ErrorResponse{
			Success: false,
			Error:   "Unauthorized - Invalid token",
		})
		c.Abort()
		return
	}

	// Rest of validation (same as token-based)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.JSON(http.StatusUnauthorized, dtos.ErrorResponse{
				Success: false,
				Error:   "Unauthorized - Token expired",
			})
			c.Abort()
			return
		}

		var user models.User
		initializers.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.JSON(http.StatusUnauthorized, dtos.ErrorResponse{
				Success: false,
				Error:   "Unauthorized - User not found",
			})
			c.Abort()
			return
		}

		// Convert user model to context struct
		contextUser := controllers.ContextUserStruct{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
			Role:  user.Role,
		}

		c.Set("user", contextUser)
		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, dtos.ErrorResponse{
			Success: false,
			Error:   "Unauthorized - Invalid token claims",
		})
		c.Abort()
		return
	}
}

