# Frontend API Documentation

## Authentication
All endpoints except `/dropdown/*` require Bearer token in Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

---

## 1. BULK CREATE CONTENT

### Endpoint
```
POST /api/v1/contents/bulk
```

### Headers
```
Content-Type: application/json
Authorization: Bearer <jwt_token>
```

### Request Body
```json
{
  "contents": [
    {
      "type": "question",
      "program_id": "550e8400-e29b-41d4-a716-446655440000",
      "topic_id": "550e8400-e29b-41d4-a716-446655440001",
      "subtopic_id": "550e8400-e29b-41d4-a716-446655440002",
      "difficulty": "easy",
      "estimated_time_minutes": 5,
      "tags": ["algebra", "basics"],
      "data": {
        "type": "question",
        "title": "Simple Addition",
        "description": "Solve simple addition",
        "question_type": "mcq",
        "question_text": "What is 2+2?",
        "hints": ["Count on your fingers", "It's between 3 and 5", "The answer is even"],
        "options": ["3", "4", "5", "6"],
        "correct_options": [1]
      }
    },
    {
      "type": "code_problem",
      "program_id": "550e8400-e29b-41d4-a716-446655440000",
      "topic_id": "550e8400-e29b-41d4-a716-446655440003",
      "subtopic_id": "550e8400-e29b-41d4-a716-446655440004",
      "difficulty": "medium",
      "estimated_time_minutes": 20,
      "tags": ["python", "strings"],
      "data": {
        "type": "code_problem",
        "title": "Reverse a String",
        "description": "Write a function to reverse a string",
        "code_problem_data": {
          "starter_code": "def reverse_string(s):\n    pass",
          "solution_code": "def reverse_string(s):\n    return s[::-1]",
          "test_cases": [
            {
              "input": "'hello'",
              "expected_output": "'olleh'",
              "is_hidden": false
            }
          ]
        }
      }
    },
    {
      "type": "documentation",
      "program_id": "550e8400-e29b-41d4-a716-446655440000",
      "topic_id": "550e8400-e29b-41d4-a716-446655440005",
      "subtopic_id": "550e8400-e29b-41d4-a716-446655440006",
      "difficulty": "easy",
      "estimated_time_minutes": 10,
      "tags": ["guide"],
      "data": {
        "type": "documentation",
        "title": "Python Basics Guide",
        "description": "Introduction to Python",
        "documentation_data": {
          "markdown_content": "# Python Basics\n\n## Variables\nPython variables..."
        }
      }
    }
  ]
}
```

### Response (Success 200)
```json
{
  "success": true,
  "data": {
    "created_count": 3,
    "failed_count": 0,
    "created_ids": [
      "550e8400-e29b-41d4-a716-446655440100",
      "550e8400-e29b-41d4-a716-446655440101",
      "550e8400-e29b-41d4-a716-446655440102"
    ],
    "errors": []
  },
  "message": "Bulk content creation completed"
}
```

### Response (Partial Success 207)
```json
{
  "success": false,
  "data": {
    "created_count": 2,
    "failed_count": 1,
    "created_ids": [
      "550e8400-e29b-41d4-a716-446655440100",
      "550e8400-e29b-41d4-a716-446655440101"
    ],
    "errors": [
      {
        "index": 2,
        "error": "Invalid difficulty level"
      }
    ]
  },
  "message": "Bulk content creation completed with errors"
}
```

### Response (Error 400)
```json
{
  "success": false,
  "message": "Invalid request",
  "code": "VALIDATION_ERROR",
  "errors": [
    "contents array is required",
    "contents array cannot be empty",
    "maximum 100 items per request"
  ]
}
```

### Validation Rules
- **Min items**: 1
- **Max items**: 100 per request
- **Question hints**: Minimum 3 hints required
- **Code test cases**: Minimum 1 test case required
- **All UUIDs**: Must be valid UUID format
- **Time**: 1-1440 minutes
- **Difficulty**: "easy", "medium", or "hard"
- **Tags**: Max 20 per content

---

## 2. GET PROGRAMS (Dropdown)

### Endpoint
```
GET /api/v1/dropdown/programs
```

### Headers
```
Content-Type: application/json
```

### Query Parameters
None

### Response (200)
```json
{
  "success": true,
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Python Programming",
      "description": "Complete Python learning roadmap from beginner to senior developer level"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440010",
      "name": "JavaScript Mastery",
      "description": "From basics to advanced JavaScript"
    }
  ],
  "message": "Programs retrieved successfully"
}
```

### Response (Error 500)
```json
{
  "success": false,
  "message": "Failed to fetch programs",
  "code": "DATABASE_ERROR"
}
```

---

## 3. GET TOPICS BY PROGRAM (Dropdown)

### Endpoint
```
GET /api/v1/dropdown/topics?program_id=<uuid>
```

### Headers
```
Content-Type: application/json
```

### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| program_id | UUID | Yes | The program ID to get topics for |

### Example Request
```
GET /api/v1/dropdown/topics?program_id=550e8400-e29b-41d4-a716-446655440000
```

### Response (200)
```json
{
  "success": true,
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "name": "Foundations",
      "description": "Programming basics and environment setup"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440002",
      "name": "Core Basics",
      "description": "Variables, operators, conditionals, loops"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440003",
      "name": "Data Structures",
      "description": "Strings, Lists, Tuples, Sets, Dictionaries"
    }
  ],
  "message": "Topics retrieved successfully"
}
```

### Response (Error 400)
```json
{
  "success": false,
  "message": "program_id is required",
  "code": "MISSING_PARAMETER"
}
```

### Response (Error 404)
```json
{
  "success": false,
  "message": "Program not found",
  "code": "NOT_FOUND"
}
```

---

## 4. GET SUBTOPICS BY TOPIC (Dropdown)

### Endpoint
```
GET /api/v1/dropdown/subtopics?topic_id=<uuid>
```

### Headers
```
Content-Type: application/json
```

### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| topic_id | UUID | Yes | The topic ID to get subtopics for |

### Example Request
```
GET /api/v1/dropdown/subtopics?topic_id=550e8400-e29b-41d4-a716-446655440001
```

### Response (200)
```json
{
  "success": true,
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440100",
      "name": "What is Programming",
      "description": "How computers execute code, Interpreted vs compiled languages"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440101",
      "name": "Setting Up Environment",
      "description": "Installing Python, VS Code, Virtual environments"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440102",
      "name": "First Program",
      "description": "print(), Comments, Basic syntax"
    }
  ],
  "message": "Subtopics retrieved successfully"
}
```

### Response (Error 400)
```json
{
  "success": false,
  "message": "topic_id is required",
  "code": "MISSING_PARAMETER"
}
```

### Response (Error 404)
```json
{
  "success": false,
  "message": "Topic not found",
  "code": "NOT_FOUND"
}
```

---

## Frontend Integration Examples

### JavaScript/React Example - Bulk Upload

```javascript
async function bulkUploadContent(contents) {
  try {
    const response = await fetch('/api/v1/contents/bulk', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${getToken()}`
      },
      body: JSON.stringify({ contents })
    });

    const result = await response.json();
    
    if (response.ok) {
      console.log(`✅ Created ${result.data.created_count} items`);
      return result.data.created_ids;
    } else if (response.status === 207) {
      console.warn(`⚠️ Created ${result.data.created_count}, Failed ${result.data.failed_count}`);
      console.error('Errors:', result.data.errors);
      return result.data.created_ids;
    } else {
      throw new Error(result.message);
    }
  } catch (error) {
    console.error('Bulk upload failed:', error);
    throw error;
  }
}
```

### JavaScript/React Example - Dropdown Flow

```javascript
// 1. Load programs on component mount
useEffect(() => {
  loadPrograms();
}, []);

async function loadPrograms() {
  const response = await fetch('/api/v1/dropdown/programs');
  const result = await response.json();
  setPrograms(result.data);
}

// 2. When user selects program
async function handleProgramChange(programId) {
  setSelectedProgram(programId);
  const response = await fetch(
    `/api/v1/dropdown/topics?program_id=${programId}`
  );
  const result = await response.json();
  setTopics(result.data);
  setSelectedTopic(null); // Reset topic
  setSubtopics([]); // Reset subtopics
}

// 3. When user selects topic
async function handleTopicChange(topicId) {
  setSelectedTopic(topicId);
  const response = await fetch(
    `/api/v1/dropdown/subtopics?topic_id=${topicId}`
  );
  const result = await response.json();
  setSubtopics(result.data);
  setSelectedSubtopic(null); // Reset subtopic
}
```

### Form Structure for Single Content Creation

```javascript
const contentForm = {
  type: 'question', // 'question' | 'code_problem' | 'documentation'
  program_id: selectedProgram,
  topic_id: selectedTopic,
  subtopic_id: selectedSubtopic,
  difficulty: 'medium', // 'easy' | 'medium' | 'hard'
  estimated_time_minutes: 15,
  tags: [],
  data: {
    type: 'question',
    title: '',
    description: '',
    question_type: 'mcq', // 'mcq' | 'msq' | 'fill' | 'short'
    question_text: '',
    hints: [], // Minimum 3
    options: [],
    correct_options: []
  }
};
```

---

## Error Status Codes

| Code | Meaning |
|------|---------|
| 200 | Success |
| 207 | Partial success (for bulk operations) |
| 400 | Bad request / Validation error |
| 401 | Unauthorized (missing/invalid token) |
| 403 | Forbidden (insufficient permissions) |
| 404 | Not found |
| 500 | Server error |

---

## Rate Limiting Notes
- Bulk upload: Max 100 items per request
- Dropdown endpoints: No rate limit
- Requests: Standard HTTP rate limiting applies

---

## Testing with cURL

### Bulk Upload
```bash
curl -X POST http://localhost:8080/api/v1/contents/bulk \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "contents": [
      {
        "type": "question",
        "program_id": "YOUR_PROGRAM_ID",
        "topic_id": "YOUR_TOPIC_ID",
        "subtopic_id": "YOUR_SUBTOPIC_ID",
        "difficulty": "easy",
        "estimated_time_minutes": 5,
        "tags": ["test"],
        "data": {
          "type": "question",
          "title": "Test Question",
          "description": "Test",
          "question_type": "mcq",
          "question_text": "What is 1+1?",
          "hints": ["Hint 1", "Hint 2", "Hint 3"],
          "options": ["1", "2", "3"],
          "correct_options": [1]
        }
      }
    ]
  }'
```

### Get Programs
```bash
curl -X GET http://localhost:8080/api/v1/dropdown/programs \
  -H "Content-Type: application/json"
```

### Get Topics
```bash
curl -X GET "http://localhost:8080/api/v1/dropdown/topics?program_id=YOUR_PROGRAM_ID" \
  -H "Content-Type: application/json"
```

### Get Subtopics
```bash
curl -X GET "http://localhost:8080/api/v1/dropdown/subtopics?topic_id=YOUR_TOPIC_ID" \
  -H "Content-Type: application/json"
```

