package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

// sendDiscordNotification sends a message and an attachment to a user on Discord.
func sendDiscordNotification(userID string, token string, message string, attachmentPath string) {
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

	// Prepare the attachment
	// If attachmentPath is not empty, attach the file
	var files []*discordgo.File
	if attachmentPath != "" {
		file, err := os.Open(attachmentPath)
		if err != nil {
			log.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		// Add the file to the files slice
		files = append(files, &discordgo.File{
			Name:   "attachment", // You can specify the file name here
			Reader: file,
		})
	}

	// Send the message with the attachment to the DM channel
	_, err = dg.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: message,
		Files:   files, // Attach files here
	})
	if err != nil {
		log.Println("Failed to send message:", err)
	} else {
		fmt.Println("Message and attachment sent successfully!")
	}
}
