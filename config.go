package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	TelegramChatID string `json:"telegr_chatid"`
	TelegramToken  string `json:"telegram_token"`
	TelegramEnable bool   `json:"telegram_enable"`

	MailHost      string `json:"mail_host"`
	MailPort      int    `json:"mail_port"`
	MailUser      string `json:"mail_user"`
	MailPassword  string `json:"mail_password"`
	ToMail        string `json:"to_mail"`
	MailEnable    bool   `json:"mail_enable"`
	DiscordChatID string `json:"discord_chat_id"`
	DiscordToken  string `json:"discord_token"`
	DiscordEnable bool   `json:"discord_enable"`
	DBFilePath    string `json:"dbfile_path"`
}

var configFilePath = filepath.Join(os.Getenv("HOME"), ".evilginx_monitor", "config.json")

func loadConfig() (*Config, error) {
	file, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func saveConfig(config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(configFilePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func initConfig() {

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		config := &Config{}
		saveConfig(config)
		fmt.Println("Config initialized.")
	}
}

func showConfig() error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println("Current configuration:")
	fmt.Println(string(configData))

	return nil
}
