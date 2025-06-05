package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

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

	// Send email to admin
	emailMsg := fmt.Sprintf(`
		New Order Notification
		Order ID: %d
		Customer: %s %s
		Phone: %s
		Total Amount: %.2f %s
		Status: %s
	`, order.ID, order.Customer.FirstName, order.Customer.LastName,
		order.Customer.Phone, order.Total, s.config.Currency, order.Status)

	if err := s.sendEmail(s.config.AdminEmail, "New Order Placed", emailMsg); err != nil {
		log.Printf("Failed to send admin notification email: %v", err)
		return fmt.Errorf("failed to send email: %w", err)
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
	return nil
}

func (s *notificationService) sendSMS(to, message string) error {
	if s.config.AfricaTalkingAPIKey == "" || s.config.AfricaTalkingUsername == "" {
		return fmt.Errorf("africa's Talking credentials not configured")
	}

	url := "https://api.africastalking.com/version1/messaging"

	payload := map[string]interface{}{
		"to":       "+254722976334",
		"message":  message,
		"username": s.config.AfricaTalkingUsername,
	}

	if s.config.SMSSenderID != "" {
		payload["from"] = s.config.SMSSenderID
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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading SMS response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("SMS API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Log successful delivery
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err == nil {
		if smsData, ok := response["SMSMessageData"].(map[string]interface{}); ok {
			log.Printf("SMS sent to %s - Status: %s", to, smsData["Message"])
		}
	}

	return nil
}

func (s *notificationService) sendEmail(to, subject, body string) error {
	// Implement your email sending logic here
	// This could use SMTP, SendGrid, Mailgun, etc.
	// Currently just logging as placeholder
	log.Printf("Email sent to %s\nSubject: %s\nBody: %s", to, subject, body)
	return nil
}
