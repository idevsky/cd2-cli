//go:build integration
// +build integration

package tests

import (
	"context"
	"testing"
	"time"

	"github.com/clouddrive/cd2-cli/internal/client"
	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

func TestIntegrationPublicAPI(t *testing.T) {
	c, err := client.NewClient(client.Config{Address: "localhost:19798"})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	info, err := c.Public().GetSystemInfo(ctx)
	if err != nil {
		t.Fatalf("GetSystemInfo failed: %v", err)
	}

	if !info.SystemReady {
		t.Logf("System not ready yet: %v", info)
	}

	t.Logf("System info: IsLogin=%v, UserName=%s, SystemReady=%v", info.IsLogin, info.UserName, info.SystemReady)
}

func TestIntegrationGetToken(t *testing.T) {
	c, _ := client.NewClient(client.Config{Address: "localhost:19798"})
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := c.Public().GetToken(ctx, &pb.GetTokenRequest{
		UserName: "admin",
		Password: "admin123",
	})
	if err != nil {
		t.Logf("GetToken error (CloudFS not login expected for first run): %v", err)
		return
	}

	if !resp.Success {
		t.Logf("Login not successful: %s", resp.ErrorMessage)
	} else {
		t.Logf("Token obtained, expires at: %v", resp.Expiration)
		c.SetToken(resp.Token)
	}
}
