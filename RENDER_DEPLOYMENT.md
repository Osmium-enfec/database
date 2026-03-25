# Deploying to Render

This guide explains how to deploy the Content Review API to Render.

## Services Required

You need to create **TWO** services on Render:

1. **Web Service** - for the Go backend API
2. **Postgres** - for the database

## Step-by-Step Deployment

### Step 1: Create Postgres Service

1. Go to [Render Dashboard](https://dashboard.render.com)
2. Click **"New +"** → **"Postgres"**
3. Configure:
   - **Name**: `content-review-db`
   - **Database**: `content_review`
   - **User**: `postgres`
   - **Region**: Choose closest to your users
   - Leave other settings as default
4. Click **"Create Database"**
5. Wait for the database to be ready (status shows "Available")
6. **Copy the Internal Database URL** (you'll need this for the Web Service)

### Step 2: Create Web Service

1. Click **"New +"** → **"Web Service"**
2. Connect your GitHub repository:
   - Select `Osmium-enfec/database`
   - Choose `main` branch
3. Configure:
   - **Name**: `content-review-api`
   - **Environment**: `Docker` (Render auto-detects Dockerfile)
   - **Region**: Same as database
   - **Plan**: Free tier is fine for testing
4. Click **"Advanced"** and add Environment Variables:

   ```
   ENVIRONMENT=production
   DATABASE_URL=[Internal Database URL from Step 1]
   JWT_SECRET=[Generate a random 32+ character string]
   JWT_REFRESH_SECRET=[Generate a random 32+ character string]
   SERVER_PORT=8080
   ```

5. Click **"Create Web Service"**
6. Wait for deployment to complete

### Step 3: Verify Deployment

Once deployed, your API will be accessible at:
```
https://content-review-api.onrender.com
```

Test with:
```bash
curl https://content-review-api.onrender.com/health
```

Access Swagger UI:
```
https://content-review-api.onrender.com/swagger/index.html
```

## Environment Variables

### Required for Render:
- `DATABASE_URL` - Automatically provided by Render's Postgres service
- `JWT_SECRET` - Must be 32+ characters (generate random string)
- `JWT_REFRESH_SECRET` - Must be 32+ characters (generate random string)

### Optional:
- `ENVIRONMENT` - Set to `production`
- `JWT_EXPIRY_HOURS` - Default: 720 (30 days)
- `DB_MAX_CONN` - Default: 25

## Generating Secure Secrets

Generate secure random strings for JWT secrets:

```bash
# On macOS/Linux
openssl rand -base64 32

# Or use this online: https://www.random.org/strings/
```

## Database Migrations

Render will run the Dockerfile which builds your Go app. The database schema needs to be initialized.

You can either:
1. **Run migrations manually** after deployment (connect to Render Postgres and run the SQL)
2. **Add auto-migration code** to your Go application startup

## Troubleshooting

### Service fails to build
- Check the build logs in Render dashboard
- Ensure `Dockerfile` is present and valid
- Verify Go version compatibility (Go 1.21)

### Database connection fails
- Verify `DATABASE_URL` is set correctly
- Check database is in "Available" state
- Ensure Web Service is linked to correct Postgres service

### API returns 500 errors
- Check logs in Render dashboard
- Verify all environment variables are set
- Check database schema is initialized

## Notes

- **Render provides free tier** with limitations (may sleep during inactivity)
- **For production**, upgrade to paid plan for guaranteed uptime
- **Database backups** are included in Render Postgres
- **SSL/TLS** is automatically enabled for all Render services
