package cmd

import (
	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

var cloudapiCmd = &cobra.Command{
	Use:   "cloudapi",
	Short: "Cloud API management",
}

func init() {
	rootCmd.AddCommand(cloudapiCmd)

	cloudapiCmd.AddCommand(listCloudApisCmd)
	cloudapiCmd.AddCommand(canAddCloudApiCmd)
	cloudapiCmd.AddCommand(removeCloudApiCmd)
	cloudapiCmd.AddCommand(getCloudApiConfigCmd)
	cloudapiCmd.AddCommand(setCloudApiConfigCmd)

	setCloudApiConfigCmd.Flags().String("config", "", "JSON string for CloudAPIConfig")
	setCloudApiConfigCmd.MarkFlagRequired("config")

	setCommandID(listCloudApisCmd, "cloudapi.list")
	setCommandID(canAddCloudApiCmd, "cloudapi.can-add")
	setCommandID(removeCloudApiCmd, "cloudapi.remove")
	setCommandID(getCloudApiConfigCmd, "cloudapi.config")
	setCommandID(setCloudApiConfigCmd, "cloudapi.set-config")
}

var listCloudApisCmd = &cobra.Command{
	Use:   "list",
	Short: "List all cloud APIs",
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

var canAddCloudApiCmd = &cobra.Command{
	Use:   "can-add",
	Short: "Check if can add more cloud APIs",
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

var removeCloudApiCmd = &cobra.Command{
	Use:   "remove [cloud-name] [user-name]",
	Short: "Remove a cloud API",
	Args:  cobra.ExactArgs(2),
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

var getCloudApiConfigCmd = &cobra.Command{
	Use:   "config [cloud-name] [user-name]",
	Short: "Get cloud API configuration",
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

var setCloudApiConfigCmd = &cobra.Command{
	Use:   "set-config <cloud-name> <user-name>",
	Short: "Set cloud API configuration (accept JSON)",
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
