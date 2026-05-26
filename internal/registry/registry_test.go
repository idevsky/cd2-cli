package registry

import (
	"testing"
)

func TestRegistry_RegisterAndGet(t *testing.T) {
	r := NewRegistry()
	spec := &CommandSpec{
		ID:          "file.delete",
		Category:    "file",
		RPC:         "DeleteFile",
		Description: "Delete a file",
		RiskLevel:   RiskHigh,
		DefaultOpen: BoolPtr(false),
	}
	r.Register(spec)

	got, exists := r.Get("file.delete")
	if !exists {
		t.Fatal("expected command to exist")
	}
	if got.ID != spec.ID {
		t.Errorf("expected ID %s, got %s", spec.ID, got.ID)
	}
	if got.Category != spec.Category {
		t.Errorf("expected Category %s, got %s", spec.Category, got.Category)
	}
	if got.RPC != spec.RPC {
		t.Errorf("expected RPC %s, got %s", spec.RPC, got.RPC)
	}
	if got.Description != spec.Description {
		t.Errorf("expected Description %s, got %s", spec.Description, got.Description)
	}
	if got.RiskLevel != spec.RiskLevel {
		t.Errorf("expected RiskLevel %s, got %s", spec.RiskLevel, got.RiskLevel)
	}
	if *got.DefaultOpen != *spec.DefaultOpen {
		t.Errorf("expected DefaultOpen %v, got %v", *spec.DefaultOpen, *got.DefaultOpen)
	}
}

func TestRegistry_GetNonExistent(t *testing.T) {
	r := NewRegistry()
	_, exists := r.Get("nonexistent")
	if exists {
		t.Fatal("expected command not to exist")
	}
}

func TestRegistry_RegisterNil(t *testing.T) {
	r := NewRegistry()
	r.Register(nil)
	if len(r.List()) != 0 {
		t.Fatal("expected no commands after registering nil")
	}
}

func TestRegistry_RegisterEmptyID(t *testing.T) {
	r := NewRegistry()
	r.Register(&CommandSpec{ID: ""})
	if len(r.List()) != 0 {
		t.Fatal("expected no commands after registering empty ID")
	}
}

func TestRegistry_List(t *testing.T) {
	r := NewRegistry()
	r.Register(&CommandSpec{ID: "system.restart", Category: "system", RiskLevel: RiskCritical, DefaultOpen: BoolPtr(false)})
	r.Register(&CommandSpec{ID: "file.delete", Category: "file", RiskLevel: RiskHigh, DefaultOpen: BoolPtr(false)})
	r.Register(&CommandSpec{ID: "auth.login", Category: "auth", RiskLevel: RiskLow, DefaultOpen: BoolPtr(true)})

	list := r.List()
	if len(list) != 3 {
		t.Fatalf("expected 3 commands, got %d", len(list))
	}
	if list[0].ID != "auth.login" {
		t.Errorf("expected first command to be auth.login, got %s", list[0].ID)
	}
	if list[1].ID != "file.delete" {
		t.Errorf("expected second command to be file.delete, got %s", list[1].ID)
	}
	if list[2].ID != "system.restart" {
		t.Errorf("expected third command to be system.restart, got %s", list[2].ID)
	}
}

func TestRegistry_ListByCategory(t *testing.T) {
	r := NewRegistry()
	r.Register(&CommandSpec{ID: "file.delete", Category: "file", RiskLevel: RiskHigh})
	r.Register(&CommandSpec{ID: "file.list", Category: "file", RiskLevel: RiskLow})
	r.Register(&CommandSpec{ID: "system.restart", Category: "system", RiskLevel: RiskCritical})
	r.Register(&CommandSpec{ID: "auth.login", Category: "auth", RiskLevel: RiskLow})

	fileCommands := r.ListByCategory("file")
	if len(fileCommands) != 2 {
		t.Fatalf("expected 2 file commands, got %d", len(fileCommands))
	}
	for _, cmd := range fileCommands {
		if cmd.Category != "file" {
			t.Errorf("expected category file, got %s", cmd.Category)
		}
	}

	systemCommands := r.ListByCategory("SYSTEM")
	if len(systemCommands) != 1 {
		t.Fatalf("expected 1 system command, got %d", len(systemCommands))
	}
	if systemCommands[0].ID != "system.restart" {
		t.Errorf("expected system.restart, got %s", systemCommands[0].ID)
	}
}

func TestRegistry_ListByRiskLevel(t *testing.T) {
	r := NewRegistry()
	r.Register(&CommandSpec{ID: "file.delete", Category: "file", RiskLevel: RiskHigh})
	r.Register(&CommandSpec{ID: "system.restart", Category: "system", RiskLevel: RiskCritical})
	r.Register(&CommandSpec{ID: "file.list", Category: "file", RiskLevel: RiskLow})
	r.Register(&CommandSpec{ID: "auth.login", Category: "auth", RiskLevel: RiskLow})

	highCommands := r.ListByRiskLevel(RiskHigh)
	if len(highCommands) != 1 {
		t.Fatalf("expected 1 high risk command, got %d", len(highCommands))
	}
	if highCommands[0].ID != "file.delete" {
		t.Errorf("expected file.delete, got %s", highCommands[0].ID)
	}

	lowCommands := r.ListByRiskLevel(RiskLow)
	if len(lowCommands) != 2 {
		t.Fatalf("expected 2 low risk commands, got %d", len(lowCommands))
	}

	criticalCommands := r.ListByRiskLevel(RiskCritical)
	if len(criticalCommands) != 1 {
		t.Fatalf("expected 1 critical risk command, got %d", len(criticalCommands))
	}
}

func TestRegistry_GetReturnsCopy(t *testing.T) {
	r := NewRegistry()
	original := &CommandSpec{
		ID:          "test.command",
		Category:    "test",
		RiskLevel:   RiskLow,
		DefaultOpen: BoolPtr(true),
	}
	r.Register(original)

	got1, _ := r.Get("test.command")
	got1.Category = "modified"

	got2, _ := r.Get("test.command")
	if got2.Category == "modified" {
		t.Error("expected Get to return a copy, but it modified the stored command")
	}
}

func TestGlobalRegistry(t *testing.T) {
	ResetGlobalRegistry()

	Register(&CommandSpec{ID: "global.test", Category: "global", RiskLevel: RiskLow})

	spec, exists := Get("global.test")
	if !exists {
		t.Fatal("expected global command to exist")
	}
	if spec.ID != "global.test" {
		t.Errorf("expected ID global.test, got %s", spec.ID)
	}

	ResetGlobalRegistry()
	_, exists = Get("global.test")
	if exists {
		t.Fatal("expected global command to be reset")
	}
}

func ResetGlobalRegistry() {
	globalRegistry = &Registry{
		commands: make(map[string]*CommandSpec),
	}
}

func TestResolveAliasGroup(t *testing.T) {
	ResetGlobalRegistry()

	Register(&CommandSpec{ID: "file.delete", Category: "file", RiskLevel: RiskHigh, AliasGroup: "file.delete"})
	Register(&CommandSpec{ID: "fs.rm", Category: "fs", RiskLevel: RiskHigh, AliasGroup: "file.delete"})
	Register(&CommandSpec{ID: "system.info", Category: "system", RiskLevel: RiskLow})

	tests := []struct {
		input    string
		expected string
	}{
		{"file.delete", "file.delete"},
		{"fs.rm", "file.delete"},
		{"system.info", "system.info"},
		{"unknown.cmd", "unknown.cmd"},
	}

	for _, tt := range tests {
		result := ResolveAliasGroup(tt.input)
		if result != tt.expected {
			t.Errorf("ResolveAliasGroup(%s) = %s, want %s", tt.input, result, tt.expected)
		}
	}
}
