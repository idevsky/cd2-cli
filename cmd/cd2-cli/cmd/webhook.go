package cmd

import (
	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

var webhookCmd = &cobra.Command{
	Use:   "webhook",
	Short: "Webhook management",
}

var listWebhooksCmd = &cobra.Command{
	Use:   "list",
	Short: "List all webhooks",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.Webhook().GetWebhookConfigs(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var webhookTemplateCmd = &cobra.Command{
	Use:   "template",
	Short: "Get webhook config template",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.Webhook().GetWebhookConfigTemplate(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var addWebhookCmd = &cobra.Command{
	Use:   "add [filename] [content]",
	Short: "Add a webhook config",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		filename := args[0]
		content := args[1]

		err := cd2Client.Webhook().AddWebhookConfig(ctx, &pb.WebhookRequest{
			FileName: filename,
			Content:  content,
		})
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "added"})
	},
}

var changeWebhookCmd = &cobra.Command{
	Use:   "change <filename>",
	Short: "Change a webhook config (accept JSON content)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		filename := args[0]
		content, _ := cmd.Flags().GetString("content")

		err := cd2Client.Webhook().ChangeWebhookConfig(ctx, &pb.WebhookRequest{
			FileName: filename,
			Content:  content,
		})
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "changed"})
	},
}

var removeWebhookCmd = &cobra.Command{
	Use:   "remove [filename]",
	Short: "Remove a webhook config",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		filename := args[0]

		err := cd2Client.Webhook().RemoveWebhookConfig(ctx, filename)
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "removed"})
	},
}

func init() {
	rootCmd.AddCommand(webhookCmd)

	webhookCmd.AddCommand(listWebhooksCmd)
	webhookCmd.AddCommand(webhookTemplateCmd)
	webhookCmd.AddCommand(addWebhookCmd)
	webhookCmd.AddCommand(changeWebhookCmd)
	webhookCmd.AddCommand(removeWebhookCmd)

	changeWebhookCmd.Flags().String("content", "", "JSON content for webhook config")

	setCommandID(listWebhooksCmd, "webhook.list")
	setCommandID(webhookTemplateCmd, "webhook.template")
	setCommandID(addWebhookCmd, "webhook.add")
	setCommandID(changeWebhookCmd, "webhook.change")
	setCommandID(removeWebhookCmd, "webhook.remove")
}
