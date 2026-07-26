package cmd

import (
	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

func init() {
	cloudapiCmd.AddCommand(loginWebdavCmd)
	cloudapiCmd.AddCommand(loginS3Cmd)
	cloudapiCmd.AddCommand(loginLocalCmd)
	cloudapiCmd.AddCommand(loginClouddriveCmd)
	cloudapiCmd.AddCommand(loginSftpCmd)
	cloudapiCmd.AddCommand(loginFtpCmd)
	cloudapiCmd.AddCommand(loginSmbCmd)
	cloudapiCmd.AddCommand(discoverSmbServersCmd)
	cloudapiCmd.AddCommand(discoverSmbSharesCmd)

	setCommandID(loginWebdavCmd, "cloudapi.login-webdav")
	setCommandID(loginS3Cmd, "cloudapi.login-s3")
	setCommandID(loginLocalCmd, "cloudapi.login-local")
	setCommandID(loginClouddriveCmd, "cloudapi.login-clouddrive")
	setCommandID(loginSftpCmd, "cloudapi.login-sftp")
	setCommandID(loginFtpCmd, "cloudapi.login-ftp")
	setCommandID(loginSmbCmd, "cloudapi.login-smb")
	setCommandID(discoverSmbServersCmd, "cloudapi.discover-smb-servers")
	setCommandID(discoverSmbSharesCmd, "cloudapi.discover-smb-shares")
}

var loginWebdavCmd = &cobra.Command{
	Use:   "login-webdav <json>",
	Short: "Login to WebDAV (accept JSON)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		var req pb.LoginWebDavRequest
		if err := parseProtoJSON([]byte(args[0]), &req); err != nil {
			return err
		}

		result, err := cd2Client.CloudAPI().APILoginWebDav(ctx, &req)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var loginS3Cmd = &cobra.Command{
	Use:   "login-s3 <json>",
	Short: "Login to S3 (accept JSON)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		var req pb.LoginS3Request
		if err := parseProtoJSON([]byte(args[0]), &req); err != nil {
			return err
		}

		result, err := cd2Client.CloudAPI().APILoginS3(ctx, &req)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var loginLocalCmd = &cobra.Command{
	Use:   "login-local <path>",
	Short: "Add local folder",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		path := args[0]

		result, err := cd2Client.CloudAPI().APIAddLocalFolder(ctx, &pb.AddLocalFolderRequest{
			LocalFolderPath: path,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var loginClouddriveCmd = &cobra.Command{
	Use:   "login-clouddrive <json>",
	Short: "Login to CloudDrive (accept JSON)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		var req pb.LoginCloudDriveRequest
		if err := parseProtoJSON([]byte(args[0]), &req); err != nil {
			return err
		}

		result, err := cd2Client.CloudAPI().APILoginCloudDrive(ctx, &req)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var loginSftpCmd = &cobra.Command{
	Use:   "login-sftp <json>",
	Short: "Login to SFTP (accept JSON)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		var req pb.LoginSftpRequest
		if err := parseProtoJSON([]byte(args[0]), &req); err != nil {
			return err
		}

		result, err := cd2Client.CloudAPI().APILoginSftp(ctx, &req)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var loginFtpCmd = &cobra.Command{
	Use:   "login-ftp <json>",
	Short: "Login to FTP (accept JSON)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		var req pb.LoginFtpRequest
		if err := parseProtoJSON([]byte(args[0]), &req); err != nil {
			return err
		}

		result, err := cd2Client.CloudAPI().APILoginFtp(ctx, &req)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var loginSmbCmd = &cobra.Command{
	Use:   "login-smb <json>",
	Short: "Login to SMB (accept JSON)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		var req pb.LoginSmbRequest
		if err := parseProtoJSON([]byte(args[0]), &req); err != nil {
			return err
		}

		result, err := cd2Client.CloudAPI().APILoginSmb(ctx, &req)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var discoverSmbServersCmd = &cobra.Command{
	Use:   "discover-smb-servers",
	Short: "Discover SMB servers",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.CloudAPI().DiscoverSmbServers(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var discoverSmbSharesCmd = &cobra.Command{
	Use:   "discover-smb-shares <server>",
	Short: "Discover SMB shares on a server",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		server := args[0]

		result, err := cd2Client.CloudAPI().DiscoverSmbShares(ctx, &pb.DiscoverSmbSharesRequest{
			Server: server,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}
