#!/bin/bash

# Savor Server Deployment Script
# This script helps with deployment preparation and common deployment tasks

set -e

echo "🚀 Savor Server Deployment Helper"
echo "================================="

# Function to check if git is clean
check_git_status() {
    if ! git diff --quiet HEAD; then
        echo "⚠️  You have uncommitted changes. Please commit or stash them first."
        echo "Run: git add . && git commit -m 'Prepare for deployment'"
        exit 1
    fi
}

# Function to validate environment variables
validate_env() {
    echo "🔍 Validating environment setup..."
    
    # Check for required files
    if [ ! -f "go.mod" ]; then
        echo "❌ go.mod not found. Are you in the correct directory?"
        exit 1
    fi
    
    if [ ! -f "main.go" ]; then
        echo "❌ main.go not found. Are you in the correct directory?"
        exit 1
    fi
    
    echo "✅ Project structure looks good!"
}

# Function to test local build
test_build() {
    echo "🔨 Testing local build..."
    if go build -o savor-server-test .; then
        echo "✅ Build successful!"
        rm -f savor-server-test
    else
        echo "❌ Build failed. Please fix the errors before deploying."
        exit 1
    fi
}

# Function to show deployment checklist
show_checklist() {
    echo ""
    echo "📋 Pre-deployment Checklist:"
    echo "============================"
    echo "□ Firebase project created and service account downloaded"
    echo "□ Stripe account set up with API keys"
    echo "□ Google Maps API key obtained"
    echo "□ Railway account created"
    echo "□ Repository pushed to GitHub"
    echo ""
    echo "🎯 Next Steps:"
    echo "1. Go to https://railway.app/dashboard"
    echo "2. Click 'New Project' → 'Deploy from GitHub repo'"
    echo "3. Select your repository"
    echo "4. Add PostgreSQL service"
    echo "5. Set environment variables (see DEPLOYMENT.md)"
    echo ""
}

# Function to generate environment variables template
generate_env_template() {
    echo "📝 Generating environment variables template..."
    cat > .env.railway << 'EOF'
# Firebase Configuration (Required)
FIREBASE_PROJECT_ID=your-firebase-project-id
FIREBASE_PRIVATE_KEY_ID=your-private-key-id
FIREBASE_PRIVATE_KEY="-----BEGIN PRIVATE KEY-----\nYOUR_PRIVATE_KEY_HERE\n-----END PRIVATE KEY-----\n"
FIREBASE_CLIENT_EMAIL=firebase-adminsdk-xxxxx@your-project.iam.gserviceaccount.com
FIREBASE_CLIENT_ID=your-client-id
FIREBASE_CLIENT_X509_CERT_URL=https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-xxxxx%40your-project.iam.gserviceaccount.com

# Stripe Configuration (Required)
STRIPE_SECRET_KEY=sk_test_your_stripe_secret_key
STRIPE_PUBLISHABLE_KEY=pk_test_your_stripe_publishable_key

# Google Maps Configuration (Required)
GOOGLE_MAPS_API_KEY=your_google_maps_api_key

# Application Configuration
GIN_MODE=release
SESSION_SECRET=your-session-secret-key-at-least-32-characters

# Database (Automatically set by Railway)
# DATABASE_URL=postgresql://...
# PORT=8080
EOF
    echo "✅ Environment template created: .env.railway"
    echo "   Edit this file with your actual values and copy to Railway dashboard"
}

# Main script
case "$1" in
    "check")
        validate_env
        check_git_status
        test_build
        echo "✅ All checks passed! Ready for deployment."
        ;;
    "env")
        generate_env_template
        ;;
    "prepare")
        validate_env
        test_build
        echo "📦 Preparing for deployment..."
        if [ -n "$(git status --porcelain)" ]; then
            echo "📝 Committing changes..."
            git add .
            git commit -m "Prepare for deployment - $(date)"
        fi
        echo "📤 Pushing to GitHub..."
        git push origin main
        echo "✅ Ready for Railway deployment!"
        show_checklist
        ;;
    "help"|"")
        echo "Usage: $0 [command]"
        echo ""
        echo "Commands:"
        echo "  check    - Validate project and check if ready for deployment"
        echo "  env      - Generate environment variables template"
        echo "  prepare  - Prepare and push code for deployment"
        echo "  help     - Show this help message"
        echo ""
        echo "For detailed deployment instructions, see DEPLOYMENT.md"
        ;;
    *)
        echo "Unknown command: $1"
        echo "Run '$0 help' for usage information"
        exit 1
        ;;
esac 