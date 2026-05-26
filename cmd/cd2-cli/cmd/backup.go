package cmd

import (
	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup management",
}

func init() {
	rootCmd.AddCommand(backupCmd)

	backupCmd.AddCommand(listBackupsCmd)
	backupCmd.AddCommand(statusBackupCmd)
	backupCmd.AddCommand(removeBackupCmd)
	backupCmd.AddCommand(restartBackupCmd)
	backupCmd.AddCommand(addBackupCmd)
	backupCmd.AddCommand(updateBackupCmd)
	backupCmd.AddCommand(enabledBackupCmd)
	backupCmd.AddCommand(watchBackupCmd)
	backupCmd.AddCommand(strategiesBackupCmd)
	backupCmd.AddCommand(canAddBackupCmd)
	backupCmd.AddCommand(notifyBackupCmd)
	backupCmd.AddCommand(destinationBackupCmd)

	setCommandID(listBackupsCmd, "backup.list")
	setCommandID(statusBackupCmd, "backup.status")
	setCommandID(removeBackupCmd, "backup.remove")
	setCommandID(restartBackupCmd, "backup.restart")
	setCommandID(addBackupCmd, "backup.add")
	setCommandID(updateBackupCmd, "backup.update")
	setCommandID(enabledBackupCmd, "backup.enabled")
	setCommandID(watchBackupCmd, "backup.watch")
	setCommandID(strategiesBackupCmd, "backup.strategies")
	setCommandID(canAddBackupCmd, "backup.can-add")
	setCommandID(notifyBackupCmd, "backup.notify")
}

var listBackupsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all backups",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.Backup().BackupGetAll(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var statusBackupCmd = &cobra.Command{
	Use:   "status [source-path]",
	Short: "Get backup status",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		sourcePath := args[0]
		status, err := cd2Client.Backup().BackupGetStatus(ctx, sourcePath)
		if err != nil {
			return err
		}
		return outputResult(status)
	},
}

var removeBackupCmd = &cobra.Command{
	Use:   "remove [source-path]",
	Short: "Remove a backup",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		sourcePath := args[0]
		err := cd2Client.Backup().BackupRemove(ctx, sourcePath)
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "removed"})
	},
}

var restartBackupCmd = &cobra.Command{
	Use:   "restart [source-path]",
	Short: "Restart backup walking through",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		sourcePath := args[0]
		err := cd2Client.Backup().BackupRestartWalkingThrough(ctx, sourcePath)
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "restarted"})
	},
}

var addBackupCmd = &cobra.Command{
	Use:   "add <json>",
	Short: "Add a backup (accept JSON for Backup message)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		var backup pb.Backup
		if err := parseProtoJSON([]byte(args[0]), &backup); err != nil {
			return err
		}
		err := cd2Client.Backup().BackupAdd(ctx, &backup)
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "added"})
	},
}

var updateBackupCmd = &cobra.Command{
	Use:   "update <json>",
	Short: "Update a backup (accept JSON)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		var backup pb.Backup
		if err := parseProtoJSON([]byte(args[0]), &backup); err != nil {
			return err
		}
		err := cd2Client.Backup().BackupUpdate(ctx, &backup)
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "updated"})
	},
}

var enabledBackupCmd = &cobra.Command{
	Use:   "enabled <source-path> <true/false>",
	Short: "Enable or disable a backup",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		sourcePath := args[0]
		isEnabled := args[1] == "true"
		err := cd2Client.Backup().BackupSetEnabled(ctx, &pb.BackupSetEnabledRequest{
			SourcePath: sourcePath,
			IsEnabled:  isEnabled,
		})
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "updated"})
	},
}

var watchBackupCmd = &cobra.Command{
	Use:   "watch <source-path> <true/false>",
	Short: "Enable or disable filesystem watch for backup",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		sourcePath := args[0]
		isEnabled := args[1] == "true"
		err := cd2Client.Backup().BackupSetFileSystemWatchEnabled(ctx, &pb.BackupModifyRequest{
			SourcePath:             sourcePath,
			FileSystemWatchEnabled: &isEnabled,
		})
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "updated"})
	},
}

var strategiesBackupCmd = &cobra.Command{
	Use:   "strategies <json>",
	Short: "Update backup strategies (deprecated, use backup update instead)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		var req pb.BackupModifyRequest
		if err := parseProtoJSON([]byte(args[0]), &req); err != nil {
			return err
		}
		err := cd2Client.Backup().BackupUpdateStrategies(ctx, &req)
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "updated"})
	},
}

var canAddBackupCmd = &cobra.Command{
	Use:   "can-add",
	Short: "Check if can add more backups",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.Backup().CanAddMoreBackups(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var notifyBackupCmd = &cobra.Command{
	Use:   "notify <json>",
	Short: "Notify photo library changes",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		var req pb.PhotoLibraryChangeList
		if err := parseProtoJSON([]byte(args[0]), &req); err != nil {
			return err
		}
		err := cd2Client.Backup().NotifyPhotoLibraryChanges(ctx, &req)
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "notified"})
	},
}

var destinationBackupCmd = &cobra.Command{
	Use:   "destination",
	Short: "Backup destination management",
}

func init() {
	destinationBackupCmd.AddCommand(destinationAddBackupCmd)
	destinationBackupCmd.AddCommand(destinationRemoveBackupCmd)

	setCommandID(destinationAddBackupCmd, "backup.destination-add")
	setCommandID(destinationRemoveBackupCmd, "backup.destination-remove")
}

var destinationAddBackupCmd = &cobra.Command{
	Use:   "add <source-path> <dest-path>",
	Short: "Add a backup destination",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		sourcePath := args[0]
		destPath := args[1]
		err := cd2Client.Backup().BackupAddDestination(ctx, &pb.BackupModifyRequest{
			SourcePath: sourcePath,
			Destinations: []*pb.BackupDestination{
				{DestinationPath: destPath},
			},
		})
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "added"})
	},
}

var destinationRemoveBackupCmd = &cobra.Command{
	Use:   "remove <source-path> <dest-path>",
	Short: "Remove a backup destination",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		sourcePath := args[0]
		destPath := args[1]
		err := cd2Client.Backup().BackupRemoveDestination(ctx, &pb.BackupModifyRequest{
			SourcePath: sourcePath,
			Destinations: []*pb.BackupDestination{
				{DestinationPath: destPath},
			},
		})
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "removed"})
	},
}
