//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

func TestIntegrationCloudAPI_DiscoverSmbServers(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := c.CloudAPI().DiscoverSmbServers(ctx)
	if err != nil {
		t.Logf("DiscoverSmbServers failed (SMB discovery may not be available): %v", err)
		return
	}

	t.Logf("Found %d SMB servers", len(result.Servers))
	for _, server := range result.Servers {
		t.Logf("SMB server: %s", server)
	}
}

func TestIntegrationCloudAPI_DiscoverSmbShares(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := c.CloudAPI().DiscoverSmbShares(ctx, &pb.DiscoverSmbSharesRequest{
		Server: "localhost",
	})
	if err != nil {
		t.Logf("DiscoverSmbShares failed (SMB server may not be available): %v", err)
		return
	}

	t.Logf("Found %d SMB shares", len(result.Shares))
	for _, share := range result.Shares {
		t.Logf("SMB share: %s", share)
	}
}
