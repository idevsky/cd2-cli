package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	storageCmd.AddCommand(storageListCmd)
	setCommandID(storageListCmd, "storage.list")
}

var storageListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all storage providers",
	Long:    "List all configured cloud storage providers",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.CloudAPI().GetAllCloudApis(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}
