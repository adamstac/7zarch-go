package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func UploadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload <archive>",
		Short: "Upload archive to TrueNAS",
		Long:  `Upload an archive to TrueNAS storage via SSH/SFTP.`,
		Args:  cobra.ExactArgs(1),
		RunE:  runUpload,
	}

	// Add flags
	cmd.Flags().String("path", "", "Remote path on TrueNAS")
	cmd.Flags().String("storage", "truenas", "Storage backend to use")
	cmd.Flags().Bool("skip-existing", true, "Skip if file already exists")

	return cmd
}

func runUpload(cmd *cobra.Command, args []string) error {
	archivePath := args[0]
	remotePath, _ := cmd.Flags().GetString("path")

<<<<<<< HEAD
	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "Uploading %s to TrueNAS...\n", archivePath)
=======
	fmt.Printf("Uploading %s to TrueNAS...\n", archivePath)
>>>>>>> origin/main
	if remotePath != "" {
		fmt.Fprintf(out, "Remote path: %s\n", remotePath)
	}

	// TODO: Implement TrueNAS upload
<<<<<<< HEAD
	fmt.Fprintf(out, "\n⚠️  Upload functionality coming soon!\n")
	fmt.Fprintf(out, "This will upload to TrueNAS via SSH/SFTP.\n")
=======
	fmt.Printf("\n⚠️  Upload functionality coming soon!\n")
	fmt.Printf("This will upload to TrueNAS via SSH/SFTP.\n")
>>>>>>> origin/main

	return nil
}
