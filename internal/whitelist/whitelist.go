package whitelist

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/clouddrive/cd2-cli/internal/registry"

	"gopkg.in/yaml.v3"
)

var (
	ErrWildcardNotAllowed = errors.New("wildcard patterns are not allowed in whitelist; use explicit command names")
)

// RiskLevel defines the risk classification of an API
type RiskLevel string

const (
	RiskLow      RiskLevel = "low"
	RiskMedium   RiskLevel = "medium"
	RiskHigh     RiskLevel = "high"
	RiskCritical RiskLevel = "critical"
)

func boolPtr(v bool) *bool {
	return &v
}

func containsWildcard(name string) bool {
	return strings.Contains(name, "*")
}

func validateCommandName(name string) error {
	if containsWildcard(name) {
		return fmt.Errorf("invalid command name '%s': %w", name, ErrWildcardNotAllowed)
	}
	return nil
}

// CommandInfo holds metadata about a CLI command
type CommandInfo struct {
	Name        string    `yaml:"name" json:"name"`
	Category    string    `yaml:"category" json:"category"`
	Description string    `yaml:"description" json:"description"`
	RiskLevel   RiskLevel `yaml:"risk_level" json:"risk_level"`
	Enabled     bool      `yaml:"enabled" json:"enabled"`
	Deprecated  bool      `yaml:"deprecated" json:"deprecated"`
}

// Config holds the whitelist configuration
type Config struct {
	// Enabled controls whether whitelist checking is active
	// If nil (not set in config), defaults to true for security
	Enabled *bool `yaml:"whitelist_enabled" json:"whitelist_enabled"`
	// Commands maps command names to their configuration
	Commands map[string]*CommandInfo `yaml:"commands" json:"commands"`
	// DefaultRisk is the default risk level for unknown commands
	DefaultRisk RiskLevel `yaml:"default_risk" json:"default_risk"`
}

// Manager handles whitelist operations
type Manager struct {
	config     *Config
	configPath string
	mu         sync.RWMutex
}

// NewManager creates a new whitelist manager
func NewManager(configPath string) (*Manager, error) {
	m := &Manager{
		configPath: configPath,
		config: &Config{
			Enabled:     boolPtr(true),
			Commands:    make(map[string]*CommandInfo),
			DefaultRisk: RiskLow,
		},
	}

	err := m.Load()
	if err != nil {
		if os.IsNotExist(err) {
			m.initDefaultConfig()
			if err := m.Save(); err != nil {
				return nil, fmt.Errorf("failed to save default config: %w", err)
			}
			return m, nil
		}
		return nil, err
	}

	if err := m.MergeRegistryDefaults(); err != nil {
		return nil, fmt.Errorf("failed to merge registry defaults: %w", err)
	}

	return m, nil
}

// initDefaultConfig initializes the default configuration with all commands
// from the registry. Low and medium risk commands are enabled by default,
// high and critical risk commands are disabled by default.
func (m *Manager) initDefaultConfig() {
	m.config.Commands = make(map[string]*CommandInfo)

	registeredGroups := make(map[string]bool)

	for _, spec := range registry.List() {
		canonicalID := spec.ID
		if spec.AliasGroup != "" {
			canonicalID = spec.AliasGroup
		}

		if registeredGroups[canonicalID] {
			continue
		}
		registeredGroups[canonicalID] = true

		enabled := false
		if spec.DefaultOpen != nil {
			enabled = *spec.DefaultOpen
		} else {
			enabled = spec.RiskLevel == registry.RiskLow || spec.RiskLevel == registry.RiskMedium
		}
		m.config.Commands[canonicalID] = &CommandInfo{
			Name:        canonicalID,
			Category:    spec.Category,
			Description: spec.Description,
			RiskLevel:   RiskLevel(spec.RiskLevel),
			Enabled:     enabled,
		}
	}
}

// Load reads the whitelist configuration from disk
func (m *Manager) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse whitelist config: %w", err)
	}

	if config.Commands == nil {
		config.Commands = make(map[string]*CommandInfo)
	}

	for name := range config.Commands {
		if err := validateCommandName(name); err != nil {
			return fmt.Errorf("invalid command in whitelist config: %w", err)
		}
	}

	if config.Enabled == nil {
		config.Enabled = boolPtr(true)
	}

	m.config = &config

	return nil
}

// Save writes the whitelist configuration to disk
// This method acquires a read lock and should be called when no lock is held.
func (m *Manager) Save() error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.saveUnlocked()
}

// saveUnlocked writes the whitelist configuration to disk.
// Must be called with lock already held (either read or write).
func (m *Manager) saveUnlocked() error {
	dir := filepath.Dir(m.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(m.config)
	if err != nil {
		return fmt.Errorf("failed to marshal whitelist config: %w", err)
	}

	if err := os.WriteFile(m.configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write whitelist config: %w", err)
	}

	return nil
}

// IsEnabled returns whether whitelist checking is enabled
func (m *Manager) IsEnabled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.config.Enabled == nil {
		return true
	}
	return *m.config.Enabled
}

// SetEnabled enables or disables whitelist checking
func (m *Manager) SetEnabled(enabled bool) error {
	m.mu.Lock()
	m.config.Enabled = boolPtr(enabled)
	m.mu.Unlock()
	return m.Save()
}

// IsAllowed checks if a command is allowed to execute
func (m *Manager) IsAllowed(commandName string) (bool, string) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	enabled := true
	if m.config.Enabled != nil {
		enabled = *m.config.Enabled
	}

	if !enabled {
		return true, ""
	}

	canonicalID := registry.ResolveAliasGroup(commandName)

	info, exists := m.config.Commands[canonicalID]
	if !exists {
		return false, fmt.Sprintf("command '%s' is not registered in whitelist", canonicalID)
	}

	if !info.Enabled {
		return false, fmt.Sprintf("command '%s' is disabled (risk level: %s)", canonicalID, info.RiskLevel)
	}

	return true, ""
}

// RegisterCommand registers a command in the whitelist
func (m *Manager) RegisterCommand(info *CommandInfo) error {
	if err := validateCommandName(info.Name); err != nil {
		return err
	}
	m.mu.Lock()
	m.config.Commands[info.Name] = info
	m.mu.Unlock()
	return m.Save()
}

// EnableCommand enables a command
func (m *Manager) EnableCommand(commandName string) error {
	if err := validateCommandName(commandName); err != nil {
		return err
	}
	canonicalID := registry.ResolveAliasGroup(commandName)
	m.mu.Lock()
	info, exists := m.config.Commands[canonicalID]
	if !exists {
		m.mu.Unlock()
		return fmt.Errorf("command '%s' not found in whitelist", canonicalID)
	}
	info.Enabled = true
	m.mu.Unlock()
	return m.Save()
}

// DisableCommand disables a command
func (m *Manager) DisableCommand(commandName string) error {
	if err := validateCommandName(commandName); err != nil {
		return err
	}
	canonicalID := registry.ResolveAliasGroup(commandName)
	m.mu.Lock()
	info, exists := m.config.Commands[canonicalID]
	if !exists {
		m.mu.Unlock()
		return fmt.Errorf("command '%s' not found in whitelist", canonicalID)
	}
	info.Enabled = false
	m.mu.Unlock()
	return m.Save()
}

// RemoveCommand removes a command from the whitelist
func (m *Manager) RemoveCommand(commandName string) error {
	m.mu.Lock()
	delete(m.config.Commands, commandName)
	m.mu.Unlock()
	return m.Save()
}

// GetCommand returns command info
func (m *Manager) GetCommand(commandName string) (*CommandInfo, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	canonicalID := registry.ResolveAliasGroup(commandName)

	info, exists := m.config.Commands[canonicalID]
	if !exists {
		return nil, false
	}

	copy := *info
	return &copy, true
}

// ListCommands returns all registered commands
func (m *Manager) ListCommands() []*CommandInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var commands []*CommandInfo
	for _, info := range m.config.Commands {
		copy := *info
		commands = append(commands, &copy)
	}

	// Sort by name
	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Name < commands[j].Name
	})

	return commands
}

// ListCommandsByCategory returns commands filtered by category
func (m *Manager) ListCommandsByCategory(category string) []*CommandInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var commands []*CommandInfo
	for _, info := range m.config.Commands {
		if strings.EqualFold(info.Category, category) {
			copy := *info
			commands = append(commands, &copy)
		}
	}

	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Name < commands[j].Name
	})

	return commands
}

// ListCommandsByRisk returns commands filtered by risk level
func (m *Manager) ListCommandsByRisk(risk RiskLevel) []*CommandInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var commands []*CommandInfo
	for _, info := range m.config.Commands {
		if info.RiskLevel == risk {
			copy := *info
			commands = append(commands, &copy)
		}
	}

	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Name < commands[j].Name
	})

	return commands
}

// GetConfigPath returns the configuration file path
func (m *Manager) GetConfigPath() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.configPath
}

// GetConfig returns a copy of the current config
func (m *Manager) GetConfig() *Config {
	m.mu.RLock()
	defer m.mu.RUnlock()

	configCopy := &Config{
		Enabled:     m.config.Enabled,
		Commands:    make(map[string]*CommandInfo),
		DefaultRisk: m.config.DefaultRisk,
	}

	for k, v := range m.config.Commands {
		copy := *v
		configCopy.Commands[k] = &copy
	}

	return configCopy
}

// Reset rebuilds the whitelist config from the registry defaults.
// This will reinitialize all commands based on their risk levels.
func (m *Manager) Reset() error {
	m.mu.Lock()
	m.config.Enabled = boolPtr(true)
	m.initDefaultConfig()
	m.mu.Unlock()
	return m.Save()
}

// MergeRegistryDefaults merges the current whitelist config with the registry.
// - Existing commands preserve their Enabled state
// - New registry commands are added with default policy
// - Commands not in registry are marked deprecated (not removed)
func (m *Manager) MergeRegistryDefaults() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	registryCmds := registry.List()
	registryMap := make(map[string]*registry.CommandSpec)
	for _, spec := range registryCmds {
		canonicalID := spec.ID
		if spec.AliasGroup != "" {
			canonicalID = spec.AliasGroup
		}
		registryMap[canonicalID] = spec
	}

	for name, info := range m.config.Commands {
		spec, inRegistry := registryMap[name]
		if !inRegistry {
			info.Deprecated = true
			continue
		}
		info.Category = spec.Category
		info.Description = spec.Description
		info.RiskLevel = RiskLevel(spec.RiskLevel)
		info.Deprecated = false
	}

	registeredGroups := make(map[string]bool)
	for _, spec := range registryCmds {
		canonicalID := spec.ID
		if spec.AliasGroup != "" {
			canonicalID = spec.AliasGroup
		}

		if registeredGroups[canonicalID] {
			continue
		}
		registeredGroups[canonicalID] = true

		if _, exists := m.config.Commands[canonicalID]; !exists {
			enabled := false
			if spec.DefaultOpen != nil {
				enabled = *spec.DefaultOpen
			} else {
				enabled = spec.RiskLevel == registry.RiskLow || spec.RiskLevel == registry.RiskMedium
			}
			m.config.Commands[canonicalID] = &CommandInfo{
				Name:        canonicalID,
				Category:    spec.Category,
				Description: spec.Description,
				RiskLevel:   RiskLevel(spec.RiskLevel),
				Enabled:     enabled,
				Deprecated:  false,
			}
		}
	}

	return m.saveUnlocked()
}
