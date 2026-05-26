package cmd

import (
	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "File and folder operations",
}

func init() {
	rootCmd.AddCommand(fileCmd)
	fileCmd.AddCommand(listFilesCmd)
	fileCmd.AddCommand(findFileCmd)

	listFilesCmd.Flags().Bool("refresh", false, "Force refresh cache")

	setCommandID(listFilesCmd, "file.list")
	setCommandID(findFileCmd, "file.find")
}

var listFilesCmd = &cobra.Command{
	Use:   "list [path]",
	Short: "List files in a directory",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		path := args[0]
		forceRefresh, _ := cmd.Flags().GetBool("refresh")

		files, err := cd2Client.File().GetSubFiles(ctx, &pb.ListSubFileRequest{
			Path:         path,
			ForceRefresh: forceRefresh,
		})
		if err != nil {
			return err
		}
		return outputResult(files)
	},
}

var findFileCmd = &cobra.Command{
	Use:   "find [path]",
	Short: "Find file by path",
	Args:  cobra.ExactArgs(1),
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
