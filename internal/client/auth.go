package client

import (
	"context"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

type AuthAPI struct {
	c *Client
}

func (a *AuthAPI) SendConfirmEmail(ctx context.Context) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.SendConfirmEmail(ctx, &emptypb.Empty{})
	return err
}

func (a *AuthAPI) ConfirmEmail(ctx context.Context, req *pb.ConfirmEmailRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.ConfirmEmail(ctx, req)
	return err
}

func (a *AuthAPI) GetAccountStatus(ctx context.Context) (*pb.AccountStatusResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetAccountStatus(ctx, &emptypb.Empty{})
}

func (a *AuthAPI) Logout(ctx context.Context, req *pb.UserLogoutRequest) (*pb.FileOperationResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.Logout(ctx, req)
}

func (a *AuthAPI) Check2FAStatus(ctx context.Context) (*pb.TwoFactorAuthStatusResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.Check2FAStatus(ctx, &emptypb.Empty{})
}

func (a *AuthAPI) Setup2FA(ctx context.Context, req *pb.Setup2FARequest) (*pb.TwoFactorAuthSetupResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.Setup2FA(ctx, req)
}

func (a *AuthAPI) Enable2FA(ctx context.Context, req *pb.TwoFactorAuthCodeRequest) (*pb.TwoFactorAuthEnableResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.Enable2FA(ctx, req)
}

func (a *AuthAPI) Disable2FA(ctx context.Context, req *pb.TwoFactorAuthCodeRequest) (*pb.TwoFactorAuthMessageResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.Disable2FA(ctx, req)
}

func (a *AuthAPI) GetRecoveryCodes(ctx context.Context, req *pb.TwoFactorAuthCodeRequest) (*pb.TwoFactorAuthRecoveryCodesResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetRecoveryCodes(ctx, req)
}

func (a *AuthAPI) RegenerateRecoveryCodes(ctx context.Context, req *pb.TwoFactorAuthCodeRequest) (*pb.TwoFactorAuthRecoveryCodesResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.RegenerateRecoveryCodes(ctx, req)
}

func (a *AuthAPI) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.FileOperationResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.ChangePassword(ctx, req)
}

func (a *AuthAPI) SendChangeEmailCode(ctx context.Context, req *pb.SendChangeEmailCodeRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.SendChangeEmailCode(ctx, req)
	return err
}

func (a *AuthAPI) ChangeEmail(ctx context.Context, req *pb.ChangeEmailRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.ChangeEmail(ctx, req)
	return err
}

func (a *AuthAPI) ChangeEmailAndPassword(ctx context.Context, req *pb.ChangeEmailAndPasswordRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.ChangeEmailAndPassword(ctx, req)
	return err
}
