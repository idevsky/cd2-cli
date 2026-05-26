package client

import (
	"context"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

type RemoteUploadAPI struct {
	c *Client
}

func (a *RemoteUploadAPI) StartRemoteUpload(ctx context.Context, req *pb.StartRemoteUploadRequest) (*pb.RemoteUploadStarted, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.StartRemoteUpload(ctx, req)
}

func (a *RemoteUploadAPI) RemoteUploadControl(ctx context.Context, req *pb.RemoteUploadControlRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.RemoteUploadControl(ctx, req)
	return err
}

func (a *RemoteUploadAPI) RemoteReadData(ctx context.Context, req *pb.RemoteReadDataUpload) (*pb.RemoteReadDataReply, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.RemoteReadData(ctx, req)
}

func (a *RemoteUploadAPI) RemoteHashProgress(ctx context.Context, req *pb.RemoteHashProgressUpload) (*pb.RemoteHashProgressReply, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.RemoteHashProgress(ctx, req)
}
