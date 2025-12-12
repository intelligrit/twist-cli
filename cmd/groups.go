package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/robertmeta/twist-cli/internal/auth"
	"github.com/robertmeta/twist-cli/pkg/api"
	"github.com/spf13/cobra"
)

var (
	groupDescriptionFlag string
	groupUserIDsFlag     string
	groupNameFlag        string
)

var groupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "Manage groups",
	Long:  `View and manage groups in Twist workspaces.`,
}

var groupsListCmd = &cobra.Command{
	Use:   "list [workspace-id]",
	Short: "List all groups in a workspace",
	Long:  `List all groups in a specific workspace.`,
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
		groups, err := client.GetGroups(workspaceID)
		if err != nil {
			return fmt.Errorf("failed to get groups: %w", err)
		}

		if len(groups) == 0 {
			fmt.Println("No groups found.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tMEMBERS")
		fmt.Fprintln(w, "--\t----\t-------")
		for _, g := range groups {
			fmt.Fprintf(w, "%d\t%s\t%d\n", g.ID, g.Name, len(g.UserIDs))
		}
		w.Flush()

		return nil
	},
}

var groupsShowCmd = &cobra.Command{
	Use:   "show [group-id]",
	Short: "Show group details",
	Long:  `Display detailed information about a specific group.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		groupID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid group ID: %w", err)
		}

		token, err := auth.GetToken(tokenFlag)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		client := api.NewClient(token)
		group, err := client.GetGroup(groupID)
		if err != nil {
			return fmt.Errorf("failed to get group: %w", err)
		}

		fmt.Printf("ID: %d\n", group.ID)
		fmt.Printf("Name: %s\n", group.Name)
		fmt.Printf("Description: %s\n", group.Description)
		fmt.Printf("Workspace ID: %d\n", group.WorkspaceID)
		fmt.Printf("Members: %d\n", len(group.UserIDs))
		if len(group.UserIDs) > 0 {
			fmt.Printf("User IDs: %v\n", group.UserIDs)
		}

		return nil
	},
}

var groupsCreateCmd = &cobra.Command{
	Use:   "create [workspace-id] [name]",
	Short: "Create a new group",
	Long:  `Create a new group in a workspace. Use flags to set optional properties.`,
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
		if groupDescriptionFlag != "" {
			opts["description"] = groupDescriptionFlag
		}
		if groupUserIDsFlag != "" {
			userIDs := []int{}
			for _, idStr := range strings.Split(groupUserIDsFlag, ",") {
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
		group, err := client.CreateGroup(workspaceID, name, opts)
		if err != nil {
			return fmt.Errorf("failed to create group: %w", err)
		}

		fmt.Printf("Group created successfully!\n")
		fmt.Printf("Group ID: %d\n", group.ID)
		fmt.Printf("Name: %s\n", group.Name)

		return nil
	},
}

var groupsUpdateCmd = &cobra.Command{
	Use:   "update [group-id]",
	Short: "Update a group",
	Long:  `Update group properties. Use flags to specify what to update.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		groupID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid group ID: %w", err)
		}

		token, err := auth.GetToken(tokenFlag)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		updates := make(map[string]interface{})
		if groupNameFlag != "" {
			updates["name"] = groupNameFlag
		}
		if groupDescriptionFlag != "" {
			updates["description"] = groupDescriptionFlag
		}

		if len(updates) == 0 {
			return fmt.Errorf("no updates specified; use flags like --name or --description")
		}

		client := api.NewClient(token)
		group, err := client.UpdateGroup(groupID, updates)
		if err != nil {
			return fmt.Errorf("failed to update group: %w", err)
		}

		fmt.Printf("Group updated successfully!\n")
		fmt.Printf("Group ID: %d\n", group.ID)
		fmt.Printf("Name: %s\n", group.Name)

		return nil
	},
}

var groupsDeleteCmd = &cobra.Command{
	Use:   "delete [group-id]",
	Short: "Delete a group",
	Long:  `Delete a group permanently.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		groupID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid group ID: %w", err)
		}

		token, err := auth.GetToken(tokenFlag)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		client := api.NewClient(token)
		if err := client.DeleteGroup(groupID); err != nil {
			return fmt.Errorf("failed to delete group: %w", err)
		}

		fmt.Printf("Group %d deleted successfully\n", groupID)
		return nil
	},
}

var groupsAddUserCmd = &cobra.Command{
	Use:   "add-user [group-id] [user-id]",
	Short: "Add a user to a group",
	Long:  `Add a user to a group by user ID.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		groupID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid group ID: %w", err)
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
		if err := client.AddGroupUser(groupID, userID); err != nil {
			return fmt.Errorf("failed to add user to group: %w", err)
		}

		fmt.Printf("User %d added to group %d successfully\n", userID, groupID)
		return nil
	},
}

var groupsRemoveUserCmd = &cobra.Command{
	Use:   "remove-user [group-id] [user-id]",
	Short: "Remove a user from a group",
	Long:  `Remove a user from a group by user ID.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		groupID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid group ID: %w", err)
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
		if err := client.RemoveGroupUser(groupID, userID); err != nil {
			return fmt.Errorf("failed to remove user from group: %w", err)
		}

		fmt.Printf("User %d removed from group %d successfully\n", userID, groupID)
		return nil
	},
}

func init() {
	groupsCreateCmd.Flags().StringVar(&groupDescriptionFlag, "description", "", "Group description")
	groupsCreateCmd.Flags().StringVar(&groupUserIDsFlag, "user-ids", "", "Comma-separated user IDs to add")

	groupsUpdateCmd.Flags().StringVar(&groupNameFlag, "name", "", "Group name")
	groupsUpdateCmd.Flags().StringVar(&groupDescriptionFlag, "description", "", "Group description")

	groupsCmd.AddCommand(groupsListCmd)
	groupsCmd.AddCommand(groupsShowCmd)
	groupsCmd.AddCommand(groupsCreateCmd)
	groupsCmd.AddCommand(groupsUpdateCmd)
	groupsCmd.AddCommand(groupsDeleteCmd)
	groupsCmd.AddCommand(groupsAddUserCmd)
	groupsCmd.AddCommand(groupsRemoveUserCmd)
}
