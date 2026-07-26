package cmd

import (
	"fmt"
	"os"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

func init() {
	fileCmd.AddCommand(deleteFileCmd)
	fileCmd.AddCommand(renameFileCmd)
	fileCmd.AddCommand(renameFilesCmd)
	fileCmd.AddCommand(moveFileCmd)
	fileCmd.AddCommand(copyFileCmd)

	deleteFileCmd.Flags().Bool("permanent", false, "Delete permanently")
	deleteFileCmd.Flags().Bool("batch", false, "Delete multiple files (provide paths as arguments)")
	deleteFileCmd.Flags().String("request-json", "", "JSON request body for complex operations")
	deleteFileCmd.Flags().String("request-file", "", "File containing JSON request body")
	renameFilesCmd.Flags().String("request-json", "", "JSON request body containing rename operations")
	renameFilesCmd.Flags().String("request-file", "", "File containing JSON request body")
	moveFileCmd.Flags().String("conflict", "overwrite", "Conflict policy: overwrite, rename, skip")
	moveFileCmd.Flags().String("request-json", "", "JSON request body for complex operations")
	moveFileCmd.Flags().String("request-file", "", "File containing JSON request body")
	copyFileCmd.Flags().String("conflict", "overwrite", "Conflict policy: overwrite, rename, skip")
	copyFileCmd.Flags().String("request-json", "", "JSON request body for complex operations")
	copyFileCmd.Flags().String("request-file", "", "File containing JSON request body")

	setCommandID(deleteFileCmd, "file.delete")
	setCommandID(renameFileCmd, "file.rename")
	setCommandID(renameFilesCmd, "file.rename-batch")
	setCommandID(moveFileCmd, "file.move")
	setCommandID(copyFileCmd, "file.copy")
}

var deleteFileCmd = &cobra.Command{
	Use:   "delete [path...]",
	Short: "Delete files or folders",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		permanent, _ := cmd.Flags().GetBool("permanent")
		batch, _ := cmd.Flags().GetBool("batch")

		requestJSON, _ := cmd.Flags().GetString("request-json")
		requestFile, _ := cmd.Flags().GetString("request-file")

		if requestJSON != "" || requestFile != "" {
			var paths []string
			if requestJSON != "" {
				var req pb.MultiFileRequest
				if err := parseProtoJSON([]byte(requestJSON), &req); err != nil {
					return fmt.Errorf("failed to parse request-json: %w", err)
				}
				paths = req.Path
			} else {
				data, err := os.ReadFile(requestFile)
				if err != nil {
					return fmt.Errorf("failed to read request-file: %w", err)
				}
				var req pb.MultiFileRequest
				if err := parseProtoJSON(data, &req); err != nil {
					return fmt.Errorf("failed to parse request-file: %w", err)
				}
				paths = req.Path
			}

			if permanent {
				result, err := cd2Client.File().DeleteFilesPermanently(ctx, paths)
				if err != nil {
					return err
				}
				return outputResult(result)
			} else {
				result, err := cd2Client.File().DeleteFiles(ctx, paths)
				if err != nil {
					return err
				}
				return outputResult(result)
			}
		}

		if batch && len(args) > 1 {
			if permanent {
				result, err := cd2Client.File().DeleteFilesPermanently(ctx, args)
				if err != nil {
					return err
				}
				return outputResult(result)
			} else {
				result, err := cd2Client.File().DeleteFiles(ctx, args)
				if err != nil {
					return err
				}
				return outputResult(result)
			}
		} else {
			path := args[0]
			if permanent {
				result, err := cd2Client.File().DeleteFilePermanently(ctx, path)
				if err != nil {
					return err
				}
				return outputResult(result)
			} else {
				result, err := cd2Client.File().DeleteFile(ctx, path)
				if err != nil {
					return err
				}
				return outputResult(result)
			}
		}
	},
}

var renameFileCmd = &cobra.Command{
	Use:   "rename [path] [new-name]",
	Short: "Rename a file or folder",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		path := args[0]
		newName := args[1]

		result, err := cd2Client.File().RenameFile(ctx, &pb.RenameFileRequest{
			TheFilePath: path,
			NewName:     newName,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var renameFilesCmd = &cobra.Command{
	Use:   "rename-batch",
	Short: "Batch rename files",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()

		req := &pb.RenameFilesRequest{}

		requestJSON, _ := cmd.Flags().GetString("request-json")
		requestFile, _ := cmd.Flags().GetString("request-file")

		if requestJSON != "" {
			if err := parseProtoJSON([]byte(requestJSON), req); err != nil {
				return fmt.Errorf("failed to parse request-json: %w", err)
			}
		} else if requestFile != "" {
			data, err := os.ReadFile(requestFile)
			if err != nil {
				return fmt.Errorf("failed to read request-file: %w", err)
			}
			if err := parseProtoJSON(data, req); err != nil {
				return fmt.Errorf("failed to parse request-file: %w", err)
			}
		} else {
			return fmt.Errorf("must provide --request-json or --request-file")
		}

		result, err := cd2Client.File().RenameFiles(ctx, req)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var moveFileCmd = &cobra.Command{
	Use:   "move [source-path...] [dest-path]",
	Short: "Move file(s) to destination",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()

		requestJSON, _ := cmd.Flags().GetString("request-json")
		requestFile, _ := cmd.Flags().GetString("request-file")

		if requestJSON != "" || requestFile != "" {
			req := &pb.MoveFileRequest{}
			if requestJSON != "" {
				if err := parseProtoJSON([]byte(requestJSON), req); err != nil {
					return fmt.Errorf("failed to parse request-json: %w", err)
				}
			} else {
				data, err := os.ReadFile(requestFile)
				if err != nil {
					return fmt.Errorf("failed to read request-file: %w", err)
				}
				if err := parseProtoJSON(data, req); err != nil {
					return fmt.Errorf("failed to parse request-file: %w", err)
				}
			}

			result, err := cd2Client.File().MoveFile(ctx, req)
			if err != nil {
				return err
			}
			return outputResult(result)
		}

		sourcePaths := args[:len(args)-1]
		destPath := args[len(args)-1]
		conflictPolicy, _ := cmd.Flags().GetString("conflict")

		var policy pb.MoveFileRequest_ConflictPolicy
		switch conflictPolicy {
		case "overwrite":
			policy = pb.MoveFileRequest_Overwrite
		case "rename":
			policy = pb.MoveFileRequest_Rename
		case "skip":
			policy = pb.MoveFileRequest_Skip
		default:
			policy = pb.MoveFileRequest_Overwrite
		}

		result, err := cd2Client.File().MoveFile(ctx, &pb.MoveFileRequest{
			TheFilePaths:   sourcePaths,
			DestPath:       destPath,
			ConflictPolicy: &policy,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var copyFileCmd = &cobra.Command{
	Use:   "copy [source-path...] [dest-path]",
	Short: "Copy file(s) to destination",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()

		requestJSON, _ := cmd.Flags().GetString("request-json")
		requestFile, _ := cmd.Flags().GetString("request-file")

		if requestJSON != "" || requestFile != "" {
			req := &pb.CopyFileRequest{}
			if requestJSON != "" {
				if err := parseProtoJSON([]byte(requestJSON), req); err != nil {
					return fmt.Errorf("failed to parse request-json: %w", err)
				}
			} else {
				data, err := os.ReadFile(requestFile)
				if err != nil {
					return fmt.Errorf("failed to read request-file: %w", err)
				}
				if err := parseProtoJSON(data, req); err != nil {
					return fmt.Errorf("failed to parse request-file: %w", err)
				}
			}

			result, err := cd2Client.File().CopyFile(ctx, req)
			if err != nil {
				return err
			}
			return outputResult(result)
		}

		sourcePaths := args[:len(args)-1]
		destPath := args[len(args)-1]
		conflictPolicy, _ := cmd.Flags().GetString("conflict")

		var policy pb.CopyFileRequest_ConflictPolicy
		switch conflictPolicy {
		case "overwrite":
			policy = pb.CopyFileRequest_Overwrite
		case "rename":
			policy = pb.CopyFileRequest_Rename
		case "skip":
			policy = pb.CopyFileRequest_Skip
		default:
			policy = pb.CopyFileRequest_Overwrite
		}

		result, err := cd2Client.File().CopyFile(ctx, &pb.CopyFileRequest{
			TheFilePaths:   sourcePaths,
			DestPath:       destPath,
			ConflictPolicy: &policy,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}
