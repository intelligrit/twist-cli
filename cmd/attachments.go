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

var attachmentsCmd = &cobra.Command{
	Use:   "attachments",
	Short: "Manage attachments",
	Long:  `Upload, download, and view attachments on threads, comments, and conversations.`,
}

var attachmentsUploadCmd = &cobra.Command{
	Use:   "upload [target-type] [target-id] [file-path]",
	Short: "Upload a file attachment",
	Long:  `Upload a file to a thread, comment, or conversation. Target type must be 'thread', 'comment', or 'conversation'.`,
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetType := args[0]
		targetID, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid target ID: %w", err)
		}
		filePath := args[2]

		if targetType != "thread" && targetType != "comment" && targetType != "conversation" {
			return fmt.Errorf("invalid target type: must be 'thread', 'comment', or 'conversation'")
		}

		token, err := auth.GetToken(tokenFlag)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		client := api.NewClient(token)
		attachment, err := client.UploadAttachment(targetType, targetID, filePath)
		if err != nil {
			return fmt.Errorf("failed to upload attachment: %w", err)
		}

		fmt.Printf("Attachment uploaded successfully!\n")
		fmt.Printf("Attachment ID: %d\n", attachment.ID)
		fmt.Printf("Title: %s\n", attachment.Title)
		fmt.Printf("Size: %d bytes\n", attachment.Size)

		return nil
	},
}

var attachmentsDownloadCmd = &cobra.Command{
	Use:   "download [attachment-id] [output-path]",
	Short: "Download an attachment",
	Long:  `Download an attachment by ID to a local file.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		attachmentID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid attachment ID: %w", err)
		}
		outputPath := args[1]

		token, err := auth.GetToken(tokenFlag)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		client := api.NewClient(token)
		if err := client.DownloadAttachment(attachmentID, outputPath); err != nil {
			return fmt.Errorf("failed to download attachment: %w", err)
		}

		fmt.Printf("Attachment downloaded successfully to %s\n", outputPath)
		return nil
	},
}

var attachmentsListCmd = &cobra.Command{
	Use:   "list [target-type] [target-id]",
	Short: "List all attachments",
	Long:  `List all attachments on a thread, comment, or conversation. Target type must be 'thread', 'comment', or 'conversation'.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetType := args[0]
		targetID, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid target ID: %w", err)
		}

		if targetType != "thread" && targetType != "comment" && targetType != "conversation" {
			return fmt.Errorf("invalid target type: must be 'thread', 'comment', or 'conversation'")
		}

		token, err := auth.GetToken(tokenFlag)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		client := api.NewClient(token)
		attachments, err := client.GetAttachments(targetType, targetID)
		if err != nil {
			return fmt.Errorf("failed to get attachments: %w", err)
		}

		if len(attachments) == 0 {
			fmt.Println("No attachments found.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tTITLE\tSIZE\tTYPE")
		fmt.Fprintln(w, "--\t-----\t----\t----")
		for _, a := range attachments {
			sizeStr := fmt.Sprintf("%d bytes", a.Size)
			if a.Size > 1024*1024 {
				sizeStr = fmt.Sprintf("%.2f MB", float64(a.Size)/(1024*1024))
			} else if a.Size > 1024 {
				sizeStr = fmt.Sprintf("%.2f KB", float64(a.Size)/1024)
			}
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", a.ID, a.Title, sizeStr, a.MimeType)
		}
		w.Flush()

		return nil
	},
}

func init() {
	attachmentsCmd.AddCommand(attachmentsUploadCmd)
	attachmentsCmd.AddCommand(attachmentsDownloadCmd)
	attachmentsCmd.AddCommand(attachmentsListCmd)
}
