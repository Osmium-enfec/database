# Frontend Development Guide - Content Review API

## 📋 Overview

This document provides a comprehensive guide for building the frontend application for the Content Review System. It includes all required features, API endpoints, data models, and integration examples.

**API Base URL:** `http://localhost:8080` (development) or `https://api.yourdomain.com` (production)

---

## 🎯 User Roles & Workflows

### 1. **Creator** (Content Author)
- Create and draft content (MCQ, MSQ, code problems, documentation)
- Submit content for review
- View submission history
- Edit rejected content
- See approval status

### 2. **Reviewer**
- View pending content submissions
- Review content with comments
- Approve or reject submissions
- See review history

### 3. **Admin**
- All reviewer permissions
- Manage users
- View system analytics
- Configure system settings

---

## 🏗️ Frontend Architecture

```
frontend/
├── src/
│   ├── components/
│   │   ├── Auth/
│   │   │   ├── LoginForm.jsx
│   │   │   ├── RegisterForm.jsx
│   │   │   └── ProtectedRoute.jsx
│   │   ├── Layout/
│   │   │   ├── Header.jsx
│   │   │   ├── Sidebar.jsx
│   │   │   └── Dashboard.jsx
│   │   ├── Content/
│   │   │   ├── ContentList.jsx
│   │   │   ├── ContentEditor.jsx
│   │   │   ├── ContentViewer.jsx
│   │   │   └── ContentForm.jsx
│   │   ├── Review/
│   │   │   ├── ReviewList.jsx
│   │   │   ├── ReviewDetail.jsx
│   │   │   ├── ReviewApproval.jsx
│   │   │   └── ReviewRejection.jsx
│   │   └── Common/
│   │       ├── LoadingSpinner.jsx
│   │       ├── ErrorBoundary.jsx
│   │       └── NotificationToast.jsx
│   ├── pages/
│   │   ├── Home.jsx
│   │   ├── Login.jsx
│   │   ├── Register.jsx
│   │   ├── Dashboard.jsx
│   │   ├── CreateContent.jsx
│   │   ├── ContentDetail.jsx
│   │   ├── ReviewQueue.jsx
│   │   └── NotFound.jsx
│   ├── services/
│   │   ├── api.js (Axios instance with auth)
│   │   ├── authService.js
│   │   ├── contentService.js
│   │   ├── reviewService.js
│   │   └── tokenManager.js
│   ├── store/ (Redux or Context)
│   │   ├── authSlice.js
│   │   ├── contentSlice.js
│   │   ├── reviewSlice.js
│   │   └── store.js
│   ├── hooks/
│   │   ├── useAuth.js
│   │   ├── useContent.js
│   │   └── useReview.js
│   ├── types/
│   │   └── index.ts
│   ├── utils/
│   │   ├── formatting.js
│   │   ├── validation.js
│   │   └── constants.js
│   └── App.jsx
└── package.json
```

---

## 🔐 Authentication Flow

### 1. Registration
```javascript
// POST /api/v1/auth/register
const registerUser = async (credentials) => {
  const response = await fetch('http://localhost:8080/api/v1/auth/register', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      email: "user@example.com",
      password: "securePassword123",
      name: "John Doe",
      role: "creator" // "creator", "reviewer", or "admin"
    })
  });
  
  const data = await response.json();
  // data.data.access_token
  // data.data.refresh_token
  // Store both tokens securely
};
```

### 2. Login
```javascript
// POST /api/v1/auth/login
const loginUser = async (email, password) => {
  const response = await fetch('http://localhost:8080/api/v1/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password })
  });
  
  const data = await response.json();
  // Store tokens in localStorage or secure cookie
  localStorage.setItem('accessToken', data.data.access_token);
  localStorage.setItem('refreshToken', data.data.refresh_token);
  return data.data.user;
};
```

### 3. Token Management
```javascript
// All authenticated requests need Authorization header
const authenticatedRequest = async (endpoint, options = {}) => {
  const token = localStorage.getItem('accessToken');
  
  const response = await fetch(`http://localhost:8080${endpoint}`, {
    ...options,
    headers: {
      ...options.headers,
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    }
  });
  
  // Handle 401 - token expired
  if (response.status === 401) {
    // Call refresh token endpoint
    // Then retry request
  }
  
  return response.json();
};
```

### 4. Token Refresh (Auto-renewal)
```javascript
// POST /api/v1/auth/refresh
const refreshAccessToken = async () => {
  const refreshToken = localStorage.getItem('refreshToken');
  
  const response = await fetch('http://localhost:8080/api/v1/auth/refresh', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ refresh_token: refreshToken })
  });
  
  const data = await response.json();
  localStorage.setItem('accessToken', data.data.access_token);
  return data.data.access_token;
};
```

---

## 📝 Content Management APIs

### 1. Create Content
```javascript
// POST /api/v1/contents/create
const createContent = async (contentData) => {
  return authenticatedRequest('/api/v1/contents/create', {
    method: 'POST',
    body: JSON.stringify({
      title: "Understanding Variables in Programming",
      description: "A comprehensive guide to variables",
      content_type: "documentation", // "question" or "code_problem" or "documentation"
      tag_ids: ["tag1-uuid", "tag2-uuid"],
      data: {
        // Different based on content_type
        // See Content Types section below
      }
    })
  });
};
```

### 2. List Content
```javascript
// GET /api/v1/contents?page=1&per_page=20&status=draft
const listContent = async (filters = {}) => {
  const params = new URLSearchParams({
    page: filters.page || 1,
    per_page: filters.per_page || 20,
    status: filters.status || 'all', // 'draft', 'pending', 'approved', 'rejected'
    content_type: filters.contentType || 'all' // 'question', 'code_problem', 'documentation'
  });
  
  return authenticatedRequest(`/api/v1/contents?${params}`);
};
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "title": "Content Title",
      "description": "Description",
      "content_type": "question",
      "status": "draft",
      "creator_id": "uuid",
      "current_version_id": "uuid",
      "created_at": "2026-03-25T20:44:44Z",
      "updated_at": "2026-03-25T20:44:44Z"
    }
  ],
  "pagination": {
    "page": 1,
    "per_page": 20,
    "total": 50,
    "last_page": 3
  }
}
```

### 3. Get Content Detail
```javascript
// GET /api/v1/contents/{id}
const getContentDetail = async (contentId) => {
  return authenticatedRequest(`/api/v1/contents/${contentId}`);
};
```

**Response includes current version with full data:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "title": "Question Title",
    "description": "Description",
    "content_type": "question",
    "status": "pending_review",
    "creator": {
      "id": "uuid",
      "name": "John Doe",
      "email": "john@example.com"
    },
    "current_version": {
      "id": "uuid",
      "version_number": 1,
      "data": {
        "question_type": "mcq",
        "question_text": "What is 2+2?",
        "options": ["3", "4", "5", "6"],
        "correct_answer": 1,
        "explanation": "2+2 equals 4"
      },
      "created_at": "2026-03-25T20:44:44Z",
      "status": "pending"
    },
    "versions": [
      // Array of all versions
    ]
  }
}
```

### 4. Update Content
```javascript
// PUT /api/v1/contents/{id}
const updateContent = async (contentId, updates) => {
  return authenticatedRequest(`/api/v1/contents/${contentId}`, {
    method: 'PUT',
    body: JSON.stringify({
      title: "Updated Title",
      description: "Updated Description",
      data: {
        // Updated data
      }
    })
  });
};
```

### 5. Submit for Review
```javascript
// POST /api/v1/contents/{id}/submit
const submitContentForReview = async (contentId) => {
  return authenticatedRequest(`/api/v1/contents/${contentId}/submit`, {
    method: 'POST'
  });
};
```

---

## 👀 Content Review APIs

### 1. Get Pending Reviews (Reviewer Only)
```javascript
// GET /api/v1/reviews/pending?page=1&per_page=20
const getPendingReviews = async (page = 1) => {
  return authenticatedRequest(`/api/v1/reviews/pending?page=${page}&per_page=20`);
};
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "content_id": "uuid",
      "content_title": "Question Title",
      "creator_name": "John Doe",
      "version_number": 1,
      "submitted_at": "2026-03-25T20:44:44Z",
      "current_version": {
        "data": { /* content data */ }
      }
    }
  ],
  "pagination": { /* ... */ }
}
```

### 2. Approve Content Version
```javascript
// POST /api/v1/reviews/{version_id}/approve
const approveContent = async (versionId, feedback = "") => {
  return authenticatedRequest(`/api/v1/reviews/${versionId}/approve`, {
    method: 'POST',
    body: JSON.stringify({
      feedback: feedback, // Optional reviewer feedback
      notes: "Content approved - ready for publication"
    })
  });
};
```

### 3. Reject Content Version
```javascript
// POST /api/v1/reviews/{version_id}/reject
const rejectContent = async (versionId, reason) => {
  return authenticatedRequest(`/api/v1/reviews/${versionId}/reject`, {
    method: 'POST',
    body: JSON.stringify({
      rejection_reason: reason,
      feedback: "Please fix the following issues..."
    })
  });
};
```

---

## 📦 Content Types & Data Structures

### 1. Multiple Choice Question (MCQ)
```javascript
{
  question_type: "mcq",
  question_text: "What is the capital of France?",
  options: [
    "London",
    "Paris",
    "Berlin",
    "Madrid"
  ],
  correct_answer: 1, // Index of correct option (0-based)
  explanation: "Paris is the capital of France",
  difficulty: "easy", // "easy", "medium", "hard"
  tags: ["geography", "capitals"]
}
```

### 2. Multiple Select Question (MSQ)
```javascript
{
  question_type: "msq",
  question_text: "Which of these are programming languages?",
  options: [
    "Python",
    "JavaScript",
    "Docker",
    "Java"
  ],
  correct_answers: [0, 1, 3], // Indices of correct options
  explanation: "Python, JavaScript, and Java are programming languages",
  difficulty: "medium"
}
```

### 3. Fill in the Blank (FIB)
```javascript
{
  question_type: "fib",
  question_text: "The chemical formula for water is ___.",
  correct_answers: ["H2O", "H20", "water"], // Multiple acceptable answers
  explanation: "Water is composed of 2 hydrogen and 1 oxygen atom",
  case_sensitive: false
}
```

### 4. Short Answer (SA)
```javascript
{
  question_type: "sa",
  question_text: "Describe the water cycle in 2-3 sentences",
  expected_answer: "Water evaporates from surface, forms clouds, and precipitates",
  evaluation_notes: "Check for understanding of evaporation, condensation, and precipitation",
  difficulty: "medium"
}
```

### 5. Code Problem
```javascript
{
  problem_type: "coding",
  problem_title: "Reverse a String",
  problem_description: "Write a function that reverses a string",
  starter_code: `function reverseString(str) {
  // Your code here
}`,
  solution_code: `function reverseString(str) {
  return str.split('').reverse().join('');
}`,
  language: "javascript", // "python", "java", "cpp", "javascript"
  test_cases: [
    {
      input: "hello",
      expected_output: "olleh",
      hidden: false
    },
    {
      input: "world",
      expected_output: "dlrow",
      hidden: false
    }
  ],
  difficulty: "easy"
}
```

### 6. Documentation
```javascript
{
  content_type: "documentation",
  title: "Understanding Async/Await",
  sections: [
    {
      heading: "What is Async/Await?",
      content: "Async/await is syntactic sugar over promises...",
      code_snippets: [
        {
          language: "javascript",
          code: "async function fetchData() { ... }"
        }
      ]
    }
  ],
  difficulty: "medium"
}
```

---

## 🎨 UI Components to Build

### Authentication
- [ ] Login Form (email, password, remember me)
- [ ] Registration Form (email, password, name, role selection)
- [ ] Password Recovery
- [ ] Session Management

### Dashboard
- [ ] User Profile Card
- [ ] Quick Stats (created, pending, approved)
- [ ] Recent Activities
- [ ] Navigation Menu (role-based)

### Content Management
- [ ] Content List with Filtering
  - Filter by status (draft, pending, approved, rejected)
  - Filter by type (question, code problem, documentation)
  - Search by title/tags
  - Sorting options
- [ ] Content Editor
  - Title and description inputs
  - Content type selector
  - Tag selector/multi-select
  - WYSIWYG or Markdown editor for documentation
  - Question builder with option management
  - Code editor with syntax highlighting
- [ ] Content Viewer (read-only)
- [ ] Version History Viewer
  - Compare versions
  - Rollback option

### Review Interface
- [ ] Review Queue (for reviewers)
  - Pending content list
  - Sorting by submission date
- [ ] Review Detail View
  - Content display
  - Reviewer feedback input
  - Approve/Reject buttons
  - Comments section
- [ ] Review History
  - Show all reviews for a content
  - Display feedback and reviewer info

### Common
- [ ] Loading Spinner
- [ ] Error Messages
- [ ] Success Notifications
- [ ] Confirmation Dialogs
- [ ] Pagination
- [ ] Tags Display

---

## 🔄 State Management (Redux/Context)

### Auth State
```javascript
{
  user: {
    id: "uuid",
    name: "John Doe",
    email: "john@example.com",
    role: "creator", // "creator", "reviewer", "admin"
    is_active: true,
    created_at: "2026-03-25T20:44:44Z"
  },
  accessToken: "jwt_token_here",
  refreshToken: "refresh_token_here",
  isAuthenticated: true,
  isLoading: false,
  error: null
}
```

### Content State
```javascript
{
  contents: [
    // Array of content items
  ],
  currentContent: null, // Currently viewed content
  filters: {
    page: 1,
    per_page: 20,
    status: "all",
    content_type: "all",
    search: ""
  },
  pagination: {
    total: 50,
    page: 1,
    per_page: 20,
    last_page: 3
  },
  isLoading: false,
  error: null
}
```

### Review State
```javascript
{
  pendingReviews: [],
  reviewedContent: [],
  currentReview: null,
  filters: {
    page: 1,
    sort_by: "submitted_at" // "submitted_at", "title"
  },
  isLoading: false,
  error: null
}
```

---

## 🛠️ API Integration Example (Axios)

```javascript
// services/api.js
import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json'
  }
});

// Add Authorization header
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('accessToken');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Handle token refresh on 401
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;
    
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      try {
        const response = await apiClient.post('/api/v1/auth/refresh', {
          refresh_token: localStorage.getItem('refreshToken')
        });
        localStorage.setItem('accessToken', response.data.data.access_token);
        originalRequest.headers.Authorization = `Bearer ${response.data.data.access_token}`;
        return apiClient(originalRequest);
      } catch (err) {
        // Redirect to login
        window.location.href = '/login';
      }
    }
    
    return Promise.reject(error);
  }
);

export default apiClient;
```

---

## 📱 Feature Checklist

### Creator Features
- [x] Register/Login
- [x] Create content (MCQ, MSQ, FIB, SA, Code Problems, Documentation)
- [x] View own content
- [x] Edit draft content
- [x] Submit for review
- [x] View submission status
- [x] View reviewer feedback
- [x] Resubmit after rejection
- [x] See version history
- [x] Manage tags

### Reviewer Features
- [x] View pending submissions queue
- [x] Review content details
- [x] Provide feedback
- [x] Approve content
- [x] Reject with reason
- [x] See review history
- [x] Compare versions
- [x] Filter by content type

### Admin Features
- [x] All reviewer features
- [x] Manage users
- [x] View system statistics
- [x] Configure system settings

---

## 🚀 Performance Optimization

### Frontend Best Practices
1. **Lazy Loading**
   ```javascript
   const ReviewList = lazy(() => import('./components/Review/ReviewList'));
   ```

2. **Pagination**
   - Always use pagination for lists
   - Load content on demand
   - Implement infinite scroll or "Load More"

3. **Caching**
   - Cache user profile
   - Cache content tags
   - Invalidate on updates

4. **Request Debouncing**
   ```javascript
   const [search, setSearch] = useState('');
   const debouncedSearch = useCallback(
     debounce((query) => searchContents(query), 500),
     []
   );
   ```

5. **Image Optimization**
   - Use responsive images
   - Lazy load images

---

## 🔒 Security Considerations

1. **Token Storage**
   - Use httpOnly cookies or localStorage
   - Never expose in URL
   - Clear on logout

2. **HTTPS**
   - Always use HTTPS in production
   - Implement CORS properly

3. **Input Validation**
   ```javascript
   const validateEmail = (email) => {
     return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
   };
   ```

4. **XSS Prevention**
   - Sanitize HTML content
   - Use `textContent` instead of `innerHTML`

5. **CSRF Protection**
   - Include CSRF tokens if needed
   - Use SameSite cookie attribute

---

## 📊 Error Handling

```javascript
// API responses follow this format
{
  "success": false,
  "message": "Error description",
  "errors": {
    "field_name": ["validation error 1", "validation error 2"]
  }
}
```

### Common Error Codes
- `400` - Bad Request (validation errors)
- `401` - Unauthorized (invalid/expired token)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found
- `422` - Unprocessable Entity
- `500` - Server Error

```javascript
const handleApiError = (error) => {
  if (error.response?.status === 401) {
    // Redirect to login
  } else if (error.response?.status === 403) {
    // Show permission denied
  } else if (error.response?.data?.errors) {
    // Show validation errors
    showFieldErrors(error.response.data.errors);
  } else {
    // Show generic error
    showNotification(error.message, 'error');
  }
};
```

---

## 📚 Environment Configuration

Create `.env` file:
```bash
REACT_APP_API_URL=http://localhost:8080
REACT_APP_ENVIRONMENT=development
REACT_APP_LOG_LEVEL=debug
```

For production:
```bash
REACT_APP_API_URL=https://api.yourdomain.com
REACT_APP_ENVIRONMENT=production
REACT_APP_LOG_LEVEL=error
```

---

## 🧪 Testing Strategy

```javascript
// Example test for authentication
describe('Authentication', () => {
  test('should login successfully with valid credentials', async () => {
    const response = await authService.login('user@example.com', 'password');
    expect(response.data.access_token).toBeDefined();
    expect(localStorage.getItem('accessToken')).toBeDefined();
  });

  test('should refresh token automatically on 401', async () => {
    // Mock API to return 401
    // Verify refresh token is called
    // Verify request is retried
  });
});

describe('Content Management', () => {
  test('should create content successfully', async () => {
    const content = await contentService.create({
      title: 'Test',
      content_type: 'question',
      data: { /* ... */ }
    });
    expect(content.id).toBeDefined();
  });
});
```

---

## 📞 Support & Resources

- **API Documentation**: See [API_REFERENCE.md](API_REFERENCE.md)
- **Backend Setup**: See [QUICK_START.md](QUICK_START.md)
- **Architecture**: See [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)
- **Deployment**: See [DEPLOYMENT.md](DEPLOYMENT.md)

---

## 🎓 Example: Complete Create Content Workflow

```javascript
// 1. Creator navigates to Create Content
<Route path="/create-content" element={<CreateContent />} />

// 2. Form collects data
const handleCreateContent = async (formData) => {
  const content = await contentService.create({
    title: formData.title,
    description: formData.description,
    content_type: formData.type, // "question", "code_problem", "documentation"
    tag_ids: formData.tags,
    data: formData.data // Question/code/doc specific data
  });
  
  // 3. Save to draft (auto-save)
  // 4. Navigate to edit page
  navigate(`/content/${content.id}`);
};

// 5. Creator can edit multiple times
// 6. When ready, submit for review
const handleSubmitForReview = async (contentId) => {
  await contentService.submitForReview(contentId);
  showNotification('Content submitted for review');
};

// 7. Reviewer sees in pending list
// 8. Reviewer approves or rejects
const handleApproveContent = async (versionId) => {
  await reviewService.approve(versionId, feedback);
  showNotification('Content approved');
};

// 9. Creator can see status and proceed
```

---

**Last Updated**: March 26, 2026  
**Version**: 1.0  
**Status**: Production Ready
