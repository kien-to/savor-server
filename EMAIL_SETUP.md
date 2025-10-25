# Email & SMS Configuration Setup

The Savor backend supports sending email and SMS confirmations for reservations.

## Architecture

- **Email Service** (`services/email.go`): Handles email confirmations with beautiful HTML templates
- **Notification Service** (`services/notifications.go`): Handles SMS confirmations via Twilio
- **Note**: The old email functionality in `notifications.go` is deprecated - all emails use the new service

## Email Configuration

To enable email confirmations, configure SMTP settings:

## Required Environment Variables

Add the following variables to your `.env` file:

```env
# Email Configuration (SMTP)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your_email@gmail.com
SMTP_PASSWORD=your_app_password
FROM_EMAIL=noreply@savor.com
FROM_NAME=Savor
```

## Email Service Options

### Option 1: Gmail (For Development/Testing)

1. Go to Google Account settings
2. Enable 2-Factor Authentication
3. Generate an App Password:
   - Go to Security → 2-Step Verification → App passwords
   - Select "Mail" and your device
   - Copy the generated password

4. Configure:
```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your.gmail@gmail.com
SMTP_PASSWORD=your_16_character_app_password
FROM_EMAIL=your.gmail@gmail.com
FROM_NAME=Savor
```

### Option 2: SendGrid (Recommended for Production)

1. Sign up at https://sendgrid.com
2. Create an API Key from Settings → API Keys
3. Configure:
```env
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USER=apikey
SMTP_PASSWORD=your_sendgrid_api_key
FROM_EMAIL=noreply@yourdomain.com
FROM_NAME=Savor
```

### Option 3: AWS SES (For Production at Scale)

1. Set up AWS SES and verify your domain
2. Create SMTP credentials
3. Configure:
```env
SMTP_HOST=email-smtp.us-east-1.amazonaws.com
SMTP_PORT=587
SMTP_USER=your_aws_smtp_username
SMTP_PASSWORD=your_aws_smtp_password
FROM_EMAIL=noreply@yourdomain.com
FROM_NAME=Savor
```

## Email Features

When configured, the system will automatically send email confirmations containing:
- Store information and image
- Pickup time and location  
- Reservation details (quantity, price, reservation ID)
- Order status
- Savings information

## Testing

To test if email is working:
1. Create a reservation (authenticated or guest)
2. Check the server logs for:
   ```
   Email service initialized successfully
   Email confirmation sent successfully to email@example.com for reservation xxx
   ```
3. Check the recipient's inbox

## SMS Configuration (Optional)

To enable SMS notifications via Twilio:

```env
# Twilio Configuration (for SMS)
TWILIO_ACCOUNT_SID=your_twilio_account_sid
TWILIO_AUTH_TOKEN=your_twilio_auth_token
TWILIO_PHONE_NUMBER=+1234567890
```

### Setting up Twilio:
1. Sign up at https://www.twilio.com
2. Get your Account SID and Auth Token from the dashboard
3. Purchase a phone number or use the trial number
4. Add credentials to your `.env` file

## Notes

- **Emails** are sent via `services/email.go` (new, better templates)
- **SMS** messages are sent via `services/notifications.go` (Twilio integration)
- Both are sent asynchronously (in goroutines) so they don't block reservation creation
- If email/SMS service is not configured, reservations will still work normally
- Email/SMS failures will not cause reservation creation to fail
- All notification status is logged in the server console

## Migration Note

The old email functionality in `notifications.go` has been deprecated. The system now uses:
- `services/email.go` for emails (better HTML templates, images, detailed info)
- `services/notifications.go` for SMS only

No action needed - the migration is automatic!

