# Deployment Guide for Savor Server

## Railway Deployment (Recommended - Easiest & Cheapest)

Railway is the recommended deployment platform for this Go application because it offers:
- Simple Git-based deployment
- Built-in PostgreSQL database
- Environment variable management
- Automatic builds and deployments
- Generous free tier ($5/month credit)

### Prerequisites

1. GitHub account with your code repository
2. Railway account (sign up at [railway.app](https://railway.app))
3. Required API keys (Firebase, Stripe, Google Maps)

### Step 1: Deploy to Railway

1. **Connect Repository:**
   - Go to [Railway Dashboard](https://railway.app/dashboard)
   - Click "New Project"
   - Select "Deploy from GitHub repo"
   - Choose your `savor-server` repository

2. **Add PostgreSQL Database:**
   - In your Railway project, click "Add Service"
   - Select "Database" → "PostgreSQL"
   - Railway will automatically provision a PostgreSQL database
   - The `DATABASE_URL` will be automatically set

### Step 2: Configure Environment Variables

In your Railway project dashboard, go to "Variables" and add these environment variables:

#### Required Variables:

**Firebase Configuration:**
```
FIREBASE_PROJECT_ID=your-firebase-project-id
FIREBASE_PRIVATE_KEY_ID=your-private-key-id
FIREBASE_PRIVATE_KEY="-----BEGIN PRIVATE KEY-----\nYOUR_PRIVATE_KEY\n-----END PRIVATE KEY-----\n"
FIREBASE_CLIENT_EMAIL=your-service-account@your-project.iam.gserviceaccount.com
FIREBASE_CLIENT_ID=your-client-id
FIREBASE_AUTH_URI=https://accounts.google.com/o/oauth2/auth
FIREBASE_TOKEN_URI=https://oauth2.googleapis.com/token
FIREBASE_AUTH_PROVIDER_X509_CERT_URL=https://www.googleapis.com/oauth2/v1/certs
FIREBASE_CLIENT_X509_CERT_URL=https://www.googleapis.com/robot/v1/metadata/x509/your-service-account%40your-project.iam.gserviceaccount.com
```

**Stripe Configuration:**
```
STRIPE_SECRET_KEY=sk_test_your_stripe_secret_key
STRIPE_PUBLISHABLE_KEY=pk_test_your_stripe_publishable_key
```

**Google Maps Configuration:**
```
GOOGLE_MAPS_API_KEY=your_google_maps_api_key
```

**Application Configuration:**
```
GIN_MODE=release
SESSION_SECRET=your-session-secret-key-at-least-32-characters
```

#### Automatic Variables (Set by Railway):
- `DATABASE_URL` - Automatically configured when you add PostgreSQL
- `PORT` - Automatically set by Railway

### Step 3: Deploy

1. Push your code to GitHub
2. Railway will automatically build and deploy your application
3. Your app will be available at `https://your-app-name.railway.app`

### Step 4: Test Deployment

Once deployed, test these endpoints:
- Health check: `https://your-app-name.railway.app/api/health`
- Swagger docs: `https://your-app-name.railway.app/swagger/index.html`

### Step 5: Update CORS Configuration

Update the CORS configuration in your code to allow your frontend domain:
```go
r.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"https://your-frontend-domain.com", "http://localhost:3000"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
}))
```

## Alternative Deployment Options

### Docker Deployment
Use the included `Dockerfile` to deploy to any Docker-compatible platform:
```bash
docker build -t savor-server .
docker run -p 8080:8080 --env-file .env savor-server
```

### Manual Server Deployment
1. Build the application: `go build -o savor-server`
2. Set environment variables
3. Run: `./savor-server`

## Database Setup

### Railway PostgreSQL (Automatic)
When you add PostgreSQL to Railway, it automatically:
- Creates a database instance
- Sets the `DATABASE_URL` environment variable
- Handles SSL connections

### Manual Database Setup
If using a different database provider:
1. Create a PostgreSQL database
2. Set these environment variables:
   ```
   DB_HOST=your-db-host
   DB_PORT=5432
   DB_USER=your-db-user
   DB_PASSWORD=your-db-password
   DB_NAME=your-db-name
   DB_SSLMODE=require
   ```

## Cost Estimation

### Railway Pricing:
- **Free Tier**: $5/month credit (usually sufficient for development)
- **Pro Plan**: $20/month (includes more resources)
- **PostgreSQL**: ~$5-10/month depending on usage

### Total Monthly Cost:
- **Development**: Free (using $5 credit)
- **Production**: ~$20-30/month

## Security Considerations

1. **Environment Variables**: Never commit sensitive data to version control
2. **HTTPS**: Railway provides HTTPS automatically
3. **Database**: Use SSL connections (Railway handles this)
4. **API Keys**: Use environment variables for all API keys
5. **Session Security**: Use a strong session secret (32+ characters)

## Monitoring

Railway provides:
- Application logs
- Metrics dashboard
- Health checks
- Deployment history

## Support

If you encounter issues:
1. Check Railway logs in the dashboard
2. Verify all environment variables are set
3. Test database connection
4. Check API key permissions 