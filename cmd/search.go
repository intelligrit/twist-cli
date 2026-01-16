package cmd

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/intelligrit/twist-cli/internal/auth"
	"github.com/intelligrit/twist-cli/pkg/api"
	"github.com/spf13/cobra"
)

var (
	searchChannelIDFlag int
	searchLimitFlag     int
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search content",
	Long:  `Search threads, messages, and conversations.`,
}

var searchThreadsCmd = &cobra.Command{
	Use:   "threads [workspace-id] [query]",
	Short: "Search threads",
	Long:  `Search for threads in a workspace. Use --channel-id to limit to a specific channel.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		workspaceID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid workspace ID: %w", err)
		}
		query := args[1]

		token, err := auth.GetToken(tokenFlag)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		opts := make(map[string]interface{})
		if searchChannelIDFlag > 0 {
			opts["channel_id"] = searchChannelIDFlag
		}
		if searchLimitFlag > 0 {
			opts["limit"] = searchLimitFlag
		}

		client := api.NewClient(token)
		threads, err := client.SearchThreads(workspaceID, query, opts)
		if err != nil {
			return fmt.Errorf("failed to search threads: %w", err)
		}

		if len(threads) == 0 {
			fmt.Println("No threads found matching the query.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tTITLE\tCHANNEL\tLAST UPDATED")
		fmt.Fprintln(w, "--\t-----\t-------\t------------")
		for _, t := range threads {
			lastUpdated := time.Unix(t.LastUpdatedTS, 0).Format("2006-01-02 15:04")
			title := t.Title
			if title == "" {
				title = "(no title)"
			}
			if len(title) > 40 {
				title = title[:37] + "..."
			}
			fmt.Fprintf(w, "%d\t%s\t%d\t%s\n", t.ID, title, t.ChannelID, lastUpdated)
		}
		w.Flush()

		fmt.Printf("\nFound %d thread(s)\n", len(threads))

		return nil
	},
}

var searchMessagesCmd = &cobra.Command{
	Use:   "messages [workspace-id] [query]",
	Short: "Search messages",
	Long:  `Search for messages/comments in a workspace.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		workspaceID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid workspace ID: %w", err)
		}
		query := args[1]

		token, err := auth.GetToken(tokenFlag)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		opts := make(map[string]interface{})
		if searchLimitFlag > 0 {
			opts["limit"] = searchLimitFlag
		}

		client := api.NewClient(token)
		comments, err := client.SearchMessages(workspaceID, query, opts)
		if err != nil {
			return fmt.Errorf("failed to search messages: %w", err)
		}

		if len(comments) == 0 {
			fmt.Println("No messages found matching the query.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tTHREAD\tCONTENT\tPOSTED")
		fmt.Fprintln(w, "--\t------\t-------\t------")
		for _, c := range comments {
			posted := time.Unix(c.PostedTS, 0).Format("2006-01-02 15:04")
			content := c.Content
			if len(content) > 50 {
				content = content[:47] + "..."
			}
			fmt.Fprintf(w, "%d\t%d\t%s\t%s\n", c.ID, c.ThreadID, content, posted)
		}
		w.Flush()

		fmt.Printf("\nFound %d message(s)\n", len(comments))

		return nil
	},
}

var searchConversationsCmd = &cobra.Command{
	Use:   "conversations [query]",
	Short: "Search conversation messages",
	Long:  `Search for messages in direct message conversations.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]

		token, err := auth.GetToken(tokenFlag)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		opts := make(map[string]interface{})
		if searchLimitFlag > 0 {
			opts["limit"] = searchLimitFlag
		}

		client := api.NewClient(token)
		messages, err := client.SearchConversations(query, opts)
		if err != nil {
			return fmt.Errorf("failed to search conversations: %w", err)
		}

		if len(messages) == 0 {
			fmt.Println("No conversation messages found matching the query.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tCONVERSATION\tCONTENT\tPOSTED")
		fmt.Fprintln(w, "--\t------------\t-------\t------")
		for _, m := range messages {
			posted := time.Unix(m.CreatedTS, 0).Format("2006-01-02 15:04")
			content := m.Content
			if len(content) > 50 {
				content = content[:47] + "..."
			}
			fmt.Fprintf(w, "%d\t%d\t%s\t%s\n", m.ID, m.ConversationID, content, posted)
		}
		w.Flush()

		fmt.Printf("\nFound %d message(s)\n", len(messages))

		return nil
	},
}

func init() {
	searchThreadsCmd.Flags().IntVar(&searchChannelIDFlag, "channel-id", 0, "Limit search to specific channel")
	searchThreadsCmd.Flags().IntVar(&searchLimitFlag, "limit", 0, "Maximum number of results")

	searchMessagesCmd.Flags().IntVar(&searchLimitFlag, "limit", 0, "Maximum number of results")

	searchConversationsCmd.Flags().IntVar(&searchLimitFlag, "limit", 0, "Maximum number of results")

	searchCmd.AddCommand(searchThreadsCmd)
	searchCmd.AddCommand(searchMessagesCmd)
	searchCmd.AddCommand(searchConversationsCmd)
}
