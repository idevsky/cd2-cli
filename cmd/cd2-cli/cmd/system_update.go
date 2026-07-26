package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	systemCmd.AddCommand(updateCmd)

	updateCmd.AddCommand(updateCheckCmd)
	updateCmd.AddCommand(updateDownloadCmd)
	updateCmd.AddCommand(updateApplyCmd)
	updateCmd.AddCommand(hasUpdateCmd)

	setCommandID(updateCheckCmd, "system.update-check")
	setCommandID(updateDownloadCmd, "system.update-download")
	setCommandID(updateApplyCmd, "system.update-apply")
	setCommandID(hasUpdateCmd, "system.has-update")
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "System update commands",
}

var updateCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for updates",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.System().CheckUpdate(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var updateDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download update",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		err := cd2Client.System().DownloadUpdate(ctx)
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "downloaded"})
	},
}

var updateApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply update",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		err := cd2Client.System().UpdateSystem(ctx)
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "updated"})
	},
}

var hasUpdateCmd = &cobra.Command{
	Use:   "has-update",
	Short: "Check if update is available",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.System().HasUpdate(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}
