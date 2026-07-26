package cmd

import (
	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Token management",
}

var listTokensCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tokens",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.Token().ListTokens(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var createTokenCmd = &cobra.Command{
	Use:   "create [root-dir] [name]",
	Short: "Create a new token",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		rootDir := args[0]
		name := args[1]

		result, err := cd2Client.Token().CreateToken(ctx, &pb.CreateTokenRequest{
			RootDir:      rootDir,
			FriendlyName: name,
			Permissions:  &pb.TokenPermissions{},
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var modifyTokenCmd = &cobra.Command{
	Use:   "modify <token>",
	Short: "Modify a token",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		tokenId := args[0]
		name, _ := cmd.Flags().GetString("name")
		rootDir, _ := cmd.Flags().GetString("root-dir")
		permissionsJson, _ := cmd.Flags().GetString("permissions")

		req := &pb.ModifyTokenRequest{
			Token: tokenId,
		}
		if name != "" {
			req.FriendlyName = &name
		}
		if rootDir != "" {
			req.RootDir = &rootDir
		}
		if permissionsJson != "" {
			var permissions pb.TokenPermissions
			if err := parseProtoJSON([]byte(permissionsJson), &permissions); err != nil {
				return err
			}
			req.Permissions = &permissions
		}

		result, err := cd2Client.Token().ModifyToken(ctx, req)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var removeTokenCmd = &cobra.Command{
	Use:   "remove [token]",
	Short: "Remove a token",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		token := args[0]

		err := cd2Client.Token().RemoveToken(ctx, token)
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "removed"})
	},
}

var infoTokenCmd = &cobra.Command{
	Use:   "info <token-id>",
	Short: "Get token info",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		tokenId := args[0]

		result, err := cd2Client.Token().GetTokenInfo(ctx, tokenId)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

func init() {
	rootCmd.AddCommand(tokenCmd)

	tokenCmd.AddCommand(listTokensCmd)
	tokenCmd.AddCommand(createTokenCmd)
	tokenCmd.AddCommand(modifyTokenCmd)
	tokenCmd.AddCommand(removeTokenCmd)
	tokenCmd.AddCommand(infoTokenCmd)

	modifyTokenCmd.Flags().String("name", "", "Friendly name")
	modifyTokenCmd.Flags().String("root-dir", "", "Root directory")
	modifyTokenCmd.Flags().String("permissions", "", "JSON string for TokenPermissions")

	setCommandID(listTokensCmd, "token.list")
	setCommandID(createTokenCmd, "token.create")
	setCommandID(modifyTokenCmd, "token.modify")
	setCommandID(removeTokenCmd, "token.remove")
	setCommandID(infoTokenCmd, "token.info")
}
