package cmd

import (
	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy task management",
}

func init() {
	rootCmd.AddCommand(copyCmd)

	copyCmd.AddCommand(copyTasksCmd)
	copyCmd.AddCommand(copyMergeTasksCmd)
	copyCmd.AddCommand(copyCancelCmd)
	copyCmd.AddCommand(copyPauseCmd)
	copyCmd.AddCommand(copyRestartCmd)
	copyCmd.AddCommand(copyRemoveCmd)
	copyCmd.AddCommand(copyResumeCmd)

	setCommandID(copyTasksCmd, "copy.tasks")
	setCommandID(copyMergeTasksCmd, "copy.merge-tasks")
	setCommandID(copyCancelCmd, "copy.cancel")
	setCommandID(copyPauseCmd, "copy.pause")
	setCommandID(copyRestartCmd, "copy.restart")
	setCommandID(copyRemoveCmd, "copy.remove")
	setCommandID(copyResumeCmd, "copy.resume")
}

var copyTasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "List all copy tasks",
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

var copyMergeTasksCmd = &cobra.Command{
	Use:   "merge-tasks",
	Short: "List all merge tasks",
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

var copyCancelCmd = &cobra.Command{
	Use:   "cancel [source-path] [dest-path]",
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

var copyPauseCmd = &cobra.Command{
	Use:   "pause [source-path] [dest-path]",
	Short: "Pause a copy task",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		sourcePath := args[0]
		destPath := args[1]

		err := cd2Client.Copy().PauseCopyTask(ctx, &pb.PauseCopyTaskRequest{
			SourcePath: sourcePath,
			DestPath:   destPath,
			Pause:      true,
		})
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "paused"})
	},
}

var copyRestartCmd = &cobra.Command{
	Use:   "restart [source-path] [dest-path]",
	Short: "Restart a copy task",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		sourcePath := args[0]
		destPath := args[1]

		err := cd2Client.Copy().RestartCopyTask(ctx, &pb.CopyTaskRequest{
			SourcePath: sourcePath,
			DestPath:   destPath,
		})
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "restarted"})
	},
}

var copyRemoveCmd = &cobra.Command{
	Use:   "remove [task-keys...]",
	Short: "Remove copy tasks by their keys",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()

		result, err := cd2Client.Copy().RemoveCopyTasks(ctx, &pb.CopyTaskBatchRequest{
			TaskKeys: args,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var copyResumeCmd = &cobra.Command{
	Use:   "resume [source-path] [dest-path]",
	Short: "Resume a paused copy task",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		sourcePath := args[0]
		destPath := args[1]

		err := cd2Client.Copy().PauseCopyTask(ctx, &pb.PauseCopyTaskRequest{
			SourcePath: sourcePath,
			DestPath:   destPath,
			Pause:      false,
		})
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "resumed"})
	},
}
