//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

func TestIntegrationTransfer_UploadCount(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := c.Transfer().GetUploadFileCount(ctx)
	if err != nil {
		t.Fatalf("GetUploadFileCount failed: %v", err)
	}

	t.Logf("Upload files: %v", result)
}

func TestIntegrationTransfer_DownloadCount(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := c.Transfer().GetDownloadFileCount(ctx)
	if err != nil {
		t.Fatalf("GetDownloadFileCount failed: %v", err)
	}

	t.Logf("Download files: %v", result)
}

func TestIntegrationTransfer_GetAllTasksCount(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := c.Transfer().GetAllTasksCount(ctx)
	if err != nil {
		t.Fatalf("GetAllTasksCount failed: %v", err)
	}

	t.Logf("All tasks count: %v", result)
}

func TestIntegrationTransfer_UploadList(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tasks, err := c.Transfer().GetUploadFileList(ctx, &pb.GetUploadFileListRequest{})
	if err != nil {
		t.Fatalf("GetUploadFileList failed: %v", err)
	}

	t.Logf("Upload tasks: %v", tasks)
}

func TestIntegrationTransfer_DownloadList(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tasks, err := c.Transfer().GetDownloadFileList(ctx)
	if err != nil {
		t.Fatalf("GetDownloadFileList failed: %v", err)
	}

	t.Logf("Download tasks: %v", tasks)
}

func TestIntegrationTransfer_GetDownloadUrl(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	files, err := c.File().GetSubFiles(ctx, &pb.ListSubFileRequest{
		Path:         "/",
		ForceRefresh: false,
	})
	if err != nil {
		t.Fatalf("GetSubFiles failed: %v", err)
	}

	if len(files) == 0 {
		t.Log("No files found for download test")
		return
	}

	var filePath string
	for _, f := range files {
		if !f.IsDirectory {
			filePath = f.FullPathName
			break
		}
	}

	if filePath == "" {
		t.Log("No files available for download test")
		return
	}

	urlResp, err := c.File().GetDownloadUrl(ctx, &pb.GetDownloadUrlPathRequest{
		Path: filePath,
	})
	if err != nil {
		t.Fatalf("GetDownloadUrl failed: %v", err)
	}

	t.Logf("Download URL: %v", urlResp)
}
