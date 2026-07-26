package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	mountCmd.AddCommand(mountStatusCmd)
	setCommandID(mountStatusCmd, "mount.status")
}

var mountStatusCmd = &cobra.Command{
	Use:   "status [mount-point]",
	Short: "Show mount point status",
	Long:  "Show status of mount points. Without arguments, shows all mount points with their status. With mount-point argument, shows detailed status of that mount point.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.Mount().GetMountPoints(ctx)
		if err != nil {
			return err
		}

		if len(args) == 0 {
			return outputResult(result)
		} else {
			mountPoint := args[0]
			for _, mp := range result.MountPoints {
				if mp.MountPoint == mountPoint {
					return outputResult(mp)
				}
			}
			return outputResult(map[string]string{"error": "mount point not found"})
		}
	},
}
