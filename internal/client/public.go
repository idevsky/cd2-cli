package client

import (
	"context"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

type PublicAPI struct {
	c *Client
}

func (a *PublicAPI) GetSystemInfo(ctx context.Context) (*pb.CloudDriveSystemInfo, error) {
	return a.c.client.GetSystemInfo(ctx, &emptypb.Empty{})
}

func (a *PublicAPI) GetToken(ctx context.Context, req *pb.GetTokenRequest) (*pb.JWTToken, error) {
	resp, err := a.c.client.GetToken(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Success && resp.Token != "" {
		a.c.SetToken(resp.Token)
	}
	return resp, nil
}

func (a *PublicAPI) Login(ctx context.Context, req *pb.UserLoginRequest) (*pb.FileOperationResult, error) {
	return a.c.client.Login(ctx, req)
}

func (a *PublicAPI) LoginWithThirdPartyAccount(ctx context.Context, req *pb.LoginWithThirdPartyAccountRequest) (*pb.JWTToken, error) {
	resp, err := a.c.client.LoginWithThirdPartyAccount(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Success && resp.Token != "" {
		a.c.SetToken(resp.Token)
	}
	return resp, nil
}

func (a *PublicAPI) Register(ctx context.Context, req *pb.UserRegisterRequest) (*pb.FileOperationResult, error) {
	return a.c.client.Register(ctx, req)
}

func (a *PublicAPI) SendResetAccountEmail(ctx context.Context, req *pb.SendResetAccountEmailRequest) error {
	_, err := a.c.client.SendResetAccountEmail(ctx, req)
	return err
}

func (a *PublicAPI) ResetAccount(ctx context.Context, req *pb.ResetAccountRequest) error {
	_, err := a.c.client.ResetAccount(ctx, req)
	return err
}

func (a *PublicAPI) GetApiTokenInfo(ctx context.Context, token string) (*pb.TokenInfo, error) {
	req := &pb.StringValue{Value: token}
	return a.c.client.GetApiTokenInfo(ctx, req)
}

func (a *PublicAPI) LoginWith2FA(ctx context.Context, req *pb.LoginWith2FARequest) (*pb.JWTToken, error) {
	resp, err := a.c.client.LoginWith2FA(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Success && resp.Token != "" {
		a.c.SetToken(resp.Token)
	}
	return resp, nil
}
