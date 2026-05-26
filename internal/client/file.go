package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

type FileAPI struct {
	c *Client
}

func (a *FileAPI) GetSubFiles(ctx context.Context, req *pb.ListSubFileRequest) ([]*pb.CloudDriveFile, error) {
	ctx = a.c.withAuth(ctx)
	stream, err := a.c.client.GetSubFiles(ctx, req)
	if err != nil {
		return nil, err
	}

	var files []*pb.CloudDriveFile
	for {
		reply, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		files = append(files, reply.SubFiles...)
	}
	return files, nil
}

func (a *FileAPI) GetSearchResults(ctx context.Context, req *pb.SearchRequest) ([]*pb.CloudDriveFile, error) {
	ctx = a.c.withAuth(ctx)
	stream, err := a.c.client.GetSearchResults(ctx, req)
	if err != nil {
		return nil, err
	}

	var files []*pb.CloudDriveFile
	for {
		reply, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		files = append(files, reply.SubFiles...)
	}
	return files, nil
}

func (a *FileAPI) FindFileByPath(ctx context.Context, path string) (*pb.CloudDriveFile, error) {
	ctx = a.c.withAuth(ctx)
	req := &pb.FindFileByPathRequest{Path: path}
	return a.c.client.FindFileByPath(ctx, req)
}

func (a *FileAPI) CreateFolder(ctx context.Context, req *pb.CreateFolderRequest) (*pb.CreateFolderResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.CreateFolder(ctx, req)
}

func (a *FileAPI) CreateEncryptedFolder(ctx context.Context, req *pb.CreateEncryptedFolderRequest) (*pb.CreateFolderResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.CreateEncryptedFolder(ctx, req)
}

func (a *FileAPI) UnlockEncryptedFile(ctx context.Context, req *pb.UnlockEncryptedFileRequest) (*pb.FileOperationResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.UnlockEncryptedFile(ctx, req)
}

func (a *FileAPI) LockEncryptedFile(ctx context.Context, path string) (*pb.FileOperationResult, error) {
	ctx = a.c.withAuth(ctx)
	req := &pb.FileRequest{Path: path}
	return a.c.client.LockEncryptedFile(ctx, req)
}

func (a *FileAPI) RenameFile(ctx context.Context, req *pb.RenameFileRequest) (*pb.FileOperationResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.RenameFile(ctx, req)
}

func (a *FileAPI) RenameFiles(ctx context.Context, req *pb.RenameFilesRequest) (*pb.FileOperationResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.RenameFiles(ctx, req)
}

func (a *FileAPI) MoveFile(ctx context.Context, req *pb.MoveFileRequest) (*pb.FileOperationResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.MoveFile(ctx, req)
}

func (a *FileAPI) CopyFile(ctx context.Context, req *pb.CopyFileRequest) (*pb.FileOperationResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.CopyFile(ctx, req)
}

func (a *FileAPI) DeleteFile(ctx context.Context, path string) (*pb.FileOperationResult, error) {
	ctx = a.c.withAuth(ctx)
	req := &pb.FileRequest{Path: path}
	return a.c.client.DeleteFile(ctx, req)
}

func (a *FileAPI) DeleteFilePermanently(ctx context.Context, path string) (*pb.FileOperationResult, error) {
	ctx = a.c.withAuth(ctx)
	req := &pb.FileRequest{Path: path}
	return a.c.client.DeleteFilePermanently(ctx, req)
}

func (a *FileAPI) DeleteFiles(ctx context.Context, paths []string) (*pb.FileOperationResult, error) {
	ctx = a.c.withAuth(ctx)
	req := &pb.MultiFileRequest{Path: paths}
	return a.c.client.DeleteFiles(ctx, req)
}

func (a *FileAPI) DeleteFilesPermanently(ctx context.Context, paths []string) (*pb.FileOperationResult, error) {
	ctx = a.c.withAuth(ctx)
	req := &pb.MultiFileRequest{Path: paths}
	return a.c.client.DeleteFilesPermanently(ctx, req)
}

func (a *FileAPI) AddSharedLink(ctx context.Context, req *pb.AddSharedLinkRequest) error {
	ctx = a.c.withAuth(ctx)
	_, err := a.c.client.AddSharedLink(ctx, req)
	return err
}

func (a *FileAPI) GetFileDetailProperties(ctx context.Context, path string) (*pb.FileDetailProperties, error) {
	ctx = a.c.withAuth(ctx)
	req := &pb.FileRequest{Path: path}
	return a.c.client.GetFileDetailProperties(ctx, req)
}

func (a *FileAPI) GetSpaceInfo(ctx context.Context, path string) (*pb.SpaceInfo, error) {
	ctx = a.c.withAuth(ctx)
	req := &pb.FileRequest{Path: path}
	return a.c.client.GetSpaceInfo(ctx, req)
}

func (a *FileAPI) GetCloudMemberships(ctx context.Context, path string) (*pb.CloudMemberships, error) {
	ctx = a.c.withAuth(ctx)
	req := &pb.FileRequest{Path: path}
	return a.c.client.GetCloudMemberships(ctx, req)
}

func (a *FileAPI) GetMetaData(ctx context.Context, path string) (*pb.FileMetaData, error) {
	ctx = a.c.withAuth(ctx)
	req := &pb.FileRequest{Path: path}
	return a.c.client.GetMetaData(ctx, req)
}

func (a *FileAPI) GetOriginalPath(ctx context.Context, path string) (*pb.StringResult, error) {
	ctx = a.c.withAuth(ctx)
	req := &pb.FileRequest{Path: path}
	return a.c.client.GetOriginalPath(ctx, req)
}

func (a *FileAPI) CreateFile(ctx context.Context, req *pb.CreateFileRequest) (*pb.CreateFileResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.CreateFile(ctx, req)
}

func (a *FileAPI) CloseFile(ctx context.Context, req *pb.CloseFileRequest) (*pb.FileOperationResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.CloseFile(ctx, req)
}

func (a *FileAPI) WriteToFile(ctx context.Context, req *pb.WriteFileRequest) (*pb.WriteFileResult, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.WriteToFile(ctx, req)
}

func (a *FileAPI) WriteToFileStream(ctx context.Context) (pb.CloudDriveFileSrv_WriteToFileStreamClient, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.WriteToFileStream(ctx)
}

func (a *FileAPI) GetDownloadUrl(ctx context.Context, req *pb.GetDownloadUrlPathRequest) (*pb.DownloadUrlPathInfo, error) {
	ctx = a.c.withAuth(ctx)
	return a.c.client.GetDownloadUrlPath(ctx, req)
}

func (a *FileAPI) UploadLocalFile(ctx context.Context, localPath, remotePath string) (map[string]interface{}, error) {
	ctx = a.c.withAuth(ctx)

	file, err := os.Open(localPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	parentPath := filepath.Dir(remotePath)
	fileName := filepath.Base(remotePath)

	createResult, err := a.c.client.CreateFile(ctx, &pb.CreateFileRequest{
		ParentPath: parentPath,
		FileName:   fileName,
	})
	if err != nil {
		return nil, err
	}

	fileHandle := createResult.FileHandle
	streamClosed := false
	fileClosed := false

	defer func() {
		if !fileClosed && fileHandle != 0 {
			closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			_, closeErr := a.c.client.CloseFile(closeCtx, &pb.CloseFileRequest{
				FileHandle: fileHandle,
			})
			if closeErr != nil {
				fmt.Printf("warning: failed to close file handle on cleanup: %v\n", closeErr)
			}
		}
	}()

	stream, err := a.c.client.WriteToFileStream(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if !streamClosed {
			stream.CloseSend()
		}
	}()

	buf := make([]byte, 64*1024)
	var pos uint64 = 0
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		err = stream.Send(&pb.WriteFileRequest{
			FileHandle: fileHandle,
			StartPos:   pos,
			Length:     uint64(n),
			Buffer:     buf[:n],
		})
		if err != nil {
			return nil, err
		}
		pos += uint64(n)
	}

	err = stream.Send(&pb.WriteFileRequest{
		FileHandle: fileHandle,
		CloseFile:  true,
	})
	if err != nil {
		return nil, err
	}

	writeResult, err := stream.CloseAndRecv()
	if err != nil {
		return nil, err
	}

	streamClosed = true
	fileClosed = true

	return map[string]interface{}{
		"success":    true,
		"bytes":      pos,
		"fileHandle": fileHandle,
		"result":     writeResult,
	}, nil
}

func (a *FileAPI) DownloadRemoteFile(ctx context.Context, remotePath, localPath string) (map[string]interface{}, error) {
	ctx = a.c.withAuth(ctx)

	urlInfo, err := a.c.client.GetDownloadUrlPath(ctx, &pb.GetDownloadUrlPathRequest{
		Path:         remotePath,
		GetDirectUrl: true,
	})
	if err != nil {
		return nil, err
	}

	downloadURL := urlInfo.DownloadUrlPath
	if urlInfo.DirectUrl != nil && *urlInfo.DirectUrl != "" {
		downloadURL = *urlInfo.DirectUrl
	}

	if downloadURL == "" {
		return nil, fmt.Errorf("no download URL available for remote path: %s", remotePath)
	}

	dir := filepath.Dir(localPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create local directory: %w", err)
	}

	tmpPath := localPath + ".tmp"

	httpClient := &http.Client{
		Timeout: 30 * time.Minute,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: a.c.skipVerifyTLS,
			},
		},
	}

	req, err := http.NewRequestWithContext(ctx, "GET", downloadURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed with status: %s", resp.Status)
	}

	tmpFile, err := os.Create(tmpPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}

	written, err := io.Copy(tmpFile, resp.Body)
	tmpFile.Close()
	if err != nil {
		os.Remove(tmpPath)
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	if err := os.Rename(tmpPath, localPath); err != nil {
		os.Remove(tmpPath)
		return nil, fmt.Errorf("failed to rename temp file: %w", err)
	}

	return map[string]interface{}{
		"success":     true,
		"bytes":       written,
		"localPath":   localPath,
		"remotePath":  remotePath,
		"downloadUrl": urlInfo.DownloadUrlPath,
		"directUrl":   urlInfo.DirectUrl,
		"expiresIn":   urlInfo.ExpiresIn,
	}, nil
}
