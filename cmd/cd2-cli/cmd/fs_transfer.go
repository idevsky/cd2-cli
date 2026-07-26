package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(uploadCmd)
	rootCmd.AddCommand(downloadCmd)

	setCommandID(uploadCmd, "fs.upload")
	setCommandID(downloadCmd, "fs.download")
}

var uploadCmd = &cobra.Command{
	Use:     "upload [local-file] [remote-path]",
	Short:   "Upload a local file to remote path",
	Aliases: []string{"up"},
	Args:    cobra.ExactArgs(2),
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

var downloadCmd = &cobra.Command{
	Use:     "download [remote-path] [local-file]",
	Short:   "Download a remote file to local path",
	Aliases: []string{"dl"},
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		remotePath := args[0]
		localPath := args[1]

		result, err := cd2Client.File().DownloadRemoteFile(ctx, remotePath, localPath)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}
