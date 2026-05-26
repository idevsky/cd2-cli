//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

func TestIntegrationFile_GetCloudMemberships(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	memberships, err := c.File().GetCloudMemberships(ctx, "/")
	if err != nil {
		t.Fatalf("GetCloudMemberships failed: %v", err)
	}

	t.Logf("Memberships: %d", len(memberships.Memberships))
}

func TestIntegrationFile_GetOriginalPath(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := c.File().GetOriginalPath(ctx, "/")
	if err != nil {
		t.Fatalf("GetOriginalPath failed: %v", err)
	}

	t.Logf("Original path: %s", result.Result)
}

func TestIntegrationFile_GetDownloadUrl(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	url, err := c.File().GetDownloadUrl(ctx, &pb.GetDownloadUrlPathRequest{
		Path:         "/",
		Preview:      false,
		GetDirectUrl: false,
	})
	if err != nil {
		t.Fatalf("GetDownloadUrl failed: %v", err)
	}

	t.Logf("Download URL path: %s", url.DownloadUrlPath)
}

func TestIntegrationFile_CreateEncryptedFolder(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	folderName := "test-encrypted-" + time.Now().Format("150405")

	result, err := c.File().CreateEncryptedFolder(ctx, &pb.CreateEncryptedFolderRequest{
		ParentPath:   "/",
		FolderName:   folderName,
		Password:     "testpassword123",
		SavePassword: false,
	})
	if err != nil {
		t.Fatalf("CreateEncryptedFolder failed: %v", err)
	}

	if result.Result.Success {
		t.Logf("Encrypted folder created: %s", result.FolderCreated.FullPathName)
		// Clean up
		c.File().DeleteFile(ctx, result.FolderCreated.FullPathName)
	}
}

func TestIntegrationFile_DeleteFilesBatch(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create multiple folders
	folder1 := "test-batch-del-1-" + time.Now().Format("150405")
	folder2 := "test-batch-del-2-" + time.Now().Format("150405")

	_, err := c.File().CreateFolder(ctx, &pb.CreateFolderRequest{
		ParentPath: "/",
		FolderName: folder1,
	})
	if err != nil {
		t.Fatalf("CreateFolder 1 failed: %v", err)
	}

	_, err = c.File().CreateFolder(ctx, &pb.CreateFolderRequest{
		ParentPath: "/",
		FolderName: folder2,
	})
	if err != nil {
		t.Fatalf("CreateFolder 2 failed: %v", err)
	}

	// Batch delete
	result, err := c.File().DeleteFiles(ctx, []string{"/" + folder1, "/" + folder2})
	if err != nil {
		t.Fatalf("DeleteFiles failed: %v", err)
	}

	if !result.Success {
		t.Errorf("DeleteFiles not successful: %s", result.ErrorMessage)
	} else {
		t.Log("Batch delete successful")
	}
}

func TestIntegrationFile_RenameFilesBatch(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a test folder
	folderName := "test-rename-batch-" + time.Now().Format("150405")
	newFolderName := "test-rename-batch-new-" + time.Now().Format("150405")

	_, err := c.File().CreateFolder(ctx, &pb.CreateFolderRequest{
		ParentPath: "/",
		FolderName: folderName,
	})
	if err != nil {
		t.Fatalf("CreateFolder failed: %v", err)
	}

	// Batch rename (single file)
	result, err := c.File().RenameFiles(ctx, &pb.RenameFilesRequest{
		RenameFiles: []*pb.RenameFileRequest{
			{
				TheFilePath: "/" + folderName,
				NewName:     newFolderName,
			},
		},
	})
	if err != nil {
		t.Fatalf("RenameFiles failed: %v", err)
	}

	if !result.Success {
		t.Errorf("RenameFiles not successful: %s", result.ErrorMessage)
	} else {
		t.Logf("Batch rename successful")
	}

	// Clean up
	c.File().DeleteFile(ctx, "/"+newFolderName)
}

func TestIntegrationFile_LockUnlockEncrypted(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create encrypted folder
	folderName := "test-lock-unlock-" + time.Now().Format("150405")

	result, err := c.File().CreateEncryptedFolder(ctx, &pb.CreateEncryptedFolderRequest{
		ParentPath:   "/",
		FolderName:   folderName,
		Password:     "testpassword123",
		SavePassword: false,
	})
	if err != nil {
		t.Fatalf("CreateEncryptedFolder failed: %v", err)
	}

	folderPath := result.FolderCreated.FullPathName

	// Lock the folder
	lockResult, err := c.File().LockEncryptedFile(ctx, folderPath)
	if err != nil {
		t.Logf("LockEncryptedFile failed (may be expected): %v", err)
	} else {
		t.Logf("Lock result: success=%v", lockResult.Success)
	}

	// Unlock the folder
	unlockResult, err := c.File().UnlockEncryptedFile(ctx, &pb.UnlockEncryptedFileRequest{
		Path:            folderPath,
		Password:        "testpassword123",
		PermanentUnlock: false,
	})
	if err != nil {
		t.Logf("UnlockEncryptedFile failed (may be expected): %v", err)
	} else {
		t.Logf("Unlock result: success=%v", unlockResult.Success)
	}

	// Clean up
	c.File().DeleteFile(ctx, folderPath)
}

func TestIntegrationFile_CreateFileAndClose(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	testDir := "/test-file-create-" + time.Now().Format("150405")

	// Create test directory
	_, err := c.File().CreateFolder(ctx, &pb.CreateFolderRequest{
		ParentPath: "/",
		FolderName: testDir[1:],
	})
	if err != nil {
		t.Fatalf("CreateFolder failed: %v", err)
	}

	// Create file
	createResult, err := c.File().CreateFile(ctx, &pb.CreateFileRequest{
		ParentPath: testDir,
		FileName:   "test-file.txt",
	})
	if err != nil {
		t.Fatalf("CreateFile failed: %v", err)
	}

	t.Logf("File created, handle: %d", createResult.FileHandle)

	// Close file
	closeResult, err := c.File().CloseFile(ctx, &pb.CloseFileRequest{
		FileHandle: createResult.FileHandle,
	})
	if err != nil {
		t.Fatalf("CloseFile failed: %v", err)
	}

	t.Logf("Close result: success=%v", closeResult.Success)

	// Clean up
	c.File().DeleteFile(ctx, testDir)
}

func TestIntegrationFile_WriteToFile(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	testDir := "/test-file-write-" + time.Now().Format("150405")

	// Create test directory
	_, err := c.File().CreateFolder(ctx, &pb.CreateFolderRequest{
		ParentPath: "/",
		FolderName: testDir[1:],
	})
	if err != nil {
		t.Fatalf("CreateFolder failed: %v", err)
	}

	// Create file
	createResult, err := c.File().CreateFile(ctx, &pb.CreateFileRequest{
		ParentPath: testDir,
		FileName:   "test-write.txt",
	})
	if err != nil {
		t.Fatalf("CreateFile failed: %v", err)
	}

	testContent := []byte("Hello, CloudDrive2!")

	// Write to file
	writeResult, err := c.File().WriteToFile(ctx, &pb.WriteFileRequest{
		FileHandle: createResult.FileHandle,
		StartPos:   0,
		Length:     uint64(len(testContent)),
		Buffer:     testContent,
		CloseFile:  true,
	})
	if err != nil {
		t.Fatalf("WriteToFile failed: %v", err)
	}

	t.Logf("Write result: bytesWritten=%d", writeResult.BytesWritten)

	// Verify file exists
	file, err := c.File().FindFileByPath(ctx, testDir+"/test-write.txt")
	if err != nil {
		t.Errorf("FindFileByPath failed: %v", err)
	} else {
		t.Logf("File verified: %s, size=%d", file.Name, file.Size)
	}

	// Clean up
	c.File().DeleteFile(ctx, testDir)
}
