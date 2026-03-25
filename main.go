package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"project/config"
	"project/handlers"
	"project/middleware"
	"project/repositories"
	"project/services"

	_ "project/docs"

	_ "github.com/lib/pq"
)

// @title Content Review API
// @version 1.0
// @description A production-grade content management and review system
// @contact.name API Support
// @contact.url http://www.example.com/support
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host content-review-api-bnkf.onrender.com
// @basePath /api/v1
// @schemes http https

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	db, err := connectDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	contentRepo := repositories.NewContentRepository(db)
	versionRepo := repositories.NewContentVersionRepository(db)
	tagRepo := repositories.NewTagRepository(db)

	// Initialize services
	authService := services.NewAuthService(
		userRepo,
		cfg.JWT.Secret,
		cfg.JWT.RefreshSecret,
	)

	contentService := services.NewContentService(
		contentRepo,
		versionRepo,
		tagRepo,
		userRepo,
	)

	reviewService := services.NewReviewService(
		versionRepo,
		contentRepo,
		userRepo,
	)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	contentHandler := handlers.NewContentHandler(contentService, reviewService)
	reviewHandler := handlers.NewReviewHandler(reviewService)

	// Setup routes
	mux := http.NewServeMux()

	// Configure CORS middleware
	corsHandler := middleware.CORSMiddleware(cfg.CORS.AllowedOrigins)

	// Health check
	mux.HandleFunc("/health", healthCheck)

	// Swagger documentation - custom handler with token persistence
	mux.HandleFunc("/swagger/", serveSwaggerUI)
	mux.HandleFunc("/swagger/index.html", serveSwaggerUI)

	// Swagger doc.json endpoint
	mux.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	// Serve swagger static assets (yaml, etc)
	mux.HandleFunc("/swagger/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/yaml")
		http.ServeFile(w, r, "./docs/swagger.yaml")
	})

	// Auth endpoints (no auth required)
	mux.HandleFunc("/api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("/api/v1/auth/login", authHandler.Login)

	// Content endpoints (auth required)
	mux.HandleFunc("/api/v1/contents", withAuth(authService, contentHandler.ListContents))
	mux.HandleFunc("/api/v1/contents/create", withAuth(authService, contentHandler.CreateContent))
	mux.HandleFunc("/api/v1/contents/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			withAuth(authService, contentHandler.GetContent)(w, r)
		} else {
			http.Error(w, `{"success": false, "message": "invalid method"}`, http.StatusBadRequest)
		}
	})

	// Review endpoints (auth + role required)
	mux.HandleFunc("/api/v1/reviews/pending", func(w http.ResponseWriter, r *http.Request) {
		roleMiddleware := withRole("reviewer", "admin")
		withAuth(authService, roleMiddleware(reviewHandler.GetPendingReviews))(w, r)
	})
	mux.HandleFunc("/api/v1/reviews/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if len(r.URL.Path) <= len("/api/v1/reviews/") {
				http.Error(w, `{"success": false, "message": "invalid path"}`, http.StatusBadRequest)
				return
			}
			roleMiddleware := withRole("reviewer", "admin")
			if r.URL.Path[len(r.URL.Path)-len("/approve"):] == "/approve" {
				withAuth(authService, roleMiddleware(reviewHandler.ApproveVersion))(w, r)
			} else if r.URL.Path[len(r.URL.Path)-len("/reject"):] == "/reject" {
				withAuth(authService, roleMiddleware(reviewHandler.RejectVersion))(w, r)
			} else {
				http.Error(w, `{"success": false, "message": "invalid path"}`, http.StatusBadRequest)
			}
		}
	})

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	log.Printf("CORS enabled for: %v", cfg.CORS.AllowedOrigins)

	// Wrap mux with CORS middleware
	handler := corsHandler(mux)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func connectDatabase(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.Database.MaxConn)
	db.SetMaxIdleConns(cfg.Database.MaxConn / 2)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func serveSwaggerUI(w http.ResponseWriter, r *http.Request) {
	// Serve custom Swagger HTML with token persistence
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Read and serve custom swagger HTML
	htmlContent := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Content Review API - Swagger UI</title>
  <link rel="stylesheet" type="text/css" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.15.5/swagger-ui.min.css" >
  <link rel="icon" type="image/png" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.15.5/favicon-32x32.png" sizes="32x32" />
  <link rel="icon" type="image/png" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.15.5/favicon-16x16.png" sizes="16x16" />
  <style>
    html {
      box-sizing: border-box;
      overflow: -moz-scrollbars-vertical;
      overflow-y: scroll;
    }
    *,
    *:before,
    *:after {
      box-sizing: inherit;
    }
    body {
      margin: 0;
      background: #fafafa;
    }
    .topbar {
      background-color: #0d47a1;
      padding: 10px 0;
      border-bottom: 3px solid #00bcd4;
    }
    .topbar-title {
      color: white;
      font-size: 24px;
      font-weight: bold;
      padding: 10px 20px;
      margin: 0;
    }
    .token-info {
      background: #e3f2fd;
      padding: 10px 20px;
      border-left: 4px solid #2196f3;
      margin: 10px 20px;
      border-radius: 4px;
      font-size: 13px;
      color: #1565c0;
    }
    .token-info strong {
      color: #0d47a1;
    }
  </style>
</head>

<body>
  <div class="topbar">
    <h1 class="topbar-title">📚 Content Review API Documentation</h1>
  </div>
  <div class="token-info">
    <strong>💡 Token Persistence:</strong> Your Bearer token is automatically saved for 30 days. 
    Use the "Authorize" button to add your token - it will be remembered!
  </div>
  <div id="swagger-ui"></div>

  <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.15.5/swagger-ui-bundle.min.js" charset="UTF-8"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.15.5/swagger-ui-standalone-preset.min.js" charset="UTF-8"></script>
  <script>
    const STORAGE_KEY = 'swagger_bearer_token';
    const TOKEN_EXPIRY_KEY = 'swagger_token_expiry';
    
    function getStoredToken() {
      const token = localStorage.getItem(STORAGE_KEY);
      const expiry = localStorage.getItem(TOKEN_EXPIRY_KEY);
      
      if (!token || !expiry) return null;
      if (new Date().getTime() > parseInt(expiry)) {
        localStorage.removeItem(STORAGE_KEY);
        localStorage.removeItem(TOKEN_EXPIRY_KEY);
        return null;
      }
      return token;
    }
    
    function saveToken(token) {
      if (token) {
        const thirtyDaysInMs = 30 * 24 * 60 * 60 * 1000;
        const expiryTime = new Date().getTime() + thirtyDaysInMs;
        localStorage.setItem(STORAGE_KEY, token);
        localStorage.setItem(TOKEN_EXPIRY_KEY, expiryTime.toString());
        console.log('✅ Token saved! Valid for 30 days.');
      }
    }
    
    function clearToken() {
      localStorage.removeItem(STORAGE_KEY);
      localStorage.removeItem(TOKEN_EXPIRY_KEY);
      console.log('🗑️  Token cleared');
    }
    
    window.onload = function() {
      const ui = SwaggerUIBundle({
        url: "/swagger/doc.json",
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout",
        onComplete: function() {
          const storedToken = getStoredToken();
          if (storedToken) {
            console.log('✅ Token auto-loaded from storage (30-day expiry)!');
          }
        }
      });
      
      // Monitor for authorization changes
      const observer = new MutationObserver(function(mutations) {
        const authButton = document.querySelector('button[aria-label="authorize"]') || 
                          document.querySelector('button.auth');
        if (authButton && !authButton.hasAttribute('data-intercepted')) {
          authButton.setAttribute('data-intercepted', 'true');
          authButton.addEventListener('click', function() {
            setTimeout(function() {
              const authInput = document.querySelector('input[placeholder*="Bearer"]') ||
                               document.querySelector('input[type="password"]');
              if (authInput && authInput.value && authInput.value.startsWith('eyJ')) {
                saveToken(authInput.value);
              }
            }, 500);
          });
        }
      });
      
      observer.observe(document.body, { childList: true, subtree: true });
      
      // Expose utilities
      window.swaggerTokenUtils = {
        saveToken: saveToken,
        getStoredToken: getStoredToken,
        clearToken: clearToken,
        status: function() {
          const token = getStoredToken();
          console.log(token ? '✅ Token stored' : '❌ No token stored');
        }
      };
    }
  </script>
</body>
</html>`

	w.Write([]byte(htmlContent))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}

// Middleware to add auth
func withAuth(authService services.AuthService, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		middleware.AuthMiddleware(authService)(handler).ServeHTTP(w, r)
	}
}

// Middleware to check role
func withRole(roles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			middleware.RequireRole(roles...)(handler).ServeHTTP(w, r)
		}
	}
}
