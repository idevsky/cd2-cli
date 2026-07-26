//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

func TestIntegration2FA_Status(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	status, err := c.Auth().Check2FAStatus(ctx)
	if err != nil {
		t.Fatalf("Check2FAStatus failed: %v", err)
	}

	t.Logf("2FA enabled: %v", status.TwoFactorEnabled)
}

func TestIntegrationBackup_GetAll(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	backups, err := c.Backup().BackupGetAll(ctx)
	if err != nil {
		t.Fatalf("BackupGetAll failed: %v", err)
	}

	t.Logf("Backups: %d", len(backups.Backups))
}

func TestIntegrationBackup_CanAddMore(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := c.Backup().CanAddMoreBackups(ctx)
	if err != nil {
		t.Fatalf("CanAddMoreBackups failed: %v", err)
	}

	t.Logf("Can add more backups: %v", result.Success)
}

func TestIntegrationSession_List(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sessions, err := c.Session().GetSessions(ctx)
	if err != nil {
		t.Fatalf("GetSessions failed: %v", err)
	}

	t.Logf("Sessions: %d", len(sessions.Sessions))
}

func TestIntegrationToken_List(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tokens, err := c.Token().ListTokens(ctx)
	if err != nil {
		t.Fatalf("ListTokens failed: %v", err)
	}

	t.Logf("Tokens: %d", len(tokens.Tokens))
}

func TestIntegrationSystem_ListDiskCacheFolders(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	folders, err := c.System().ListDiskCacheFolders(ctx)
	if err != nil {
		t.Fatalf("ListDiskCacheFolders failed: %v", err)
	}

	t.Logf("Disk cache folders: %d", len(folders.Folders))
}

func TestIntegrationSystem_GetDirCacheDbSize(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	size, err := c.System().GetDirCacheDbSize(ctx)
	if err != nil {
		t.Fatalf("GetDirCacheDbSize failed: %v", err)
	}

	t.Logf("Dir cache DB: %d bytes, vacuuming: %v", size.TotalSizeBytes, size.IsVacuuming)
}

func TestIntegrationCopy_GetTasks(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tasks, err := c.Copy().GetCopyTasks(ctx)
	if err != nil {
		t.Fatalf("GetCopyTasks failed: %v", err)
	}

	t.Logf("Copy tasks: %d", len(tasks.CopyTasks))
}

func TestIntegrationCopy_GetMergeTasks(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tasks, err := c.Copy().GetMergeTasks(ctx)
	if err != nil {
		t.Fatalf("GetMergeTasks failed: %v", err)
	}

	t.Logf("Merge tasks: %d", len(tasks.MergeTasks))
}

func TestIntegrationLocal_GetSubFiles(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	files, err := c.Local().LocalGetSubFiles(ctx, &pb.LocalGetSubFilesRequest{
		ParentFolder:      "/",
		IncludeCloudDrive: true,
	})
	if err != nil {
		t.Fatalf("LocalGetSubFiles failed: %v", err)
	}

	t.Logf("Local files: %d", len(files))
}

func TestIntegrationOffline_ListAll(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	files, err := c.Offline().ListAllOfflineFiles(ctx, &pb.OfflineFileListAllRequest{})
	if err != nil {
		t.Fatalf("ListAllOfflineFiles failed: %v", err)
	}

	t.Logf("Offline files: %d", len(files.OfflineFiles))
}

func TestIntegrationWebhook_GetConfigs(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	hooks, err := c.Webhook().GetWebhookConfigs(ctx)
	if err != nil {
		t.Fatalf("GetWebhookConfigs failed: %v", err)
	}

	t.Logf("Webhooks: %d", len(hooks.Webhooks))
}

func TestIntegrationSystem_GetOpenFileHandles(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	handles, err := c.System().GetOpenFileHandles(ctx)
	if err != nil {
		t.Fatalf("GetOpenFileHandles failed: %v", err)
	}

	t.Logf("Open file handles: %v", handles)
}

func TestIntegrationSystem_GetTempFileTable(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	table, err := c.System().GetTempFileTable(ctx)
	if err != nil {
		t.Fatalf("GetTempFileTable failed: %v", err)
	}

	t.Logf("Temp files: %v", table)
}

func TestIntegrationSystem_GetOpenFileTable(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	table, err := c.System().GetOpenFileTable(ctx, &pb.GetOpenFileTableRequest{})
	if err != nil {
		t.Fatalf("GetOpenFileTable failed: %v", err)
	}

	t.Logf("Open files: %v", table)
}

func TestIntegrationSystem_ListLogFiles(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logs, err := c.System().ListLogFiles(ctx)
	if err != nil {
		t.Fatalf("ListLogFiles failed: %v", err)
	}

	t.Logf("Log files: %d", len(logs.LogFiles))
}

func TestIntegrationSystem_GetFileBufferDiskCacheStats(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stats, err := c.System().GetFileBufferDiskCacheStats(ctx)
	if err != nil {
		t.Fatalf("GetFileBufferDiskCacheStats failed: %v", err)
	}

	t.Logf("Cache stats: %v", stats)
}
