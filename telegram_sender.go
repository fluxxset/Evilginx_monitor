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

func sendTelegramNotification(chatID string, token string, message string, zipFilePath string) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Println("Failed to create Telegram bot:", err)
		return
	}

	// Convert chatID string to int64
	chatIDInt, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		log.Println("Invalid chat ID format:", err)
		return
	}

	// Send the message with the zip file as a document (all in one message)
	sendMessageWithtxt(bot, chatIDInt, message, zipFilePath)
}

func sendMessageWithtxt(bot *tgbotapi.BotAPI, chatID int64, message string, zipFilePath string) {
	// Open the zip file
	file, err := os.Open(zipFilePath)
	if err != nil {
		log.Println("Error opening zip file:", err)
		return
	}
	defer file.Close()

	// Create a new document message with the zip file
	doc := tgbotapi.NewDocument(chatID, tgbotapi.FileReader{
		Name:   zipFilePath,
		Reader: file,
	})

	// Add the message as the caption for the zip file
	doc.Caption = message // The message will appear as the caption to the file

	// Send the document (zip file) with the message caption
	_, err = bot.Send(doc)
	if err != nil {
		log.Println("Error sending zip file:", err)
		return
	}

	fmt.Println("Message with zip file sent successfully")
}
