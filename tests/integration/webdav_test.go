//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

func TestIntegrationWebDAV_ServerGet(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	config, err := c.WebDAV().GetDavServerConfig(ctx)
	if err != nil {
		t.Fatalf("GetDavServerConfig failed: %v", err)
	}

	t.Logf("WebDAV server config retrieved: enabled=%v", config.DavServerEnabled)
}

func TestIntegrationWebDAV_UserAddGetRemove(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	testUsername := "test-webdav-user-cli"
	testPassword := "test-password-123"

	err := c.WebDAV().AddDavUser(ctx, &pb.AddDavUserRequest{
		UserName: testUsername,
		Password: testPassword,
	})
	if err != nil {
		t.Fatalf("AddDavUser failed: %v", err)
	}
	t.Logf("Added WebDAV user: %s", testUsername)

	user, err := c.WebDAV().GetDavUser(ctx, testUsername)
	if err != nil {
		t.Fatalf("GetDavUser failed: %v", err)
	}

	if user.UserName != testUsername {
		t.Errorf("Expected username %s, got %s", testUsername, user.UserName)
	}
	t.Logf("Retrieved WebDAV user: username=%s, enabled=%v", user.UserName, user.Enabled)

	err = c.WebDAV().RemoveDavUser(ctx, testUsername)
	if err != nil {
		t.Fatalf("RemoveDavUser failed: %v", err)
	}
	t.Logf("Removed WebDAV user: %s", testUsername)
}

func TestIntegrationWebDAV_UserModify(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	testUsername := "test-webdav-modify-cli"
	testPassword := "test-password-123"

	err := c.WebDAV().AddDavUser(ctx, &pb.AddDavUserRequest{
		UserName: testUsername,
		Password: testPassword,
	})
	if err != nil {
		t.Fatalf("AddDavUser failed: %v", err)
	}

	defer func() {
		err := c.WebDAV().RemoveDavUser(ctx, testUsername)
		if err != nil {
			t.Logf("Cleanup RemoveDavUser failed: %v", err)
		}
	}()

	newPassword := "new-password-456"
	newRootPath := "/test-root"
	readOnly := true

	err = c.WebDAV().ModifyDavUser(ctx, &pb.ModifyDavUserRequest{
		UserName: testUsername,
		Password: &newPassword,
		RootPath: &newRootPath,
		ReadOnly: &readOnly,
	})
	if err != nil {
		t.Fatalf("ModifyDavUser failed: %v", err)
	}

	user, err := c.WebDAV().GetDavUser(ctx, testUsername)
	if err != nil {
		t.Fatalf("GetDavUser after modify failed: %v", err)
	}

	if user.RootPath != newRootPath {
		t.Errorf("Expected root path %s, got %s", newRootPath, user.RootPath)
	}
	if user.ReadOnly != readOnly {
		t.Errorf("Expected read-only %v, got %v", readOnly, user.ReadOnly)
	}

	t.Logf("Modified WebDAV user: username=%s, rootPath=%s, readOnly=%v", user.UserName, user.RootPath, user.ReadOnly)
}
