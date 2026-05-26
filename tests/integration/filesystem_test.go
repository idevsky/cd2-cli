//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

func TestIntegrationFS_ListRoot(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	files, err := c.File().GetSubFiles(ctx, &pb.ListSubFileRequest{
		Path:         "/",
		ForceRefresh: false,
	})
	if err != nil {
		t.Fatalf("GetSubFiles failed: %v", err)
	}

	t.Logf("Root directory has %d items", len(files))
	for i, file := range files {
		if i >= 5 {
			break
		}
		t.Logf("  - %s (dir: %v, size: %d)", file.Name, file.IsDirectory, file.Size)
	}
}

func TestIntegrationFS_CreateFolder(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	folderName := "test-integration-folder-" + time.Now().Format("150405")

	result, err := c.File().CreateFolder(ctx, &pb.CreateFolderRequest{
		ParentPath: "/",
		FolderName: folderName,
	})
	if err != nil {
		t.Fatalf("CreateFolder failed: %v", err)
	}

	if !result.Result.Success {
		t.Errorf("CreateFolder not successful: %s", result.Result.ErrorMessage)
	} else {
		t.Logf("Folder created: %s", folderName)
	}
}

func TestIntegrationFS_FindFileByPath(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	file, err := c.File().FindFileByPath(ctx, "/")
	if err != nil {
		t.Fatalf("FindFileByPath failed: %v", err)
	}

	t.Logf("Root directory found: name=%s, path=%s", file.Name, file.FullPathName)
}

func TestIntegrationFS_UploadFile(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	testDir := "/test-upload-" + time.Now().Format("20060102-150405")

	_, err := c.File().CreateFolder(ctx, &pb.CreateFolderRequest{
		ParentPath: "/",
		FolderName: testDir[1:],
	})
	if err != nil {
		t.Logf("CreateFolder error (may already exist): %v", err)
	}

	t.Logf("Upload directory created: %s", testDir)

	createResp, err := c.File().CreateFile(ctx, &pb.CreateFileRequest{
		ParentPath: testDir,
		FileName:   "test-upload.txt",
	})
	if err != nil {
		t.Fatalf("CreateFile failed: %v", err)
	}

	t.Logf("File created, handle: %d", createResp.FileHandle)

	stream, err := c.File().WriteToFileStream(ctx)
	if err != nil {
		t.Fatalf("WriteToFileStream failed: %v", err)
	}

	testContent := []byte("Hello, CloudDrive2 Integration Test!")

	err = stream.Send(&pb.WriteFileRequest{
		FileHandle: createResp.FileHandle,
		StartPos:   0,
		Length:     uint64(len(testContent)),
		Buffer:     testContent,
		CloseFile:  true,
	})
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}

	writeResult, err := stream.CloseAndRecv()
	if err != nil {
		t.Fatalf("CloseAndRecv failed: %v", err)
	}

	t.Logf("File uploaded successfully, bytes written: %d", writeResult.BytesWritten)

	uploadedPath := testDir + "/test-upload.txt"
	_, err = c.File().FindFileByPath(ctx, uploadedPath)
	if err != nil {
		t.Errorf("Uploaded file not found: %v", err)
	} else {
		t.Log("Uploaded file verified")
	}
}

func TestIntegrationFS_DeleteFile(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	folderName := "test-delete-folder-" + time.Now().Format("150405")

	_, err := c.File().CreateFolder(ctx, &pb.CreateFolderRequest{
		ParentPath: "/",
		FolderName: folderName,
	})
	if err != nil {
		t.Fatalf("CreateFolder failed: %v", err)
	}

	result, err := c.File().DeleteFile(ctx, "/"+folderName)
	if err != nil {
		t.Fatalf("DeleteFile failed: %v", err)
	}

	if !result.Success {
		t.Errorf("DeleteFile not successful: %s", result.ErrorMessage)
	} else {
		t.Log("Folder deleted successfully")
	}
}

func TestIntegrationFS_RenameFile(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	oldName := "test-rename-old-" + time.Now().Format("150405")
	newName := "test-rename-new-" + time.Now().Format("150405")

	_, err := c.File().CreateFolder(ctx, &pb.CreateFolderRequest{
		ParentPath: "/",
		FolderName: oldName,
	})
	if err != nil {
		t.Fatalf("CreateFolder failed: %v", err)
	}

	result, err := c.File().RenameFile(ctx, &pb.RenameFileRequest{
		TheFilePath: "/" + oldName,
		NewName:     newName,
	})
	if err != nil {
		t.Fatalf("RenameFile failed: %v", err)
	}

	if !result.Success {
		t.Errorf("RenameFile not successful: %s", result.ErrorMessage)
	} else {
		t.Logf("Folder renamed: %s -> %s", oldName, newName)
	}

	c.File().DeleteFile(ctx, "/"+newName)
}

func TestIntegrationFS_GetSpaceInfo(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	info, err := c.File().GetSpaceInfo(ctx, "/")
	if err != nil {
		t.Fatalf("GetSpaceInfo failed: %v", err)
	}

	t.Logf("Space info: total=%d, used=%d", info.TotalSpace, info.UsedSpace)
}

func TestIntegrationFS_GetMetaData(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	meta, err := c.File().GetMetaData(ctx, "/")
	if err != nil {
		t.Fatalf("GetMetaData failed: %v", err)
	}

	t.Logf("Metadata: %v", meta)
}

func TestIntegrationFS_CopyFile(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	srcName := "test-copy-src-" + time.Now().Format("150405")
	dstDir := "test-copy-dst-" + time.Now().Format("150405")

	_, err := c.File().CreateFolder(ctx, &pb.CreateFolderRequest{
		ParentPath: "/",
		FolderName: srcName,
	})
	if err != nil {
		t.Fatalf("CreateFolder (src) failed: %v", err)
	}

	_, err = c.File().CreateFolder(ctx, &pb.CreateFolderRequest{
		ParentPath: "/",
		FolderName: dstDir,
	})
	if err != nil {
		t.Fatalf("CreateFolder (dst) failed: %v", err)
	}

	result, err := c.File().CopyFile(ctx, &pb.CopyFileRequest{
		TheFilePaths: []string{"/" + srcName},
		DestPath:     "/" + dstDir,
	})
	if err != nil {
		t.Fatalf("CopyFile failed: %v", err)
	}

	if !result.Success {
		t.Errorf("CopyFile not successful: %s", result.ErrorMessage)
	} else {
		t.Logf("Copied: /%s -> /%s", srcName, dstDir)
	}

	c.File().DeleteFile(ctx, "/"+dstDir+"/"+srcName)
	c.File().DeleteFile(ctx, "/"+srcName)
	c.File().DeleteFile(ctx, "/"+dstDir)
}

func TestIntegrationFS_MoveFile(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	srcName := "test-move-src-" + time.Now().Format("150405")
	dstDir := "test-move-dst-" + time.Now().Format("150405")

	_, err := c.File().CreateFolder(ctx, &pb.CreateFolderRequest{
		ParentPath: "/",
		FolderName: srcName,
	})
	if err != nil {
		t.Fatalf("CreateFolder (src) failed: %v", err)
	}

	_, err = c.File().CreateFolder(ctx, &pb.CreateFolderRequest{
		ParentPath: "/",
		FolderName: dstDir,
	})
	if err != nil {
		t.Fatalf("CreateFolder (dst) failed: %v", err)
	}

	result, err := c.File().MoveFile(ctx, &pb.MoveFileRequest{
		TheFilePaths: []string{"/" + srcName},
		DestPath:     "/" + dstDir,
	})
	if err != nil {
		t.Fatalf("MoveFile failed: %v", err)
	}

	if !result.Success {
		t.Errorf("MoveFile not successful: %s", result.ErrorMessage)
	} else {
		t.Logf("Moved: /%s -> /%s", srcName, dstDir)
	}

	c.File().DeleteFile(ctx, "/"+dstDir+"/"+srcName)
	c.File().DeleteFile(ctx, "/"+dstDir)
}
