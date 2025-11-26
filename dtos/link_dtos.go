package dtos

type CreateLinkRequest struct {
	// @notice The desired custom short code. Must be alphanumeric.
	// We've changed the validation from "email" to "alphanum".
	ShortCode   string `json:"shortCode" binding:"omitempty,alphanum,min=4,max=20"` 

	// @notice The full, long URL. Required and must be a valid URL format.
	OriginalURL string `json:"originalUrl" binding:"required,url"`

	UserID uint `json:"userId" binding:"omitempty,min=1"` 
}

type LinkUpdateRequest struct {
	// @notice The new target URL (optional for updates).
	OriginalURL string `json:"originalUrl" binding:"omitempty,url"`

	// @notice Flag to activate/deactivate the link (example field).
	IsActive *bool  `json:"isActive" binding:"omitempty"` 
}
type LinkResponse struct {
	// @notice The unique short identifier.
	ShortCode   string `json:"shortCode"`

	// @notice The full target URL.
	OriginalURL string `json:"originalUrl"`

	// @notice The number of times the link has been clicked.
	Clicks int `json:"clicks"`

	// @notice The URL of the favicon, can be null.
	Favicon  *string `json:"favicon"`

	// @notice The User ID this link belongs to.
	UserID uint `json:"userId"`
}
