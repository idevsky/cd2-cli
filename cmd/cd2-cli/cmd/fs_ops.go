package cmd

import (
	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rmCmd)
	rootCmd.AddCommand(mvCmd)
	rootCmd.AddCommand(cpCmd)

	rmCmd.Flags().Bool("permanent", false, "Delete permanently (bypass trash)")
	mvCmd.Flags().String("conflict", "overwrite", "Conflict policy: overwrite, rename, skip")
	cpCmd.Flags().String("conflict", "overwrite", "Conflict policy: overwrite, rename, skip")

	setCommandID(rmCmd, "fs.rm")
	setCommandID(mvCmd, "fs.mv")
	setCommandID(cpCmd, "fs.cp")
}

var rmCmd = &cobra.Command{
	Use:     "rm [path...]",
	Short:   "Remove files or directories",
	Aliases: []string{"delete", "remove"},
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		permanent, _ := cmd.Flags().GetBool("permanent")

		if len(args) > 1 {
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

var mvCmd = &cobra.Command{
	Use:     "mv [source...] [destination]",
	Short:   "Move/rename files or directories",
	Aliases: []string{"move", "rename"},
	Args:    cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		sources := args[:len(args)-1]
		dest := args[len(args)-1]
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
			TheFilePaths:   sources,
			DestPath:       dest,
			ConflictPolicy: &policy,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var cpCmd = &cobra.Command{
	Use:     "cp [source...] [destination]",
	Short:   "Copy files or directories",
	Aliases: []string{"copy"},
	Args:    cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		sources := args[:len(args)-1]
		dest := args[len(args)-1]
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
			TheFilePaths:   sources,
			DestPath:       dest,
			ConflictPolicy: &policy,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}
