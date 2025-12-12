package cmd

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/robertmeta/twist-cli/internal/auth"
	"github.com/robertmeta/twist-cli/pkg/api"
	"github.com/spf13/cobra"
)

var reactionsCmd = &cobra.Command{
	Use:   "reactions",
	Short: "Manage reactions",
	Long:  `Add, remove, and view reactions on threads and comments.`,
}

var reactionsAddCmd = &cobra.Command{
	Use:   "add [target-type] [target-id] [emoji]",
	Short: "Add a reaction",
	Long:  `Add an emoji reaction to a thread or comment. Target type must be 'thread' or 'comment'.`,
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetType := args[0]
		targetID, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid target ID: %w", err)
		}
		emoji := args[2]

		if targetType != "thread" && targetType != "comment" {
			return fmt.Errorf("invalid target type: must be 'thread' or 'comment'")
		}

		token, err := auth.GetToken(tokenFlag)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		client := api.NewClient(token)
		reaction, err := client.AddReaction(targetType, targetID, emoji)
		if err != nil {
			return fmt.Errorf("failed to add reaction: %w", err)
		}

		fmt.Printf("Reaction added successfully!\n")
		fmt.Printf("Reaction ID: %d\n", reaction.ID)
		fmt.Printf("Emoji: %s\n", reaction.Emoji)

		return nil
	},
}

var reactionsRemoveCmd = &cobra.Command{
	Use:   "remove [target-type] [target-id] [emoji]",
	Short: "Remove a reaction",
	Long:  `Remove an emoji reaction from a thread or comment. Target type must be 'thread' or 'comment'.`,
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetType := args[0]
		targetID, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid target ID: %w", err)
		}
		emoji := args[2]

		if targetType != "thread" && targetType != "comment" {
			return fmt.Errorf("invalid target type: must be 'thread' or 'comment'")
		}

		token, err := auth.GetToken(tokenFlag)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		client := api.NewClient(token)
		if err := client.RemoveReaction(targetType, targetID, emoji); err != nil {
			return fmt.Errorf("failed to remove reaction: %w", err)
		}

		fmt.Printf("Reaction removed successfully\n")

		return nil
	},
}

var reactionsListCmd = &cobra.Command{
	Use:   "list [target-type] [target-id]",
	Short: "List all reactions",
	Long:  `List all reactions on a thread or comment. Target type must be 'thread' or 'comment'.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetType := args[0]
		targetID, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid target ID: %w", err)
		}

		if targetType != "thread" && targetType != "comment" {
			return fmt.Errorf("invalid target type: must be 'thread' or 'comment'")
		}

		token, err := auth.GetToken(tokenFlag)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		client := api.NewClient(token)
		reactions, err := client.GetReactions(targetType, targetID)
		if err != nil {
			return fmt.Errorf("failed to get reactions: %w", err)
		}

		if len(reactions) == 0 {
			fmt.Println("No reactions found.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tEMOJI\tUSER ID")
		fmt.Fprintln(w, "--\t-----\t-------")
		for _, r := range reactions {
			fmt.Fprintf(w, "%d\t%s\t%d\n", r.ID, r.Emoji, r.UserID)
		}
		w.Flush()

		return nil
	},
}

func init() {
	reactionsCmd.AddCommand(reactionsAddCmd)
	reactionsCmd.AddCommand(reactionsRemoveCmd)
	reactionsCmd.AddCommand(reactionsListCmd)
}
