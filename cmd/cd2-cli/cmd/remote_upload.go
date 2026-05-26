package cmd

import (
	"fmt"
	"strconv"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

var remoteUploadCmd = &cobra.Command{
	Use:   "remote-upload",
	Short: "Remote upload management",
}

func init() {
	rootCmd.AddCommand(remoteUploadCmd)

	remoteUploadCmd.AddCommand(remoteUploadStartCmd)
	remoteUploadCmd.AddCommand(remoteUploadControlCmd)
	remoteUploadCmd.AddCommand(remoteUploadReadDataCmd)
	remoteUploadCmd.AddCommand(remoteUploadHashProgressCmd)

	setCommandID(remoteUploadStartCmd, "remote-upload.start")
	setCommandID(remoteUploadControlCmd, "remote-upload.control")
	setCommandID(remoteUploadReadDataCmd, "remote-upload.read-data")
	setCommandID(remoteUploadHashProgressCmd, "remote-upload.hash-progress")
}

var remoteUploadStartCmd = &cobra.Command{
	Use:   "start <file-path> <file-size>",
	Short: "Start a remote upload",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		filePath := args[0]
		fileSize, err := strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid file size: %v", err)
		}

		clientCanCalcHashes, _ := cmd.Flags().GetBool("client-can-calculate-hashes")

		result, err := cd2Client.RemoteUpload().StartRemoteUpload(ctx, &pb.StartRemoteUploadRequest{
			FilePath:                 filePath,
			FileSize:                 fileSize,
			ClientCanCalculateHashes: clientCanCalcHashes,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var remoteUploadControlCmd = &cobra.Command{
	Use:   "control <upload-id> <action>",
	Short: "Control a remote upload (cancel/pause/resume)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		uploadID := args[0]
		action := args[1]

		var req *pb.RemoteUploadControlRequest
		switch action {
		case "cancel":
			req = &pb.RemoteUploadControlRequest{
				UploadId: uploadID,
				Control:  &pb.RemoteUploadControlRequest_Cancel{Cancel: &pb.CancelRemoteUpload{}},
			}
		case "pause":
			req = &pb.RemoteUploadControlRequest{
				UploadId: uploadID,
				Control:  &pb.RemoteUploadControlRequest_Pause{Pause: &pb.PauseRemoteUpload{}},
			}
		case "resume":
			req = &pb.RemoteUploadControlRequest{
				UploadId: uploadID,
				Control:  &pb.RemoteUploadControlRequest_Resume{Resume: &pb.ResumeRemoteUpload{}},
			}
		default:
			return fmt.Errorf("invalid action: must be cancel, pause, or resume")
		}

		err := cd2Client.RemoteUpload().RemoteUploadControl(ctx, req)
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": action + "ed"})
	},
}

var remoteUploadReadDataCmd = &cobra.Command{
	Use:   "read-data <upload-id> <offset> <length>",
	Short: "Send read data for remote upload",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		uploadID := args[0]
		offset, err := strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid offset: %v", err)
		}
		length, err := strconv.ParseUint(args[2], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid length: %v", err)
		}

		lazyRead, _ := cmd.Flags().GetBool("lazy-read")
		isLastChunk, _ := cmd.Flags().GetBool("is-last-chunk")

		result, err := cd2Client.RemoteUpload().RemoteReadData(ctx, &pb.RemoteReadDataUpload{
			UploadId:    uploadID,
			Offset:      offset,
			Length:      length,
			LazyRead:    lazyRead,
			IsLastChunk: isLastChunk,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var remoteUploadHashProgressCmd = &cobra.Command{
	Use:   "hash-progress <upload-id> <bytes-hashed> <total-bytes>",
	Short: "Send hash progress for remote upload",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		uploadID := args[0]
		bytesHashed, err := strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid bytes-hashed: %v", err)
		}
		totalBytes, err := strconv.ParseUint(args[2], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid total-bytes: %v", err)
		}

		hashTypeStr, _ := cmd.Flags().GetString("hash-type")
		hashValueStr, _ := cmd.Flags().GetString("hash-value")

		var hashType pb.CloudDriveFile_HashType
		switch hashTypeStr {
		case "sha1":
			hashType = pb.CloudDriveFile_Sha1
		case "md5":
			hashType = pb.CloudDriveFile_Md5
		case "pikpak-sha1":
			hashType = pb.CloudDriveFile_PikPakSha1
		default:
			hashType = pb.CloudDriveFile_Sha1
		}

		req := &pb.RemoteHashProgressUpload{
			UploadId:    uploadID,
			BytesHashed: bytesHashed,
			TotalBytes:  totalBytes,
			HashType:    hashType,
		}
		if hashValueStr != "" {
			req.HashValue = &hashValueStr
		}

		result, err := cd2Client.RemoteUpload().RemoteHashProgress(ctx, req)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

func init() {
	remoteUploadStartCmd.Flags().Bool("client-can-calculate-hashes", false, "Client can calculate hashes locally")
	remoteUploadReadDataCmd.Flags().Bool("lazy-read", false, "Lazy read mode")
	remoteUploadReadDataCmd.Flags().Bool("is-last-chunk", false, "Is the last chunk of data")
	remoteUploadHashProgressCmd.Flags().String("hash-type", "sha1", "Hash type (sha1, md5, pikpak-sha1)")
	remoteUploadHashProgressCmd.Flags().String("hash-value", "", "Hash value")
}
