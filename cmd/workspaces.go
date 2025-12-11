package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/robertmeta/twist-cli/internal/auth"
	"github.com/robertmeta/twist-cli/pkg/api"
	"github.com/spf13/cobra"
)

var workspacesCmd = &cobra.Command{
	Use:   "workspaces",
	Short: "Manage Twist workspaces",
	Long:  `View and manage your Twist workspaces.`,
}

var workspacesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all workspaces",
	Long:  `List all workspaces that you have access to.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := auth.EnsureToken()
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		client := api.NewClient(token)
		workspaces, err := client.GetWorkspaces()
		if err != nil {
			return fmt.Errorf("failed to get workspaces: %w", err)
		}

		if len(workspaces) == 0 {
			fmt.Println("No workspaces found.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tPLAN")
		fmt.Fprintln(w, "--\t----\t----")
		for _, ws := range workspaces {
			fmt.Fprintf(w, "%d\t%s\t%s\n", ws.ID, ws.Name, ws.Plan)
		}
		w.Flush()

		return nil
	},
}

func init() {
	workspacesCmd.AddCommand(workspacesListCmd)
}
