package cmd

import (
	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

func init() {
	systemCmd.AddCommand(webserverCmd)

	webserverCmd.AddCommand(webserverGetCmd)
	webserverCmd.AddCommand(webserverSetCmd)
	webserverCmd.AddCommand(generateCertCmd)

	setCommandID(webserverGetCmd, "system.webserver-get")
	setCommandID(webserverSetCmd, "system.webserver-set")
	setCommandID(generateCertCmd, "system.generate-cert")
}

var webserverCmd = &cobra.Command{
	Use:   "webserver",
	Short: "WebServer configuration commands",
}

var webserverGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get WebServer configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		config, err := cd2Client.System().GetWebServerConfig(ctx)
		if err != nil {
			return err
		}
		return outputResult(config)
	},
}

var webserverSetCmd = &cobra.Command{
	Use:   "set <json>",
	Short: "Set WebServer configuration",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		var req pb.SetWebServerConfigRequest
		if err := parseProtoJSON([]byte(args[0]), &req); err != nil {
			return err
		}
		err := cd2Client.System().SetWebServerConfig(ctx, &req)
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "updated"})
	},
}

var generateCertCmd = &cobra.Command{
	Use:   "generate-cert",
	Short: "Generate self-signed certificate",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		var req pb.GenerateSelfSignedCertRequest
		err := cd2Client.System().GenerateSelfSignedCert(ctx, &req)
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "generated"})
	},
}
