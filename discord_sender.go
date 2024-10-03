package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func sendDiscordNotification(userID string, token string, message string) {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Println("Failed to create Discord session:", err)
		return
	}

	// Open a websocket connection to Discord
	err = dg.Open()
	if err != nil {
		log.Println("Error opening connection:", err)
		return
	}
	defer dg.Close()

	// Create a DM (Direct Message) channel with the user
	channel, err := dg.UserChannelCreate(userID)
	if err != nil {
		log.Println("Error creating DM channel:", err)
		return
	}

	// Send the message to the DM channel
	_, err = dg.ChannelMessageSend(channel.ID, message)
	if err != nil {
		log.Println("Failed to send message:", err)
	} else {
		fmt.Println("Message sent using Discord")
	}
}
