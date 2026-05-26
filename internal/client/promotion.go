package client

import (
	"context"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

type PromotionAPI struct {
	c *Client
}

func (a *PromotionAPI) GetPromotions(ctx context.Context) (*pb.GetPromotionsResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetPromotions(ctx, &emptypb.Empty{})
}

func (a *PromotionAPI) GetPromotionsByCloud(ctx context.Context, req *pb.CloudAPIRequest) (*pb.GetPromotionsResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetPromotionsByCloud(ctx, req)
}

func (a *PromotionAPI) UpdatePromotionResult(ctx context.Context) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.UpdatePromotionResult(ctx, &emptypb.Empty{})
	return err
}

func (a *PromotionAPI) UpdatePromotionResultByCloud(ctx context.Context, req *pb.UpdatePromotionResultByCloudRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.UpdatePromotionResultByCloud(ctx, req)
	return err
}

func (a *PromotionAPI) SendPromotionAction(ctx context.Context, req *pb.SendPromotionActionRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.SendPromotionAction(ctx, req)
	return err
}

func (a *PromotionAPI) GetCloudDrivePlans(ctx context.Context) (*pb.GetCloudDrivePlansResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetCloudDrivePlans(ctx, &emptypb.Empty{})
}

func (a *PromotionAPI) JoinPlan(ctx context.Context, req *pb.JoinPlanRequest) (*pb.JoinPlanResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.JoinPlan(ctx, req)
}

func (a *PromotionAPI) BindCloudAccount(ctx context.Context, req *pb.BindCloudAccountRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.BindCloudAccount(ctx, req)
	return err
}

func (a *PromotionAPI) TransferBalance(ctx context.Context, req *pb.TransferBalanceRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.TransferBalance(ctx, req)
	return err
}

func (a *PromotionAPI) GetBalanceLog(ctx context.Context) (*pb.BalanceLogResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetBalanceLog(ctx, &emptypb.Empty{})
}

func (a *PromotionAPI) CheckActivationCode(ctx context.Context, code string) (*pb.CheckActivationCodeResult, error) {
	ctx = a.c.withAuth(ctx)
	req := &pb.StringValue{Value: code}
	return a.c.client.CheckActivationCode(ctx, req)
}

func (a *PromotionAPI) ActivatePlan(ctx context.Context, code string) (*pb.JoinPlanResult, error) {
	ctx = a.c.withAuth(ctx)
	req := &pb.StringValue{Value: code}
	return a.c.client.ActivatePlan(ctx, req)
}

func (a *PromotionAPI) CheckCouponCode(ctx context.Context, req *pb.CheckCouponCodeRequest) (*pb.CouponCodeResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.CheckCouponCode(ctx, req)
}

func (a *PromotionAPI) GetReferralCode(ctx context.Context) (*pb.StringValue, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetReferralCode(ctx, &emptypb.Empty{})
}
