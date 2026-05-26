package cmd

import (
	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

func init() {
	storageCmd.AddCommand(storageConfigCmd)
	storageCmd.AddCommand(storageSetConfigCmd)

	storageSetConfigCmd.Flags().String("config", "", "JSON string for CloudAPIConfig")
	storageSetConfigCmd.MarkFlagRequired("config")

	setCommandID(storageConfigCmd, "storage.config")
	setCommandID(storageSetConfigCmd, "storage.set-config")
}

var storageConfigCmd = &cobra.Command{
	Use:   "config <cloud-name> <user-name>",
	Short: "Get storage provider configuration",
	Long:  "Get configuration for a storage provider by cloud name and user name",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		cloudName := args[0]
		userName := args[1]

		config, err := cd2Client.CloudAPI().GetCloudAPIConfig(ctx, &pb.GetCloudAPIConfigRequest{
			CloudName: cloudName,
			UserName:  userName,
		})
		if err != nil {
			return err
		}
		return outputResult(config)
	},
}

var storageSetConfigCmd = &cobra.Command{
	Use:   "set-config <cloud-name> <user-name>",
	Short: "Set storage provider configuration",
	Long:  "Set configuration for a storage provider (accepts JSON configuration)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		cloudName := args[0]
		userName := args[1]
		jsonStr, _ := cmd.Flags().GetString("config")

		var config pb.CloudAPIConfig
		if err := parseProtoJSON([]byte(jsonStr), &config); err != nil {
			return err
		}

		err := cd2Client.CloudAPI().SetCloudAPIConfig(ctx, &pb.SetCloudAPIConfigRequest{
			CloudName: cloudName,
			UserName:  userName,
			Config:    &config,
		})
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "updated"})
	},
}
