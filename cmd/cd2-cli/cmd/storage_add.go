package cmd

import (
	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

func init() {
	storageCmd.AddCommand(storageAddCmd)
	storageAddCmd.AddCommand(addS3Cmd)
	storageAddCmd.AddCommand(addWebdavCmd)
	storageAddCmd.AddCommand(addLocalCmd)
	storageAddCmd.AddCommand(addSftpCmd)
	storageAddCmd.AddCommand(addFtpCmd)
	storageAddCmd.AddCommand(addSmbCmd)

	setCommandID(storageAddCmd, "storage.add")
	setCommandID(addS3Cmd, "storage.add.s3")
	setCommandID(addWebdavCmd, "storage.add.webdav")
	setCommandID(addLocalCmd, "storage.add.local")
	setCommandID(addSftpCmd, "storage.add.sftp")
	setCommandID(addFtpCmd, "storage.add.ftp")
	setCommandID(addSmbCmd, "storage.add.smb")
}

var storageAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a storage provider",
	Long:  "Add a new cloud storage provider (S3, WebDAV, Local, SFTP, FTP, SMB)",
}

var addS3Cmd = &cobra.Command{
	Use:   "s3 <json>",
	Short: "Add S3 or S3-compatible storage",
	Long:  "Add Amazon S3 or S3-compatible storage provider (accepts JSON configuration)",
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

var addWebdavCmd = &cobra.Command{
	Use:   "webdav <json>",
	Short: "Add WebDAV storage",
	Long:  "Add WebDAV storage provider (accepts JSON configuration)",
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

var addLocalCmd = &cobra.Command{
	Use:   "local <path>",
	Short: "Add local folder as storage",
	Long:  "Add a local folder as a storage provider",
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

var addSftpCmd = &cobra.Command{
	Use:   "sftp <json>",
	Short: "Add SFTP storage",
	Long:  "Add SFTP storage provider (accepts JSON configuration)",
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

var addFtpCmd = &cobra.Command{
	Use:   "ftp <json>",
	Short: "Add FTP storage",
	Long:  "Add FTP storage provider (accepts JSON configuration)",
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

var addSmbCmd = &cobra.Command{
	Use:   "smb <json>",
	Short: "Add SMB/CIFS storage",
	Long:  "Add SMB/CIFS storage provider (accepts JSON configuration)",
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
