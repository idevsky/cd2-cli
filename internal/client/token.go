package client

import (
	"context"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

type TokenAPI struct {
	c *Client
}

func (a *TokenAPI) CreateToken(ctx context.Context, req *pb.CreateTokenRequest) (*pb.TokenInfo, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.CreateToken(ctx, req)
}

func (a *TokenAPI) ModifyToken(ctx context.Context, req *pb.ModifyTokenRequest) (*pb.TokenInfo, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.ModifyToken(ctx, req)
}

func (a *TokenAPI) RemoveToken(ctx context.Context, tokenId string) error {
	ctx = a.c.withAuth(ctx)
	req := &pb.StringValue{Value: tokenId}
	_, err := a.c.client.RemoveToken(ctx, req)
	return err
}

func (a *TokenAPI) ListTokens(ctx context.Context) (*pb.ListTokensResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.ListTokens(ctx, &emptypb.Empty{})
}

func (a *TokenAPI) GetTokenInfo(ctx context.Context, tokenId string) (*pb.TokenInfo, error) {
	ctx = a.c.withAuth(ctx)
	req := &pb.StringValue{Value: tokenId}
	return a.c.client.GetApiTokenInfo(ctx, req)
}
