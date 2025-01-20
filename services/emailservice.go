package services

import (
	"fmt"
	"net/smtp"
)

// Send OTP on Email
func SendOTPEmail(email, otp string) error {

	//Set UP Email Sender
	sender := "firman.hariri@idstar.co.id"
	password := "Gintama26"
	host := "smtp.gmail.com"
	port := "587"

	//Set Up Receivent Email
	to := []string{email}

	subject := "Email Verification OTP"
	body := fmt.Sprint("Your OTP for email verification is: %s\"", otp)

	// Set up the message
	message := []byte("Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n\r\n" +
		body)

	// Authentication
	auth := smtp.PlainAuth("", sender, password, host)

	// Send email
	err := smtp.SendMail(host+":"+port, auth, sender, to, message)
	if err != nil {
		return err
	}

	return nil

}
