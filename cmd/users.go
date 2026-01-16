package cmd

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/intelligrit/twist-cli/internal/auth"
	"github.com/intelligrit/twist-cli/pkg/api"
	"github.com/spf13/cobra"
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Manage users",
	Long:  `View users in your Twist workspaces.`,
}

var usersListCmd = &cobra.Command{
	Use:   "list [workspace-id]",
	Short: "List all users in a workspace",
	Long:  `List all users in a specific workspace.`,
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
		users, err := client.GetWorkspaceUsers(workspaceID)
		if err != nil {
			return fmt.Errorf("failed to get users: %w", err)
		}

		if len(users) == 0 {
			fmt.Println("No users found in this workspace.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tEMAIL\tTYPE\tBOT\tREMOVED")
		fmt.Fprintln(w, "--\t----\t-----\t----\t---\t-------")
		for _, u := range users {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%t\t%t\n",
				u.ID, u.Name, u.Email, u.UserType, u.Bot, u.Removed)
		}
		w.Flush()

		return nil
	},
}

func init() {
	usersCmd.AddCommand(usersListCmd)
}
