//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

func TestIntegrationStorage_List(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	apis, err := c.CloudAPI().GetAllCloudApis(ctx)
	if err != nil {
		t.Fatalf("GetAllCloudApis failed: %v", err)
	}

	t.Logf("Found %d storage providers", len(apis.Apis))
}

func TestIntegrationStorage_CanAddMore(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := c.CloudAPI().CanAddMoreCloudApis(ctx)
	if err != nil {
		t.Fatalf("CanAddMoreCloudApis failed: %v", err)
	}

	t.Logf("Can add more APIs: %v", result.Success)
}

func TestIntegrationStorage_AddS3(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	addResult, err := c.CloudAPI().APILoginS3(ctx, &pb.LoginS3Request{
		AccessKeyId:     "minioadmin",
		SecretAccessKey: "minioadmin123",
	})
	if err != nil {
		t.Logf("APILoginS3 failed (MinIO may not be available): %v", err)
		return
	}

	t.Logf("S3 storage test result: %v", addResult)
}

func TestIntegrationStorage_Remove(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	apis, err := c.CloudAPI().GetAllCloudApis(ctx)
	if err != nil {
		t.Fatalf("GetAllCloudApis failed: %v", err)
	}

	if len(apis.Apis) == 0 {
		t.Skip("No storage providers available to remove")
	}

	api := apis.Apis[0]
	result, err := c.CloudAPI().RemoveCloudAPI(ctx, &pb.RemoveCloudAPIRequest{
		CloudName:       api.Name,
		UserName:        api.UserName,
		PermanentRemove: false,
	})
	if err != nil {
		t.Fatalf("RemoveCloudAPI failed: %v", err)
	}

	t.Logf("Remove storage result: success=%v, cloud=%s, user=%s", result.Success, api.Name, api.UserName)
}

func TestIntegrationStorage_UpdateConfig(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	apis, err := c.CloudAPI().GetAllCloudApis(ctx)
	if err != nil {
		t.Fatalf("GetAllCloudApis failed: %v", err)
	}

	if len(apis.Apis) == 0 {
		t.Skip("No storage providers available to update config")
	}

	api := apis.Apis[0]

	config, err := c.CloudAPI().GetCloudAPIConfig(ctx, &pb.GetCloudAPIConfigRequest{
		CloudName: api.Name,
		UserName:  api.UserName,
	})
	if err != nil {
		t.Fatalf("GetCloudAPIConfig failed: %v", err)
	}

	t.Logf("Current config: maxDownloadThreads=%d, maxBufferPoolSizeMB=%d",
		config.MaxDownloadThreads, config.MaxBufferPoolSizeMB)

	err = c.CloudAPI().SetCloudAPIConfig(ctx, &pb.SetCloudAPIConfigRequest{
		CloudName: api.Name,
		UserName:  api.UserName,
		Config: &pb.CloudAPIConfig{
			MaxDownloadThreads:  8,
			MaxBufferPoolSizeMB: 256,
		},
	})
	if err != nil {
		t.Fatalf("SetCloudAPIConfig failed: %v", err)
	}

	t.Logf("Successfully updated storage config for %s/%s", api.Name, api.UserName)
}
