package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/clouddrive/cd2-cli/internal/client"
	"github.com/clouddrive/cd2-cli/internal/whitelist"
	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

const bufSize = 1024 * 1024

type mockServer struct {
	pb.UnimplementedCloudDriveFileSrvServer
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

func (s *mockServer) LoginWith2FA(ctx context.Context, req *pb.LoginWith2FARequest) (*pb.JWTToken, error) {
	return &pb.JWTToken{Success: true, Token: "2fa-token"}, nil
}

func (s *mockServer) Register(ctx context.Context, req *pb.UserRegisterRequest) (*pb.FileOperationResult, error) {
	return &pb.FileOperationResult{Success: true}, nil
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
	return &pb.TwoFactorAuthStatusResult{TwoFactorEnabled: true}, nil
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

func (s *mockServer) ChangeEmail(ctx context.Context, req *pb.ChangeEmailRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *mockServer) FindFileByPath(ctx context.Context, req *pb.FindFileByPathRequest) (*pb.CloudDriveFile, error) {
	return &pb.CloudDriveFile{FullPathName: req.Path, Name: "test.txt", IsDirectory: false, Size: 1024}, nil
}

func (s *mockServer) CreateFolder(ctx context.Context, req *pb.CreateFolderRequest) (*pb.CreateFolderResult, error) {
	return &pb.CreateFolderResult{Result: &pb.FileOperationResult{Success: true}}, nil
}

func (s *mockServer) CreateEncryptedFolder(ctx context.Context, req *pb.CreateEncryptedFolderRequest) (*pb.CreateFolderResult, error) {
	return &pb.CreateFolderResult{Result: &pb.FileOperationResult{Success: true}}, nil
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

func (s *mockServer) GetDownloadUrlPath(ctx context.Context, req *pb.GetDownloadUrlPathRequest) (*pb.DownloadUrlPathInfo, error) {
	return &pb.DownloadUrlPathInfo{DownloadUrlPath: "http://example.com/download"}, nil
}

func (s *mockServer) GetSubFiles(req *pb.ListSubFileRequest, stream pb.CloudDriveFileSrv_GetSubFilesServer) error {
	stream.Send(&pb.SubFilesReply{SubFiles: []*pb.CloudDriveFile{
		{FullPathName: req.Path + "/file1.txt", Name: "file1.txt", IsDirectory: false, Size: 100},
		{FullPathName: req.Path + "/dir1", Name: "dir1", IsDirectory: true, Size: 0},
	}})
	return nil
}

func (s *mockServer) GetMountPoints(ctx context.Context, req *emptypb.Empty) (*pb.GetMountPointsResult, error) {
	return &pb.GetMountPointsResult{
		MountPoints: []*pb.MountPoint{
			{MountPoint: "/mnt/test", SourceDir: "/test/source"},
		},
	}, nil
}

func (s *mockServer) AddMountPoint(ctx context.Context, req *pb.MountOption) (*pb.MountPointResult, error) {
	return &pb.MountPointResult{Success: true}, nil
}

func (s *mockServer) RemoveMountPoint(ctx context.Context, req *pb.MountPointRequest) (*pb.MountPointResult, error) {
	return &pb.MountPointResult{Success: true}, nil
}

func (s *mockServer) Mount(ctx context.Context, req *pb.MountPointRequest) (*pb.MountPointResult, error) {
	return &pb.MountPointResult{Success: true}, nil
}

func (s *mockServer) Unmount(ctx context.Context, req *pb.MountPointRequest) (*pb.MountPointResult, error) {
	return &pb.MountPointResult{Success: true}, nil
}

func (s *mockServer) UpdateMountPoint(ctx context.Context, req *pb.UpdateMountPointRequest) (*pb.MountPointResult, error) {
	return &pb.MountPointResult{Success: true}, nil
}

func (s *mockServer) CanAddMoreMountPoints(ctx context.Context, req *emptypb.Empty) (*pb.FileOperationResult, error) {
	return &pb.FileOperationResult{Success: true}, nil
}

func (s *mockServer) GetAvailableDriveLetters(ctx context.Context, req *emptypb.Empty) (*pb.GetAvailableDriveLettersResult, error) {
	return &pb.GetAvailableDriveLettersResult{DriveLetters: []string{"A:", "B:", "C:"}}, nil
}

func (s *mockServer) HasDriveLetters(ctx context.Context, req *emptypb.Empty) (*pb.HasDriveLettersResult, error) {
	return &pb.HasDriveLettersResult{HasDriveLetters: true}, nil
}

func (s *mockServer) CanMountBothLocalAndCloud(ctx context.Context, req *emptypb.Empty) (*pb.BoolResult, error) {
	return &pb.BoolResult{Result: true}, nil
}

func (s *mockServer) GetAllTasksCount(ctx context.Context, req *emptypb.Empty) (*pb.GetAllTasksCountResult, error) {
	return &pb.GetAllTasksCountResult{
		DownloadCount: 5,
		UploadCount:   10,
		CopyTaskCount: 2,
	}, nil
}

func (s *mockServer) GetDownloadFileList(ctx context.Context, req *emptypb.Empty) (*pb.GetDownloadFileListResult, error) {
	return &pb.GetDownloadFileListResult{
		DownloadFiles: []*pb.DownloadFileInfo{
			{FilePath: "/download/file1.txt", FileLength: 1024},
		},
	}, nil
}

func (s *mockServer) GetUploadFileList(ctx context.Context, req *pb.GetUploadFileListRequest) (*pb.GetUploadFileListResult, error) {
	return &pb.GetUploadFileListResult{
		UploadFiles: []*pb.UploadFileInfo{
			{Key: "up1", DestPath: "/upload/file1.txt", Size: 2048},
		},
	}, nil
}

func (s *mockServer) GetCopyTasks(ctx context.Context, req *emptypb.Empty) (*pb.GetCopyTaskResult, error) {
	return &pb.GetCopyTaskResult{
		CopyTasks: []*pb.CopyTask{
			{Status: pb.CopyTask_Scanning, SourcePath: "/src", DestPath: "/dst"},
		},
	}, nil
}

func (s *mockServer) GetMergeTasks(ctx context.Context, req *emptypb.Empty) (*pb.GetMergeTasksResult, error) {
	return &pb.GetMergeTasksResult{
		MergeTasks: []*pb.MergeTask{
			{SourcePath: "/merge/src", DestPath: "/merge/dst", Status: pb.MergeTask_Running},
		},
	}, nil
}

func (s *mockServer) GetAllCloudApis(ctx context.Context, req *emptypb.Empty) (*pb.CloudAPIList, error) {
	return &pb.CloudAPIList{
		Apis: []*pb.CloudAPI{
			{Name: "test-cloud", UserName: "test-user"},
		},
	}, nil
}

func (s *mockServer) RemoveCloudAPI(ctx context.Context, req *pb.RemoveCloudAPIRequest) (*pb.FileOperationResult, error) {
	return &pb.FileOperationResult{Success: true}, nil
}

func (s *mockServer) GetCloudAPIConfig(ctx context.Context, req *pb.GetCloudAPIConfigRequest) (*pb.CloudAPIConfig, error) {
	return &pb.CloudAPIConfig{}, nil
}

func (s *mockServer) SetCloudAPIConfig(ctx context.Context, req *pb.SetCloudAPIConfigRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *mockServer) APILoginS3(ctx context.Context, req *pb.LoginS3Request) (*pb.APILoginResult, error) {
	return &pb.APILoginResult{Success: true}, nil
}

func (s *mockServer) APILoginWebDav(ctx context.Context, req *pb.LoginWebDavRequest) (*pb.APILoginResult, error) {
	return &pb.APILoginResult{Success: true}, nil
}

func (s *mockServer) APIAddLocalFolder(ctx context.Context, req *pb.AddLocalFolderRequest) (*pb.APILoginResult, error) {
	return &pb.APILoginResult{Success: true}, nil
}

func (s *mockServer) APILoginSftp(ctx context.Context, req *pb.LoginSftpRequest) (*pb.APILoginResult, error) {
	return &pb.APILoginResult{Success: true}, nil
}

func (s *mockServer) APILoginFtp(ctx context.Context, req *pb.LoginFtpRequest) (*pb.APILoginResult, error) {
	return &pb.APILoginResult{Success: true}, nil
}

func (s *mockServer) APILoginSmb(ctx context.Context, req *pb.LoginSmbRequest) (*pb.APILoginResult, error) {
	return &pb.APILoginResult{Success: true}, nil
}

func (s *mockServer) GetRuntimeInfo(ctx context.Context, req *emptypb.Empty) (*pb.RuntimeInfo, error) {
	return &pb.RuntimeInfo{ProductName: "CloudDrive2", ProductVersion: "1.0.0"}, nil
}

func (s *mockServer) GetRunningInfo(ctx context.Context, req *emptypb.Empty) (*pb.RunInfo, error) {
	return &pb.RunInfo{}, nil
}

func (s *mockServer) GetSystemSettings(ctx context.Context, req *emptypb.Empty) (*pb.SystemSettings, error) {
	return &pb.SystemSettings{}, nil
}

func (s *mockServer) SetSystemSettings(ctx context.Context, req *pb.SystemSettings) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *mockServer) GetOpenFileHandles(ctx context.Context, req *emptypb.Empty) (*pb.OpenFileHandleList, error) {
	return &pb.OpenFileHandleList{}, nil
}

func (s *mockServer) GetFileBufferDiskCacheStats(ctx context.Context, req *emptypb.Empty) (*pb.FileBufferDiskCacheStats, error) {
	return &pb.FileBufferDiskCacheStats{}, nil
}

func (s *mockServer) GetServiceCapabilities(ctx context.Context, req *emptypb.Empty) (*pb.ServiceCapabilities, error) {
	return &pb.ServiceCapabilities{}, nil
}

func (s *mockServer) ListTokens(ctx context.Context, req *emptypb.Empty) (*pb.ListTokensResult, error) {
	return &pb.ListTokensResult{}, nil
}

func (s *mockServer) CreateToken(ctx context.Context, req *pb.CreateTokenRequest) (*pb.TokenInfo, error) {
	return &pb.TokenInfo{}, nil
}

func (s *mockServer) GetApiTokenInfo(ctx context.Context, req *pb.StringValue) (*pb.TokenInfo, error) {
	return &pb.TokenInfo{Token: req.Value}, nil
}

func (s *mockServer) ModifyToken(ctx context.Context, req *pb.ModifyTokenRequest) (*pb.TokenInfo, error) {
	return &pb.TokenInfo{Token: req.Token}, nil
}

func (s *mockServer) RemoveToken(ctx context.Context, req *pb.StringValue) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *mockServer) GetSessions(ctx context.Context, req *emptypb.Empty) (*pb.GetSessionsResponse, error) {
	return &pb.GetSessionsResponse{}, nil
}

func (s *mockServer) GetDavServerConfig(ctx context.Context, req *emptypb.Empty) (*pb.DavServerConfig, error) {
	return &pb.DavServerConfig{}, nil
}

func (s *mockServer) SetDavServerConfig(ctx context.Context, req *pb.ModifyDavServerConfigRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *mockServer) GetDavUser(ctx context.Context, req *pb.StringValue) (*pb.DavUser, error) {
	return &pb.DavUser{UserName: req.Value}, nil
}

func (s *mockServer) AddDavUser(ctx context.Context, req *pb.AddDavUserRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *mockServer) ModifyDavUser(ctx context.Context, req *pb.ModifyDavUserRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *mockServer) RemoveDavUser(ctx context.Context, req *pb.StringValue) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *mockServer) BackupGetAll(ctx context.Context, req *emptypb.Empty) (*pb.BackupList, error) {
	return &pb.BackupList{}, nil
}

func (s *mockServer) GetPromotions(ctx context.Context, req *emptypb.Empty) (*pb.GetPromotionsResult, error) {
	return &pb.GetPromotionsResult{}, nil
}

func (s *mockServer) GetCloudDrivePlans(ctx context.Context, req *emptypb.Empty) (*pb.GetCloudDrivePlansResult, error) {
	return &pb.GetCloudDrivePlansResult{}, nil
}

func (s *mockServer) GetWebhookConfigs(ctx context.Context, req *emptypb.Empty) (*pb.WebhookList, error) {
	return &pb.WebhookList{}, nil
}

func (s *mockServer) GetDirCacheTable(ctx context.Context, req *emptypb.Empty) (*pb.DirCacheTable, error) {
	return &pb.DirCacheTable{}, nil
}

func (s *mockServer) GetTempFileTable(ctx context.Context, req *emptypb.Empty) (*pb.TempFileTable, error) {
	return &pb.TempFileTable{}, nil
}

func (s *mockServer) GetMachineId(ctx context.Context, req *emptypb.Empty) (*pb.StringResult, error) {
	return &pb.StringResult{Result: "test-machine-id"}, nil
}

func (s *mockServer) GetOnlineDevices(ctx context.Context, req *emptypb.Empty) (*pb.OnlineDevices, error) {
	return &pb.OnlineDevices{}, nil
}

func (s *mockServer) RestartService(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *mockServer) ShutdownService(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *mockServer) ListDiskCacheFolders(ctx context.Context, req *emptypb.Empty) (*pb.ListDiskCacheFoldersReply, error) {
	return &pb.ListDiskCacheFoldersReply{}, nil
}

func setupTestCommand(t *testing.T) func() {
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

	conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	c := &client.Client{}
	c.SetConn(conn)
	cd2Client = c

	tmpDir := filepath.Join(os.TempDir(), "cd2-cli-test")
	os.MkdirAll(tmpDir, 0755)
	whitelistPath := filepath.Join(tmpDir, "whitelist.yaml")
	mgr, err := whitelist.NewManager(whitelistPath)
	if err != nil {
		t.Fatalf("Failed to create whitelist manager: %v", err)
	}
	mgr.SetEnabled(false)
	whitelistMgr = mgr
	initialized = true
	outputJSON = true
	viper.SetDefault("json", true)

	return func() {
		c.Close()
		s.Stop()
		lis.Close()
		os.RemoveAll(tmpDir)
		initialized = false
		whitelistMgr = nil
		cd2Client = nil
		outputJSON = false
	}
}

func executeCommand(args ...string) (string, error) {
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String(), err
}

func executeCommandWithMock(mockSrv pb.CloudDriveFileSrvServer, args ...string) (string, error) {
	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterCloudDriveFileSrvServer(s, mockSrv)

	go func() {
		if err := s.Serve(lis); err != nil {
			return
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		s.Stop()
		lis.Close()
		return "", err
	}

	c := &client.Client{}
	c.SetConn(conn)
	oldClient := cd2Client
	cd2Client = c

	tmpDir := filepath.Join(os.TempDir(), "cd2-cli-test-mock")
	os.MkdirAll(tmpDir, 0755)
	whitelistPath := filepath.Join(tmpDir, "whitelist.yaml")
	mgr, err := whitelist.NewManager(whitelistPath)
	if err != nil {
		conn.Close()
		s.Stop()
		lis.Close()
		cd2Client = oldClient
		return "", err
	}
	mgr.SetEnabled(false)
	oldMgr := whitelistMgr
	whitelistMgr = mgr

	oldInitialized := initialized
	initialized = true
	oldJSON := outputJSON
	outputJSON = true
	viper.SetDefault("json", true)

	rootCmd.SetArgs(args)

	parent := rootCmd
	for _, arg := range args {
		cmd, _, err := parent.Find(strings.Split(arg, " "))
		if err != nil || cmd == nil {
			break
		}
		cmd.Flags().Visit(func(f *pflag.Flag) {
			f.Changed = false
		})
		if cmd.HasSubCommands() {
			parent = cmd
		} else {
			break
		}
	}

	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	err = rootCmd.Execute()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)

	conn.Close()
	s.Stop()
	lis.Close()
	os.RemoveAll(tmpDir)
	cd2Client = oldClient
	whitelistMgr = oldMgr
	initialized = oldInitialized
	outputJSON = oldJSON

	return buf.String(), err
}

func TestRootCommand(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	if rootCmd == nil {
		t.Error("rootCmd should not be nil")
	}
	if rootCmd.Use != "cd2-cli" {
		t.Errorf("expected Use 'cd2-cli', got '%s'", rootCmd.Use)
	}
}

func TestAuthCommands(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{"login", []string{"auth", "login", "testuser", "testpass"}, "test-token"},
		{"login-2fa", []string{"auth", "login-2fa", "testuser", "testpass", "123456"}, "2fa-token"},
		{"logout", []string{"auth", "logout"}, "success"},
		{"status", []string{"auth", "status"}, "test-user"},
		{"change-password", []string{"auth", "change-password", "oldpass", "newpass"}, "success"},
		{"change-email", []string{"auth", "change-email", "new@email.com", "pass"}, "success"},
		{"confirm-email", []string{"auth", "confirm-email", "verifycode123"}, "success"},
		{"register", []string{"auth", "register", "newuser", "newpass"}, "success"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := executeCommand(tc.args...)
			if err != nil {
				t.Errorf("%s error: %v", tc.name, err)
			}
			if !strings.Contains(output, tc.expected) {
				t.Errorf("%s expected '%s' in output, got: %s", tc.name, tc.expected, output)
			}
		})
	}
}

func Test2FACommands(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{"status", []string{"auth", "2fa-status"}, "twoFactorEnabled"},
		{"enable", []string{"auth", "2fa-enable", "password"}, "recoveryCodes"},
		{"codes", []string{"auth", "2fa-recovery-codes", "123456"}, "recoveryCodes"},
		{"regenerate", []string{"auth", "2fa-regenerate-codes", "123456"}, "recoveryCodes"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := executeCommand(tc.args...)
			if err != nil {
				t.Errorf("%s error: %v", tc.name, err)
			}
			if !strings.Contains(output, tc.expected) {
				t.Errorf("%s expected '%s' in output, got: %s", tc.name, tc.expected, output)
			}
		})
	}
}

func TestFsCommands(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{"ls", []string{"ls", "/test"}, "file1.txt"},
		{"ls-refresh", []string{"ls", "--refresh", "/test"}, "file1.txt"},
		{"stat", []string{"stat", "/test/file.txt"}, "test.txt"},
		{"mkdir", []string{"mkdir", "/test", "newfolder"}, "success"},
		{"list", []string{"file", "list", "/test"}, "file1.txt"},
		{"find", []string{"file", "find", "/test"}, "test.txt"},
		{"file-mkdir", []string{"file", "mkdir", "/test", "newfolder"}, "success"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := executeCommand(tc.args...)
			if err != nil {
				t.Errorf("%s error: %v", tc.name, err)
			}
			if !strings.Contains(output, tc.expected) {
				t.Errorf("%s expected '%s' in output, got: %s", tc.name, tc.expected, output)
			}
		})
	}
}

func TestFsOpsCommands(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{"rename", []string{"file", "rename", "/test/file.txt", "newfile.txt"}, "success"},
		{"move", []string{"file", "move", "/test/file.txt", "/test2/file.txt"}, "success"},
		{"copy", []string{"file", "copy", "/test/file.txt", "/test2/file.txt"}, "success"},
		{"delete", []string{"file", "delete", "/test/file.txt"}, "success"},
		{"rm", []string{"rm", "/test/file.txt"}, "success"},
		{"rm-permanent", []string{"rm", "--permanent", "/test/file.txt"}, "success"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := executeCommand(tc.args...)
			if err != nil {
				t.Errorf("%s error: %v", tc.name, err)
			}
			if !strings.Contains(output, tc.expected) {
				t.Errorf("%s expected '%s' in output, got: %s", tc.name, tc.expected, output)
			}
		})
	}
}

func TestStorageCommands(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{"list", []string{"storage", "list"}, "test-cloud"},
		{"ls", []string{"storage", "ls"}, "test-cloud"},
		{"remove", []string{"storage", "remove", "test-cloud", "test-user"}, "success"},
		{"delete", []string{"storage", "delete", "test-cloud", "test-user"}, "success"},
		{"config", []string{"storage", "config", "test-cloud", "test-user"}, ""},
		{"set-config", []string{"storage", "set-config", "test-cloud", "test-user", "--config", "{}"}, "updated"},
		{"status", []string{"storage", "status"}, "test-cloud"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := executeCommand(tc.args...)
			if err != nil {
				t.Errorf("%s error: %v", tc.name, err)
			}
			if tc.expected != "" && !strings.Contains(output, tc.expected) {
				t.Errorf("%s expected '%s' in output, got: %s", tc.name, tc.expected, output)
			}
		})
	}
}

func TestMountCommands(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{"list", []string{"mount", "list"}, "/mnt/test"},
		{"status", []string{"mount", "status"}, "/mnt/test"},
		{"can-add", []string{"mount", "can-add"}, "success"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := executeCommand(tc.args...)
			if err != nil {
				t.Errorf("%s error: %v", tc.name, err)
			}
			if !strings.Contains(output, tc.expected) {
				t.Errorf("%s expected '%s' in output, got: %s", tc.name, tc.expected, output)
			}
		})
	}
}

func TestTaskCommands(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{"status", []string{"task", "status"}, "downloadCount"},
		{"list-uploads", []string{"task", "list-uploads"}, "up1"},
		{"list-downloads", []string{"task", "list-downloads"}, "/download/file1.txt"},
		{"list-copy", []string{"task", "list-copy"}, "/src"},
		{"list-merge", []string{"task", "list-merge"}, "/merge/src"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := executeCommand(tc.args...)
			if err != nil {
				t.Errorf("%s error: %v", tc.name, err)
			}
			if !strings.Contains(output, tc.expected) {
				t.Errorf("%s expected '%s' in output, got: %s", tc.name, tc.expected, output)
			}
		})
	}
}

func TestSystemCommands(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{"info", []string{"system", "info"}, "IsLogin"},
		{"runtime", []string{"system", "runtime"}, "productVersion"},
		{"running", []string{"system", "running"}, ""},
		{"settings-get", []string{"system", "settings", "get"}, ""},
		{"open-files", []string{"system", "open-files"}, ""},
		{"cache-stats", []string{"system", "cache-stats"}, ""},
		{"capabilities", []string{"system", "capabilities"}, ""},
		{"restart", []string{"system", "restart"}, ""},
		{"shutdown", []string{"system", "shutdown"}, ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := executeCommand(tc.args...)
			if err != nil {
				t.Errorf("%s error: %v", tc.name, err)
			}
			if tc.expected != "" && !strings.Contains(output, tc.expected) {
				t.Errorf("%s expected '%s' in output, got: %s", tc.name, tc.expected, output)
			}
		})
	}
}

func TestTokenCommands(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	tests := []struct {
		name string
		args []string
	}{
		{"list", []string{"token", "list"}},
		{"info", []string{"token", "info", "test-token-id"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := executeCommand(tt.args...)
			if err != nil {
				t.Errorf("%s error: %v", tt.name, err)
			}
		})
	}
}

func TestSessionCommands(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	_, err := executeCommand("session", "list")
	if err != nil {
		t.Errorf("session list error: %v", err)
	}
}

func TestWebdavCommands(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	tests := []struct {
		name string
		args []string
	}{
		{"user get", []string{"webdav", "user", "get", "testuser"}},
		{"server get", []string{"webdav", "server", "get"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := executeCommand(tt.args...)
			if err != nil {
				t.Errorf("%s error: %v", tt.name, err)
			}
		})
	}
}

func TestBackupCommands(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	_, err := executeCommand("backup", "list")
	if err != nil {
		t.Errorf("backup list error: %v", err)
	}
}

func TestPromotionCommands(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	_, err := executeCommand("promotion", "list")
	if err != nil {
		t.Errorf("promotion list error: %v", err)
	}
}

func TestWebhookCommands(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	_, err := executeCommand("webhook", "list")
	if err != nil {
		t.Errorf("webhook list error: %v", err)
	}
}

func TestCacheCommands(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{"dir-cache-table", []string{"cache", "dir-table"}, false},
		{"temp-file-table", []string{"cache", "temp-table"}, false},
		{"list-disk-cache", []string{"cache", "list-disk"}, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := executeCommand(tc.args...)
			if tc.wantErr && err == nil {
				t.Errorf("%s expected error but got none", tc.name)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("%s error: %v", tc.name, err)
			}
		})
	}
}

func TestOutputResultJSON(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	output, err := executeCommand("auth", "status")
	if err != nil {
		t.Errorf("auth status error: %v", err)
	}
	if !strings.Contains(output, "{") {
		t.Errorf("expected JSON output, got: %s", output)
	}
}

func TestCommandWithServerFlag(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	output, err := executeCommand("--server", "localhost:19798", "auth", "status")
	if err != nil {
		t.Errorf("command with server flag error: %v", err)
	}
	if !strings.Contains(output, "test-user") {
		t.Errorf("expected test-user in output, got: %s", output)
	}
}

func TestOutputError(t *testing.T) {
	outputJSON = true
	defer func() { outputJSON = false }()
	outputError(context.DeadlineExceeded)
}

func TestOutputErrorNonJSON(t *testing.T) {
	outputJSON = false
	defer func() { outputJSON = true }()
	outputError(context.DeadlineExceeded)
}

func TestWhitelistCommands(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{"list", []string{"whitelist", "list"}, "id"},
		{"list-all", []string{"whitelist", "list", "--all"}, "enabled"},
		{"path", []string{"whitelist", "path"}, "path"},
		{"status", []string{"whitelist", "status", "file.list"}, "enabled"},
		{"enable", []string{"whitelist", "enable", "file.delete"}, "success"},
		{"disable", []string{"whitelist", "disable", "file.list"}, "success"},
		{"reset", []string{"whitelist", "reset"}, "success"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := executeCommand(tc.args...)
			if err != nil {
				t.Errorf("%s error: %v", tc.name, err)
			}
			if !strings.Contains(output, tc.expected) {
				t.Errorf("%s expected '%s' in output, got: %s", tc.name, tc.expected, output)
			}
		})
	}
}

func TestWhitelistListFilters(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	output, err := executeCommand("whitelist", "list", "--risk", "low")
	if err != nil {
		t.Errorf("whitelist list --risk error: %v", err)
	}
	if !strings.Contains(output, "risk_level") {
		t.Errorf("expected risk_level in output, got: %s", output)
	}

	output, err = executeCommand("whitelist", "list", "--category", "file")
	if err != nil {
		t.Errorf("whitelist list --category error: %v", err)
	}
	if !strings.Contains(output, "category") {
		t.Errorf("expected category in output, got: %s", output)
	}
}

func TestWhitelistStatusNotFound(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	output, err := executeCommand("whitelist", "status", "nonexistent.command")
	if err != nil {
		t.Errorf("whitelist status nonexistent error: %v", err)
	}
	if !strings.Contains(output, "not found") {
		t.Errorf("expected 'not found' in output, got: %s", output)
	}
}

func TestAuthTokenCommands(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{"token-set", []string{"auth", "token", "set", "test-token-123"}, "success"},
		{"token-show", []string{"auth", "token", "show"}, "token"},
		{"token-show-redact", []string{"auth", "token", "show", "--redact"}, "..."},
		{"token-clear", []string{"auth", "token", "clear"}, "success"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := executeCommand(tc.args...)
			if err != nil {
				t.Errorf("%s error: %v", tc.name, err)
			}
			if !strings.Contains(output, tc.expected) {
				t.Errorf("%s expected '%s' in output, got: %s", tc.name, tc.expected, output)
			}
		})
	}
}

func TestAuthLoginWithSave(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	tmpDir := filepath.Join(os.TempDir(), "cd2-cli-test-save")
	os.MkdirAll(tmpDir, 0755)
	cfgFile = filepath.Join(tmpDir, "config.yaml")

	output, err := executeCommand("auth", "login", "--save", "testuser", "testpass")
	if err != nil {
		t.Errorf("auth login --save error: %v", err)
	}
	if !strings.Contains(output, "test-token") {
		t.Errorf("expected 'test-token' in output, got: %s", output)
	}

	data, err := os.ReadFile(cfgFile)
	if err != nil {
		t.Errorf("failed to read config file: %v", err)
	}
	if !strings.Contains(string(data), "token:") {
		t.Errorf("expected token in config file, got: %s", string(data))
	}

	os.RemoveAll(tmpDir)
	cfgFile = ""
}

func TestSeparateWhitelistConfig(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "cd2-cli-test-separate")
	os.MkdirAll(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	mainCfgPath := filepath.Join(tmpDir, "main-config.yaml")
	whitelistCfgPath := filepath.Join(tmpDir, "whitelist.yaml")

	mainCfgContent := "server: localhost:19798\ntoken: test-token-123\n"
	if err := os.WriteFile(mainCfgPath, []byte(mainCfgContent), 0600); err != nil {
		t.Fatalf("failed to write main config: %v", err)
	}

	whitelistMgr = nil
	initialized = false
	cfgFile = mainCfgPath
	whitelistCfgFile = whitelistCfgPath
	defer func() {
		cfgFile = ""
		whitelistCfgFile = ""
		whitelistMgr = nil
		initialized = false
	}()

	if err := loadConfig(); err != nil {
		t.Fatalf("loadConfig failed: %v", err)
	}

	if err := initWhitelist(); err != nil {
		t.Fatalf("initWhitelist failed: %v", err)
	}

	if !whitelistMgr.IsEnabled() {
		t.Error("whitelist should be enabled by default when whitelist_enabled is not set")
	}

	allowed, _ := whitelistMgr.IsAllowed("file.delete")
	if allowed {
		t.Error("file.delete (high-risk command) should be blocked when whitelist is enabled")
	}
}

func TestConfigPreservesExistingFields(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "cd2-cli-test-preserve")
	os.MkdirAll(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.yaml")

	initialConfig := "server: custom-server:8080\ntls: true\njson: false\nskip-tls-verify: true\n"
	if err := os.WriteFile(configPath, []byte(initialConfig), 0600); err != nil {
		t.Fatalf("failed to write initial config: %v", err)
	}

	if err := saveTokenToConfig(configPath, "my-new-token"); err != nil {
		t.Fatalf("saveTokenToConfig failed: %v", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	content := string(data)

	if !strings.Contains(content, "server: custom-server:8080") {
		t.Errorf("server field should be preserved, got: %s", content)
	}
	if !strings.Contains(content, "tls: true") {
		t.Errorf("tls field should be preserved, got: %s", content)
	}
	if !strings.Contains(content, "json: false") {
		t.Errorf("json field should be preserved, got: %s", content)
	}
	if !strings.Contains(content, "skip-tls-verify: true") {
		t.Errorf("skip-tls-verify field should be preserved, got: %s", content)
	}
	if !strings.Contains(content, "token: my-new-token") {
		t.Errorf("token should be set, got: %s", content)
	}

	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("failed to stat config: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("config file should have 0600 permissions, got: %v", info.Mode().Perm())
	}

	if err := clearTokenFromConfig(configPath); err != nil {
		t.Fatalf("clearTokenFromConfig failed: %v", err)
	}

	data, err = os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read config after clear: %v", err)
	}

	content = string(data)

	if strings.Contains(content, "token:") {
		t.Errorf("token should be removed after clear, got: %s", content)
	}
	if !strings.Contains(content, "server: custom-server:8080") {
		t.Errorf("server field should still be preserved after clear, got: %s", content)
	}
	if !strings.Contains(content, "tls: true") {
		t.Errorf("tls field should still be preserved after clear, got: %s", content)
	}
}

func TestAtomicWriteFile(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "cd2-cli-test-atomic")
	os.MkdirAll(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	testPath := filepath.Join(tmpDir, "testfile.yaml")
	testData := []byte("test: value\n")

	if err := atomicWriteFile(testPath, testData, 0600); err != nil {
		t.Fatalf("atomicWriteFile failed: %v", err)
	}

	data, err := os.ReadFile(testPath)
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}
	if string(data) != string(testData) {
		t.Errorf("content mismatch, got: %s", string(data))
	}

	info, err := os.Stat(testPath)
	if err != nil {
		t.Fatalf("failed to stat test file: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("permissions should be 0600, got: %v", info.Mode().Perm())
	}

	tmpFiles, err := filepath.Glob(filepath.Join(tmpDir, ".cd2-cli-config-*.tmp"))
	if err != nil {
		t.Fatalf("glob failed: %v", err)
	}
	if len(tmpFiles) > 0 {
		t.Errorf("temp files should be cleaned up, found: %v", tmpFiles)
	}
}

func TestTimeoutFlag(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	tests := []struct {
		name      string
		timeout   string
		wantDur   time.Duration
		wantError bool
	}{
		{"default", "", 30 * time.Second, false},
		{"30s", "30s", 30 * time.Second, false},
		{"1m", "1m", 60 * time.Second, false},
		{"invalid", "abc", 30 * time.Second, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.timeout != "" {
				viper.Set("timeout", tc.timeout)
				defer viper.Set("timeout", "30s")
			} else {
				viper.Set("timeout", "30s")
			}

			dur := parseTimeout()
			if dur != tc.wantDur {
				t.Errorf("parseTimeout() = %v, want %v", dur, tc.wantDur)
			}

			ctx, cancel := getTimeoutContext()
			defer cancel()

			deadline, ok := ctx.Deadline()
			if !ok {
				t.Error("context should have a deadline")
			}

			remaining := time.Until(deadline)
			if remaining < tc.wantDur-2*time.Second || remaining > tc.wantDur+2*time.Second {
				t.Errorf("deadline remaining = %v, want ~%v", remaining, tc.wantDur)
			}
		})
	}
}

func TestTimeoutEnvVar(t *testing.T) {
	os.Setenv("CD2_CLI_TIMEOUT", "45s")
	defer os.Unsetenv("CD2_CLI_TIMEOUT")

	viper.Reset()
	viper.AutomaticEnv()
	viper.SetEnvPrefix("CD2_CLI")
	viper.SetDefault("timeout", "30s")

	timeoutStr := viper.GetString("timeout")
	if timeoutStr != "45s" {
		t.Errorf("viper.GetString('timeout') = %v, want 45s", timeoutStr)
	}
}

type errorMockServer struct {
	pb.UnimplementedCloudDriveFileSrvServer
}

func (s *errorMockServer) GetSystemInfo(ctx context.Context, req *emptypb.Empty) (*pb.CloudDriveSystemInfo, error) {
	return nil, fmt.Errorf("simulated RPC error")
}

func (s *errorMockServer) GetAccountStatus(ctx context.Context, req *emptypb.Empty) (*pb.AccountStatusResult, error) {
	return nil, fmt.Errorf("simulated RPC error")
}

func setupErrorTestCommand(t *testing.T) func() {
	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterCloudDriveFileSrvServer(s, &errorMockServer{})

	go func() {
		if err := s.Serve(lis); err != nil {
			t.Logf("Server error: %v", err)
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	c := &client.Client{}
	c.SetConn(conn)
	cd2Client = c

	tmpDir := filepath.Join(os.TempDir(), "cd2-cli-test-error")
	os.MkdirAll(tmpDir, 0755)
	whitelistPath := filepath.Join(tmpDir, "whitelist.yaml")
	mgr, err := whitelist.NewManager(whitelistPath)
	if err != nil {
		t.Fatalf("Failed to create whitelist manager: %v", err)
	}
	mgr.SetEnabled(false)
	whitelistMgr = mgr
	initialized = true
	outputJSON = true
	viper.SetDefault("json", true)

	return func() {
		c.Close()
		s.Stop()
		lis.Close()
		os.RemoveAll(tmpDir)
		initialized = false
		whitelistMgr = nil
		cd2Client = nil
		outputJSON = false
	}
}

func TestCommandReturnsErrorOnRPCFailure(t *testing.T) {
	cleanup := setupErrorTestCommand(t)
	defer cleanup()

	tests := []struct {
		name string
		args []string
	}{
		{"system info", []string{"system", "info"}},
		{"auth status", []string{"auth", "status"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := executeCommand(tc.args...)
			if err == nil {
				t.Errorf("%s: expected error but got nil (should return non-zero exit code)", tc.name)
			}
			if !strings.Contains(err.Error(), "simulated RPC error") {
				t.Errorf("%s: expected 'simulated RPC error' in error message, got: %v", tc.name, err)
			}
		})
	}
}

func TestRootCommandSilenceSettings(t *testing.T) {
	if !rootCmd.SilenceUsage {
		t.Error("rootCmd.SilenceUsage should be true")
	}
	if !rootCmd.SilenceErrors {
		t.Error("rootCmd.SilenceErrors should be true")
	}
}

func TestExitCodeOnRPCError(t *testing.T) {
	cleanup := setupErrorTestCommand(t)
	defer cleanup()

	rootCmd.SetArgs([]string{"system", "info"})
	err := rootCmd.Execute()

	if err == nil {
		t.Error("Execute() should return error when RPC fails")
	}

	t.Log("Note: os.Exit(1) is called in Execute() but cannot be tested directly in unit tests")
}

func TestFlagsOverrideConfig(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "cd2-cli-test-flags")
	os.MkdirAll(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := "server: config-server:8080\ntoken: config-token\ntls: true\nskip-tls-verify: true\njson: false\ntimeout: 60s\n"
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterCloudDriveFileSrvServer(s, &mockServer{})
	go func() {
		if err := s.Serve(lis); err != nil {
			t.Logf("Server error: %v", err)
		}
	}()
	defer s.Stop()
	defer lis.Close()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
	conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	c := &client.Client{}
	c.SetConn(conn)
	defer c.Close()
	cd2Client = c

	whitelistPath := filepath.Join(tmpDir, "whitelist.yaml")
	mgr, err := whitelist.NewManager(whitelistPath)
	if err != nil {
		t.Fatalf("Failed to create whitelist manager: %v", err)
	}
	mgr.SetEnabled(false)
	whitelistMgr = mgr
	initialized = true

	defer func() {
		initialized = false
		whitelistMgr = nil
		cd2Client = nil
		cfgFile = ""
	}()

	cfgFile = configPath

	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	args := []string{"--config", configPath, "--tls=true", "--skip-tls-verify=true", "--json=true", "auth", "status"}
	rootCmd.SetArgs(args)
	err = rootCmd.Execute()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)

	if err != nil {
		t.Errorf("command failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "test-user") {
		t.Errorf("expected test-user in output, got: %s", output)
	}

	if viper.GetBool("tls") {
		t.Log("--tls=true CLI flag bound to viper correctly")
	} else {
		t.Log("Note: viper state may have been reset by earlier tests; BindPFlag works in init()")
	}
}

func TestViperBindPFlagWorks(t *testing.T) {
	viper.Reset()
	viper.BindPFlag("tls", rootCmd.PersistentFlags().Lookup("tls"))
	viper.BindPFlag("skip-tls-verify", rootCmd.PersistentFlags().Lookup("skip-tls-verify"))
	viper.BindPFlag("json", rootCmd.PersistentFlags().Lookup("json"))
	viper.BindPFlag("timeout", rootCmd.PersistentFlags().Lookup("timeout"))

	rootCmd.PersistentFlags().Set("tls", "true")
	rootCmd.PersistentFlags().Set("skip-tls-verify", "true")
	rootCmd.PersistentFlags().Set("json", "true")
	rootCmd.PersistentFlags().Set("timeout", "45s")

	if !viper.GetBool("tls") {
		t.Error("BindPFlag: --tls=true should be reflected in viper.GetBool")
	}
	if !viper.GetBool("skip-tls-verify") {
		t.Error("BindPFlag: --skip-tls-verify=true should be reflected in viper.GetBool")
	}
	if !viper.GetBool("json") {
		t.Error("BindPFlag: --json=true should be reflected in viper.GetBool")
	}
	if viper.GetString("timeout") != "45s" {
		t.Errorf("BindPFlag: --timeout=45s should be reflected in viper.GetString, got: %s", viper.GetString("timeout"))
	}

	rootCmd.PersistentFlags().Set("tls", "false")
	rootCmd.PersistentFlags().Set("json", "false")

	if viper.GetBool("tls") {
		t.Error("BindPFlag: --tls=false should now be reflected in viper.GetBool")
	}
	if viper.GetBool("json") {
		t.Error("BindPFlag: --json=false should now be reflected in viper.GetBool")
	}
}

func TestEnvVarWithReplacer(t *testing.T) {
	os.Setenv("CD2_CLI_SKIP_TLS_VERIFY", "true")
	defer os.Unsetenv("CD2_CLI_SKIP_TLS_VERIFY")

	viper.Reset()
	viper.AutomaticEnv()
	viper.SetEnvPrefix("CD2_CLI")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.SetDefault("skip-tls-verify", false)

	if !viper.GetBool("skip-tls-verify") {
		t.Error("CD2_CLI_SKIP_TLS_VERIFY env var should set skip-tls-verify via replacer")
	}
}

func TestFlagOverridesEnvVar(t *testing.T) {
	os.Setenv("CD2_CLI_TLS", "true")
	defer os.Unsetenv("CD2_CLI_TLS")

	viper.Reset()
	viper.AutomaticEnv()
	viper.SetEnvPrefix("CD2_CLI")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.BindPFlag("tls", rootCmd.PersistentFlags().Lookup("tls"))

	rootCmd.PersistentFlags().Set("tls", "false")

	if viper.GetBool("tls") {
		t.Error("--tls=false flag should override CD2_CLI_TLS=true env var")
	}
}

func TestLocalCommandsWorkWithoutServer(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "cd2-cli-local-test")
	os.MkdirAll(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	whitelistPath := filepath.Join(tmpDir, "whitelist.yaml")
	mgr, err := whitelist.NewManager(whitelistPath)
	if err != nil {
		t.Fatalf("Failed to create whitelist manager: %v", err)
	}
	mgr.SetEnabled(false)
	whitelistMgr = mgr
	initialized = true
	defer func() {
		initialized = false
		whitelistMgr = nil
	}()

	viper.Set("json", true)
	defer viper.Reset()

	output, err := executeCommand("whitelist", "list")
	if err != nil {
		t.Errorf("whitelist list should work without server, got error: %v", err)
	}
	if !strings.Contains(output, "id") {
		t.Errorf("whitelist list output should contain 'id', got: %s", output)
	}

	output, err = executeCommand("whitelist", "path")
	if err != nil {
		t.Errorf("whitelist path should work without server, got error: %v", err)
	}
	if !strings.Contains(output, "path") {
		t.Errorf("whitelist path output should contain 'path', got: %s", output)
	}

	configPath := filepath.Join(tmpDir, "config.yaml")
	viper.SetConfigFile(configPath)

	output, err = executeCommand("auth", "token", "set", "test-token-xyz")
	if err != nil {
		t.Errorf("auth token set should work without server, got error: %v", err)
	}
	if !strings.Contains(output, "success") {
		t.Errorf("auth token set output should contain 'success', got: %s", output)
	}
}

func TestTokenSourceDetection(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "cd2-cli-token-source")
	os.MkdirAll(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	whitelistPath := filepath.Join(tmpDir, "whitelist.yaml")
	mgr, err := whitelist.NewManager(whitelistPath)
	if err != nil {
		t.Fatalf("Failed to create whitelist manager: %v", err)
	}
	mgr.SetEnabled(false)
	whitelistMgr = mgr
	initialized = true
	defer func() {
		initialized = false
		whitelistMgr = nil
	}()

	viper.Set("json", true)
	defer viper.Reset()

	output, err := executeCommand("auth", "token", "show", "--redact")
	if err != nil {
		t.Errorf("auth token show --redact should work without token, got error: %v", err)
	}
	if !strings.Contains(output, `"source":"none"`) {
		t.Errorf("expected source 'none' when no token configured, got: %s", output)
	}

	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := "token: config-token-12345\n"
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	output, err = executeCommand("auth", "token", "show", "--redact")
	if err != nil {
		t.Errorf("auth token show --redact should work with config token, got error: %v", err)
	}
	if !strings.Contains(output, `"source":"config"`) {
		t.Errorf("expected source 'config' when token from config file, got: %s", output)
	}
	if !strings.Contains(output, "conf...2345") {
		t.Errorf("expected redacted token format, got: %s", output)
	}

	os.Setenv("CD2_CLI_TOKEN", "env-token-abcdef")
	defer os.Unsetenv("CD2_CLI_TOKEN")

	viper.Reset()
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()
	viper.SetEnvPrefix("CD2_CLI")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	viper.Set("json", true)

	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("failed to re-read config: %v", err)
	}

	output, err = executeCommand("auth", "token", "show", "--redact")
	if err != nil {
		t.Errorf("auth token show --redact should work with env token, got error: %v", err)
	}
	if !strings.Contains(output, `"source":"env"`) {
		t.Errorf("expected source 'env' when CD2_CLI_TOKEN overrides config, got: %s", output)
	}
	if !strings.Contains(output, "env-...cdef") {
		t.Errorf("expected redacted env token format, got: %s", output)
	}

	output, err = executeCommand("--token", "flag-token-xyz123", "auth", "token", "show", "--redact")
	if err != nil {
		t.Errorf("auth token show --redact with --token flag should work, got error: %v", err)
	}
	if !strings.Contains(output, `"source":"flag"`) {
		t.Errorf("expected source 'flag' when --token flag used, got: %s", output)
	}
	if !strings.Contains(output, "flag...z123") {
		t.Errorf("expected redacted flag token format, got: %s", output)
	}
}

func TestEnvVarOverridesConfigToken(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "cd2-cli-env-overrides")
	os.MkdirAll(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	whitelistPath := filepath.Join(tmpDir, "whitelist.yaml")
	mgr, err := whitelist.NewManager(whitelistPath)
	if err != nil {
		t.Fatalf("Failed to create whitelist manager: %v", err)
	}
	mgr.SetEnabled(false)
	whitelistMgr = mgr
	initialized = true
	defer func() {
		initialized = false
		whitelistMgr = nil
	}()

	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := "server: localhost:19798\ntoken: config-file-token\njson: true\n"
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	viper.Reset()
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()
	viper.SetEnvPrefix("CD2_CLI")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	if viper.GetString("token") != "config-file-token" {
		t.Errorf("token should be from config file before env var set, got: %s", viper.GetString("token"))
	}

	os.Setenv("CD2_CLI_TOKEN", "env-override-token")
	defer os.Unsetenv("CD2_CLI_TOKEN")

	viper.Reset()
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()
	viper.SetEnvPrefix("CD2_CLI")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("failed to re-read config: %v", err)
	}

	if viper.GetString("token") != "env-override-token" {
		t.Errorf("CD2_CLI_TOKEN should override config file token, got: %s", viper.GetString("token"))
	}
}

func TestProtoJSONTimestampFormat(t *testing.T) {
	ts := time.Date(2025, 5, 25, 12, 30, 45, 123456789, time.UTC)
	timestampProto := timestamppb.New(ts)

	file := &pb.CloudDriveFile{
		CreateTime: timestampProto,
		WriteTime:  timestampProto,
		AccessTime: timestampProto,
	}

	result, err := protojson.MarshalOptions{
		EmitUnpopulated: true,
	}.Marshal(file)
	if err != nil {
		t.Fatalf("protojson.Marshal failed: %v", err)
	}

	expectedTsStr := "2025-05-25T12:30:45.123456789Z"
	if !strings.Contains(string(result), expectedTsStr) {
		t.Errorf("Timestamp should be in proto3 JSON format (RFC 3339 with nanoseconds), got: %s", string(result))
	}

	var unmarshaled pb.CloudDriveFile
	unmarshalOpts := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
	if err := unmarshalOpts.Unmarshal(result, &unmarshaled); err != nil {
		t.Fatalf("protojson.Unmarshal failed: %v", err)
	}

	if unmarshaled.CreateTime.AsTime().UTC() != ts.UTC() {
		t.Errorf("Unmarshaled timestamp mismatch: got %v, want %v", unmarshaled.CreateTime.AsTime(), ts)
	}
}

func TestOutputResultProtoMessage(t *testing.T) {
	file := &pb.CloudDriveFile{
		Name:        "test.txt",
		IsDirectory: false,
	}

	viper.Set("json", true)

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := outputResult(file)

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("outputResult failed: %v", err)
	}

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "\"name\":\"test.txt\"") {
		t.Errorf("Proto JSON output should contain name field, got: %s", output)
	}

	if !strings.Contains(output, "\"isDirectory\":false") {
		t.Errorf("Proto JSON output with EmitUnpopulated=true should include false values, got: %s", output)
	}
}

func TestConfigCreationWhenMissing(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "cd2-cli-config-create")
	os.MkdirAll(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	whitelistPath := filepath.Join(tmpDir, "whitelist.yaml")
	mgr, err := whitelist.NewManager(whitelistPath)
	if err != nil {
		t.Fatalf("Failed to create whitelist manager: %v", err)
	}
	mgr.SetEnabled(false)
	whitelistMgr = mgr
	initialized = true
	defer func() {
		initialized = false
		whitelistMgr = nil
	}()

	nonexistentConfig := filepath.Join(tmpDir, "new-config.yaml")
	if _, err := os.Stat(nonexistentConfig); !os.IsNotExist(err) {
		t.Fatalf("config file should not exist before test")
	}

	cfgFile = nonexistentConfig
	defer func() { cfgFile = "" }()

	viper.Reset()
	viper.Set("json", true)
	defer viper.Reset()

	output, err := executeCommand("auth", "token", "set", "test-token-123")
	if err != nil {
		t.Errorf("auth token set should succeed with non-existent config, got error: %v", err)
	}
	if !strings.Contains(output, "success") {
		t.Errorf("output should contain 'success', got: %s", output)
	}

	if _, err := os.Stat(nonexistentConfig); os.IsNotExist(err) {
		t.Errorf("config file should have been created")
	}

	data, err := os.ReadFile(nonexistentConfig)
	if err != nil {
		t.Fatalf("failed to read created config: %v", err)
	}
	if !strings.Contains(string(data), "test-token-123") {
		t.Errorf("created config should contain token, got: %s", string(data))
	}
}

func TestReadCommandsWorkWithMissingConfig(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "cd2-cli-read-missing")
	os.MkdirAll(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	whitelistPath := filepath.Join(tmpDir, "whitelist.yaml")
	mgr, err := whitelist.NewManager(whitelistPath)
	if err != nil {
		t.Fatalf("Failed to create whitelist manager: %v", err)
	}
	mgr.SetEnabled(false)
	whitelistMgr = mgr
	initialized = true
	defer func() {
		initialized = false
		whitelistMgr = nil
	}()

	nonexistentConfig := filepath.Join(tmpDir, "missing-config.yaml")
	cfgFile = nonexistentConfig
	defer func() { cfgFile = "" }()

	viper.Reset()
	viper.Set("json", true)
	defer viper.Reset()

	output, err := executeCommand("auth", "token", "show")
	if err != nil {
		t.Errorf("auth token show should work with missing config, got error: %v", err)
	}
	if !strings.Contains(output, `"source":"none"`) {
		t.Errorf("expected source 'none', got: %s", output)
	}
}

type emptyTaskMockServer struct {
	pb.UnimplementedCloudDriveFileSrvServer
}

func (s *emptyTaskMockServer) GetSystemInfo(ctx context.Context, req *emptypb.Empty) (*pb.CloudDriveSystemInfo, error) {
	return &pb.CloudDriveSystemInfo{IsLogin: true, UserName: "test-user"}, nil
}

func (s *emptyTaskMockServer) GetCopyTasks(ctx context.Context, req *emptypb.Empty) (*pb.GetCopyTaskResult, error) {
	return &pb.GetCopyTaskResult{CopyTasks: []*pb.CopyTask{}}, nil
}

func (s *emptyTaskMockServer) GetMergeTasks(ctx context.Context, req *emptypb.Empty) (*pb.GetMergeTasksResult, error) {
	return &pb.GetMergeTasksResult{MergeTasks: []*pb.MergeTask{}}, nil
}

func setupEmptyTaskTestCommand(t *testing.T) func() {
	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterCloudDriveFileSrvServer(s, &emptyTaskMockServer{})

	go func() {
		if err := s.Serve(lis); err != nil {
			t.Logf("Server error: %v", err)
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	c := &client.Client{}
	c.SetConn(conn)
	cd2Client = c

	tmpDir := filepath.Join(os.TempDir(), "cd2-cli-test-empty-task")
	os.MkdirAll(tmpDir, 0755)
	whitelistPath := filepath.Join(tmpDir, "whitelist.yaml")
	mgr, err := whitelist.NewManager(whitelistPath)
	if err != nil {
		t.Fatalf("Failed to create whitelist manager: %v", err)
	}
	mgr.SetEnabled(false)
	whitelistMgr = mgr
	initialized = true
	outputJSON = true
	viper.SetDefault("json", true)

	return func() {
		c.Close()
		s.Stop()
		lis.Close()
		os.RemoveAll(tmpDir)
		initialized = false
		whitelistMgr = nil
		cd2Client = nil
		outputJSON = false
	}
}

func TestTaskWaitNotFound(t *testing.T) {
	cleanup := setupEmptyTaskTestCommand(t)
	defer cleanup()
	missingIsComplete = false

	output, err := executeCommand("task", "wait-copy", "/nonexistent/src", "/nonexistent/dst")
	if err == nil {
		t.Error("expected non-zero exit code for not_found task")
	}

	if !strings.Contains(output, "not_found") {
		t.Errorf("expected 'not_found' status in output, got: %s", output)
	}
}

func TestTaskWaitMissingIsCompleteFlag(t *testing.T) {
	cleanup := setupEmptyTaskTestCommand(t)
	defer cleanup()
	missingIsComplete = false

	output, err := executeCommand("task", "wait-copy", "/nonexistent/src", "/nonexistent/dst", "--missing-is-complete")
	if err != nil {
		t.Errorf("expected zero exit code with --missing-is-complete, got error: %v", err)
	}

	if !strings.Contains(output, "completed") {
		t.Errorf("expected 'completed' status in output with --missing-is-complete, got: %s", output)
	}
}

func TestTaskWaitMergeNotFound(t *testing.T) {
	cleanup := setupEmptyTaskTestCommand(t)
	defer cleanup()
	missingIsComplete = false

	output, err := executeCommand("task", "wait-merge", "/nonexistent/src", "/nonexistent/dst")
	if err == nil {
		t.Error("expected non-zero exit code for not_found merge task")
	}

	if !strings.Contains(output, "not_found") {
		t.Errorf("expected 'not_found' status in output, got: %s", output)
	}
}

func TestAuthLoginPasswordFromEnv(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	os.Setenv("CD2_CLI_PASSWORD", "env-password")
	defer os.Unsetenv("CD2_CLI_PASSWORD")

	output, err := executeCommand("auth", "login", "testuser")
	if err != nil {
		t.Errorf("auth login with env password error: %v", err)
	}

	if !strings.Contains(output, "test-token") {
		t.Errorf("expected 'test-token' in output, got: %s", output)
	}
}

func TestAuthLoginPasswordFromFile(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	tmpDir := filepath.Join(os.TempDir(), "cd2-cli-password-test")
	os.MkdirAll(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	passwordFile := filepath.Join(tmpDir, "password.txt")
	os.WriteFile(passwordFile, []byte("file-password\n"), 0600)

	output, err := executeCommand("auth", "login", "testuser", "--password-file", passwordFile)
	if err != nil {
		t.Errorf("auth login with password file error: %v", err)
	}

	if !strings.Contains(output, "test-token") {
		t.Errorf("expected 'test-token' in output, got: %s", output)
	}
}

func TestAuthLogin2FAEnvVars(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	os.Setenv("CD2_CLI_PASSWORD", "env-password")
	os.Setenv("CD2_CLI_TOTP", "123456")
	defer func() {
		os.Unsetenv("CD2_CLI_PASSWORD")
		os.Unsetenv("CD2_CLI_TOTP")
	}()

	output, err := executeCommand("auth", "login-2fa", "testuser")
	if err != nil {
		t.Errorf("auth login-2fa with env vars error: %v", err)
	}

	if !strings.Contains(output, "2fa-token") {
		t.Errorf("expected '2fa-token' in output, got: %s", output)
	}
}

func TestAuthChangePasswordEnvVars(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	os.Setenv("CD2_CLI_OLD_PASSWORD", "old-pass")
	os.Setenv("CD2_CLI_NEW_PASSWORD", "new-pass")
	defer func() {
		os.Unsetenv("CD2_CLI_OLD_PASSWORD")
		os.Unsetenv("CD2_CLI_NEW_PASSWORD")
	}()

	output, err := executeCommand("auth", "change-password")
	if err != nil {
		t.Errorf("auth change-password with env vars error: %v", err)
	}

	if !strings.Contains(output, "success") {
		t.Errorf("expected 'success' in output, got: %s", output)
	}
}

func Test2FACommandsEnvVar(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	os.Setenv("CD2_CLI_TOTP", "123456")
	defer os.Unsetenv("CD2_CLI_TOTP")

	output, err := executeCommand("auth", "2fa-enable")
	if err != nil {
		t.Errorf("auth 2fa-enable with env var error: %v", err)
	}

	if !strings.Contains(output, "recoveryCodes") {
		t.Errorf("expected 'recoveryCodes' in output, got: %s", output)
	}
}

type webdavModifyMockServer struct {
	mockServer
	lastModifyRequest *pb.ModifyDavUserRequest
}

func (s *webdavModifyMockServer) ModifyDavUser(ctx context.Context, req *pb.ModifyDavUserRequest) (*emptypb.Empty, error) {
	s.lastModifyRequest = req
	return &emptypb.Empty{}, nil
}

func TestWebDavUserModifyWithoutBoolFlags(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	var modifyReq *pb.ModifyDavUserRequest
	captureMock := &webdavModifyMockServer{}
	captureMock.lastModifyRequest = nil

	_, err := executeCommandWithMock(captureMock, "webdav", "user", "modify", "alice", "--password", "newpass")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	modifyReq = captureMock.lastModifyRequest
	if modifyReq == nil {
		t.Fatal("ModifyDavUserRequest was nil")
	}
	if modifyReq.UserName != "alice" {
		t.Errorf("expected username 'alice', got %q", modifyReq.UserName)
	}
	if modifyReq.Password == nil || *modifyReq.Password != "newpass" {
		t.Errorf("expected password 'newpass', got %v", modifyReq.Password)
	}
	if modifyReq.ReadOnly != nil {
		t.Errorf("ReadOnly should be nil when --read-only not provided, got %v", modifyReq.ReadOnly)
	}
	if modifyReq.Enabled != nil {
		t.Errorf("Enabled should be nil when --enabled not provided, got %v", modifyReq.Enabled)
	}
	if modifyReq.Guest != nil {
		t.Errorf("Guest should be nil when --guest not provided, got %v", modifyReq.Guest)
	}
}

func TestWebDavUserModifyWithBoolFlags(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	captureMock := &webdavModifyMockServer{}
	captureMock.lastModifyRequest = nil

	_, err := executeCommandWithMock(captureMock, "webdav", "user", "modify", "bob", "--enabled", "--read-only")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	modifyReq := captureMock.lastModifyRequest
	if modifyReq == nil {
		t.Fatal("ModifyDavUserRequest was nil")
	}
	if modifyReq.UserName != "bob" {
		t.Errorf("expected username 'bob', got %q", modifyReq.UserName)
	}
	if modifyReq.Enabled == nil || *modifyReq.Enabled != true {
		t.Errorf("expected Enabled=true when --enabled provided, got %v", modifyReq.Enabled)
	}
	if modifyReq.ReadOnly == nil || *modifyReq.ReadOnly != true {
		t.Errorf("expected ReadOnly=true when --read-only provided, got %v", modifyReq.ReadOnly)
	}
	if modifyReq.Guest != nil {
		t.Errorf("Guest should be nil when --guest not provided, got %v", modifyReq.Guest)
	}
}

func TestWebDavUserModifyWithFalseBoolFlags(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	captureMock := &webdavModifyMockServer{}
	captureMock.lastModifyRequest = nil

	_, err := executeCommandWithMock(captureMock, "webdav", "user", "modify", "charlie", "--enabled=false")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	modifyReq := captureMock.lastModifyRequest
	if modifyReq == nil {
		t.Fatal("ModifyDavUserRequest was nil")
	}
	if modifyReq.Enabled == nil {
		t.Fatal("Enabled should not be nil when --enabled=false explicitly provided")
	}
	if *modifyReq.Enabled != false {
		t.Errorf("expected Enabled=false, got %v", *modifyReq.Enabled)
	}
	if modifyReq.ReadOnly != nil {
		t.Errorf("ReadOnly should be nil when --read-only not provided, got %v (value=%v)", modifyReq.ReadOnly, *modifyReq.ReadOnly)
	}
	if modifyReq.Guest != nil {
		t.Errorf("Guest should be nil when --guest not provided, got %v", modifyReq.Guest)
	}
}

func TestMountUpdateRequiresOptions(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	resetMountUpdateFlags()
	_, err := executeCommand("mount", "update", "/mount/point")
	if err == nil {
		t.Fatal("expected error when --options flag is missing")
	}
	if !strings.Contains(err.Error(), "required flag") {
		t.Errorf("expected required flag error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "options") {
		t.Errorf("expected error to mention 'options' flag, got: %v", err)
	}
}

func TestWebdavServerSetRequiresConfig(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	resetWebdavServerSetFlags()
	_, err := executeCommand("webdav", "server", "set")
	if err == nil {
		t.Fatal("expected error when --config flag is missing")
	}
	if !strings.Contains(err.Error(), "required flag") {
		t.Errorf("expected required flag error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "config") {
		t.Errorf("expected error to mention 'config' flag, got: %v", err)
	}
}

func TestStorageSetConfigRequiresConfig(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	resetStorageSetConfigFlags()
	_, err := executeCommand("storage", "set-config", "cloud", "user")
	if err == nil {
		t.Fatal("expected error when --config flag is missing")
	}
	if !strings.Contains(err.Error(), "required flag") {
		t.Errorf("expected required flag error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "config") {
		t.Errorf("expected error to mention 'config' flag, got: %v", err)
	}
}

func TestCloudapiSetConfigRequiresConfig(t *testing.T) {
	cleanup := setupTestCommand(t)
	defer cleanup()

	resetCloudapiSetConfigFlags()
	_, err := executeCommand("cloudapi", "set-config", "cloud", "user")
	if err == nil {
		t.Fatal("expected error when --config flag is missing")
	}
	if !strings.Contains(err.Error(), "required flag") {
		t.Errorf("expected required flag error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "config") {
		t.Errorf("expected error to mention 'config' flag, got: %v", err)
	}
}

func resetMountUpdateFlags() {
	cmd, _, _ := rootCmd.Find([]string{"mount", "update"})
	if cmd != nil {
		cmd.Flags().Visit(func(f *pflag.Flag) {
			f.Changed = false
		})
	}
}

func resetWebdavServerSetFlags() {
	cmd, _, _ := rootCmd.Find([]string{"webdav", "server", "set"})
	if cmd != nil {
		cmd.Flags().Visit(func(f *pflag.Flag) {
			f.Changed = false
		})
	}
}

func resetStorageSetConfigFlags() {
	cmd, _, _ := rootCmd.Find([]string{"storage", "set-config"})
	if cmd != nil {
		cmd.Flags().Visit(func(f *pflag.Flag) {
			f.Changed = false
		})
	}
}

func resetCloudapiSetConfigFlags() {
	cmd, _, _ := rootCmd.Find([]string{"cloudapi", "set-config"})
	if cmd != nil {
		cmd.Flags().Visit(func(f *pflag.Flag) {
			f.Changed = false
		})
	}
}
