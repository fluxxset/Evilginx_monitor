package main

import (
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// MaxTelegramMessageLength is the maximum size of a Telegram message (4096 characters).
const MaxTelegramMessageLength = 4096

func sendTelegramNotification(chatID string, token string, message string) {
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

	// Split and send the message in chunks
	sendInChunks(bot, chatIDInt, message)
}

func sendInChunks(bot *tgbotapi.BotAPI, chatID int64, message string) {
	// Split the message into chunks if it exceeds MaxTelegramMessageLength
	for len(message) > 0 {
		part := message
		if len(part) > MaxTelegramMessageLength {
			part = message[:MaxTelegramMessageLength]    // Get the first 4096 characters
			message = message[MaxTelegramMessageLength:] // Remaining part of the message
		} else {
			message = "" // We're done with the whole message
		}

		msg := tgbotapi.NewMessage(chatID, part)
		msg.ParseMode = "Markdown" // Enable Markdown formatting if needed

		_, err := bot.Send(msg)
		if err != nil {
			log.Println("Failed to send message part:", err)
			return
		}
		fmt.Println("Message part sent")
	}
}
