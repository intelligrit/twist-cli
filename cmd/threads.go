package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/intelligrit/twist-cli/internal/auth"
	"github.com/intelligrit/twist-cli/pkg/api"
	"github.com/spf13/cobra"
)

var (
	createNotifyFlag string
	replyNotifyFlag  string
	titleFlag        string
	contentFlag      string
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
	Long:  `Post a comment/reply to an existing thread. Use --notify to specify user IDs to notify (comma-separated).`,
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		threadID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid thread ID: %w", err)
		}

		content := strings.Join(args[1:], " ")

		token, err := auth.GetToken(tokenFlag)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		var recipients []int
		if replyNotifyFlag != "" {
			userIDs := strings.Split(replyNotifyFlag, ",")
			for _, idStr := range userIDs {
				idStr = strings.TrimSpace(idStr)
				if idStr == "" {
					continue
				}
				id, err := strconv.Atoi(idStr)
				if err != nil {
					return fmt.Errorf("invalid user ID: %s", idStr)
				}
				recipients = append(recipients, id)
			}
		}

		client := api.NewClient(token)
		comment, err := client.PostComment(threadID, content, recipients)
		if err != nil {
			return fmt.Errorf("failed to post reply: %w", err)
		}

		fmt.Printf("Reply posted successfully (comment #%d)\n", comment.ID)
		if len(recipients) > 0 {
			fmt.Printf("Notified %d user(s)\n", len(recipients))
		}
		return nil
	},
}

var threadsCreateCmd = &cobra.Command{
	Use:   "create [channel-id] [title] [content]",
	Short: "Create a new thread",
	Long:  `Create a new thread in a channel. Use --notify to specify user IDs to notify (comma-separated).`,
	Args:  cobra.MinimumNArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		channelID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid channel ID: %w", err)
		}

		title := args[1]
		content := strings.Join(args[2:], " ")

		token, err := auth.GetToken(tokenFlag)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		var recipients []int
		if createNotifyFlag != "" {
			userIDs := strings.Split(createNotifyFlag, ",")
			for _, idStr := range userIDs {
				idStr = strings.TrimSpace(idStr)
				if idStr == "" {
					continue
				}
				id, err := strconv.Atoi(idStr)
				if err != nil {
					return fmt.Errorf("invalid user ID: %s", idStr)
				}
				recipients = append(recipients, id)
			}
		}

		client := api.NewClient(token)
		thread, err := client.CreateThread(channelID, title, content, recipients)
		if err != nil {
			return fmt.Errorf("failed to create thread: %w", err)
		}

		fmt.Printf("Thread created successfully!\n")
		fmt.Printf("Thread ID: %d\n", thread.ID)
		fmt.Printf("Title: %s\n", thread.Title)
		if len(recipients) > 0 {
			fmt.Printf("Notified %d user(s)\n", len(recipients))
		}

		return nil
	},
}

var threadsUpdateCmd = &cobra.Command{
	Use:   "update [thread-id]",
	Short: "Update a thread",
	Long:  `Update thread title and/or content. Use flags to specify what to update.`,
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

		updates := make(map[string]interface{})
		if titleFlag != "" {
			updates["title"] = titleFlag
		}
		if contentFlag != "" {
			updates["content"] = contentFlag
		}

		if len(updates) == 0 {
			return fmt.Errorf("no updates specified; use --title or --content flags")
		}

		client := api.NewClient(token)
		thread, err := client.UpdateThread(threadID, updates)
		if err != nil {
			return fmt.Errorf("failed to update thread: %w", err)
		}

		fmt.Printf("Thread updated successfully!\n")
		fmt.Printf("Thread ID: %d\n", thread.ID)
		fmt.Printf("Title: %s\n", thread.Title)

		return nil
	},
}

var threadsDeleteCmd = &cobra.Command{
	Use:   "delete [thread-id]",
	Short: "Delete a thread",
	Long:  `Delete a thread permanently.`,
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
		if err := client.DeleteThread(threadID); err != nil {
			return fmt.Errorf("failed to delete thread: %w", err)
		}

		fmt.Printf("Thread %d deleted successfully\n", threadID)
		return nil
	},
}

var threadsPinCmd = &cobra.Command{
	Use:   "pin [thread-id]",
	Short: "Pin a thread",
	Long:  `Pin a thread to the top of the channel.`,
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
		if err := client.PinThread(threadID); err != nil {
			return fmt.Errorf("failed to pin thread: %w", err)
		}

		fmt.Printf("Thread %d pinned successfully\n", threadID)
		return nil
	},
}

var threadsUnpinCmd = &cobra.Command{
	Use:   "unpin [thread-id]",
	Short: "Unpin a thread",
	Long:  `Remove the pin from a thread.`,
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
		if err := client.UnpinThread(threadID); err != nil {
			return fmt.Errorf("failed to unpin thread: %w", err)
		}

		fmt.Printf("Thread %d unpinned successfully\n", threadID)
		return nil
	},
}

var threadsStarCmd = &cobra.Command{
	Use:   "star [thread-id]",
	Short: "Star a thread",
	Long:  `Add a star to a thread for quick access.`,
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
		if err := client.StarThread(threadID); err != nil {
			return fmt.Errorf("failed to star thread: %w", err)
		}

		fmt.Printf("Thread %d starred successfully\n", threadID)
		return nil
	},
}

var threadsUnstarCmd = &cobra.Command{
	Use:   "unstar [thread-id]",
	Short: "Unstar a thread",
	Long:  `Remove the star from a thread.`,
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
		if err := client.UnstarThread(threadID); err != nil {
			return fmt.Errorf("failed to unstar thread: %w", err)
		}

		fmt.Printf("Thread %d unstarred successfully\n", threadID)
		return nil
	},
}

var threadsArchiveCmd = &cobra.Command{
	Use:   "archive [thread-id]",
	Short: "Archive a thread",
	Long:  `Archive a thread to hide it from active lists.`,
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
		if err := client.ArchiveThread(threadID); err != nil {
			return fmt.Errorf("failed to archive thread: %w", err)
		}

		fmt.Printf("Thread %d archived successfully\n", threadID)
		return nil
	},
}

var threadsUnarchiveCmd = &cobra.Command{
	Use:   "unarchive [thread-id]",
	Short: "Unarchive a thread",
	Long:  `Unarchive a previously archived thread.`,
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
		if err := client.UnarchiveThread(threadID); err != nil {
			return fmt.Errorf("failed to unarchive thread: %w", err)
		}

		fmt.Printf("Thread %d unarchived successfully\n", threadID)
		return nil
	},
}

var commentsUpdateCmd = &cobra.Command{
	Use:   "update [comment-id] [content...]",
	Short: "Update a comment",
	Long:  `Update the content of a comment.`,
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		commentID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid comment ID: %w", err)
		}

		content := strings.Join(args[1:], " ")

		token, err := auth.GetToken(tokenFlag)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		client := api.NewClient(token)
		comment, err := client.UpdateComment(commentID, content)
		if err != nil {
			return fmt.Errorf("failed to update comment: %w", err)
		}

		fmt.Printf("Comment updated successfully!\n")
		fmt.Printf("Comment ID: %d\n", comment.ID)

		return nil
	},
}

var commentsDeleteCmd = &cobra.Command{
	Use:   "delete [comment-id]",
	Short: "Delete a comment",
	Long:  `Delete a comment permanently.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		commentID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid comment ID: %w", err)
		}

		token, err := auth.GetToken(tokenFlag)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		client := api.NewClient(token)
		if err := client.DeleteComment(commentID); err != nil {
			return fmt.Errorf("failed to delete comment: %w", err)
		}

		fmt.Printf("Comment %d deleted successfully\n", commentID)
		return nil
	},
}

var commentsCmd = &cobra.Command{
	Use:   "comments",
	Short: "Manage comments",
	Long:  `Manage comments on threads.`,
}

func init() {
	threadsUpdateCmd.Flags().StringVar(&titleFlag, "title", "", "Thread title")
	threadsUpdateCmd.Flags().StringVar(&contentFlag, "content", "", "Thread content")

	threadsCreateCmd.Flags().StringVar(&createNotifyFlag, "notify", "", "Comma-separated user IDs to notify")
	threadsReplyCmd.Flags().StringVar(&replyNotifyFlag, "notify", "", "Comma-separated user IDs to notify")

	threadsCmd.AddCommand(threadsListCmd)
	threadsCmd.AddCommand(threadsShowCmd)
	threadsCmd.AddCommand(threadsReplyCmd)
	threadsCmd.AddCommand(threadsCreateCmd)
	threadsCmd.AddCommand(threadsUpdateCmd)
	threadsCmd.AddCommand(threadsDeleteCmd)
	threadsCmd.AddCommand(threadsPinCmd)
	threadsCmd.AddCommand(threadsUnpinCmd)
	threadsCmd.AddCommand(threadsStarCmd)
	threadsCmd.AddCommand(threadsUnstarCmd)
	threadsCmd.AddCommand(threadsArchiveCmd)
	threadsCmd.AddCommand(threadsUnarchiveCmd)

	commentsCmd.AddCommand(commentsUpdateCmd)
	commentsCmd.AddCommand(commentsDeleteCmd)
}
