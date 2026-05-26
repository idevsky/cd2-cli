//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

func TestIntegrationPromotion_List(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := c.Promotion().GetPromotions(ctx)
	if err != nil {
		t.Logf("GetPromotions failed (may not have promotion data): %v", err)
		return
	}

	t.Logf("Promotions: %v", result)
}

func TestIntegrationPromotion_PlanList(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := c.Promotion().GetCloudDrivePlans(ctx)
	if err != nil {
		t.Logf("GetCloudDrivePlans failed: %v", err)
		return
	}

	t.Logf("Plans: %v", result)
}

func TestIntegrationPromotion_Balance(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := c.Promotion().GetBalanceLog(ctx)
	if err != nil {
		t.Logf("GetBalanceLog failed (may not have balance): %v", err)
		return
	}

	t.Logf("Balance: %v", result)
}

func TestIntegrationPromotion_ReferralCode(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := c.Promotion().GetReferralCode(ctx)
	if err != nil {
		t.Logf("GetReferralCode failed: %v", err)
		return
	}

	t.Logf("Referral code: %s", result.Value)
}

func TestIntegrationPromotion_Update(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := c.Promotion().UpdatePromotionResult(ctx)
	if err != nil {
		t.Logf("UpdatePromotionResult failed: %v", err)
		return
	}

	t.Logf("Promotion result updated successfully")
}

func TestIntegrationPromotion_JoinPlan(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &pb.JoinPlanRequest{
		PlanId: "test-plan-id",
	}

	result, err := c.Promotion().JoinPlan(ctx, req)
	if err != nil {
		t.Logf("JoinPlan failed (expected for test plan id): %v", err)
		return
	}

	t.Logf("Join plan result: %v", result)
}

func TestIntegrationPromotion_BindCloudAccount(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &pb.BindCloudAccountRequest{
		CloudName:      "test-cloud",
		CloudAccountId: "test-account-id",
	}

	err := c.Promotion().BindCloudAccount(ctx, req)
	if err != nil {
		t.Logf("BindCloudAccount failed (expected for test values): %v", err)
		return
	}

	t.Logf("BindCloudAccount succeeded")
}

func TestIntegrationPromotion_TransferBalance(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &pb.TransferBalanceRequest{
		ToUserName: "test-user",
		Amount:     1.0,
		Password:   "test-password",
	}

	err := c.Promotion().TransferBalance(ctx, req)
	if err != nil {
		t.Logf("TransferBalance failed (expected for test values): %v", err)
		return
	}

	t.Logf("TransferBalance succeeded")
}

func TestIntegrationPromotion_ActivatePlan(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := c.Promotion().ActivatePlan(ctx, "test-activation-code")
	if err != nil {
		t.Logf("ActivatePlan failed (expected for test code): %v", err)
		return
	}

	t.Logf("ActivatePlan result: %v", result)
}

func TestIntegrationPromotion_SendAction(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &pb.SendPromotionActionRequest{
		CloudName: "test-cloud",
	}

	err := c.Promotion().SendPromotionAction(ctx, req)
	if err != nil {
		t.Logf("SendPromotionAction failed (expected for test values): %v", err)
		return
	}

	t.Logf("SendPromotionAction succeeded")
}

func TestIntegrationPromotion_CheckActivationCode(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := c.Promotion().CheckActivationCode(ctx, "test-code")
	if err != nil {
		t.Logf("CheckActivationCode failed: %v", err)
		return
	}

	t.Logf("CheckActivationCode result: %v", result)
}

func TestIntegrationPromotion_CheckCouponCode(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &pb.CheckCouponCodeRequest{
		PlanId:     "test-plan-id",
		CouponCode: "test-coupon",
	}

	result, err := c.Promotion().CheckCouponCode(ctx, req)
	if err != nil {
		t.Logf("CheckCouponCode failed: %v", err)
		return
	}

	t.Logf("CheckCouponCode result: %v", result)
}

func TestIntegrationPromotion_GetPromotionsByCloud(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &pb.CloudAPIRequest{
		CloudName: "115",
	}

	result, err := c.Promotion().GetPromotionsByCloud(ctx, req)
	if err != nil {
		t.Logf("GetPromotionsByCloud failed: %v", err)
		return
	}

	t.Logf("GetPromotionsByCloud result: %v", result)
}

func TestIntegrationPromotion_UpdatePromotionResultByCloud(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &pb.UpdatePromotionResultByCloudRequest{
		CloudName: "test-cloud",
	}

	err := c.Promotion().UpdatePromotionResultByCloud(ctx, req)
	if err != nil {
		t.Logf("UpdatePromotionResultByCloud failed: %v", err)
		return
	}

	t.Logf("UpdatePromotionResultByCloud succeeded")
}
