//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"
)

func TestIntegrationSync_FileChanges(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := c.Sync().SyncFileChangesFromCloud(ctx, "/")
	if err != nil {
		t.Fatalf("SyncFileChangesFromCloud failed: %v", err)
	}

	t.Logf("Sync result: %+v", result)
}

func TestIntegrationSync_WalkThroughTest(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := c.Sync().WalkThroughFolderTest(ctx, "/")
	if err != nil {
		t.Fatalf("WalkThroughFolderTest failed: %v", err)
	}

	t.Logf("Walk-through result: %+v", result)
}

func TestIntegrationSync_GetCD1UserData(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := c.Sync().GetCloudDrive1UserData(ctx)
	if err != nil {
		t.Fatalf("GetCloudDrive1UserData failed: %v", err)
	}

	t.Logf("CD1 user data: %s", result.Result)
}
