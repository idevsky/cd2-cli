package cmd

import (
	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

func init() {
	authCmd.AddCommand(twoFAStatusCmd)
	authCmd.AddCommand(twoFASetupCmd)
	authCmd.AddCommand(twoFAEnableCmd)
	authCmd.AddCommand(twoFADisableCmd)
	authCmd.AddCommand(twoFARecoveryCodesCmd)
	authCmd.AddCommand(twoFARegenerateCodesCmd)

	twoFASetupCmd.Flags().String("password-file", "", "Read password from file")
	twoFAEnableCmd.Flags().String("totp-file", "", "Read TOTP code from file")
	twoFADisableCmd.Flags().String("totp-file", "", "Read TOTP code from file")
	twoFARecoveryCodesCmd.Flags().String("totp-file", "", "Read TOTP code from file")
	twoFARegenerateCodesCmd.Flags().String("totp-file", "", "Read TOTP code from file")

	setCommandID(twoFAStatusCmd, "auth.2fa-status")
	setCommandID(twoFASetupCmd, "auth.2fa-setup")
	setCommandID(twoFAEnableCmd, "auth.2fa-enable")
	setCommandID(twoFADisableCmd, "auth.2fa-disable")
	setCommandID(twoFARecoveryCodesCmd, "auth.2fa-recovery-codes")
	setCommandID(twoFARegenerateCodesCmd, "auth.2fa-regenerate-codes")
}

var twoFAStatusCmd = &cobra.Command{
	Use:   "2fa-status",
	Short: "Check 2FA status",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()

		status, err := cd2Client.Auth().Check2FAStatus(ctx)
		if err != nil {
			return err
		}
		return outputResult(status)
	},
}

var twoFASetupCmd = &cobra.Command{
	Use:   "2fa-setup [password]",
	Short: "Setup 2FA (generate TOTP secret)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		var password string
		if len(args) >= 1 {
			password = args[0]
		}
		password, err := readSensitiveValue(password, "CD2_CLI_PASSWORD", "password-file", cmd)
		if err != nil {
			return err
		}

		result, err := cd2Client.Auth().Setup2FA(ctx, &pb.Setup2FARequest{
			Password: password,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var twoFAEnableCmd = &cobra.Command{
	Use:   "2fa-enable [totp-code]",
	Short: "Enable 2FA with TOTP code",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		var totpCode string
		if len(args) >= 1 {
			totpCode = args[0]
		}
		totpCode, err := readSensitiveValue(totpCode, "CD2_CLI_TOTP", "totp-file", cmd)
		if err != nil {
			return err
		}

		result, err := cd2Client.Auth().Enable2FA(ctx, &pb.TwoFactorAuthCodeRequest{
			TotpCode: totpCode,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var twoFADisableCmd = &cobra.Command{
	Use:   "2fa-disable [totp-code]",
	Short: "Disable 2FA with TOTP code",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		var totpCode string
		if len(args) >= 1 {
			totpCode = args[0]
		}
		totpCode, err := readSensitiveValue(totpCode, "CD2_CLI_TOTP", "totp-file", cmd)
		if err != nil {
			return err
		}

		result, err := cd2Client.Auth().Disable2FA(ctx, &pb.TwoFactorAuthCodeRequest{
			TotpCode: totpCode,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var twoFARecoveryCodesCmd = &cobra.Command{
	Use:   "2fa-recovery-codes [totp-code]",
	Short: "Get recovery codes",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		var totpCode string
		if len(args) >= 1 {
			totpCode = args[0]
		}
		totpCode, err := readSensitiveValue(totpCode, "CD2_CLI_TOTP", "totp-file", cmd)
		if err != nil {
			return err
		}

		result, err := cd2Client.Auth().GetRecoveryCodes(ctx, &pb.TwoFactorAuthCodeRequest{
			TotpCode: totpCode,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var twoFARegenerateCodesCmd = &cobra.Command{
	Use:   "2fa-regenerate-codes [totp-code]",
	Short: "Regenerate recovery codes",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		var totpCode string
		if len(args) >= 1 {
			totpCode = args[0]
		}
		totpCode, err := readSensitiveValue(totpCode, "CD2_CLI_TOTP", "totp-file", cmd)
		if err != nil {
			return err
		}

		result, err := cd2Client.Auth().RegenerateRecoveryCodes(ctx, &pb.TwoFactorAuthCodeRequest{
			TotpCode: totpCode,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}
