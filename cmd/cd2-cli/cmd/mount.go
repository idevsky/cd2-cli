package cmd

import (
	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

var mountCmd = &cobra.Command{
	Use:   "mount",
	Short: "Mount point management",
}

func init() {
	rootCmd.AddCommand(mountCmd)

	mountCmd.AddCommand(listMountsCmd)
	mountCmd.AddCommand(addMountCmd)
	mountCmd.AddCommand(updateMountCmd)
	mountCmd.AddCommand(removeMountCmd)

	updateMountCmd.Flags().String("options", "", "JSON string for MountOption")
	updateMountCmd.MarkFlagRequired("options")

	setCommandID(listMountsCmd, "mount.list")
	setCommandID(addMountCmd, "mount.add")
	setCommandID(updateMountCmd, "mount.update")
	setCommandID(removeMountCmd, "mount.remove")
}

var listMountsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all mount points",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.Mount().GetMountPoints(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var addMountCmd = &cobra.Command{
	Use:   "add [source-path] [mount-point]",
	Short: "Add a mount point",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		sourcePath := args[0]
		mountPoint := args[1]

		result, err := cd2Client.Mount().AddMountPoint(ctx, &pb.MountOption{
			SourceDir:  sourcePath,
			MountPoint: mountPoint,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var updateMountCmd = &cobra.Command{
	Use:   "update <mount-point>",
	Short: "Update a mount point (accept JSON for MountOption)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		mountPoint := args[0]
		jsonStr, _ := cmd.Flags().GetString("options")

		var option pb.MountOption
		if err := parseProtoJSON([]byte(jsonStr), &option); err != nil {
			return err
		}

		result, err := cd2Client.Mount().UpdateMountPoint(ctx, &pb.UpdateMountPointRequest{
			MountPoint:     mountPoint,
			NewMountOption: &option,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var removeMountCmd = &cobra.Command{
	Use:   "remove [mount-point]",
	Short: "Remove a mount point",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		mountPoint := args[0]

		result, err := cd2Client.Mount().RemoveMountPoint(ctx, &pb.MountPointRequest{
			MountPoint: mountPoint,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}
