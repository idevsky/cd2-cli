package whitelist

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/clouddrive/cd2-cli/internal/registry"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestInitDefaultConfig_FromRegistry(t *testing.T) {
	registry.Register(&registry.CommandSpec{
		ID:          "file.list",
		Category:    "file",
		RPC:         "ListFiles",
		Description: "List files",
		RiskLevel:   registry.RiskLow,
	})
	registry.Register(&registry.CommandSpec{
		ID:          "file.delete",
		Category:    "file",
		RPC:         "DeleteFile",
		Description: "Delete file",
		RiskLevel:   registry.RiskHigh,
	})
	registry.Register(&registry.CommandSpec{
		ID:          "system.restart",
		Category:    "system",
		RPC:         "Restart",
		Description: "Restart system",
		RiskLevel:   registry.RiskCritical,
	})
	registry.Register(&registry.CommandSpec{
		ID:          "auth.login",
		Category:    "auth",
		RPC:         "Login",
		Description: "Login",
		RiskLevel:   registry.RiskMedium,
	})

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "whitelist.yaml")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	cmd, ok := m.GetCommand("file.list")
	if !ok {
		t.Fatal("file.list command not found")
	}
	if !cmd.Enabled {
		t.Error("file.list (low risk) should be enabled by default")
	}

	cmd, ok = m.GetCommand("auth.login")
	if !ok {
		t.Fatal("auth.login command not found")
	}
	if !cmd.Enabled {
		t.Error("auth.login (medium risk) should be enabled by default")
	}

	cmd, ok = m.GetCommand("file.delete")
	if !ok {
		t.Fatal("file.delete command not found")
	}
	if cmd.Enabled {
		t.Error("file.delete (high risk) should be disabled by default")
	}

	cmd, ok = m.GetCommand("system.restart")
	if !ok {
		t.Fatal("system.restart command not found")
	}
	if cmd.Enabled {
		t.Error("system.restart (critical risk) should be disabled by default")
	}
}

func TestRegisterCommand_NoDeadlock(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "whitelist.yaml")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	done := make(chan bool)
	go func() {
		err := m.RegisterCommand(&CommandInfo{
			Name:        "test.command",
			Category:    "test",
			Description: "Test command",
			RiskLevel:   RiskLow,
			Enabled:     true,
		})
		if err != nil {
			t.Errorf("RegisterCommand failed: %v", err)
		}
		done <- true
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("RegisterCommand deadlocked")
	}
}

func TestEnableCommand_NoDeadlock(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "whitelist.yaml")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	m.RegisterCommand(&CommandInfo{
		Name:        "test.enable",
		Category:    "test",
		Description: "Test",
		RiskLevel:   RiskLow,
		Enabled:     false,
	})

	done := make(chan bool)
	go func() {
		err := m.EnableCommand("test.enable")
		if err != nil {
			t.Errorf("EnableCommand failed: %v", err)
		}
		done <- true
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("EnableCommand deadlocked")
	}
}

func TestDisableCommand_NoDeadlock(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "whitelist.yaml")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	m.RegisterCommand(&CommandInfo{
		Name:        "test.disable",
		Category:    "test",
		Description: "Test",
		RiskLevel:   RiskLow,
		Enabled:     true,
	})

	done := make(chan bool)
	go func() {
		err := m.DisableCommand("test.disable")
		if err != nil {
			t.Errorf("DisableCommand failed: %v", err)
		}
		done <- true
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("DisableCommand deadlocked")
	}
}

func TestConcurrentAccess(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "whitelist.yaml")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(4)
		go func(i int) {
			defer wg.Done()
			m.RegisterCommand(&CommandInfo{
				Name:        "concurrent." + string(rune('a'+i)),
				Category:    "test",
				Description: "Test",
				RiskLevel:   RiskLow,
				Enabled:     true,
			})
		}(i)
		go func(i int) {
			defer wg.Done()
			m.GetCommand("concurrent." + string(rune('a'+i)))
		}(i)
		go func(i int) {
			defer wg.Done()
			m.ListCommands()
		}(i)
		go func(i int) {
			defer wg.Done()
			m.IsEnabled()
		}(i)
	}
	wg.Wait()
}

func TestIsAllowed_WhenDisabled(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "whitelist.yaml")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	m.SetEnabled(false)

	allowed, _ := m.IsAllowed("unknown.command")
	if !allowed {
		t.Error("When whitelist is disabled, all commands should be allowed")
	}
}

func TestConfigStructure(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "whitelist.yaml")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	cfg := m.GetConfig()
	if cfg.Commands == nil {
		t.Error("Commands map should be initialized")
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "whitelist.yaml")

	m1, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	m1.RegisterCommand(&CommandInfo{
		Name:        "save.test",
		Category:    "test",
		Description: "Test save",
		RiskLevel:   RiskMedium,
		Enabled:     true,
	})
	m1.SetEnabled(true)

	m2, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager for reload failed: %v", err)
	}

	cfg := m2.GetConfig()
	if cfg.Enabled == nil || !*cfg.Enabled {
		t.Error("whitelist_enabled should be true after reload")
	}

	cmd, ok := m2.GetCommand("save.test")
	if !ok {
		t.Fatal("save.test command not found after reload")
	}
	if !cmd.Enabled {
		t.Error("save.test should be enabled after reload")
	}
}

func TestGetCommand_ReturnsCopy(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "whitelist.yaml")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	m.RegisterCommand(&CommandInfo{
		Name:        "copy.test",
		Category:    "test",
		Description: "Original",
		RiskLevel:   RiskLow,
		Enabled:     true,
	})

	cmd1, _ := m.GetCommand("copy.test")
	cmd1.Description = "Modified"

	cmd2, _ := m.GetCommand("copy.test")
	if cmd2.Description == "Modified" {
		t.Error("GetCommand should return a copy, not a reference")
	}
}

func TestListCommands_ReturnsSorted(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "whitelist.yaml")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	m.RegisterCommand(&CommandInfo{Name: "z.cmd", Category: "test", RiskLevel: RiskLow, Enabled: true})
	m.RegisterCommand(&CommandInfo{Name: "a.cmd", Category: "test", RiskLevel: RiskLow, Enabled: true})
	m.RegisterCommand(&CommandInfo{Name: "m.cmd", Category: "test", RiskLevel: RiskLow, Enabled: true})

	cmds := m.ListCommands()
	if len(cmds) < 3 {
		t.Fatalf("expected at least 3 commands, got %d", len(cmds))
	}

	for i := 1; i < len(cmds); i++ {
		if cmds[i-1].Name > cmds[i].Name {
			t.Error("ListCommands should return sorted results")
		}
	}
}

func TestReset_NoDeadlock(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "whitelist.yaml")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	done := make(chan bool)
	go func() {
		err := m.Reset()
		if err != nil {
			t.Errorf("Reset failed: %v", err)
		}
		done <- true
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Reset deadlocked")
	}
}

func TestReset_RebuildsFromRegistry(t *testing.T) {
	registry.Register(&registry.CommandSpec{
		ID:          "reset.test.low",
		Category:    "reset",
		RPC:         "Test",
		Description: "Low risk test",
		RiskLevel:   registry.RiskLow,
	})
	registry.Register(&registry.CommandSpec{
		ID:          "reset.test.high",
		Category:    "reset",
		RPC:         "Test",
		Description: "High risk test",
		RiskLevel:   registry.RiskHigh,
	})

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "whitelist.yaml")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	m.DisableCommand("reset.test.low")
	m.EnableCommand("reset.test.high")

	err = m.Reset()
	if err != nil {
		t.Fatalf("Reset failed: %v", err)
	}

	m2, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager for reload failed: %v", err)
	}

	lowCmd, ok := m2.GetCommand("reset.test.low")
	if !ok {
		t.Fatal("reset.test.low command not found after reset")
	}
	if !lowCmd.Enabled {
		t.Error("reset.test.low (low risk) should be enabled after reset")
	}

	highCmd, ok := m2.GetCommand("reset.test.high")
	if !ok {
		t.Fatal("reset.test.high command not found after reset")
	}
	if highCmd.Enabled {
		t.Error("reset.test.high (high risk) should be disabled after reset")
	}

	if !m2.IsEnabled() {
		t.Error("whitelist should be enabled after reset")
	}
}

func TestRegisterCommand_RejectsWildcard(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "whitelist.yaml")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	err = m.RegisterCommand(&CommandInfo{
		Name:        "file.*",
		Category:    "file",
		Description: "Wildcard test",
		RiskLevel:   RiskLow,
		Enabled:     true,
	})
	if err == nil {
		t.Error("RegisterCommand should reject wildcard patterns")
	}
	if !errors.Is(err, ErrWildcardNotAllowed) {
		t.Errorf("expected ErrWildcardNotAllowed, got: %v", err)
	}
}

func TestEnableCommand_RejectsWildcard(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "whitelist.yaml")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	m.RegisterCommand(&CommandInfo{
		Name:        "test.enable",
		Category:    "test",
		Description: "Test",
		RiskLevel:   RiskLow,
		Enabled:     false,
	})

	err = m.EnableCommand("test.*")
	if err == nil {
		t.Error("EnableCommand should reject wildcard patterns")
	}
	if !errors.Is(err, ErrWildcardNotAllowed) {
		t.Errorf("expected ErrWildcardNotAllowed, got: %v", err)
	}
}

func TestDisableCommand_RejectsWildcard(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "whitelist.yaml")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	m.RegisterCommand(&CommandInfo{
		Name:        "test.disable",
		Category:    "test",
		Description: "Test",
		RiskLevel:   RiskLow,
		Enabled:     true,
	})

	err = m.DisableCommand("test.*")
	if err == nil {
		t.Error("DisableCommand should reject wildcard patterns")
	}
	if !errors.Is(err, ErrWildcardNotAllowed) {
		t.Errorf("expected ErrWildcardNotAllowed, got: %v", err)
	}
}

func TestLoadConfig_RejectsWildcard(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "whitelist.yaml")

	wildcardConfig := `
whitelist_enabled: true
commands:
  "file.*":
    name: "file.*"
    category: file
    description: "Wildcard pattern"
    risk_level: low
    enabled: true
`
	if err := os.WriteFile(configPath, []byte(wildcardConfig), 0600); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	_, err := NewManager(configPath)
	if err == nil {
		t.Error("Load should reject wildcard patterns in config")
	}
	if !errors.Is(err, ErrWildcardNotAllowed) {
		t.Errorf("expected ErrWildcardNotAllowed, got: %v", err)
	}
}

func TestDefaultOpen_OverridesRiskLevel(t *testing.T) {
	registry.Register(&registry.CommandSpec{
		ID:          "test.low.write",
		Category:    "test",
		RPC:         "Test",
		Description: "Low risk write operation",
		RiskLevel:   registry.RiskLow,
		DefaultOpen: registry.BoolPtr(false),
	})
	registry.Register(&registry.CommandSpec{
		ID:          "test.high.read",
		Category:    "test",
		RPC:         "Test",
		Description: "High risk read operation",
		RiskLevel:   registry.RiskHigh,
		DefaultOpen: registry.BoolPtr(true),
	})

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "whitelist.yaml")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	cmd, ok := m.GetCommand("test.low.write")
	if !ok {
		t.Fatal("test.low.write command not found")
	}
	if cmd.Enabled {
		t.Error("test.low.write (low risk, DefaultOpen=false) should be disabled by default")
	}

	cmd, ok = m.GetCommand("test.high.read")
	if !ok {
		t.Fatal("test.high.read command not found")
	}
	if !cmd.Enabled {
		t.Error("test.high.read (high risk, DefaultOpen=true) should be enabled by default")
	}
}

func TestFileMoveBlockedByDefault(t *testing.T) {
	registry.Register(&registry.CommandSpec{
		ID:          "file.move",
		Category:    "file",
		RPC:         "MoveFile",
		Description: "Move files",
		RiskLevel:   registry.RiskMedium,
		DefaultOpen: registry.BoolPtr(false),
	})

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "whitelist.yaml")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	allowed, reason := m.IsAllowed("file.move")
	if allowed {
		t.Error("file.move should be blocked by default whitelist")
	}
	if reason == "" {
		t.Error("reason should be provided when command is blocked")
	}
}

func TestNilDefaultOpen_UsesRiskLevel(t *testing.T) {
	registry.Register(&registry.CommandSpec{
		ID:          "test.nil.low",
		Category:    "test",
		RPC:         "Test",
		Description: "Low risk with nil DefaultOpen",
		RiskLevel:   registry.RiskLow,
	})
	registry.Register(&registry.CommandSpec{
		ID:          "test.nil.high",
		Category:    "test",
		RPC:         "Test",
		Description: "High risk with nil DefaultOpen",
		RiskLevel:   registry.RiskHigh,
	})

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "whitelist.yaml")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	cmd, ok := m.GetCommand("test.nil.low")
	if !ok {
		t.Fatal("test.nil.low command not found")
	}
	if !cmd.Enabled {
		t.Error("test.nil.low (low risk, nil DefaultOpen) should be enabled by default (risk level fallback)")
	}

	cmd, ok = m.GetCommand("test.nil.high")
	if !ok {
		t.Fatal("test.nil.high command not found")
	}
	if cmd.Enabled {
		t.Error("test.nil.high (high risk, nil DefaultOpen) should be disabled by default (risk level fallback)")
	}
}

func TestTaskCancelUploadBlockedByDefault(t *testing.T) {
	registry.Register(&registry.CommandSpec{
		ID:          "task.cancel.upload",
		Category:    "task",
		RPC:         "CancelUploadFiles",
		Description: "Cancel upload tasks",
		RiskLevel:   registry.RiskHigh,
	})

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "whitelist.yaml")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	allowed, reason := m.IsAllowed("task.cancel.upload")
	if allowed {
		t.Error("task.cancel.upload (high risk) should be blocked by default whitelist")
	}
	if reason == "" {
		t.Error("reason should be provided when command is blocked")
	}
}

func TestAliasGroup_ShareWhitelistPolicy(t *testing.T) {
	registry.Register(&registry.CommandSpec{
		ID:          "file.delete",
		Category:    "file",
		RPC:         "DeleteFile",
		Description: "Delete file",
		RiskLevel:   registry.RiskHigh,
		AliasGroup:  "file.delete",
	})
	registry.Register(&registry.CommandSpec{
		ID:          "fs.rm",
		Category:    "fs",
		RPC:         "DeleteFile",
		Description: "Remove file",
		RiskLevel:   registry.RiskHigh,
		AliasGroup:  "file.delete",
	})

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "whitelist.yaml")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	allowed, _ := m.IsAllowed("file.delete")
	allowedAlias, _ := m.IsAllowed("fs.rm")

	if allowed != allowedAlias {
		t.Errorf("file.delete (allowed=%v) and fs.rm (allowed=%v) should have same whitelist policy", allowed, allowedAlias)
	}

	_, ok := m.GetCommand("file.delete")
	_, okAlias := m.GetCommand("fs.rm")

	if !ok {
		t.Error("file.delete should exist in whitelist (canonical)")
	}
	if !okAlias {
		t.Error("fs.rm should resolve to file.delete in whitelist")
	}

	m.DisableCommand("file.delete")

	allowed, _ = m.IsAllowed("file.delete")
	allowedAlias, _ = m.IsAllowed("fs.rm")

	if allowed || allowedAlias {
		t.Error("Both file.delete and fs.rm should be blocked after disabling file.delete")
	}
}

func TestAliasGroup_InitDefaultConfigOnlyCanonical(t *testing.T) {
	registry.Register(&registry.CommandSpec{
		ID:          "file.list",
		Category:    "file",
		RPC:         "GetSubFiles",
		Description: "List files",
		RiskLevel:   registry.RiskLow,
		AliasGroup:  "file.list",
	})
	registry.Register(&registry.CommandSpec{
		ID:          "fs.ls",
		Category:    "fs",
		RPC:         "GetSubFiles",
		Description: "List directory",
		RiskLevel:   registry.RiskLow,
		AliasGroup:  "file.list",
	})

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "whitelist.yaml")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	cfg := m.GetConfig()

	if _, ok := cfg.Commands["file.list"]; !ok {
		t.Error("file.list (canonical) should be in whitelist config")
	}

	if _, ok := cfg.Commands["fs.ls"]; ok {
		t.Error("fs.ls (alias) should NOT be in whitelist config - only canonical IDs")
	}
}

func TestMergeRegistryDefaults_PreservesExistingSettings(t *testing.T) {
	registry.Register(&registry.CommandSpec{
		ID:          "existing.cmd",
		Category:    "test",
		RPC:         "Test",
		Description: "Existing command",
		RiskLevel:   registry.RiskLow,
	})
	registry.Register(&registry.CommandSpec{
		ID:          "existing.disabled",
		Category:    "test",
		RPC:         "Test",
		Description: "Disabled command",
		RiskLevel:   registry.RiskLow,
	})

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "whitelist.yaml")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	m.DisableCommand("existing.cmd")
	m.EnableCommand("existing.disabled")

	m2, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager for reload failed: %v", err)
	}

	cmd, ok := m2.GetCommand("existing.cmd")
	if !ok {
		t.Fatal("existing.cmd not found")
	}
	if cmd.Enabled {
		t.Error("existing.cmd should remain disabled after merge (user setting preserved)")
	}

	cmd, ok = m2.GetCommand("existing.disabled")
	if !ok {
		t.Fatal("existing.disabled not found")
	}
	if !cmd.Enabled {
		t.Error("existing.disabled should remain enabled after merge (user setting preserved)")
	}
}

func TestMergeRegistryDefaults_AddsNewCommands(t *testing.T) {
	registry.Register(&registry.CommandSpec{
		ID:          "new.low",
		Category:    "test",
		RPC:         "Test",
		Description: "New low risk command",
		RiskLevel:   registry.RiskLow,
	})
	registry.Register(&registry.CommandSpec{
		ID:          "new.high",
		Category:    "test",
		RPC:         "Test",
		Description: "New high risk command",
		RiskLevel:   registry.RiskHigh,
	})

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "whitelist.yaml")

	oldConfig := `whitelist_enabled: true
commands:
  existing.old:
    name: existing.old
    category: test
    description: Old command
    risk_level: low
    enabled: true
`
	if err := os.WriteFile(configPath, []byte(oldConfig), 0600); err != nil {
		t.Fatalf("failed to write old config: %v", err)
	}

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	cmd, ok := m.GetCommand("new.low")
	if !ok {
		t.Fatal("new.low should be added from registry")
	}
	if !cmd.Enabled {
		t.Error("new.low (low risk) should be enabled by default")
	}

	cmd, ok = m.GetCommand("new.high")
	if !ok {
		t.Fatal("new.high should be added from registry")
	}
	if cmd.Enabled {
		t.Error("new.high (high risk) should be disabled by default")
	}

	cmd, ok = m.GetCommand("existing.old")
	if !ok {
		t.Fatal("existing.old should still exist")
	}
	if !cmd.Enabled {
		t.Error("existing.old should remain enabled (user setting preserved)")
	}
}

func TestMergeRegistryDefaults_MarksDeprecatedCommands(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "whitelist.yaml")

	oldConfig := `whitelist_enabled: true
commands:
  deleted.cmd:
    name: deleted.cmd
    category: test
    description: Deleted command
    risk_level: low
    enabled: true
`
	if err := os.WriteFile(configPath, []byte(oldConfig), 0600); err != nil {
		t.Fatalf("failed to write old config: %v", err)
	}

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	cmd, ok := m.GetCommand("deleted.cmd")
	if !ok {
		t.Fatal("deleted.cmd should NOT be removed from config")
	}
	if !cmd.Deprecated {
		t.Error("deleted.cmd should be marked deprecated")
	}
}

func TestMergeRegistryDefaults_UpdatesMetadata(t *testing.T) {
	registry.Register(&registry.CommandSpec{
		ID:          "updated.cmd",
		Category:    "updated-category",
		RPC:         "UpdatedRPC",
		Description: "Updated description",
		RiskLevel:   registry.RiskHigh,
	})

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "whitelist.yaml")

	oldConfig := `whitelist_enabled: true
commands:
  updated.cmd:
    name: updated.cmd
    category: old-category
    description: Old description
    risk_level: low
    enabled: true
`
	if err := os.WriteFile(configPath, []byte(oldConfig), 0600); err != nil {
		t.Fatalf("failed to write old config: %v", err)
	}

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	cmd, ok := m.GetCommand("updated.cmd")
	if !ok {
		t.Fatal("updated.cmd not found")
	}
	if cmd.Category != "updated-category" {
		t.Errorf("category should be updated to 'updated-category', got %q", cmd.Category)
	}
	if cmd.Description != "Updated description" {
		t.Errorf("description should be updated, got %q", cmd.Description)
	}
	if cmd.RiskLevel != RiskHigh {
		t.Errorf("risk_level should be updated to 'high', got %q", cmd.RiskLevel)
	}
	if !cmd.Enabled {
		t.Error("enabled state should be preserved (was true)")
	}
	if cmd.Deprecated {
		t.Error("updated.cmd should NOT be deprecated (it exists in registry)")
	}
}

func TestMergeRegistryDefaults_RespectsDefaultOpen(t *testing.T) {
	registry.Register(&registry.CommandSpec{
		ID:          "new.defaultopen.false",
		Category:    "test",
		RPC:         "Test",
		Description: "Low risk but closed by default",
		RiskLevel:   registry.RiskLow,
		DefaultOpen: registry.BoolPtr(false),
	})
	registry.Register(&registry.CommandSpec{
		ID:          "new.defaultopen.true",
		Category:    "test",
		RPC:         "Test",
		Description: "High risk but open by default",
		RiskLevel:   registry.RiskHigh,
		DefaultOpen: registry.BoolPtr(true),
	})

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "whitelist.yaml")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	cmd, ok := m.GetCommand("new.defaultopen.false")
	if !ok {
		t.Fatal("new.defaultopen.false not found")
	}
	if cmd.Enabled {
		t.Error("new.defaultopen.false (DefaultOpen=false) should be disabled")
	}

	cmd, ok = m.GetCommand("new.defaultopen.true")
	if !ok {
		t.Fatal("new.defaultopen.true not found")
	}
	if !cmd.Enabled {
		t.Error("new.defaultopen.true (DefaultOpen=true) should be enabled")
	}
}

func TestStorageAddAndCloudapiLoginSharePolicy(t *testing.T) {
	registry.Register(&registry.CommandSpec{
		ID:          "storage.add.s3",
		Category:    "storage",
		RPC:         "APILoginS3",
		Description: "Add S3 storage",
		RiskLevel:   registry.RiskHigh,
		DefaultOpen: registry.BoolPtr(false),
		AliasGroup:  "storage.add.s3",
	})
	registry.Register(&registry.CommandSpec{
		ID:          "cloudapi.login-s3",
		Category:    "cloudapi",
		RPC:         "APILoginS3",
		Description: "Login to S3",
		RiskLevel:   registry.RiskHigh,
		DefaultOpen: registry.BoolPtr(false),
		AliasGroup:  "storage.add.s3",
	})

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "whitelist.yaml")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	allowedS3Add, _ := m.IsAllowed("storage.add.s3")
	allowedS3Login, _ := m.IsAllowed("cloudapi.login-s3")

	if allowedS3Add || allowedS3Login {
		t.Error("Both storage.add.s3 and cloudapi.login-s3 should be blocked by default (DefaultOpen=false)")
	}

	if allowedS3Add != allowedS3Login {
		t.Errorf("storage.add.s3 (allowed=%v) and cloudapi.login-s3 (allowed=%v) should have same whitelist policy", allowedS3Add, allowedS3Login)
	}

	m.EnableCommand("storage.add.s3")

	allowedS3Add, _ = m.IsAllowed("storage.add.s3")
	allowedS3Login, _ = m.IsAllowed("cloudapi.login-s3")

	if !allowedS3Add || !allowedS3Login {
		t.Error("Both storage.add.s3 and cloudapi.login-s3 should be enabled after enabling storage.add.s3")
	}
}

func TestEnableCommand_ResolvesAliasToCanonical(t *testing.T) {
	registry.Register(&registry.CommandSpec{
		ID:          "storage.add.s3",
		Category:    "storage",
		RPC:         "APILoginS3",
		Description: "Add S3 storage",
		RiskLevel:   registry.RiskHigh,
		DefaultOpen: registry.BoolPtr(false),
		AliasGroup:  "storage.add.s3",
	})
	registry.Register(&registry.CommandSpec{
		ID:          "cloudapi.login-s3",
		Category:    "cloudapi",
		RPC:         "APILoginS3",
		Description: "Login to S3",
		RiskLevel:   registry.RiskHigh,
		DefaultOpen: registry.BoolPtr(false),
		AliasGroup:  "storage.add.s3",
	})

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "whitelist.yaml")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	err = m.EnableCommand("cloudapi.login-s3")
	if err != nil {
		t.Fatalf("EnableCommand failed: %v", err)
	}

	allowedS3Add, _ := m.IsAllowed("storage.add.s3")
	allowedS3Login, _ := m.IsAllowed("cloudapi.login-s3")

	if !allowedS3Add || !allowedS3Login {
		t.Error("Enabling cloudapi.login-s3 (alias) should enable storage.add.s3 (canonical)")
	}
}

func TestDisableCommand_ResolvesCanonicalToAlias(t *testing.T) {
	registry.Register(&registry.CommandSpec{
		ID:          "storage.add.gcs",
		Category:    "storage",
		RPC:         "APILoginGCS",
		Description: "Add GCS storage",
		RiskLevel:   registry.RiskHigh,
		DefaultOpen: registry.BoolPtr(true),
		AliasGroup:  "storage.add.gcs",
	})
	registry.Register(&registry.CommandSpec{
		ID:          "cloudapi.login-gcs",
		Category:    "cloudapi",
		RPC:         "APILoginGCS",
		Description: "Login to GCS",
		RiskLevel:   registry.RiskHigh,
		DefaultOpen: registry.BoolPtr(true),
		AliasGroup:  "storage.add.gcs",
	})

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "whitelist.yaml")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	allowedGCSAdd, _ := m.IsAllowed("storage.add.gcs")
	allowedGCSLogin, _ := m.IsAllowed("cloudapi.login-gcs")
	if !allowedGCSAdd || !allowedGCSLogin {
		t.Fatal("GCS commands should be enabled by default (DefaultOpen=true)")
	}

	err = m.DisableCommand("storage.add.gcs")
	if err != nil {
		t.Fatalf("DisableCommand failed: %v", err)
	}

	allowedGCSAdd, _ = m.IsAllowed("storage.add.gcs")
	allowedGCSLogin, _ = m.IsAllowed("cloudapi.login-gcs")

	if allowedGCSAdd || allowedGCSLogin {
		t.Error("Disabling storage.add.gcs (canonical) should disable cloudapi.login-gcs (alias)")
	}
}

func TestGetCommand_ResolvesAlias(t *testing.T) {
	registry.Register(&registry.CommandSpec{
		ID:          "file.copy",
		Category:    "file",
		RPC:         "CopyFile",
		Description: "Copy file",
		RiskLevel:   registry.RiskMedium,
		AliasGroup:  "file.copy",
	})
	registry.Register(&registry.CommandSpec{
		ID:          "fs.cp",
		Category:    "fs",
		RPC:         "CopyFile",
		Description: "Copy file",
		RiskLevel:   registry.RiskMedium,
		AliasGroup:  "file.copy",
	})

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "whitelist.yaml")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	info, exists := m.GetCommand("fs.cp")
	if !exists {
		t.Fatal("GetCommand should resolve alias fs.cp to canonical file.copy")
	}
	if info.Name != "file.copy" {
		t.Errorf("GetCommand returned name %s, expected file.copy (canonical)", info.Name)
	}
}

func TestDefaultWhitelist_BlocksWriteCommandsAllowsReadCommands(t *testing.T) {
	registry.Register(&registry.CommandSpec{
		ID:          "file.list",
		Category:    "file",
		RPC:         "GetSubFiles",
		Description: "List files",
		RiskLevel:   registry.RiskLow,
	})
	registry.Register(&registry.CommandSpec{
		ID:          "file.find",
		Category:    "file",
		RPC:         "FindFileByPath",
		Description: "Find file",
		RiskLevel:   registry.RiskLow,
	})
	registry.Register(&registry.CommandSpec{
		ID:          "file.mkdir",
		Category:    "file",
		RPC:         "CreateFolder",
		Description: "Create folder",
		RiskLevel:   registry.RiskMedium,
		DefaultOpen: registry.BoolPtr(false),
	})
	registry.Register(&registry.CommandSpec{
		ID:          "file.copy",
		Category:    "file",
		RPC:         "CopyFile",
		Description: "Copy files",
		RiskLevel:   registry.RiskMedium,
		DefaultOpen: registry.BoolPtr(false),
	})
	registry.Register(&registry.CommandSpec{
		ID:          "file.move",
		Category:    "file",
		RPC:         "MoveFile",
		Description: "Move files",
		RiskLevel:   registry.RiskMedium,
		DefaultOpen: registry.BoolPtr(false),
	})
	registry.Register(&registry.CommandSpec{
		ID:          "file.upload",
		Category:    "file",
		RPC:         "UploadLocalFile",
		Description: "Upload file",
		RiskLevel:   registry.RiskHigh,
		DefaultOpen: registry.BoolPtr(false),
	})
	registry.Register(&registry.CommandSpec{
		ID:          "auth.logout",
		Category:    "auth",
		RPC:         "Logout",
		Description: "Logout",
		RiskLevel:   registry.RiskMedium,
		DefaultOpen: registry.BoolPtr(false),
	})
	registry.Register(&registry.CommandSpec{
		ID:          "auth.register",
		Category:    "auth",
		RPC:         "Register",
		Description: "Register account",
		RiskLevel:   registry.RiskMedium,
		DefaultOpen: registry.BoolPtr(false),
	})
	registry.Register(&registry.CommandSpec{
		ID:          "auth.2fa-recovery-codes",
		Category:    "auth",
		RPC:         "GetRecoveryCodes",
		Description: "Get recovery codes",
		RiskLevel:   registry.RiskMedium,
		DefaultOpen: registry.BoolPtr(false),
	})
	registry.Register(&registry.CommandSpec{
		ID:          "file.delete",
		Category:    "file",
		RPC:         "DeleteFile",
		Description: "Delete file",
		RiskLevel:   registry.RiskHigh,
	})

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "whitelist.yaml")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	readCommands := []string{"file.list", "file.find"}
	for _, cmd := range readCommands {
		allowed, _ := m.IsAllowed(cmd)
		if !allowed {
			t.Errorf("read command %s should be allowed by default", cmd)
		}
	}

	writeCommands := []string{
		"file.mkdir",
		"file.copy",
		"file.move",
		"file.upload",
		"file.delete",
		"auth.logout",
		"auth.register",
		"auth.2fa-recovery-codes",
	}
	for _, cmd := range writeCommands {
		allowed, reason := m.IsAllowed(cmd)
		if allowed {
			t.Errorf("write/delete/modify command %s should be blocked by default", cmd)
		}
		if reason == "" {
			t.Errorf("blocked command %s should have a reason", cmd)
		}
	}
}
