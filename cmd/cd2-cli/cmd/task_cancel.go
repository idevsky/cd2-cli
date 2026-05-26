package cmd

import (
	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

func init() {
	taskCmd.AddCommand(taskCancelCmd)
	taskCancelCmd.AddCommand(taskCancelUploadCmd)
	taskCancelCmd.AddCommand(taskCancelCopyCmd)
	taskCancelCmd.AddCommand(taskCancelMergeCmd)

	setCommandID(taskCancelUploadCmd, "task.cancel.upload")
	setCommandID(taskCancelCopyCmd, "task.cancel.copy")
	setCommandID(taskCancelMergeCmd, "task.cancel.merge")
}

var taskCancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel tasks",
}

var taskCancelUploadCmd = &cobra.Command{
	Use:   "cancel-upload [key...]",
	Short: "Cancel upload tasks",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		err := cd2Client.Transfer().CancelUploadFiles(ctx, args)
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "cancelled"})
	},
}

var taskCancelCopyCmd = &cobra.Command{
	Use:   "cancel-copy [source-path] [dest-path]",
	Short: "Cancel a copy task",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		sourcePath := args[0]
		destPath := args[1]

		err := cd2Client.Copy().CancelCopyTask(ctx, &pb.CopyTaskRequest{
			SourcePath: sourcePath,
			DestPath:   destPath,
		})
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "cancelled"})
	},
}

var taskCancelMergeCmd = &cobra.Command{
	Use:   "cancel-merge [source-path] [dest-path]",
	Short: "Cancel a merge task",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		sourcePath := args[0]
		destPath := args[1]

		err := cd2Client.Copy().CancelMergeTask(ctx, &pb.CancelMergeTaskRequest{
			SourcePath: sourcePath,
			DestPath:   destPath,
		})
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "cancelled"})
	},
}
