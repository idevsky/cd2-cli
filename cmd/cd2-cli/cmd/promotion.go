package cmd

import (
	"strconv"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

var promotionCmd = &cobra.Command{
	Use:   "promotion",
	Short: "Promotion operations",
}

var promotionListCmd = &cobra.Command{
	Use:   "list",
	Short: "Get promotions",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.Promotion().GetPromotions(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var promotionListCloudCmd = &cobra.Command{
	Use:   "list-cloud <cloud-name>",
	Short: "Get promotions by cloud",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		cloudName := args[0]

		result, err := cd2Client.Promotion().GetPromotionsByCloud(ctx, &pb.CloudAPIRequest{
			CloudName: cloudName,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var promotionUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update promotion result",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		err := cd2Client.Promotion().UpdatePromotionResult(ctx)
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "updated"})
	},
}

var promotionUpdateCloudCmd = &cobra.Command{
	Use:   "update-cloud <cloud-name>",
	Short: "Update promotion result by cloud",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		cloudName := args[0]

		err := cd2Client.Promotion().UpdatePromotionResultByCloud(ctx, &pb.UpdatePromotionResultByCloudRequest{
			CloudName: cloudName,
		})
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "updated"})
	},
}

var promotionPlanListCmd = &cobra.Command{
	Use:   "plan-list",
	Short: "Get CloudDrive plans",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.Promotion().GetCloudDrivePlans(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var promotionJoinPlanCmd = &cobra.Command{
	Use:   "join-plan <plan-id>",
	Short: "Join a plan",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		planId := args[0]
		coupon, _ := cmd.Flags().GetString("coupon")

		req := &pb.JoinPlanRequest{
			PlanId: planId,
		}
		if coupon != "" {
			req.CouponCode = &coupon
		}

		result, err := cd2Client.Promotion().JoinPlan(ctx, req)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var promotionBindCloudAccountCmd = &cobra.Command{
	Use:   "bind-cloud-account <cloud-name> <account-id>",
	Short: "Bind cloud account to promotion",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		cloudName := args[0]
		accountId := args[1]

		err := cd2Client.Promotion().BindCloudAccount(ctx, &pb.BindCloudAccountRequest{
			CloudName:      cloudName,
			CloudAccountId: accountId,
		})
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "bound"})
	},
}

var promotionTransferBalanceCmd = &cobra.Command{
	Use:   "transfer-balance <to-user> <amount> <password>",
	Short: "Transfer balance to another user",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		toUser := args[0]
		amountStr := args[1]
		password := args[2]

		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return err
		}

		err = cd2Client.Promotion().TransferBalance(ctx, &pb.TransferBalanceRequest{
			ToUserName: toUser,
			Amount:     amount,
			Password:   password,
		})
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "transferred"})
	},
}

var promotionActivatePlanCmd = &cobra.Command{
	Use:   "activate-plan <code>",
	Short: "Activate a plan with code",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		code := args[0]

		result, err := cd2Client.Promotion().ActivatePlan(ctx, code)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var promotionSendActionCmd = &cobra.Command{
	Use:   "send-action <cloud-name>",
	Short: "Send promotion action",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		cloudName := args[0]
		accountId, _ := cmd.Flags().GetString("account-id")
		promotionId, _ := cmd.Flags().GetString("promotion-id")

		req := &pb.SendPromotionActionRequest{
			CloudName: cloudName,
		}
		if accountId != "" {
			req.CloudAccountId = &accountId
		}
		if promotionId != "" {
			req.PromotionId = &promotionId
		}

		err := cd2Client.Promotion().SendPromotionAction(ctx, req)
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "sent"})
	},
}

var promotionReferralCodeCmd = &cobra.Command{
	Use:   "referral-code",
	Short: "Get referral code",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.Promotion().GetReferralCode(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var promotionBalanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Get balance",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.Promotion().GetBalanceLog(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

func init() {
	rootCmd.AddCommand(promotionCmd)

	promotionCmd.AddCommand(promotionListCmd)
	promotionCmd.AddCommand(promotionListCloudCmd)
	promotionCmd.AddCommand(promotionUpdateCmd)
	promotionCmd.AddCommand(promotionUpdateCloudCmd)
	promotionCmd.AddCommand(promotionPlanListCmd)
	promotionCmd.AddCommand(promotionJoinPlanCmd)
	promotionCmd.AddCommand(promotionBindCloudAccountCmd)
	promotionCmd.AddCommand(promotionTransferBalanceCmd)
	promotionCmd.AddCommand(promotionActivatePlanCmd)
	promotionCmd.AddCommand(promotionSendActionCmd)
	promotionCmd.AddCommand(promotionReferralCodeCmd)
	promotionCmd.AddCommand(promotionBalanceCmd)

	promotionJoinPlanCmd.Flags().String("coupon", "", "Coupon code")
	promotionSendActionCmd.Flags().String("account-id", "", "Cloud account ID")
	promotionSendActionCmd.Flags().String("promotion-id", "", "Promotion ID")

	setCommandID(promotionListCmd, "promotion.list")
	setCommandID(promotionListCloudCmd, "promotion.list-cloud")
	setCommandID(promotionUpdateCmd, "promotion.update")
	setCommandID(promotionUpdateCloudCmd, "promotion.update-cloud")
	setCommandID(promotionPlanListCmd, "promotion.plan-list")
	setCommandID(promotionJoinPlanCmd, "promotion.join-plan")
	setCommandID(promotionBindCloudAccountCmd, "promotion.bind-cloud-account")
	setCommandID(promotionTransferBalanceCmd, "promotion.transfer-balance")
	setCommandID(promotionActivatePlanCmd, "promotion.activate-plan")
	setCommandID(promotionSendActionCmd, "promotion.send-action")
	setCommandID(promotionReferralCodeCmd, "promotion.referral-code")
	setCommandID(promotionBalanceCmd, "promotion.balance")
}
