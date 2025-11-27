package controllers

import (
	"crypto/md5"
	_"log"
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olujimiAdebakin/Shurl/dtos"
	"github.com/olujimiAdebakin/Shurl/initializers"
	"github.com/olujimiAdebakin/Shurl/models"
)

// CreateLink godoc
// @Summary Create a new shortened URL
// @Description Create a new shortened URL with optional custom short code
// @Tags Links
// @Security Bearer
// @Accept json
// @Produce json
// @Param input body dtos.CreateLinkRequest true "Create link request"
// @Success 201 {object} dtos.SuccessResponse{data=dtos.LinkResponse}
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 401 {object} dtos.ErrorResponse
// @Failure 409 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /links [post]
func CreateLink(c *gin.Context) {
	var req dtos.CreateLinkRequest

	// Validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			Success: false,
			Error:   "Invalid input: " + err.Error(),
		})
		return
	}

	// Get authenticated user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, dtos.ErrorResponse{
			Success: false,
			Error:   "Unauthorized - User not found",
		})
		return
	}

	contextUser, ok := user.(ContextUserStruct)
	if !ok {
		c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			Success: false,
			Error:   "Invalid user data",
		})
		return
	}

	// Generate or use provided short code
	shortCode := req.ShortCode
	if shortCode == "" {
		shortCode = generateShortCode()
	}

	// Generate hash of the original URL
	hash := generateHash(req.OriginalURL)

	// Create link model
	link := models.Link{
		ShortCode:   shortCode,
		OriginalURL: req.OriginalURL,
		Hash:        hash,
		Clicks:      0,
		UserID:      contextUser.ID,
	}

	// Save link to database
	if err := initializers.DB.Create(&link).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			c.JSON(http.StatusConflict, dtos.ErrorResponse{
				Success: false,
				Error:   "Short code already exists",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			Success: false,
			Error:   "Failed to create link",
		})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, dtos.SuccessResponse{
		Success: true,
		Data: dtos.LinkResponse{
			ShortCode:   link.ShortCode,
			OriginalURL: link.OriginalURL,
			Clicks:      link.Clicks,
			Favicon:     link.Favicon,
			UserID:      link.UserID,
		},
	})
}

// GetLink godoc
// @Summary Get link information
// @Description Retrieve link details by short code and increment click count
// @Tags Links
// @Accept json
// @Produce json
// @Param shortCode path string true "Short code of the link"
// @Success 200 {object} dtos.SuccessResponse{data=dtos.LinkResponse}
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 404 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /links/{shortCode} [get]
func GetLink(c *gin.Context) {
	shortCode := c.Param("shortCode")

	if shortCode == "" {
		c.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			Success: false,
			Error:   "Short code is required",
		})
		return
	}

	var link models.Link
	result := initializers.DB.Where("short_code = ?", shortCode).First(&link)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, dtos.ErrorResponse{
			Success: false,
			Error:   "Link not found",
		})
		return
	}

	// Increment click count
	initializers.DB.Model(&link).Update("clicks", link.Clicks+1)

	// Return the link
	c.JSON(http.StatusOK, dtos.SuccessResponse{
		Success: true,
		Data: dtos.LinkResponse{
			ShortCode:   link.ShortCode,
			OriginalURL: link.OriginalURL,
			Clicks:      link.Clicks + 1,
			Favicon:     link.Favicon,
			UserID:      link.UserID,
		},
	})
}

// RedirectLink godoc
// @Summary Redirect to original URL
// @Description Redirect to original URL using short code and increment clicks
// @Tags Redirect
// @Accept json
// @Produce json
// @Param shortCode path string true "Short code of the link"
// @Success 301 {string} string "Moved Permanently"
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 404 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /{shortCode} [get]
func RedirectLink(c *gin.Context) {
	shortCode := c.Param("shortCode")
	fmt.Println("Attempting to redirect shortCode:", shortCode)

	if shortCode == "" {
		c.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			Success: false,
			Error:   "Short code is required",
		})
		return
	}

	var link models.Link
	result := initializers.DB.Where("short_code = ?", shortCode).First(&link)
	fmt.Println("Link found:", link.OriginalURL)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, dtos.ErrorResponse{
			Success: false,
			Error:   "Link not found",
		})
		return
	}

// 	if result.Error != nil {
//     if errors.Is(result.Error, gorm.ErrRecordNotFound) {
//         c.JSON(http.StatusNotFound, dtos.ErrorResponse{...}) // 404
//     } else {
//         c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{ // 500
//             Success: false,
//             Error:   "Database error: " + result.Error.Error(),
//         })
//     }
//     return
// }

	// Increment click count asynchronously to avoid blocking the redirect
	go initializers.DB.Model(&link).Update("clicks", link.Clicks+1)

	// Redirect to original URL (301 = permanent redirect)
	c.Redirect(http.StatusMovedPermanently, link.OriginalURL)
}

// UpdateLink godoc
// @Summary Update a link
// @Description Update an existing link (owner only)
// @Tags Links
// @Security Bearer
// @Accept json
// @Produce json
// @Param shortCode path string true "Short code of the link"
// @Param input body dtos.LinkUpdateRequest true "Update link request"
// @Success 200 {object} dtos.SuccessResponse{data=dtos.LinkResponse}
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 401 {object} dtos.ErrorResponse
// @Failure 403 {object} dtos.ErrorResponse
// @Failure 404 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /links/{shortCode} [patch]
func UpdateLink(c *gin.Context) {
	shortCode := c.Param("shortCode")

	var req dtos.LinkUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			Success: false,
			Error:   "Invalid input: " + err.Error(),
		})
		return
	}

	// Get authenticated user
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, dtos.ErrorResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	contextUser, ok := user.(ContextUserStruct)
	if !ok {
		c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			Success: false,
			Error:   "Invalid user data",
		})
		return
	}

	// Find the link
	var link models.Link
	result := initializers.DB.Where("short_code = ?", shortCode).First(&link)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, dtos.ErrorResponse{
			Success: false,
			Error:   "Link not found",
		})
		return
	}

	// Check ownership
	if link.UserID != contextUser.ID {
		c.JSON(http.StatusForbidden, dtos.ErrorResponse{
			Success: false,
			Error:   "You can only update your own links",
		})
		return
	}

	// Update fields
	updates := map[string]interface{}{}
	if req.OriginalURL != "" {
		updates["original_url"] = req.OriginalURL
		updates["hash"] = generateHash(req.OriginalURL)
	}

	if req.IsActive != nil {
		// Note: Add IsActive field to Link model if needed for soft delete logic
		// For now, this is a placeholder for future implementation
	}

	if len(updates) > 0 {
		initializers.DB.Model(&link).Updates(updates)
	}

	c.JSON(http.StatusOK, dtos.SuccessResponse{
		Success: true,
		Data: dtos.LinkResponse{
			ShortCode:   link.ShortCode,
			OriginalURL: link.OriginalURL,
			Clicks:      link.Clicks,
			Favicon:     link.Favicon,
			UserID:      link.UserID,
		},
	})
}

// DeleteLink godoc
// @Summary Delete a link
// @Description Delete a link by short code (owner only)
// @Tags Links
// @Security Bearer
// @Accept json
// @Produce json
// @Param shortCode path string true "Short code of the link"
// @Success 200 {object} dtos.SuccessResponse{data=map[string]string}
// @Failure 401 {object} dtos.ErrorResponse
// @Failure 403 {object} dtos.ErrorResponse
// @Failure 404 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /links/{shortCode} [delete]
func DeleteLink(c *gin.Context) {
	shortCode := c.Param("shortCode")

	// Get authenticated user
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, dtos.ErrorResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	contextUser, ok := user.(ContextUserStruct)
	if !ok {
		c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			Success: false,
			Error:   "Invalid user data",
		})
		return
	}

	// Find the link
	var link models.Link
	result := initializers.DB.Where("short_code = ?", shortCode).First(&link)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, dtos.ErrorResponse{
			Success: false,
			Error:   "Link not found",
		})
		return
	}

	// Check ownership
	if link.UserID != contextUser.ID {
		c.JSON(http.StatusForbidden, dtos.ErrorResponse{
			Success: false,
			Error:   "You can only delete your own links",
		})
		return
	}

	// Delete the link (soft delete)
	initializers.DB.Delete(&link)

	c.JSON(http.StatusOK, dtos.SuccessResponse{
		Success: true,
		Data: map[string]string{
			"message": "Link deleted successfully",
		},
	})
}

// @title GetUserLinks
// GetUserLinks godoc
// @Summary Get user's links
// @Description Retrieve all links created by authenticated user
// @Tags Links
// @Security Bearer
// @Accept json
// @Produce json
// @Success 200 {object} dtos.SuccessResponse{data=[]dtos.LinkResponse}
// @Failure 401 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /links [get]
func GetUserLinks(c *gin.Context) {
	// Get authenticated user
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, dtos.ErrorResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	contextUser, ok := user.(ContextUserStruct)
	if !ok {
		c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			Success: false,
			Error:   "Invalid user data",
		})
		return
	}

	var links []models.Link
	initializers.DB.Where("user_id = ?", contextUser.ID).Find(&links)

	// Convert to response DTOs
	var linkResponses []dtos.LinkResponse
	for _, link := range links {
		linkResponses = append(linkResponses, dtos.LinkResponse{
			ShortCode:   link.ShortCode,
			OriginalURL: link.OriginalURL,
			Clicks:      link.Clicks,
			Favicon:     link.Favicon,
			UserID:      link.UserID,
		})
	}

	c.JSON(http.StatusOK, dtos.SuccessResponse{
		Success: true,
		Data:    linkResponses,
	})
}

// Helper function: Generate a random short code
func generateShortCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 6
	shortCode := make([]byte, length)
	for i := range shortCode {
		shortCode[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortCode)
}

// Helper function: Generate MD5 hash of URL
func generateHash(url string) string {
	hash := md5.Sum([]byte(url))
	return fmt.Sprintf("%x", hash)
}
