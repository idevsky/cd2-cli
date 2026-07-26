//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"
)

func TestIntegrationMount_List(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mounts, err := c.Mount().GetMountPoints(ctx)
	if err != nil {
		t.Fatalf("GetMountPoints failed: %v", err)
	}

	t.Logf("Found %d mount points", len(mounts.MountPoints))
}

func TestIntegrationMount_CanAddMore(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := c.Mount().CanAddMoreMountPoints(ctx)
	if err != nil {
		t.Fatalf("CanAddMoreMountPoints failed: %v", err)
	}

	t.Logf("Can add more mount points: %v", result.Success)
}

func TestIntegrationMount_HasDriveLetters(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := c.Mount().HasDriveLetters(ctx)
	if err != nil {
		t.Fatalf("HasDriveLetters failed: %v", err)
	}

	t.Logf("Has drive letters: %v", result)
}

func TestIntegrationMount_DriveLetters(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := c.Mount().GetAvailableDriveLetters(ctx)
	if err != nil {
		t.Fatalf("GetAvailableDriveLetters failed: %v", err)
	}

	t.Logf("Available drive letters: %v", result.DriveLetters)
}

func TestIntegrationMount_CanMountBoth(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := c.Mount().CanMountBothLocalAndCloud(ctx)
	if err != nil {
		t.Fatalf("CanMountBothLocalAndCloud failed: %v", err)
	}

	t.Logf("Can mount both local and cloud: %v", result.Result)
}
