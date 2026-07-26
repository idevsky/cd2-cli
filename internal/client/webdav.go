package client

import (
	"context"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

type WebDAVAPI struct {
	c *Client
}

func (a *WebDAVAPI) AddDavUser(ctx context.Context, req *pb.AddDavUserRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.AddDavUser(ctx, req)
	return err
}

func (a *WebDAVAPI) RemoveDavUser(ctx context.Context, username string) error {
	ctx = a.c.withAuth(ctx)
	req := &pb.StringValue{Value: username}
	_, err := a.c.client.RemoveDavUser(ctx, req)
	return err
}

func (a *WebDAVAPI) ModifyDavUser(ctx context.Context, req *pb.ModifyDavUserRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.ModifyDavUser(ctx, req)
	return err
}

func (a *WebDAVAPI) GetDavUser(ctx context.Context, username string) (*pb.DavUser, error) {
	ctx = a.c.withAuth(ctx)
	req := &pb.StringValue{Value: username}
	return a.c.client.GetDavUser(ctx, req)
}

func (a *WebDAVAPI) GetDavServerConfig(ctx context.Context) (*pb.DavServerConfig, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetDavServerConfig(ctx, &emptypb.Empty{})
}

func (a *WebDAVAPI) SetDavServerConfig(ctx context.Context, req *pb.ModifyDavServerConfigRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.SetDavServerConfig(ctx, req)
	return err
}
