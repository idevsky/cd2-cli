//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

func TestIntegrationTask_GetUploadFileCount(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := c.Transfer().GetUploadFileCount(ctx)
	if err != nil {
		t.Fatalf("GetUploadFileCount failed: %v", err)
	}

	t.Logf("Upload files: %v", result)
}

func TestIntegrationTask_GetDownloadFileCount(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := c.Transfer().GetDownloadFileCount(ctx)
	if err != nil {
		t.Fatalf("GetDownloadFileCount failed: %v", err)
	}

	t.Logf("Download files: %v", result)
}

func TestIntegrationTask_GetCopyTasks(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tasks, err := c.Copy().GetCopyTasks(ctx)
	if err != nil {
		t.Fatalf("GetCopyTasks failed: %v", err)
	}

	t.Logf("Copy tasks: %d", len(tasks.CopyTasks))
	for i, task := range tasks.CopyTasks {
		if i >= 5 {
			break
		}
		t.Logf("  - Source: %s, Dest: %s, Status: %v", task.SourcePath, task.DestPath, task.Status)
	}
}

func TestIntegrationTask_GetAllTasksCount(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := c.Transfer().GetAllTasksCount(ctx)
	if err != nil {
		t.Fatalf("GetAllTasksCount failed: %v", err)
	}

	t.Logf("All tasks count: %v", result)
}

func TestIntegrationTask_GetMergeTasks(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tasks, err := c.Copy().GetMergeTasks(ctx)
	if err != nil {
		t.Fatalf("GetMergeTasks failed: %v", err)
	}

	t.Logf("Merge tasks: %d", len(tasks.MergeTasks))
	for i, task := range tasks.MergeTasks {
		if i >= 5 {
			break
		}
		t.Logf("  - Source: %s, Dest: %s, Status: %v", task.SourcePath, task.DestPath, task.Status)
	}
}

func TestIntegrationTask_GetUploadFileList(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tasks, err := c.Transfer().GetUploadFileList(ctx, &pb.GetUploadFileListRequest{})
	if err != nil {
		t.Fatalf("GetUploadFileList failed: %v", err)
	}

	t.Logf("Upload file list: %d items", len(tasks.UploadFiles))
	for i, file := range tasks.UploadFiles {
		if i >= 5 {
			break
		}
		progress := uint64(0)
		if file.Size > 0 {
			progress = (file.TransferedBytes * 100) / file.Size
		}
		t.Logf("  - Dest: %s, Size: %d, Progress: %d%%, Status: %s", file.DestPath, file.Size, progress, file.Status)
	}
}

func TestIntegrationTask_GetDownloadFileList(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tasks, err := c.Transfer().GetDownloadFileList(ctx)
	if err != nil {
		t.Fatalf("GetDownloadFileList failed: %v", err)
	}

	t.Logf("Download file list: %d items, global speed: %.2f bytes/sec", len(tasks.DownloadFiles), tasks.GlobalBytesPerSecond)
	for i, file := range tasks.DownloadFiles {
		if i >= 5 {
			break
		}
		t.Logf("  - Path: %s, Size: %d, Threads: %d, Speed: %.2f bytes/sec", file.FilePath, file.FileLength, file.DownloadThreadCount, file.BytesPerSecond)
	}
}

func TestIntegrationTask_CancelCopyTask(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	copyTasks, err := c.Copy().GetCopyTasks(ctx)
	if err != nil {
		t.Fatalf("GetCopyTasks failed: %v", err)
	}

	if len(copyTasks.CopyTasks) == 0 {
		t.Log("No copy tasks available to cancel - skipping test")
		t.Skip("No copy tasks available to cancel")
	}

	task := copyTasks.CopyTasks[0]
	err = c.Copy().CancelCopyTask(ctx, &pb.CopyTaskRequest{
		SourcePath: task.SourcePath,
		DestPath:   task.DestPath,
	})
	if err != nil {
		t.Errorf("CancelCopyTask failed: %v", err)
	} else {
		t.Logf("Cancelled copy task: %s -> %s", task.SourcePath, task.DestPath)
	}
}

func TestIntegrationTask_CancelMergeTask(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mergeTasks, err := c.Copy().GetMergeTasks(ctx)
	if err != nil {
		t.Fatalf("GetMergeTasks failed: %v", err)
	}

	if len(mergeTasks.MergeTasks) == 0 {
		t.Log("No merge tasks available to cancel - skipping test")
		t.Skip("No merge tasks available to cancel")
	}

	task := mergeTasks.MergeTasks[0]
	err = c.Copy().CancelMergeTask(ctx, &pb.CancelMergeTaskRequest{
		SourcePath: task.SourcePath,
		DestPath:   task.DestPath,
	})
	if err != nil {
		t.Errorf("CancelMergeTask failed: %v", err)
	} else {
		t.Logf("Cancelled merge task: %s -> %s", task.SourcePath, task.DestPath)
	}
}

func TestIntegrationTask_PauseCopyTask(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	copyTasks, err := c.Copy().GetCopyTasks(ctx)
	if err != nil {
		t.Fatalf("GetCopyTasks failed: %v", err)
	}

	if len(copyTasks.CopyTasks) == 0 {
		t.Log("No copy tasks available to pause - skipping test")
		t.Skip("No copy tasks available to pause")
	}

	task := copyTasks.CopyTasks[0]
	err = c.Copy().PauseCopyTask(ctx, &pb.PauseCopyTaskRequest{
		SourcePath: task.SourcePath,
		DestPath:   task.DestPath,
		Pause:      true,
	})
	if err != nil {
		t.Errorf("PauseCopyTask failed: %v", err)
	} else {
		t.Logf("Paused copy task: %s -> %s", task.SourcePath, task.DestPath)
	}
}

func TestIntegrationTask_RestartCopyTask(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	copyTasks, err := c.Copy().GetCopyTasks(ctx)
	if err != nil {
		t.Fatalf("GetCopyTasks failed: %v", err)
	}

	if len(copyTasks.CopyTasks) == 0 {
		t.Log("No copy tasks available to restart - skipping test")
		t.Skip("No copy tasks available to restart")
	}

	task := copyTasks.CopyTasks[0]
	err = c.Copy().RestartCopyTask(ctx, &pb.CopyTaskRequest{
		SourcePath: task.SourcePath,
		DestPath:   task.DestPath,
	})
	if err != nil {
		t.Errorf("RestartCopyTask failed: %v", err)
	} else {
		t.Logf("Restarted copy task: %s -> %s", task.SourcePath, task.DestPath)
	}
}

func TestIntegrationTask_RemoveCompletedCopyTasks(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := c.Copy().RemoveCompletedCopyTasks(ctx)
	if err != nil {
		t.Errorf("RemoveCompletedCopyTasks failed: %v", err)
	} else {
		t.Log("Removed completed copy tasks")
	}
}

func TestIntegrationTask_CancelUploadFiles(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	uploadList, err := c.Transfer().GetUploadFileList(ctx, &pb.GetUploadFileListRequest{})
	if err != nil {
		t.Fatalf("GetUploadFileList failed: %v", err)
	}

	if len(uploadList.UploadFiles) == 0 {
		t.Log("No upload tasks available to cancel - skipping test")
		t.Skip("No upload tasks available to cancel")
	}

	keys := []string{}
	for i, file := range uploadList.UploadFiles {
		if i >= 3 {
			break
		}
		keys = append(keys, file.Key)
	}

	err = c.Transfer().CancelUploadFiles(ctx, keys)
	if err != nil {
		t.Errorf("CancelUploadFiles failed: %v", err)
	} else {
		t.Logf("Cancelled %d upload files", len(keys))
	}
}

func TestIntegrationTask_TaskStatusSummary(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	countResult, err := c.Transfer().GetAllTasksCount(ctx)
	if err != nil {
		t.Fatalf("GetAllTasksCount failed: %v", err)
	}

	t.Logf("Task Status Summary:")
	t.Logf("  Upload count: %d", countResult.UploadCount)
	t.Logf("  Download count: %d", countResult.DownloadCount)

	copyTasks, err := c.Copy().GetCopyTasks(ctx)
	if err != nil {
		t.Errorf("GetCopyTasks failed: %v", err)
	} else {
		t.Logf("  Copy tasks: %d", len(copyTasks.CopyTasks))
	}

	mergeTasks, err := c.Copy().GetMergeTasks(ctx)
	if err != nil {
		t.Errorf("GetMergeTasks failed: %v", err)
	} else {
		t.Logf("  Merge tasks: %d", len(mergeTasks.MergeTasks))
	}
}
