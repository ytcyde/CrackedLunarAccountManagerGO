package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const version = "v1.0"

type MinecraftProfile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Account struct {
	AccessToken           string           `json:"accessToken"`
	AccessTokenExpiresAt  string           `json:"accessTokenExpiresAt"`
	EligibleForMigration  bool             `json:"eligibleForMigration"`
	HasMultipleProfiles   bool             `json:"hasMultipleProfiles"`
	Legacy                bool             `json:"legacy"`
	Persistent            bool             `json:"persistent"`
	UserProperties        []interface{}    `json:"userProperites"`
	LocalID               string           `json:"localId"`
	MinecraftProfile      MinecraftProfile `json:"minecraftProfile"`
	RemoteID              string           `json:"remoteId"`
	Type                  string           `json:"type"`
	Username              string           `json:"username"`
}

type AccountsData struct {
	Accounts map[string]Account `json:"accounts"`
}

var accountsData AccountsData
var lunarAccountsPath string

func init() {
	homeDir, _ := os.UserHomeDir()
	lunarAccountsPath = filepath.Join(homeDir, ".lunarclient", "settings", "game", "accounts.json")
}

func loadJSON() error {
	if _, err := os.Stat(lunarAccountsPath); os.IsNotExist(err) {
		accountsData = AccountsData{Accounts: make(map[string]Account)}
		return nil
	}

	file, err := os.Open(lunarAccountsPath)
	if err != nil {
		return fmt.Errorf("failed to open accounts file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&accountsData)
	if err != nil {
		return fmt.Errorf("failed to parse accounts file: %v", err)
	}

	if accountsData.Accounts == nil {
		accountsData.Accounts = make(map[string]Account)
	}

	return nil
}

func saveJSON() error {
	file, err := os.Create(lunarAccountsPath)
	if err != nil {
		return fmt.Errorf("failed to create accounts file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(accountsData)
	if err != nil {
		return fmt.Errorf("failed to save accounts file: %v", err)
	}

	return nil
}

func createAccount(username, uuid string) {
	newAccount := Account{
		AccessToken:          uuid,
		AccessTokenExpiresAt: "2050-07-02T10:56:30.717167800Z",
		EligibleForMigration: false,
		HasMultipleProfiles:  false,
		Legacy:               true,
		Persistent:           true,
		UserProperties:       []interface{}{},
		LocalID:              uuid,
		MinecraftProfile: MinecraftProfile{
			ID:   uuid,
			Name: username,
		},
		RemoteID: uuid,
		Type:     "Xbox",
		Username: username,
	}

	accountsData.Accounts[uuid] = newAccount
	printLine("SUCCESS", "Your account has successfully been created.", "\033[38;5;39m")
}

func removeAllAccounts() {
	accountsData.Accounts = make(map[string]Account)
	printLine("SUCCESS", "All accounts have been successfully removed.", "\033[38;5;39m")
}

func removeCrackedAccounts() {
	keysToRemove := []string{}
	for uuid, account := range accountsData.Accounts {
		if isValidUUID(account.AccessToken) {
			keysToRemove = append(keysToRemove, uuid)
		}
	}

	for _, key := range keysToRemove {
		delete(accountsData.Accounts, key)
	}

	printLine("SUCCESS", "Cracked accounts have been successfully removed.", "\033[38;5;39m")
}

func removePremiumAccounts() {
	keysToRemove := []string{}
	for uuid, account := range accountsData.Accounts {
		if !isValidUUID(account.AccessToken) {
			keysToRemove = append(keysToRemove, uuid)
		}
	}

	for _, key := range keysToRemove {
		delete(accountsData.Accounts, key)
	}

	printLine("SUCCESS", "Premium accounts have been successfully removed.", "\033[38;5;39m")
}

func viewInstalledAccounts() {
	printLine("INFO", "Installed Accounts:", "\033[38;5;39m")
	for uuid, account := range accountsData.Accounts {
		printLine("ACCOUNT", fmt.Sprintf("%s: %s", uuid, account.Username), "\033[38;5;39m")
	}
}

func isValidMinecraftUsername(username string) bool {
	if len(username) < 3 || len(username) > 16 {
		return false
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, username)
	return matched
}

func isValidUUID(uuid string) bool {
	matched, _ := regexp.MatchString(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, uuid)
	return matched
}

func print(info, text, color string) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf(" [%s] > [%s] %s", timestamp, info, color+text+"\033[0m")
}

func printLine(info, text, color string) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf(" [%s] > [%s] %s%s\n", timestamp, info, color, text+"\033[0m")
}

func main() {
	err := loadJSON()
	if err != nil {
		printLine("ERROR", fmt.Sprintf("Failed to load accounts file: %v", err), "\033[31m")
		printLine("NOTICE", "Please check that you have Lunar Client installed.", "\033[31m")
		printLine("NOTICE", "Exiting in 3 seconds...", "\033[33m")
		time.Sleep(3 * time.Second)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	continueProgram := true

	for continueProgram {
		clearScreen()
		fmt.Printf("Cracked Lunar Account Tool (GOlang) %s\n\n", version)

		printMenu()

		fmt.Print("Please type your option (1-4) here: ")
		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			createAccountPrompt(scanner)
		case "2":
			removeAccountsMenu(scanner)
		case "3":
			viewInstalledAccounts()
		case "4":
			printLine("INFO", "Exiting the program.", "\033[38;5;39m")
			continueProgram = false
		default:
			printLine("ERROR", "Your choice is invalid. Please pick an option (1-4).", "\033[31m")
		}

		if continueProgram {
			printLine("INFO", "Press any key to return to the main menu...", "\033[38;5;39m")
			scanner.Scan()
		}
	}

	err = saveJSON()
	if err != nil {
		printLine("ERROR", fmt.Sprintf("Failed to save accounts file: %v", err), "\033[31m")
	}
}

func printMenu() {
	printLine("?", "What would you like to do:", "\033[38;5;39m")
	printLine("OPTION", "1. Create Account", "\033[38;5;39m")
	printLine("OPTION", "2. Remove Accounts", "\033[38;5;39m")
	printLine("OPTION", "3. View Installed Accounts", "\033[38;5;39m")
	printLine("OPTION", "4. Exit the program", "\033[38;5;39m")
}

func removeAccountsMenu(scanner *bufio.Scanner) {
	clearScreen()
	printLine("?", "Choose an option to remove accounts:", "\033[38;5;39m")
	printLine("OPTION", "1. Remove All Accounts", "\033[38;5;39m")
	printLine("OPTION", "2. Remove Cracked Accounts (accessToken is not a UUID)", "\033[38;5;39m")
	printLine("OPTION", "3. Remove Premium Accounts (accessToken is a UUID)", "\033[38;5;39m")
	fmt.Print("Please type your option (1-3) here: ")
	scanner.Scan()
	choice := strings.TrimSpace(scanner.Text())

	switch choice {
	case "1":
		removeAllAccounts()
	case "2":
		removeCrackedAccounts()
	case "3":
		removePremiumAccounts()
	default:
		printLine("ERROR", "Invalid option. Returning to main menu.", "\033[31m")
	}

	saveJSON()
}

func createAccountPrompt(scanner *bufio.Scanner) {
	fmt.Print("Enter your desired username: ")
	scanner.Scan()
	username := strings.TrimSpace(scanner.Text())

	if !isValidMinecraftUsername(username) {
		printLine("WARNING", "You may experience issues joining servers because of your username being invalid.", "\033[31m")
	}

	for {
		fmt.Print("Enter a valid UUID: ")
		scanner.Scan()
		uuid := strings.TrimSpace(scanner.Text())

		if !isValidUUID(uuid) {
			printLine("WARNING", "The UUID you entered is invalid. Please ensure it follows the correct format.", "\033[31m")
			fmt.Print("Would you like to try again? (y/n): ")
			scanner.Scan()
			retry := strings.ToLower(strings.TrimSpace(scanner.Text()))
			if retry == "n" {
				printLine("INFO", "Returning to main menu.", "\033[38;5;39m")
				return
			}
		} else {
			createAccount(username, uuid)
			saveJSON()
			break
		}
	}
}

func clearScreen() {
	fmt.Print("\033[2J\033[1;1H")
}