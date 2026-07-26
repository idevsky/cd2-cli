//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"
)

func TestIntegrationSystem_GetRuntimeInfo(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	info, err := c.System().GetRuntimeInfo(ctx)
	if err != nil {
		t.Fatalf("GetRuntimeInfo failed: %v", err)
	}

	t.Logf("Runtime info: %v", info)
}

func TestIntegrationSystem_GetRunningInfo(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	info, err := c.System().GetRunningInfo(ctx)
	if err != nil {
		t.Fatalf("GetRunningInfo failed: %v", err)
	}

	t.Logf("Running info: %v", info)
}

func TestIntegrationSystem_GetSettings(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	settings, err := c.System().GetSystemSettings(ctx)
	if err != nil {
		t.Fatalf("GetSystemSettings failed: %v", err)
	}

	t.Logf("Settings retrieved: %v", settings)
}

func TestIntegrationSystem_GetCapabilities(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	caps, err := c.System().GetServiceCapabilities(ctx)
	if err != nil {
		t.Fatalf("GetServiceCapabilities failed: %v", err)
	}

	t.Logf("Capabilities: %v", caps)
}

func TestIntegrationSystem_GetMachineID(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	id, err := c.System().GetMachineId(ctx)
	if err != nil {
		t.Fatalf("GetMachineId failed: %v", err)
	}

	t.Logf("Machine ID: %v", id)
}

func TestIntegrationSystem_GetOnlineDevices(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	handles, err := c.System().GetOpenFileHandles(ctx)
	if err != nil {
		t.Fatalf("GetOpenFileHandles failed: %v", err)
	}

	t.Logf("Open file handles: %v", handles)
}

func TestIntegrationSystem_GetDevices(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	devices, err := c.System().GetOnlineDevices(ctx)
	if err != nil {
		t.Fatalf("GetOnlineDevices failed: %v", err)
	}

	t.Logf("Online devices: %d", len(devices.Devices))
}

func TestIntegrationSystem_GetTempFiles(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	files, err := c.System().GetTempFileTable(ctx)
	if err != nil {
		t.Fatalf("GetTempFileTable failed: %v", err)
	}

	t.Logf("Temp files: %v", files)
}

func TestIntegrationSystem_GetLogs(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logs, err := c.System().ListLogFiles(ctx)
	if err != nil {
		t.Fatalf("ListLogFiles failed: %v", err)
	}

	t.Logf("Log files: %d", len(logs.LogFiles))
}

func TestIntegrationSystem_UpdateCheck(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	info, err := c.System().CheckUpdate(ctx)
	if err != nil {
		t.Fatalf("CheckUpdate failed: %v", err)
	}

	t.Logf("Update info: %v", info)
}

func TestIntegrationSystem_HasUpdate(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := c.System().HasUpdate(ctx)
	if err != nil {
		t.Fatalf("HasUpdate failed: %v", err)
	}

	t.Logf("Has update: %v", result)
}

func TestIntegrationSystem_GetWebServerConfig(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	config, err := c.System().GetWebServerConfig(ctx)
	if err != nil {
		t.Fatalf("GetWebServerConfig failed: %v", err)
	}

	t.Logf("WebServer config: %v", config)
}

func TestIntegrationSystem_GetDirCacheTable(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	table, err := c.System().GetDirCacheTable(ctx)
	if err != nil {
		t.Fatalf("GetDirCacheTable failed: %v", err)
	}

	t.Logf("Dir cache table: %v", table)
}

func TestIntegrationSystem_GetVacuumProgress(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	progress, err := c.System().GetVacuumProgress(ctx)
	if err != nil {
		t.Fatalf("GetVacuumProgress failed: %v", err)
	}

	t.Logf("Vacuum progress: %v", progress)
}
