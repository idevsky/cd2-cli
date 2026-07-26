package cmd

import (
	"fmt"
	"os"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

func init() {
	fileCmd.AddCommand(createFileCmd)
	fileCmd.AddCommand(writeFileCmd)
	fileCmd.AddCommand(closeFileCmd)

	createFileCmd.Flags().String("request-json", "", "JSON request body")
	createFileCmd.Flags().String("request-file", "", "File containing JSON request body")
	writeFileCmd.Flags().String("request-json", "", "JSON request body")
	writeFileCmd.Flags().String("request-file", "", "File containing JSON request body")
	writeFileCmd.Flags().Uint64("handle", 0, "File handle from create command")
	writeFileCmd.Flags().Uint64("start", 0, "Start position")
	writeFileCmd.Flags().Bool("close", false, "Close file after write")
	writeFileCmd.Flags().String("data", "", "Data to write (base64 encoded for binary)")
	closeFileCmd.Flags().Uint64("handle", 0, "File handle to close")

	setCommandID(createFileCmd, "file.create")
	setCommandID(writeFileCmd, "file.write")
	setCommandID(closeFileCmd, "file.close")
}

var createFileCmd = &cobra.Command{
	Use:   "create [parent-path] [file-name]",
	Short: "Create a new file",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()

		req := &pb.CreateFileRequest{}

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
		} else if len(args) >= 2 {
			req.ParentPath = args[0]
			req.FileName = args[1]
		} else {
			return fmt.Errorf("must provide either request-json, request-file, or parent-path and file-name arguments")
		}

		result, err := cd2Client.File().CreateFile(ctx, req)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var writeFileCmd = &cobra.Command{
	Use:   "write",
	Short: "Write data to an open file",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()

		req := &pb.WriteFileRequest{}

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
			handle, _ := cmd.Flags().GetUint64("handle")
			start, _ := cmd.Flags().GetUint64("start")
			closeFile, _ := cmd.Flags().GetBool("close")
			dataStr, _ := cmd.Flags().GetString("data")

			if handle == 0 {
				return fmt.Errorf("must provide --handle or request-json/request-file")
			}

			req.FileHandle = handle
			req.StartPos = start
			req.CloseFile = closeFile

			if dataStr != "" {
				req.Buffer = []byte(dataStr)
				req.Length = uint64(len(req.Buffer))
			}
		}

		result, err := cd2Client.File().WriteToFile(ctx, req)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var closeFileCmd = &cobra.Command{
	Use:   "close --handle [handle]",
	Short: "Close an open file handle",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()

		handle, _ := cmd.Flags().GetUint64("handle")
		if handle == 0 {
			return fmt.Errorf("must provide --handle")
		}

		result, err := cd2Client.File().CloseFile(ctx, &pb.CloseFileRequest{
			FileHandle: handle,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}
