package client

import (
	"context"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

type SyncAPI struct {
	c *Client
}

func (a *SyncAPI) SyncFileChangesFromCloud(ctx context.Context, path string) (*pb.FileSystemChangeStatistics, error) {
	ctx = a.c.withAuth(ctx)
	req := &pb.FileRequest{Path: path}
	return a.c.client.SyncFileChangesFromCloud(ctx, req)
}

func (a *SyncAPI) StartCloudEventListener(ctx context.Context, path string) error {
	ctx = a.c.withAuth(ctx)
	req := &pb.FileRequest{Path: path}
	_, err := a.c.client.StartCloudEventListener(ctx, req)
	return err
}

func (a *SyncAPI) StopCloudEventListener(ctx context.Context, path string) error {
	ctx = a.c.withAuth(ctx)
	req := &pb.FileRequest{Path: path}
	_, err := a.c.client.StopCloudEventListener(ctx, req)
	return err
}

func (a *SyncAPI) WalkThroughFolderTest(ctx context.Context, path string) (*pb.WalkThroughFolderResult, error) {
	ctx = a.c.withAuth(ctx)
	req := &pb.FileRequest{Path: path}
	return a.c.client.WalkThroughFolderTest(ctx, req)
}

func (a *SyncAPI) GetCloudDrive1UserData(ctx context.Context) (*pb.StringResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetCloudDrive1UserData(ctx, &emptypb.Empty{})
}
