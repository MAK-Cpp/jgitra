package config

import (
	"bufio"
	"fmt"
	"os"
	"unicode/utf8"

	"golang.org/x/term"
)

func initConfig(currConf *Config) error {
	fmt.Println("Initializing configuration: ")
	fmt.Println("If some value is exists, leave it empty to use previous configuration")

	initBaseURL(currConf)
	initUserEmail(currConf)
	if err := initAPIToken(currConf); err != nil {
		return err
	}

	return nil
}

func initBaseURL(config *Config) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter base url of your jira (like https://example.atlassian.net)")
	if config.Jira.BaseURL != "" {
		fmt.Printf(" [current: %s]", config.Jira.BaseURL)
	}
	fmt.Print(": ")
	scanner.Scan()
	if baseURL := scanner.Text(); baseURL != "" {
		config.Jira.BaseURL = baseURL
	}
}

func initUserEmail(config *Config) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter user email")
	if config.Jira.UserEmail != "" {
		fmt.Printf(" [current: %s]", config.Jira.UserEmail)
	}
	fmt.Print(": ")
	scanner.Scan()
	if email := scanner.Text(); email != "" {
		config.Jira.UserEmail = email
	}
}

func initAPIToken(config *Config) error {
	fmt.Print("Enter API token")
	if config.Jira.APIToken != "" {
		tokenLen := utf8.RuneCountInString(config.Jira.APIToken)
		fmt.Printf(
			" [current: %s...%s]",
			config.Jira.APIToken[0:2],
			config.Jira.APIToken[tokenLen-2:tokenLen],
		)
	}
	fmt.Print(": ")
	APITokenBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println("")
	if err != nil {
		return err
	}
	config.Jira.APIToken = string(APITokenBytes)
	return nil
}
