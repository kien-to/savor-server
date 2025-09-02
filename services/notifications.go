package services

import (
	"bytes"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"time"
)

// NotificationService handles email and SMS notifications
type NotificationService struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	TwilioSID    string
	TwilioToken  string
	TwilioPhone  string
}

// ReservationNotificationData contains data for notification templates
type ReservationNotificationData struct {
	CustomerName  string
	StoreName     string
	StoreAddress  string
	Quantity      int
	TotalAmount   float64
	PickupTime    string
	ReservationID string
	Email         string
	Phone         string
}

// Global notification service instance
var NotificationSvc *NotificationService

// InitializeNotificationService initializes the notification service
func InitializeNotificationService() {
	NotificationSvc = &NotificationService{
		SMTPHost:     getEnvOrDefault("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnvOrDefault("SMTP_PORT", "587"),
		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		TwilioSID:    os.Getenv("TWILIO_ACCOUNT_SID"),
		TwilioToken:  os.Getenv("TWILIO_AUTH_TOKEN"),
		TwilioPhone:  os.Getenv("TWILIO_PHONE_NUMBER"),
	}
}

// SendReservationConfirmation sends both email and SMS notifications
func (ns *NotificationService) SendReservationConfirmation(data ReservationNotificationData) error {
	var errors []string

	// Send email if email is provided
	if data.Email != "" && ns.SMTPUsername != "" && ns.SMTPPassword != "" {
		if err := ns.sendEmailConfirmation(data); err != nil {
			errors = append(errors, fmt.Sprintf("Email error: %v", err))
		}
	}

	// Send SMS if phone is provided
	if data.Phone != "" && ns.TwilioSID != "" && ns.TwilioToken != "" {
		if err := ns.sendSMSConfirmation(data); err != nil {
			errors = append(errors, fmt.Sprintf("SMS error: %v", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("notification errors: %v", errors)
	}

	return nil
}

// sendEmailConfirmation sends email confirmation
func (ns *NotificationService) sendEmailConfirmation(data ReservationNotificationData) error {
	subject := "Xác nhận đặt hàng - Savor"
	body := ns.generateEmailTemplate(data)

	// Setup authentication
	auth := smtp.PlainAuth("", ns.SMTPUsername, ns.SMTPPassword, ns.SMTPHost)

	// Create message
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", data.Email, subject, body))

	// Send email
	err := smtp.SendMail(ns.SMTPHost+":"+ns.SMTPPort, auth, ns.SMTPUsername, []string{data.Email}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	fmt.Printf("Email confirmation sent to: %s\n", data.Email)
	return nil
}

// sendSMSConfirmation sends SMS confirmation via Twilio
func (ns *NotificationService) sendSMSConfirmation(data ReservationNotificationData) error {
	message := ns.generateSMSTemplate(data)

	// Twilio API endpoint
	url := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", ns.TwilioSID)

	// Prepare form data
	payload := fmt.Sprintf("From=%s&To=%s&Body=%s", ns.TwilioPhone, data.Phone, message)

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(payload))
	if err != nil {
		return fmt.Errorf("failed to create SMS request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(ns.TwilioSID, ns.TwilioToken)

	// Send request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return fmt.Errorf("SMS API returned status: %d", resp.StatusCode)
	}

	fmt.Printf("SMS confirmation sent to: %s\n", data.Phone)
	return nil
}

// generateEmailTemplate creates HTML email template
func (ns *NotificationService) generateEmailTemplate(data ReservationNotificationData) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Xác nhận đặt hàng</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #ef4444; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background-color: #f9f9f9; }
        .order-details { background-color: white; padding: 15px; margin: 15px 0; border-radius: 5px; }
        .footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>SAVOR</h1>
            <h2>Xác nhận đặt hàng</h2>
        </div>
        <div class="content">
            <p>Xin chào <strong>%s</strong>,</p>
            <p>Cảm ơn bạn đã đặt hàng tại Savor! Đơn hàng của bạn đã được xác nhận.</p>
            
            <div class="order-details">
                <h3>Chi tiết đơn hàng:</h3>
                <p><strong>Mã đơn hàng:</strong> %s</p>
                <p><strong>Cửa hàng:</strong> %s</p>
                <p><strong>Địa chỉ:</strong> %s</p>
                <p><strong>Số lượng:</strong> %d túi</p>
                <p><strong>Tổng tiền:</strong> %.2fk đ</p>
                <p><strong>Thời gian nhận hàng:</strong> %s</p>
            </div>
            
            <p>Vui lòng đến cửa hàng đúng giờ để nhận hàng. Cảm ơn bạn đã sử dụng dịch vụ của chúng tôi!</p>
        </div>
        <div class="footer">
            <p>© 2024 Savor. Tất cả quyền được bảo lưu.</p>
        </div>
    </div>
</body>
</html>`,
		data.CustomerName,
		data.ReservationID,
		data.StoreName,
		data.StoreAddress,
		data.Quantity,
		data.TotalAmount,
		data.PickupTime,
	)
}

// generateSMSTemplate creates SMS message
func (ns *NotificationService) generateSMSTemplate(data ReservationNotificationData) string {
	return fmt.Sprintf("SAVOR - Xác nhận đặt hàng\n\nXin chào %s!\n\nĐơn hàng #%s đã được xác nhận:\n- Cửa hàng: %s\n- Số lượng: %d túi\n- Tổng tiền: %.2fk đ\n- Nhận hàng: %s\n\nCảm ơn bạn!",
		data.CustomerName,
		data.ReservationID,
		data.StoreName,
		data.Quantity,
		data.TotalAmount,
		data.PickupTime,
	)
}

// getEnvOrDefault returns environment variable or default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
