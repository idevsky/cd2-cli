package cmd

import (
	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

func init() {
	mountCmd.AddCommand(mountMountCmd)
	mountCmd.AddCommand(unmountCmd)
	mountCmd.AddCommand(driveLettersCmd)
	mountCmd.AddCommand(canAddMountCmd)
	mountCmd.AddCommand(hasDriveLettersCmd)
	mountCmd.AddCommand(canMountBothCmd)

	setCommandID(mountMountCmd, "mount.start")
	setCommandID(unmountCmd, "mount.stop")
	setCommandID(driveLettersCmd, "mount.drive-letters")
	setCommandID(canAddMountCmd, "mount.can-add")
	setCommandID(hasDriveLettersCmd, "mount.has-drive-letters")
	setCommandID(canMountBothCmd, "mount.can-mount-both")
}

var mountMountCmd = &cobra.Command{
	Use:   "start [mount-point]",
	Short: "Mount a mount point",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		mountPoint := args[0]

		result, err := cd2Client.Mount().Mount(ctx, &pb.MountPointRequest{
			MountPoint: mountPoint,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var unmountCmd = &cobra.Command{
	Use:   "stop [mount-point]",
	Short: "Unmount a mount point",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		mountPoint := args[0]

		result, err := cd2Client.Mount().Unmount(ctx, &pb.MountPointRequest{
			MountPoint: mountPoint,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var driveLettersCmd = &cobra.Command{
	Use:   "drive-letters",
	Short: "Get available drive letters (Windows)",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.Mount().GetAvailableDriveLetters(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var canAddMountCmd = &cobra.Command{
	Use:   "can-add",
	Short: "Check if can add more mount points",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.Mount().CanAddMoreMountPoints(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var hasDriveLettersCmd = &cobra.Command{
	Use:   "has-drive-letters",
	Short: "Check if has drive letters",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.Mount().HasDriveLetters(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var canMountBothCmd = &cobra.Command{
	Use:   "can-mount-both",
	Short: "Check if can mount both local and cloud",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.Mount().CanMountBothLocalAndCloud(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}
