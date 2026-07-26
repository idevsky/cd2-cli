//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

func TestIntegrationAuth_Login(t *testing.T) {
	c := getTestClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := c.Public().GetToken(ctx, &pb.GetTokenRequest{
		UserName: "admin",
		Password: "admin123",
	})
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if !resp.Success {
		t.Errorf("Login not successful: %s", resp.ErrorMessage)
	}

	if resp.Token == "" {
		t.Error("Empty token returned")
	}

	t.Logf("Login successful, token expires at: %v", resp.Expiration)
}

func TestIntegrationAuth_Status(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	status, err := c.Auth().GetAccountStatus(ctx)
	if err != nil {
		t.Fatalf("GetAccountStatus failed: %v", err)
	}

	if status.UserName != "admin" {
		t.Errorf("Expected username 'admin', got '%s'", status.UserName)
	}

	t.Logf("Account status: UserName=%s", status.UserName)
}

func TestIntegrationAuth_Logout(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := c.Auth().Logout(ctx, &pb.UserLogoutRequest{
		LogoutFromCloudFS: false,
	})
	if err != nil {
		t.Fatalf("Logout failed: %v", err)
	}

	if !result.Success {
		t.Errorf("Logout not successful")
	}

	t.Log("Logout successful")
}

func TestIntegrationAuth_SystemInfo(t *testing.T) {
	c := getTestClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	info, err := c.Public().GetSystemInfo(ctx)
	if err != nil {
		t.Fatalf("GetSystemInfo failed: %v", err)
	}

	t.Logf("System info: IsLogin=%v, UserName=%s, SystemReady=%v",
		info.IsLogin, info.UserName, info.SystemReady)
}

func TestIntegrationAuth_SendConfirmEmail(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := c.Auth().SendConfirmEmail(ctx)
	if err != nil {
		t.Logf("SendConfirmEmail returned error (may be expected): %v", err)
	} else {
		t.Log("SendConfirmEmail succeeded")
	}
}

func TestIntegrationAuth_SendChangeEmailCode(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := c.Auth().SendChangeEmailCode(ctx, &pb.SendChangeEmailCodeRequest{
		NewEmail: "newemail@example.com",
		Password: "testpassword",
	})
	if err != nil {
		t.Logf("SendChangeEmailCode returned error (may be expected): %v", err)
	} else {
		t.Log("SendChangeEmailCode succeeded")
	}
}

func TestIntegrationAuth_SendResetEmail(t *testing.T) {
	c := getTestClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := c.Public().SendResetAccountEmail(ctx, &pb.SendResetAccountEmailRequest{
		Email: "test@example.com",
	})
	if err != nil {
		t.Logf("SendResetAccountEmail returned error (may be expected): %v", err)
	} else {
		t.Log("SendResetAccountEmail succeeded")
	}
}

func TestIntegrationAuth_LoginWithThirdPartyAccount(t *testing.T) {
	c := getTestClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := c.Public().LoginWithThirdPartyAccount(ctx, &pb.LoginWithThirdPartyAccountRequest{
		CloudName:    "testcloud",
		RefreshToken: "test-token",
	})
	if err != nil {
		t.Logf("LoginWithThirdPartyAccount returned error (may be expected): %v", err)
	} else {
		t.Log("LoginWithThirdPartyAccount succeeded")
	}
}
