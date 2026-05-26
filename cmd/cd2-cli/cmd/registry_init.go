package cmd

import (
	"sync"

	"github.com/clouddrive/cd2-cli/internal/registry"
)

var registryOnce sync.Once

func registerCommands() {
	registryOnce.Do(func() {
		registerSystemCommands()
		registerAuthCommands()
		registerFileCommands()
		registerMountCommands()
		registerStorageCommands()
		registerBackupCommands()
		registerTransferCommands()
		registerTokenCommands()
		registerSessionCommands()
		registerWebdavCommands()
		registerOfflineCommands()
		registerWebhookCommands()
		registerSyncCommands()
		registerCopyCommands()
		registerCacheCommands()
		registerPromotionCommands()
		registerRemoteUploadCommands()
		registerCloudapiCommands()
		registerWhitelistCommands()
		registerTaskCommands()

		markCompletionAsLocal()
	})
}

func markCompletionAsLocal() {
	if completionCmd, _, err := rootCmd.Find([]string{"completion"}); err == nil {
		markAsLocal(completionCmd)
	}
	if helpCmd, _, err := rootCmd.Find([]string{"help"}); err == nil {
		markAsLocal(helpCmd)
	}
}

func registerSystemCommands() {
	registry.MustRegister(&registry.CommandSpec{
		ID:          "system.info",
		Category:    "system",
		RPC:         "GetSystemInfo",
		Description: "Get public system information",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "system.runtime",
		Category:    "system",
		RPC:         "GetRuntimeInfo",
		Description: "Get runtime information",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "system.capabilities",
		Category:    "system",
		RPC:         "GetServiceCapabilities",
		Description: "Get service capabilities",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "system.restart",
		Category:    "system",
		RPC:         "RestartService",
		Description: "Restart CloudDrive2 service",
		RiskLevel:   registry.RiskCritical,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "system.shutdown",
		Category:    "system",
		RPC:         "ShutdownService",
		Description: "Shutdown CloudDrive2 service",
		RiskLevel:   registry.RiskCritical,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "system.running",
		Category:    "system",
		RPC:         "GetRunningInfo",
		Description: "Get running information",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "system.settings-get",
		Category:    "system",
		RPC:         "GetSystemSettings",
		Description: "Get system settings",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "system.settings-set",
		Category:    "system",
		RPC:         "SetSystemSettings",
		Description: "Set system settings",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "system.open-files",
		Category:    "system",
		RPC:         "GetOpenFileHandles",
		Description: "Get open file handles",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "system.open-file-table",
		Category:    "system",
		RPC:         "GetOpenFileTable",
		Description: "Get open file table",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "system.cache-stats",
		Category:    "system",
		RPC:         "GetFileBufferDiskCacheStats",
		Description: "Get file buffer disk cache stats",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "system.purge-cache",
		Category:    "system",
		RPC:         "PurgeFileBufferDiskCache",
		Description: "Purge file buffer disk cache",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "system.list-disk-cache",
		Category:    "system",
		RPC:         "ListDiskCacheFolders",
		Description: "List disk cache folders",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "system.set-folder-disk-cache",
		Category:    "system",
		RPC:         "SetFolderDiskCache",
		Description: "Set folder disk cache configuration",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "system.remove-folder-disk-cache",
		Category:    "system",
		RPC:         "RemoveFolderDiskCache",
		Description: "Remove folder disk cache",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "system.set-eviction-strategy",
		Category:    "system",
		RPC:         "SetDiskCacheEvictionStrategy",
		Description: "Set disk cache eviction strategy",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "system.force-expire-dir-cache",
		Category:    "system",
		RPC:         "ForceExpireDirCache",
		Description: "Force expire directory cache",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "system.machine-id",
		Category:    "system",
		RPC:         "GetMachineId",
		Description: "Get machine ID",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "system.devices",
		Category:    "system",
		RPC:         "GetOnlineDevices",
		Description: "Get online devices",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "system.kickout-device",
		Category:    "system",
		RPC:         "KickoutDevice",
		Description: "Kick out a device",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "system.logs",
		Category:    "system",
		RPC:         "ListLogFiles",
		Description: "List log files",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "system.temp-files",
		Category:    "system",
		RPC:         "GetTempFileTable",
		Description: "Get temp file table",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "system.dir-cache-size",
		Category:    "system",
		RPC:         "GetDirCacheDbSize",
		Description: "Get directory cache DB size",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "system.dir-cache-vacuum",
		Category:    "system",
		RPC:         "VacuumDirCache",
		Description: "Vacuum directory cache",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "system.dir-cache-progress",
		Category:    "system",
		RPC:         "GetVacuumProgress",
		Description: "Get vacuum progress",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "system.dir-cache-table",
		Category:    "system",
		RPC:         "GetDirCacheTable",
		Description: "Get directory cache table",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "system.update-check",
		Category:    "system",
		RPC:         "CheckUpdate",
		Description: "Check for system updates",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "system.update-download",
		Category:    "system",
		RPC:         "DownloadUpdate",
		Description: "Download system update",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "system.update-apply",
		Category:    "system",
		RPC:         "UpdateSystem",
		Description: "Apply system update",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "system.has-update",
		Category:    "system",
		RPC:         "HasUpdate",
		Description: "Check if update is available",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "system.webserver-get",
		Category:    "system",
		RPC:         "GetWebServerConfig",
		Description: "Get WebServer configuration",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "system.webserver-set",
		Category:    "system",
		RPC:         "SetWebServerConfig",
		Description: "Set WebServer configuration",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "system.generate-cert",
		Category:    "system",
		RPC:         "GenerateSelfSignedCert",
		Description: "Generate self-signed certificate",
		RiskLevel:   registry.RiskHigh,
	})
}

func registerAuthCommands() {
	registry.MustRegister(&registry.CommandSpec{
		ID:          "auth.login",
		Category:    "auth",
		RPC:         "GetToken",
		Description: "Login and get JWT token",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "auth.login-2fa",
		Category:    "auth",
		RPC:         "LoginWith2FA",
		Description: "Login with 2FA code",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "auth.login-third-party",
		Category:    "auth",
		RPC:         "LoginWithThirdPartyAccount",
		Description: "Login with third-party account",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "auth.logout",
		Category:    "auth",
		RPC:         "Logout",
		Description: "Logout current session",
		RiskLevel:   registry.RiskMedium,
		DefaultOpen: registry.BoolPtr(false),
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "auth.status",
		Category:    "auth",
		RPC:         "GetAccountStatus",
		Description: "Get account status",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "auth.change-password",
		Category:    "auth",
		RPC:         "ChangePassword",
		Description: "Change account password",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "auth.change-email",
		Category:    "auth",
		RPC:         "ChangeEmail",
		Description: "Change account email",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "auth.confirm-email",
		Category:    "auth",
		RPC:         "ConfirmEmail",
		Description: "Confirm email with verification code",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "auth.send-confirm-email",
		Category:    "auth",
		RPC:         "SendConfirmEmail",
		Description: "Request email confirmation code",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "auth.send-change-email-code",
		Category:    "auth",
		RPC:         "SendChangeEmailCode",
		Description: "Request change email verification code",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "auth.change-email-and-password",
		Category:    "auth",
		RPC:         "ChangeEmailAndPassword",
		Description: "Change email and password (trusted devices only)",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "auth.register",
		Category:    "auth",
		RPC:         "Register",
		Description: "Register new account",
		RiskLevel:   registry.RiskMedium,
		DefaultOpen: registry.BoolPtr(false),
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "auth.send-reset-email",
		Category:    "auth",
		RPC:         "SendResetAccountEmail",
		Description: "Request account reset email",
		RiskLevel:   registry.RiskMedium,
		DefaultOpen: registry.BoolPtr(false),
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "auth.reset-account",
		Category:    "auth",
		RPC:         "ResetAccount",
		Description: "Reset account with reset code",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "auth.2fa-status",
		Category:    "auth",
		RPC:         "Check2FAStatus",
		Description: "Check 2FA status",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "auth.2fa-setup",
		Category:    "auth",
		RPC:         "Setup2FA",
		Description: "Setup 2FA",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "auth.2fa-enable",
		Category:    "auth",
		RPC:         "Enable2FA",
		Description: "Enable 2FA",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "auth.2fa-disable",
		Category:    "auth",
		RPC:         "Disable2FA",
		Description: "Disable 2FA",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "auth.2fa-recovery-codes",
		Category:    "auth",
		RPC:         "GetRecoveryCodes",
		Description: "Get 2FA recovery codes",
		RiskLevel:   registry.RiskMedium,
		DefaultOpen: registry.BoolPtr(false),
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "auth.2fa-regenerate-codes",
		Category:    "auth",
		RPC:         "RegenerateRecoveryCodes",
		Description: "Regenerate 2FA recovery codes",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "auth.token-set",
		Category:    "auth",
		RPC:         "",
		Description: "Set authentication token in config",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "auth.token-clear",
		Category:    "auth",
		RPC:         "",
		Description: "Clear authentication token from config",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "auth.token-show",
		Category:    "auth",
		RPC:         "",
		Description: "Show current authentication token",
		RiskLevel:   registry.RiskLow,
	})
}

func registerFileCommands() {
	registry.MustRegister(&registry.CommandSpec{
		ID:          "file.list",
		Category:    "file",
		RPC:         "GetSubFiles",
		Description: "List files in directory",
		RiskLevel:   registry.RiskLow,
		AliasGroup:  "file.list",
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "file.find",
		Category:    "file",
		RPC:         "FindFileByPath",
		Description: "Find file by path",
		RiskLevel:   registry.RiskLow,
		AliasGroup:  "file.find",
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "file.search",
		Category:    "file",
		RPC:         "GetSearchResults",
		Description: "Search for files",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "file.mkdir",
		Category:    "file",
		RPC:         "CreateFolder",
		Description: "Create folder",
		RiskLevel:   registry.RiskMedium,
		DefaultOpen: registry.BoolPtr(false),
		AliasGroup:  "file.mkdir",
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "file.mkdir-encrypted",
		Category:    "file",
		RPC:         "CreateEncryptedFolder",
		Description: "Create encrypted folder",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "file.lock",
		Category:    "file",
		RPC:         "LockEncryptedFile",
		Description: "Lock encrypted file/folder",
		RiskLevel:   registry.RiskMedium,
		DefaultOpen: registry.BoolPtr(false),
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "file.unlock",
		Category:    "file",
		RPC:         "UnlockEncryptedFile",
		Description: "Unlock encrypted file/folder",
		RiskLevel:   registry.RiskMedium,
		DefaultOpen: registry.BoolPtr(false),
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "file.delete",
		Category:    "file",
		RPC:         "DeleteFile",
		Description: "Delete file",
		RiskLevel:   registry.RiskHigh,
		AliasGroup:  "file.delete",
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "file.delete-permanently",
		Category:    "file",
		RPC:         "DeleteFilePermanently",
		Description: "Delete file permanently",
		RiskLevel:   registry.RiskCritical,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "file.rename",
		Category:    "file",
		RPC:         "RenameFile",
		Description: "Rename file",
		RiskLevel:   registry.RiskMedium,
		DefaultOpen: registry.BoolPtr(false),
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "file.rename-batch",
		Category:    "file",
		RPC:         "RenameFiles",
		Description: "Batch rename files",
		RiskLevel:   registry.RiskMedium,
		DefaultOpen: registry.BoolPtr(false),
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "file.move",
		Category:    "file",
		RPC:         "MoveFile",
		Description: "Move files",
		RiskLevel:   registry.RiskMedium,
		DefaultOpen: registry.BoolPtr(false),
		AliasGroup:  "file.move",
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "file.copy",
		Category:    "file",
		RPC:         "CopyFile",
		Description: "Copy files",
		RiskLevel:   registry.RiskMedium,
		DefaultOpen: registry.BoolPtr(false),
		AliasGroup:  "file.copy",
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "file.download-url",
		Category:    "file",
		RPC:         "GetDownloadUrlPath",
		Description: "Get download URL for file",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "file.properties",
		Category:    "file",
		RPC:         "GetFileDetailProperties",
		Description: "Get file detail properties",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "file.space",
		Category:    "file",
		RPC:         "GetSpaceInfo",
		Description: "Get space info for path",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "file.memberships",
		Category:    "file",
		RPC:         "GetCloudMemberships",
		Description: "Get cloud account memberships",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "file.metadata",
		Category:    "file",
		RPC:         "GetMetaData",
		Description: "Get file metadata",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "file.original-path",
		Category:    "file",
		RPC:         "GetOriginalPath",
		Description: "Get original path from search result",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "file.add-shared-link",
		Category:    "file",
		RPC:         "AddSharedLink",
		Description: "Add shared link to folder",
		RiskLevel:   registry.RiskMedium,
		DefaultOpen: registry.BoolPtr(false),
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "file.upload",
		Category:    "file",
		RPC:         "UploadLocalFile",
		Description: "Upload local file to remote",
		RiskLevel:   registry.RiskHigh,
		DefaultOpen: registry.BoolPtr(false),
		AliasGroup:  "file.upload",
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "file.download",
		Category:    "file",
		RPC:         "DownloadFile",
		Description: "Download remote file to local",
		RiskLevel:   registry.RiskLow,
		AliasGroup:  "file.download",
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "file.create",
		Category:    "file",
		RPC:         "CreateFile",
		Description: "Create a new file",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "file.write",
		Category:    "file",
		RPC:         "WriteToFile",
		Description: "Write data to file",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "file.close",
		Category:    "file",
		RPC:         "CloseFile",
		Description: "Close file handle",
		RiskLevel:   registry.RiskMedium,
		DefaultOpen: registry.BoolPtr(false),
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "file.stat",
		Category:    "file",
		RPC:         "FindFileByPath",
		Description: "Get file information",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "fs.ls",
		Category:    "fs",
		RPC:         "GetSubFiles",
		Description: "List directory contents",
		RiskLevel:   registry.RiskLow,
		AliasGroup:  "file.list",
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "fs.stat",
		Category:    "fs",
		RPC:         "FindFileByPath",
		Description: "Get file info",
		RiskLevel:   registry.RiskLow,
		AliasGroup:  "file.find",
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "fs.mkdir",
		Category:    "fs",
		RPC:         "CreateFolder",
		Description: "Create directory",
		RiskLevel:   registry.RiskMedium,
		DefaultOpen: registry.BoolPtr(false),
		AliasGroup:  "file.mkdir",
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "fs.rm",
		Category:    "fs",
		RPC:         "DeleteFile",
		Description: "Remove file or directory",
		RiskLevel:   registry.RiskHigh,
		AliasGroup:  "file.delete",
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "fs.rm-permanent",
		Category:    "fs",
		RPC:         "DeleteFilePermanently",
		Description: "Remove file or directory permanently",
		RiskLevel:   registry.RiskCritical,
		AliasGroup:  "file.delete-permanently",
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "fs.mv",
		Category:    "fs",
		RPC:         "MoveFile",
		Description: "Move file or directory",
		RiskLevel:   registry.RiskMedium,
		DefaultOpen: registry.BoolPtr(false),
		AliasGroup:  "file.move",
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "fs.cp",
		Category:    "fs",
		RPC:         "CopyFile",
		Description: "Copy file or directory",
		RiskLevel:   registry.RiskMedium,
		DefaultOpen: registry.BoolPtr(false),
		AliasGroup:  "file.copy",
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "fs.upload",
		Category:    "fs",
		RPC:         "UploadLocalFile",
		Description: "Upload local file",
		RiskLevel:   registry.RiskMedium,
		DefaultOpen: registry.BoolPtr(false),
		AliasGroup:  "file.upload",
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "fs.download",
		Category:    "fs",
		RPC:         "DownloadFile",
		Description: "Download file",
		RiskLevel:   registry.RiskLow,
		AliasGroup:  "file.download",
	})
}

func registerMountCommands() {
	registry.MustRegister(&registry.CommandSpec{
		ID:          "mount.list",
		Category:    "mount",
		RPC:         "GetMountPoints",
		Description: "List mount points",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "mount.add",
		Category:    "mount",
		RPC:         "AddMountPoint",
		Description: "Add mount point",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "mount.remove",
		Category:    "mount",
		RPC:         "RemoveMountPoint",
		Description: "Remove mount point",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "mount.start",
		Category:    "mount",
		RPC:         "Mount",
		Description: "Start mount point",
		RiskLevel:   registry.RiskMedium,
		DefaultOpen: registry.BoolPtr(false),
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "mount.stop",
		Category:    "mount",
		RPC:         "Unmount",
		Description: "Stop mount point",
		RiskLevel:   registry.RiskMedium,
		DefaultOpen: registry.BoolPtr(false),
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "mount.status",
		Category:    "mount",
		RPC:         "GetMountPoints",
		Description: "Get mount point status",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "mount.can-add",
		Category:    "mount",
		RPC:         "CanAddMoreMountPoints",
		Description: "Check if can add more mount points",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "mount.update",
		Category:    "mount",
		RPC:         "UpdateMountPoint",
		Description: "Update mount point",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "mount.drive-letters",
		Category:    "mount",
		RPC:         "GetAvailableDriveLetters",
		Description: "Get available drive letters",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "mount.has-drive-letters",
		Category:    "mount",
		RPC:         "HasDriveLetters",
		Description: "Check if has drive letters",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "mount.can-mount-both",
		Category:    "mount",
		RPC:         "CanMountBothLocalAndCloud",
		Description: "Check if can mount both local and cloud",
		RiskLevel:   registry.RiskLow,
	})
}

func registerStorageCommands() {
	registry.MustRegister(&registry.CommandSpec{
		ID:          "storage.list",
		Category:    "storage",
		RPC:         "GetAllCloudApis",
		Description: "List cloud storage APIs",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "storage.config",
		Category:    "storage",
		RPC:         "GetCloudAPIConfig",
		Description: "Get storage config",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "storage.set-config",
		Category:    "storage",
		RPC:         "SetCloudAPIConfig",
		Description: "Set storage config",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "storage.remove",
		Category:    "storage",
		RPC:         "RemoveCloudAPI",
		Description: "Remove storage",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "storage.add",
		Category:    "storage",
		RPC:         "APILogin",
		Description: "Add storage",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "storage.status",
		Category:    "storage",
		RPC:         "GetAllCloudApis",
		Description: "Get storage status",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "storage.can-add",
		Category:    "storage",
		RPC:         "CanAddMoreCloudApis",
		Description: "Check if can add more storage",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "storage.add.s3",
		Category:    "storage",
		RPC:         "APILoginS3",
		Description: "Add S3 storage",
		RiskLevel:   registry.RiskHigh,
		DefaultOpen: registry.BoolPtr(false),
		AliasGroup:  "storage.add.s3",
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "storage.add.webdav",
		Category:    "storage",
		RPC:         "APILoginWebDav",
		Description: "Add WebDAV storage",
		RiskLevel:   registry.RiskHigh,
		DefaultOpen: registry.BoolPtr(false),
		AliasGroup:  "storage.add.webdav",
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "storage.add.local",
		Category:    "storage",
		RPC:         "APIAddLocalFolder",
		Description: "Add local storage",
		RiskLevel:   registry.RiskHigh,
		DefaultOpen: registry.BoolPtr(false),
		AliasGroup:  "storage.add.local",
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "storage.add.sftp",
		Category:    "storage",
		RPC:         "APILoginSftp",
		Description: "Add SFTP storage",
		RiskLevel:   registry.RiskHigh,
		DefaultOpen: registry.BoolPtr(false),
		AliasGroup:  "storage.add.sftp",
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "storage.add.ftp",
		Category:    "storage",
		RPC:         "APILoginFtp",
		Description: "Add FTP storage",
		RiskLevel:   registry.RiskHigh,
		DefaultOpen: registry.BoolPtr(false),
		AliasGroup:  "storage.add.ftp",
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "storage.add.smb",
		Category:    "storage",
		RPC:         "APILoginSmb",
		Description: "Add SMB storage",
		RiskLevel:   registry.RiskHigh,
		DefaultOpen: registry.BoolPtr(false),
		AliasGroup:  "storage.add.smb",
	})
}

func registerBackupCommands() {
	registry.MustRegister(&registry.CommandSpec{
		ID:          "backup.list",
		Category:    "backup",
		RPC:         "BackupGetAll",
		Description: "List backups",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "backup.status",
		Category:    "backup",
		RPC:         "BackupGetStatus",
		Description: "Get backup status",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "backup.remove",
		Category:    "backup",
		RPC:         "BackupRemove",
		Description: "Remove backup",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "backup.restart",
		Category:    "backup",
		RPC:         "BackupRestartWalkingThrough",
		Description: "Restart backup scan",
		RiskLevel:   registry.RiskMedium,
		DefaultOpen: registry.BoolPtr(false),
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "backup.add",
		Category:    "backup",
		RPC:         "BackupAdd",
		Description: "Add a backup",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "backup.update",
		Category:    "backup",
		RPC:         "BackupUpdate",
		Description: "Update backup configuration",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "backup.enabled",
		Category:    "backup",
		RPC:         "BackupSetEnabled",
		Description: "Enable or disable a backup",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "backup.watch",
		Category:    "backup",
		RPC:         "BackupSetFileSystemWatchEnabled",
		Description: "Enable or disable filesystem watch for backup",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "backup.strategies",
		Category:    "backup",
		RPC:         "BackupUpdateStrategies",
		Description: "Update backup strategies",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "backup.can-add",
		Category:    "backup",
		RPC:         "CanAddMoreBackups",
		Description: "Check if can add more backups",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "backup.notify",
		Category:    "backup",
		RPC:         "NotifyPhotoLibraryChanges",
		Description: "Notify photo library changes",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "backup.destination-add",
		Category:    "backup",
		RPC:         "BackupAddDestination",
		Description: "Add backup destination",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "backup.destination-remove",
		Category:    "backup",
		RPC:         "BackupRemoveDestination",
		Description: "Remove backup destination",
		RiskLevel:   registry.RiskHigh,
	})
}

func registerTransferCommands() {
	registry.MustRegister(&registry.CommandSpec{
		ID:          "transfer.count",
		Category:    "transfer",
		RPC:         "GetAllTasksCount",
		Description: "Get task counts",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "transfer.download",
		Category:    "transfer",
		Description: "Download task operations",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "transfer.download.count",
		Category:    "transfer",
		RPC:         "GetDownloadFileCount",
		Description: "Get download file count",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "transfer.download.list",
		Category:    "transfer",
		RPC:         "GetDownloadFileList",
		Description: "List download tasks",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "transfer.upload",
		Category:    "transfer",
		Description: "Upload task operations",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "transfer.upload.count",
		Category:    "transfer",
		RPC:         "GetUploadFileCount",
		Description: "Get upload file count",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "transfer.upload.list",
		Category:    "transfer",
		RPC:         "GetUploadFileList",
		Description: "List upload tasks",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "transfer.upload.cancel",
		Category:    "transfer",
		RPC:         "CancelAllUploadFiles,CancelUploadFiles",
		Description: "Cancel upload tasks",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "transfer.upload.pause",
		Category:    "transfer",
		RPC:         "PauseAllUploadFiles,PauseUploadFiles",
		Description: "Pause upload tasks",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "transfer.upload.resume",
		Category:    "transfer",
		RPC:         "ResumeAllUploadFiles,ResumeUploadFiles",
		Description: "Resume upload tasks",
		RiskLevel:   registry.RiskHigh,
	})
}

func registerTokenCommands() {
	registry.MustRegister(&registry.CommandSpec{
		ID:          "token.list",
		Category:    "token",
		RPC:         "ListTokens",
		Description: "List tokens",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "token.create",
		Category:    "token",
		RPC:         "CreateToken",
		Description: "Create token",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "token.modify",
		Category:    "token",
		RPC:         "ModifyToken",
		Description: "Modify token",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "token.remove",
		Category:    "token",
		RPC:         "RemoveToken",
		Description: "Remove token",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "token.info",
		Category:    "token",
		RPC:         "GetApiTokenInfo",
		Description: "Get token info",
		RiskLevel:   registry.RiskLow,
	})
}

func registerSessionCommands() {
	registry.MustRegister(&registry.CommandSpec{
		ID:          "session.list",
		Category:    "session",
		RPC:         "GetSessions",
		Description: "List sessions",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "session.revoke",
		Category:    "session",
		RPC:         "RevokeSession",
		Description: "Revoke session",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "session.revoke-others",
		Category:    "session",
		RPC:         "RevokeOtherSessions",
		Description: "Revoke other sessions",
		RiskLevel:   registry.RiskHigh,
	})
}

func registerWebdavCommands() {
	registry.MustRegister(&registry.CommandSpec{
		ID:          "webdav.user.get",
		Category:    "webdav",
		RPC:         "GetDavUser",
		Description: "Get WebDAV user info",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "webdav.user.add",
		Category:    "webdav",
		RPC:         "AddDavUser",
		Description: "Add WebDAV user",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "webdav.user.modify",
		Category:    "webdav",
		RPC:         "ModifyDavUser",
		Description: "Modify WebDAV user",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "webdav.user.remove",
		Category:    "webdav",
		RPC:         "RemoveDavUser",
		Description: "Remove WebDAV user",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "webdav.server.get",
		Category:    "webdav",
		RPC:         "GetDavServerConfig",
		Description: "Get WebDAV server config",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "webdav.server.set",
		Category:    "webdav",
		RPC:         "SetDavServerConfig",
		Description: "Set WebDAV server config",
		RiskLevel:   registry.RiskHigh,
	})
}

func registerOfflineCommands() {
	registry.MustRegister(&registry.CommandSpec{
		ID:          "offline.list",
		Category:    "offline",
		RPC:         "ListOfflineFilesByPath",
		Description: "List offline files by path",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "offline.list-all",
		Category:    "offline",
		RPC:         "ListAllOfflineFiles",
		Description: "List all offline files with pagination",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "offline.quota",
		Category:    "offline",
		RPC:         "GetOfflineQuotaInfo",
		Description: "Get offline download quota info",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "offline.add",
		Category:    "offline",
		RPC:         "AddOfflineFiles",
		Description: "Add offline download task",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "offline.remove",
		Category:    "offline",
		RPC:         "RemoveOfflineFiles",
		Description: "Remove offline download task",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "offline.clear",
		Category:    "offline",
		RPC:         "ClearOfflineFiles",
		Description: "Clear offline download tasks",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "offline.restart",
		Category:    "offline",
		RPC:         "RestartOfflineTask",
		Description: "Restart offline download task",
		RiskLevel:   registry.RiskHigh,
	})
}

func registerWebhookCommands() {
	registry.MustRegister(&registry.CommandSpec{
		ID:          "webhook.list",
		Category:    "webhook",
		RPC:         "GetWebhookConfigs",
		Description: "List webhooks",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "webhook.template",
		Category:    "webhook",
		RPC:         "GetWebhookConfigTemplate",
		Description: "Get webhook config template",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "webhook.add",
		Category:    "webhook",
		RPC:         "AddWebhookConfig",
		Description: "Add webhook",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "webhook.remove",
		Category:    "webhook",
		RPC:         "RemoveWebhookConfig",
		Description: "Remove webhook",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "webhook.change",
		Category:    "webhook",
		RPC:         "ChangeWebhookConfig",
		Description: "Change webhook",
		RiskLevel:   registry.RiskHigh,
	})
}

func registerSyncCommands() {
	registry.MustRegister(&registry.CommandSpec{
		ID:          "sync.file-changes",
		Category:    "sync",
		RPC:         "SyncFileChangesFromCloud",
		Description: "Sync file changes from cloud",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "sync.start-listener",
		Category:    "sync",
		RPC:         "StartCloudEventListener",
		Description: "Start sync listener",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "sync.stop-listener",
		Category:    "sync",
		RPC:         "StopCloudEventListener",
		Description: "Stop sync listener",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "sync.walk-test",
		Category:    "sync",
		RPC:         "WalkThroughFolderTest",
		Description: "Walk through folder test",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "sync.cd1-user-data",
		Category:    "sync",
		RPC:         "GetCloudDrive1UserData",
		Description: "Get CloudDrive1 user data",
		RiskLevel:   registry.RiskLow,
	})
}

func registerCopyCommands() {
	registry.MustRegister(&registry.CommandSpec{
		ID:          "copy.tasks",
		Category:    "copy",
		RPC:         "GetCopyTasks",
		Description: "List copy tasks",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "copy.merge-tasks",
		Category:    "copy",
		RPC:         "GetMergeTasks",
		Description: "List merge tasks",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "copy.cancel",
		Category:    "copy",
		RPC:         "CancelCopyTask",
		Description: "Cancel copy task",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "copy.pause",
		Category:    "copy",
		RPC:         "PauseCopyTask",
		Description: "Pause copy task",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "copy.restart",
		Category:    "copy",
		RPC:         "RestartCopyTask",
		Description: "Restart copy task",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "copy.remove",
		Category:    "copy",
		RPC:         "CancelCopyTask",
		Description: "Remove copy task",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "copy.resume",
		Category:    "copy",
		RPC:         "PauseCopyTask",
		Description: "Resume paused copy task",
		RiskLevel:   registry.RiskHigh,
	})
}

func registerCacheCommands() {
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cache.stats",
		Category:    "cache",
		RPC:         "GetFileBufferDiskCacheStats",
		Description: "Get disk cache stats",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cache.purge",
		Category:    "cache",
		RPC:         "PurgeFileBufferDiskCache",
		Description: "Purge disk cache",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cache.dir-table",
		Category:    "cache",
		RPC:         "GetDirCacheTable",
		Description: "Get directory cache table",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cache.temp-table",
		Category:    "cache",
		RPC:         "GetTempFileTable",
		Description: "Get temp file table",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cache.list-disk",
		Category:    "cache",
		RPC:         "ListDiskCacheFolders",
		Description: "List disk cache folders",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cache.set-eviction",
		Category:    "cache",
		RPC:         "SetDiskCacheEvictionStrategy",
		Description: "Set disk cache eviction strategy",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cache.set-folder",
		Category:    "cache",
		RPC:         "SetFolderDiskCache",
		Description: "Set folder disk cache configuration",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cache.remove-folder",
		Category:    "cache",
		RPC:         "RemoveFolderDiskCache",
		Description: "Remove folder disk cache",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cache.prefetch",
		Category:    "cache",
		RPC:         "PrefetchFileRanges",
		Description: "Prefetch file ranges",
		RiskLevel:   registry.RiskMedium,
		DefaultOpen: registry.BoolPtr(false),
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cache.prefetch-start",
		Category:    "cache",
		RPC:         "StartFilePrefetch",
		Description: "Start file prefetch",
		RiskLevel:   registry.RiskMedium,
		DefaultOpen: registry.BoolPtr(false),
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cache.cancel-prefetch",
		Category:    "cache",
		RPC:         "CancelFilePrefetch",
		Description: "Cancel file prefetch",
		RiskLevel:   registry.RiskMedium,
		DefaultOpen: registry.BoolPtr(false),
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cache.close-reader",
		Category:    "cache",
		RPC:         "CloseFileReader",
		Description: "Close file reader",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cache.prefetch-hints",
		Category:    "cache",
		RPC:         "GetActivePrefetchHints",
		Description: "Get active prefetch hints",
		RiskLevel:   registry.RiskLow,
	})
}

func registerPromotionCommands() {
	registry.MustRegister(&registry.CommandSpec{
		ID:          "promotion.list",
		Category:    "promotion",
		RPC:         "GetPromotions",
		Description: "List promotions",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "promotion.list-cloud",
		Category:    "promotion",
		RPC:         "GetPromotionsByCloud",
		Description: "List promotions by cloud",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "promotion.update",
		Category:    "promotion",
		RPC:         "UpdatePromotionResult",
		Description: "Update promotion result",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "promotion.update-cloud",
		Category:    "promotion",
		RPC:         "UpdatePromotionResultByCloud",
		Description: "Update promotion result by cloud",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "promotion.plan-list",
		Category:    "promotion",
		RPC:         "GetCloudDrivePlans",
		Description: "List CloudDrive plans",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "promotion.referral-code",
		Category:    "promotion",
		RPC:         "GetReferralCode",
		Description: "Get referral code",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "promotion.balance",
		Category:    "promotion",
		RPC:         "GetBalanceLog",
		Description: "Get balance",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "promotion.join-plan",
		Category:    "promotion",
		RPC:         "JoinPlan",
		Description: "Join a plan",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "promotion.bind-cloud-account",
		Category:    "promotion",
		RPC:         "BindCloudAccount",
		Description: "Bind cloud account to promotion",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "promotion.transfer-balance",
		Category:    "promotion",
		RPC:         "TransferBalance",
		Description: "Transfer balance to another user",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "promotion.activate-plan",
		Category:    "promotion",
		RPC:         "ActivatePlan",
		Description: "Activate a plan with code",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "promotion.send-action",
		Category:    "promotion",
		RPC:         "SendPromotionAction",
		Description: "Send promotion action",
		RiskLevel:   registry.RiskHigh,
	})
}

func registerRemoteUploadCommands() {
	registry.MustRegister(&registry.CommandSpec{
		ID:          "remote-upload.start",
		Category:    "remote-upload",
		RPC:         "StartRemoteUpload",
		Description: "Start a remote upload",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "remote-upload.control",
		Category:    "remote-upload",
		RPC:         "RemoteUploadControl",
		Description: "Control a remote upload (cancel/pause/resume)",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "remote-upload.read-data",
		Category:    "remote-upload",
		RPC:         "RemoteReadData",
		Description: "Send read data for remote upload",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "remote-upload.hash-progress",
		Category:    "remote-upload",
		RPC:         "RemoteHashProgress",
		Description: "Send hash progress for remote upload",
		RiskLevel:   registry.RiskLow,
	})
}

func registerCloudapiCommands() {
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cloudapi.list",
		Category:    "cloudapi",
		RPC:         "GetAllCloudApis",
		Description: "List all cloud APIs",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cloudapi.can-add",
		Category:    "cloudapi",
		RPC:         "CanAddMoreCloudApis",
		Description: "Check if can add more cloud APIs",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cloudapi.config",
		Category:    "cloudapi",
		RPC:         "GetCloudAPIConfig",
		Description: "Get cloud API configuration",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cloudapi.remove",
		Category:    "cloudapi",
		RPC:         "RemoveCloudAPI",
		Description: "Remove a cloud API",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cloudapi.set-config",
		Category:    "cloudapi",
		RPC:         "SetCloudAPIConfig",
		Description: "Set cloud API configuration",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cloudapi.discover-smb-servers",
		Category:    "cloudapi",
		RPC:         "DiscoverSmbServers",
		Description: "Discover SMB servers",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cloudapi.discover-smb-shares",
		Category:    "cloudapi",
		RPC:         "DiscoverSmbShares",
		Description: "Discover SMB shares on a server",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cloudapi.login-webdav",
		Category:    "cloudapi",
		RPC:         "APILoginWebDav",
		Description: "Login to WebDAV",
		RiskLevel:   registry.RiskHigh,
		DefaultOpen: registry.BoolPtr(false),
		AliasGroup:  "storage.add.webdav",
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cloudapi.login-s3",
		Category:    "cloudapi",
		RPC:         "APILoginS3",
		Description: "Login to S3",
		RiskLevel:   registry.RiskHigh,
		DefaultOpen: registry.BoolPtr(false),
		AliasGroup:  "storage.add.s3",
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cloudapi.login-local",
		Category:    "cloudapi",
		RPC:         "APIAddLocalFolder",
		Description: "Add local folder",
		RiskLevel:   registry.RiskHigh,
		DefaultOpen: registry.BoolPtr(false),
		AliasGroup:  "storage.add.local",
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cloudapi.login-clouddrive",
		Category:    "cloudapi",
		RPC:         "APILoginCloudDrive",
		Description: "Login to CloudDrive",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cloudapi.login-sftp",
		Category:    "cloudapi",
		RPC:         "APILoginSftp",
		Description: "Login to SFTP",
		RiskLevel:   registry.RiskHigh,
		DefaultOpen: registry.BoolPtr(false),
		AliasGroup:  "storage.add.sftp",
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cloudapi.login-ftp",
		Category:    "cloudapi",
		RPC:         "APILoginFtp",
		Description: "Login to FTP",
		RiskLevel:   registry.RiskHigh,
		DefaultOpen: registry.BoolPtr(false),
		AliasGroup:  "storage.add.ftp",
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cloudapi.login-smb",
		Category:    "cloudapi",
		RPC:         "APILoginSmb",
		Description: "Login to SMB",
		RiskLevel:   registry.RiskHigh,
		DefaultOpen: registry.BoolPtr(false),
		AliasGroup:  "storage.add.smb",
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cloudapi.login-onedrive-oauth",
		Category:    "cloudapi",
		RPC:         "APILoginOneDriveOAuth",
		Description: "Login to OneDrive with OAuth",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cloudapi.login-google-oauth",
		Category:    "cloudapi",
		RPC:         "ApiLoginGoogleDriveOAuth",
		Description: "Login to Google Drive with OAuth",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cloudapi.login-google-refresh",
		Category:    "cloudapi",
		RPC:         "ApiLoginGoogleDriveRefreshToken",
		Description: "Login to Google Drive with refresh token",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cloudapi.login-xunlei-oauth",
		Category:    "cloudapi",
		RPC:         "ApiLoginXunleiOAuth",
		Description: "Login to Xunlei with OAuth",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cloudapi.login-xunleiopen-oauth",
		Category:    "cloudapi",
		RPC:         "ApiLoginXunleiOpenOAuth",
		Description: "Login to Xunlei Open with OAuth",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cloudapi.login-123pan-oauth",
		Category:    "cloudapi",
		RPC:         "ApiLogin123PanOAuth",
		Description: "Login to 123Pan with OAuth",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cloudapi.login-115-cookie",
		Category:    "cloudapi",
		RPC:         "APILogin115Editthiscookie",
		Description: "Login to 115 with editthiscookie string",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cloudapi.login-115-qrcode",
		Category:    "cloudapi",
		RPC:         "APILogin115QRCode",
		Description: "Login to 115 with QR code",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cloudapi.login-115open-oauth",
		Category:    "cloudapi",
		RPC:         "APILogin115OpenOAuth",
		Description: "Login to 115 Open with OAuth",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cloudapi.login-115open-qrcode",
		Category:    "cloudapi",
		RPC:         "APILogin115OpenQRCode",
		Description: "Login to 115 Open with QR code",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cloudapi.login-aliyun-oauth",
		Category:    "cloudapi",
		RPC:         "APILoginAliyundriveOAuth",
		Description: "Login to Aliyun with OAuth",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cloudapi.login-aliyun-refresh",
		Category:    "cloudapi",
		RPC:         "APILoginAliyundriveRefreshtoken",
		Description: "Login to Aliyun with refresh token",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cloudapi.login-aliyun-qrcode",
		Category:    "cloudapi",
		RPC:         "APILoginAliyunDriveQRCode",
		Description: "Login to Aliyun with QR code",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cloudapi.login-baidu-oauth",
		Category:    "cloudapi",
		RPC:         "APILoginBaiduPanOAuth",
		Description: "Login to Baidu with OAuth",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "cloudapi.login-189-qrcode",
		Category:    "cloudapi",
		RPC:         "APILogin189QRCode",
		Description: "Login to 189 with QR code",
		RiskLevel:   registry.RiskHigh,
	})
}

func registerWhitelistCommands() {
	registry.MustRegister(&registry.CommandSpec{
		ID:          "whitelist.list",
		Category:    "whitelist",
		RPC:         "",
		Description: "List whitelisted commands",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "whitelist.status",
		Category:    "whitelist",
		RPC:         "",
		Description: "Get status of a command",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "whitelist.enable",
		Category:    "whitelist",
		RPC:         "",
		Description: "Enable a command",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "whitelist.disable",
		Category:    "whitelist",
		RPC:         "",
		Description: "Disable a command",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "whitelist.reset",
		Category:    "whitelist",
		RPC:         "",
		Description: "Reset whitelist to defaults",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "whitelist.path",
		Category:    "whitelist",
		RPC:         "",
		Description: "Show whitelist config path",
		RiskLevel:   registry.RiskLow,
	})
}

func registerTaskCommands() {
	registry.MustRegister(&registry.CommandSpec{
		ID:          "task.list",
		Category:    "task",
		RPC:         "GetAllTasksCount,GetUploadFileList,GetDownloadFileList,GetCopyTasks,GetMergeTasks",
		Description: "List all tasks",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "task.list.upload",
		Category:    "task",
		RPC:         "GetUploadFileList",
		Description: "List upload tasks",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "task.list.download",
		Category:    "task",
		RPC:         "GetDownloadFileList",
		Description: "List download tasks",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "task.list.copy",
		Category:    "task",
		RPC:         "GetCopyTasks",
		Description: "List copy tasks",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "task.list.merge",
		Category:    "task",
		RPC:         "GetMergeTasks",
		Description: "List merge tasks",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "task.status",
		Category:    "task",
		RPC:         "GetAllTasksCount",
		Description: "Get task status summary",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "task.cancel.upload",
		Category:    "task",
		RPC:         "CancelUploadFiles",
		Description: "Cancel upload tasks",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "task.cancel.copy",
		Category:    "task",
		RPC:         "CancelCopyTask",
		Description: "Cancel copy task",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "task.cancel.merge",
		Category:    "task",
		RPC:         "CancelMergeTask",
		Description: "Cancel merge task",
		RiskLevel:   registry.RiskHigh,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "task.wait.copy",
		Category:    "task",
		RPC:         "GetCopyTasks",
		Description: "Wait for copy task completion",
		RiskLevel:   registry.RiskLow,
	})
	registry.MustRegister(&registry.CommandSpec{
		ID:          "task.wait.merge",
		Category:    "task",
		RPC:         "GetMergeTasks",
		Description: "Wait for merge task completion",
		RiskLevel:   registry.RiskLow,
	})
}
