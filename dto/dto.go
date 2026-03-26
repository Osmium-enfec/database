package dto

import (
	"project/models"
	"time"
)

// ============================================================================
// Auth DTOs
// ============================================================================

// LoginRequest represents user login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,max=255"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"omitempty,oneof=creator reviewer admin"` // Optional: defaults to 'creator'
}

// AuthResponse represents authentication response with tokens
type AuthResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	User         *UserResponse `json:"user"`
	ExpiresIn    int           `json:"expires_in"`
}

// ============================================================================
// User DTOs
// ============================================================================

// UserResponse represents a user in API response
type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUserResponse converts model to response DTO
func NewUserResponse(user *models.User) *UserResponse {
	if user == nil {
		return nil
	}
	return &UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ============================================================================
// Content DTOs
// ============================================================================

// CreateContentRequest represents content creation request
type CreateContentRequest struct {
	Type                 string             `json:"type" validate:"required,oneof=question code_problem documentation"`
	ProgramID            string             `json:"program_id" validate:"required,uuid"`
	TopicID              string             `json:"topic_id" validate:"required,uuid"`
	SubtopicID           string             `json:"subtopic_id" validate:"required,uuid"`
	Difficulty           string             `json:"difficulty" validate:"required,oneof=easy medium hard"`
	EstimatedTimeMinutes int                `json:"estimated_time_minutes" validate:"required,min=1,max=1440"`
	Tags                 []string           `json:"tags" validate:"max=20"`
	Data                 models.ContentData `json:"data" validate:"required"`
}

// UpdateContentRequest represents content update request
type UpdateContentRequest struct {
	Difficulty           string             `json:"difficulty" validate:"omitempty,oneof=easy medium hard"`
	EstimatedTimeMinutes *int               `json:"estimated_time_minutes" validate:"omitempty,min=1,max=1440"`
	Tags                 []string           `json:"tags" validate:"max=20"`
	Data                 models.ContentData `json:"data" validate:"required"`
}

// SubmitForReviewRequest represents content submission for review
type SubmitForReviewRequest struct {
	Comment    string `json:"comment" validate:"max=1000"`
	ReviewerID string `json:"reviewer_id" validate:"required"` // Reviewer to assign this for review
}

// ContentResponse represents content in API response
type ContentResponse struct {
	ID                   string                  `json:"id"`
	Type                 string                  `json:"type"`
	ProgramID            string                  `json:"program_id"`
	TopicID              string                  `json:"topic_id"`
	SubtopicID           string                  `json:"subtopic_id"`
	Difficulty           string                  `json:"difficulty"`
	EstimatedTimeMinutes int                     `json:"estimated_time_minutes"`
	Status               string                  `json:"status"`
	CreatedBy            string                  `json:"created_by"`
	CurrentVersionID     *string                 `json:"current_version_id"`
	IsActive             bool                    `json:"is_active"`
	Tags                 []string                `json:"tags"`
	Version              *ContentVersionResponse `json:"version,omitempty"`
	Creator              *UserResponse           `json:"creator,omitempty"`
	CreatedAt            time.Time               `json:"created_at"`
	UpdatedAt            time.Time               `json:"updated_at"`
}

// NewContentResponse converts model to response DTO
func NewContentResponse(content *models.Content, tags []string) *ContentResponse {
	if content == nil {
		return nil
	}
	resp := &ContentResponse{
		ID:                   content.ID,
		Type:                 content.Type,
		ProgramID:            content.ProgramID,
		TopicID:              content.TopicID,
		SubtopicID:           content.SubtopicID,
		Difficulty:           content.Difficulty,
		EstimatedTimeMinutes: content.EstimatedTimeMinutes,
		Status:               content.Status,
		CreatedBy:            content.CreatedBy,
		CurrentVersionID:     content.CurrentVersionID,
		IsActive:             content.IsActive,
		Tags:                 tags,
		CreatedAt:            content.CreatedAt,
		UpdatedAt:            content.UpdatedAt,
	}

	if content.Creator != nil {
		resp.Creator = NewUserResponse(content.Creator)
	}

	if content.Version != nil {
		resp.Version = NewContentVersionResponse(content.Version)
	}

	return resp
}

// ============================================================================
// Content Version DTOs
// ============================================================================

// ContentVersionResponse represents a content version in API response
type ContentVersionResponse struct {
	ID            string             `json:"id"`
	ContentID     string             `json:"content_id"`
	VersionNumber int                `json:"version_number"`
	Data          models.ContentData `json:"data"`
	CreatedBy     string             `json:"created_by"`
	ReviewStatus  string             `json:"review_status"`
	ReviewComment string             `json:"review_comment"`
	ReviewedBy    *string            `json:"reviewed_by"`
	Creator       *UserResponse      `json:"creator,omitempty"`
	Reviewer      *UserResponse      `json:"reviewer,omitempty"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
}

// NewContentVersionResponse converts model to response DTO
func NewContentVersionResponse(version *models.ContentVersion) *ContentVersionResponse {
	if version == nil {
		return nil
	}
	resp := &ContentVersionResponse{
		ID:            version.ID,
		ContentID:     version.ContentID,
		VersionNumber: version.VersionNumber,
		Data:          version.Data,
		CreatedBy:     version.CreatedBy,
		ReviewStatus:  version.ReviewStatus,
		ReviewComment: version.ReviewComment,
		ReviewedBy:    version.ReviewedBy,
		CreatedAt:     version.CreatedAt,
		UpdatedAt:     version.UpdatedAt,
	}

	if version.Creator != nil {
		resp.Creator = NewUserResponse(version.Creator)
	}

	if version.Reviewer != nil {
		resp.Reviewer = NewUserResponse(version.Reviewer)
	}

	return resp
}

// ApproveVersionRequest represents version approval request
type ApproveVersionRequest struct {
	Comment string `json:"comment" validate:"max=1000"`
}

// RejectVersionRequest represents version rejection request
type RejectVersionRequest struct {
	Comment string `json:"comment" validate:"required,max=1000"`
}

// ============================================================================
// List/Pagination DTOs
// ============================================================================

// PaginationMeta represents pagination information
type PaginationMeta struct {
	Page     int `json:"page"`
	PerPage  int `json:"per_page"`
	Total    int `json:"total"`
	LastPage int `json:"last_page"`
}

// ContentListRequest represents list query parameters
type ContentListRequest struct {
	Page       int    `json:"page" validate:"min=1"`
	PerPage    int    `json:"per_page" validate:"min=1,max=100"`
	Type       string `json:"type" validate:"omitempty,oneof=question code_problem documentation"`
	ProgramID  string `json:"program_id" validate:"omitempty,uuid"`
	TopicID    string `json:"topic_id" validate:"omitempty,uuid"`
	Difficulty string `json:"difficulty" validate:"omitempty,oneof=easy medium hard"`
	Status     string `json:"status" validate:"omitempty,oneof=draft pending_review approved rejected"`
	CreatedBy  string `json:"created_by" validate:"omitempty,uuid"`
	SortBy     string `json:"sort_by" validate:"omitempty,oneof=created_at updated_at difficulty"`
	SortOrder  string `json:"sort_order" validate:"omitempty,oneof=asc desc"`
}

// ReviewListRequest represents review list query parameters
type ReviewListRequest struct {
	Page      int    `json:"page" validate:"min=1"`
	PerPage   int    `json:"per_page" validate:"min=1,max=100"`
	Status    string `json:"status" validate:"omitempty,oneof=pending approved rejected"`
	SortBy    string `json:"sort_by" validate:"omitempty,oneof=created_at updated_at"`
	SortOrder string `json:"sort_order" validate:"omitempty,oneof=asc desc"`
}

// ============================================================================
// API Response Wrappers
// ============================================================================

// APIResponse represents a generic API response
type APIResponse struct {
	Success    bool            `json:"success"`
	Data       interface{}     `json:"data,omitempty"`
	Message    string          `json:"message,omitempty"`
	Code       string          `json:"code,omitempty"`
	Errors     []string        `json:"errors,omitempty"`
	Pagination *PaginationMeta `json:"pagination,omitempty"`
}

// NewSuccessResponse creates a success response
func NewSuccessResponse(data interface{}, message string) *APIResponse {
	return &APIResponse{
		Success: true,
		Data:    data,
		Message: message,
	}
}

// NewSuccessResponseWithPagination creates a paginated success response
func NewSuccessResponseWithPagination(data interface{}, pagination *PaginationMeta) *APIResponse {
	return &APIResponse{
		Success:    true,
		Data:       data,
		Pagination: pagination,
	}
}

// NewErrorResponse creates an error response
func NewErrorResponse(message string, code string, errors []string) *APIResponse {
	return &APIResponse{
		Success: false,
		Message: message,
		Code:    code,
		Errors:  errors,
	}
}

// ============================================================================
// Program/Topic/Subtopic DTOs (for dropdowns)
// ============================================================================

// ProgramResponse represents a program for dropdown
type ProgramResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// TopicResponse represents a topic for dropdown
type TopicResponse struct {
	ID          string `json:"id"`
	ProgramID   string `json:"program_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// SubtopicResponse represents a subtopic for dropdown
type SubtopicResponse struct {
	ID          string `json:"id"`
	TopicID     string `json:"topic_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ReviewerResponse represents a reviewer for dropdown
type ReviewerResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ============================================================================
// Bulk Upload DTOs
// ============================================================================

// BulkCreateContentRequest represents bulk content creation request
type BulkCreateContentRequest struct {
	Contents []CreateContentRequest `json:"contents" validate:"required,min=1,max=100"`
}

// BulkContentResultItem represents result for a single content creation
type BulkContentResultItem struct {
	Index   int    `json:"index"`
	Title   string `json:"title"`
	Success bool   `json:"success"`
	ID      string `json:"id,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

// BulkCreateContentResponse represents bulk content creation response
type BulkCreateContentResponse struct {
	Success        bool                    `json:"success"`
	TotalRequested int                     `json:"total_requested"`
	TotalCreated   int                     `json:"total_created"`
	TotalFailed    int                     `json:"total_failed"`
	Results        []BulkContentResultItem `json:"results"`
	Message        string                  `json:"message,omitempty"`
}
