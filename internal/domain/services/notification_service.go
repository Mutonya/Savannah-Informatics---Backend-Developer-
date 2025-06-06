package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Mutonya/Savanah/internal/utils/templates"
	"io"
	"log"
	"net/http"
	"net/smtp"

	"github.com/Mutonya/Savanah/internal/config"
	"github.com/Mutonya/Savanah/internal/domain/models"
)

type NotificationService interface {
	SendOrderConfirmation(order *models.Order) error
	SendStatusUpdate(order *models.Order) error
}

type notificationService struct {
	config *config.Config
}

func NewNotificationService(config *config.Config) NotificationService {
	return &notificationService{config: config}
}

func (s *notificationService) SendOrderConfirmation(order *models.Order) error {
	// Send SMS to customer
	smsMsg := fmt.Sprintf("Hello %s, your order #%d has been received. Total: %.2f %s",
		order.Customer.FirstName, order.ID, order.Total, s.config.Currency)
	if err := s.sendSMS(order.Customer.Phone, smsMsg); err != nil {
		log.Printf("Failed to send order confirmation SMS: %v", err)
		return fmt.Errorf("failed to send SMS: %w", err)
	}

	// Prepare email data
	emailData := struct {
		Order  *models.Order
		Config *config.Config
	}{
		Order:  order,
		Config: s.config,
	}

	// Send email to customer
	customerSubject := fmt.Sprintf("Order #%d Confirmation", order.ID)

	if err := s.sendHTMLEmail(
		order.Customer.Email,
		customerSubject,
		"order_confirmation",
		emailData,
	); err != nil {
		log.Printf("Failed to send customer confirmation email: %v", err)
		return fmt.Errorf("failed to send customer email: %w", err)
	}

	// Send email to admin
	adminSubject := fmt.Sprintf("New Order #%d Received", order.ID)
	if err := s.sendHTMLEmail(
		s.config.AdminEmail,
		adminSubject,
		"admin_notification",
		emailData,
	); err != nil {
		log.Printf("Failed to send admin notification email: %v", err)
		return fmt.Errorf("failed to send admin email: %w", err)
	}

	return nil
}

func (s *notificationService) SendStatusUpdate(order *models.Order) error {
	// Send SMS to customer about status change
	smsMsg := fmt.Sprintf("Hello %s, your order #%d status is now: %s",
		order.Customer.FirstName, order.ID, order.Status)
	if err := s.sendSMS(order.Customer.Phone, smsMsg); err != nil {
		return fmt.Errorf("failed to send status update SMS: %w", err)
	}

	// Send status update email
	if err := s.sendHTMLEmail(
		order.Customer.Email,
		fmt.Sprintf("Order #%d Status Update", order.ID),
		"status_update",
		struct {
			Order *models.Order
		}{order},
	); err != nil {
		return fmt.Errorf("failed to send status email: %w", err)
	}

	return nil
}

func (s *notificationService) sendSMS(to, message string) error {
	if s.config.AfricaTalkingAPIKey == "" || s.config.AfricaTalkingUsername == "" {
		return fmt.Errorf("Africa's Talking credentials not configured")
	}

	url := "https://api.africastalking.com/version1/messaging/bulk"

	// Hardcoded phone number as a slice
	phoneNumbers := []string{"+254722976334"}

	payload := map[string]interface{}{
		"username":     s.config.AfricaTalkingUsername,
		"message":      message,
		"phoneNumbers": phoneNumbers,
	}

	if s.config.SMSSenderID != "" {
		payload["senderId"] = s.config.SMSSenderID
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling SMS payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("error creating SMS request: %w", err)
	}

	req.Header.Set("apiKey", s.config.AfricaTalkingAPIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending SMS request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading SMS response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("SMS API error (status %d): %s", resp.StatusCode, string(body))
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err == nil {
		if smsData, ok := response["SMSMessageData"].(map[string]interface{}); ok {
			log.Printf("SMS sent to %v - Status: %s", phoneNumbers, smsData["Message"])
		}
	}

	return nil
}

func (s *notificationService) sendHTMLEmail(to, subject, templateName string, data interface{}) error {

	log.Printf("Sending email to: %s using host %s:%d", to, s.config.SMTPHost, s.config.SMTPPort)

	if s.config.SMTPHost == "" {
		return fmt.Errorf("SMTP configuration not set")
	}

	// Get email template
	template, err := templates.GetEmailTemplate(templateName)
	if err != nil {
		return fmt.Errorf("error getting email template: %w", err)
	}

	// Parse and execute template
	body, err := templates.ParseTemplate(template.Body, data)
	if err != nil {
		return fmt.Errorf("error parsing email template: %w", err)
	}
	log.Println("Email body generated successfully")

	// Use template.Subject if no subject was provided
	if subject == "" {
		subject = template.Subject
	}

	// Conditionally set up authentication (skip for MailHog)
	var auth smtp.Auth
	if s.config.SMTPUsername != "" && s.config.SMTPPassword != "" {
		auth = smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPHost)
	}

	// Set up message headers
	headers := map[string]string{
		"From":         s.config.AdminEmail,
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=\"UTF-8\"",
	}

	// Build message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Send email
	err = smtp.SendMail(
		fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort),
		auth,
		s.config.AdminEmail,
		[]string{to},
		[]byte(message),
	)

	if err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	log.Printf("Email sent to %s with subject: %s", to, subject)
	return nil
}
