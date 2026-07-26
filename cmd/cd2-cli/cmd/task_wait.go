package cmd

import (
	"context"
	"fmt"
	"time"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

var waitTimeout int
var missingIsComplete bool

func init() {
	taskCmd.AddCommand(taskWaitCopyCmd)
	taskCmd.AddCommand(taskWaitMergeCmd)

	taskWaitCopyCmd.Flags().IntVar(&waitTimeout, "timeout", 0, "Timeout in seconds (0 = no timeout)")
	taskWaitCopyCmd.Flags().BoolVar(&missingIsComplete, "missing-is-complete", false, "Treat missing task as completed (backward compatibility)")
	taskWaitMergeCmd.Flags().IntVar(&waitTimeout, "timeout", 0, "Timeout in seconds (0 = no timeout)")
	taskWaitMergeCmd.Flags().BoolVar(&missingIsComplete, "missing-is-complete", false, "Treat missing task as completed (backward compatibility)")

	setCommandID(taskWaitCopyCmd, "task.wait.copy")
	setCommandID(taskWaitMergeCmd, "task.wait.merge")
}

var taskWaitCmd = &cobra.Command{
	Use:   "wait",
	Short: "Wait for task completion",
}

var taskWaitCopyCmd = &cobra.Command{
	Use:   "wait-copy [source-path] [dest-path]",
	Short: "Wait for copy task to complete",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		sourcePath := args[0]
		destPath := args[1]

		timeout := time.Duration(waitTimeout) * time.Second
		if timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}

		firstPoll := true
		for {
			result, err := cd2Client.Copy().GetCopyTasks(ctx)
			if err != nil {
				return err
			}

			var foundTask *pb.CopyTask
			for _, task := range result.CopyTasks {
				if task.SourcePath == sourcePath && task.DestPath == destPath {
					foundTask = task
					break
				}
			}

			if foundTask == nil {
				if firstPoll {
					if missingIsComplete {
						return outputResult(map[string]string{"status": "completed", "message": "task not found"})
					}
					if err := outputResult(map[string]interface{}{"status": "not_found"}); err != nil {
						return err
					}
					return exitCode(map[string]interface{}{"status": "not_found"})
				}
				return outputResult(map[string]string{"status": "completed", "message": "task completed"})
			}

			firstPoll = false

			if foundTask.Status == pb.CopyTask_Completed {
				return outputResult(map[string]string{"status": "completed"})
			}

			if foundTask.Status == pb.CopyTask_Failed {
				if err := outputResult(map[string]interface{}{"status": "failed", "error": fmt.Sprintf("%v", foundTask.Errors)}); err != nil {
					return err
				}
				return exitCode(map[string]interface{}{"status": "failed", "error": fmt.Sprintf("%v", foundTask.Errors)})
			}

			select {
			case <-ctx.Done():
				if err := outputResult(map[string]interface{}{"status": "timeout"}); err != nil {
					return err
				}
				return exitCode(map[string]interface{}{"status": "timeout"})
			case <-time.After(2 * time.Second):
			}
		}
	},
}

var taskWaitMergeCmd = &cobra.Command{
	Use:   "wait-merge [source-path] [dest-path]",
	Short: "Wait for merge task to complete",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		sourcePath := args[0]
		destPath := args[1]

		timeout := time.Duration(waitTimeout) * time.Second
		if timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}

		firstPoll := true
		for {
			result, err := cd2Client.Copy().GetMergeTasks(ctx)
			if err != nil {
				return err
			}

			var foundTask *pb.MergeTask
			for _, task := range result.MergeTasks {
				if task.SourcePath == sourcePath && task.DestPath == destPath {
					foundTask = task
					break
				}
			}

			if foundTask == nil {
				if firstPoll {
					if missingIsComplete {
						return outputResult(map[string]string{"status": "completed", "message": "task not found"})
					}
					if err := outputResult(map[string]interface{}{"status": "not_found"}); err != nil {
						return err
					}
					return exitCode(map[string]interface{}{"status": "not_found"})
				}
				return outputResult(map[string]string{"status": "completed", "message": "task completed"})
			}

			firstPoll = false

			if foundTask.Status == pb.MergeTask_Completed {
				return outputResult(map[string]string{"status": "completed"})
			}

			if foundTask.Status == pb.MergeTask_Failed {
				errMsg := ""
				if foundTask.ErrorMessage != nil {
					errMsg = *foundTask.ErrorMessage
				}
				if err := outputResult(map[string]interface{}{"status": "failed", "error": errMsg}); err != nil {
					return err
				}
				return exitCode(map[string]interface{}{"status": "failed", "error": errMsg})
			}

			if foundTask.Status == pb.MergeTask_Cancelled {
				if err := outputResult(map[string]interface{}{"status": "cancelled"}); err != nil {
					return err
				}
				return exitCode(map[string]interface{}{"status": "cancelled"})
			}

			select {
			case <-ctx.Done():
				if err := outputResult(map[string]interface{}{"status": "timeout"}); err != nil {
					return err
				}
				return exitCode(map[string]interface{}{"status": "timeout"})
			case <-time.After(2 * time.Second):
			}
		}
	},
}
