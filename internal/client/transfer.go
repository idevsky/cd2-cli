package client

import (
	"context"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

type TransferAPI struct {
	c *Client
}

func (a *TransferAPI) GetAllTasksCount(ctx context.Context) (*pb.GetAllTasksCountResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetAllTasksCount(ctx, &emptypb.Empty{})
}

func (a *TransferAPI) GetDownloadFileCount(ctx context.Context) (*pb.GetDownloadFileCountResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetDownloadFileCount(ctx, &emptypb.Empty{})
}

func (a *TransferAPI) GetDownloadFileList(ctx context.Context) (*pb.GetDownloadFileListResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetDownloadFileList(ctx, &emptypb.Empty{})
}

func (a *TransferAPI) GetUploadFileCount(ctx context.Context) (*pb.GetUploadFileCountResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetUploadFileCount(ctx, &emptypb.Empty{})
}

func (a *TransferAPI) GetUploadFileList(ctx context.Context, req *pb.GetUploadFileListRequest) (*pb.GetUploadFileListResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetUploadFileList(ctx, req)
}

func (a *TransferAPI) CancelAllUploadFiles(ctx context.Context) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.CancelAllUploadFiles(ctx, &emptypb.Empty{})
	return err
}

func (a *TransferAPI) CancelUploadFiles(ctx context.Context, keys []string) error {
	ctx = a.c.withAuth(ctx)
	req := &pb.MultpleUploadFileKeyRequest{Keys: keys}
	_, err := a.c.client.CancelUploadFiles(ctx, req)
	return err
}

func (a *TransferAPI) PauseAllUploadFiles(ctx context.Context) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.PauseAllUploadFiles(ctx, &emptypb.Empty{})
	return err
}

func (a *TransferAPI) PauseUploadFiles(ctx context.Context, keys []string) error {
	ctx = a.c.withAuth(ctx)
	req := &pb.MultpleUploadFileKeyRequest{Keys: keys}
	_, err := a.c.client.PauseUploadFiles(ctx, req)
	return err
}

func (a *TransferAPI) ResumeAllUploadFiles(ctx context.Context) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.ResumeAllUploadFiles(ctx, &emptypb.Empty{})
	return err
}

func (a *TransferAPI) ResumeUploadFiles(ctx context.Context, keys []string) error {
	ctx = a.c.withAuth(ctx)
	req := &pb.MultpleUploadFileKeyRequest{Keys: keys}
	_, err := a.c.client.ResumeUploadFiles(ctx, req)
	return err
}
