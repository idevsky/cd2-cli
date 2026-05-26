package client

import (
	"context"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

type WebhookAPI struct {
	c *Client
}

func (a *WebhookAPI) GetWebhookConfigTemplate(ctx context.Context) (*pb.StringResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetWebhookConfigTemplate(ctx, &emptypb.Empty{})
}

func (a *WebhookAPI) GetWebhookConfigs(ctx context.Context) (*pb.WebhookList, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetWebhookConfigs(ctx, &emptypb.Empty{})
}

func (a *WebhookAPI) AddWebhookConfig(ctx context.Context, req *pb.WebhookRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.AddWebhookConfig(ctx, req)
	return err
}

func (a *WebhookAPI) RemoveWebhookConfig(ctx context.Context, webhookId string) error {
	ctx = a.c.withAuth(ctx)
	req := &pb.StringValue{Value: webhookId}
	_, err := a.c.client.RemoveWebhookConfig(ctx, req)
	return err
}

func (a *WebhookAPI) ChangeWebhookConfig(ctx context.Context, req *pb.WebhookRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.ChangeWebhookConfig(ctx, req)
	return err
}
