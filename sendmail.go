package main

import (
	"log"
	"net/smtp"
)

func main() {
	smtpHost := "localhost" // not "mailhog"
	smtpPort := "1025"

	auth := smtp.PlainAuth("", "", "", smtpHost)

	to := []string{"test@example.com"}
	msg := []byte("To: test@example.com\r\n" +
		"Subject: Hello from Go\r\n" +
		"\r\n" +
		"This is a test email sent from Go.\r\n")

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, "sender@example.com", to, msg)
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	log.Println("Email sent successfully!")
}
