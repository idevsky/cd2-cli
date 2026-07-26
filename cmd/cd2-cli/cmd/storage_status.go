package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	storageCmd.AddCommand(storageStatusCmd)
	storageCmd.AddCommand(storageCanAddCmd)

	setCommandID(storageStatusCmd, "storage.status")
	setCommandID(storageCanAddCmd, "storage.can-add")
}

var storageStatusCmd = &cobra.Command{
	Use:   "status [cloud-name] [user-name]",
	Short: "Show storage provider status",
	Long:  "Show status of storage providers. Without arguments, shows overall storage status.",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()

		if len(args) == 0 {
			result, err := cd2Client.CloudAPI().GetAllCloudApis(ctx)
			if err != nil {
				return err
			}
			return outputResult(result)
		} else if len(args) == 2 {
			cloudName := args[0]
			userName := args[1]
			result, err := cd2Client.CloudAPI().GetAllCloudApis(ctx)
			if err != nil {
				return err
			}

			for _, api := range result.Apis {
				if api.Name == cloudName && api.UserName == userName {
					return outputResult(api)
				}
			}
			return outputResult(map[string]string{"error": "storage provider not found"})
		} else {
			return outputResult(map[string]string{"error": "please provide both cloud-name and user-name, or no arguments"})
		}
	},
}

var storageCanAddCmd = &cobra.Command{
	Use:   "can-add",
	Short: "Check if more storage providers can be added",
	Long:  "Check if the system can add more cloud storage providers",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.CloudAPI().CanAddMoreCloudApis(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}
