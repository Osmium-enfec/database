# Content Review API - Production-Grade Backend System

A production-grade backend system built in Go for managing educational content with a GitHub PR-like review workflow. This system handles content creation, versioning, and peer review with comprehensive role-based access control.

## 🎯 Overview

This is a full-stack content management and review system designed for educational platforms. It implements:

- **Content Management**: Create and manage questions, code problems, and documentation
- **Versioning**: Every change creates a new version with full history
- **Review Workflow**: Multi-step approval process similar to GitHub Pull Requests
- **Role-Based Access**: Creator, Reviewer, and Admin roles with specific permissions
- **JWT Authentication**: Secure token-based authentication with refresh tokens
- **PostgreSQL Database**: Fully normalized schema with proper indexing
- **Docker Support**: Complete containerization for easy deployment
- **Swagger Documentation**: Auto-generated API documentation

## 📋 Tech Stack

- **Language**: Go 1.21+
- **Framework**: Standard library (scale-library compliant patterns)
- **Database**: PostgreSQL 15+
- **Authentication**: JWT (golang-jwt)
- **Encryption**: bcrypt for passwords, crypto/sha256 for tokens
- **Containerization**: Docker + Docker Compose
- **Documentation**: Swagger/OpenAPI

## 🏗️ Project Structure

```
content-review-api/
├── main.go                 # Application entry point
├── go.mod                  # Go module definition
├── Dockerfile              # Docker image configuration
├── docker-compose.yml      # Docker Compose orchestration
├── migrations/
│   └── 001_initial_schema.sql  # Database schema
├── bootstrap/
│   ├── services.go         # Service registration
│   ├── middleware.go       # Middleware setup
│   └── routes.go           # Route definitions
├── handlers/
│   └── handlers.go         # HTTP request handlers
├── services/
│   └── services.go         # Business logic
├── repositories/
│   └── repositories.go     # Data access layer
├── models/
│   └── models.go           # Domain models
├── middleware/
│   └── auth.go             # JWT authentication middleware
├── config/
│   └── config.go           # Configuration management
└── dto/
    └── dto.go              # Data transfer objects
```

## 🗄️ Database Schema

### Core Tables

1. **users** - User accounts with roles
2. **programs** - Educational programs
3. **topics** - Topics within programs
4. **subtopics** - Subtopics within topics
5. **tags** - Content tags
6. **contents** - Main content records
7. **content_versions** - Version history with JSONB storage
8. **content_tags** - Many-to-many relationship
9. **refresh_tokens** - JWT refresh token management
10. **audit_logs** - Comprehensive audit trail

### Key Features

- **UUID Primary Keys**: Globally unique identifiers
- **Timestamps**: created_at, updated_at for all records
- **Soft Delete Pattern**: is_active flag for logical deletion
- **JSONB Storage**: Version data stored as JSON for flexibility
- **Foreign Keys**: Referential integrity with cascade delete
- **Indexes**: Performance optimization on frequently queried fields

## 👥 User Roles

### Creator
- Create new content
- Edit their own draft content
- Submit content for review
- View review feedback

### Reviewer
- View content pending review
- Approve content with optional comments
- Reject content with required feedback
- Cannot edit content directly

### Admin
- All permissions of Creator and Reviewer
- Manage users and programs
- Full access to all content
- Can override any workflow

## 📝 Content Types

### 1. Multiple Choice Questions (MCQ)
- Single correct answer
- Multiple options
- Metadata: difficulty, time estimate, tags

### 2. Multiple Select Questions (MSQ)
- Multiple correct answers
- Multiple options
- Same metadata as MCQ

### 3. Fill in the Blanks (FILL)
- Pattern matching or exact match
- Partial word completion
- Same metadata

### 4. Short Answer (SHORT)
- Free text answer
- Custom validation rules
- Same metadata

### 5. Code Problems
- Starter code provided
- Solution code
- Execution template
- Multiple test cases (visible and hidden)

### 6. Documentation
- Markdown content
- Syntax highlighting support
- Versioning and review workflow

## 🔄 Content Versioning Logic

1. **Draft Creation**: Initial content version created as draft
2. **Edits**: Each update creates a new version (never overwrites)
3. **Submission**: Creator submits for review (status → pending_review)
4. **Review**: Reviewer approves or rejects
5. **Activation**: Approved version becomes current_version_id
6. **History**: All versions retained for audit trail

### Example Version Flow
```
Content Created
  ↓
Version 1 (draft) - Created by user
  ↓
Edit by creator
  ↓
Version 2 (draft) - Edits applied
  ↓
Submit for review
  ↓
Status: pending_review
  ↓
Reviewer Approves
  ↓
Version 2 becomes current_version_id
Status: approved
  ↓
Creator edits again
  ↓
Version 3 (draft) - New edits
  ↓
Submit again → Status: pending_review
```

## 🔐 Authentication Flow

### Registration
```
POST /auth/register
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "secure_password_min_8_chars"
}

Response:
{
  "success": true,
  "data": {
    "access_token": "eyJ...",
    "refresh_token": "eyJ...",
    "user": { ... },
    "expires_in": 86400
  }
}
```

### Login
```
POST /auth/login
{
  "email": "john@example.com",
  "password": "secure_password"
}

Response: Same as registration
```

### Token Format
- **Access Token**: 24-hour validity, contains user_id:role in subject
- **Refresh Token**: 7-day validity, used to obtain new access token
- **Algorithm**: HS256 (HMAC with SHA-256)

### Using Tokens
```
Authorization: Bearer <access_token>
```

## 📡 API Endpoints

### Authentication Endpoints

#### Register User
```
POST /api/v1/auth/register
Content-Type: application/json

{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "password": "password123"
}

Response: 200 OK
{
  "success": true,
  "data": {
    "access_token": "...",
    "refresh_token": "...",
    "user": { "id": "uuid", "name": "Jane Doe", "email": "jane@example.com", "role": "creator" },
    "expires_in": 86400
  }
}
```

#### Login User
```
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "jane@example.com",
  "password": "password123"
}

Response: 200 OK (same format as register)
```

### Content Endpoints

#### Create Content
```
POST /api/v1/contents
Authorization: Bearer <token>
Content-Type: application/json

{
  "type": "question",
  "program_id": "uuid",
  "topic_id": "uuid",
  "subtopic_id": "uuid",
  "difficulty": "medium",
  "estimated_time_minutes": 15,
  "tags": ["algebra", "equations"],
  "data": {
    "type": "question",
    "title": "Solve for x",
    "question_text": "What is 2x + 5 = 15?",
    "question_type": "mcq",
    "options": ["x=5", "x=10", "x=15", "x=20"],
    "correct_options": [0]
  }
}

Response: 201 Created
{
  "success": true,
  "data": {
    "id": "uuid",
    "type": "question",
    "program_id": "uuid",
    "status": "draft",
    "created_by": "uuid",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

#### Get Content
```
GET /api/v1/contents/{id}
Authorization: Bearer <token>

Response: 200 OK
{
  "success": true,
  "data": {
    "id": "uuid",
    "type": "question",
    "status": "approved",
    "creator": { "id": "uuid", "name": "Jane Doe", "email": "jane@example.com", "role": "creator" },
    "tags": ["algebra", "equations"],
    "current_version": { ... }
  }
}
```

#### List Contents
```
GET /api/v1/contents?program_id=uuid&topic_id=uuid&difficulty=medium&type=question&status=approved&page=1&per_page=20
Authorization: Bearer <token>

Response: 200 OK
{
  "success": true,
  "data": [ ... ],
  "pagination": {
    "page": 1,
    "per_page": 20,
    "total": 150,
    "last_page": 8
  }
}
```

#### Update Content
```
PUT /api/v1/contents/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "difficulty": "hard",
  "estimated_time_minutes": 20,
  "tags": ["algebra", "advanced-equations"],
  "data": { ... }
}

Response: 200 OK
{
  "success": true,
  "data": { ... }
}
```

#### Submit Content for Review
```
POST /api/v1/contents/{id}/submit
Authorization: Bearer <token>
Content-Type: application/json

{
  "version_comment": "Ready for review"
}

Response: 200 OK
{
  "success": true,
  "data": {
    "id": "uuid",
    "status": "pending_review"
  }
}
```

### Review Endpoints

#### Get Pending Reviews
```
GET /api/v1/reviews/pending?page=1&per_page=20
Authorization: Bearer <token>
Roles: reviewer, admin

Response: 200 OK
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "content_id": "uuid",
      "version_number": 1,
      "data": { ... },
      "created_by": "uuid",
      "creator": { ... },
      "created_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

#### Approve Version
```
POST /api/v1/reviews/{version_id}/approve
Authorization: Bearer <token>
Content-Type: application/json
Roles: reviewer, admin

{
  "review_comment": "Looks good, approved!"
}

Response: 200 OK
{
  "success": true,
  "message": "version approved"
}
```

#### Reject Version
```
POST /api/v1/reviews/{version_id}/reject
Authorization: Bearer <token>
Content-Type: application/json
Roles: reviewer, admin

{
  "review_comment": "Please fix the grammar and clarify the question"
}

Response: 200 OK
{
  "success": true,
  "message": "version rejected"
}
```

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 15 or higher
- Docker and Docker Compose (optional)
- Git

### Local Development

#### 1. Clone Repository
```bash
git clone https://github.com/yourusername/content-review-api.git
cd content-review-api
```

#### 2. Setup Environment
```bash
cp .env.example .env
# Edit .env with your configuration
```

#### 3. Install Dependencies
```bash
go mod download
go mod tidy
```

#### 4. Setup Database
```bash
# Create database
createdb content_review

# Run migrations
psql content_review < migrations/001_initial_schema.sql
```

#### 5. Run Application
```bash
go run main.go
```

The API will be available at `http://localhost:8080`

### Docker Deployment

#### 1. Build and Run
```bash
docker-compose up -d
```

#### 2. Check Logs
```bash
docker-compose logs -f api
```

#### 3. Stop Services
```bash
docker-compose down
```

## 🐳 Docker Commands

### Build Image
```bash
docker build -t content-review-api:latest .
```

### Run Container
```bash
docker run -p 8080:8080 \
  -e DB_HOST=postgres \
  -e DB_PASSWORD=postgres \
  content-review-api:latest
```

### Docker Compose with Custom Port
```bash
docker-compose -f docker-compose.yml up -d
```

## 📊 Database Migrations

### Running Migrations
```bash
psql content_review -f migrations/001_initial_schema.sql
```

### Backup Database
```bash
pg_dump content_review > backup_$(date +%Y%m%d).sql
```

### Restore Database
```bash
psql content_review < backup_YYYYMMDD.sql
```

## 🧪 Testing

### Run Tests (Example)
```bash
go test ./...
```

### Test Coverage
```bash
go test -cover ./...
```

### Integration Tests
```bash
go test -tags=integration ./...
```

## 📚 API Documentation

Access Swagger documentation after starting the server:
- **Local**: http://localhost:8080/swagger/
- **Production**: https://api.example.com/swagger/

## 🔍 Monitoring & Logging

### Health Check Endpoint
```bash
curl http://localhost:8080/health
```

### Docker Container Logs
```bash
docker-compose logs -f api
```

### Database Monitoring
```bash
# Connect to PostgreSQL container
docker-compose exec postgres psql -U postgres -d content_review

# View table sizes
SELECT schemaname, tablename, pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename))
FROM pg_tables
WHERE schemaname != 'pg_catalog'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

## 🔐 Security Best Practices

1. **Environment Variables**: Never commit .env file with production secrets
2. **JWT Secrets**: Generate strong random secrets (64+ characters)
3. **HTTPS**: Use HTTPS in production (configure reverse proxy like Nginx)
4. **Database**: Use strong passwords, restrict network access
5. **CORS**: Configure appropriate CORS policies
6. **Rate Limiting**: Implement rate limiting for login endpoint
7. **Input Validation**: All inputs validated on server side
8. **SQL Injection**: Using parameterized queries prevents injection

## 📈 Performance Considerations

1. **Database Indexes**: Optimized for common queries (see schema)
2. **Connection Pooling**: Configured with 25 max connections
3. **Caching**: Ready for Redis integration
4. **Pagination**: All list endpoints support pagination (default 20 items)
5. **Query Optimization**: Efficient queries with appropriate joins

## 🛠️ Deployment Checklist

- [ ] Generate strong JWT secrets
- [ ] Configure PostgreSQL with proper backup strategy
- [ ] Set up SSL/TLS certificates
- [ ] Configure reverse proxy (Nginx/Apache)
- [ ] Setup monitoring and alerting
- [ ] Configure log aggregation (ELK/Splunk)
- [ ] Setup automated backups
- [ ] Configure CI/CD pipeline
- [ ] Load test the application
- [ ] Setup health checks

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## 📄 License

This project is licensed under the MIT License - see LICENSE file for details.

## 📞 Support

For issues and questions:
- Create an issue on GitHub
- Contact: support@example.com
- Documentation: https://docs.example.com

## 🎓 Architecture Decisions

1. **Standard Library First**: Uses Go standard library for compatibility
2. **Scalable Design**: Modular structure allows horizontal scaling
3. **Cloud-Ready**: Fully containerized and environment-configurable
4. **Version Control**: Every change is versioned for audit trail
5. **Separation of Concerns**: Clear layers (handlers, services, repositories)

## 📊 Scalability Roadmap

- [ ] Add caching layer (Redis)
- [ ] Implement message queue (RabbitMQ)
- [ ] Add search indexing (Elasticsearch)
- [ ] Database sharding support
- [ ] Microservices migration path
- [ ] GraphQL API option
- [ ] Real-time updates (WebSocket)
- [ ] Advanced analytics

---

**Built with ❤️ using Go**
