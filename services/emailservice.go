package services

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendOTPEmail(email, otp string) error {

	username := getEnv("SMTP_USERNAME", "")
	password := getEnv("SMTP_PASSWORD", "")
	host := getEnv("SMTP_HOST", "smtp-relay.brevo.com")
	port := getEnv("SMTP_PORT", "587")

	from := getEnv("SMTP_FROM_EMAIL", "")
	to := []string{email}

	subject := getEnv("EMAIL_SUBJECT", "Your OTP Code")
	expiryMinutes := getEnv("OTP_EXPIRY_MINUTES", "5")
	body := fmt.Sprintf("Your OTP verification code is: %s\nThis OTP expires in %s minutes.", otp, expiryMinutes)

	message := []byte(
		"From: " + from + "\r\n" +
			"To: " + email + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=UTF-8\r\n\r\n" +
			body)

	auth := smtp.PlainAuth("", username, password, host)

	err := smtp.SendMail(host+":"+port, auth, from, to, message)
	if err != nil {
		return err
	}

	return nil
}

// getEnv retrieves the value of the environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
