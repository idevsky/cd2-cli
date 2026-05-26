package cmd

import (
	"github.com/spf13/cobra"
)

var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "System information and management commands",
}

func init() {
	rootCmd.AddCommand(systemCmd)

	systemCmd.AddCommand(systemInfoCmd)
	systemCmd.AddCommand(runtimeInfoCmd)
	systemCmd.AddCommand(runningInfoCmd)
	systemCmd.AddCommand(restartCmd)
	systemCmd.AddCommand(shutdownCmd)
	systemCmd.AddCommand(capabilitiesCmd)

	setCommandID(systemInfoCmd, "system.info")
	setCommandID(runtimeInfoCmd, "system.runtime")
	setCommandID(runningInfoCmd, "system.running")
	setCommandID(restartCmd, "system.restart")
	setCommandID(shutdownCmd, "system.shutdown")
	setCommandID(capabilitiesCmd, "system.capabilities")
}

var systemInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get system information",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		info, err := cd2Client.Public().GetSystemInfo(ctx)
		if err != nil {
			return err
		}
		return outputResult(info)
	},
}

var runtimeInfoCmd = &cobra.Command{
	Use:   "runtime",
	Short: "Get runtime information",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		info, err := cd2Client.System().GetRuntimeInfo(ctx)
		if err != nil {
			return err
		}
		return outputResult(info)
	},
}

var runningInfoCmd = &cobra.Command{
	Use:   "running",
	Short: "Get running information",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		info, err := cd2Client.System().GetRunningInfo(ctx)
		if err != nil {
			return err
		}
		return outputResult(info)
	},
}

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart the service",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		err := cd2Client.System().RestartService(ctx)
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "restarted"})
	},
}

var shutdownCmd = &cobra.Command{
	Use:   "shutdown",
	Short: "Shutdown the service",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		err := cd2Client.System().ShutdownService(ctx)
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "shutdown"})
	},
}

var capabilitiesCmd = &cobra.Command{
	Use:   "capabilities",
	Short: "Get service capabilities",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		caps, err := cd2Client.System().GetServiceCapabilities(ctx)
		if err != nil {
			return err
		}
		return outputResult(caps)
	},
}
