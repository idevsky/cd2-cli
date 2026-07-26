package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	setCommandID(taskStatusCmd, "task.status")
}

var taskStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get task status summary",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.Transfer().GetAllTasksCount(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}
