package cmd

import (
	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

var webdavCmd = &cobra.Command{
	Use:   "webdav",
	Short: "WebDAV server management",
}

var webdavUserCmd = &cobra.Command{
	Use:   "user",
	Short: "WebDAV user management",
}

var webdavServerCmd = &cobra.Command{
	Use:   "server",
	Short: "WebDAV server configuration",
}

var webdavUserGetCmd = &cobra.Command{
	Use:   "get <username>",
	Short: "Get WebDAV user info",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		username := args[0]

		result, err := cd2Client.WebDAV().GetDavUser(ctx, username)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var webdavUserAddCmd = &cobra.Command{
	Use:   "add [username] [password]",
	Short: "Add a WebDAV user",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		username := args[0]
		password := args[1]

		err := cd2Client.WebDAV().AddDavUser(ctx, &pb.AddDavUserRequest{
			UserName: username,
			Password: password,
		})
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "added"})
	},
}

var webdavUserModifyCmd = &cobra.Command{
	Use:   "modify <username>",
	Short: "Modify a WebDAV user",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		username := args[0]
		password, _ := cmd.Flags().GetString("password")
		rootPath, _ := cmd.Flags().GetString("root-path")

		req := &pb.ModifyDavUserRequest{
			UserName: username,
		}
		if password != "" {
			req.Password = &password
		}
		if rootPath != "" {
			req.RootPath = &rootPath
		}
		if cmd.Flags().Changed("read-only") {
			readOnly, _ := cmd.Flags().GetBool("read-only")
			req.ReadOnly = &readOnly
		}
		if cmd.Flags().Changed("enabled") {
			enabled, _ := cmd.Flags().GetBool("enabled")
			req.Enabled = &enabled
		}
		if cmd.Flags().Changed("guest") {
			guest, _ := cmd.Flags().GetBool("guest")
			req.Guest = &guest
		}

		err := cd2Client.WebDAV().ModifyDavUser(ctx, req)
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "modified"})
	},
}

var webdavUserRemoveCmd = &cobra.Command{
	Use:   "remove [username]",
	Short: "Remove a WebDAV user",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		username := args[0]

		err := cd2Client.WebDAV().RemoveDavUser(ctx, username)
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "removed"})
	},
}

var webdavServerGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get WebDAV server config",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.WebDAV().GetDavServerConfig(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var webdavServerSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set WebDAV server config (accept JSON)",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		jsonStr, _ := cmd.Flags().GetString("config")

		var config pb.ModifyDavServerConfigRequest
		if err := parseProtoJSON([]byte(jsonStr), &config); err != nil {
			return err
		}

		err := cd2Client.WebDAV().SetDavServerConfig(ctx, &config)
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "updated"})
	},
}

func init() {
	rootCmd.AddCommand(webdavCmd)

	webdavCmd.AddCommand(webdavUserCmd)
	webdavCmd.AddCommand(webdavServerCmd)

	webdavUserCmd.AddCommand(webdavUserGetCmd)
	webdavUserCmd.AddCommand(webdavUserAddCmd)
	webdavUserCmd.AddCommand(webdavUserModifyCmd)
	webdavUserCmd.AddCommand(webdavUserRemoveCmd)

	webdavServerCmd.AddCommand(webdavServerGetCmd)
	webdavServerCmd.AddCommand(webdavServerSetCmd)

	webdavUserModifyCmd.Flags().String("password", "", "New password")
	webdavUserModifyCmd.Flags().String("root-path", "", "Root path")
	webdavUserModifyCmd.Flags().Bool("read-only", false, "Read only")
	webdavUserModifyCmd.Flags().Bool("enabled", true, "Enabled")
	webdavUserModifyCmd.Flags().Bool("guest", false, "Guest")

	webdavServerSetCmd.Flags().String("config", "", "JSON string for ModifyDavServerConfigRequest")
	webdavServerSetCmd.MarkFlagRequired("config")

	setCommandID(webdavUserGetCmd, "webdav.user.get")
	setCommandID(webdavUserAddCmd, "webdav.user.add")
	setCommandID(webdavUserModifyCmd, "webdav.user.modify")
	setCommandID(webdavUserRemoveCmd, "webdav.user.remove")
	setCommandID(webdavServerGetCmd, "webdav.server.get")
	setCommandID(webdavServerSetCmd, "webdav.server.set")
}
