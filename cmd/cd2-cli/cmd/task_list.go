package cmd

import (
	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

func init() {
	taskCmd.AddCommand(taskListUploadsCmd)
	taskCmd.AddCommand(taskListDownloadsCmd)
	taskCmd.AddCommand(taskListCopyCmd)
	taskCmd.AddCommand(taskListMergeCmd)

	setCommandID(taskListCmd, "task.list")
	setCommandID(taskListUploadsCmd, "task.list.upload")
	setCommandID(taskListDownloadsCmd, "task.list.download")
	setCommandID(taskListCopyCmd, "task.list.copy")
	setCommandID(taskListMergeCmd, "task.list.merge")
}

var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()

		uploads, err := cd2Client.Transfer().GetUploadFileList(ctx, &pb.GetUploadFileListRequest{})
		if err != nil {
			return err
		}

		downloads, err := cd2Client.Transfer().GetDownloadFileList(ctx)
		if err != nil {
			return err
		}

		copyTasks, err := cd2Client.Copy().GetCopyTasks(ctx)
		if err != nil {
			return err
		}

		mergeTasks, err := cd2Client.Copy().GetMergeTasks(ctx)
		if err != nil {
			return err
		}

		result := map[string]interface{}{
			"uploads":    uploads,
			"downloads":  downloads,
			"copyTasks":  copyTasks,
			"mergeTasks": mergeTasks,
		}
		return outputResult(result)
	},
}

var taskListUploadsCmd = &cobra.Command{
	Use:   "list-uploads",
	Short: "List upload tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.Transfer().GetUploadFileList(ctx, &pb.GetUploadFileListRequest{})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var taskListDownloadsCmd = &cobra.Command{
	Use:   "list-downloads",
	Short: "List download tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.Transfer().GetDownloadFileList(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var taskListCopyCmd = &cobra.Command{
	Use:   "list-copy",
	Short: "List copy tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.Copy().GetCopyTasks(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var taskListMergeCmd = &cobra.Command{
	Use:   "list-merge",
	Short: "List merge tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.Copy().GetMergeTasks(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}
