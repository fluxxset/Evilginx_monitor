package main

import (
	"fmt"
	"net/smtp"
	"strconv"
)

func sendMailNotification(smtpHost string, smtpPort int, smtpUser, smtpPassword, to, body string) error {
	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)
	smtpAddress := smtpHost + ":" + strconv.Itoa(smtpPort)
	from := fmt.Sprintf("%s", smtpUser)
	subject := "New Session Captured."
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	err := smtp.SendMail(smtpAddress, auth, from, []string{to}, []byte(msg))
	if err != nil {
		fmt.Println("Failed to send email:", err)
		return err
	}
	fmt.Println("Email sent successfully.")
	return nil
}
