package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"project/models"
	"project/services"
)

// AuthMiddleware validates JWT tokens from Authorization header
func AuthMiddleware(authService services.AuthService) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeAuthError(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			// Extract token from "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				writeAuthError(w, http.StatusUnauthorized, "invalid authorization header format")
				return
			}

			token := parts[1]

			// Validate token
			user, err := authService.ValidateToken(r.Context(), token)
			if err != nil {
				writeAuthError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			// Add user to context
			ctx := context.WithValue(r.Context(), "user", user)
			next(w, r.WithContext(ctx))
		}
	}
}

// RequireRole middleware checks if user has one of the required roles
// It returns a function that takes a handler and returns a handler
func RequireRole(roles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user := r.Context().Value("user").(*models.User)
			if user == nil {
				writeAuthError(w, http.StatusUnauthorized, "user not found in context")
				return
			}

			// Check if user role is in allowed roles
			hasRole := false
			for _, role := range roles {
				if user.Role == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				writeAuthError(w, http.StatusForbidden, fmt.Sprintf("user role '%s' not authorized for this action", user.Role))
				return
			}

			next(w, r)
		}
	}
}

// GetUserFromContext retrieves the user from request context
func GetUserFromContext(ctx context.Context) (*models.User, error) {
	user, ok := ctx.Value("user").(*models.User)
	if !ok || user == nil {
		return nil, fmt.Errorf("user not found in context")
	}
	return user, nil
}

// GetUserIDFromContext retrieves the user ID from request context
func GetUserIDFromContext(ctx context.Context) (string, error) {
	user, err := GetUserFromContext(ctx)
	if err != nil {
		return "", err
	}
	return user.ID, nil
}

// writeAuthError writes an authentication error response
func writeAuthError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, `{"success": false, "message": "%s", "code": "AUTH_ERROR"}`, message)
}

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			allowed := false

			// Check if origin is in allowed list
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin || allowedOrigin == "*" {
					allowed = true
					break
				}
			}

			// Always set CORS headers for requests with Origin header
			if origin != "" && allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token, Origin")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Max-Age", "86400")
			}

			// Handle preflight OPTIONS requests
			if r.Method == http.MethodOptions {
				if origin != "" && allowed {
					w.WriteHeader(http.StatusNoContent)
				} else if origin != "" {
					// Return 403 for disallowed origins on preflight
					w.WriteHeader(http.StatusForbidden)
				}
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
