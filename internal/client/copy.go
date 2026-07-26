package client

import (
	"context"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

type CopyAPI struct {
	c *Client
}

func (a *CopyAPI) GetCopyTasks(ctx context.Context) (*pb.GetCopyTaskResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetCopyTasks(ctx, &emptypb.Empty{})
}

func (a *CopyAPI) GetMergeTasks(ctx context.Context) (*pb.GetMergeTasksResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetMergeTasks(ctx, &emptypb.Empty{})
}

func (a *CopyAPI) CancelMergeTask(ctx context.Context, req *pb.CancelMergeTaskRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.CancelMergeTask(ctx, req)
	return err
}

func (a *CopyAPI) CancelCopyTask(ctx context.Context, req *pb.CopyTaskRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.CancelCopyTask(ctx, req)
	return err
}

func (a *CopyAPI) PauseCopyTask(ctx context.Context, req *pb.PauseCopyTaskRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.PauseCopyTask(ctx, req)
	return err
}

func (a *CopyAPI) RestartCopyTask(ctx context.Context, req *pb.CopyTaskRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.RestartCopyTask(ctx, req)
	return err
}

func (a *CopyAPI) RemoveCompletedCopyTasks(ctx context.Context) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.RemoveCompletedCopyTasks(ctx, &emptypb.Empty{})
	return err
}

func (a *CopyAPI) RemoveAllCopyTasks(ctx context.Context) (*pb.BatchOperationResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.RemoveAllCopyTasks(ctx, &emptypb.Empty{})
}

func (a *CopyAPI) RemoveCopyTasks(ctx context.Context, req *pb.CopyTaskBatchRequest) (*pb.BatchOperationResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.RemoveCopyTasks(ctx, req)
}

func (a *CopyAPI) PauseAllCopyTasks(ctx context.Context, req *pb.PauseAllCopyTasksRequest) (*pb.BatchOperationResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.PauseAllCopyTasks(ctx, req)
}

func (a *CopyAPI) PauseCopyTasks(ctx context.Context, req *pb.PauseCopyTasksRequest) (*pb.BatchOperationResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.PauseCopyTasks(ctx, req)
}

func (a *CopyAPI) ResumeAllCopyTasks(ctx context.Context) (*pb.BatchOperationResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.ResumeAllCopyTasks(ctx, &emptypb.Empty{})
}

func (a *CopyAPI) ResumeCopyTasks(ctx context.Context, req *pb.CopyTaskBatchRequest) (*pb.BatchOperationResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.ResumeCopyTasks(ctx, req)
}
