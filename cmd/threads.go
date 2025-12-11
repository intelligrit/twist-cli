package cmd

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/robertmeta/twist-cli/internal/auth"
	"github.com/robertmeta/twist-cli/pkg/api"
	"github.com/spf13/cobra"
)

var threadsCmd = &cobra.Command{
	Use:   "threads",
	Short: "Manage Twist threads",
	Long:  `View and manage threads (conversations) in Twist channels.`,
}

var threadsListCmd = &cobra.Command{
	Use:   "list [channel-id]",
	Short: "List all threads in a channel",
	Long:  `List all threads in a specific channel by providing the channel ID.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		channelID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid channel ID: %w", err)
		}

		token, err := auth.GetToken(tokenFlag)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		client := api.NewClient(token)
		threads, err := client.GetThreads(channelID)
		if err != nil {
			return fmt.Errorf("failed to get threads: %w", err)
		}

		if len(threads) == 0 {
			fmt.Println("No threads found in this channel.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tTITLE\tCOMMENTS\tLAST UPDATED")
		fmt.Fprintln(w, "--\t-----\t--------\t------------")
		for _, t := range threads {
			lastUpdated := time.Unix(t.LastUpdatedTS, 0).Format("2006-01-02 15:04")
			title := t.Title
			if title == "" {
				title = "(no title)"
			}
			if len(title) > 50 {
				title = title[:47] + "..."
			}
			fmt.Fprintf(w, "%d\t%s\t%d\t%s\n", t.ID, title, t.CommentCount, lastUpdated)
		}
		w.Flush()

		return nil
	},
}

var threadsShowCmd = &cobra.Command{
	Use:   "show [thread-id]",
	Short: "Show a thread with its content and replies",
	Long:  `Display the full content of a thread including all comments/replies.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		threadID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid thread ID: %w", err)
		}

		token, err := auth.GetToken(tokenFlag)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		client := api.NewClient(token)

		thread, err := client.GetThread(threadID)
		if err != nil {
			return fmt.Errorf("failed to get thread: %w", err)
		}

		comments, err := client.GetComments(threadID)
		if err != nil {
			return fmt.Errorf("failed to get comments: %w", err)
		}

		fmt.Println("================================================================================")
		fmt.Printf("Thread #%d: %s\n", thread.ID, thread.Title)
		fmt.Printf("Posted: %s\n", time.Unix(thread.PostedTS, 0).Format("2006-01-02 15:04:05"))
		fmt.Printf("Comments: %d\n", thread.CommentCount)
		fmt.Println("================================================================================")
		fmt.Println()
		fmt.Println(thread.Content)
		fmt.Println()

		if len(comments) > 0 {
			fmt.Println("--------------------------------------------------------------------------------")
			fmt.Printf("Replies (%d):\n", len(comments))
			fmt.Println("--------------------------------------------------------------------------------")
			for i, comment := range comments {
				fmt.Printf("\n[%d] User %d • %s\n", i+1, comment.Creator,
					time.Unix(comment.PostedTS, 0).Format("2006-01-02 15:04:05"))
				fmt.Println(comment.Content)
			}
		}

		return nil
	},
}

var threadsReplyCmd = &cobra.Command{
	Use:   "reply [thread-id] [message]",
	Short: "Reply to a thread",
	Long:  `Post a comment/reply to an existing thread.`,
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		threadID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid thread ID: %w", err)
		}

		content := ""
		for i := 1; i < len(args); i++ {
			if i > 1 {
				content += " "
			}
			content += args[i]
		}

		token, err := auth.GetToken(tokenFlag)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		client := api.NewClient(token)
		comment, err := client.PostComment(threadID, content)
		if err != nil {
			return fmt.Errorf("failed to post reply: %w", err)
		}

		fmt.Printf("Reply posted successfully (comment #%d)\n", comment.ID)
		return nil
	},
}

func init() {
	threadsCmd.AddCommand(threadsListCmd)
	threadsCmd.AddCommand(threadsShowCmd)
	threadsCmd.AddCommand(threadsReplyCmd)
}
