package client

import (
	"context"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

type MountAPI struct {
	c *Client
}

func (a *MountAPI) CanAddMoreMountPoints(ctx context.Context) (*pb.FileOperationResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.CanAddMoreMountPoints(ctx, &emptypb.Empty{})
}

func (a *MountAPI) GetMountPoints(ctx context.Context) (*pb.GetMountPointsResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetMountPoints(ctx, &emptypb.Empty{})
}

func (a *MountAPI) AddMountPoint(ctx context.Context, req *pb.MountOption) (*pb.MountPointResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.AddMountPoint(ctx, req)
}

func (a *MountAPI) RemoveMountPoint(ctx context.Context, req *pb.MountPointRequest) (*pb.MountPointResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.RemoveMountPoint(ctx, req)
}

func (a *MountAPI) Mount(ctx context.Context, req *pb.MountPointRequest) (*pb.MountPointResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.Mount(ctx, req)
}

func (a *MountAPI) Unmount(ctx context.Context, req *pb.MountPointRequest) (*pb.MountPointResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.Unmount(ctx, req)
}

func (a *MountAPI) UpdateMountPoint(ctx context.Context, req *pb.UpdateMountPointRequest) (*pb.MountPointResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.UpdateMountPoint(ctx, req)
}

func (a *MountAPI) GetAvailableDriveLetters(ctx context.Context) (*pb.GetAvailableDriveLettersResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetAvailableDriveLetters(ctx, &emptypb.Empty{})
}

func (a *MountAPI) HasDriveLetters(ctx context.Context) (*pb.HasDriveLettersResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.HasDriveLetters(ctx, &emptypb.Empty{})
}

func (a *MountAPI) CanMountBothLocalAndCloud(ctx context.Context) (*pb.BoolResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.CanMountBothLocalAndCloud(ctx, &emptypb.Empty{})
}
