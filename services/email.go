package services

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"time"
)

type EmailService struct {
	smtpHost     string
	smtpPort     string
	smtpUser     string
	smtpPassword string
	fromEmail    string
	fromName     string
}

type ReservationEmailData struct {
	CustomerName    string
	StoreName       string
	StoreAddress    string
	StoreImage      string
	Quantity        int
	TotalAmount     float64
	PickupTime      string
	ReservationID   string
	Status          string
	PaymentType     string
	CreatedAt       time.Time
	OriginalPrice   float64
	DiscountedPrice float64
}

var emailService *EmailService

// InitializeEmailService initializes the email service with environment variables
func InitializeEmailService() {
	emailService = &EmailService{
		smtpHost:     os.Getenv("SMTP_HOST"),
		smtpPort:     os.Getenv("SMTP_PORT"),
		smtpUser:     os.Getenv("SMTP_USER"),
		smtpPassword: os.Getenv("SMTP_PASSWORD"),
		fromEmail:    os.Getenv("FROM_EMAIL"),
		fromName:     os.Getenv("FROM_NAME"),
	}

	// Check if email service is properly configured
	if emailService.smtpHost == "" || emailService.smtpPort == "" {
		log.Println("Warning: Email service not configured. Set SMTP_HOST and SMTP_PORT environment variables.")
	} else {
		log.Println("Email service initialized successfully")
	}
}

// maskPassword masks the password for logging
func maskPassword(password string) string {
	if len(password) == 0 {
		return "(empty)"
	}
	if len(password) <= 4 {
		return "****"
	}
	return password[:4] + "****"
}

// GetEmailService returns the singleton email service instance
func GetEmailService() *EmailService {
	if emailService == nil {
		InitializeEmailService()
	}
	return emailService
}

// IsConfigured checks if the email service is properly configured
func (e *EmailService) IsConfigured() bool {
	return e.smtpHost != "" && e.smtpPort != "" && e.smtpUser != "" && e.smtpPassword != ""
}

// SendReservationConfirmation sends a confirmation email for a new reservation
func (e *EmailService) SendReservationConfirmation(toEmail string, data ReservationEmailData) error {
	if !e.IsConfigured() {
		log.Println("Email service not configured, skipping email send")
		return nil // Don't fail reservation creation if email isn't configured
	}

	subject := fmt.Sprintf("Xác nhận đặt chỗ tại %s - Savor", data.StoreName)
	body, err := e.generateReservationEmail(data)
	if err != nil {
		log.Printf("Failed to generate email template: %v", err)
		return err
	}

	return e.sendEmail(toEmail, subject, body)
}

// sendEmail sends an email using SMTP
func (e *EmailService) sendEmail(to, subject, body string) error {
	// Set up authentication
	auth := smtp.PlainAuth("", e.smtpUser, e.smtpPassword, e.smtpHost)

	// Build email message with proper headers to reduce spam likelihood
	from := fmt.Sprintf("%s <%s>", e.fromName, e.fromEmail)
	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"Reply-To: %s\r\n"+
		"X-Mailer: Savor App\r\n"+
		"MIME-version: 1.0;\r\n"+
		"Content-Type: text/html; charset=\"UTF-8\";\r\n"+
		"X-Priority: 3\r\n"+
		"\r\n"+
		"%s\r\n", from, to, subject, e.fromEmail, body))

	// Send email
	addr := fmt.Sprintf("%s:%s", e.smtpHost, e.smtpPort)
	err := smtp.SendMail(addr, auth, e.fromEmail, []string{to}, msg)
	if err != nil {
		log.Printf("Failed to send email to %s: %v", to, err)
		return err
	}

	log.Printf("Sent reservation confirmation email to %s", to)
	return nil
}

// generateReservationEmail generates the HTML email template
func (e *EmailService) generateReservationEmail(data ReservationEmailData) (string, error) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            background-color: #ffffff;
            border-radius: 10px;
            padding: 30px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .header {
            background-color: #036B52;
            color: #fff6e7;
            padding: 20px;
            border-radius: 10px 10px 0 0;
            text-align: center;
            margin: -30px -30px 30px -30px;
        }
        h1 {
            margin: 0;
            font-size: 24px;
        }
        .store-info {
            margin: 20px 0;
            padding: 20px;
            background-color: #f8f9fa;
            border-radius: 8px;
        }
        .store-image {
            width: 100%;
            max-height: 200px;
            object-fit: cover;
            border-radius: 8px;
            margin-bottom: 15px;
        }
        .info-row {
            display: flex;
            justify-content: space-between;
            padding: 10px 0;
            border-bottom: 1px solid #e0e0e0;
        }
        .info-row:last-child {
            border-bottom: none;
        }
        .label {
            font-weight: 600;
            color: #036B52;
        }
        .value {
            color: #333;
        }
        .price-section {
            background-color: #e8f5e9;
            padding: 15px;
            border-radius: 8px;
            margin: 20px 0;
        }
        .price {
            font-size: 24px;
            font-weight: bold;
            color: #036B52;
        }
        .original-price {
            text-decoration: line-through;
            color: #999;
            font-size: 16px;
            margin-right: 10px;
        }
        .footer {
            margin-top: 30px;
            padding-top: 20px;
            border-top: 2px solid #e0e0e0;
            text-align: center;
            color: #666;
            font-size: 14px;
        }
        .button {
            display: inline-block;
            padding: 12px 30px;
            background-color: #036B52;
            color: #fff6e7;
            text-decoration: none;
            border-radius: 6px;
            margin: 20px 0;
            font-weight: 600;
        }
        .status-badge {
            display: inline-block;
            padding: 6px 12px;
            background-color: #4CAF50;
            color: white;
            border-radius: 4px;
            font-size: 14px;
            font-weight: 600;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎉 Đặt chỗ thành công!</h1>
        </div>

        <p>Xin chào <strong>{{.CustomerName}}</strong>,</p>
        <p>Cảm ơn bạn đã đặt túi bất ngờ tại Savor! Dưới đây là chi tiết đơn hàng của bạn:</p>

        <div class="store-info">
            {{if .StoreImage}}
            <img src="{{.StoreImage}}" alt="{{.StoreName}}" class="store-image">
            {{end}}
            <h2 style="margin-top: 0; color: #036B52;">{{.StoreName}}</h2>
            
            <div class="info-row">
                <span class="label">📍 Địa chỉ:</span>
                <span class="value">{{.StoreAddress}}</span>
            </div>
            
            <div class="info-row">
                <span class="label">🕐 Thời gian lấy hàng:</span>
                <span class="value">{{.PickupTime}}</span>
            </div>
            
            <div class="info-row">
                <span class="label">📦 Số lượng:</span>
                <span class="value">{{.Quantity}} túi</span>
            </div>
            
            <div class="info-row">
                <span class="label">💳 Hình thức thanh toán:</span>
                <span class="value">{{.PaymentType}}</span>
            </div>
            
            <div class="info-row">
                <span class="label">🔢 Mã đặt chỗ:</span>
                <span class="value">{{.ReservationID}}</span>
            </div>
            
            <div class="info-row">
                <span class="label">📅 Trạng thái:</span>
                <span class="value"><span class="status-badge">{{.Status}}</span></span>
            </div>
        </div>

        <div class="price-section">
            <div class="info-row">
                <span class="label">Tổng tiền:</span>
                <div>
                    {{if gt .OriginalPrice .DiscountedPrice}}
                    <span class="original-price">{{printf "%.0f" .OriginalPrice}}.000đ</span>
                    {{end}}
                    <span class="price">{{printf "%.0f" .TotalAmount}}.000đ</span>
                </div>
            </div>
            {{if gt .OriginalPrice .DiscountedPrice}}
            <p style="margin: 10px 0 0 0; color: #4CAF50; font-weight: 600;">
                🎊 Bạn tiết kiệm được {{printf "%.0f" (sub .OriginalPrice .TotalAmount)}}.000đ!
            </p>
            {{end}}
        </div>

        <p><strong>Lưu ý quan trọng:</strong></p>
        <ul>
            <li>Vui lòng đến đúng giờ để lấy hàng</li>
            <li>Mang theo mã đặt chỗ khi đến lấy hàng</li>
            <li>Liên hệ cửa hàng nếu có bất kỳ thắc mắc nào</li>
        </ul>

        <div style="text-align: center;">
            <p>Xem chi tiết đơn hàng trong ứng dụng Savor</p>
        </div>

        <div class="footer">
            <p><strong>Savor</strong> - Giảm lãng phí thực phẩm, tiết kiệm chi phí</p>
            <p>Email này được gửi tự động, vui lòng không trả lời.</p>
            <p style="font-size: 12px; color: #999;">
                Nếu bạn không thực hiện đặt chỗ này, vui lòng bỏ qua email này.
            </p>
        </div>
    </div>
</body>
</html>
`

	// Parse and execute template
	t, err := template.New("reservation").Funcs(template.FuncMap{
		"sub": func(a, b float64) float64 { return a - b },
	}).Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
