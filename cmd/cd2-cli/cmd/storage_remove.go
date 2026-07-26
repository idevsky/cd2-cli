package cmd

import (
	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

func init() {
	storageCmd.AddCommand(storageRemoveCmd)
	setCommandID(storageRemoveCmd, "storage.remove")
}

var storageRemoveCmd = &cobra.Command{
	Use:     "remove <cloud-name> <user-name>",
	Aliases: []string{"rm", "delete"},
	Short:   "Remove a storage provider",
	Long:    "Remove a storage provider by cloud name and user name",
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		cloudName := args[0]
		userName := args[1]

		result, err := cd2Client.CloudAPI().RemoveCloudAPI(ctx, &pb.RemoveCloudAPIRequest{
			CloudName: cloudName,
			UserName:  userName,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}
