# Email & SMS Configuration Setup

The Savor backend supports sending email and SMS confirmations for reservations.

## Architecture

- **Notification Service** (`services/notifications.go`): Handles SMS confirmations via Twilio
- **Note**: The old email functionality in `notifications.go` is deprecated - all emails use the new service

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

- **SMS** messages are sent via `services/notifications.go` (Twilio integration)
- Both are sent asynchronously (in goroutines) so they don't block reservation creation
- All notification status is logged in the server console

