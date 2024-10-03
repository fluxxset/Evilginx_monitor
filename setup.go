package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/user"
	"path/filepath"
	"time"

	_ "github.com/glebarez/sqlite"
)

const (
	configDir  = ".evilginx_monitor"
	configFile = "config.json"
	keysFile   = "keys.json"
	credsFile  = "creds.json"
	dbFile     = "record_tracker.db"
)

type ChatConfig struct {
	ChatID  string `json:"chat_id"`
	Token   string `json:"token"`
	Enabled bool   `json:"enabled"`
}

type MailConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	ToMail   string `json:"To"`
	Enabled  bool   `json:"enabled"`
}

type DefaultConfig struct {
	TelegramChatID  string `json:"telegr_chatid"`
	TelegramToken   string `json:"telegram_token"`
	TelegramEnabled bool   `json:"telegram_enable"`
	MailHost        string `json:"mail_host"`
	MailPort        int    `json:"mail_port"`
	MailUser        string `json:"mail_user"`
	MailPassword    string `json:"mail_password"`
	ToMail          string `json:"ToMail"`
	MailEnabled     bool   `json:"mail_enable"`
	DiscordChatID   string `json:"discord_chat_id"`
	DiscordToken    string `json:"discord_token"`
	DiscordEnabled  bool   `json:"discord_enable"`
	DBFilePath      string `json:"dbfile_path"`
}

func Setup() error {

	userHomeDir, err := user.Current()
	if err != nil {
		return fmt.Errorf("error getting user home directory: %v", err)
	}

	configDirPath := filepath.Join(userHomeDir.HomeDir, configDir)
	configFilePath := filepath.Join(configDirPath, configFile)
	keysFilePath := filepath.Join(configDirPath, keysFile)
	credsFilePath := filepath.Join(configDirPath, credsFile)

	if _, err := os.Stat(configDirPath); os.IsNotExist(err) {
		// Create the directory
		err := os.Mkdir(configDirPath, 0755) // You can change the permissions as needed
		if err != nil {
			fmt.Println("Error creating directory:", err)
		}
		fmt.Println("Directory created:", configDirPath)
	} else {
		fmt.Println("Directory already exists:", configDirPath)
	}

	for _, filePath := range []string{configFilePath, keysFilePath, credsFilePath} {
		if err := createFileIfNotExists(filePath); err != nil {
			return fmt.Errorf("error creating file: %s, %v", filePath, err)
		}
	}

	if err := setDefaultConfig(configFilePath); err != nil {
		return fmt.Errorf("error setting default config: %v", err)
	}

	if err := generateKeys(keysFilePath); err != nil {
		return fmt.Errorf("error generating keys: %v", err)
	}

	fmt.Println("Setup completed successfully.")
	return nil
}

func setDefaultConfig(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) || isFileEmpty(filePath) {
		defaultConfig := DefaultConfig{
			TelegramChatID:  "",
			TelegramToken:   "",
			TelegramEnabled: false,
			MailHost:        "",
			MailPort:        0,
			MailUser:        "",
			MailPassword:    "",
			ToMail:          "",
			MailEnabled:     false,
			DiscordChatID:   "",
			DiscordToken:    "",
			DiscordEnabled:  false,
			DBFilePath:      "/root/.evilginx/data.db",
		}
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		if err := encoder.Encode(defaultConfig); err != nil {
			return err
		}

		fmt.Printf("Default configuration created at %s\n", filePath)
	}

	return nil
}

func isFileEmpty(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return true
	}
	return info.Size() == 0
}

func generateKeys(filePath string) error {
	key := generateRandomKey()
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(key)
	if err != nil {
		return err
	}
	fmt.Printf("Generated key saved at %s\n", filePath)
	return nil
}

func generateRandomKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	key := make([]byte, 32) // Example: 32 characters long
	for i := range key {
		key[i] = charset[rand.Intn(len(charset))]
	}
	return string(key)
}

func createFileIfNotExists(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("Creating file: %s\n", filePath)
		_, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("could not create file %s: %v", filePath, err)
		}
	}
	// fmt.Printf("File already exists: %s\n", filePath)
	return nil
}
