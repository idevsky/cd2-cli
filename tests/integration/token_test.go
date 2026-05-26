//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

func TestIntegrationToken_CreateAndRemove(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tokensBefore, err := c.Token().ListTokens(ctx)
	if err != nil {
		t.Fatalf("ListTokens failed: %v", err)
	}
	countBefore := len(tokensBefore.Tokens)

	token, err := c.Token().CreateToken(ctx, &pb.CreateTokenRequest{
		RootDir:      "/",
		FriendlyName: "test-token",
		Permissions:  &pb.TokenPermissions{},
	})
	if err != nil {
		t.Fatalf("CreateToken failed: %v", err)
	}

	if token.Token == "" {
		t.Error("Expected non-empty token")
	}

	tokensAfter, err := c.Token().ListTokens(ctx)
	if err != nil {
		t.Fatalf("ListTokens after create failed: %v", err)
	}
	countAfter := len(tokensAfter.Tokens)

	if countAfter != countBefore+1 {
		t.Errorf("Expected %d tokens after create, got %d", countBefore+1, countAfter)
	}

	err = c.Token().RemoveToken(ctx, token.Token)
	if err != nil {
		t.Fatalf("RemoveToken failed: %v", err)
	}

	tokensFinal, err := c.Token().ListTokens(ctx)
	if err != nil {
		t.Fatalf("ListTokens after remove failed: %v", err)
	}
	countFinal := len(tokensFinal.Tokens)

	if countFinal != countBefore {
		t.Errorf("Expected %d tokens after remove, got %d", countBefore, countFinal)
	}

	t.Logf("Created and removed token: %s", token.Token)
}

func TestIntegrationToken_Info(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	token, err := c.Token().CreateToken(ctx, &pb.CreateTokenRequest{
		RootDir:      "/",
		FriendlyName: "test-info-token",
		Permissions:  &pb.TokenPermissions{},
	})
	if err != nil {
		t.Fatalf("CreateToken failed: %v", err)
	}
	defer c.Token().RemoveToken(ctx, token.Token)

	info, err := c.Token().GetTokenInfo(ctx, token.Token)
	if err != nil {
		t.Fatalf("GetTokenInfo failed: %v", err)
	}

	if info.Token != token.Token {
		t.Errorf("Expected token %s, got %s", token.Token, info.Token)
	}

	if info.FriendlyName != "test-info-token" {
		t.Errorf("Expected friendly name 'test-info-token', got '%s'", info.FriendlyName)
	}

	t.Logf("Token info: token=%s, name=%s, root=%s", info.Token, info.FriendlyName, info.RootDir)
}

func TestIntegrationToken_Modify(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	token, err := c.Token().CreateToken(ctx, &pb.CreateTokenRequest{
		RootDir:      "/",
		FriendlyName: "test-modify-token",
		Permissions:  &pb.TokenPermissions{},
	})
	if err != nil {
		t.Fatalf("CreateToken failed: %v", err)
	}
	defer c.Token().RemoveToken(ctx, token.Token)

	newName := "modified-name"
	newRoot := "/test"
	modified, err := c.Token().ModifyToken(ctx, &pb.ModifyTokenRequest{
		Token:        token.Token,
		FriendlyName: &newName,
		RootDir:      &newRoot,
	})
	if err != nil {
		t.Fatalf("ModifyToken failed: %v", err)
	}

	if modified.FriendlyName != newName {
		t.Errorf("Expected friendly name '%s', got '%s'", newName, modified.FriendlyName)
	}

	if modified.RootDir != newRoot {
		t.Errorf("Expected root dir '%s', got '%s'", newRoot, modified.RootDir)
	}

	t.Logf("Modified token: name=%s, root=%s", modified.FriendlyName, modified.RootDir)
}
