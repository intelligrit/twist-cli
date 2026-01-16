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

var conversationsCmd = &cobra.Command{
	Use:   "conversations",
	Short: "Manage direct message conversations",
	Long:  `View and manage direct message conversations with other Twist users.`,
}

var conversationsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all conversations",
	Long:  `List all direct message conversations.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := auth.GetToken(tokenFlag)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		client := api.NewClient(token)
		conversations, err := client.GetConversations()
		if err != nil {
			return fmt.Errorf("failed to get conversations: %w", err)
		}

		if len(conversations) == 0 {
			fmt.Println("No conversations found.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tPARTICIPANTS\tMESSAGES\tCREATED")
		fmt.Fprintln(w, "--\t------------\t--------\t-------")
		for _, conv := range conversations {
			created := time.Unix(conv.CreatedTS, 0).Format("2006-01-02")
			userIDsStr := ""
			for i, uid := range conv.UserIDs {
				if i > 0 {
					userIDsStr += ", "
				}
				userIDsStr += fmt.Sprintf("%d", uid)
			}
			fmt.Fprintf(w, "%d\t%s\t%d\t%s\n", conv.ID, userIDsStr, conv.MessageCount, created)
		}
		w.Flush()

		return nil
	},
}

var conversationsShowCmd = &cobra.Command{
	Use:   "show [conversation-id]",
	Short: "Show messages in a conversation",
	Long:  `Display all messages in a direct message conversation.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		conversationID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid conversation ID: %w", err)
		}

		token, err := auth.GetToken(tokenFlag)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		client := api.NewClient(token)
		messages, err := client.GetConversationMessages(conversationID)
		if err != nil {
			return fmt.Errorf("failed to get messages: %w", err)
		}

		fmt.Println("================================================================================")
		fmt.Printf("Conversation #%d\n", conversationID)
		fmt.Println("================================================================================")

		if len(messages) == 0 {
			fmt.Println("\nNo messages yet.")
			return nil
		}

		for _, msg := range messages {
			fmt.Printf("\n[User %d] • %s\n", msg.UserID,
				time.Unix(msg.CreatedTS, 0).Format("2006-01-02 15:04:05"))
			fmt.Println(msg.Content)
		}

		return nil
	},
}

var conversationsSendCmd = &cobra.Command{
	Use:   "send [user-id] [message...]",
	Short: "Send a direct message",
	Long:  `Send a direct message to a user by user ID.`,
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		userID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid user ID: %w", err)
		}

		content := strings.Join(args[1:], " ")

		token, err := auth.GetToken(tokenFlag)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		client := api.NewClient(token)

		// Get or create conversation
		conversation, err := client.GetOrCreateConversation([]int{userID})
		if err != nil {
			return fmt.Errorf("failed to create conversation: %w", err)
		}

		// Send message
		message, err := client.SendConversationMessage(conversation.ID, content, nil)
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}

		fmt.Printf("Message sent successfully (message #%d in conversation #%d)\n",
			message.ID, conversation.ID)

		return nil
	},
}

var conversationsArchiveCmd = &cobra.Command{
	Use:   "archive [conversation-id]",
	Short: "Archive a conversation",
	Long:  `Archive a conversation to hide it from active lists.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		conversationID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid conversation ID: %w", err)
		}

		token, err := auth.GetToken(tokenFlag)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		client := api.NewClient(token)
		if err := client.ArchiveConversation(conversationID); err != nil {
			return fmt.Errorf("failed to archive conversation: %w", err)
		}

		fmt.Printf("Conversation %d archived successfully\n", conversationID)
		return nil
	},
}

var conversationsUnarchiveCmd = &cobra.Command{
	Use:   "unarchive [conversation-id]",
	Short: "Unarchive a conversation",
	Long:  `Unarchive a previously archived conversation.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		conversationID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid conversation ID: %w", err)
		}

		token, err := auth.GetToken(tokenFlag)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		client := api.NewClient(token)
		if err := client.UnarchiveConversation(conversationID); err != nil {
			return fmt.Errorf("failed to unarchive conversation: %w", err)
		}

		fmt.Printf("Conversation %d unarchived successfully\n", conversationID)
		return nil
	},
}

var conversationsMuteCmd = &cobra.Command{
	Use:   "mute [conversation-id]",
	Short: "Mute a conversation",
	Long:  `Mute notifications from a conversation.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		conversationID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid conversation ID: %w", err)
		}

		token, err := auth.GetToken(tokenFlag)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		client := api.NewClient(token)
		if err := client.MuteConversation(conversationID); err != nil {
			return fmt.Errorf("failed to mute conversation: %w", err)
		}

		fmt.Printf("Conversation %d muted successfully\n", conversationID)
		return nil
	},
}

var conversationsUnmuteCmd = &cobra.Command{
	Use:   "unmute [conversation-id]",
	Short: "Unmute a conversation",
	Long:  `Unmute notifications from a conversation.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		conversationID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid conversation ID: %w", err)
		}

		token, err := auth.GetToken(tokenFlag)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		client := api.NewClient(token)
		if err := client.UnmuteConversation(conversationID); err != nil {
			return fmt.Errorf("failed to unmute conversation: %w", err)
		}

		fmt.Printf("Conversation %d unmuted successfully\n", conversationID)
		return nil
	},
}

var conversationsMarkReadCmd = &cobra.Command{
	Use:   "mark-read [conversation-id]",
	Short: "Mark a conversation as read",
	Long:  `Mark all messages in a conversation as read.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		conversationID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid conversation ID: %w", err)
		}

		token, err := auth.GetToken(tokenFlag)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		client := api.NewClient(token)
		if err := client.MarkConversationRead(conversationID); err != nil {
			return fmt.Errorf("failed to mark conversation as read: %w", err)
		}

		fmt.Printf("Conversation %d marked as read\n", conversationID)
		return nil
	},
}

var conversationsMarkUnreadCmd = &cobra.Command{
	Use:   "mark-unread [conversation-id]",
	Short: "Mark a conversation as unread",
	Long:  `Mark a conversation as unread.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		conversationID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid conversation ID: %w", err)
		}

		token, err := auth.GetToken(tokenFlag)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		client := api.NewClient(token)
		if err := client.MarkConversationUnread(conversationID); err != nil {
			return fmt.Errorf("failed to mark conversation as unread: %w", err)
		}

		fmt.Printf("Conversation %d marked as unread\n", conversationID)
		return nil
	},
}

func init() {
	conversationsCmd.AddCommand(conversationsListCmd)
	conversationsCmd.AddCommand(conversationsShowCmd)
	conversationsCmd.AddCommand(conversationsSendCmd)
	conversationsCmd.AddCommand(conversationsArchiveCmd)
	conversationsCmd.AddCommand(conversationsUnarchiveCmd)
	conversationsCmd.AddCommand(conversationsMuteCmd)
	conversationsCmd.AddCommand(conversationsUnmuteCmd)
	conversationsCmd.AddCommand(conversationsMarkReadCmd)
	conversationsCmd.AddCommand(conversationsMarkUnreadCmd)
}
