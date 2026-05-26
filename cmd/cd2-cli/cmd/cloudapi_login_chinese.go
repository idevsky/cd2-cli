package cmd

import (
	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

func init() {
	cloudapiCmd.AddCommand(login115CookieCmd)
	cloudapiCmd.AddCommand(login115QrcodeCmd)
	cloudapiCmd.AddCommand(login115OpenOAuthCmd)
	cloudapiCmd.AddCommand(login115OpenQrcodeCmd)
	cloudapiCmd.AddCommand(loginAliyunOAuthCmd)
	cloudapiCmd.AddCommand(loginAliyunRefreshCmd)
	cloudapiCmd.AddCommand(loginAliyunQrcodeCmd)
	cloudapiCmd.AddCommand(loginBaiduOAuthCmd)
	cloudapiCmd.AddCommand(login189QrcodeCmd)

	loginAliyunRefreshCmd.Flags().Bool("use-open-api", false, "Use Open API")
	login115QrcodeCmd.Flags().String("platform", "", "Platform string")

	login115OpenOAuthCmd.Flags().String("refresh-token-file", "", "Read refresh token from file")
	loginAliyunOAuthCmd.Flags().String("refresh-token-file", "", "Read refresh token from file")
	loginAliyunRefreshCmd.Flags().String("refresh-token-file", "", "Read refresh token from file")
	loginBaiduOAuthCmd.Flags().String("refresh-token-file", "", "Read refresh token from file")

	setCommandID(login115CookieCmd, "cloudapi.login-115-cookie")
	setCommandID(login115QrcodeCmd, "cloudapi.login-115-qrcode")
	setCommandID(login115OpenOAuthCmd, "cloudapi.login-115open-oauth")
	setCommandID(login115OpenQrcodeCmd, "cloudapi.login-115open-qrcode")
	setCommandID(loginAliyunOAuthCmd, "cloudapi.login-aliyun-oauth")
	setCommandID(loginAliyunRefreshCmd, "cloudapi.login-aliyun-refresh")
	setCommandID(loginAliyunQrcodeCmd, "cloudapi.login-aliyun-qrcode")
	setCommandID(loginBaiduOAuthCmd, "cloudapi.login-baidu-oauth")
	setCommandID(login189QrcodeCmd, "cloudapi.login-189-qrcode")
}

var login115CookieCmd = &cobra.Command{
	Use:   "login-115-cookie <cookie>",
	Short: "Login to 115 with editthiscookie string",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		cookie := args[0]

		result, err := cd2Client.CloudAPI().APILogin115Editthiscookie(ctx, &pb.Login115EditthiscookieRequest{
			EditThiscookieString: cookie,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var login115QrcodeCmd = &cobra.Command{
	Use:   "login-115-qrcode",
	Short: "Login to 115 with QR code",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		platform, _ := cmd.Flags().GetString("platform")

		result, err := cd2Client.CloudAPI().APILogin115QRCode(ctx, &pb.Login115QrCodeRequest{
			PlatformString: &platform,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var login115OpenOAuthCmd = &cobra.Command{
	Use:   "login-115open-oauth [refresh-token]",
	Short: "Login to 115 Open with OAuth",
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

		result, err := cd2Client.CloudAPI().APILogin115OpenOAuth(ctx, &pb.Login115OpenOAuthRequest{
			RefreshToken: refreshToken,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var login115OpenQrcodeCmd = &cobra.Command{
	Use:   "login-115open-qrcode",
	Short: "Login to 115 Open with QR code",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()

		result, err := cd2Client.CloudAPI().APILogin115OpenQRCode(ctx, &pb.Login115OpenQRCodeRequest{})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var loginAliyunOAuthCmd = &cobra.Command{
	Use:   "login-aliyun-oauth [refresh-token]",
	Short: "Login to Aliyun with OAuth",
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

		result, err := cd2Client.CloudAPI().APILoginAliyundriveOAuth(ctx, &pb.LoginAliyundriveOAuthRequest{
			RefreshToken: refreshToken,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var loginAliyunRefreshCmd = &cobra.Command{
	Use:   "login-aliyun-refresh [refresh-token]",
	Short: "Login to Aliyun with refresh token",
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
		useOpenAPI, _ := cmd.Flags().GetBool("use-open-api")

		result, err := cd2Client.CloudAPI().APILoginAliyundriveRefreshtoken(ctx, &pb.LoginAliyundriveRefreshtokenRequest{
			RefreshToken: refreshToken,
			UseOpenAPI:   useOpenAPI,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var loginAliyunQrcodeCmd = &cobra.Command{
	Use:   "login-aliyun-qrcode",
	Short: "Login to Aliyun with QR code",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()

		result, err := cd2Client.CloudAPI().APILoginAliyunDriveQRCode(ctx, &pb.LoginAliyundriveQRCodeRequest{})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var loginBaiduOAuthCmd = &cobra.Command{
	Use:   "login-baidu-oauth [refresh-token]",
	Short: "Login to Baidu with OAuth",
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

		result, err := cd2Client.CloudAPI().APILoginBaiduPanOAuth(ctx, &pb.LoginBaiduPanOAuthRequest{
			RefreshToken: refreshToken,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var login189QrcodeCmd = &cobra.Command{
	Use:   "login-189-qrcode",
	Short: "Login to 189 with QR code",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()

		result, err := cd2Client.CloudAPI().APILogin189QRCode(ctx, &pb.Login189QRCodeRequest{})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}
