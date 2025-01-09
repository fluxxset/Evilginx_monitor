package main

import (
	"fmt"
	"log"
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
