package client

import (
	"context"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

type LocalAPI struct {
	c *Client
}

func (a *LocalAPI) LocalGetSubFiles(ctx context.Context, req *pb.LocalGetSubFilesRequest) ([]*pb.LocalGetSubFilesResult, error) {
	ctx = a.c.withAuth(ctx)
	stream, err := a.c.client.LocalGetSubFiles(ctx, req)
	if err != nil {
		return nil, err
	}
	return collectStream[*pb.LocalGetSubFilesResult](stream)
}

func (a *LocalAPI) LocalCreateFolder(ctx context.Context, req *pb.LocalCreateFolderRequest) (*pb.LocalCreateFolderResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.LocalCreateFolder(ctx, req)
}
