package cmd

import (
	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

func init() {
	fileCmd.AddCommand(createFolderCmd)
	fileCmd.AddCommand(createEncryptedFolderCmd)
	fileCmd.AddCommand(unlockFileCmd)
	fileCmd.AddCommand(lockFileCmd)

	createEncryptedFolderCmd.Flags().Bool("save-password", false, "Save password permanently")
	unlockFileCmd.Flags().Bool("permanent", false, "Unlock permanently (save password)")

	setCommandID(createFolderCmd, "file.mkdir")
	setCommandID(createEncryptedFolderCmd, "file.mkdir-encrypted")
	setCommandID(unlockFileCmd, "file.unlock")
	setCommandID(lockFileCmd, "file.lock")
}

var createFolderCmd = &cobra.Command{
	Use:   "mkdir [parent-path] [folder-name]",
	Short: "Create a folder",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		parentPath := args[0]
		folderName := args[1]

		result, err := cd2Client.File().CreateFolder(ctx, &pb.CreateFolderRequest{
			ParentPath: parentPath,
			FolderName: folderName,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var createEncryptedFolderCmd = &cobra.Command{
	Use:   "mkdir-encrypted [parent-path] [folder-name] [password]",
	Short: "Create an encrypted folder",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		parentPath := args[0]
		folderName := args[1]
		password := args[2]
		savePassword, _ := cmd.Flags().GetBool("save-password")

		result, err := cd2Client.File().CreateEncryptedFolder(ctx, &pb.CreateEncryptedFolderRequest{
			ParentPath:   parentPath,
			FolderName:   folderName,
			Password:     password,
			SavePassword: savePassword,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var unlockFileCmd = &cobra.Command{
	Use:   "unlock [path] [password]",
	Short: "Unlock an encrypted file or folder",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		path := args[0]
		password := args[1]
		permanent, _ := cmd.Flags().GetBool("permanent")

		result, err := cd2Client.File().UnlockEncryptedFile(ctx, &pb.UnlockEncryptedFileRequest{
			Path:            path,
			Password:        password,
			PermanentUnlock: permanent,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var lockFileCmd = &cobra.Command{
	Use:   "lock [path]",
	Short: "Lock an encrypted file or folder",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		path := args[0]

		result, err := cd2Client.File().LockEncryptedFile(ctx, path)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}
