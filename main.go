package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/fatih/color"
)

var (
	watcher    *fsnotify.Watcher
	stopChan   chan bool
	monitoring bool
)
func showToolName() {
	// Colors
	white := color.New(color.FgWhite)
	lightBlue := color.New(color.FgHiBlue)
	red := color.New(color.FgRed)
	lightBlack := color.New(color.FgHiBlack)

	// Print the logo with colors
	white.Println(`
                                             ___________      __ __           __               
                                             \_   _____/__  _|__|  |    ____ |__| ____ ___  ___
                                              |    __)_\  \/ /  |  |   / __ \|  |/    \\  \/  /
                                              |        \\   /|  |  |__/ /_/  >  |   |  \>    < 
                                             /_______  / \_/ |__|____/\___  /|__|___|  /__/\_ \
                                                     \/              /_____/         \/      \/
	`)

	// Print additional details with colors
	lightBlue.Printf("                                         Evilginx Monitor Community Edition\n\n")
	fmt.Printf("                                               by %s %s     version %s\n\n",
		lightBlue.Sprint("Fluxxset"),
		red.Sprint("(@fluxxset)"),
		lightBlack.Sprint("1.0.0"))

	// Course link
	fmt.Println("Check out the course here: " + red.Sprint("https://shop.fluxxset.com/product/evilginx-training-course/"))
}

func showHelp() {
	fmt.Println(`
Evilginx Credential Monitor

Usage:
  ./evilginx_monitor [OPTIONS]

Options:
  --help            Show this message and exit.
  --config          Show current configuration.

Interactive Commands:
  start                   Start Monitoring.
  config                  Show current configuration.
  help                    Show this help message.

  tele token <value>      Set Telegram token.
  tele chatid <value>     Set Telegram chat ID.
  tele enable             Enable Telegram notifications.
  tele disable            Disable Telegram notifications.

  mail host <value>       Set SMTP mail host.
  mail port <value>       Set SMTP mail port.
  mail user <value>       Set SMTP mail user.
  mail password <value>   Set SMTP mail password.
  mail to <value>         Set Email to recive alerts.
  mail enable             Enable email notifications.
  mail disable            Disable email notifications.

  discord token <value>   Set Discord token.
  discord chatid <value>  Set Discord chat ID.
  discord enable          Enable Discord notifications.
  discord disable         Disable Discord notifications.

  dbfile path <value>     Set database file path.

  exit                    Exit interactive mode.
`)
}

func createDirIfNotExists(filePath string) error {
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
func reloadConfig(filePath string) {
	config, err := loadConfigx(filePath)
	if err != nil {
		fmt.Printf("Error reloading config: %v\n", err)
		return
	}
	fmt.Println("Config reloaded:", config)
}

func loadConfigx(filePath string) (Config, error) {
	// Simulate a config load for demo purposes
	var config Config

	fmt.Printf("Loading config from %s...\n", filePath)
	time.Sleep(1 * time.Second)

	config1, err := loadConfig()
	if err != nil {
		fmt.Println(err)
	}
	config.DBFilePath = config1.DBFilePath
	return config, nil
}

func StartMonitoring(filePath string) error {
	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("error: file %s does not exist", filePath)
	}

	if monitoring {
		fmt.Println("Already monitoring file.")
		return nil
	}

	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("error creating file watcher: %v", err)
	}

	stopChan = make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					// Reload the configuration when the file is modified
					readFile()
					// reloadConfig(filePath)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Printf("Error watching file: %v\n", err)
			case <-stopChan:
				// Stop the watcher and exit the goroutine
				fmt.Println("Stopping monitoring...")
				return
			}
		}
	}()

	err = watcher.Add(filePath)
	if err != nil {
		return fmt.Errorf("error adding file to watcher: %v", err)
	}

	fmt.Printf("Started monitoring changes to %s...\n", filePath)
	monitoring = true
	return nil
}

// StopMonitoring stops watching the file for changes.
func StopMonitoring() {
	if !monitoring {
		fmt.Println("No file monitoring is currently active.")
		return
	}
	stopChan <- true
	close(stopChan)
	watcher.Close()
	monitoring = false

	fmt.Println("File monitoring stopped.")
}

func UpdateConfig(config *Config) error {
	if err := createDirIfNotExists(filepath.Dir(configFilePath)); err != nil {
		return err
	}
	file, err := os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening config file: %v", err)
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(config); err != nil {
		return fmt.Errorf("error writing to config file: %v", err)
	}
	fmt.Println("Configuration updated successfully.")
	return nil
}

func interactiveMode() {
	reader := bufio.NewReader(os.Stdin)
	showToolName()

	// showHelp()
	fmt.Println("Interactive Mode. Type 'exit' to quit.")

	config := &Config{}

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			fmt.Println("Exiting...")
			os.Exit(0)
		}

		// Handle commands here
		switch {
		case input == "":
			fmt.Println("")
		case strings.HasPrefix(input, "help"):
			showHelp()
		case strings.HasPrefix(input, "start"):
			config, err := loadConfig()
			if err != nil {
				fmt.Println(err)
			}
			var Mfile string
			Mfile = config.DBFilePath

			err = StartMonitoring(Mfile)
			if err != nil {
				log.Println("Error starting file monitoring: %v", err)
			}
		case strings.HasPrefix(input, "stop"):

			StopMonitoring()

		case strings.HasPrefix(input, "config"):
			showConfig()

		case strings.HasPrefix(input, "tele token"):
			token := strings.TrimSpace(input[len("tele token "):])
			config.TelegramToken = token
			if err := UpdateConfig(config); err != nil {
				fmt.Println(err)
			}

		case strings.HasPrefix(input, "tele chatid"):
			chatID := strings.TrimSpace(input[len("tele chatid "):])
			config.TelegramChatID = chatID
			if err := UpdateConfig(config); err != nil {
				fmt.Println(err)
			}

		case input == "tele enable":
			config.TelegramEnable = true
			if err := UpdateConfig(config); err != nil {
				fmt.Println(err)
			}
		case input == "tele disable":
			config.TelegramEnable = false
			if err := UpdateConfig(config); err != nil {
				fmt.Println(err)
			}

		case strings.HasPrefix(input, "mail host"):
			host := strings.TrimSpace(input[len("mail host "):])
			config.MailHost = host
			if err := UpdateConfig(config); err != nil {
				fmt.Println(err)
			}
		case strings.HasPrefix(input, "mail to"):
			mailto := strings.TrimSpace(input[len("mail to "):])
			config.ToMail = mailto
			if err := UpdateConfig(config); err != nil {
				fmt.Println(err)
			}

		case strings.HasPrefix(input, "mail port"):
			portStr := strings.TrimSpace(input[len("mail port "):])
			if port, err := strconv.Atoi(portStr); err == nil {
				config.MailPort = port
				if err := UpdateConfig(config); err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Println("Invalid port number")
			}

		case strings.HasPrefix(input, "mail user"):
			user := strings.TrimSpace(input[len("mail user "):])
			config.MailUser = user
			if err := UpdateConfig(config); err != nil {
				fmt.Println(err)
			}

		case strings.HasPrefix(input, "mail password"):
			password := strings.TrimSpace(input[len("mail password "):])
			config.MailPassword = password
			if err := UpdateConfig(config); err != nil {
				fmt.Println(err)
			}

		case input == "mail enable":
			config.MailEnable = true
			if err := UpdateConfig(config); err != nil {
				fmt.Println(err)
			}
		case input == "mail disable":
			config.MailEnable = false
			if err := UpdateConfig(config); err != nil {
				fmt.Println(err)
			}

		case strings.HasPrefix(input, "discord token"):
			token := strings.TrimSpace(input[len("discord token "):])
			config.DiscordToken = token
			if err := UpdateConfig(config); err != nil {
				fmt.Println(err)
			}

		case strings.HasPrefix(input, "discord chatid"):
			chatID := strings.TrimSpace(input[len("discord chatid "):])
			config.DiscordChatID = chatID
			if err := UpdateConfig(config); err != nil {
				fmt.Println(err)
			}

		case input == "discord enable":
			config.DiscordEnable = true
			if err := UpdateConfig(config); err != nil {
				fmt.Println(err)
			}
		case input == "discord disable":
			config.DiscordEnable = false
			if err := UpdateConfig(config); err != nil {
				fmt.Println(err)
			}

		case strings.HasPrefix(input, "dbfile path"):
			dbFilePath := strings.TrimSpace(input[len("dbfile path "):])
			config.DBFilePath = dbFilePath
			if err := UpdateConfig(config); err != nil {
				fmt.Println(err)
			}

		default:
			fmt.Println("Unknown command. Type 'help' for a list of commands.")
		}
	}
}

func main() {
	if err := Setup(); err != nil {
		fmt.Printf("Setup failed: %v\n", err)
		return
	}
	helpFlag := flag.Bool("help", false, "Show help message")
	configFlag := flag.Bool("config", false, "Show current configuration")
	flag.Parse()

	if *helpFlag {
		showHelp()
		return
	}

	if *configFlag {
		fmt.Println("Current configuration loaded from config.json")
		return
	}

	config, err := loadConfig()
	if err != nil {
		fmt.Println(err)
	}
	var Mfile string
	Mfile = config.DBFilePath

	err = StartMonitoring(Mfile)
	if err != nil {
		log.Println("Error starting file monitoring: %v", err)
	}

	interactiveMode()
}
