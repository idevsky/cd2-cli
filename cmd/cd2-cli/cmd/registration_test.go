package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/clouddrive/cd2-cli/internal/registry"
	"github.com/clouddrive/cd2-cli/internal/whitelist"
	"github.com/spf13/cobra"
)

func TestCommandRegistrationCompleteness(t *testing.T) {
	registerCommands()

	cobraCommandIDs := collectCommandIDs(rootCmd)

	registryIDs := make(map[string]bool)
	for _, spec := range registry.List() {
		registryIDs[spec.ID] = true
	}

	var unregistered []string
	for id := range cobraCommandIDs {
		if !registryIDs[id] {
			unregistered = append(unregistered, id)
		}
	}

	if len(unregistered) > 0 {
		sort.Strings(unregistered)
		t.Errorf("Found %d Cobra commands not registered in registry:\n%s",
			len(unregistered),
			formatCommandList(unregistered))
	}
}

func TestAllRegistryCommandsInWhitelist(t *testing.T) {
	registerCommands()

	tmpDir := filepath.Join(os.TempDir(), "cd2-cli-whitelist-test")
	os.MkdirAll(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	whitelistPath := filepath.Join(tmpDir, "whitelist.yaml")
	mgr, err := whitelist.NewManager(whitelistPath)
	if err != nil {
		t.Fatalf("Failed to create whitelist manager: %v", err)
	}

	registryIDs := make(map[string]bool)
	for _, spec := range registry.List() {
		canonicalID := registry.ResolveAliasGroup(spec.ID)
		registryIDs[canonicalID] = true
	}

	whitelistCommands := mgr.ListCommands()
	whitelistIDs := make(map[string]bool)
	for _, cmd := range whitelistCommands {
		whitelistIDs[cmd.Name] = true
	}

	var missingFromWhitelist []string
	for id := range registryIDs {
		if !whitelistIDs[id] {
			missingFromWhitelist = append(missingFromWhitelist, id)
		}
	}

	if len(missingFromWhitelist) > 0 {
		sort.Strings(missingFromWhitelist)
		t.Errorf("Found %d registry commands missing from whitelist:\n%s",
			len(missingFromWhitelist),
			formatCommandList(missingFromWhitelist))
	}
}

func collectCommandIDs(cmd *cobra.Command) map[string]bool {
	result := make(map[string]bool)
	walkCommands(cmd, result)
	return result
}

func walkCommands(cmd *cobra.Command, ids map[string]bool) {
	if cmd.Annotations != nil {
		if _, isLocal := cmd.Annotations[annotationIsLocal]; !isLocal {
			if id, hasID := cmd.Annotations[annotationCommandID]; hasID && id != "" {
				ids[id] = true
			}
		}
	}

	for _, child := range cmd.Commands() {
		walkCommands(child, ids)
	}
}

func formatCommandList(commands []string) string {
	var lines []string
	for _, cmd := range commands {
		lines = append(lines, fmt.Sprintf("  - %s", cmd))
	}
	return strings.Join(lines, "\n")
}

func TestRegistryHasNoDuplicateIDs(t *testing.T) {
	registerCommands()

	ids := make(map[string]int)
	for _, spec := range registry.List() {
		ids[spec.ID]++
	}

	var duplicates []string
	for id, count := range ids {
		if count > 1 {
			duplicates = append(duplicates, fmt.Sprintf("%s (count: %d)", id, count))
		}
	}

	if len(duplicates) > 0 {
		sort.Strings(duplicates)
		t.Errorf("Found %d duplicate command IDs in registry:\n%s",
			len(duplicates),
			formatCommandList(duplicates))
	}
}

func TestCobraHasNoDuplicateCommandIDs(t *testing.T) {
	registerCommands()

	ids := make(map[string]int)
	walkAndCountIDs(rootCmd, ids)

	var duplicates []string
	for id, count := range ids {
		if count > 1 {
			duplicates = append(duplicates, fmt.Sprintf("%s (count: %d)", id, count))
		}
	}

	if len(duplicates) > 0 {
		sort.Strings(duplicates)
		t.Errorf("Found %d duplicate command IDs in Cobra commands:\n%s",
			len(duplicates),
			formatCommandList(duplicates))
	}
}

func walkAndCountIDs(cmd *cobra.Command, ids map[string]int) {
	if cmd.Annotations != nil {
		if _, isLocal := cmd.Annotations[annotationIsLocal]; !isLocal {
			if id, hasID := cmd.Annotations[annotationCommandID]; hasID && id != "" {
				ids[id]++
			}
		}
	}

	for _, child := range cmd.Commands() {
		walkAndCountIDs(child, ids)
	}
}

func TestLocalCommandsAreMarked(t *testing.T) {
	registerCommands()

	localCommands := []string{
		"whitelist.list",
		"whitelist.status",
		"whitelist.enable",
		"whitelist.disable",
		"whitelist.reset",
		"whitelist.path",
		"auth.token-set",
		"auth.token-clear",
		"auth.token-show",
		"local.list",
		"local.mkdir",
	}

	cobraIDs := collectCommandIDs(rootCmd)

	for _, localID := range localCommands {
		if cobraIDs[localID] {
			t.Errorf("Local command %s should be marked with annotationIsLocal and not collected as remote API command", localID)
		}
	}
}

func TestAllNonLocalCommandsHaveIDs(t *testing.T) {
	registerCommands()

	commandsMissingID := collectCommandsMissingID(rootCmd)

	if len(commandsMissingID) > 0 {
		sort.Strings(commandsMissingID)
		t.Errorf("Found %d non-local commands with Run/RunE missing command IDs:\n%s",
			len(commandsMissingID),
			formatCommandList(commandsMissingID))
	}
}

func TestCobraRemoteCommandsHaveRegistryEntries(t *testing.T) {
	registerCommands()

	cobraCommandIDs := collectCommandIDs(rootCmd)

	registryIDs := make(map[string]bool)
	for _, spec := range registry.List() {
		registryIDs[spec.ID] = true
	}

	var missingFromRegistry []string
	for id := range cobraCommandIDs {
		if !registryIDs[id] {
			missingFromRegistry = append(missingFromRegistry, id)
		}
	}

	if len(missingFromRegistry) > 0 {
		sort.Strings(missingFromRegistry)
		t.Errorf("Found %d Cobra remote commands missing from registry:\n%s",
			len(missingFromRegistry),
			formatCommandList(missingFromRegistry))
	}
}

func containsString(list []string, target string) bool {
	for _, s := range list {
		if s == target {
			return true
		}
	}
	return false
}

func TestTaskHelpNoDuplicates(t *testing.T) {
	registerCommands()

	taskCmd := getCommandByPath(rootCmd, "task")
	if taskCmd == nil {
		t.Fatal("task command not found")
	}

	commands := taskCmd.Commands()
	commandNames := make(map[string]int)
	for _, cmd := range commands {
		commandNames[cmd.Name()]++
	}

	var duplicates []string
	for name, count := range commandNames {
		if count > 1 {
			duplicates = append(duplicates, fmt.Sprintf("%s (count: %d)", name, count))
		}
	}

	if len(duplicates) > 0 {
		sort.Strings(duplicates)
		t.Errorf("Found %d duplicate command names under 'task':\n%s",
			len(duplicates),
			formatCommandList(duplicates))
	}
}

func getCommandByPath(root *cobra.Command, path string) *cobra.Command {
	parts := strings.Split(path, " ")
	current := root
	for _, part := range parts {
		for _, cmd := range current.Commands() {
			if cmd.Name() == part {
				current = cmd
				break
			}
		}
	}
	if current.Name() == parts[len(parts)-1] {
		return current
	}
	return nil
}

func collectCommandsMissingID(cmd *cobra.Command) []string {
	var result []string

	if cmd.Annotations != nil {
		if _, isLocal := cmd.Annotations[annotationIsLocal]; !isLocal {
			if _, skipWhitelist := cmd.Annotations[annotationSkipWhitelist]; !skipWhitelist {
				if cmd.Run != nil || cmd.RunE != nil {
					if _, hasID := cmd.Annotations[annotationCommandID]; !hasID {
						result = append(result, cmd.CommandPath())
					}
				}
			}
		}
	} else if cmd.Run != nil || cmd.RunE != nil {
		result = append(result, cmd.CommandPath())
	}

	for _, child := range cmd.Commands() {
		result = append(result, collectCommandsMissingID(child)...)
	}

	return result
}
