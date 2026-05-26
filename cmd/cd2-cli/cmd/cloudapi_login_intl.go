package cmd

import (
	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

func init() {
	cloudapiCmd.AddCommand(loginOnedriveOAuthCmd)
	cloudapiCmd.AddCommand(loginGoogleOAuthCmd)
	cloudapiCmd.AddCommand(loginGoogleRefreshCmd)
	cloudapiCmd.AddCommand(loginXunleiOAuthCmd)
	cloudapiCmd.AddCommand(loginXunleiOpenOAuthCmd)
	cloudapiCmd.AddCommand(login123panOAuthCmd)

	loginOnedriveOAuthCmd.Flags().String("refresh-token-file", "", "Read refresh token from file")
	loginGoogleOAuthCmd.Flags().String("refresh-token-file", "", "Read refresh token from file")
	loginGoogleRefreshCmd.Flags().String("refresh-token-file", "", "Read refresh token from file")
	loginXunleiOAuthCmd.Flags().String("refresh-token-file", "", "Read refresh token from file")
	loginXunleiOpenOAuthCmd.Flags().String("refresh-token-file", "", "Read refresh token from file")
	login123panOAuthCmd.Flags().String("refresh-token-file", "", "Read refresh token from file")

	setCommandID(loginOnedriveOAuthCmd, "cloudapi.login-onedrive-oauth")
	setCommandID(loginGoogleOAuthCmd, "cloudapi.login-google-oauth")
	setCommandID(loginGoogleRefreshCmd, "cloudapi.login-google-refresh")
	setCommandID(loginXunleiOAuthCmd, "cloudapi.login-xunlei-oauth")
	setCommandID(loginXunleiOpenOAuthCmd, "cloudapi.login-xunleiopen-oauth")
	setCommandID(login123panOAuthCmd, "cloudapi.login-123pan-oauth")
}

var loginOnedriveOAuthCmd = &cobra.Command{
	Use:   "login-onedrive-oauth [refresh-token]",
	Short: "Login to OneDrive with OAuth",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		var refreshToken string
		if len(args) >= 1 {
			refreshToken = args[0]
		}
		refreshToken, err := readSensitiveValue(refreshToken, "CD2_CLI_REFRESH_TOKEN", "refresh-token-file", cmd)
		if err != nil {
			return err
		}

		result, err := cd2Client.CloudAPI().APILoginOneDriveOAuth(ctx, &pb.LoginOneDriveOAuthRequest{
			RefreshToken: refreshToken,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var loginGoogleOAuthCmd = &cobra.Command{
	Use:   "login-google-oauth [refresh-token]",
	Short: "Login to Google Drive with OAuth",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		var refreshToken string
		if len(args) >= 1 {
			refreshToken = args[0]
		}
		refreshToken, err := readSensitiveValue(refreshToken, "CD2_CLI_REFRESH_TOKEN", "refresh-token-file", cmd)
		if err != nil {
			return err
		}

		result, err := cd2Client.CloudAPI().ApiLoginGoogleDriveOAuth(ctx, &pb.LoginGoogleDriveOAuthRequest{
			RefreshToken: refreshToken,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var loginGoogleRefreshCmd = &cobra.Command{
	Use:   "login-google-refresh [refresh-token]",
	Short: "Login to Google Drive with refresh token",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		var refreshToken string
		if len(args) >= 1 {
			refreshToken = args[0]
		}
		refreshToken, err := readSensitiveValue(refreshToken, "CD2_CLI_REFRESH_TOKEN", "refresh-token-file", cmd)
		if err != nil {
			return err
		}

		result, err := cd2Client.CloudAPI().ApiLoginGoogleDriveRefreshToken(ctx, &pb.LoginGoogleDriveRefreshTokenRequest{
			RefreshToken: refreshToken,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var loginXunleiOAuthCmd = &cobra.Command{
	Use:   "login-xunlei-oauth [refresh-token]",
	Short: "Login to Xunlei with OAuth",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		var refreshToken string
		if len(args) >= 1 {
			refreshToken = args[0]
		}
		refreshToken, err := readSensitiveValue(refreshToken, "CD2_CLI_REFRESH_TOKEN", "refresh-token-file", cmd)
		if err != nil {
			return err
		}

		result, err := cd2Client.CloudAPI().ApiLoginXunleiOAuth(ctx, &pb.LoginXunleiOAuthRequest{
			RefreshToken: refreshToken,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var loginXunleiOpenOAuthCmd = &cobra.Command{
	Use:   "login-xunleiopen-oauth [refresh-token]",
	Short: "Login to Xunlei Open with OAuth",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		var refreshToken string
		if len(args) >= 1 {
			refreshToken = args[0]
		}
		refreshToken, err := readSensitiveValue(refreshToken, "CD2_CLI_REFRESH_TOKEN", "refresh-token-file", cmd)
		if err != nil {
			return err
		}

		result, err := cd2Client.CloudAPI().ApiLoginXunleiOpenOAuth(ctx, &pb.LoginXunleiOpenOAuthRequest{
			RefreshToken: refreshToken,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var login123panOAuthCmd = &cobra.Command{
	Use:   "login-123pan-oauth [refresh-token]",
	Short: "Login to 123Pan with OAuth",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		var refreshToken string
		if len(args) >= 1 {
			refreshToken = args[0]
		}
		refreshToken, err := readSensitiveValue(refreshToken, "CD2_CLI_REFRESH_TOKEN", "refresh-token-file", cmd)
		if err != nil {
			return err
		}

		result, err := cd2Client.CloudAPI().ApiLogin123PanOAuth(ctx, &pb.Login123PanOAuthRequest{
			RefreshToken: refreshToken,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}
