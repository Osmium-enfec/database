# Content Review API - Complete API Reference

## Overview

Base URL: `http://localhost:8080/api/v1`

All requests (except auth endpoints) require `Authorization: Bearer <access_token>` header.

## Response Format

### Success Response
```json
{
  "success": true,
  "data": { ... },
  "message": "optional message"
}
```

### Error Response
```json
{
  "success": false,
  "message": "error description",
  "code": "ERROR_CODE",
  "errors": [ ... ]
}
```

### Pagination Response
```json
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

## HTTP Status Codes

- `200 OK` - Request successful
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request parameters
- `401 Unauthorized` - Missing or invalid authentication
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `422 Unprocessable Entity` - Validation error
- `500 Internal Server Error` - Server error

---

## Authentication

### POST /auth/register

Register a new user account.

**Request Body:**
```json
{
  "name": "string (required, max 255)",
  "email": "string (required, valid email)",
  "password": "string (required, min 8 chars)"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "access_token": "eyJ0eXAiOiJKV1QiLCJhbGc...",
    "refresh_token": "eyJ0eXAiOiJKV1QiLCJhbGc...",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "creator",
      "is_active": true,
      "created_at": "2024-01-15T10:30:00Z"
    },
    "expires_in": 86400
  }
}
```

**Example cURL:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Doe",
    "email": "jane@example.com",
    "password": "SecurePassword123"
  }'
```

---

### POST /auth/login

Authenticate user and receive JWT tokens.

**Request Body:**
```json
{
  "email": "string (required, valid email)",
  "password": "string (required)"
}
```

**Response (200 OK):**
Same as register endpoint

**Example cURL:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "jane@example.com",
    "password": "SecurePassword123"
  }'
```

---

## Content Management

### POST /contents

Create new content (questions, code problems, or documentation).

**Authentication:** Required

**Request Body:**
```json
{
  "type": "question|code_problem|documentation (required)",
  "program_id": "uuid (required)",
  "topic_id": "uuid (required)",
  "subtopic_id": "uuid (required)",
  "difficulty": "easy|medium|hard (required)",
  "estimated_time_minutes": "integer 1-1440 (required)",
  "tags": ["string", "max 20 tags"],
  "data": {
    "type": "question",
    "title": "Question Title",
    "description": "Detailed description",
    "question_type": "mcq|msq|fill|short",
    "question_text": "What is 2+2?",
    "options": ["3", "4", "5"],
    "correct_options": [1]
  }
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "type": "question",
    "program_id": "550e8400-e29b-41d4-a716-446655440001",
    "topic_id": "550e8400-e29b-41d4-a716-446655440002",
    "subtopic_id": "550e8400-e29b-41d4-a716-446655440003",
    "difficulty": "medium",
    "estimated_time_minutes": 15,
    "status": "draft",
    "created_by": "550e8400-e29b-41d4-a716-446655440004",
    "current_version_id": null,
    "is_active": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

**Example cURL:**
```bash
curl -X POST http://localhost:8080/api/v1/contents \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{
    "type": "question",
    "program_id": "550e8400-e29b-41d4-a716-446655440001",
    "topic_id": "550e8400-e29b-41d4-a716-446655440002",
    "subtopic_id": "550e8400-e29b-41d4-a716-446655440003",
    "difficulty": "medium",
    "estimated_time_minutes": 15,
    "tags": ["algebra", "equations"],
    "data": {
      "type": "question",
      "title": "Solve the Equation",
      "question_type": "mcq",
      "question_text": "What is 2x + 5 = 15?",
      "options": ["x=5", "x=10", "x=15", "x=20"],
      "correct_options": [0]
    }
  }'
```

---

### GET /contents

List all content with filtering and pagination.

**Authentication:** Required

**Query Parameters:**
```
program_id: string (uuid) - Filter by program
topic_id: string (uuid) - Filter by topic
subtopic_id: string (uuid) - Filter by subtopic
difficulty: string - Filter by difficulty (easy|medium|hard)
type: string - Filter by type (question|code_problem|documentation)
status: string - Filter by status (draft|pending_review|approved|rejected)
page: integer (default: 1) - Page number
per_page: integer (default: 20, max: 100) - Items per page
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "type": "question",
      "program_id": "550e8400-e29b-41d4-a716-446655440001",
      "topic_id": "550e8400-e29b-41d4-a716-446655440002",
      "subtopic_id": "550e8400-e29b-41d4-a716-446655440003",
      "difficulty": "medium",
      "estimated_time_minutes": 15,
      "status": "approved",
      "created_by": "550e8400-e29b-41d4-a716-446655440004",
      "current_version_id": "550e8400-e29b-41d4-a716-446655440005",
      "is_active": true,
      "tags": ["algebra", "equations"],
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T11:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "per_page": 20,
    "total": 150,
    "last_page": 8
  }
}
```

**Example cURL:**
```bash
curl -X GET "http://localhost:8080/api/v1/contents?program_id=550e8400-e29b-41d4-a716-446655440001&difficulty=medium&page=1&per_page=20" \
  -H "Authorization: Bearer <access_token>"
```

---

### GET /contents/{id}

Get detailed information about specific content.

**Authentication:** Required

**URL Parameters:**
```
id: string (uuid) - Content ID
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "type": "question",
    "program_id": "550e8400-e29b-41d4-a716-446655440001",
    "topic_id": "550e8400-e29b-41d4-a716-446655440002",
    "subtopic_id": "550e8400-e29b-41d4-a716-446655440003",
    "difficulty": "medium",
    "estimated_time_minutes": 15,
    "status": "approved",
    "created_by": "550e8400-e29b-41d4-a716-446655440004",
    "current_version_id": "550e8400-e29b-41d4-a716-446655440005",
    "is_active": true,
    "tags": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440006",
        "name": "algebra",
        "is_active": true
      }
    ],
    "creator": {
      "id": "550e8400-e29b-41d4-a716-446655440004",
      "name": "Jane Doe",
      "email": "jane@example.com",
      "role": "creator",
      "is_active": true,
      "created_at": "2024-01-10T08:00:00Z"
    },
    "current_version": {
      "id": "550e8400-e29b-41d4-a716-446655440005",
      "content_id": "550e8400-e29b-41d4-a716-446655440000",
      "version_number": 1,
      "data": { ... },
      "created_by": "550e8400-e29b-41d4-a716-446655440004",
      "review_status": "approved",
      "review_comment": "Looks good!",
      "reviewed_by": "550e8400-e29b-41d4-a716-446655440007",
      "created_at": "2024-01-15T10:30:00Z"
    },
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T11:00:00Z"
  }
}
```

**Example cURL:**
```bash
curl -X GET http://localhost:8080/api/v1/contents/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer <access_token>"
```

---

### PUT /contents/{id}

Update existing content (only draft status).

**Authentication:** Required
**Authorization:** Creator of the content or admin

**URL Parameters:**
```
id: string (uuid) - Content ID
```

**Request Body:**
```json
{
  "difficulty": "easy|medium|hard (optional)",
  "estimated_time_minutes": "integer 1-1440 (optional)",
  "tags": ["string", "max 20 tags (optional)"],
  "data": {
    "type": "question",
    "title": "Updated Title",
    ...
  }
}
```

**Response (200 OK):**
Same as GET /contents/{id}

**Example cURL:**
```bash
curl -X PUT http://localhost:8080/api/v1/contents/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{
    "difficulty": "hard",
    "estimated_time_minutes": 20,
    "data": {
      "type": "question",
      "title": "Updated Question Title"
    }
  }'
```

---

### POST /contents/{id}/submit

Submit content for review (changes status from draft to pending_review).

**Authentication:** Required
**Authorization:** Creator of the content or admin

**URL Parameters:**
```
id: string (uuid) - Content ID
```

**Request Body:**
```json
{
  "version_comment": "string (optional, max 1000)"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "pending_review",
    ...
  }
}
```

**Example cURL:**
```bash
curl -X POST http://localhost:8080/api/v1/contents/550e8400-e29b-41d4-a716-446655440000/submit \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{
    "version_comment": "Ready for review, please check"
  }'
```

---

## Review Workflow

### GET /reviews/pending

Get all pending content versions awaiting review.

**Authentication:** Required
**Authorization:** Reviewer or Admin

**Query Parameters:**
```
page: integer (default: 1)
per_page: integer (default: 20, max: 100)
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440005",
      "content_id": "550e8400-e29b-41d4-a716-446655440000",
      "version_number": 2,
      "data": { ... },
      "created_by": "550e8400-e29b-41d4-a716-446655440004",
      "creator": {
        "id": "550e8400-e29b-41d4-a716-446655440004",
        "name": "Jane Doe",
        "email": "jane@example.com",
        "role": "creator"
      },
      "created_at": "2024-01-15T10:30:00Z"
    }
  ],
  "pagination": { ... }
}
```

**Example cURL:**
```bash
curl -X GET "http://localhost:8080/api/v1/reviews/pending?page=1&per_page=20" \
  -H "Authorization: Bearer <reviewer_token>"
```

---

### POST /reviews/{version_id}/approve

Approve a version (makes it the active version).

**Authentication:** Required
**Authorization:** Reviewer or Admin

**URL Parameters:**
```
version_id: string (uuid) - Version ID
```

**Request Body:**
```json
{
  "review_comment": "string (optional, max 1000)"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "version approved"
}
```

**Example cURL:**
```bash
curl -X POST http://localhost:8080/api/v1/reviews/550e8400-e29b-41d4-a716-446655440005/approve \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <reviewer_token>" \
  -d '{
    "review_comment": "Approved! Looks perfect."
  }'
```

---

### POST /reviews/{version_id}/reject

Reject a version with feedback.

**Authentication:** Required
**Authorization:** Reviewer or Admin

**URL Parameters:**
```
version_id: string (uuid) - Version ID
```

**Request Body:**
```json
{
  "review_comment": "string (required, max 1000)"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "version rejected"
}
```

**Example cURL:**
```bash
curl -X POST http://localhost:8080/api/v1/reviews/550e8400-e29b-41d4-a716-446655440005/reject \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <reviewer_token>" \
  -d '{
    "review_comment": "Please clarify the question wording and fix grammar errors on line 3"
  }'
```

---

## Programs, Topics & Subtopics

### POST /programs

Create new program (Admin only).

**Authentication:** Required
**Authorization:** Admin

**Request Body:**
```json
{
  "name": "string (required, max 255)",
  "description": "string (optional, max 1000)"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "name": "Python Fundamentals",
    "description": "Learn Python programming basics",
    "is_active": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

---

### GET /programs

List all active programs.

**Authentication:** Required

**Response (200 OK):**
```json
{
  "success": true,
  "data": [ ... ]
}
```

---

### POST /programs/{id}/topics

Create topic under a program.

**Authentication:** Required
**Authorization:** Admin

**Request Body:**
```json
{
  "name": "string (required, max 255)",
  "description": "string (optional, max 1000)"
}
```

---

### POST /topics/{id}/subtopics

Create subtopic under a topic.

**Authentication:** Required
**Authorization:** Admin

**Request Body:**
```json
{
  "name": "string (required, max 255)",
  "description": "string (optional, max 1000)"
}
```

---

## Tags

### GET /tags

List all tags.

**Authentication:** Required

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440006",
      "name": "algebra",
      "description": "Algebra topics",
      "is_active": true,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

---

### POST /tags

Create new tag (Admin only).

**Authentication:** Required
**Authorization:** Admin

**Request Body:**
```json
{
  "name": "string (required, max 255, unique)",
  "description": "string (optional, max 1000)"
}
```

**Response (201 Created):**
Same as tag object in list

---

## Error Codes

| Code | HTTP | Description |
|------|------|-------------|
| VALIDATION_ERROR | 422 | Invalid input validation |
| UNAUTHORIZED | 401 | Missing or invalid token |
| FORBIDDEN | 403 | Insufficient permissions |
| NOT_FOUND | 404 | Resource not found |
| DUPLICATE_EMAIL | 400 | Email already registered |
| INVALID_CREDENTIALS | 401 | Wrong email or password |
| INVALID_TOKEN | 401 | Token is invalid or expired |
| CONTENT_STATUS_ERROR | 400 | Content status doesn't allow operation |
| VERSION_ALREADY_REVIEWED | 400 | Version has already been reviewed |
| INTERNAL_ERROR | 500 | Server error |

---

## Rate Limiting

Authentication endpoints:
- 5 requests per minute per IP

Other endpoints:
- 100 requests per minute per user

---

## Pagination

Default pagination: 20 items per page
Max items per page: 100

Response includes:
- `page` - Current page (1-indexed)
- `per_page` - Items returned
- `total` - Total items available
- `last_page` - Last page number

---

## Sorting

Default sorting: `created_at DESC`

Add `?sort=field&order=asc|desc` to override

---

## Versioning

API versioning via URL: `/api/v1/...`

Current version: v1
Previous versions: None

---

## Webhooks

Planned for future release:
- `content.created`
- `content.submitted_for_review`
- `version.approved`
- `version.rejected`

---

For more details, see [README.md](./README.md)
