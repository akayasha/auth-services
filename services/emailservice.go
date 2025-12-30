package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type brevoEmailRequest struct {
	Sender struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"sender"`
	To []struct {
		Email string `json:"email"`
	} `json:"to"`
	Subject     string `json:"subject"`
	TextContent string `json:"textContent"`
}

func SendOTPEmail(email, otp string) error {
	apiKey := getEnv("BREVO_API_KEY", "")
	if apiKey == "" {
		return fmt.Errorf("BREVO_API_KEY missing")
	}

	senderEmail := getEnv("BREVO_SENDER_EMAIL", "")
	senderName := getEnv("BREVO_SENDER_NAME", "BloomBudy")
	subject := getEnv("EMAIL_SUBJECT", "Your OTP Code")
	expiry := getEnv("OTP_EXPIRY_MINUTES", "5")

	bodyText := fmt.Sprintf(
		"Your OTP verification code is: %s\n\nThis OTP expires in %s minutes.",
		otp,
		expiry,
	)

	reqBody := brevoEmailRequest{}
	reqBody.Sender.Email = senderEmail
	reqBody.Sender.Name = senderName
	reqBody.Subject = subject
	reqBody.TextContent = bodyText
	reqBody.To = append(reqBody.To, struct {
		Email string `json:"email"`
	}{Email: email})

	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest(
		"POST",
		"https://api.brevo.com/v3/smtp/email",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("brevo send failed, status: %s", resp.Status)
	}

	return nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
