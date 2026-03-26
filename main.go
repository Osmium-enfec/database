package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

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
// @schemes https http
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	db, err := connectDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize database schema
	if err := initializeSchema(db); err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}

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
	mux.HandleFunc("/api/v1/contents/bulk", withAuth(authService, contentHandler.BulkCreateContent))
	mux.HandleFunc("/api/v1/contents/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			withAuth(authService, contentHandler.GetContent)(w, r)
		} else if r.Method == http.MethodPost {
			// Check if it's a submit request
			if strings.HasSuffix(r.URL.Path, "/submit") {
				withAuth(authService, contentHandler.SubmitForReview)(w, r)
			} else {
				http.Error(w, `{"success": false, "message": "invalid path"}`, http.StatusBadRequest)
			}
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

	// Dropdown endpoints (public - no auth required)
	dropdownHandler := handlers.NewDropdownHandler(contentService)
	mux.HandleFunc("/api/v1/dropdown/programs", func(w http.ResponseWriter, r *http.Request) {
		// Inject db into context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "db", db)
		dropdownHandler.GetPrograms(w, r.WithContext(ctx))
	})
	mux.HandleFunc("/api/v1/dropdown/topics", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, "db", db)
		dropdownHandler.GetTopicsByProgram(w, r.WithContext(ctx))
	})
	mux.HandleFunc("/api/v1/dropdown/subtopics", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, "db", db)
		dropdownHandler.GetSubtopicsByTopic(w, r.WithContext(ctx))
	})
	mux.HandleFunc("/api/v1/dropdown/reviewers", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, "db", db)
		dropdownHandler.GetReviewers(w, r.WithContext(ctx))
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

func initializeSchema(db *sql.DB) error {
	// SQL schema - all tables
	schema := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	CREATE EXTENSION IF NOT EXISTS "pgcrypto";

	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		password_hash VARCHAR(255) NOT NULL,
		role VARCHAR(50) NOT NULL DEFAULT 'creator' CHECK (role IN ('creator', 'reviewer', 'admin')),
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS programs (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		name VARCHAR(255) NOT NULL,
		description TEXT,
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS topics (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		program_id UUID NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS subtopics (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		topic_id UUID NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS contents (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		type VARCHAR(50) NOT NULL CHECK (type IN ('question', 'code_problem', 'documentation')),
		program_id UUID NOT NULL REFERENCES programs(id),
		topic_id UUID NOT NULL REFERENCES topics(id),
		subtopic_id UUID NOT NULL REFERENCES subtopics(id),
		difficulty VARCHAR(50) NOT NULL CHECK (difficulty IN ('easy', 'medium', 'hard')),
		estimated_time_minutes INT NOT NULL,
		status VARCHAR(50) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'pending_review', 'approved', 'rejected')),
		created_by UUID NOT NULL REFERENCES users(id),
		current_version_id UUID,
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS content_versions (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		content_id UUID NOT NULL REFERENCES contents(id) ON DELETE CASCADE,
		version_number INT NOT NULL,
		data JSONB NOT NULL,
		created_by UUID NOT NULL REFERENCES users(id),
		review_status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (review_status IN ('pending', 'approved', 'rejected')),
		review_comment TEXT,
		reviewed_by UUID REFERENCES users(id),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(content_id, version_number)
	);

	CREATE TABLE IF NOT EXISTS content_tags (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		content_id UUID NOT NULL REFERENCES contents(id) ON DELETE CASCADE,
		tag VARCHAR(100) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(content_id, tag)
	);

	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	CREATE INDEX IF NOT EXISTS idx_programs_is_active ON programs(is_active);
	CREATE INDEX IF NOT EXISTS idx_topics_program_id ON topics(program_id);
	CREATE INDEX IF NOT EXISTS idx_contents_created_by ON contents(created_by);
	CREATE INDEX IF NOT EXISTS idx_contents_status ON contents(status);
	CREATE INDEX IF NOT EXISTS idx_content_versions_content_id ON content_versions(content_id);
	CREATE INDEX IF NOT EXISTS idx_content_tags_content_id ON content_tags(content_id);
	`

	_, err := db.Exec(schema)
	if err != nil {
		log.Printf("Warning: Failed to initialize schema (may already exist): %v", err)
		// Don't fail - schema might already exist
		return nil
	}

	log.Println("Database schema initialized successfully")

	// Seed Python roadmap data
	if err := seedPythonRoadmap(db); err != nil {
		log.Printf("Warning: Failed to seed Python roadmap: %v", err)
		// Don't fail - seed data might already exist
	}

	// Seed Java roadmap data
	if err := seedJavaRoadmap(db); err != nil {
		log.Printf("Warning: Failed to seed Java roadmap: %v", err)
		// Don't fail - seed data might already exist
	}

	return nil
}

func seedPythonRoadmap(db *sql.DB) error {
	// Check if Python program already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM programs WHERE name = 'Python: Zero to Senior Developer'").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Println("Python roadmap already seeded")
		return nil
	}

	// Python program with all 16 phases
	pythonSQL := `
	INSERT INTO programs (name, description, is_active) 
	VALUES ('Python: Zero to Senior Developer', 'Complete Python learning roadmap from beginner to senior developer level', true)
	ON CONFLICT DO NOTHING;

	-- Phase 0: Foundations
	INSERT INTO topics (program_id, name, description, is_active)
	SELECT id, 'Phase 0: Foundations of Programming', 'Beginner Mindset - What is programming and setting up environment', true
	FROM programs WHERE name = 'Python: Zero to Senior Developer' LIMIT 1
	ON CONFLICT DO NOTHING;

	INSERT INTO subtopics (topic_id, name, description, is_active) VALUES
	((SELECT id FROM topics WHERE name = 'Phase 0: Foundations of Programming' LIMIT 1), 'What is Programming?', 'How computers execute code, Interpreted vs compiled languages', true),
	((SELECT id FROM topics WHERE name = 'Phase 0: Foundations of Programming' LIMIT 1), 'Setting Up Environment', 'Installing Python, VS Code, Virtual environments', true),
	((SELECT id FROM topics WHERE name = 'Phase 0: Foundations of Programming' LIMIT 1), 'First Program', 'print(), Comments, Basic syntax', true)
	ON CONFLICT DO NOTHING;

	-- Phase 1: Core Python Basics
	INSERT INTO topics (program_id, name, description, is_active)
	SELECT id, 'Phase 1: Core Python Basics', 'Variables, operators, conditionals, loops', true
	FROM programs WHERE name = 'Python: Zero to Senior Developer' LIMIT 1
	ON CONFLICT DO NOTHING;

	INSERT INTO subtopics (topic_id, name, description, is_active) VALUES
	((SELECT id FROM topics WHERE name = 'Phase 1: Core Python Basics' LIMIT 1), 'Variables & Data Types', 'Integers, Floats, Strings, Booleans, Type casting', true),
	((SELECT id FROM topics WHERE name = 'Phase 1: Core Python Basics' LIMIT 1), 'Operators', 'Arithmetic, Comparison, Logical, Assignment operators', true),
	((SELECT id FROM topics WHERE name = 'Phase 1: Core Python Basics' LIMIT 1), 'Input/Output', 'input() function, Formatting strings (f-strings)', true),
	((SELECT id FROM topics WHERE name = 'Phase 1: Core Python Basics' LIMIT 1), 'Conditional Statements', 'if, elif, else, Nested conditions', true),
	((SELECT id FROM topics WHERE name = 'Phase 1: Core Python Basics' LIMIT 1), 'Loops', 'for, while, break, continue, pass', true)
	ON CONFLICT DO NOTHING;

	-- Phase 2: Data Structures
	INSERT INTO topics (program_id, name, description, is_active)
	SELECT id, 'Phase 2: Data Structures', 'Core data structures - Lists, Tuples, Sets, Dictionaries', true
	FROM programs WHERE name = 'Python: Zero to Senior Developer' LIMIT 1
	ON CONFLICT DO NOTHING;

	INSERT INTO subtopics (topic_id, name, description, is_active) VALUES
	((SELECT id FROM topics WHERE name = 'Phase 2: Data Structures' LIMIT 1), 'Strings Deep Dive', 'Indexing, slicing, String methods', true),
	((SELECT id FROM topics WHERE name = 'Phase 2: Data Structures' LIMIT 1), 'Lists', 'CRUD operations, List slicing, List comprehensions', true),
	((SELECT id FROM topics WHERE name = 'Phase 2: Data Structures' LIMIT 1), 'Tuples', 'Immutability, Packing/unpacking', true),
	((SELECT id FROM topics WHERE name = 'Phase 2: Data Structures' LIMIT 1), 'Sets', 'Unique elements, Set operations', true),
	((SELECT id FROM topics WHERE name = 'Phase 2: Data Structures' LIMIT 1), 'Dictionaries', 'Key-value structure, Nested dictionaries', true)
	ON CONFLICT DO NOTHING;

	-- Phase 3: Functions & Modular Code
	INSERT INTO topics (program_id, name, description, is_active)
	SELECT id, 'Phase 3: Functions & Modular Code', 'Functions, Lambda, Recursion, Modules & Packages', true
	FROM programs WHERE name = 'Python: Zero to Senior Developer' LIMIT 1
	ON CONFLICT DO NOTHING;

	INSERT INTO subtopics (topic_id, name, description, is_active) VALUES
	((SELECT id FROM topics WHERE name = 'Phase 3: Functions & Modular Code' LIMIT 1), 'Functions Basics', 'Defining & calling, Arguments, Keyword arguments', true),
	((SELECT id FROM topics WHERE name = 'Phase 3: Functions & Modular Code' LIMIT 1), 'Advanced Arguments', '*args, **kwargs, Default parameters', true),
	((SELECT id FROM topics WHERE name = 'Phase 3: Functions & Modular Code' LIMIT 1), 'Lambda Functions', 'Anonymous functions, Using lambda with map/filter', true),
	((SELECT id FROM topics WHERE name = 'Phase 3: Functions & Modular Code' LIMIT 1), 'Recursion', 'Base cases, Call stack', true),
	((SELECT id FROM topics WHERE name = 'Phase 3: Functions & Modular Code' LIMIT 1), 'Modules & Packages', 'Importing modules, Creating modules', true)
	ON CONFLICT DO NOTHING;

	-- Phase 4: Object-Oriented Programming
	INSERT INTO topics (program_id, name, description, is_active)
	SELECT id, 'Phase 4: Object-Oriented Programming', 'Classes, Inheritance, Polymorphism, Encapsulation', true
	FROM programs WHERE name = 'Python: Zero to Senior Developer' LIMIT 1
	ON CONFLICT DO NOTHING;

	INSERT INTO subtopics (topic_id, name, description, is_active) VALUES
	((SELECT id FROM topics WHERE name = 'Phase 4: Object-Oriented Programming' LIMIT 1), 'Classes & Objects', 'Attributes & methods, Creating instances', true),
	((SELECT id FROM topics WHERE name = 'Phase 4: Object-Oriented Programming' LIMIT 1), 'Constructors', '__init__ method, Initialization', true),
	((SELECT id FROM topics WHERE name = 'Phase 4: Object-Oriented Programming' LIMIT 1), 'Encapsulation', 'Private variables, Getters/Setters', true),
	((SELECT id FROM topics WHERE name = 'Phase 4: Object-Oriented Programming' LIMIT 1), 'Inheritance', 'Single & multiple inheritance, super()', true),
	((SELECT id FROM topics WHERE name = 'Phase 4: Object-Oriented Programming' LIMIT 1), 'Polymorphism', 'Method overriding, Duck typing', true),
	((SELECT id FROM topics WHERE name = 'Phase 4: Object-Oriented Programming' LIMIT 1), 'Magic Methods', '__str__, __repr__, __eq__', true),
	((SELECT id FROM topics WHERE name = 'Phase 4: Object-Oriented Programming' LIMIT 1), 'Abstract Classes', 'Interface design, ABC module', true)
	ON CONFLICT DO NOTHING;

	-- Phase 5: Error Handling
	INSERT INTO topics (program_id, name, description, is_active)
	SELECT id, 'Phase 5: Error Handling & Debugging', 'Exceptions, Debugging, Custom exceptions', true
	FROM programs WHERE name = 'Python: Zero to Senior Developer' LIMIT 1
	ON CONFLICT DO NOTHING;

	INSERT INTO subtopics (topic_id, name, description, is_active) VALUES
	((SELECT id FROM topics WHERE name = 'Phase 5: Error Handling & Debugging' LIMIT 1), 'Exceptions', 'try, except, finally, else blocks', true),
	((SELECT id FROM topics WHERE name = 'Phase 5: Error Handling & Debugging' LIMIT 1), 'Custom Exceptions', 'Creating custom exception classes', true),
	((SELECT id FROM topics WHERE name = 'Phase 5: Error Handling & Debugging' LIMIT 1), 'Debugging Techniques', 'Breakpoints, Stack traces, Logging', true)
	ON CONFLICT DO NOTHING;

	-- Phase 6: File Handling
	INSERT INTO topics (program_id, name, description, is_active)
	SELECT id, 'Phase 6: File Handling & Data Processing', 'Files, JSON, CSV, Logging', true
	FROM programs WHERE name = 'Python: Zero to Senior Developer' LIMIT 1
	ON CONFLICT DO NOTHING;

	INSERT INTO subtopics (topic_id, name, description, is_active) VALUES
	((SELECT id FROM topics WHERE name = 'Phase 6: File Handling & Data Processing' LIMIT 1), 'File Operations', 'Read/write files, Context managers', true),
	((SELECT id FROM topics WHERE name = 'Phase 6: File Handling & Data Processing' LIMIT 1), 'Working with JSON', 'json.load, json.dump, Parsing', true),
	((SELECT id FROM topics WHERE name = 'Phase 6: File Handling & Data Processing' LIMIT 1), 'CSV Handling', 'csv module, pandas basics', true),
	((SELECT id FROM topics WHERE name = 'Phase 6: File Handling & Data Processing' LIMIT 1), 'Logging', 'logging module, Log levels', true)
	ON CONFLICT DO NOTHING;

	-- Phase 7: Advanced Concepts
	INSERT INTO topics (program_id, name, description, is_active)
	SELECT id, 'Phase 7: Advanced Python Concepts', 'Iterators, Generators, Decorators, Context Managers', true
	FROM programs WHERE name = 'Python: Zero to Senior Developer' LIMIT 1
	ON CONFLICT DO NOTHING;

	INSERT INTO subtopics (topic_id, name, description, is_active) VALUES
	((SELECT id FROM topics WHERE name = 'Phase 7: Advanced Python Concepts' LIMIT 1), 'Iterators & Generators', 'yield, generator functions', true),
	((SELECT id FROM topics WHERE name = 'Phase 7: Advanced Python Concepts' LIMIT 1), 'Decorators', 'Function wrapping, Parameterized decorators', true),
	((SELECT id FROM topics WHERE name = 'Phase 7: Advanced Python Concepts' LIMIT 1), 'Context Managers', 'with statement, __enter__ and __exit__', true),
	((SELECT id FROM topics WHERE name = 'Phase 7: Advanced Python Concepts' LIMIT 1), 'Closures', 'Nested functions, Variable capture', true)
	ON CONFLICT DO NOTHING;

	-- Phase 8: APIs & Networking
	INSERT INTO topics (program_id, name, description, is_active)
	SELECT id, 'Phase 8: APIs & Networking', 'HTTP, REST, requests library, Authentication', true
	FROM programs WHERE name = 'Python: Zero to Senior Developer' LIMIT 1
	ON CONFLICT DO NOTHING;

	INSERT INTO subtopics (topic_id, name, description, is_active) VALUES
	((SELECT id FROM topics WHERE name = 'Phase 8: APIs & Networking' LIMIT 1), 'HTTP Basics', 'GET, POST, PUT, DELETE, Status codes', true),
	((SELECT id FROM topics WHERE name = 'Phase 8: APIs & Networking' LIMIT 1), 'Using APIs', 'requests library, JSON responses', true),
	((SELECT id FROM topics WHERE name = 'Phase 8: APIs & Networking' LIMIT 1), 'REST Concepts', 'RESTful principles, HTTP methods', true),
	((SELECT id FROM topics WHERE name = 'Phase 8: APIs & Networking' LIMIT 1), 'Authentication', 'API keys, JWT basics', true)
	ON CONFLICT DO NOTHING;

	-- Phase 9: Concurrency
	INSERT INTO topics (program_id, name, description, is_active)
	SELECT id, 'Phase 9: Concurrency & Performance', 'Threading, Multiprocessing, Async', true
	FROM programs WHERE name = 'Python: Zero to Senior Developer' LIMIT 1
	ON CONFLICT DO NOTHING;

	INSERT INTO subtopics (topic_id, name, description, is_active) VALUES
	((SELECT id FROM topics WHERE name = 'Phase 9: Concurrency & Performance' LIMIT 1), 'Multithreading', 'Thread class, Synchronization, Locks', true),
	((SELECT id FROM topics WHERE name = 'Phase 9: Concurrency & Performance' LIMIT 1), 'Multiprocessing', 'Process class, Process pools', true),
	((SELECT id FROM topics WHERE name = 'Phase 9: Concurrency & Performance' LIMIT 1), 'Async Programming', 'async/await, asyncio, Event loops', true)
	ON CONFLICT DO NOTHING;

	-- Phase 10: Databases
	INSERT INTO topics (program_id, name, description, is_active)
	SELECT id, 'Phase 10: Databases', 'SQL, PostgreSQL, ORM, Transactions', true
	FROM programs WHERE name = 'Python: Zero to Senior Developer' LIMIT 1
	ON CONFLICT DO NOTHING;

	INSERT INTO subtopics (topic_id, name, description, is_active) VALUES
	((SELECT id FROM topics WHERE name = 'Phase 10: Databases' LIMIT 1), 'SQL Basics', 'SELECT, INSERT, UPDATE, DELETE, Joins', true),
	((SELECT id FROM topics WHERE name = 'Phase 10: Databases' LIMIT 1), 'PostgreSQL Integration', 'psycopg2, Connection strings', true),
	((SELECT id FROM topics WHERE name = 'Phase 10: Databases' LIMIT 1), 'ORM', 'SQLAlchemy, Django ORM', true),
	((SELECT id FROM topics WHERE name = 'Phase 10: Databases' LIMIT 1), 'Transactions', 'ACID properties, Commit/Rollback', true)
	ON CONFLICT DO NOTHING;

	-- Phase 11: Testing
	INSERT INTO topics (program_id, name, description, is_active)
	SELECT id, 'Phase 11: Testing & Code Quality', 'Unit Testing, Mocking, Coverage', true
	FROM programs WHERE name = 'Python: Zero to Senior Developer' LIMIT 1
	ON CONFLICT DO NOTHING;

	INSERT INTO subtopics (topic_id, name, description, is_active) VALUES
	((SELECT id FROM topics WHERE name = 'Phase 11: Testing & Code Quality' LIMIT 1), 'Unit Testing', 'unittest, pytest, Writing tests', true),
	((SELECT id FROM topics WHERE name = 'Phase 11: Testing & Code Quality' LIMIT 1), 'Mocking', 'unittest.mock, Mocking dependencies', true),
	((SELECT id FROM topics WHERE name = 'Phase 11: Testing & Code Quality' LIMIT 1), 'Code Coverage', 'Coverage.py', true),
	((SELECT id FROM topics WHERE name = 'Phase 11: Testing & Code Quality' LIMIT 1), 'Linting & Formatting', 'flake8, black, Code style', true)
	ON CONFLICT DO NOTHING;

	-- Phase 12: Backend Development
	INSERT INTO topics (program_id, name, description, is_active)
	SELECT id, 'Phase 12: Backend Development', 'Web frameworks, APIs, Authentication', true
	FROM programs WHERE name = 'Python: Zero to Senior Developer' LIMIT 1
	ON CONFLICT DO NOTHING;

	INSERT INTO subtopics (topic_id, name, description, is_active) VALUES
	((SELECT id FROM topics WHERE name = 'Phase 12: Backend Development' LIMIT 1), 'Web Basics', 'HTTP lifecycle, Request/Response', true),
	((SELECT id FROM topics WHERE name = 'Phase 12: Backend Development' LIMIT 1), 'Frameworks', 'Flask, FastAPI, Django basics', true),
	((SELECT id FROM topics WHERE name = 'Phase 12: Backend Development' LIMIT 1), 'Building APIs', 'RESTful APIs, Validation', true),
	((SELECT id FROM topics WHERE name = 'Phase 12: Backend Development' LIMIT 1), 'Authentication Systems', 'JWT, Sessions, OAuth', true),
	((SELECT id FROM topics WHERE name = 'Phase 12: Backend Development' LIMIT 1), 'Middleware', 'CORS, Rate limiting', true)
	ON CONFLICT DO NOTHING;

	-- Phase 13: System Design
	INSERT INTO topics (program_id, name, description, is_active)
	SELECT id, 'Phase 13: System Design with Python', 'Architecture, Design Patterns, Scaling', true
	FROM programs WHERE name = 'Python: Zero to Senior Developer' LIMIT 1
	ON CONFLICT DO NOTHING;

	INSERT INTO subtopics (topic_id, name, description, is_active) VALUES
	((SELECT id FROM topics WHERE name = 'Phase 13: System Design with Python' LIMIT 1), 'Project Structure', 'Folder organization, Modularity', true),
	((SELECT id FROM topics WHERE name = 'Phase 13: System Design with Python' LIMIT 1), 'Design Patterns', 'Singleton, Factory, Observer', true),
	((SELECT id FROM topics WHERE name = 'Phase 13: System Design with Python' LIMIT 1), 'Scaling Applications', 'Load balancing, Microservices', true),
	((SELECT id FROM topics WHERE name = 'Phase 13: System Design with Python' LIMIT 1), 'Caching', 'Redis, Memcached, Strategies', true),
	((SELECT id FROM topics WHERE name = 'Phase 13: System Design with Python' LIMIT 1), 'Message Queues', 'Kafka, RabbitMQ, Event-driven', true)
	ON CONFLICT DO NOTHING;

	-- Phase 14: DevOps & Deployment
	INSERT INTO topics (program_id, name, description, is_active)
	SELECT id, 'Phase 14: DevOps & Deployment', 'Git, Docker, CI/CD, Cloud Deployment', true
	FROM programs WHERE name = 'Python: Zero to Senior Developer' LIMIT 1
	ON CONFLICT DO NOTHING;

	INSERT INTO subtopics (topic_id, name, description, is_active) VALUES
	((SELECT id FROM topics WHERE name = 'Phase 14: DevOps & Deployment' LIMIT 1), 'Git & GitHub', 'Version control, Branching, PRs', true),
	((SELECT id FROM topics WHERE name = 'Phase 14: DevOps & Deployment' LIMIT 1), 'Docker', 'Containers, Images, Docker Compose', true),
	((SELECT id FROM topics WHERE name = 'Phase 14: DevOps & Deployment' LIMIT 1), 'CI/CD', 'GitHub Actions, Jenkins', true),
	((SELECT id FROM topics WHERE name = 'Phase 14: DevOps & Deployment' LIMIT 1), 'Cloud Deployment', 'AWS, GCP, Heroku, Render', true)
	ON CONFLICT DO NOTHING;

	-- Phase 15: Specialization Tracks
	INSERT INTO topics (program_id, name, description, is_active)
	SELECT id, 'Phase 15: Specialization Tracks', 'Backend, Data Engineering, AI/ML, Automation', true
	FROM programs WHERE name = 'Python: Zero to Senior Developer' LIMIT 1
	ON CONFLICT DO NOTHING;

	INSERT INTO subtopics (topic_id, name, description, is_active) VALUES
	((SELECT id FROM topics WHERE name = 'Phase 15: Specialization Tracks' LIMIT 1), 'Backend Engineer Track', 'FastAPI, Microservices, Distributed systems', true),
	((SELECT id FROM topics WHERE name = 'Phase 15: Specialization Tracks' LIMIT 1), 'Data Engineering Track', 'Pandas, ETL pipelines, Airflow', true),
	((SELECT id FROM topics WHERE name = 'Phase 15: Specialization Tracks' LIMIT 1), 'AI/ML Track', 'Scikit-learn, TensorFlow, Deep learning', true),
	((SELECT id FROM topics WHERE name = 'Phase 15: Specialization Tracks' LIMIT 1), 'Automation Track', 'Web scraping, Task automation', true)
	ON CONFLICT DO NOTHING;

	-- Phase 16: Senior-Level Skills
	INSERT INTO topics (program_id, name, description, is_active)
	SELECT id, 'Phase 16: Senior-Level Skills', 'Architecture, System Design, Mentoring', true
	FROM programs WHERE name = 'Python: Zero to Senior Developer' LIMIT 1
	ON CONFLICT DO NOTHING;

	INSERT INTO subtopics (topic_id, name, description, is_active) VALUES
	((SELECT id FROM topics WHERE name = 'Phase 16: Senior-Level Skills' LIMIT 1), 'Code Architecture', 'Clean code, SOLID principles', true),
	((SELECT id FROM topics WHERE name = 'Phase 16: Senior-Level Skills' LIMIT 1), 'System Design Interviews', 'Scalable systems, Trade-offs', true),
	((SELECT id FROM topics WHERE name = 'Phase 16: Senior-Level Skills' LIMIT 1), 'Open Source Contribution', 'Contributing to projects', true),
	((SELECT id FROM topics WHERE name = 'Phase 16: Senior-Level Skills' LIMIT 1), 'Mentoring & Code Reviews', 'Effective reviews, Mentoring', true),
	((SELECT id FROM topics WHERE name = 'Phase 16: Senior-Level Skills' LIMIT 1), 'Production Systems', 'Reliability, Performance, Monitoring', true)
	ON CONFLICT DO NOTHING;
	`

	// Split by semicolon and execute each statement
	statements := strings.Split(pythonSQL, ";")
	for _, stmt := range statements {
		trimmed := strings.TrimSpace(stmt)
		if trimmed == "" {
			continue
		}
		if _, err := db.Exec(trimmed); err != nil {
			log.Printf("Warning: Error executing statement: %v", err)
			// Continue on error as some inserts may conflict
		}
	}

	log.Println("Python roadmap seeded successfully with all topics")
	return nil
}

func seedJavaRoadmap(db *sql.DB) error {
	// Check if Java program already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM programs WHERE name = 'Java Programming'").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Println("Java roadmap already seeded")
		return nil
	}

	// Read and execute Java seed data from migration file
	javaSQL := `
	-- Insert Java program
	INSERT INTO programs (name, description, is_active) 
	VALUES ('Java Programming', 'Comprehensive Java learning from basics to advanced Java 8 features', true);

	-- TOPIC 1: Core Concepts
	INSERT INTO topics (program_id, name, description, is_active)
	SELECT id, 'Core Concepts', 'Java fundamentals, data types, operators, and control flow', true
	FROM programs WHERE name = 'Java Programming' LIMIT 1;

	INSERT INTO subtopics (topic_id, name, description, is_active) VALUES
	((SELECT id FROM topics WHERE name = 'Core Concepts' LIMIT 1), 'Basics & Installation', 'Intro, Installation, program execution flow, types of code editors, Compilation and Execution', true),
	((SELECT id FROM topics WHERE name = 'Core Concepts' LIMIT 1), 'Variables and Data Types', 'Primitives (byte, short, int, long, double, float, char, boolean) and Non-Primitives (arrays, Strings, classes)', true),
	((SELECT id FROM topics WHERE name = 'Core Concepts' LIMIT 1), 'Identifiers, Rules & Keywords', 'Identifier rules and Java keywords', true),
	((SELECT id FROM topics WHERE name = 'Core Concepts' LIMIT 1), 'Comments and Naming Conventions', 'Code documentation and naming best practices', true),
	((SELECT id FROM topics WHERE name = 'Core Concepts' LIMIT 1), 'Scanner Class', 'Reading input from keyboard using Scanner class', true),
	((SELECT id FROM topics WHERE name = 'Core Concepts' LIMIT 1), 'Operators', 'Arithmetic, logical, relational, new, assignment, and instanceof operators', true),
	((SELECT id FROM topics WHERE name = 'Core Concepts' LIMIT 1), 'Control Statements', 'if, if-else, if-else ladder, switch, while, do-while, for, for-each loops, and break/continue/return', true),
	((SELECT id FROM topics WHERE name = 'Core Concepts' LIMIT 1), 'Arrays & Its Types', '1D and 2D arrays', true),
	((SELECT id FROM topics WHERE name = 'Core Concepts' LIMIT 1), 'Strings', 'String creation, manipulation methods, StringBuffer, StringBuilder, StringTokenizer', true),
	((SELECT id FROM topics WHERE name = 'Core Concepts' LIMIT 1), 'Type Casting', 'Widening/Narrowing with primitives, UpCasting/DownCasting with non-primitives', true),
	((SELECT id FROM topics WHERE name = 'Core Concepts' LIMIT 1), 'Wrapper Classes', 'Wrapper class uses and importance', true);

	-- TOPIC 2: Object-Oriented Programming
	INSERT INTO topics (program_id, name, description, is_active)
	SELECT id, 'Object-Oriented Programming', 'Classes, objects, inheritance, polymorphism, encapsulation, abstraction', true
	FROM programs WHERE name = 'Java Programming' LIMIT 1;

	INSERT INTO subtopics (topic_id, name, description, is_active) VALUES
	((SELECT id FROM topics WHERE name = 'Object-Oriented Programming' LIMIT 1), 'Classes & Objects', 'Class definition and object creation', true),
	((SELECT id FROM topics WHERE name = 'Object-Oriented Programming' LIMIT 1), 'Packages & Types', 'Package organization and types', true),
	((SELECT id FROM topics WHERE name = 'Object-Oriented Programming' LIMIT 1), 'Inheritance', 'Single, multiple, multilevel, and hierarchical inheritance', true),
	((SELECT id FROM topics WHERE name = 'Object-Oriented Programming' LIMIT 1), 'SUPER Keyword', 'Using super keyword for parent class access', true),
	((SELECT id FROM topics WHERE name = 'Object-Oriented Programming' LIMIT 1), 'Encapsulation', 'Data hiding and access control', true),
	((SELECT id FROM topics WHERE name = 'Object-Oriented Programming' LIMIT 1), 'Polymorphism', 'Method overloading and overriding', true),
	((SELECT id FROM topics WHERE name = 'Object-Oriented Programming' LIMIT 1), 'Abstraction', 'Abstract concept implementation', true),
	((SELECT id FROM topics WHERE name = 'Object-Oriented Programming' LIMIT 1), 'Access Modifiers', 'public, private, protected, default access levels', true),
	((SELECT id FROM topics WHERE name = 'Object-Oriented Programming' LIMIT 1), 'FINAL Keyword', 'Using final keyword for classes, methods, and variables', true),
	((SELECT id FROM topics WHERE name = 'Object-Oriented Programming' LIMIT 1), 'Abstract Classes', 'Abstract class definition and usage', true),
	((SELECT id FROM topics WHERE name = 'Object-Oriented Programming' LIMIT 1), 'Object Class', 'Object class methods and functionality', true),
	((SELECT id FROM topics WHERE name = 'Object-Oriented Programming' LIMIT 1), 'Interfaces', 'Interface definition, implementation, abstract methods, marker interfaces', true),
	((SELECT id FROM topics WHERE name = 'Object-Oriented Programming' LIMIT 1), 'Types of Variables', 'Instance, local, and static variables', true),
	((SELECT id FROM topics WHERE name = 'Object-Oriented Programming' LIMIT 1), 'Methods', 'Method combinations and types', true),
	((SELECT id FROM topics WHERE name = 'Object-Oriented Programming' LIMIT 1), 'Constructors', 'Constructor types and THIS keyword usage', true);

	-- TOPIC 3: Exception Handling
	INSERT INTO topics (program_id, name, description, is_active)
	SELECT id, 'Exception Handling', 'Checked/unchecked exceptions, try-catch-finally, custom exceptions', true
	FROM programs WHERE name = 'Java Programming' LIMIT 1;

	INSERT INTO subtopics (topic_id, name, description, is_active) VALUES
	((SELECT id FROM topics WHERE name = 'Exception Handling' LIMIT 1), 'Exception Types', 'Checked and unchecked exceptions', true),
	((SELECT id FROM topics WHERE name = 'Exception Handling' LIMIT 1), 'Try-Catch-Finally', 'try, catch, finally, throw, throws keywords', true),
	((SELECT id FROM topics WHERE name = 'Exception Handling' LIMIT 1), 'Try with Resources', 'Automatic resource management', true),
	((SELECT id FROM topics WHERE name = 'Exception Handling' LIMIT 1), 'Custom Exceptions', 'Creating and using user-defined exceptions', true);

	-- TOPIC 4: Collections Framework
	INSERT INTO topics (program_id, name, description, is_active)
	SELECT id, 'Collections Framework', 'List, Set, Map, Queues, and related interfaces', true
	FROM programs WHERE name = 'Java Programming' LIMIT 1;

	INSERT INTO subtopics (topic_id, name, description, is_active) VALUES
	((SELECT id FROM topics WHERE name = 'Collections Framework' LIMIT 1), 'Cursors & Interfaces', 'Iterators, Comparable, Comparator interfaces', true),
	((SELECT id FROM topics WHERE name = 'Collections Framework' LIMIT 1), 'List Interface', 'ArrayList, LinkedList, Vector, and Stack', true),
	((SELECT id FROM topics WHERE name = 'Collections Framework' LIMIT 1), 'Queue Interface', 'Queue implementations and operations', true),
	((SELECT id FROM topics WHERE name = 'Collections Framework' LIMIT 1), 'Set Interface', 'HashSet, LinkedHashSet, SortedSet, TreeSet', true),
	((SELECT id FROM topics WHERE name = 'Collections Framework' LIMIT 1), 'Map Interface', 'HashMap, LinkedHashMap, TreeMap, Hashtable, IdentityHashMap, WeakHashMap', true);

	-- TOPIC 5: Multi-Threading
	INSERT INTO topics (program_id, name, description, is_active)
	SELECT id, 'Multi-Threading', 'Threads, synchronization, concurrent programming', true
	FROM programs WHERE name = 'Java Programming' LIMIT 1;

	INSERT INTO subtopics (topic_id, name, description, is_active) VALUES
	((SELECT id FROM topics WHERE name = 'Multi-Threading' LIMIT 1), 'Thread Basics', 'Default threads and user-defined threads', true),
	((SELECT id FROM topics WHERE name = 'Multi-Threading' LIMIT 1), 'Thread Life Cycle', 'Thread states and transitions', true),
	((SELECT id FROM topics WHERE name = 'Multi-Threading' LIMIT 1), 'Callable & Executor', 'Callable interface and ExecutorService', true),
	((SELECT id FROM topics WHERE name = 'Multi-Threading' LIMIT 1), 'Daemon Threads', 'Daemon thread creation and usage', true),
	((SELECT id FROM topics WHERE name = 'Multi-Threading' LIMIT 1), 'Synchronization', 'Synchronization techniques and methods', true);

	-- TOPIC 6: File IO & Serialization
	INSERT INTO topics (program_id, name, description, is_active)
	SELECT id, 'File IO & Serialization', 'File operations, streams, serialization/deserialization', true
	FROM programs WHERE name = 'Java Programming' LIMIT 1;

	INSERT INTO subtopics (topic_id, name, description, is_active) VALUES
	((SELECT id FROM topics WHERE name = 'File IO & Serialization' LIMIT 1), 'File Handling', 'Create, write, read, delete file operations', true),
	((SELECT id FROM topics WHERE name = 'File IO & Serialization' LIMIT 1), 'IO Streams', 'FileWriter, FileReader, and stream operations', true),
	((SELECT id FROM topics WHERE name = 'File IO & Serialization' LIMIT 1), 'Serialization', 'SerialVersionUID, transient keyword, serialization/deserialization', true);

	-- TOPIC 7: Generics
	INSERT INTO topics (program_id, name, description, is_active)
	SELECT id, 'Generics', 'Generic types, wildcards, type parameters, bounded types', true
	FROM programs WHERE name = 'Java Programming' LIMIT 1;

	INSERT INTO subtopics (topic_id, name, description, is_active) VALUES
	((SELECT id FROM topics WHERE name = 'Generics' LIMIT 1), 'Generic Basics', 'Generic type parameters and type safety', true),
	((SELECT id FROM topics WHERE name = 'Generics' LIMIT 1), 'Wildcards', 'Wildcard types and bounded types', true),
	((SELECT id FROM topics WHERE name = 'Generics' LIMIT 1), 'Generic Methods', 'Generic method definition and usage', true);

	-- TOPIC 8: Java 8 Features
	INSERT INTO topics (program_id, name, description, is_active)
	SELECT id, 'Java 8 Features', 'Lambda expressions, functional interfaces, streams, and new APIs', true
	FROM programs WHERE name = 'Java Programming' LIMIT 1;

	INSERT INTO subtopics (topic_id, name, description, is_active) VALUES
	((SELECT id FROM topics WHERE name = 'Java 8 Features' LIMIT 1), 'Interface Changes', 'default and static methods in interfaces', true),
	((SELECT id FROM topics WHERE name = 'Java 8 Features' LIMIT 1), 'Lambda Expressions', 'Lambda syntax and functional programming', true),
	((SELECT id FROM topics WHERE name = 'Java 8 Features' LIMIT 1), 'Functional Interfaces', 'Consumer, Supplier, Predicate, Function', true),
	((SELECT id FROM topics WHERE name = 'Java 8 Features' LIMIT 1), 'Stream API', 'Creation, filtering, mapping, slicing, matching, finding, collectors', true),
	((SELECT id FROM topics WHERE name = 'Java 8 Features' LIMIT 1), 'Parallel Streams', 'Parallel stream processing', true),
	((SELECT id FROM topics WHERE name = 'Java 8 Features' LIMIT 1), 'Date & Time API', 'New date/time API changes', true),
	((SELECT id FROM topics WHERE name = 'Java 8 Features' LIMIT 1), 'Optional Class', 'Optional class and its methods', true),
	((SELECT id FROM topics WHERE name = 'Java 8 Features' LIMIT 1), 'Method References', 'Method and constructor references', true),
	((SELECT id FROM topics WHERE name = 'Java 8 Features' LIMIT 1), 'Utilities', 'forEach() method, Spliterator, StringJoiner', true);
	`

	_, err = db.Exec(javaSQL)
	if err != nil {
		return fmt.Errorf("failed to seed Java roadmap: %v", err)
	}

	log.Println("Java roadmap seeded successfully")
	return nil
}

func seedPythonRoadmapFull(db *sql.DB) error {
	// Placeholder for future full seeding
	return nil
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
