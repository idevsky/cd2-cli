package client

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{"empty", Config{Address: ""}, true},
		{"valid", Config{Address: "localhost:19798"}, false},
		{"tls", Config{Address: "localhost:19798", UseTLS: true, SkipVerifyTLS: true}, false},
		{"timeout", Config{Address: "localhost:19798", Timeout: 5 * time.Second}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewClient(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if c != nil {
				c.Close()
			}
		})
	}
}

func TestClientToken(t *testing.T) {
	c, _ := NewClient(Config{Address: "localhost:19798"})
	defer c.Close()
	c.SetToken("test")
	if c.GetToken() != "test" {
		t.Error("token mismatch")
	}
}

func TestClientAccessors(t *testing.T) {
	c, _ := NewClient(Config{Address: "localhost:19798"})
	defer c.Close()
	apis := []interface{}{
		c.Public(), c.Auth(), c.File(), c.Mount(), c.Transfer(),
		c.CloudAPI(), c.Backup(), c.WebDAV(), c.Token(), c.Session(),
		c.System(), c.Offline(), c.Webhook(), c.Local(), c.RemoteUpload(),
		c.Copy(), c.Stream(), c.Sync(), c.Promotion(),
	}
	for i, api := range apis {
		if api == nil {
			t.Errorf("API %d is nil", i)
		}
	}
}

func TestWithTimeout(t *testing.T) {
	c, _ := NewClient(Config{Address: "localhost:19798"})
	defer c.Close()
	ctx, cancel := c.withTimeout(context.Background(), 0)
	defer cancel()
	_, ok := ctx.Deadline()
	if !ok {
		t.Error("missing deadline")
	}
}

func TestWithTimeoutUsesConfigDefault(t *testing.T) {
	c, _ := NewClient(Config{Address: "localhost:19798", Timeout: 5 * time.Second})
	defer c.Close()

	start := time.Now()
	ctx, cancel := c.withTimeout(context.Background(), 0)
	defer cancel()

	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatal("missing deadline")
	}

	remaining := deadline.Sub(start)
	if remaining < 4*time.Second || remaining > 6*time.Second {
		t.Fatalf("deadline = %v from start, want around 5s", remaining)
	}
}

func TestWithAuth(t *testing.T) {
	c, _ := NewClient(Config{Address: "localhost:19798"})
	defer c.Close()
	ctx1 := c.withAuth(context.Background())
	if ctx1 != context.Background() {
		t.Error("without token should return same ctx")
	}
	c.SetToken("test")
	ctx2 := c.withAuth(context.Background())
	if ctx2 == context.Background() {
		t.Error("with token should return different ctx")
	}
}

func TestWithAuthPreservesMetadata(t *testing.T) {
	c, _ := NewClient(Config{Address: "localhost:19798", Token: "test"})
	defer c.Close()

	ctx := metadata.AppendToOutgoingContext(context.Background(), "x-request-id", "abc")
	ctx = c.withAuth(ctx)

	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		t.Fatal("missing outgoing metadata")
	}
	if got := md.Get("x-request-id"); len(got) != 1 || got[0] != "abc" {
		t.Fatalf("x-request-id metadata = %v, want [abc]", got)
	}
	if got := md.Get("authorization"); len(got) != 1 || got[0] != "Bearer test" {
		t.Fatalf("authorization metadata = %v, want bearer token", got)
	}
}

func TestClose(t *testing.T) {
	c, _ := NewClient(Config{Address: "localhost:19798"})
	if err := c.Close(); err != nil {
		t.Error("Close error:", err)
	}
	c2 := &Client{}
	if err := c2.Close(); err != nil {
		t.Error("Close nil conn error:", err)
	}
}

func TestAPICalls(t *testing.T) {
	c, _ := NewClient(Config{Address: "localhost:19798"})
	defer c.Close()
	c.SetToken("test")
	ctx := context.Background()

	_, _ = c.Public().GetSystemInfo(ctx)
	_, _ = c.Public().GetToken(ctx, &pb.GetTokenRequest{})
	_, _ = c.Public().Login(ctx, &pb.UserLoginRequest{})
	_, _ = c.Public().LoginWithThirdPartyAccount(ctx, &pb.LoginWithThirdPartyAccountRequest{})
	_, _ = c.Public().Register(ctx, &pb.UserRegisterRequest{})
	_ = c.Public().SendResetAccountEmail(ctx, &pb.SendResetAccountEmailRequest{})
	_ = c.Public().ResetAccount(ctx, &pb.ResetAccountRequest{})
	_, _ = c.Public().GetApiTokenInfo(ctx, "test")
	_, _ = c.Public().LoginWith2FA(ctx, &pb.LoginWith2FARequest{})

	_ = c.Auth().SendConfirmEmail(ctx)
	_ = c.Auth().ConfirmEmail(ctx, &pb.ConfirmEmailRequest{})
	_, _ = c.Auth().GetAccountStatus(ctx)
	_, _ = c.Auth().Logout(ctx, &pb.UserLogoutRequest{})
	_, _ = c.Auth().Check2FAStatus(ctx)
	_, _ = c.Auth().Setup2FA(ctx, &pb.Setup2FARequest{})
	_, _ = c.Auth().Enable2FA(ctx, &pb.TwoFactorAuthCodeRequest{})
	_, _ = c.Auth().Disable2FA(ctx, &pb.TwoFactorAuthCodeRequest{})
	_, _ = c.Auth().GetRecoveryCodes(ctx, &pb.TwoFactorAuthCodeRequest{})
	_, _ = c.Auth().RegenerateRecoveryCodes(ctx, &pb.TwoFactorAuthCodeRequest{})
	_, _ = c.Auth().ChangePassword(ctx, &pb.ChangePasswordRequest{})
	_ = c.Auth().SendChangeEmailCode(ctx, &pb.SendChangeEmailCodeRequest{})
	_ = c.Auth().ChangeEmail(ctx, &pb.ChangeEmailRequest{})
	_ = c.Auth().ChangeEmailAndPassword(ctx, &pb.ChangeEmailAndPasswordRequest{})

	_, _ = c.File().GetSubFiles(ctx, &pb.ListSubFileRequest{})
	_, _ = c.File().GetSearchResults(ctx, &pb.SearchRequest{})
	_, _ = c.File().FindFileByPath(ctx, "/test")
	_, _ = c.File().CreateFolder(ctx, &pb.CreateFolderRequest{})
	_, _ = c.File().CreateEncryptedFolder(ctx, &pb.CreateEncryptedFolderRequest{})
	_, _ = c.File().UnlockEncryptedFile(ctx, &pb.UnlockEncryptedFileRequest{})
	_, _ = c.File().LockEncryptedFile(ctx, "/test")
	_, _ = c.File().RenameFile(ctx, &pb.RenameFileRequest{})
	_, _ = c.File().RenameFiles(ctx, &pb.RenameFilesRequest{})
	_, _ = c.File().MoveFile(ctx, &pb.MoveFileRequest{})
	_, _ = c.File().CopyFile(ctx, &pb.CopyFileRequest{})
	_, _ = c.File().DeleteFile(ctx, "/test")
	_, _ = c.File().DeleteFilePermanently(ctx, "/test")
	_, _ = c.File().DeleteFiles(ctx, []string{})
	_, _ = c.File().DeleteFilesPermanently(ctx, []string{})
	_ = c.File().AddSharedLink(ctx, &pb.AddSharedLinkRequest{})
	_, _ = c.File().GetFileDetailProperties(ctx, "/test")
	_, _ = c.File().GetSpaceInfo(ctx, "/test")
	_, _ = c.File().GetCloudMemberships(ctx, "/test")
	_, _ = c.File().GetMetaData(ctx, "/test")
	_, _ = c.File().GetOriginalPath(ctx, "/test")
	_, _ = c.File().CreateFile(ctx, &pb.CreateFileRequest{})
	_, _ = c.File().CloseFile(ctx, &pb.CloseFileRequest{})
	_, _ = c.File().WriteToFile(ctx, &pb.WriteFileRequest{})
	_, _ = c.File().WriteToFileStream(ctx)
	_, _ = c.File().GetDownloadUrl(ctx, &pb.GetDownloadUrlPathRequest{})

	_, _ = c.Mount().CanAddMoreMountPoints(ctx)
	_, _ = c.Mount().GetMountPoints(ctx)
	_, _ = c.Mount().AddMountPoint(ctx, &pb.MountOption{})
	_, _ = c.Mount().RemoveMountPoint(ctx, &pb.MountPointRequest{})
	_, _ = c.Mount().Mount(ctx, &pb.MountPointRequest{})
	_, _ = c.Mount().Unmount(ctx, &pb.MountPointRequest{})
	_, _ = c.Mount().UpdateMountPoint(ctx, &pb.UpdateMountPointRequest{})
	_, _ = c.Mount().GetAvailableDriveLetters(ctx)
	_, _ = c.Mount().HasDriveLetters(ctx)
	_, _ = c.Mount().CanMountBothLocalAndCloud(ctx)

	_, _ = c.Transfer().GetAllTasksCount(ctx)
	_, _ = c.Transfer().GetDownloadFileCount(ctx)
	_, _ = c.Transfer().GetDownloadFileList(ctx)
	_, _ = c.Transfer().GetUploadFileCount(ctx)
	_, _ = c.Transfer().GetUploadFileList(ctx, &pb.GetUploadFileListRequest{})
	_ = c.Transfer().CancelAllUploadFiles(ctx)
	_ = c.Transfer().CancelUploadFiles(ctx, []string{})
	_ = c.Transfer().PauseAllUploadFiles(ctx)
	_ = c.Transfer().PauseUploadFiles(ctx, []string{})
	_ = c.Transfer().ResumeAllUploadFiles(ctx)
	_ = c.Transfer().ResumeUploadFiles(ctx, []string{})

	_, _ = c.CloudAPI().CanAddMoreCloudApis(ctx)
	_, _ = c.CloudAPI().APILogin115Editthiscookie(ctx, &pb.Login115EditthiscookieRequest{})
	_, _ = c.CloudAPI().APILogin115QRCode(ctx, &pb.Login115QrCodeRequest{})
	_, _ = c.CloudAPI().APILogin115OpenOAuth(ctx, &pb.Login115OpenOAuthRequest{})
	_, _ = c.CloudAPI().APILogin115OpenQRCode(ctx, &pb.Login115OpenQRCodeRequest{})
	_, _ = c.CloudAPI().APILoginAliyundriveOAuth(ctx, &pb.LoginAliyundriveOAuthRequest{})
	_, _ = c.CloudAPI().APILoginAliyundriveRefreshtoken(ctx, &pb.LoginAliyundriveRefreshtokenRequest{})
	_, _ = c.CloudAPI().APILoginAliyunDriveQRCode(ctx, &pb.LoginAliyundriveQRCodeRequest{})
	_, _ = c.CloudAPI().APILoginBaiduPanOAuth(ctx, &pb.LoginBaiduPanOAuthRequest{})
	_, _ = c.CloudAPI().APILoginOneDriveOAuth(ctx, &pb.LoginOneDriveOAuthRequest{})
	_, _ = c.CloudAPI().ApiLoginGoogleDriveOAuth(ctx, &pb.LoginGoogleDriveOAuthRequest{})
	_, _ = c.CloudAPI().ApiLoginGoogleDriveRefreshToken(ctx, &pb.LoginGoogleDriveRefreshTokenRequest{})
	_, _ = c.CloudAPI().ApiLoginXunleiOAuth(ctx, &pb.LoginXunleiOAuthRequest{})
	_, _ = c.CloudAPI().ApiLoginXunleiOpenOAuth(ctx, &pb.LoginXunleiOpenOAuthRequest{})
	_, _ = c.CloudAPI().ApiLogin123PanOAuth(ctx, &pb.Login123PanOAuthRequest{})
	_, _ = c.CloudAPI().APILogin189QRCode(ctx, &pb.Login189QRCodeRequest{})
	_, _ = c.CloudAPI().APILoginWebDav(ctx, &pb.LoginWebDavRequest{})
	_, _ = c.CloudAPI().APILoginS3(ctx, &pb.LoginS3Request{})
	_, _ = c.CloudAPI().APIAddLocalFolder(ctx, &pb.AddLocalFolderRequest{})
	_, _ = c.CloudAPI().APILoginCloudDrive(ctx, &pb.LoginCloudDriveRequest{})
	_, _ = c.CloudAPI().APILoginSftp(ctx, &pb.LoginSftpRequest{})
	_, _ = c.CloudAPI().APILoginFtp(ctx, &pb.LoginFtpRequest{})
	_, _ = c.CloudAPI().APILoginSmb(ctx, &pb.LoginSmbRequest{})
	_, _ = c.CloudAPI().DiscoverSmbServers(ctx)
	_, _ = c.CloudAPI().DiscoverSmbShares(ctx, &pb.DiscoverSmbSharesRequest{})
	_, _ = c.CloudAPI().RemoveCloudAPI(ctx, &pb.RemoveCloudAPIRequest{})
	_, _ = c.CloudAPI().GetAllCloudApis(ctx)
	_, _ = c.CloudAPI().GetCloudAPIConfig(ctx, &pb.GetCloudAPIConfigRequest{})
	_ = c.CloudAPI().SetCloudAPIConfig(ctx, &pb.SetCloudAPIConfigRequest{})

	_, _ = c.Backup().BackupGetAll(ctx)
	_, _ = c.Backup().BackupGetStatus(ctx, "test")
	_ = c.Backup().BackupAdd(ctx, &pb.Backup{})
	_ = c.Backup().BackupRemove(ctx, "test")
	_ = c.Backup().BackupUpdate(ctx, &pb.Backup{})
	_ = c.Backup().BackupAddDestination(ctx, &pb.BackupModifyRequest{})
	_ = c.Backup().BackupRemoveDestination(ctx, &pb.BackupModifyRequest{})
	_ = c.Backup().BackupSetEnabled(ctx, &pb.BackupSetEnabledRequest{})
	_ = c.Backup().BackupSetFileSystemWatchEnabled(ctx, &pb.BackupModifyRequest{})
	_ = c.Backup().BackupUpdateStrategies(ctx, &pb.BackupModifyRequest{})
	_ = c.Backup().BackupRestartWalkingThrough(ctx, "test")
	_, _ = c.Backup().CanAddMoreBackups(ctx)
	_ = c.Backup().NotifyPhotoLibraryChanges(ctx, &pb.PhotoLibraryChangeList{})

	_ = c.WebDAV().AddDavUser(ctx, &pb.AddDavUserRequest{})
	_ = c.WebDAV().RemoveDavUser(ctx, "test")
	_ = c.WebDAV().ModifyDavUser(ctx, &pb.ModifyDavUserRequest{})
	_, _ = c.WebDAV().GetDavUser(ctx, "test")
	_, _ = c.WebDAV().GetDavServerConfig(ctx)
	_ = c.WebDAV().SetDavServerConfig(ctx, &pb.ModifyDavServerConfigRequest{})

	_, _ = c.Token().CreateToken(ctx, &pb.CreateTokenRequest{})
	_, _ = c.Token().ModifyToken(ctx, &pb.ModifyTokenRequest{})
	_ = c.Token().RemoveToken(ctx, "test")
	_, _ = c.Token().ListTokens(ctx)

	_, _ = c.Session().GetSessions(ctx)
	_ = c.Session().RevokeSession(ctx, "test")
	_ = c.Session().RevokeOtherSessions(ctx)

	_, _ = c.System().GetRuntimeInfo(ctx)
	_, _ = c.System().GetRunningInfo(ctx)
	_, _ = c.System().GetOpenFileHandles(ctx)
	_, _ = c.System().GetFileBufferDiskCacheStats(ctx)
	_ = c.System().PurgeFileBufferDiskCache(ctx)
	_ = c.System().SetDiskCacheEvictionStrategy(ctx, &pb.SetDiskCacheEvictionStrategyRequest{})
	_ = c.System().SetFolderDiskCache(ctx, &pb.SetFolderDiskCacheRequest{})
	_ = c.System().RemoveFolderDiskCache(ctx, "/test")
	_, _ = c.System().ListDiskCacheFolders(ctx)
	_, _ = c.System().PrefetchFileRanges(ctx, &pb.PrefetchFileRangesRequest{})
	_ = c.System().CancelFilePrefetch(ctx, &pb.CancelFilePrefetchRequest{})
	_ = c.System().CloseFileReader(ctx, "/test")
	_, _ = c.System().GetActivePrefetchHints(ctx)
	_, _ = c.System().GetSystemSettings(ctx)
	_ = c.System().SetSystemSettings(ctx, &pb.SystemSettings{})
	_ = c.System().SetDirCacheTimeSecs(ctx, &pb.SetDirCacheTimeRequest{})
	_, _ = c.System().GetEffectiveDirCacheTimeSecs(ctx, &pb.GetEffectiveDirCacheTimeRequest{})
	_ = c.System().ForceExpireDirCache(ctx, "/test")
	_ = c.System().VacuumDirCache(ctx)
	_, _ = c.System().GetVacuumProgress(ctx)
	_, _ = c.System().GetDirCacheDbSize(ctx)
	_, _ = c.System().GetOpenFileTable(ctx, &pb.GetOpenFileTableRequest{})
	_, _ = c.System().GetDirCacheTable(ctx)
	_, _ = c.System().GetReferencedEntryPaths(ctx, "/test")
	_, _ = c.System().GetTempFileTable(ctx)
	_, _ = c.System().GetServiceCapabilities(ctx)
	_ = c.System().RestartService(ctx)
	_ = c.System().ShutdownService(ctx)
	_, _ = c.System().HasUpdate(ctx)
	_, _ = c.System().CheckUpdate(ctx)
	_ = c.System().DownloadUpdate(ctx)
	_ = c.System().UpdateSystem(ctx)
	_, _ = c.System().GetWebServerConfig(ctx)
	_ = c.System().SetWebServerConfig(ctx, &pb.SetWebServerConfigRequest{})
	_ = c.System().GenerateSelfSignedCert(ctx, &pb.GenerateSelfSignedCertRequest{})
	_, _ = c.System().GetMachineId(ctx)
	_, _ = c.System().GetOnlineDevices(ctx)
	_ = c.System().KickoutDevice(ctx, &pb.DeviceRequest{})
	_, _ = c.System().ListLogFiles(ctx)
	_ = c.System().TestUpdate(ctx, &pb.FileRequest{})

	_, _ = c.Offline().AddOfflineFiles(ctx, &pb.AddOfflineFileRequest{})
	_, _ = c.Offline().RemoveOfflineFiles(ctx, &pb.RemoveOfflineFilesRequest{})
	_, _ = c.Offline().ListOfflineFilesByPath(ctx, "/test")
	_, _ = c.Offline().ListAllOfflineFiles(ctx, &pb.OfflineFileListAllRequest{})
	_, _ = c.Offline().GetOfflineQuotaInfo(ctx, &pb.OfflineQuotaRequest{})
	_ = c.Offline().ClearOfflineFiles(ctx, &pb.ClearOfflineFileRequest{})
	_ = c.Offline().RestartOfflineTask(ctx, &pb.RestartOfflineFileRequest{})

	_, _ = c.Webhook().GetWebhookConfigTemplate(ctx)
	_, _ = c.Webhook().GetWebhookConfigs(ctx)
	_ = c.Webhook().AddWebhookConfig(ctx, &pb.WebhookRequest{})
	_ = c.Webhook().RemoveWebhookConfig(ctx, "test")
	_ = c.Webhook().ChangeWebhookConfig(ctx, &pb.WebhookRequest{})

	_, _ = c.Local().LocalGetSubFiles(ctx, &pb.LocalGetSubFilesRequest{})
	_, _ = c.Local().LocalCreateFolder(ctx, &pb.LocalCreateFolderRequest{})

	_, _ = c.RemoteUpload().StartRemoteUpload(ctx, &pb.StartRemoteUploadRequest{})
	_ = c.RemoteUpload().RemoteUploadControl(ctx, &pb.RemoteUploadControlRequest{})
	_, _ = c.RemoteUpload().RemoteReadData(ctx, &pb.RemoteReadDataUpload{})
	_, _ = c.RemoteUpload().RemoteHashProgress(ctx, &pb.RemoteHashProgressUpload{})

	_, _ = c.Copy().GetCopyTasks(ctx)
	_, _ = c.Copy().GetMergeTasks(ctx)
	_ = c.Copy().CancelMergeTask(ctx, &pb.CancelMergeTaskRequest{})
	_ = c.Copy().CancelCopyTask(ctx, &pb.CopyTaskRequest{})
	_ = c.Copy().PauseCopyTask(ctx, &pb.PauseCopyTaskRequest{})
	_ = c.Copy().RestartCopyTask(ctx, &pb.CopyTaskRequest{})
	_ = c.Copy().RemoveCompletedCopyTasks(ctx)
	_, _ = c.Copy().RemoveAllCopyTasks(ctx)
	_, _ = c.Copy().RemoveCopyTasks(ctx, &pb.CopyTaskBatchRequest{})
	_, _ = c.Copy().PauseAllCopyTasks(ctx, &pb.PauseAllCopyTasksRequest{})
	_, _ = c.Copy().PauseCopyTasks(ctx, &pb.PauseCopyTasksRequest{})
	_, _ = c.Copy().ResumeAllCopyTasks(ctx)
	_, _ = c.Copy().ResumeCopyTasks(ctx, &pb.CopyTaskBatchRequest{})

	_, _ = c.Stream().PushTaskChange(ctx)
	_, _ = c.Stream().PushMessage(ctx)
	_, _ = c.Stream().RemoteUploadChannel(ctx, &pb.RemoteUploadChannelRequest{})

	_, _ = c.Sync().SyncFileChangesFromCloud(ctx, "/test")
	_ = c.Sync().StartCloudEventListener(ctx, "/test")
	_ = c.Sync().StopCloudEventListener(ctx, "/test")
	_, _ = c.Sync().WalkThroughFolderTest(ctx, "/test")
	_, _ = c.Sync().GetCloudDrive1UserData(ctx)

	_, _ = c.Promotion().GetPromotions(ctx)
	_, _ = c.Promotion().GetPromotionsByCloud(ctx, &pb.CloudAPIRequest{})
	_ = c.Promotion().UpdatePromotionResult(ctx)
	_ = c.Promotion().UpdatePromotionResultByCloud(ctx, &pb.UpdatePromotionResultByCloudRequest{})
	_ = c.Promotion().SendPromotionAction(ctx, &pb.SendPromotionActionRequest{})
	_, _ = c.Promotion().GetCloudDrivePlans(ctx)
	_, _ = c.Promotion().JoinPlan(ctx, &pb.JoinPlanRequest{})
	_ = c.Promotion().BindCloudAccount(ctx, &pb.BindCloudAccountRequest{})
	_ = c.Promotion().TransferBalance(ctx, &pb.TransferBalanceRequest{})
	_, _ = c.Promotion().GetBalanceLog(ctx)
	_, _ = c.Promotion().CheckActivationCode(ctx, "test")
	_, _ = c.Promotion().ActivatePlan(ctx, "test")
	_, _ = c.Promotion().CheckCouponCode(ctx, &pb.CheckCouponCodeRequest{})
	_, _ = c.Promotion().GetReferralCode(ctx)
}
