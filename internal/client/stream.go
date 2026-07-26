package client

import (
	"context"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

type StreamAPI struct {
	c *Client
}

func (a *StreamAPI) PushTaskChange(ctx context.Context) ([]*pb.GetAllTasksCountResult, error) {
	ctx = a.c.withAuth(ctx)
	stream, err := a.c.client.PushTaskChange(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return collectStream[*pb.GetAllTasksCountResult](stream)
}

func (a *StreamAPI) PushMessage(ctx context.Context) ([]*pb.CloudDrivePushMessage, error) {
	ctx = a.c.withAuth(ctx)
	stream, err := a.c.client.PushMessage(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return collectStream[*pb.CloudDrivePushMessage](stream)
}

func (a *StreamAPI) RemoteUploadChannel(ctx context.Context, req *pb.RemoteUploadChannelRequest) ([]*pb.RemoteUploadChannelReply, error) {
	ctx = a.c.withAuth(ctx)
	stream, err := a.c.client.RemoteUploadChannel(ctx, req)
	if err != nil {
		return nil, err
	}
	return collectStream[*pb.RemoteUploadChannelReply](stream)
}
