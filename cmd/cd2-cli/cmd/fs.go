package cmd

import (
	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(lsCmd)
	rootCmd.AddCommand(statCmd)
	rootCmd.AddCommand(mkdirCmd)

	lsCmd.Flags().Bool("refresh", false, "Force refresh cache")

	setCommandID(lsCmd, "fs.ls")
	setCommandID(statCmd, "fs.stat")
	setCommandID(mkdirCmd, "fs.mkdir")
}

var lsCmd = &cobra.Command{
	Use:     "ls [path]",
	Short:   "List directory contents",
	Aliases: []string{"list"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		path := args[0]
		refresh, _ := cmd.Flags().GetBool("refresh")

		files, err := cd2Client.File().GetSubFiles(ctx, &pb.ListSubFileRequest{
			Path:         path,
			ForceRefresh: refresh,
		})
		if err != nil {
			return err
		}
		return outputResult(files)
	},
}

var statCmd = &cobra.Command{
	Use:     "stat [path]",
	Short:   "Get file/directory information",
	Aliases: []string{"info"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		path := args[0]

		file, err := cd2Client.File().FindFileByPath(ctx, path)
		if err != nil {
			return err
		}
		return outputResult(file)
	},
}

var mkdirCmd = &cobra.Command{
	Use:     "mkdir [parent-path] [folder-name]",
	Short:   "Create a directory",
	Aliases: []string{"makedir"},
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		parentPath := args[0]
		folderName := args[1]

		result, err := cd2Client.File().CreateFolder(ctx, &pb.CreateFolderRequest{
			ParentPath: parentPath,
			FolderName: folderName,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}
