package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/intelligrit/twist-cli/internal/auth"
	"github.com/intelligrit/twist-cli/pkg/api"
	"github.com/spf13/cobra"
)

var (
	archivedFlag    bool
	descriptionFlag string
	colorFlag       int
	iconFlag        int
	publicFlag      bool
	userIDsFlag     string
	nameFlag        string
)

var channelsCmd = &cobra.Command{
	Use:   "channels",
	Short: "Manage Twist channels",
	Long:  `View and manage channels in Twist workspaces.`,
}

var channelsListCmd = &cobra.Command{
	Use:   "list [workspace-id]",
	Short: "List all channels in a workspace",
	Long:  `List all channels in a specific workspace. Use --archived to show only archived channels.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		workspaceID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid workspace ID: %w", err)
		}

		token, err := auth.GetToken(tokenFlag)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		client := api.NewClient(token)
		channels, err := client.GetChannels(workspaceID, archivedFlag)
		if err != nil {
			return fmt.Errorf("failed to get channels: %w", err)
		}

		if len(channels) == 0 {
			fmt.Println("No channels found.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tPUBLIC\tARCHIVED")
		fmt.Fprintln(w, "--\t----\t------\t--------")
		for _, ch := range channels {
			fmt.Fprintf(w, "%d\t%s\t%t\t%t\n", ch.ID, ch.Name, ch.Public, ch.Archived)
		}
		w.Flush()

		return nil
	},
}

var channelsShowCmd = &cobra.Command{
	Use:   "show [channel-id]",
	Short: "Show channel details",
	Long:  `Display detailed information about a specific channel.`,
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
		channel, err := client.GetChannel(channelID)
		if err != nil {
			return fmt.Errorf("failed to get channel: %w", err)
		}

		fmt.Printf("ID: %d\n", channel.ID)
		fmt.Printf("Name: %s\n", channel.Name)
		fmt.Printf("Description: %s\n", channel.Description)
		fmt.Printf("Workspace ID: %d\n", channel.WorkspaceID)
		fmt.Printf("Public: %t\n", channel.Public)
		fmt.Printf("Archived: %t\n", channel.Archived)
		fmt.Printf("Color: %d\n", channel.Color)
		fmt.Printf("Icon: %d\n", channel.Icon)

		return nil
	},
}

var channelsCreateCmd = &cobra.Command{
	Use:   "create [workspace-id] [name]",
	Short: "Create a new channel",
	Long:  `Create a new channel in a workspace. Use flags to set optional properties.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		workspaceID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid workspace ID: %w", err)
		}

		name := args[1]

		token, err := auth.GetToken(tokenFlag)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		opts := make(map[string]interface{})
		if descriptionFlag != "" {
			opts["description"] = descriptionFlag
		}
		if colorFlag >= 0 {
			opts["color"] = colorFlag
		}
		if iconFlag > 0 {
			opts["icon"] = iconFlag
		}
		if cmd.Flags().Changed("public") {
			opts["public"] = publicFlag
		}
		if userIDsFlag != "" {
			userIDs := []int{}
			for _, idStr := range strings.Split(userIDsFlag, ",") {
				idStr = strings.TrimSpace(idStr)
				if idStr == "" {
					continue
				}
				id, err := strconv.Atoi(idStr)
				if err != nil {
					return fmt.Errorf("invalid user ID: %s", idStr)
				}
				userIDs = append(userIDs, id)
			}
			if len(userIDs) > 0 {
				opts["user_ids"] = userIDs
			}
		}

		client := api.NewClient(token)
		channel, err := client.CreateChannel(workspaceID, name, opts)
		if err != nil {
			return fmt.Errorf("failed to create channel: %w", err)
		}

		fmt.Printf("Channel created successfully!\n")
		fmt.Printf("Channel ID: %d\n", channel.ID)
		fmt.Printf("Name: %s\n", channel.Name)

		return nil
	},
}

var channelsUpdateCmd = &cobra.Command{
	Use:   "update [channel-id]",
	Short: "Update a channel",
	Long:  `Update channel properties. Use flags to specify what to update.`,
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

		updates := make(map[string]interface{})
		if nameFlag != "" {
			updates["name"] = nameFlag
		}
		if descriptionFlag != "" {
			updates["description"] = descriptionFlag
		}
		if colorFlag >= 0 {
			updates["color"] = colorFlag
		}
		if cmd.Flags().Changed("public") {
			updates["public"] = publicFlag
		}

		if len(updates) == 0 {
			return fmt.Errorf("no updates specified; use flags like --name, --description, --color, or --public")
		}

		client := api.NewClient(token)
		channel, err := client.UpdateChannel(channelID, updates)
		if err != nil {
			return fmt.Errorf("failed to update channel: %w", err)
		}

		fmt.Printf("Channel updated successfully!\n")
		fmt.Printf("Channel ID: %d\n", channel.ID)
		fmt.Printf("Name: %s\n", channel.Name)

		return nil
	},
}

var channelsArchiveCmd = &cobra.Command{
	Use:   "archive [channel-id]",
	Short: "Archive a channel",
	Long:  `Archive a channel. Archived channels are hidden from active channel lists.`,
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
		if err := client.ArchiveChannel(channelID); err != nil {
			return fmt.Errorf("failed to archive channel: %w", err)
		}

		fmt.Printf("Channel %d archived successfully\n", channelID)
		return nil
	},
}

var channelsUnarchiveCmd = &cobra.Command{
	Use:   "unarchive [channel-id]",
	Short: "Unarchive a channel",
	Long:  `Unarchive a previously archived channel.`,
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
		if err := client.UnarchiveChannel(channelID); err != nil {
			return fmt.Errorf("failed to unarchive channel: %w", err)
		}

		fmt.Printf("Channel %d unarchived successfully\n", channelID)
		return nil
	},
}

var channelsDeleteCmd = &cobra.Command{
	Use:   "delete [channel-id]",
	Short: "Delete an archived channel",
	Long:  `Delete a channel. The channel must be archived first before deletion.`,
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
		if err := client.DeleteChannel(channelID); err != nil {
			return fmt.Errorf("failed to delete channel: %w", err)
		}

		fmt.Printf("Channel %d deleted successfully\n", channelID)
		return nil
	},
}

var channelsAddUserCmd = &cobra.Command{
	Use:   "add-user [channel-id] [user-id]",
	Short: "Add a user to a channel",
	Long:  `Add a user to a channel by user ID.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		channelID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid channel ID: %w", err)
		}

		userID, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid user ID: %w", err)
		}

		token, err := auth.GetToken(tokenFlag)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		client := api.NewClient(token)
		if err := client.AddChannelUser(channelID, userID); err != nil {
			return fmt.Errorf("failed to add user to channel: %w", err)
		}

		fmt.Printf("User %d added to channel %d successfully\n", userID, channelID)
		return nil
	},
}

var channelsRemoveUserCmd = &cobra.Command{
	Use:   "remove-user [channel-id] [user-id]",
	Short: "Remove a user from a channel",
	Long:  `Remove a user from a channel by user ID.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		channelID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid channel ID: %w", err)
		}

		userID, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid user ID: %w", err)
		}

		token, err := auth.GetToken(tokenFlag)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		client := api.NewClient(token)
		if err := client.RemoveChannelUser(channelID, userID); err != nil {
			return fmt.Errorf("failed to remove user from channel: %w", err)
		}

		fmt.Printf("User %d removed from channel %d successfully\n", userID, channelID)
		return nil
	},
}

func init() {
	channelsListCmd.Flags().BoolVar(&archivedFlag, "archived", false, "Show only archived channels")

	channelsCreateCmd.Flags().StringVar(&descriptionFlag, "description", "", "Channel description")
	channelsCreateCmd.Flags().IntVar(&colorFlag, "color", -1, "Channel color (0-11)")
	channelsCreateCmd.Flags().IntVar(&iconFlag, "icon", 0, "Channel icon (1-255)")
	channelsCreateCmd.Flags().BoolVar(&publicFlag, "public", false, "Make channel public")
	channelsCreateCmd.Flags().StringVar(&userIDsFlag, "user-ids", "", "Comma-separated user IDs to add")

	channelsUpdateCmd.Flags().StringVar(&nameFlag, "name", "", "Channel name")
	channelsUpdateCmd.Flags().StringVar(&descriptionFlag, "description", "", "Channel description")
	channelsUpdateCmd.Flags().IntVar(&colorFlag, "color", -1, "Channel color (0-11)")
	channelsUpdateCmd.Flags().BoolVar(&publicFlag, "public", false, "Make channel public")

	channelsCmd.AddCommand(channelsListCmd)
	channelsCmd.AddCommand(channelsShowCmd)
	channelsCmd.AddCommand(channelsCreateCmd)
	channelsCmd.AddCommand(channelsUpdateCmd)
	channelsCmd.AddCommand(channelsArchiveCmd)
	channelsCmd.AddCommand(channelsUnarchiveCmd)
	channelsCmd.AddCommand(channelsDeleteCmd)
	channelsCmd.AddCommand(channelsAddUserCmd)
	channelsCmd.AddCommand(channelsRemoveUserCmd)
}
