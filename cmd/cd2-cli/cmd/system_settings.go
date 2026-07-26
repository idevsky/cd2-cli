package cmd

import (
	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

func init() {
	systemCmd.AddCommand(settingsCmd)

	settingsCmd.AddCommand(settingsGetCmd)
	settingsCmd.AddCommand(settingsSetCmd)

	setCommandID(settingsGetCmd, "system.settings-get")
	setCommandID(settingsSetCmd, "system.settings-set")
}

var settingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "System settings commands",
}

var settingsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get system settings",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		settings, err := cd2Client.System().GetSystemSettings(ctx)
		if err != nil {
			return err
		}
		return outputResult(settings)
	},
}

var settingsSetCmd = &cobra.Command{
	Use:   "set <json>",
	Short: "Set system settings from JSON",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		var settings pb.SystemSettings
		if err := parseProtoJSON([]byte(args[0]), &settings); err != nil {
			return err
		}
		err := cd2Client.System().SetSystemSettings(ctx, &settings)
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "updated"})
	},
}
