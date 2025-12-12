package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	tokenFlag string
)

var rootCmd = &cobra.Command{
	Use:   "twist",
	Short: "Twist CLI - Command line interface for Twist",
	Long: `A command line tool for interacting with the Twist API.
Authenticate using your personal access token to manage workspaces,
channels, and conversations.`,
	Version: "0.1.0",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&tokenFlag, "token", "", "Twist API token (or set TWIST_API_TOKEN env var)")
	rootCmd.AddCommand(workspacesCmd)
	rootCmd.AddCommand(channelsCmd)
	rootCmd.AddCommand(threadsCmd)
	rootCmd.AddCommand(commentsCmd)
	rootCmd.AddCommand(reactionsCmd)
	rootCmd.AddCommand(conversationsCmd)
	rootCmd.AddCommand(groupsCmd)
	rootCmd.AddCommand(attachmentsCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(usersCmd)
}
