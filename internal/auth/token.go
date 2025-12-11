package auth

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/robertmeta/twist-cli/pkg/config"
)

func PromptForToken() (string, error) {
	fmt.Println("No Twist API token found.")
	fmt.Println("To get your personal access token:")
	fmt.Println("1. Go to https://twist.com/integrations")
	fmt.Println("2. Create a new integration or select an existing one")
	fmt.Println("3. Copy your personal access token from the OAuth section")
	fmt.Println()
	fmt.Print("Enter your Twist API token: ")

	reader := bufio.NewReader(os.Stdin)
	token, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read token: %w", err)
	}

	token = strings.TrimSpace(token)
	if token == "" {
		return "", fmt.Errorf("token cannot be empty")
	}

	return token, nil
}

func EnsureToken() (string, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.Token != "" {
		return cfg.Token, nil
	}

	token, err := PromptForToken()
	if err != nil {
		return "", err
	}

	cfg.Token = token
	if err := config.SaveConfig(cfg); err != nil {
		return "", fmt.Errorf("failed to save token: %w", err)
	}

	fmt.Println("Token saved successfully!")
	return token, nil
}
