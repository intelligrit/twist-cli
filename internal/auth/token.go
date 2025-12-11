package auth

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func PromptForToken() (string, error) {
	fmt.Fprintln(os.Stderr, "No Twist API token found.")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Please provide your token using one of these methods:")
	fmt.Fprintln(os.Stderr, "  1. Set TWIST_API_TOKEN environment variable")
	fmt.Fprintln(os.Stderr, "  2. Use --token flag")
	fmt.Fprintln(os.Stderr, "  3. Enter it now (not saved)")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "To get your personal access token:")
	fmt.Fprintln(os.Stderr, "  - Go to https://twist.com/integrations")
	fmt.Fprintln(os.Stderr, "  - Create a new integration or select an existing one")
	fmt.Fprintln(os.Stderr, "  - Copy your personal access token from the OAuth section")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprint(os.Stderr, "Enter your Twist API token: ")

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

func GetToken(flagToken string) (string, error) {
	// Priority 1: --token flag
	if flagToken != "" {
		return flagToken, nil
	}

	// Priority 2: TWIST_API_TOKEN environment variable
	if envToken := os.Getenv("TWIST_API_TOKEN"); envToken != "" {
		return envToken, nil
	}

	// Priority 3: Prompt user (not saved)
	return PromptForToken()
}
