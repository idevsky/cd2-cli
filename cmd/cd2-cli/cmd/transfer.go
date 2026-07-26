package cmd

import (
	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfer task management",
}

var transferDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download task operations",
}

var transferUploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload task operations",
}

var downloadCountCmd = &cobra.Command{
	Use:   "count",
	Short: "Get download file count",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.Transfer().GetDownloadFileCount(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var uploadCountCmd = &cobra.Command{
	Use:   "count",
	Short: "Get upload file count",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.Transfer().GetUploadFileCount(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var downloadListCmd = &cobra.Command{
	Use:   "list",
	Short: "List download tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.Transfer().GetDownloadFileList(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var uploadListCmd = &cobra.Command{
	Use:   "list",
	Short: "List upload tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.Transfer().GetUploadFileList(ctx, &pb.GetUploadFileListRequest{})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var uploadCancelCmd = &cobra.Command{
	Use:   "cancel [key...]",
	Short: "Cancel upload tasks",
	Long:  "Cancel upload tasks. Without keys, cancels all uploads. With keys, cancels specific uploads.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		var err error
		if len(args) == 0 {
			err = cd2Client.Transfer().CancelAllUploadFiles(ctx)
		} else {
			err = cd2Client.Transfer().CancelUploadFiles(ctx, args)
		}
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "cancelled"})
	},
}

var uploadPauseCmd = &cobra.Command{
	Use:   "pause [key...]",
	Short: "Pause upload tasks",
	Long:  "Pause upload tasks. Without keys, pauses all uploads. With keys, pauses specific uploads.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		var err error
		if len(args) == 0 {
			err = cd2Client.Transfer().PauseAllUploadFiles(ctx)
		} else {
			err = cd2Client.Transfer().PauseUploadFiles(ctx, args)
		}
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "paused"})
	},
}

var uploadResumeCmd = &cobra.Command{
	Use:   "resume [key...]",
	Short: "Resume upload tasks",
	Long:  "Resume upload tasks. Without keys, resumes all uploads. With keys, resumes specific uploads.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		var err error
		if len(args) == 0 {
			err = cd2Client.Transfer().ResumeAllUploadFiles(ctx)
		} else {
			err = cd2Client.Transfer().ResumeUploadFiles(ctx, args)
		}
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "resumed"})
	},
}

var taskCountCmd = &cobra.Command{
	Use:   "count",
	Short: "Get all tasks count",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.Transfer().GetAllTasksCount(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

func init() {
	rootCmd.AddCommand(transferCmd)

	transferCmd.AddCommand(taskCountCmd)
	transferCmd.AddCommand(transferDownloadCmd)
	transferCmd.AddCommand(transferUploadCmd)

	transferDownloadCmd.AddCommand(downloadCountCmd)
	transferDownloadCmd.AddCommand(downloadListCmd)

	transferUploadCmd.AddCommand(uploadCountCmd)
	transferUploadCmd.AddCommand(uploadListCmd)
	transferUploadCmd.AddCommand(uploadCancelCmd)
	transferUploadCmd.AddCommand(uploadPauseCmd)
	transferUploadCmd.AddCommand(uploadResumeCmd)

	setCommandID(taskCountCmd, "transfer.count")
	setCommandID(transferDownloadCmd, "transfer.download")
	setCommandID(transferUploadCmd, "transfer.upload")
	setCommandID(downloadCountCmd, "transfer.download.count")
	setCommandID(uploadCountCmd, "transfer.upload.count")
	setCommandID(downloadListCmd, "transfer.download.list")
	setCommandID(uploadListCmd, "transfer.upload.list")
	setCommandID(uploadCancelCmd, "transfer.upload.cancel")
	setCommandID(uploadPauseCmd, "transfer.upload.pause")
	setCommandID(uploadResumeCmd, "transfer.upload.resume")
}
