package cmd

import (
	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

func init() {
	systemCmd.AddCommand(openFilesCmd)
	systemCmd.AddCommand(machineIdCmd)
	systemCmd.AddCommand(devicesCmd)
	systemCmd.AddCommand(kickoutDeviceCmd)
	systemCmd.AddCommand(logsCmd)
	systemCmd.AddCommand(tempFilesCmd)

	setCommandID(openFilesCmd, "system.open-files")
	setCommandID(machineIdCmd, "system.machine-id")
	setCommandID(devicesCmd, "system.devices")
	setCommandID(kickoutDeviceCmd, "system.kickout-device")
	setCommandID(logsCmd, "system.logs")
	setCommandID(tempFilesCmd, "system.temp-files")
}

var openFilesCmd = &cobra.Command{
	Use:   "open-files",
	Short: "Get open file handles",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.System().GetOpenFileHandles(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var machineIdCmd = &cobra.Command{
	Use:   "machine-id",
	Short: "Get machine ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.System().GetMachineId(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Get online devices",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.System().GetOnlineDevices(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var kickoutDeviceCmd = &cobra.Command{
	Use:   "kickout-device <device-id>",
	Short: "Kickout a device",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		deviceId := args[0]
		err := cd2Client.System().KickoutDevice(ctx, &pb.DeviceRequest{
			DeviceId: deviceId,
		})
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "kicked"})
	},
}

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "List log files",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.System().ListLogFiles(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var tempFilesCmd = &cobra.Command{
	Use:   "temp-files",
	Short: "Get temp file table",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.System().GetTempFileTable(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}
