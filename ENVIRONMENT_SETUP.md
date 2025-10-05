# Environment Configuration for Notifications

## Required Environment Variables

### Email Notifications (Gmail SMTP)
```bash
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your_email@gmail.com
SMTP_PASSWORD=your_app_password_here
```

**Note:** For Gmail, you need to:
1. Enable 2-factor authentication
2. Generate an App Password (not your regular Gmail password)
3. Use the App Password as SMTP_PASSWORD

### SMS Notifications (Twilio)
```bash
TWILIO_ACCOUNT_SID=your_twilio_account_sid
TWILIO_AUTH_TOKEN=your_twilio_auth_token
TWILIO_PHONE_NUMBER=+1234567890
```

**Setup Steps:**
1. Sign up at https://www.twilio.com/
2. Get your Account SID and Auth Token from the Console
3. Purchase a phone number for sending SMS

### Existing Variables
```bash
# Database, Firebase, Stripe, Google Maps, etc.
DATABASE_URL=your_database_url
FIREBASE_PROJECT_ID=your_project_id
FIREBASE_API_KEY=your_firebase_web_api_key  # Required for email/password authentication
STRIPE_SECRET_KEY=your_stripe_key
GOOGLE_MAPS_API_KEY=your_maps_key
SESSION_SECRET=your_session_secret
FRONTEND_URL=https://your-frontend-url.com
```

**Note:** The `FIREBASE_API_KEY` is your Firebase project's Web API Key, which can be found in:
1. Firebase Console → Project Settings → General → Web API Key
2. This is different from the service account credentials and is required for email/password authentication

## Testing Notifications

1. **Email Testing**: Send a test email to verify SMTP configuration
2. **SMS Testing**: Send a test SMS to verify Twilio configuration
3. **Integration Testing**: Create a reservation to test full flow

## Security Notes

- Never commit .env files to version control
- Use strong, unique passwords and tokens
- Rotate credentials regularly
- Use environment-specific configurations for development/production
