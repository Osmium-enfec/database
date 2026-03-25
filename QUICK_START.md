# Quick Start Guide

## ⚡ 5-Minute Setup

### Prerequisites
- Docker & Docker Compose
- OR: Go 1.21+, PostgreSQL 15

### Option A: Docker (Recommended - 2 minutes)

```bash
cd content-review-api
docker-compose up -d
sleep 10
curl http://localhost:8080/health
```

✅ Done! API running on http://localhost:8080

### Option B: Local Development (5 minutes)

```bash
# 1. Install Go 1.21+
# 2. Install PostgreSQL 15
# 3. Clone repository
cd content-review-api

# 4. Setup
make setup

# 5. Edit .env with database connection
nano .env

# 6. Run migrations
make db-migrate

# 7. Start application
make run
```

✅ Done! API running on http://localhost:8080

---

## 🧪 Test the API

### 1. Health Check
```bash
curl http://localhost:8080/health
# Response: {"status": "ok"}
```

### 2. Register User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Doe",
    "email": "jane@example.com",
    "password": "SecurePassword123"
  }'

# Save the access_token from response
# export TOKEN="<access_token>"
```

### 3. Create Content
```bash
curl -X POST http://localhost:8080/api/v1/contents \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "question",
    "program_id": "550e8400-e29b-41d4-a716-446655440001",
    "topic_id": "550e8400-e29b-41d4-a716-446655440002",
    "subtopic_id": "550e8400-e29b-41d4-a716-446655440003",
    "difficulty": "medium",
    "estimated_time_minutes": 15,
    "tags": ["algebra"],
    "data": {
      "type": "question",
      "title": "Solve Equation",
      "question_type": "mcq",
      "question_text": "What is 2+2?",
      "options": ["3", "4", "5"],
      "correct_options": [1]
    }
  }'
```

### 4. List Content
```bash
curl -X GET http://localhost:8080/api/v1/contents \
  -H "Authorization: Bearer $TOKEN"
```

---

## 🛠️ Useful Commands

```bash
# View logs
docker-compose logs -f api

# Stop services
docker-compose down

# Restart
docker-compose restart api

# Database shell
docker-compose exec postgres psql -U postgres -d content_review

# API logs
docker-compose logs api | tail -20
```

---

## 📚 Documentation

- **Full README**: [README.md](./README.md)
- **API Reference**: [API_REFERENCE.md](./API_REFERENCE.md)
- **Deployment**: [DEPLOYMENT.md](./DEPLOYMENT.md)
- **Project Summary**: [PROJECT_SUMMARY.md](./PROJECT_SUMMARY.md)

---

## 🔑 Key Endpoints

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | /auth/register | - | Create account |
| POST | /auth/login | - | Login |
| POST | /contents | ✅ | Create content |
| GET | /contents | ✅ | List content |
| GET | /contents/{id} | ✅ | Get content |
| PUT | /contents/{id} | ✅ | Update content |
| POST | /contents/{id}/submit | ✅ | Submit for review |
| GET | /reviews/pending | ✅ | Pending reviews |
| POST | /reviews/{version}/approve | ✅ | Approve |
| POST | /reviews/{version}/reject | ✅ | Reject |

---

## 🐛 Troubleshooting

### API won't start
```bash
# Check logs
docker-compose logs api

# Verify database is running
docker-compose logs postgres

# Restart everything
docker-compose down && docker-compose up -d
```

### Database connection error
```bash
# Check database is healthy
docker-compose exec postgres pg_isready -U postgres

# Check credentials in .env
grep DB_ .env
```

### Port already in use
```bash
# Change port in docker-compose.yml
# Or kill process using port 8080
sudo lsof -ti:8080 | xargs kill -9
```

---

## 📝 Test Flow

1. **Register** → Get tokens
2. **Create Content** → Draft status
3. **Submit** → Pending review
4. **Approve** (as reviewer) → Active
5. **View** → See published content

---

## 💡 Pro Tips

- **Keep tokens handy**: Copy access_token to terminal variable
- **Read error messages**: API returns descriptive errors
- **Check documentation**: Full API reference in API_REFERENCE.md
- **View migrations**: Database schema in migrations/001_initial_schema.sql
- **Use Makefile**: Common tasks: `make help`

---

## 🚀 Next Steps

1. ✅ Verify it's running
2. 📖 Read full [README.md](./README.md)
3. 🔑 Understand authentication flow
4. 📝 Review API endpoints in [API_REFERENCE.md](./API_REFERENCE.md)
5. 🛠️ Customize for your needs
6. 📦 Deploy with [DEPLOYMENT.md](./DEPLOYMENT.md)

---

## 📞 Need Help?

- **API issues**: See [API_REFERENCE.md](./API_REFERENCE.md#error-codes)
- **Deployment**: See [DEPLOYMENT.md](./DEPLOYMENT.md#troubleshooting)
- **Architecture**: See [PROJECT_SUMMARY.md](./PROJECT_SUMMARY.md)
- **Setup**: See [README.md](./README.md#getting-started)

---

**You're all set! The API is production-ready.** 🎉
