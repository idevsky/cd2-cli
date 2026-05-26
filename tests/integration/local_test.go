//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

func TestIntegrationLocal_List(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := c.Local().LocalGetSubFiles(ctx, &pb.LocalGetSubFilesRequest{
		ParentFolder:          "/",
		FolderOnly:            false,
		IncludeCloudDrive:     false,
		IncludeAvailableDrive: false,
	})
	if err != nil {
		t.Fatalf("LocalGetSubFiles failed: %v", err)
	}

	t.Logf("Local files: %v", result)
}

func TestIntegrationLocal_Mkdir(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := c.Local().LocalCreateFolder(ctx, &pb.LocalCreateFolderRequest{
		ParentFolder: "/tmp",
		FolderName:   "cd2-cli-test-mkdir",
	})
	if err != nil {
		t.Fatalf("LocalCreateFolder failed: %v", err)
	}

	t.Logf("Created folder: %v", result)
}
