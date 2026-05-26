package cmd

import (
	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

var openFileTableCmd = &cobra.Command{
	Use:   "open-file-table",
	Short: "Get open file table",
	RunE: func(cmd *cobra.Command, args []string) error {
		includeDir, _ := cmd.Flags().GetBool("include-dir")
		ctx, cancel := getTimeoutContext()
		defer cancel()
		table, err := cd2Client.System().GetOpenFileTable(ctx, &pb.GetOpenFileTableRequest{
			IncludeDir: includeDir,
		})
		if err != nil {
			return err
		}
		return outputResult(table)
	},
}

func init() {
	openFileTableCmd.Flags().Bool("include-dir", false, "Include directory entries")
	systemCmd.AddCommand(openFileTableCmd)
	setCommandID(openFileTableCmd, "system.open-file-table")
}
