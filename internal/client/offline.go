package client

import (
	"context"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

type OfflineAPI struct {
	c *Client
}

func (a *OfflineAPI) AddOfflineFiles(ctx context.Context, req *pb.AddOfflineFileRequest) (*pb.FileOperationResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.AddOfflineFiles(ctx, req)
}

func (a *OfflineAPI) RemoveOfflineFiles(ctx context.Context, req *pb.RemoveOfflineFilesRequest) (*pb.FileOperationResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.RemoveOfflineFiles(ctx, req)
}

func (a *OfflineAPI) ListOfflineFilesByPath(ctx context.Context, path string) (*pb.OfflineFileListResult, error) {
	ctx = a.c.withAuth(ctx)
	req := &pb.FileRequest{Path: path}
	return a.c.client.ListOfflineFilesByPath(ctx, req)
}

func (a *OfflineAPI) ListAllOfflineFiles(ctx context.Context, req *pb.OfflineFileListAllRequest) (*pb.OfflineFileListAllResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.ListAllOfflineFiles(ctx, req)
}

func (a *OfflineAPI) GetOfflineQuotaInfo(ctx context.Context, req *pb.OfflineQuotaRequest) (*pb.OfflineQuotaInfo, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetOfflineQuotaInfo(ctx, req)
}

func (a *OfflineAPI) ClearOfflineFiles(ctx context.Context, req *pb.ClearOfflineFileRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.ClearOfflineFiles(ctx, req)
	return err
}

func (a *OfflineAPI) RestartOfflineTask(ctx context.Context, req *pb.RestartOfflineFileRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.RestartOfflineTask(ctx, req)
	return err
}
