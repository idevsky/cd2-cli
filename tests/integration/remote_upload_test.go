//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

func TestIntegrationRemoteUpload_Start(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := c.RemoteUpload().StartRemoteUpload(ctx, &pb.StartRemoteUploadRequest{
		FilePath:                 "/test-upload.txt",
		FileSize:                 1024,
		ClientCanCalculateHashes: false,
	})
	if err != nil {
		t.Fatalf("StartRemoteUpload failed: %v", err)
	}

	t.Logf("Remote upload started with ID: %s", result.UploadId)
}

func TestIntegrationRemoteUpload_Control(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	startResult, err := c.RemoteUpload().StartRemoteUpload(ctx, &pb.StartRemoteUploadRequest{
		FilePath:                 "/test-control.txt",
		FileSize:                 2048,
		ClientCanCalculateHashes: false,
	})
	if err != nil {
		t.Fatalf("StartRemoteUpload failed: %v", err)
	}

	uploadID := startResult.UploadId

	err = c.RemoteUpload().RemoteUploadControl(ctx, &pb.RemoteUploadControlRequest{
		UploadId: uploadID,
		Control:  &pb.RemoteUploadControlRequest_Cancel{Cancel: &pb.CancelRemoteUpload{}},
	})
	if err != nil {
		t.Fatalf("RemoteUploadControl (cancel) failed: %v", err)
	}

	t.Logf("Remote upload %s cancelled successfully", uploadID)
}

func TestIntegrationRemoteUpload_ReadData(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := c.RemoteUpload().RemoteReadData(ctx, &pb.RemoteReadDataUpload{
		UploadId:    "test-upload-id",
		Offset:      0,
		Length:      1024,
		LazyRead:    false,
		IsLastChunk: false,
	})
	if err != nil {
		t.Fatalf("RemoteReadData failed: %v", err)
	}

	t.Logf("Remote read data result: success=%v, bytes_received=%d", result.Success, result.BytesReceived)
}

func TestIntegrationRemoteUpload_HashProgress(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	hashValue := "abc123"
	result, err := c.RemoteUpload().RemoteHashProgress(ctx, &pb.RemoteHashProgressUpload{
		UploadId:    "test-upload-id",
		BytesHashed: 512,
		TotalBytes:  1024,
		HashType:    pb.CloudDriveFile_Sha1,
		HashValue:   &hashValue,
	})
	if err != nil {
		t.Fatalf("RemoteHashProgress failed: %v", err)
	}

	t.Logf("Remote hash progress result: %+v", result)
}
