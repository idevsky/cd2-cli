package client

import (
	"context"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

type SessionAPI struct {
	c *Client
}

func (a *SessionAPI) GetSessions(ctx context.Context) (*pb.GetSessionsResponse, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetSessions(ctx, &emptypb.Empty{})
}

func (a *SessionAPI) RevokeSession(ctx context.Context, sessionId string) error {
	ctx = a.c.withAuth(ctx)
	req := &pb.RevokeSessionRequest{SessionId: sessionId}
	_, err := a.c.client.RevokeSession(ctx, req)
	return err
}

func (a *SessionAPI) RevokeOtherSessions(ctx context.Context) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.RevokeOtherSessions(ctx, &emptypb.Empty{})
	return err
}
