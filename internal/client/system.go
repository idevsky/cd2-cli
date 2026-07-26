package client

import (
	"context"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

type SystemAPI struct {
	c *Client
}

func (a *SystemAPI) GetRuntimeInfo(ctx context.Context) (*pb.RuntimeInfo, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetRuntimeInfo(ctx, &emptypb.Empty{})
}

func (a *SystemAPI) GetRunningInfo(ctx context.Context) (*pb.RunInfo, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetRunningInfo(ctx, &emptypb.Empty{})
}

func (a *SystemAPI) GetOpenFileHandles(ctx context.Context) (*pb.OpenFileHandleList, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetOpenFileHandles(ctx, &emptypb.Empty{})
}

func (a *SystemAPI) GetFileBufferDiskCacheStats(ctx context.Context) (*pb.FileBufferDiskCacheStats, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetFileBufferDiskCacheStats(ctx, &emptypb.Empty{})
}

func (a *SystemAPI) PurgeFileBufferDiskCache(ctx context.Context) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.PurgeFileBufferDiskCache(ctx, &emptypb.Empty{})
	return err
}

func (a *SystemAPI) SetDiskCacheEvictionStrategy(ctx context.Context, req *pb.SetDiskCacheEvictionStrategyRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.SetDiskCacheEvictionStrategy(ctx, req)
	return err
}

func (a *SystemAPI) SetFolderDiskCache(ctx context.Context, req *pb.SetFolderDiskCacheRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.SetFolderDiskCache(ctx, req)
	return err
}

func (a *SystemAPI) RemoveFolderDiskCache(ctx context.Context, path string) error {
	ctx = a.c.withAuth(ctx)
	req := &pb.FileRequest{Path: path}
	_, err := a.c.client.RemoveFolderDiskCache(ctx, req)
	return err
}

func (a *SystemAPI) ListDiskCacheFolders(ctx context.Context) (*pb.ListDiskCacheFoldersReply, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.ListDiskCacheFolders(ctx, &emptypb.Empty{})
}

func (a *SystemAPI) PrefetchFileRanges(ctx context.Context, req *pb.PrefetchFileRangesRequest) (*pb.PrefetchFileRangesReply, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.PrefetchFileRanges(ctx, req)
}

func (a *SystemAPI) CancelFilePrefetch(ctx context.Context, req *pb.CancelFilePrefetchRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.CancelFilePrefetch(ctx, req)
	return err
}

func (a *SystemAPI) CloseFileReader(ctx context.Context, path string) error {
	ctx = a.c.withAuth(ctx)
	req := &pb.FileRequest{Path: path}
	_, err := a.c.client.CloseFileReader(ctx, req)
	return err
}

func (a *SystemAPI) GetActivePrefetchHints(ctx context.Context) (*pb.GetActivePrefetchHintsReply, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetActivePrefetchHints(ctx, &emptypb.Empty{})
}

func (a *SystemAPI) GetSystemSettings(ctx context.Context) (*pb.SystemSettings, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetSystemSettings(ctx, &emptypb.Empty{})
}

func (a *SystemAPI) SetSystemSettings(ctx context.Context, req *pb.SystemSettings) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.SetSystemSettings(ctx, req)
	return err
}

func (a *SystemAPI) SetDirCacheTimeSecs(ctx context.Context, req *pb.SetDirCacheTimeRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.SetDirCacheTimeSecs(ctx, req)
	return err
}

func (a *SystemAPI) GetEffectiveDirCacheTimeSecs(ctx context.Context, req *pb.GetEffectiveDirCacheTimeRequest) (*pb.GetEffectiveDirCacheTimeResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetEffectiveDirCacheTimeSecs(ctx, req)
}

func (a *SystemAPI) ForceExpireDirCache(ctx context.Context, path string) error {
	ctx = a.c.withAuth(ctx)
	req := &pb.FileRequest{Path: path}
	_, err := a.c.client.ForceExpireDirCache(ctx, req)
	return err
}

func (a *SystemAPI) VacuumDirCache(ctx context.Context) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.VacuumDirCache(ctx, &emptypb.Empty{})
	return err
}

func (a *SystemAPI) GetVacuumProgress(ctx context.Context) (*pb.VacuumProgressResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetVacuumProgress(ctx, &emptypb.Empty{})
}

func (a *SystemAPI) GetDirCacheDbSize(ctx context.Context) (*pb.GetDirCacheDbSizeResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetDirCacheDbSize(ctx, &emptypb.Empty{})
}

func (a *SystemAPI) GetOpenFileTable(ctx context.Context, req *pb.GetOpenFileTableRequest) (*pb.OpenFileTable, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetOpenFileTable(ctx, req)
}

func (a *SystemAPI) GetDirCacheTable(ctx context.Context) (*pb.DirCacheTable, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetDirCacheTable(ctx, &emptypb.Empty{})
}

func (a *SystemAPI) GetReferencedEntryPaths(ctx context.Context, path string) (*pb.StringList, error) {
	ctx = a.c.withAuth(ctx)
	req := &pb.FileRequest{Path: path}
	return a.c.client.GetReferencedEntryPaths(ctx, req)
}

func (a *SystemAPI) GetTempFileTable(ctx context.Context) (*pb.TempFileTable, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetTempFileTable(ctx, &emptypb.Empty{})
}

func (a *SystemAPI) GetServiceCapabilities(ctx context.Context) (*pb.ServiceCapabilities, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetServiceCapabilities(ctx, &emptypb.Empty{})
}

func (a *SystemAPI) RestartService(ctx context.Context) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.RestartService(ctx, &emptypb.Empty{})
	return err
}

func (a *SystemAPI) ShutdownService(ctx context.Context) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.ShutdownService(ctx, &emptypb.Empty{})
	return err
}

func (a *SystemAPI) HasUpdate(ctx context.Context) (*pb.UpdateResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.HasUpdate(ctx, &emptypb.Empty{})
}

func (a *SystemAPI) CheckUpdate(ctx context.Context) (*pb.UpdateResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.CheckUpdate(ctx, &emptypb.Empty{})
}

func (a *SystemAPI) DownloadUpdate(ctx context.Context) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.DownloadUpdate(ctx, &emptypb.Empty{})
	return err
}

func (a *SystemAPI) UpdateSystem(ctx context.Context) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.UpdateSystem(ctx, &emptypb.Empty{})
	return err
}

func (a *SystemAPI) GetWebServerConfig(ctx context.Context) (*pb.WebServerConfig, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetWebServerConfig(ctx, &emptypb.Empty{})
}

func (a *SystemAPI) SetWebServerConfig(ctx context.Context, req *pb.SetWebServerConfigRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.SetWebServerConfig(ctx, req)
	return err
}

func (a *SystemAPI) GenerateSelfSignedCert(ctx context.Context, req *pb.GenerateSelfSignedCertRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.GenerateSelfSignedCert(ctx, req)
	return err
}

func (a *SystemAPI) GetMachineId(ctx context.Context) (*pb.StringResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetMachineId(ctx, &emptypb.Empty{})
}

func (a *SystemAPI) GetOnlineDevices(ctx context.Context) (*pb.OnlineDevices, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetOnlineDevices(ctx, &emptypb.Empty{})
}

func (a *SystemAPI) KickoutDevice(ctx context.Context, req *pb.DeviceRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.KickoutDevice(ctx, req)
	return err
}

func (a *SystemAPI) ListLogFiles(ctx context.Context) (*pb.ListLogFileResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.ListLogFiles(ctx, &emptypb.Empty{})
}

func (a *SystemAPI) TestUpdate(ctx context.Context, req *pb.FileRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.TestUpdate(ctx, req)
	return err
}
