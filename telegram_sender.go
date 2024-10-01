package main

import (
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func sendTelegramNotification(chatID string, token string, message string) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Println("Failed to create Telegram bot:", err)
		return
	}

	// Convert chatID string to int64 if necessary
	chatIDInt, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		log.Println("Invalid chat ID format:", err)
		return
	}

	msg := tgbotapi.NewMessage(chatIDInt, message) // Use NewMessage for direct messages to users or groups
	_, err = bot.Send(msg)
	if err != nil {
		log.Println("Failed to send message:", err)
	} else {
		fmt.Println("Message sent using Telegram")
	}
}
