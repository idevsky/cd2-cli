//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

func TestIntegrationBackup_List(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	backups, err := c.Backup().BackupGetAll(ctx)
	if err != nil {
		t.Fatalf("BackupGetAll failed: %v", err)
	}

	t.Logf("Backups: %d", len(backups.Backups))

	for _, backup := range backups.Backups {
		if backup.Backup != nil {
			t.Logf("Backup: source=%s, enabled=%v", backup.Backup.SourcePath, backup.Backup.IsEnabled)
		}
	}
}

func TestIntegrationBackup_CanAdd(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := c.Backup().CanAddMoreBackups(ctx)
	if err != nil {
		t.Fatalf("CanAddMoreBackups failed: %v", err)
	}

	t.Logf("Can add more backups: %v", result.Success)
	t.Logf("Result: %+v", result)
}

func TestIntegrationBackup_GetStatus(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	backups, err := c.Backup().BackupGetAll(ctx)
	if err != nil {
		t.Fatalf("BackupGetAll failed: %v", err)
	}

	if len(backups.Backups) == 0 {
		t.Log("No backups configured, skipping status test")
		return
	}

	sourcePath := backups.Backups[0].Backup.SourcePath
	status, err := c.Backup().BackupGetStatus(ctx, sourcePath)
	if err != nil {
		t.Fatalf("BackupGetStatus failed: %v", err)
	}

	t.Logf("Backup status: source=%s, status=%v", sourcePath, status.Status)
}

func TestIntegrationBackup_AddAndRemove(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	sourcePath := "/tmp/test-backup-source"
	destPath := "/tmp/test-backup-dest"

	backup := &pb.Backup{
		SourcePath: sourcePath,
		Destinations: []*pb.BackupDestination{
			{
				DestinationPath: destPath,
				IsEnabled:       true,
			},
		},
		IsEnabled:              true,
		FileSystemWatchEnabled: false,
	}

	err := c.Backup().BackupAdd(ctx, backup)
	if err != nil {
		t.Fatalf("BackupAdd failed: %v", err)
	}
	t.Logf("Backup added: source=%s", sourcePath)

	backups, err := c.Backup().BackupGetAll(ctx)
	if err != nil {
		t.Fatalf("BackupGetAll after add failed: %v", err)
	}

	found := false
	for _, b := range backups.Backups {
		if b.Backup != nil && b.Backup.SourcePath == sourcePath {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Backup with source path %s not found in list", sourcePath)
	}

	err = c.Backup().BackupRemove(ctx, sourcePath)
	if err != nil {
		t.Fatalf("BackupRemove failed: %v", err)
	}
	t.Logf("Backup removed: source=%s", sourcePath)

	backupsFinal, err := c.Backup().BackupGetAll(ctx)
	if err != nil {
		t.Fatalf("BackupGetAll after remove failed: %v", err)
	}

	for _, b := range backupsFinal.Backups {
		if b.Backup != nil && b.Backup.SourcePath == sourcePath {
			t.Errorf("Backup with source path %s still exists after remove", sourcePath)
		}
	}
}

func TestIntegrationBackup_SetEnabled(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	sourcePath := "/tmp/test-backup-enabled"

	backup := &pb.Backup{
		SourcePath: sourcePath,
		Destinations: []*pb.BackupDestination{
			{
				DestinationPath: "/tmp/test-backup-dest",
				IsEnabled:       true,
			},
		},
		IsEnabled:              true,
		FileSystemWatchEnabled: false,
	}

	err := c.Backup().BackupAdd(ctx, backup)
	if err != nil {
		t.Fatalf("BackupAdd failed: %v", err)
	}
	defer c.Backup().BackupRemove(ctx, sourcePath)

	err = c.Backup().BackupSetEnabled(ctx, &pb.BackupSetEnabledRequest{
		SourcePath: sourcePath,
		IsEnabled:  false,
	})
	if err != nil {
		t.Fatalf("BackupSetEnabled failed: %v", err)
	}
	t.Logf("Backup disabled: source=%s", sourcePath)

	status, err := c.Backup().BackupGetStatus(ctx, sourcePath)
	if err != nil {
		t.Fatalf("BackupGetStatus failed: %v", err)
	}

	if status.Backup.IsEnabled {
		t.Errorf("Expected backup to be disabled, but it is enabled")
	}

	err = c.Backup().BackupSetEnabled(ctx, &pb.BackupSetEnabledRequest{
		SourcePath: sourcePath,
		IsEnabled:  true,
	})
	if err != nil {
		t.Fatalf("BackupSetEnabled (re-enable) failed: %v", err)
	}
	t.Logf("Backup re-enabled: source=%s", sourcePath)
}

func TestIntegrationBackup_DestinationAddRemove(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	sourcePath := "/tmp/test-backup-dest-ops"
	destPath1 := "/tmp/test-backup-dest1"
	destPath2 := "/tmp/test-backup-dest2"

	backup := &pb.Backup{
		SourcePath: sourcePath,
		Destinations: []*pb.BackupDestination{
			{
				DestinationPath: destPath1,
				IsEnabled:       true,
			},
		},
		IsEnabled:              true,
		FileSystemWatchEnabled: false,
	}

	err := c.Backup().BackupAdd(ctx, backup)
	if err != nil {
		t.Fatalf("BackupAdd failed: %v", err)
	}
	defer c.Backup().BackupRemove(ctx, sourcePath)

	err = c.Backup().BackupAddDestination(ctx, &pb.BackupModifyRequest{
		SourcePath: sourcePath,
		Destinations: []*pb.BackupDestination{
			{
				DestinationPath: destPath2,
				IsEnabled:       true,
			},
		},
	})
	if err != nil {
		t.Fatalf("BackupAddDestination failed: %v", err)
	}
	t.Logf("Destination added: dest=%s", destPath2)

	status, err := c.Backup().BackupGetStatus(ctx, sourcePath)
	if err != nil {
		t.Fatalf("BackupGetStatus failed: %v", err)
	}

	if len(status.Backup.Destinations) != 2 {
		t.Errorf("Expected 2 destinations, got %d", len(status.Backup.Destinations))
	}

	err = c.Backup().BackupRemoveDestination(ctx, &pb.BackupModifyRequest{
		SourcePath: sourcePath,
		Destinations: []*pb.BackupDestination{
			{
				DestinationPath: destPath2,
			},
		},
	})
	if err != nil {
		t.Fatalf("BackupRemoveDestination failed: %v", err)
	}
	t.Logf("Destination removed: dest=%s", destPath2)

	statusFinal, err := c.Backup().BackupGetStatus(ctx, sourcePath)
	if err != nil {
		t.Fatalf("BackupGetStatus after remove failed: %v", err)
	}

	if len(statusFinal.Backup.Destinations) != 1 {
		t.Errorf("Expected 1 destination after remove, got %d", len(statusFinal.Backup.Destinations))
	}

	if statusFinal.Backup.Destinations[0].DestinationPath != destPath1 {
		t.Errorf("Expected remaining destination to be %s, got %s", destPath1, statusFinal.Backup.Destinations[0].DestinationPath)
	}
}
