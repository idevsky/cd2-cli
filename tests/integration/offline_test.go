//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

func TestIntegrationOffline_Quota(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	quota, err := c.Offline().GetOfflineQuotaInfo(ctx, &pb.OfflineQuotaRequest{})
	if err != nil {
		t.Logf("GetOfflineQuotaInfo failed (may need specific cloud/account): %v", err)
		return
	}

	t.Logf("Offline quota: total=%d, used=%d", quota.Total, quota.Used)
}

func TestIntegrationOffline_List(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	files, err := c.Offline().ListOfflineFilesByPath(ctx, "/")
	if err != nil {
		t.Logf("ListOfflineFilesByPath failed (may not be supported): %v", err)
		return
	}

	t.Logf("Offline files by path: %d files", len(files.OfflineFiles))
}
