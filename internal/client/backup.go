package client

import (
	"context"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

type BackupAPI struct {
	c *Client
}

func (a *BackupAPI) BackupGetAll(ctx context.Context) (*pb.BackupList, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.BackupGetAll(ctx, &emptypb.Empty{})
}

func (a *BackupAPI) BackupGetStatus(ctx context.Context, backupId string) (*pb.BackupStatus, error) {
	ctx = a.c.withAuth(ctx)
	req := &pb.StringValue{Value: backupId}
	return a.c.client.BackupGetStatus(ctx, req)
}

func (a *BackupAPI) BackupAdd(ctx context.Context, req *pb.Backup) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.BackupAdd(ctx, req)
	return err
}

func (a *BackupAPI) BackupRemove(ctx context.Context, backupId string) error {
	ctx = a.c.withAuth(ctx)
	req := &pb.StringValue{Value: backupId}
	_, err := a.c.client.BackupRemove(ctx, req)
	return err
}

func (a *BackupAPI) BackupUpdate(ctx context.Context, req *pb.Backup) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.BackupUpdate(ctx, req)
	return err
}

func (a *BackupAPI) BackupAddDestination(ctx context.Context, req *pb.BackupModifyRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.BackupAddDestination(ctx, req)
	return err
}

func (a *BackupAPI) BackupRemoveDestination(ctx context.Context, req *pb.BackupModifyRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.BackupRemoveDestination(ctx, req)
	return err
}

func (a *BackupAPI) BackupSetEnabled(ctx context.Context, req *pb.BackupSetEnabledRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.BackupSetEnabled(ctx, req)
	return err
}

func (a *BackupAPI) BackupSetFileSystemWatchEnabled(ctx context.Context, req *pb.BackupModifyRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.BackupSetFileSystemWatchEnabled(ctx, req)
	return err
}

func (a *BackupAPI) BackupUpdateStrategies(ctx context.Context, req *pb.BackupModifyRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.BackupUpdateStrategies(ctx, req)
	return err
}

func (a *BackupAPI) BackupRestartWalkingThrough(ctx context.Context, backupId string) error {
	ctx = a.c.withAuth(ctx)
	req := &pb.StringValue{Value: backupId}
	_, err := a.c.client.BackupRestartWalkingThrough(ctx, req)
	return err
}

func (a *BackupAPI) CanAddMoreBackups(ctx context.Context) (*pb.FileOperationResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.CanAddMoreBackups(ctx, &emptypb.Empty{})
}

func (a *BackupAPI) NotifyPhotoLibraryChanges(ctx context.Context, req *pb.PhotoLibraryChangeList) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.NotifyPhotoLibraryChanges(ctx, req)
	return err
}
