package client

import (
	"context"
	"io"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

type CloudAPIAPI struct {
	c *Client
}

func (a *CloudAPIAPI) CanAddMoreCloudApis(ctx context.Context) (*pb.FileOperationResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.CanAddMoreCloudApis(ctx, &emptypb.Empty{})
}

func (a *CloudAPIAPI) APILogin115Editthiscookie(ctx context.Context, req *pb.Login115EditthiscookieRequest) (*pb.APILoginResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.APILogin115Editthiscookie(ctx, req)
}

func (a *CloudAPIAPI) APILogin115QRCode(ctx context.Context, req *pb.Login115QrCodeRequest) ([]*pb.QRCodeScanMessage, error) {
	ctx = a.c.withAuth(ctx)
	stream, err := a.c.client.APILogin115QRCode(ctx, req)
	if err != nil {
		return nil, err
	}
	var results []*pb.QRCodeScanMessage
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		results = append(results, msg)
	}
	return results, nil
}

func (a *CloudAPIAPI) APILogin115OpenOAuth(ctx context.Context, req *pb.Login115OpenOAuthRequest) (*pb.APILoginResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.APILogin115OpenOAuth(ctx, req)
}

func (a *CloudAPIAPI) APILogin115OpenQRCode(ctx context.Context, req *pb.Login115OpenQRCodeRequest) ([]*pb.QRCodeScanMessage, error) {
	ctx = a.c.withAuth(ctx)
	stream, err := a.c.client.APILogin115OpenQRCode(ctx, req)
	if err != nil {
		return nil, err
	}
	var results []*pb.QRCodeScanMessage
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		results = append(results, msg)
	}
	return results, nil
}

func (a *CloudAPIAPI) APILoginAliyundriveOAuth(ctx context.Context, req *pb.LoginAliyundriveOAuthRequest) (*pb.APILoginResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.APILoginAliyundriveOAuth(ctx, req)
}

func (a *CloudAPIAPI) APILoginAliyundriveRefreshtoken(ctx context.Context, req *pb.LoginAliyundriveRefreshtokenRequest) (*pb.APILoginResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.APILoginAliyundriveRefreshtoken(ctx, req)
}

func (a *CloudAPIAPI) APILoginAliyunDriveQRCode(ctx context.Context, req *pb.LoginAliyundriveQRCodeRequest) ([]*pb.QRCodeScanMessage, error) {
	ctx = a.c.withAuth(ctx)
	stream, err := a.c.client.APILoginAliyunDriveQRCode(ctx, req)
	if err != nil {
		return nil, err
	}
	var results []*pb.QRCodeScanMessage
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		results = append(results, msg)
	}
	return results, nil
}

func (a *CloudAPIAPI) APILoginBaiduPanOAuth(ctx context.Context, req *pb.LoginBaiduPanOAuthRequest) (*pb.APILoginResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.APILoginBaiduPanOAuth(ctx, req)
}

func (a *CloudAPIAPI) APILoginOneDriveOAuth(ctx context.Context, req *pb.LoginOneDriveOAuthRequest) (*pb.APILoginResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.APILoginOneDriveOAuth(ctx, req)
}

func (a *CloudAPIAPI) ApiLoginGoogleDriveOAuth(ctx context.Context, req *pb.LoginGoogleDriveOAuthRequest) (*pb.APILoginResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.ApiLoginGoogleDriveOAuth(ctx, req)
}

func (a *CloudAPIAPI) ApiLoginGoogleDriveRefreshToken(ctx context.Context, req *pb.LoginGoogleDriveRefreshTokenRequest) (*pb.APILoginResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.ApiLoginGoogleDriveRefreshToken(ctx, req)
}

func (a *CloudAPIAPI) ApiLoginXunleiOAuth(ctx context.Context, req *pb.LoginXunleiOAuthRequest) (*pb.APILoginResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.ApiLoginXunleiOAuth(ctx, req)
}

func (a *CloudAPIAPI) ApiLoginXunleiOpenOAuth(ctx context.Context, req *pb.LoginXunleiOpenOAuthRequest) (*pb.APILoginResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.ApiLoginXunleiOpenOAuth(ctx, req)
}

func (a *CloudAPIAPI) ApiLogin123PanOAuth(ctx context.Context, req *pb.Login123PanOAuthRequest) (*pb.APILoginResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.ApiLogin123PanOAuth(ctx, req)
}

func (a *CloudAPIAPI) APILogin189QRCode(ctx context.Context, req *pb.Login189QRCodeRequest) ([]*pb.QRCodeScanMessage, error) {
	ctx = a.c.withAuth(ctx)
	stream, err := a.c.client.APILogin189QRCode(ctx, req)
	if err != nil {
		return nil, err
	}
	var results []*pb.QRCodeScanMessage
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		results = append(results, msg)
	}
	return results, nil
}

func (a *CloudAPIAPI) APILoginWebDav(ctx context.Context, req *pb.LoginWebDavRequest) (*pb.APILoginResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.APILoginWebDav(ctx, req)
}

func (a *CloudAPIAPI) APILoginS3(ctx context.Context, req *pb.LoginS3Request) (*pb.APILoginResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.APILoginS3(ctx, req)
}

func (a *CloudAPIAPI) APIAddLocalFolder(ctx context.Context, req *pb.AddLocalFolderRequest) (*pb.APILoginResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.APIAddLocalFolder(ctx, req)
}

func (a *CloudAPIAPI) APILoginCloudDrive(ctx context.Context, req *pb.LoginCloudDriveRequest) (*pb.APILoginResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.APILoginCloudDrive(ctx, req)
}

func (a *CloudAPIAPI) APILoginSftp(ctx context.Context, req *pb.LoginSftpRequest) (*pb.APILoginResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.APILoginSftp(ctx, req)
}

func (a *CloudAPIAPI) APILoginFtp(ctx context.Context, req *pb.LoginFtpRequest) (*pb.APILoginResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.APILoginFtp(ctx, req)
}

func (a *CloudAPIAPI) APILoginSmb(ctx context.Context, req *pb.LoginSmbRequest) (*pb.APILoginResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.APILoginSmb(ctx, req)
}

func (a *CloudAPIAPI) DiscoverSmbServers(ctx context.Context) (*pb.DiscoverSmbServersResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.DiscoverSmbServers(ctx, &emptypb.Empty{})
}

func (a *CloudAPIAPI) DiscoverSmbShares(ctx context.Context, req *pb.DiscoverSmbSharesRequest) (*pb.DiscoverSmbSharesResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.DiscoverSmbShares(ctx, req)
}

func (a *CloudAPIAPI) RemoveCloudAPI(ctx context.Context, req *pb.RemoveCloudAPIRequest) (*pb.FileOperationResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.RemoveCloudAPI(ctx, req)
}

func (a *CloudAPIAPI) GetAllCloudApis(ctx context.Context) (*pb.CloudAPIList, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetAllCloudApis(ctx, &emptypb.Empty{})
}

func (a *CloudAPIAPI) GetCloudAPIConfig(ctx context.Context, req *pb.GetCloudAPIConfigRequest) (*pb.CloudAPIConfig, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetCloudAPIConfig(ctx, req)
}

func (a *CloudAPIAPI) SetCloudAPIConfig(ctx context.Context, req *pb.SetCloudAPIConfigRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.SetCloudAPIConfig(ctx, req)
	return err
}
