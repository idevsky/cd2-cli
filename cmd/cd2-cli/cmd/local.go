package cmd

import (
	"os"
	"path/filepath"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Local filesystem operations",
}

var localListCmd = &cobra.Command{
	Use:   "list <parent-folder>",
	Short: "List local sub files",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		parentFolder := args[0]
		folderOnly, _ := cmd.Flags().GetBool("folder-only")

		entries, err := os.ReadDir(parentFolder)
		if err != nil {
			return err
		}

		var subFiles []string
		for _, entry := range entries {
			if folderOnly && !entry.IsDir() {
				continue
			}
			subFiles = append(subFiles, filepath.Join(parentFolder, entry.Name()))
		}

		result := &pb.LocalGetSubFilesResult{
			SubFiles: subFiles,
		}
		return outputResult(result)
	},
}

var localMkdirCmd = &cobra.Command{
	Use:   "mkdir <parent> <name>",
	Short: "Create local folder",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		parent := args[0]
		name := args[1]

		createdPath := filepath.Join(parent, name)
		err := os.MkdirAll(createdPath, 0755)
		if err != nil {
			return err
		}

		result := &pb.LocalCreateFolderResult{
			Success:     true,
			CreatedPath: createdPath,
		}
		return outputResult(result)
	},
}

func init() {
	rootCmd.AddCommand(localCmd)

	localCmd.AddCommand(localListCmd)
	localCmd.AddCommand(localMkdirCmd)

	localListCmd.Flags().Bool("folder-only", false, "List only folders")
	localListCmd.Flags().Bool("include-cd", false, "Include CloudDrive mount points (ignored for local filesystem)")
	localListCmd.Flags().Bool("include-available", false, "Include available drives (ignored for local filesystem)")

	setCommandID(localListCmd, "local.list")
	setCommandID(localMkdirCmd, "local.mkdir")

	markAsLocal(localListCmd)
	markAsLocal(localMkdirCmd)
}
