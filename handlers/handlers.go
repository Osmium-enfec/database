package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"project/dto"
	"project/models"
	"project/services"
)

// ============================================================================
// Auth Handlers
// ============================================================================

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService services.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Create a new user account with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User registration details"
// @Success 200
// @Failure 400
// @Failure 400
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" || req.Name == "" {
		writeError(w, http.StatusBadRequest, "missing required fields")
		return
	}

	user, accessToken, refreshToken, err := h.authService.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			writeError(w, http.StatusConflict, "email already registered")
		} else {
			writeError(w, http.StatusInternalServerError, "registration failed")
		}
		return
	}

	resp := &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         dto.NewUserResponse(user),
		ExpiresIn:    86400,
	}

	writeJSON(w, http.StatusCreated, dto.NewSuccessResponse(resp, "User registered successfully"))
}

// Login handles user login
// @Summary User login
// @Description Authenticate user with email and password, returns JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "User login credentials"
// @Success 200
// @Failure 400
// @Failure 400
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "missing required fields")
		return
	}

	user, accessToken, refreshToken, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	resp := &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         dto.NewUserResponse(user),
		ExpiresIn:    86400,
	}

	writeJSON(w, http.StatusOK, dto.NewSuccessResponse(resp, "Login successful"))
}

// ============================================================================
// Content Handlers
// ============================================================================

// ContentHandler handles content-related endpoints
type ContentHandler struct {
	contentService services.ContentService
	reviewService  services.ReviewService
}

// NewContentHandler creates a new content handler
func NewContentHandler(contentService services.ContentService, reviewService services.ReviewService) *ContentHandler {
	return &ContentHandler{
		contentService: contentService,
		reviewService:  reviewService,
	}
}

// CreateContent handles content creation
// @Summary Create new content
// @Description Create a new content item (MCQ, MSQ, code problem, or documentation)
// @Tags content
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.CreateContentRequest true "Content creation details"
// @Success 200
// @Failure 400
// @Failure 400
// @Router /contents/create [post]
func (h *ContentHandler) CreateContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req dto.CreateContentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user := r.Context().Value("user").(*models.User)
	if user.Role != "creator" && user.Role != "admin" {
		writeError(w, http.StatusForbidden, "insufficient permissions")
		return
	}

	content, tags, err := h.contentService.Create(r.Context(), &req, user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create content")
		return
	}

	resp := dto.NewContentResponse(content, tags)
	writeJSON(w, http.StatusCreated, dto.NewSuccessResponse(resp, "Content created successfully"))
}

// GetContent handles single content retrieval
// @Summary Get content by ID
// @Description Retrieve a specific content item with all its versions and details
// @Tags content
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Content ID"
// @Success 200
// @Failure 400
// @Failure 400
// @Router /contents/{id} [get]
func (h *ContentHandler) GetContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/v1/contents/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "content id required")
		return
	}

	user := r.Context().Value("user").(*models.User)
	content, tags, err := h.contentService.Get(r.Context(), id, user.ID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "content not found")
		} else {
			writeError(w, http.StatusInternalServerError, "failed to retrieve content")
		}
		return
	}

	resp := dto.NewContentResponse(content, tags)
	writeJSON(w, http.StatusOK, dto.NewSuccessResponse(resp, ""))
}

// ListContents handles content listing
// @Summary List all content
// @Description List content items with pagination and filtering options
// @Tags content
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(20)
// @Param type query string false "Content type filter"
// @Param status query string false "Content status filter"
// @Success 200
// @Failure 400
// @Router /contents [get]
func (h *ContentHandler) ListContents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	page := 1
	perPage := 20

	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if pp := r.URL.Query().Get("per_page"); pp != "" {
		if parsed, err := strconv.Atoi(pp); err == nil && parsed > 0 && parsed <= 100 {
			perPage = parsed
		}
	}

	filters := map[string]string{
		"type":       r.URL.Query().Get("type"),
		"program_id": r.URL.Query().Get("program_id"),
		"topic_id":   r.URL.Query().Get("topic_id"),
		"difficulty": r.URL.Query().Get("difficulty"),
		"status":     r.URL.Query().Get("status"),
	}

	contents, total, err := h.contentService.List(r.Context(), page, perPage, filters)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list contents")
		return
	}

	var responses []*dto.ContentResponse
	for _, content := range contents {
		responses = append(responses, dto.NewContentResponse(content, nil))
	}

	pagination := &dto.PaginationMeta{
		Page:     page,
		PerPage:  perPage,
		Total:    total,
		LastPage: (total + perPage - 1) / perPage,
	}

	writeJSON(w, http.StatusOK, dto.NewSuccessResponseWithPagination(responses, pagination))
}

// UpdateContent handles content updates
// @Summary Update content
// @Description Update an existing content item
// @Tags content
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Content ID"
// @Param request body dto.UpdateContentRequest true "Update details"
// @Success 200
// @Failure 400
// @Failure 400
// @Failure 400
// @Router /contents/{id} [put]
func (h *ContentHandler) UpdateContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/v1/contents/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "content id required")
		return
	}

	var req dto.UpdateContentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user := r.Context().Value("user").(*models.User)
	content, tags, err := h.contentService.Update(r.Context(), id, &req, user.ID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "content not found")
		} else if strings.Contains(err.Error(), "permission") {
			writeError(w, http.StatusForbidden, "insufficient permissions")
		} else {
			writeError(w, http.StatusInternalServerError, "failed to update content")
		}
		return
	}

	resp := dto.NewContentResponse(content, tags)
	writeJSON(w, http.StatusOK, dto.NewSuccessResponse(resp, "Content updated successfully"))
}

// SubmitForReview handles submitting content for review
// @Summary Submit content for review
// @Description Submit content for reviewer approval
// @Tags content
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Content ID"
// @Success 200
// @Failure 400
// @Failure 400
// @Router /contents/{id}/submit [post]
func (h *ContentHandler) SubmitForReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/v1/contents/")
	if id == "" || !strings.HasSuffix(id, "/submit") {
		writeError(w, http.StatusBadRequest, "invalid path")
		return
	}

	id = strings.TrimSuffix(id, "/submit")

	var req dto.SubmitForReviewRequest
	json.NewDecoder(r.Body).Decode(&req)

	user := r.Context().Value("user").(*models.User)
	content, tags, err := h.contentService.SubmitForReview(r.Context(), id, user.ID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "content not found")
		} else {
			writeError(w, http.StatusInternalServerError, "failed to submit content")
		}
		return
	}

	resp := dto.NewContentResponse(content, tags)
	writeJSON(w, http.StatusOK, dto.NewSuccessResponse(resp, "Content submitted for review"))
}

// ============================================================================
// Review Handlers
// ============================================================================

// ReviewHandler handles review-related endpoints
type ReviewHandler struct {
	reviewService services.ReviewService
}

// NewReviewHandler creates a new review handler
func NewReviewHandler(reviewService services.ReviewService) *ReviewHandler {
	return &ReviewHandler{
		reviewService: reviewService,
	}
}

// GetPendingReviews handles retrieval of pending reviews
// @Summary Get pending reviews
// @Description Retrieve list of content pending review (Reviewer/Admin only)
// @Tags review
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(20)
// @Success 200
// @Failure 400
// @Failure 400
// @Router /reviews/pending [get]
func (h *ReviewHandler) GetPendingReviews(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	page := 1
	perPage := 20

	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if pp := r.URL.Query().Get("per_page"); pp != "" {
		if parsed, err := strconv.Atoi(pp); err == nil && parsed > 0 && parsed <= 100 {
			perPage = parsed
		}
	}

	versions, total, err := h.reviewService.GetPending(r.Context(), page, perPage)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to retrieve reviews")
		return
	}

	var responses []*dto.ContentVersionResponse
	for _, version := range versions {
		responses = append(responses, dto.NewContentVersionResponse(version))
	}

	pagination := &dto.PaginationMeta{
		Page:     page,
		PerPage:  perPage,
		Total:    total,
		LastPage: (total + perPage - 1) / perPage,
	}

	writeJSON(w, http.StatusOK, dto.NewSuccessResponseWithPagination(responses, pagination))
}

// ApproveVersion handles version approval
// @Summary Approve content version
// @Description Approve a content version for publication
// @Tags review
// @Accept json
// @Produce json
// @Security Bearer
// @Param version_id path string true "Content Version ID"
// @Param request body dto.ApproveVersionRequest true "Approval details"
// @Success 200
// @Failure 400
// @Failure 400
// @Failure 400
// @Router /reviews/{version_id}/approve [post]
func (h *ReviewHandler) ApproveVersion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	versionID := extractVersionIDFromPath(r.URL.Path, "/approve")
	if versionID == "" {
		writeError(w, http.StatusBadRequest, "version id required")
		return
	}

	var req dto.ApproveVersionRequest
	json.NewDecoder(r.Body).Decode(&req)

	user := r.Context().Value("user").(*models.User)
	version, err := h.reviewService.ApproveVersion(r.Context(), versionID, user.ID, req.Comment)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "version not found")
		} else {
			writeError(w, http.StatusInternalServerError, "failed to approve version")
		}
		return
	}

	resp := dto.NewContentVersionResponse(version)
	writeJSON(w, http.StatusOK, dto.NewSuccessResponse(resp, "Version approved successfully"))
}

// RejectVersion handles version rejection
// @Summary Reject content version
// @Description Reject a content version with feedback
// @Tags review
// @Accept json
// @Produce json
// @Security Bearer
// @Param version_id path string true "Content Version ID"
// @Param request body dto.RejectVersionRequest true "Rejection details"
// @Success 200
// @Failure 400
// @Failure 400
// @Failure 400
// @Router /reviews/{version_id}/reject [post]
func (h *ReviewHandler) RejectVersion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	versionID := extractVersionIDFromPath(r.URL.Path, "/reject")
	if versionID == "" {
		writeError(w, http.StatusBadRequest, "version id required")
		return
	}

	var req dto.RejectVersionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user := r.Context().Value("user").(*models.User)
	version, err := h.reviewService.RejectVersion(r.Context(), versionID, user.ID, req.Comment)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "version not found")
		} else {
			writeError(w, http.StatusInternalServerError, "failed to reject version")
		}
		return
	}

	resp := dto.NewContentVersionResponse(version)
	writeJSON(w, http.StatusOK, dto.NewSuccessResponse(resp, "Version rejected successfully"))
}

// Helper functions
func extractVersionIDFromPath(path string, suffix string) string {
	trimmed := strings.TrimPrefix(path, "/api/v1/reviews/")
	trimmed = strings.TrimSuffix(trimmed, suffix)
	return trimmed
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, dto.NewErrorResponse(message, "", nil))
}
