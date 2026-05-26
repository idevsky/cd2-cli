package client

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

const bufSize = 1024 * 1024

type mockServer struct {
	pb.UnimplementedCloudDriveFileSrvServer
	closeFileCalled bool
	closeFileMutex  sync.Mutex
}

func (s *mockServer) GetSystemInfo(ctx context.Context, req *emptypb.Empty) (*pb.CloudDriveSystemInfo, error) {
	return &pb.CloudDriveSystemInfo{IsLogin: true, UserName: "test-user"}, nil
}

func (s *mockServer) GetToken(ctx context.Context, req *pb.GetTokenRequest) (*pb.JWTToken, error) {
	return &pb.JWTToken{Success: true, Token: "test-token"}, nil
}

func (s *mockServer) Login(ctx context.Context, req *pb.UserLoginRequest) (*pb.FileOperationResult, error) {
	return &pb.FileOperationResult{Success: true}, nil
}

func (s *mockServer) LoginWithThirdPartyAccount(ctx context.Context, req *pb.LoginWithThirdPartyAccountRequest) (*pb.JWTToken, error) {
	return &pb.JWTToken{Success: true, Token: "third-party-token"}, nil
}

func (s *mockServer) Register(ctx context.Context, req *pb.UserRegisterRequest) (*pb.FileOperationResult, error) {
	return &pb.FileOperationResult{Success: true}, nil
}

func (s *mockServer) SendResetAccountEmail(ctx context.Context, req *pb.SendResetAccountEmailRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *mockServer) ResetAccount(ctx context.Context, req *pb.ResetAccountRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *mockServer) GetApiTokenInfo(ctx context.Context, req *pb.StringValue) (*pb.TokenInfo, error) {
	return &pb.TokenInfo{Token: req.Value}, nil
}

func (s *mockServer) LoginWith2FA(ctx context.Context, req *pb.LoginWith2FARequest) (*pb.JWTToken, error) {
	return &pb.JWTToken{Success: true, Token: "2fa-token"}, nil
}

func (s *mockServer) SendConfirmEmail(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *mockServer) ConfirmEmail(ctx context.Context, req *pb.ConfirmEmailRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *mockServer) GetAccountStatus(ctx context.Context, req *emptypb.Empty) (*pb.AccountStatusResult, error) {
	return &pb.AccountStatusResult{UserName: "test-user"}, nil
}

func (s *mockServer) Logout(ctx context.Context, req *pb.UserLogoutRequest) (*pb.FileOperationResult, error) {
	return &pb.FileOperationResult{Success: true}, nil
}

func (s *mockServer) Check2FAStatus(ctx context.Context, req *emptypb.Empty) (*pb.TwoFactorAuthStatusResult, error) {
	return &pb.TwoFactorAuthStatusResult{TwoFactorEnabled: false}, nil
}

func (s *mockServer) Setup2FA(ctx context.Context, req *pb.Setup2FARequest) (*pb.TwoFactorAuthSetupResult, error) {
	return &pb.TwoFactorAuthSetupResult{Secret: "test-secret"}, nil
}

func (s *mockServer) Enable2FA(ctx context.Context, req *pb.TwoFactorAuthCodeRequest) (*pb.TwoFactorAuthEnableResult, error) {
	return &pb.TwoFactorAuthEnableResult{RecoveryCodes: []string{"code1", "code2"}}, nil
}

func (s *mockServer) Disable2FA(ctx context.Context, req *pb.TwoFactorAuthCodeRequest) (*pb.TwoFactorAuthMessageResult, error) {
	return &pb.TwoFactorAuthMessageResult{Message: "2FA disabled"}, nil
}

func (s *mockServer) GetRecoveryCodes(ctx context.Context, req *pb.TwoFactorAuthCodeRequest) (*pb.TwoFactorAuthRecoveryCodesResult, error) {
	return &pb.TwoFactorAuthRecoveryCodesResult{RecoveryCodes: []string{"code1", "code2"}, Total: 2}, nil
}

func (s *mockServer) RegenerateRecoveryCodes(ctx context.Context, req *pb.TwoFactorAuthCodeRequest) (*pb.TwoFactorAuthRecoveryCodesResult, error) {
	return &pb.TwoFactorAuthRecoveryCodesResult{RecoveryCodes: []string{"new1", "new2"}, Total: 2}, nil
}

func (s *mockServer) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.FileOperationResult, error) {
	return &pb.FileOperationResult{Success: true}, nil
}

func (s *mockServer) SendChangeEmailCode(ctx context.Context, req *pb.SendChangeEmailCodeRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *mockServer) ChangeEmail(ctx context.Context, req *pb.ChangeEmailRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *mockServer) ChangeEmailAndPassword(ctx context.Context, req *pb.ChangeEmailAndPasswordRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *mockServer) FindFileByPath(ctx context.Context, req *pb.FindFileByPathRequest) (*pb.CloudDriveFile, error) {
	return &pb.CloudDriveFile{FullPathName: req.Path, Name: "test.txt"}, nil
}

func (s *mockServer) CreateFolder(ctx context.Context, req *pb.CreateFolderRequest) (*pb.CreateFolderResult, error) {
	return &pb.CreateFolderResult{Result: &pb.FileOperationResult{Success: true}}, nil
}

func (s *mockServer) CreateEncryptedFolder(ctx context.Context, req *pb.CreateEncryptedFolderRequest) (*pb.CreateFolderResult, error) {
	return &pb.CreateFolderResult{Result: &pb.FileOperationResult{Success: true}}, nil
}

func (s *mockServer) UnlockEncryptedFile(ctx context.Context, req *pb.UnlockEncryptedFileRequest) (*pb.FileOperationResult, error) {
	return &pb.FileOperationResult{Success: true}, nil
}

func (s *mockServer) LockEncryptedFile(ctx context.Context, req *pb.FileRequest) (*pb.FileOperationResult, error) {
	return &pb.FileOperationResult{Success: true}, nil
}

func (s *mockServer) RenameFile(ctx context.Context, req *pb.RenameFileRequest) (*pb.FileOperationResult, error) {
	return &pb.FileOperationResult{Success: true}, nil
}

func (s *mockServer) RenameFiles(ctx context.Context, req *pb.RenameFilesRequest) (*pb.FileOperationResult, error) {
	return &pb.FileOperationResult{Success: true}, nil
}

func (s *mockServer) MoveFile(ctx context.Context, req *pb.MoveFileRequest) (*pb.FileOperationResult, error) {
	return &pb.FileOperationResult{Success: true}, nil
}

func (s *mockServer) CopyFile(ctx context.Context, req *pb.CopyFileRequest) (*pb.FileOperationResult, error) {
	return &pb.FileOperationResult{Success: true}, nil
}

func (s *mockServer) DeleteFile(ctx context.Context, req *pb.FileRequest) (*pb.FileOperationResult, error) {
	return &pb.FileOperationResult{Success: true}, nil
}

func (s *mockServer) DeleteFilePermanently(ctx context.Context, req *pb.FileRequest) (*pb.FileOperationResult, error) {
	return &pb.FileOperationResult{Success: true}, nil
}

func (s *mockServer) DeleteFiles(ctx context.Context, req *pb.MultiFileRequest) (*pb.FileOperationResult, error) {
	return &pb.FileOperationResult{Success: true}, nil
}

func (s *mockServer) DeleteFilesPermanently(ctx context.Context, req *pb.MultiFileRequest) (*pb.FileOperationResult, error) {
	return &pb.FileOperationResult{Success: true}, nil
}

func (s *mockServer) AddSharedLink(ctx context.Context, req *pb.AddSharedLinkRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *mockServer) GetFileDetailProperties(ctx context.Context, req *pb.FileRequest) (*pb.FileDetailProperties, error) {
	return &pb.FileDetailProperties{}, nil
}

func (s *mockServer) GetSpaceInfo(ctx context.Context, req *pb.FileRequest) (*pb.SpaceInfo, error) {
	return &pb.SpaceInfo{}, nil
}

func (s *mockServer) GetCloudMemberships(ctx context.Context, req *pb.FileRequest) (*pb.CloudMemberships, error) {
	return &pb.CloudMemberships{}, nil
}

func (s *mockServer) GetMetaData(ctx context.Context, req *pb.FileRequest) (*pb.FileMetaData, error) {
	return &pb.FileMetaData{}, nil
}

func (s *mockServer) GetOriginalPath(ctx context.Context, req *pb.FileRequest) (*pb.StringResult, error) {
	return &pb.StringResult{Result: "/original/path"}, nil
}

func (s *mockServer) CreateFile(ctx context.Context, req *pb.CreateFileRequest) (*pb.CreateFileResult, error) {
	return &pb.CreateFileResult{FileHandle: 12345}, nil
}

func (s *mockServer) CloseFile(ctx context.Context, req *pb.CloseFileRequest) (*pb.FileOperationResult, error) {
	s.closeFileMutex.Lock()
	s.closeFileCalled = true
	s.closeFileMutex.Unlock()
	return &pb.FileOperationResult{Success: true}, nil
}

func (s *mockServer) WriteToFile(ctx context.Context, req *pb.WriteFileRequest) (*pb.WriteFileResult, error) {
	return &pb.WriteFileResult{BytesWritten: req.Length}, nil
}

func (s *mockServer) WriteToFileStream(stream pb.CloudDriveFileSrv_WriteToFileStreamServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.WriteFileResult{BytesWritten: 0})
		}
		if err != nil {
			return err
		}
		if req.CloseFile {
			return stream.SendAndClose(&pb.WriteFileResult{BytesWritten: req.Length})
		}
	}
}

func (s *mockServer) GetDownloadUrlPath(ctx context.Context, req *pb.GetDownloadUrlPathRequest) (*pb.DownloadUrlPathInfo, error) {
	return &pb.DownloadUrlPathInfo{DownloadUrlPath: "http://example.com/download"}, nil
}

func (s *mockServer) GetSubFiles(req *pb.ListSubFileRequest, stream pb.CloudDriveFileSrv_GetSubFilesServer) error {
	stream.Send(&pb.SubFilesReply{SubFiles: []*pb.CloudDriveFile{{FullPathName: "/test", Name: "file.txt"}}})
	return nil
}

func (s *mockServer) GetSearchResults(req *pb.SearchRequest, stream pb.CloudDriveFileSrv_GetSearchResultsServer) error {
	stream.Send(&pb.SubFilesReply{SubFiles: []*pb.CloudDriveFile{{FullPathName: "/search", Name: "result.txt"}}})
	return nil
}

func (s *mockServer) CanAddMoreMountPoints(ctx context.Context, req *emptypb.Empty) (*pb.FileOperationResult, error) {
	return &pb.FileOperationResult{Success: true}, nil
}

func (s *mockServer) GetMountPoints(ctx context.Context, req *emptypb.Empty) (*pb.GetMountPointsResult, error) {
	return &pb.GetMountPointsResult{}, nil
}

func (s *mockServer) AddMountPoint(ctx context.Context, req *pb.MountOption) (*pb.MountPointResult, error) {
	return &pb.MountPointResult{}, nil
}

func (s *mockServer) RemoveMountPoint(ctx context.Context, req *pb.MountPointRequest) (*pb.MountPointResult, error) {
	return &pb.MountPointResult{}, nil
}

func (s *mockServer) Mount(ctx context.Context, req *pb.MountPointRequest) (*pb.MountPointResult, error) {
	return &pb.MountPointResult{}, nil
}

func (s *mockServer) Unmount(ctx context.Context, req *pb.MountPointRequest) (*pb.MountPointResult, error) {
	return &pb.MountPointResult{}, nil
}

func (s *mockServer) UpdateMountPoint(ctx context.Context, req *pb.UpdateMountPointRequest) (*pb.MountPointResult, error) {
	return &pb.MountPointResult{}, nil
}

func (s *mockServer) GetAvailableDriveLetters(ctx context.Context, req *emptypb.Empty) (*pb.GetAvailableDriveLettersResult, error) {
	return &pb.GetAvailableDriveLettersResult{DriveLetters: []string{"A:", "B:"}}, nil
}

func (s *mockServer) HasDriveLetters(ctx context.Context, req *emptypb.Empty) (*pb.HasDriveLettersResult, error) {
	return &pb.HasDriveLettersResult{}, nil
}

func (s *mockServer) CanMountBothLocalAndCloud(ctx context.Context, req *emptypb.Empty) (*pb.BoolResult, error) {
	return &pb.BoolResult{Result: true}, nil
}

func (s *mockServer) GetAllTasksCount(ctx context.Context, req *emptypb.Empty) (*pb.GetAllTasksCountResult, error) {
	return &pb.GetAllTasksCountResult{}, nil
}

func (s *mockServer) GetDownloadFileCount(ctx context.Context, req *emptypb.Empty) (*pb.GetDownloadFileCountResult, error) {
	return &pb.GetDownloadFileCountResult{}, nil
}

func (s *mockServer) GetDownloadFileList(ctx context.Context, req *emptypb.Empty) (*pb.GetDownloadFileListResult, error) {
	return &pb.GetDownloadFileListResult{}, nil
}

func (s *mockServer) GetUploadFileCount(ctx context.Context, req *emptypb.Empty) (*pb.GetUploadFileCountResult, error) {
	return &pb.GetUploadFileCountResult{}, nil
}

func (s *mockServer) GetUploadFileList(ctx context.Context, req *pb.GetUploadFileListRequest) (*pb.GetUploadFileListResult, error) {
	return &pb.GetUploadFileListResult{}, nil
}

func (s *mockServer) CancelAllUploadFiles(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *mockServer) CancelUploadFiles(ctx context.Context, req *pb.MultpleUploadFileKeyRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *mockServer) PauseAllUploadFiles(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *mockServer) PauseUploadFiles(ctx context.Context, req *pb.MultpleUploadFileKeyRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *mockServer) ResumeAllUploadFiles(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *mockServer) ResumeUploadFiles(ctx context.Context, req *pb.MultpleUploadFileKeyRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *mockServer) CanAddMoreCloudApis(ctx context.Context, req *emptypb.Empty) (*pb.FileOperationResult, error) {
	return &pb.FileOperationResult{Success: true}, nil
}

func (s *mockServer) GetAllCloudApis(ctx context.Context, req *emptypb.Empty) (*pb.CloudAPIList, error) {
	return &pb.CloudAPIList{}, nil
}

func (s *mockServer) GetCloudAPIConfig(ctx context.Context, req *pb.GetCloudAPIConfigRequest) (*pb.CloudAPIConfig, error) {
	return &pb.CloudAPIConfig{}, nil
}

func (s *mockServer) SetCloudAPIConfig(ctx context.Context, req *pb.SetCloudAPIConfigRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *mockServer) RemoveCloudAPI(ctx context.Context, req *pb.RemoveCloudAPIRequest) (*pb.FileOperationResult, error) {
	return &pb.FileOperationResult{Success: true}, nil
}

func setupMockClient(t *testing.T) (*Client, func()) {
	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterCloudDriveFileSrvServer(s, &mockServer{})

	go func() {
		if err := s.Serve(lis); err != nil {
			t.Logf("Server error: %v", err)
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(dialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	client := &Client{
		conn:   conn,
		client: pb.NewCloudDriveFileSrvClient(conn),
	}

	return client, func() {
		client.Close()
		s.Stop()
		lis.Close()
	}
}

func TestPublicAPIWithMock(t *testing.T) {
	c, cleanup := setupMockClient(t)
	defer cleanup()

	t.Run("GetSystemInfo", func(t *testing.T) {
		info, err := c.Public().GetSystemInfo(context.Background())
		if err != nil {
			t.Errorf("GetSystemInfo error: %v", err)
		}
		if info.UserName != "test-user" {
			t.Errorf("expected userName 'test-user', got '%s'", info.UserName)
		}
	})

	t.Run("GetToken", func(t *testing.T) {
		token, err := c.Public().GetToken(context.Background(), &pb.GetTokenRequest{UserName: "test", Password: "pass"})
		if err != nil {
			t.Errorf("GetToken error: %v", err)
		}
		if token.Token != "test-token" {
			t.Errorf("expected token 'test-token', got '%s'", token.Token)
		}
		if !token.Success {
			t.Error("expected success")
		}
	})

	t.Run("Login", func(t *testing.T) {
		result, err := c.Public().Login(context.Background(), &pb.UserLoginRequest{UserName: "test", Password: "pass"})
		if err != nil {
			t.Errorf("Login error: %v", err)
		}
		if !result.Success {
			t.Error("expected success")
		}
	})

	t.Run("LoginWithThirdPartyAccount", func(t *testing.T) {
		token, err := c.Public().LoginWithThirdPartyAccount(context.Background(), &pb.LoginWithThirdPartyAccountRequest{})
		if err != nil {
			t.Errorf("LoginWithThirdPartyAccount error: %v", err)
		}
		if token.Token != "third-party-token" {
			t.Errorf("expected token 'third-party-token', got '%s'", token.Token)
		}
	})

	t.Run("Register", func(t *testing.T) {
		result, err := c.Public().Register(context.Background(), &pb.UserRegisterRequest{UserName: "test", Password: "pass"})
		if err != nil {
			t.Errorf("Register error: %v", err)
		}
		if !result.Success {
			t.Error("expected success")
		}
	})

	t.Run("SendResetAccountEmail", func(t *testing.T) {
		err := c.Public().SendResetAccountEmail(context.Background(), &pb.SendResetAccountEmailRequest{Email: "test@example.com"})
		if err != nil {
			t.Errorf("SendResetAccountEmail error: %v", err)
		}
	})

	t.Run("ResetAccount", func(t *testing.T) {
		err := c.Public().ResetAccount(context.Background(), &pb.ResetAccountRequest{ResetCode: "code123"})
		if err != nil {
			t.Errorf("ResetAccount error: %v", err)
		}
	})

	t.Run("GetApiTokenInfo", func(t *testing.T) {
		info, err := c.Public().GetApiTokenInfo(context.Background(), "my-token")
		if err != nil {
			t.Errorf("GetApiTokenInfo error: %v", err)
		}
		if info.Token != "my-token" {
			t.Errorf("expected token 'my-token', got '%s'", info.Token)
		}
	})

	t.Run("LoginWith2FA", func(t *testing.T) {
		token, err := c.Public().LoginWith2FA(context.Background(), &pb.LoginWith2FARequest{TotpCode: "123456"})
		if err != nil {
			t.Errorf("LoginWith2FA error: %v", err)
		}
		if token.Token != "2fa-token" {
			t.Errorf("expected token '2fa-token', got '%s'", token.Token)
		}
	})
}

func TestAuthAPIWithMock(t *testing.T) {
	c, cleanup := setupMockClient(t)
	defer cleanup()
	c.SetToken("test-token")

	t.Run("SendConfirmEmail", func(t *testing.T) {
		err := c.Auth().SendConfirmEmail(context.Background())
		if err != nil {
			t.Errorf("SendConfirmEmail error: %v", err)
		}
	})

	t.Run("ConfirmEmail", func(t *testing.T) {
		err := c.Auth().ConfirmEmail(context.Background(), &pb.ConfirmEmailRequest{ConfirmCode: "123"})
		if err != nil {
			t.Errorf("ConfirmEmail error: %v", err)
		}
	})

	t.Run("GetAccountStatus", func(t *testing.T) {
		status, err := c.Auth().GetAccountStatus(context.Background())
		if err != nil {
			t.Errorf("GetAccountStatus error: %v", err)
		}
		if status.UserName != "test-user" {
			t.Errorf("expected userName 'test-user', got '%s'", status.UserName)
		}
	})

	t.Run("Logout", func(t *testing.T) {
		result, err := c.Auth().Logout(context.Background(), &pb.UserLogoutRequest{LogoutFromCloudFS: true})
		if err != nil {
			t.Errorf("Logout error: %v", err)
		}
		if !result.Success {
			t.Error("expected success")
		}
	})

	t.Run("Check2FAStatus", func(t *testing.T) {
		status, err := c.Auth().Check2FAStatus(context.Background())
		if err != nil {
			t.Errorf("Check2FAStatus error: %v", err)
		}
		if status.TwoFactorEnabled {
			t.Error("expected 2FA to be disabled")
		}
	})

	t.Run("Setup2FA", func(t *testing.T) {
		result, err := c.Auth().Setup2FA(context.Background(), &pb.Setup2FARequest{Password: "pass"})
		if err != nil {
			t.Errorf("Setup2FA error: %v", err)
		}
		if result.Secret != "test-secret" {
			t.Errorf("expected secret 'test-secret', got '%s'", result.Secret)
		}
	})

	t.Run("Enable2FA", func(t *testing.T) {
		result, err := c.Auth().Enable2FA(context.Background(), &pb.TwoFactorAuthCodeRequest{TotpCode: "123456"})
		if err != nil {
			t.Errorf("Enable2FA error: %v", err)
		}
		if len(result.RecoveryCodes) != 2 {
			t.Errorf("expected 2 recovery codes, got %d", len(result.RecoveryCodes))
		}
	})

	t.Run("Disable2FA", func(t *testing.T) {
		result, err := c.Auth().Disable2FA(context.Background(), &pb.TwoFactorAuthCodeRequest{TotpCode: "123456"})
		if err != nil {
			t.Errorf("Disable2FA error: %v", err)
		}
		if result.Message != "2FA disabled" {
			t.Errorf("expected message '2FA disabled', got '%s'", result.Message)
		}
	})

	t.Run("GetRecoveryCodes", func(t *testing.T) {
		result, err := c.Auth().GetRecoveryCodes(context.Background(), &pb.TwoFactorAuthCodeRequest{TotpCode: "123456"})
		if err != nil {
			t.Errorf("GetRecoveryCodes error: %v", err)
		}
		if len(result.RecoveryCodes) != 2 {
			t.Errorf("expected 2 codes, got %d", len(result.RecoveryCodes))
		}
		if result.Total != 2 {
			t.Errorf("expected total 2, got %d", result.Total)
		}
	})

	t.Run("RegenerateRecoveryCodes", func(t *testing.T) {
		result, err := c.Auth().RegenerateRecoveryCodes(context.Background(), &pb.TwoFactorAuthCodeRequest{TotpCode: "123456"})
		if err != nil {
			t.Errorf("RegenerateRecoveryCodes error: %v", err)
		}
		if len(result.RecoveryCodes) != 2 {
			t.Errorf("expected 2 codes, got %d", len(result.RecoveryCodes))
		}
	})

	t.Run("ChangePassword", func(t *testing.T) {
		result, err := c.Auth().ChangePassword(context.Background(), &pb.ChangePasswordRequest{OldPassword: "old", NewPassword: "new"})
		if err != nil {
			t.Errorf("ChangePassword error: %v", err)
		}
		if !result.Success {
			t.Error("expected success")
		}
	})

	t.Run("SendChangeEmailCode", func(t *testing.T) {
		err := c.Auth().SendChangeEmailCode(context.Background(), &pb.SendChangeEmailCodeRequest{NewEmail: "new@example.com"})
		if err != nil {
			t.Errorf("SendChangeEmailCode error: %v", err)
		}
	})

	t.Run("ChangeEmail", func(t *testing.T) {
		err := c.Auth().ChangeEmail(context.Background(), &pb.ChangeEmailRequest{NewEmail: "new@example.com"})
		if err != nil {
			t.Errorf("ChangeEmail error: %v", err)
		}
	})

	t.Run("ChangeEmailAndPassword", func(t *testing.T) {
		err := c.Auth().ChangeEmailAndPassword(context.Background(), &pb.ChangeEmailAndPasswordRequest{})
		if err != nil {
			t.Errorf("ChangeEmailAndPassword error: %v", err)
		}
	})
}

func TestFileAPIWithMock(t *testing.T) {
	c, cleanup := setupMockClient(t)
	defer cleanup()
	c.SetToken("test-token")

	t.Run("FindFileByPath", func(t *testing.T) {
		file, err := c.File().FindFileByPath(context.Background(), "/test/path")
		if err != nil {
			t.Errorf("FindFileByPath error: %v", err)
		}
		if file.FullPathName != "/test/path" {
			t.Errorf("expected path '/test/path', got '%s'", file.FullPathName)
		}
	})

	t.Run("CreateFolder", func(t *testing.T) {
		result, err := c.File().CreateFolder(context.Background(), &pb.CreateFolderRequest{ParentPath: "/new", FolderName: "folder"})
		if err != nil {
			t.Errorf("CreateFolder error: %v", err)
		}
		if !result.Result.Success {
			t.Error("expected success")
		}
	})

	t.Run("CreateEncryptedFolder", func(t *testing.T) {
		result, err := c.File().CreateEncryptedFolder(context.Background(), &pb.CreateEncryptedFolderRequest{ParentPath: "/encrypted", FolderName: "folder"})
		if err != nil {
			t.Errorf("CreateEncryptedFolder error: %v", err)
		}
		if !result.Result.Success {
			t.Error("expected success")
		}
	})

	t.Run("UnlockEncryptedFile", func(t *testing.T) {
		result, err := c.File().UnlockEncryptedFile(context.Background(), &pb.UnlockEncryptedFileRequest{Path: "/encrypted/file"})
		if err != nil {
			t.Errorf("UnlockEncryptedFile error: %v", err)
		}
		if !result.Success {
			t.Error("expected success")
		}
	})

	t.Run("LockEncryptedFile", func(t *testing.T) {
		result, err := c.File().LockEncryptedFile(context.Background(), "/encrypted/file")
		if err != nil {
			t.Errorf("LockEncryptedFile error: %v", err)
		}
		if !result.Success {
			t.Error("expected success")
		}
	})

	t.Run("RenameFile", func(t *testing.T) {
		result, err := c.File().RenameFile(context.Background(), &pb.RenameFileRequest{NewName: "newname"})
		if err != nil {
			t.Errorf("RenameFile error: %v", err)
		}
		if !result.Success {
			t.Error("expected success")
		}
	})

	t.Run("RenameFiles", func(t *testing.T) {
		result, err := c.File().RenameFiles(context.Background(), &pb.RenameFilesRequest{})
		if err != nil {
			t.Errorf("RenameFiles error: %v", err)
		}
		if !result.Success {
			t.Error("expected success")
		}
	})

	t.Run("MoveFile", func(t *testing.T) {
		result, err := c.File().MoveFile(context.Background(), &pb.MoveFileRequest{DestPath: "/dst"})
		if err != nil {
			t.Errorf("MoveFile error: %v", err)
		}
		if !result.Success {
			t.Error("expected success")
		}
	})

	t.Run("CopyFile", func(t *testing.T) {
		result, err := c.File().CopyFile(context.Background(), &pb.CopyFileRequest{DestPath: "/dst"})
		if err != nil {
			t.Errorf("CopyFile error: %v", err)
		}
		if !result.Success {
			t.Error("expected success")
		}
	})

	t.Run("DeleteFile", func(t *testing.T) {
		result, err := c.File().DeleteFile(context.Background(), "/test/file")
		if err != nil {
			t.Errorf("DeleteFile error: %v", err)
		}
		if !result.Success {
			t.Error("expected success")
		}
	})

	t.Run("DeleteFilePermanently", func(t *testing.T) {
		result, err := c.File().DeleteFilePermanently(context.Background(), "/test/file")
		if err != nil {
			t.Errorf("DeleteFilePermanently error: %v", err)
		}
		if !result.Success {
			t.Error("expected success")
		}
	})

	t.Run("DeleteFiles", func(t *testing.T) {
		result, err := c.File().DeleteFiles(context.Background(), []string{"/file1", "/file2"})
		if err != nil {
			t.Errorf("DeleteFiles error: %v", err)
		}
		if !result.Success {
			t.Error("expected success")
		}
	})

	t.Run("DeleteFilesPermanently", func(t *testing.T) {
		result, err := c.File().DeleteFilesPermanently(context.Background(), []string{"/file1", "/file2"})
		if err != nil {
			t.Errorf("DeleteFilesPermanently error: %v", err)
		}
		if !result.Success {
			t.Error("expected success")
		}
	})

	t.Run("AddSharedLink", func(t *testing.T) {
		err := c.File().AddSharedLink(context.Background(), &pb.AddSharedLinkRequest{SharedLinkUrl: "http://share"})
		if err != nil {
			t.Errorf("AddSharedLink error: %v", err)
		}
	})

	t.Run("GetFileDetailProperties", func(t *testing.T) {
		_, err := c.File().GetFileDetailProperties(context.Background(), "/test/file")
		if err != nil {
			t.Errorf("GetFileDetailProperties error: %v", err)
		}
	})

	t.Run("GetSpaceInfo", func(t *testing.T) {
		_, err := c.File().GetSpaceInfo(context.Background(), "/test")
		if err != nil {
			t.Errorf("GetSpaceInfo error: %v", err)
		}
	})

	t.Run("GetCloudMemberships", func(t *testing.T) {
		_, err := c.File().GetCloudMemberships(context.Background(), "/test")
		if err != nil {
			t.Errorf("GetCloudMemberships error: %v", err)
		}
	})

	t.Run("GetMetaData", func(t *testing.T) {
		_, err := c.File().GetMetaData(context.Background(), "/test")
		if err != nil {
			t.Errorf("GetMetaData error: %v", err)
		}
	})

	t.Run("GetOriginalPath", func(t *testing.T) {
		result, err := c.File().GetOriginalPath(context.Background(), "/test")
		if err != nil {
			t.Errorf("GetOriginalPath error: %v", err)
		}
		if result.Result != "/original/path" {
			t.Errorf("expected '/original/path', got '%s'", result.Result)
		}
	})

	t.Run("CreateFile", func(t *testing.T) {
		result, err := c.File().CreateFile(context.Background(), &pb.CreateFileRequest{ParentPath: "/parent", FileName: "file.txt"})
		if err != nil {
			t.Errorf("CreateFile error: %v", err)
		}
		if result.FileHandle != 12345 {
			t.Errorf("expected FileHandle 12345, got %d", result.FileHandle)
		}
	})

	t.Run("CloseFile", func(t *testing.T) {
		result, err := c.File().CloseFile(context.Background(), &pb.CloseFileRequest{FileHandle: 12345})
		if err != nil {
			t.Errorf("CloseFile error: %v", err)
		}
		if !result.Success {
			t.Error("expected success")
		}
	})

	t.Run("WriteToFile", func(t *testing.T) {
		result, err := c.File().WriteToFile(context.Background(), &pb.WriteFileRequest{FileHandle: 12345, Length: 1024})
		if err != nil {
			t.Errorf("WriteToFile error: %v", err)
		}
		if result.BytesWritten != 1024 {
			t.Errorf("expected BytesWritten 1024, got %d", result.BytesWritten)
		}
	})

	t.Run("GetDownloadUrl", func(t *testing.T) {
		info, err := c.File().GetDownloadUrl(context.Background(), &pb.GetDownloadUrlPathRequest{Path: "/download"})
		if err != nil {
			t.Errorf("GetDownloadUrl error: %v", err)
		}
		if info.DownloadUrlPath != "http://example.com/download" {
			t.Errorf("expected 'http://example.com/download', got '%s'", info.DownloadUrlPath)
		}
	})

	t.Run("GetSubFiles", func(t *testing.T) {
		files, err := c.File().GetSubFiles(context.Background(), &pb.ListSubFileRequest{Path: "/test"})
		if err != nil {
			t.Errorf("GetSubFiles error: %v", err)
		}
		if len(files) != 1 {
			t.Errorf("expected 1 file, got %d", len(files))
		}
		if files[0].FullPathName != "/test" {
			t.Errorf("expected path '/test', got '%s'", files[0].FullPathName)
		}
	})

	t.Run("GetSearchResults", func(t *testing.T) {
		files, err := c.File().GetSearchResults(context.Background(), &pb.SearchRequest{SearchFor: "test"})
		if err != nil {
			t.Errorf("GetSearchResults error: %v", err)
		}
		if len(files) != 1 {
			t.Errorf("expected 1 file, got %d", len(files))
		}
		if files[0].FullPathName != "/search" {
			t.Errorf("expected path '/search', got '%s'", files[0].FullPathName)
		}
	})
}

func TestMountAPIWithMock(t *testing.T) {
	c, cleanup := setupMockClient(t)
	defer cleanup()
	c.SetToken("test-token")

	t.Run("CanAddMoreMountPoints", func(t *testing.T) {
		result, err := c.Mount().CanAddMoreMountPoints(context.Background())
		if err != nil {
			t.Errorf("CanAddMoreMountPoints error: %v", err)
		}
		if !result.Success {
			t.Error("expected success")
		}
	})

	t.Run("GetMountPoints", func(t *testing.T) {
		_, err := c.Mount().GetMountPoints(context.Background())
		if err != nil {
			t.Errorf("GetMountPoints error: %v", err)
		}
	})

	t.Run("AddMountPoint", func(t *testing.T) {
		_, err := c.Mount().AddMountPoint(context.Background(), &pb.MountOption{MountPoint: "/mnt"})
		if err != nil {
			t.Errorf("AddMountPoint error: %v", err)
		}
	})

	t.Run("RemoveMountPoint", func(t *testing.T) {
		_, err := c.Mount().RemoveMountPoint(context.Background(), &pb.MountPointRequest{MountPoint: "/mnt"})
		if err != nil {
			t.Errorf("RemoveMountPoint error: %v", err)
		}
	})

	t.Run("Mount", func(t *testing.T) {
		_, err := c.Mount().Mount(context.Background(), &pb.MountPointRequest{MountPoint: "/mnt"})
		if err != nil {
			t.Errorf("Mount error: %v", err)
		}
	})

	t.Run("Unmount", func(t *testing.T) {
		_, err := c.Mount().Unmount(context.Background(), &pb.MountPointRequest{MountPoint: "/mnt"})
		if err != nil {
			t.Errorf("Unmount error: %v", err)
		}
	})

	t.Run("UpdateMountPoint", func(t *testing.T) {
		_, err := c.Mount().UpdateMountPoint(context.Background(), &pb.UpdateMountPointRequest{MountPoint: "/mnt"})
		if err != nil {
			t.Errorf("UpdateMountPoint error: %v", err)
		}
	})

	t.Run("GetAvailableDriveLetters", func(t *testing.T) {
		result, err := c.Mount().GetAvailableDriveLetters(context.Background())
		if err != nil {
			t.Errorf("GetAvailableDriveLetters error: %v", err)
		}
		if len(result.DriveLetters) != 2 {
			t.Errorf("expected 2 drive letters, got %d", len(result.DriveLetters))
		}
	})

	t.Run("HasDriveLetters", func(t *testing.T) {
		_, err := c.Mount().HasDriveLetters(context.Background())
		if err != nil {
			t.Errorf("HasDriveLetters error: %v", err)
		}
	})

	t.Run("CanMountBothLocalAndCloud", func(t *testing.T) {
		result, err := c.Mount().CanMountBothLocalAndCloud(context.Background())
		if err != nil {
			t.Errorf("CanMountBothLocalAndCloud error: %v", err)
		}
		if !result.Result {
			t.Error("expected result to be true")
		}
	})
}

func TestTransferAPIWithMock(t *testing.T) {
	c, cleanup := setupMockClient(t)
	defer cleanup()
	c.SetToken("test-token")

	t.Run("GetAllTasksCount", func(t *testing.T) {
		_, err := c.Transfer().GetAllTasksCount(context.Background())
		if err != nil {
			t.Errorf("GetAllTasksCount error: %v", err)
		}
	})

	t.Run("GetDownloadFileCount", func(t *testing.T) {
		_, err := c.Transfer().GetDownloadFileCount(context.Background())
		if err != nil {
			t.Errorf("GetDownloadFileCount error: %v", err)
		}
	})

	t.Run("GetDownloadFileList", func(t *testing.T) {
		_, err := c.Transfer().GetDownloadFileList(context.Background())
		if err != nil {
			t.Errorf("GetDownloadFileList error: %v", err)
		}
	})

	t.Run("GetUploadFileCount", func(t *testing.T) {
		_, err := c.Transfer().GetUploadFileCount(context.Background())
		if err != nil {
			t.Errorf("GetUploadFileCount error: %v", err)
		}
	})

	t.Run("GetUploadFileList", func(t *testing.T) {
		_, err := c.Transfer().GetUploadFileList(context.Background(), &pb.GetUploadFileListRequest{})
		if err != nil {
			t.Errorf("GetUploadFileList error: %v", err)
		}
	})

	t.Run("CancelAllUploadFiles", func(t *testing.T) {
		err := c.Transfer().CancelAllUploadFiles(context.Background())
		if err != nil {
			t.Errorf("CancelAllUploadFiles error: %v", err)
		}
	})

	t.Run("CancelUploadFiles", func(t *testing.T) {
		err := c.Transfer().CancelUploadFiles(context.Background(), []string{"key1", "key2"})
		if err != nil {
			t.Errorf("CancelUploadFiles error: %v", err)
		}
	})

	t.Run("PauseAllUploadFiles", func(t *testing.T) {
		err := c.Transfer().PauseAllUploadFiles(context.Background())
		if err != nil {
			t.Errorf("PauseAllUploadFiles error: %v", err)
		}
	})

	t.Run("PauseUploadFiles", func(t *testing.T) {
		err := c.Transfer().PauseUploadFiles(context.Background(), []string{"key1", "key2"})
		if err != nil {
			t.Errorf("PauseUploadFiles error: %v", err)
		}
	})

	t.Run("ResumeAllUploadFiles", func(t *testing.T) {
		err := c.Transfer().ResumeAllUploadFiles(context.Background())
		if err != nil {
			t.Errorf("ResumeAllUploadFiles error: %v", err)
		}
	})

	t.Run("ResumeUploadFiles", func(t *testing.T) {
		err := c.Transfer().ResumeUploadFiles(context.Background(), []string{"key1", "key2"})
		if err != nil {
			t.Errorf("ResumeUploadFiles error: %v", err)
		}
	})
}

func TestCloudAPIWithMock(t *testing.T) {
	c, cleanup := setupMockClient(t)
	defer cleanup()
	c.SetToken("test-token")

	t.Run("CanAddMoreCloudApis", func(t *testing.T) {
		result, err := c.CloudAPI().CanAddMoreCloudApis(context.Background())
		if err != nil {
			t.Errorf("CanAddMoreCloudApis error: %v", err)
		}
		if !result.Success {
			t.Error("expected success")
		}
	})

	t.Run("GetAllCloudApis", func(t *testing.T) {
		_, err := c.CloudAPI().GetAllCloudApis(context.Background())
		if err != nil {
			t.Errorf("GetAllCloudApis error: %v", err)
		}
	})

	t.Run("GetCloudAPIConfig", func(t *testing.T) {
		_, err := c.CloudAPI().GetCloudAPIConfig(context.Background(), &pb.GetCloudAPIConfigRequest{CloudName: "test"})
		if err != nil {
			t.Errorf("GetCloudAPIConfig error: %v", err)
		}
	})

	t.Run("SetCloudAPIConfig", func(t *testing.T) {
		err := c.CloudAPI().SetCloudAPIConfig(context.Background(), &pb.SetCloudAPIConfigRequest{})
		if err != nil {
			t.Errorf("SetCloudAPIConfig error: %v", err)
		}
	})

	t.Run("RemoveCloudAPI", func(t *testing.T) {
		result, err := c.CloudAPI().RemoveCloudAPI(context.Background(), &pb.RemoveCloudAPIRequest{CloudName: "test"})
		if err != nil {
			t.Errorf("RemoveCloudAPI error: %v", err)
		}
		if !result.Success {
			t.Error("expected success")
		}
	})
}

func TestWithTimeoutAndAuth(t *testing.T) {
	c, cleanup := setupMockClient(t)
	defer cleanup()

	t.Run("withTimeout", func(t *testing.T) {
		ctx, cancel := c.withTimeout(context.Background(), 5*time.Second)
		defer cancel()
		deadline, ok := ctx.Deadline()
		if !ok {
			t.Error("expected deadline to be set")
		}
		if deadline.IsZero() {
			t.Error("expected non-zero deadline")
		}
	})

	t.Run("withTimeout default", func(t *testing.T) {
		ctx, cancel := c.withTimeout(context.Background(), 0)
		defer cancel()
		deadline, ok := ctx.Deadline()
		if !ok {
			t.Error("expected deadline to be set")
		}
		if deadline.IsZero() {
			t.Error("expected non-zero deadline")
		}
	})

	t.Run("withAuth without token", func(t *testing.T) {
		ctx := c.withAuth(context.Background())
		if ctx != context.Background() {
			t.Error("without token should return same context")
		}
	})

	t.Run("withAuth with token", func(t *testing.T) {
		c.SetToken("test-token")
		ctx := c.withAuth(context.Background())
		if ctx == context.Background() {
			t.Error("with token should return different context")
		}
	})
}

func TestClientErrorHandling(t *testing.T) {
	c, cleanup := setupMockClient(t)
	defer cleanup()
	c.SetToken("test-token")

	t.Run("GetToken auto-sets token on success", func(t *testing.T) {
		c.SetToken("")
		token, err := c.Public().GetToken(context.Background(), &pb.GetTokenRequest{UserName: "test"})
		if err != nil {
			t.Errorf("GetToken error: %v", err)
		}
		if c.GetToken() != token.Token {
			t.Error("client token should be set after GetToken success")
		}
	})

	t.Run("LoginWithThirdPartyAccount auto-sets token", func(t *testing.T) {
		c.SetToken("")
		token, err := c.Public().LoginWithThirdPartyAccount(context.Background(), &pb.LoginWithThirdPartyAccountRequest{})
		if err != nil {
			t.Errorf("LoginWithThirdPartyAccount error: %v", err)
		}
		if c.GetToken() != token.Token {
			t.Error("client token should be set after LoginWithThirdPartyAccount success")
		}
	})

	t.Run("LoginWith2FA auto-sets token", func(t *testing.T) {
		c.SetToken("")
		token, err := c.Public().LoginWith2FA(context.Background(), &pb.LoginWith2FARequest{})
		if err != nil {
			t.Errorf("LoginWith2FA error: %v", err)
		}
		if c.GetToken() != token.Token {
			t.Error("client token should be set after LoginWith2FA success")
		}
	})

	t.Run("Token lifecycle", func(t *testing.T) {
		c.SetToken("first-token")
		if c.GetToken() != "first-token" {
			t.Error("GetToken should return first-token")
		}
		c.SetToken("second-token")
		if c.GetToken() != "second-token" {
			t.Error("GetToken should return second-token")
		}
	})
}

func TestClientClose(t *testing.T) {
	c, cleanup := setupMockClient(t)
	defer cleanup()

	err := c.Close()
	if err != nil {
		t.Errorf("Close error: %v", err)
	}
}

func TestNewClientConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{"empty address", Config{Address: ""}, true},
		{"valid address", Config{Address: "localhost:19798"}, false},
		{"with timeout", Config{Address: "localhost:19798", Timeout: 10 * time.Second}, false},
		{"with token", Config{Address: "localhost:19798", Token: "test-token"}, false},
		{"zero timeout uses default", Config{Address: "localhost:19798", Timeout: 0}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewClient(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if c != nil {
				c.Close()
			}
		})
	}
}

type mockUploadErrorServer struct {
	pb.UnimplementedCloudDriveFileSrvServer
	closeFileCalled bool
	closeFileMutex  sync.Mutex
	sendError       bool
	sendErrorMutex  sync.Mutex
}

func (s *mockUploadErrorServer) CreateFile(ctx context.Context, req *pb.CreateFileRequest) (*pb.CreateFileResult, error) {
	return &pb.CreateFileResult{FileHandle: 12345}, nil
}

func (s *mockUploadErrorServer) CloseFile(ctx context.Context, req *pb.CloseFileRequest) (*pb.FileOperationResult, error) {
	s.closeFileMutex.Lock()
	s.closeFileCalled = true
	s.closeFileMutex.Unlock()
	return &pb.FileOperationResult{Success: true}, nil
}

func (s *mockUploadErrorServer) WriteToFileStream(stream pb.CloudDriveFileSrv_WriteToFileStreamServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.WriteFileResult{BytesWritten: 0})
		}
		if err != nil {
			return err
		}
		s.sendErrorMutex.Lock()
		if s.sendError {
			s.sendErrorMutex.Unlock()
			return fmt.Errorf("simulated send error")
		}
		s.sendErrorMutex.Unlock()
		if req.CloseFile {
			return stream.SendAndClose(&pb.WriteFileResult{BytesWritten: req.Length})
		}
	}
}

func setupMockUploadErrorClient(t *testing.T) (*Client, *mockUploadErrorServer, func()) {
	lis := bufconn.Listen(bufSize)
	s := &mockUploadErrorServer{}
	grpcServer := grpc.NewServer()
	pb.RegisterCloudDriveFileSrvServer(grpcServer, s)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			t.Logf("Server error: %v", err)
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(dialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	client := &Client{
		conn:   conn,
		client: pb.NewCloudDriveFileSrvClient(conn),
	}

	return client, s, func() {
		client.Close()
		grpcServer.Stop()
		lis.Close()
	}
}

func TestUploadLocalFileClosesHandleOnError(t *testing.T) {
	c, mockServer, cleanup := setupMockUploadErrorClient(t)
	defer cleanup()
	c.SetToken("test-token")

	tmpFile, err := os.CreateTemp("", "upload-test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString("test content for upload")
	tmpFile.Close()

	mockServer.sendErrorMutex.Lock()
	mockServer.sendError = true
	mockServer.sendErrorMutex.Unlock()

	_, err = c.File().UploadLocalFile(context.Background(), tmpFile.Name(), "/test/upload.txt")
	if err == nil {
		t.Error("expected upload to fail with simulated error")
	}

	mockServer.closeFileMutex.Lock()
	closeCalled := mockServer.closeFileCalled
	mockServer.closeFileMutex.Unlock()

	if !closeCalled {
		t.Error("expected CloseFile to be called on upload failure")
	}
}
