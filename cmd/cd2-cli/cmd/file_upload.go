package cmd

import (
	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

func init() {
	fileCmd.AddCommand(uploadFileCmd)
	fileCmd.AddCommand(downloadFileCmd)

	setCommandID(uploadFileCmd, "file.upload")
	setCommandID(downloadFileCmd, "file.download")
}

var uploadFileCmd = &cobra.Command{
	Use:   "upload [local-path] [remote-path]",
	Short: "Upload a local file to remote",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		localPath := args[0]
		remotePath := args[1]

		result, err := cd2Client.File().UploadLocalFile(ctx, localPath, remotePath)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var downloadFileCmd = &cobra.Command{
	Use:   "download [remote-path] [local-path]",
	Short: "Download a remote file to local",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		remotePath := args[0]
		localPath := args[1]

		urlInfo, err := cd2Client.File().GetDownloadUrl(ctx, &pb.GetDownloadUrlPathRequest{
			Path:         remotePath,
			GetDirectUrl: true,
		})
		if err != nil {
			return err
		}

		if urlInfo.DirectUrl != nil && *urlInfo.DirectUrl != "" {
			return outputResult(map[string]interface{}{
				"success":     true,
				"downloadUrl": urlInfo.DownloadUrlPath,
				"directUrl":   *urlInfo.DirectUrl,
				"expiresIn":   urlInfo.ExpiresIn,
				"localPath":   localPath,
				"remotePath":  remotePath,
				"message":     "Direct URL available. Use external tool to download.",
			})
		}

		return outputResult(map[string]interface{}{
			"success":     true,
			"downloadUrl": urlInfo.DownloadUrlPath,
			"expiresIn":   urlInfo.ExpiresIn,
			"localPath":   localPath,
			"remotePath":  remotePath,
			"message":     "Download URL generated. Use external tool to download.",
		})
	},
}
