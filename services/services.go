package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"project/dto"
	"project/models"
	"project/repositories"

	"github.com/google/uuid"
)

// ============================================================================
// Auth Service
// ============================================================================

// AuthService defines authentication operations
type AuthService interface {
	Register(ctx context.Context, name, email, password string) (*models.User, string, string, error)
	Login(ctx context.Context, email, password string) (*models.User, string, string, error)
	ValidateToken(ctx context.Context, token string) (*models.User, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
}

type authService struct {
	userRepo      repositories.UserRepository
	jwtSecret     string
	refreshSecret string
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repositories.UserRepository, jwtSecret, refreshSecret string) AuthService {
	return &authService{
		userRepo:      userRepo,
		jwtSecret:     jwtSecret,
		refreshSecret: refreshSecret,
	}
}

// Register creates a new user account
func (s *authService) Register(ctx context.Context, name, email, password string) (*models.User, string, string, error) {
	// Check if user already exists
	_, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil {
		return nil, "", "", fmt.Errorf("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", "", err
	}

	// Create user
	user := &models.User{
		ID:           uuid.New().String(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         "creator", // Default role
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, "", "", err
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

// Login authenticates a user and returns tokens
func (s *authService) Login(ctx context.Context, email, password string) (*models.User, string, string, error) {
	// Find user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", "", fmt.Errorf("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", "", fmt.Errorf("invalid credentials")
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

// ValidateToken validates a JWT token and returns the user
func (s *authService) ValidateToken(ctx context.Context, tokenString string) (*models.User, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Extract subject (user ID) from claims
	userID, ok := claims["sub"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Get user from database
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// RefreshToken generates a new access token from a refresh token
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.refreshSecret), nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid refresh token")
	}

	// Extract subject (user ID) from claims
	userID, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", err
	}

	// Generate new access token
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

// generateAccessToken creates a JWT access token
func (s *authService) generateAccessToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// generateRefreshToken creates a JWT refresh token
func (s *authService) generateRefreshToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.refreshSecret))
}

// ============================================================================
// Content Service
// ============================================================================

// ContentService defines content operations
type ContentService interface {
	Create(ctx context.Context, req *dto.CreateContentRequest, userID string) (*models.Content, []string, error)
	Get(ctx context.Context, contentID, userID string) (*models.Content, []string, error)
	Update(ctx context.Context, contentID string, req *dto.UpdateContentRequest, userID string) (*models.Content, []string, error)
	List(ctx context.Context, page, perPage int, filters map[string]string) ([]*models.Content, int, error)
	SubmitForReview(ctx context.Context, contentID, userID string) (*models.Content, []string, error)
}

type contentService struct {
	contentRepo repositories.ContentRepository
	versionRepo repositories.ContentVersionRepository
	tagRepo     repositories.TagRepository
	userRepo    repositories.UserRepository
}

// NewContentService creates a new content service
func NewContentService(contentRepo repositories.ContentRepository, versionRepo repositories.ContentVersionRepository,
	tagRepo repositories.TagRepository, userRepo repositories.UserRepository) ContentService {
	return &contentService{
		contentRepo: contentRepo,
		versionRepo: versionRepo,
		tagRepo:     tagRepo,
		userRepo:    userRepo,
	}
}

// Create creates new content
func (s *contentService) Create(ctx context.Context, req *dto.CreateContentRequest, userID string) (*models.Content, []string, error) {
	// Create content
	content := &models.Content{
		ID:                   uuid.New().String(),
		Type:                 req.Type,
		ProgramID:            req.ProgramID,
		TopicID:              req.TopicID,
		SubtopicID:           req.SubtopicID,
		Difficulty:           req.Difficulty,
		EstimatedTimeMinutes: req.EstimatedTimeMinutes,
		Status:               "draft",
		CreatedBy:            userID,
		IsActive:             true,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	if err := s.contentRepo.Create(ctx, content); err != nil {
		return nil, nil, err
	}

	// Create initial version
	version := &models.ContentVersion{
		ID:            uuid.New().String(),
		ContentID:     content.ID,
		VersionNumber: 1,
		Data:          req.Data,
		CreatedBy:     userID,
		ReviewStatus:  "pending",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.versionRepo.Create(ctx, version); err != nil {
		return nil, nil, err
	}

	// Handle tags
	var tags []string
	for _, tagName := range req.Tags {
		tag, err := s.tagRepo.GetByName(ctx, tagName)
		if err != nil {
			// Create tag if it doesn't exist
			tag = &models.Tag{
				ID:        uuid.New().String(),
				Name:      tagName,
				IsActive:  true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			s.tagRepo.Create(ctx, tag)
		}
		s.tagRepo.AssignTag(ctx, content.ID, tag.ID)
		tags = append(tags, tagName)
	}

	return content, tags, nil
}

// Get retrieves a specific content
func (s *contentService) Get(ctx context.Context, contentID, userID string) (*models.Content, []string, error) {
	content, err := s.contentRepo.GetByID(ctx, contentID)
	if err != nil {
		return nil, nil, err
	}

	// Get tags
	tags, err := s.tagRepo.GetTagsByContentID(ctx, contentID)
	if err != nil {
		tags = []string{}
	}

	// Get current version if exists
	if content.CurrentVersionID != nil {
		version, err := s.versionRepo.GetByID(ctx, *content.CurrentVersionID)
		if err == nil {
			content.Version = version
		}
	}

	// Get creator info
	creator, err := s.userRepo.GetByID(ctx, content.CreatedBy)
	if err == nil {
		content.Creator = creator
	}

	return content, tags, nil
}

// Update updates content
func (s *contentService) Update(ctx context.Context, contentID string, req *dto.UpdateContentRequest, userID string) (*models.Content, []string, error) {
	// Get existing content
	content, err := s.contentRepo.GetByID(ctx, contentID)
	if err != nil {
		return nil, nil, err
	}

	// Check permission - only creator or admin can update
	if content.CreatedBy != userID {
		return nil, nil, fmt.Errorf("permission denied")
	}

	// Create new version
	latestVersion, err := s.versionRepo.GetLatestByContentID(ctx, contentID)
	if err != nil {
		return nil, nil, err
	}

	newVersion := &models.ContentVersion{
		ID:            uuid.New().String(),
		ContentID:     contentID,
		VersionNumber: latestVersion.VersionNumber + 1,
		Data:          req.Data,
		CreatedBy:     userID,
		ReviewStatus:  "pending",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.versionRepo.Create(ctx, newVersion); err != nil {
		return nil, nil, err
	}

	// Update content fields if provided
	if req.Difficulty != "" {
		content.Difficulty = req.Difficulty
	}
	if req.EstimatedTimeMinutes != nil {
		content.EstimatedTimeMinutes = *req.EstimatedTimeMinutes
	}

	// Update status back to draft
	content.Status = "draft"
	content.UpdatedAt = time.Now()

	if err := s.contentRepo.Update(ctx, content); err != nil {
		return nil, nil, err
	}

	// Handle tags
	var tags []string
	for _, tagName := range req.Tags {
		tag, err := s.tagRepo.GetByName(ctx, tagName)
		if err != nil {
			// Create tag if it doesn't exist
			tag = &models.Tag{
				ID:        uuid.New().String(),
				Name:      tagName,
				IsActive:  true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			s.tagRepo.Create(ctx, tag)
		}
		s.tagRepo.AssignTag(ctx, contentID, tag.ID)
		tags = append(tags, tagName)
	}

	return content, tags, nil
}

// List retrieves paginated content
func (s *contentService) List(ctx context.Context, page, perPage int, filters map[string]string) ([]*models.Content, int, error) {
	return s.contentRepo.List(ctx, page, perPage, filters)
}

// SubmitForReview submits content for review
func (s *contentService) SubmitForReview(ctx context.Context, contentID, userID string) (*models.Content, []string, error) {
	content, err := s.contentRepo.GetByID(ctx, contentID)
	if err != nil {
		return nil, nil, err
	}

	// Check permission
	if content.CreatedBy != userID {
		return nil, nil, fmt.Errorf("permission denied")
	}

	// Update status
	content.Status = "pending_review"
	content.UpdatedAt = time.Now()

	if err := s.contentRepo.Update(ctx, content); err != nil {
		return nil, nil, err
	}

	// Get tags
	tags, err := s.tagRepo.GetTagsByContentID(ctx, contentID)
	if err != nil {
		tags = []string{}
	}

	return content, tags, nil
}

// ============================================================================
// Review Service
// ============================================================================

// ReviewService defines review operations
type ReviewService interface {
	GetPending(ctx context.Context, page, perPage int) ([]*models.ContentVersion, int, error)
	ApproveVersion(ctx context.Context, versionID, reviewerID, comment string) (*models.ContentVersion, error)
	RejectVersion(ctx context.Context, versionID, reviewerID, comment string) (*models.ContentVersion, error)
}

type reviewService struct {
	versionRepo repositories.ContentVersionRepository
	contentRepo repositories.ContentRepository
	userRepo    repositories.UserRepository
}

// NewReviewService creates a new review service
func NewReviewService(versionRepo repositories.ContentVersionRepository,
	contentRepo repositories.ContentRepository, userRepo repositories.UserRepository) ReviewService {
	return &reviewService{
		versionRepo: versionRepo,
		contentRepo: contentRepo,
		userRepo:    userRepo,
	}
}

// GetPending retrieves pending reviews
func (s *reviewService) GetPending(ctx context.Context, page, perPage int) ([]*models.ContentVersion, int, error) {
	versions, total, err := s.versionRepo.GetPending(ctx, page, perPage)
	if err != nil {
		return nil, 0, err
	}

	// Enrich with creator info
	for _, version := range versions {
		if creator, err := s.userRepo.GetByID(ctx, version.CreatedBy); err == nil {
			version.Creator = creator
		}
	}

	return versions, total, nil
}

// ApproveVersion approves a content version
func (s *reviewService) ApproveVersion(ctx context.Context, versionID, reviewerID, comment string) (*models.ContentVersion, error) {
	version, err := s.versionRepo.GetByID(ctx, versionID)
	if err != nil {
		return nil, err
	}

	// Update version
	version.ReviewStatus = "approved"
	version.ReviewComment = comment
	version.ReviewedBy = &reviewerID
	version.UpdatedAt = time.Now()

	if err := s.versionRepo.Update(ctx, version); err != nil {
		return nil, err
	}

	// Update content status and set current version
	content, err := s.contentRepo.GetByID(ctx, version.ContentID)
	if err != nil {
		return nil, err
	}

	content.Status = "approved"
	content.CurrentVersionID = &versionID
	content.UpdatedAt = time.Now()

	if err := s.contentRepo.Update(ctx, content); err != nil {
		return nil, err
	}

	// Get reviewer info
	if reviewer, err := s.userRepo.GetByID(ctx, reviewerID); err == nil {
		version.Reviewer = reviewer
	}

	// Get creator info
	if creator, err := s.userRepo.GetByID(ctx, version.CreatedBy); err == nil {
		version.Creator = creator
	}

	return version, nil
}

// RejectVersion rejects a content version
func (s *reviewService) RejectVersion(ctx context.Context, versionID, reviewerID, comment string) (*models.ContentVersion, error) {
	version, err := s.versionRepo.GetByID(ctx, versionID)
	if err != nil {
		return nil, err
	}

	// Update version
	version.ReviewStatus = "rejected"
	version.ReviewComment = comment
	version.ReviewedBy = &reviewerID
	version.UpdatedAt = time.Now()

	if err := s.versionRepo.Update(ctx, version); err != nil {
		return nil, err
	}

	// Update content status to rejected
	content, err := s.contentRepo.GetByID(ctx, version.ContentID)
	if err != nil {
		return nil, err
	}

	content.Status = "rejected"
	content.UpdatedAt = time.Now()

	if err := s.contentRepo.Update(ctx, content); err != nil {
		return nil, err
	}

	// Get reviewer info
	if reviewer, err := s.userRepo.GetByID(ctx, reviewerID); err == nil {
		version.Reviewer = reviewer
	}

	// Get creator info
	if creator, err := s.userRepo.GetByID(ctx, version.CreatedBy); err == nil {
		version.Creator = creator
	}

	return version, nil
}

// hashToken creates a SHA256 hash of a token for storage
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
