package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// MaxTelegramMessageLength is the maximum size of a Telegram message (4096 characters).
const MaxTelegramMessageLength = 4096

func sendTelegramNotification(chatID string, token string, message string, txtFilePath string) (int, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return 0, fmt.Errorf("failed to create Telegram bot: %v", err)
	}

	chatIDInt, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid chat ID format: %v", err)
	}

	file, err := os.Open(txtFilePath)
	if err != nil {
		return 0, fmt.Errorf("error opening TXT file: %v", err)
	}
	defer file.Close()

	doc := tgbotapi.NewDocument(chatIDInt, tgbotapi.FileReader{
		Name:   txtFilePath,
		Reader: file,
	})
	doc.Caption = message

	msg, err := bot.Send(doc)
	if err != nil {
		return 0, fmt.Errorf("error sending TXT file: %v", err)
	}

	fmt.Println("Message with TXT file sent successfully")
	return msg.MessageID, nil
}

func sendMessageWithtxt(bot *tgbotapi.BotAPI, chatID int64, message string, txtFilePath string) {
	// Open the TXT file
	file, err := os.Open(txtFilePath)
	if err != nil {
		log.Println("Error opening TXT file:", err)
		return
	}
	defer file.Close()

	// Create a new document message with the TXT file
	doc := tgbotapi.NewDocument(chatID, tgbotapi.FileReader{
		Name:   txtFilePath,
		Reader: file,
	})

	// Add the message as the caption for the TXT file
	doc.Caption = message // The message will appear as the caption to the file

	// Send the document (TXT file) with the message caption
	_, err = bot.Send(doc)
	if err != nil {
		log.Println("Error sending TXT file:", err)
		return
	}

	fmt.Println("Message with TXT file sent successfully")
}

func updateMessageFile(chatID string, token string, originalMessageID int, txtFilePath string, message_body string) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return fmt.Errorf("failed to create Telegram bot: %v", err)
	}

	chatIDInt, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid chat ID format: %v", err)
	}

	file, err := os.Open(txtFilePath)
	if err != nil {
		return fmt.Errorf("error opening TXT file: %v", err)
	}
	defer file.Close()

	reply := tgbotapi.NewDocument(chatIDInt, tgbotapi.FileReader{
		Name:   txtFilePath,
		Reader: file,
	})
	reply.ReplyToMessageID = originalMessageID
	reply.Caption = "Updated file attached."
	reply.Caption = message_body

	_, err = bot.Send(reply)
	if err != nil {
		return fmt.Errorf("error sending updated file: %v", err)
	}

	fmt.Println("Updated file sent successfully")
	return nil
}
func editMessageFile(chatID string, token string, messageID int, txtFilePath string, msg_body string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/editMessageMedia", token)

	// Open the TXT file
	file, err := os.Open(txtFilePath)
	if err != nil {
		return fmt.Errorf("error opening TXT file: %v", err)
	}
	defer file.Close()

	// Create a multipart form request
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add chat_id, message_id, and media fields
	_ = writer.WriteField("chat_id", chatID)
	_ = writer.WriteField("message_id", fmt.Sprintf("%d", messageID))
	media := map[string]interface{}{
		"type":    "document",
		"media":   "attach://file",
		"caption": "Note - Message has been updated .\n \n" + msg_body,
	}
	mediaJSON, _ := json.Marshal(media)
	_ = writer.WriteField("media", string(mediaJSON))

	// Add the file as a form field
	filePart, err := writer.CreateFormFile("file", txtFilePath)
	if err != nil {
		return fmt.Errorf("error creating form file: %v", err)
	}
	_, err = io.Copy(filePart, file)
	if err != nil {
		return fmt.Errorf("error copying file to form: %v", err)
	}

	// Close the writer
	writer.Close()

	// Send the request
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Check for success response
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to edit message: %s", string(body))
	}

	fmt.Println("Message edited successfully with updated file.")
	return nil
}
