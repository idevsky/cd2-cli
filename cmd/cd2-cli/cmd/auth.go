package cmd

import (
	"fmt"
	"os"
	"strings"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication operations",
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(login2faCmd)
	authCmd.AddCommand(loginThirdPartyCmd)
	authCmd.AddCommand(logoutCmd)
	authCmd.AddCommand(statusCmd)
	authCmd.AddCommand(changePasswordCmd)
	authCmd.AddCommand(changeEmailCmd)
	authCmd.AddCommand(confirmEmailCmd)
	authCmd.AddCommand(sendConfirmEmailCmd)
	authCmd.AddCommand(sendChangeEmailCodeCmd)
	authCmd.AddCommand(changeEmailAndPasswordCmd)
	authCmd.AddCommand(registerCmd)
	authCmd.AddCommand(sendResetEmailCmd)
	authCmd.AddCommand(resetAccountCmd)

	loginCmd.Flags().Bool("sync", false, "Sync data to cloud")
	loginCmd.Flags().Bool("save", false, "Save token to config file after successful login")
	loginCmd.Flags().String("password-file", "", "Read password from file")
	login2faCmd.Flags().Bool("sync", false, "Sync data to cloud")
	login2faCmd.Flags().Bool("save", false, "Save token to config file after successful login")
	login2faCmd.Flags().String("password-file", "", "Read password from file")
	login2faCmd.Flags().String("totp-file", "", "Read TOTP code from file")
	loginThirdPartyCmd.Flags().Bool("sync", false, "Sync data to cloud")
	loginThirdPartyCmd.Flags().Bool("save", false, "Save token to config file after successful login")
	loginThirdPartyCmd.Flags().String("access-token", "", "Access token for OAuth login")
	loginThirdPartyCmd.Flags().Uint64("expires-in", 0, "Token expiration time in seconds")
	logoutCmd.Flags().Bool("cloudfs", false, "Logout from CloudFS server too")
	changePasswordCmd.Flags().String("totp", "", "TOTP code for 2FA-enabled accounts")
	changePasswordCmd.Flags().String("old-password-file", "", "Read old password from file")
	changePasswordCmd.Flags().String("new-password-file", "", "Read new password from file")
	changeEmailCmd.Flags().String("totp", "", "TOTP code for 2FA-enabled accounts")
	changeEmailCmd.Flags().String("code", "", "Change email verification code")
	sendChangeEmailCodeCmd.Flags().String("password", "", "Account password")
	registerCmd.Flags().Bool("sync", false, "Sync data to cloud (unused)")
	changeEmailAndPasswordCmd.Flags().Bool("sync", false, "Sync user data with cloud")

	setCommandID(loginCmd, "auth.login")
	setCommandID(login2faCmd, "auth.login-2fa")
	setCommandID(loginThirdPartyCmd, "auth.login-third-party")
	setCommandID(logoutCmd, "auth.logout")
	setCommandID(statusCmd, "auth.status")
	setCommandID(changePasswordCmd, "auth.change-password")
	setCommandID(changeEmailCmd, "auth.change-email")
	setCommandID(confirmEmailCmd, "auth.confirm-email")
	setCommandID(sendConfirmEmailCmd, "auth.send-confirm-email")
	setCommandID(sendChangeEmailCodeCmd, "auth.send-change-email-code")
	setCommandID(changeEmailAndPasswordCmd, "auth.change-email-and-password")
	setCommandID(registerCmd, "auth.register")
	setCommandID(sendResetEmailCmd, "auth.send-reset-email")
	setCommandID(resetAccountCmd, "auth.reset-account")
}

func readSensitiveValue(value string, envVar string, fileFlag string, cmd *cobra.Command) (string, error) {
	if value != "" {
		return strings.TrimSpace(value), nil
	}

	if envVar != "" {
		envVal := os.Getenv(envVar)
		if envVal != "" {
			return strings.TrimSpace(envVal), nil
		}
	}

	if fileFlag != "" {
		filePath, err := cmd.Flags().GetString(fileFlag)
		if err != nil {
			return "", err
		}
		if filePath != "" {
			data, err := os.ReadFile(filePath)
			if err != nil {
				return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
			}
			return strings.TrimSpace(string(data)), nil
		}
	}

	return "", fmt.Errorf("required value not provided (argument, %s, or %s flag)", envVar, fileFlag)
}

var loginCmd = &cobra.Command{
	Use:   "login [username] [password]",
	Short: "Login to CloudDrive",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		username := args[0]
		var password string
		if len(args) >= 2 {
			password = args[1]
		}
		password, err := readSensitiveValue(password, "CD2_CLI_PASSWORD", "password-file", cmd)
		if err != nil {
			return err
		}
		saveToken, _ := cmd.Flags().GetBool("save")

		resp, err := cd2Client.Public().GetToken(ctx, &pb.GetTokenRequest{
			UserName: username,
			Password: password,
		})
		if err != nil {
			return err
		}

		if saveToken && resp.Token != "" {
			configPath := getConfigPath()
			if err := saveTokenToConfig(configPath, resp.Token); err != nil {
				return fmt.Errorf("login succeeded but failed to save token: %w", err)
			}
		}

		return outputResult(resp)
	},
}

var login2faCmd = &cobra.Command{
	Use:   "login-2fa [username] [password] [totp-code]",
	Short: "Login with 2FA code",
	Args:  cobra.RangeArgs(1, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		username := args[0]
		var password string
		var totpCode string
		if len(args) >= 2 {
			password = args[1]
		}
		if len(args) >= 3 {
			totpCode = args[2]
		}
		password, err := readSensitiveValue(password, "CD2_CLI_PASSWORD", "password-file", cmd)
		if err != nil {
			return err
		}
		totpCode, err = readSensitiveValue(totpCode, "CD2_CLI_TOTP", "totp-file", cmd)
		if err != nil {
			return err
		}
		sync, _ := cmd.Flags().GetBool("sync")
		saveToken, _ := cmd.Flags().GetBool("save")

		resp, err := cd2Client.Public().LoginWith2FA(ctx, &pb.LoginWith2FARequest{
			UserName:       username,
			Password:       password,
			TotpCode:       totpCode,
			SynDataToCloud: sync,
		})
		if err != nil {
			return err
		}

		if saveToken && resp.Token != "" {
			configPath := getConfigPath()
			if err := saveTokenToConfig(configPath, resp.Token); err != nil {
				return fmt.Errorf("login succeeded but failed to save token: %w", err)
			}
		}

		return outputResult(resp)
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from CloudDrive",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		fromCloudFS, _ := cmd.Flags().GetBool("cloudfs")

		result, err := cd2Client.Auth().Logout(ctx, &pb.UserLogoutRequest{
			LogoutFromCloudFS: fromCloudFS,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get account status",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()

		status, err := cd2Client.Auth().GetAccountStatus(ctx)
		if err != nil {
			return err
		}
		return outputResult(status)
	},
}

var changePasswordCmd = &cobra.Command{
	Use:   "change-password [old-password] [new-password]",
	Short: "Change account password",
	Args:  cobra.RangeArgs(0, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		var oldPassword string
		var newPassword string
		if len(args) >= 1 {
			oldPassword = args[0]
		}
		if len(args) >= 2 {
			newPassword = args[1]
		}
		oldPassword, err := readSensitiveValue(oldPassword, "CD2_CLI_OLD_PASSWORD", "old-password-file", cmd)
		if err != nil {
			return err
		}
		newPassword, err = readSensitiveValue(newPassword, "CD2_CLI_NEW_PASSWORD", "new-password-file", cmd)
		if err != nil {
			return err
		}
		totpCode, _ := cmd.Flags().GetString("totp")

		result, err := cd2Client.Auth().ChangePassword(ctx, &pb.ChangePasswordRequest{
			OldPassword: oldPassword,
			NewPassword: newPassword,
			TotpCode:    &totpCode,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var changeEmailCmd = &cobra.Command{
	Use:   "change-email [new-email] [password]",
	Short: "Change account email",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		newEmail := args[0]
		password := args[1]
		totpCode, _ := cmd.Flags().GetString("totp")
		changeCode, _ := cmd.Flags().GetString("code")

		err := cd2Client.Auth().ChangeEmail(ctx, &pb.ChangeEmailRequest{
			NewEmail:   newEmail,
			Password:   password,
			TotpCode:   &totpCode,
			ChangeCode: &changeCode,
		})
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "success"})
	},
}

var confirmEmailCmd = &cobra.Command{
	Use:   "confirm-email [code]",
	Short: "Confirm email with verification code",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		code := args[0]

		err := cd2Client.Auth().ConfirmEmail(ctx, &pb.ConfirmEmailRequest{
			ConfirmCode: code,
		})
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "success"})
	},
}

var registerCmd = &cobra.Command{
	Use:   "register [username] [password]",
	Short: "Register new account",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		username := args[0]
		password := args[1]

		result, err := cd2Client.Public().Register(ctx, &pb.UserRegisterRequest{
			UserName: username,
			Password: password,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var loginThirdPartyCmd = &cobra.Command{
	Use:   "login-third-party [cloud-name] [refresh-token]",
	Short: "Login with third-party account OAuth tokens",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		cloudName := args[0]
		refreshToken := args[1]
		accessToken, _ := cmd.Flags().GetString("access-token")
		expiresIn, _ := cmd.Flags().GetUint64("expires-in")
		sync, _ := cmd.Flags().GetBool("sync")
		saveToken, _ := cmd.Flags().GetBool("save")

		resp, err := cd2Client.Public().LoginWithThirdPartyAccount(ctx, &pb.LoginWithThirdPartyAccountRequest{
			CloudName:      cloudName,
			RefreshToken:   refreshToken,
			AccessToken:    accessToken,
			ExpiresIn:      expiresIn,
			SynDataToCloud: sync,
		})
		if err != nil {
			return err
		}

		if saveToken && resp.Token != "" {
			configPath := getConfigPath()
			if err := saveTokenToConfig(configPath, resp.Token); err != nil {
				return fmt.Errorf("login succeeded but failed to save token: %w", err)
			}
		}

		return outputResult(resp)
	},
}

var sendConfirmEmailCmd = &cobra.Command{
	Use:   "send-confirm-email",
	Short: "Request email confirmation code",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()

		err := cd2Client.Auth().SendConfirmEmail(ctx)
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "success", "message": "confirmation email sent"})
	},
}

var sendChangeEmailCodeCmd = &cobra.Command{
	Use:   "send-change-email-code [new-email]",
	Short: "Request change email verification code",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		newEmail := args[0]
		password, _ := cmd.Flags().GetString("password")

		err := cd2Client.Auth().SendChangeEmailCode(ctx, &pb.SendChangeEmailCodeRequest{
			NewEmail: newEmail,
			Password: password,
		})
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "success", "message": "change email code sent"})
	},
}

var changeEmailAndPasswordCmd = &cobra.Command{
	Use:   "change-email-and-password [new-email] [new-password]",
	Short: "Change email and password (trusted devices only)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		newEmail := args[0]
		newPassword := args[1]
		sync, _ := cmd.Flags().GetBool("sync")

		err := cd2Client.Auth().ChangeEmailAndPassword(ctx, &pb.ChangeEmailAndPasswordRequest{
			NewEmail:              newEmail,
			NewPassword:           newPassword,
			SyncUserDataWithCloud: sync,
		})
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "success"})
	},
}

var sendResetEmailCmd = &cobra.Command{
	Use:   "send-reset-email [email]",
	Short: "Request account reset email",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		email := args[0]

		err := cd2Client.Public().SendResetAccountEmail(ctx, &pb.SendResetAccountEmailRequest{
			Email: email,
		})
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "success", "message": "reset email sent"})
	},
}

var resetAccountCmd = &cobra.Command{
	Use:   "reset-account [reset-code] [new-password]",
	Short: "Reset account with reset code from email",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		resetCode := args[0]
		newPassword := args[1]

		err := cd2Client.Public().ResetAccount(ctx, &pb.ResetAccountRequest{
			ResetCode:   resetCode,
			NewPassword: newPassword,
		})
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "success", "message": "account reset"})
	},
}
