package registry

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

func BoolPtr(v bool) *bool {
	return &v
}

type RiskLevel string

const (
	RiskLow      RiskLevel = "low"
	RiskMedium   RiskLevel = "medium"
	RiskHigh     RiskLevel = "high"
	RiskCritical RiskLevel = "critical"
)

type CommandSpec struct {
	ID          string    `json:"id" yaml:"id"`
	Category    string    `json:"category" yaml:"category"`
	RPC         string    `json:"rpc" yaml:"rpc"`
	Description string    `json:"description" yaml:"description"`
	RiskLevel   RiskLevel `json:"risk_level" yaml:"risk_level"`
	DefaultOpen *bool     `json:"default_open" yaml:"default_open"`
	AliasGroup  string    `json:"alias_group" yaml:"alias_group"`
}

type Registry struct {
	mu       sync.RWMutex
	commands map[string]*CommandSpec
}

var globalRegistry = &Registry{
	commands: make(map[string]*CommandSpec),
}

func NewRegistry() *Registry {
	return &Registry{
		commands: make(map[string]*CommandSpec),
	}
}

func (r *Registry) Register(spec *CommandSpec) {
	if spec == nil || spec.ID == "" {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.commands[spec.ID] = spec
}

func (r *Registry) MustRegister(spec *CommandSpec) {
	if spec == nil || spec.ID == "" {
		panic("cannot register nil or empty ID command spec")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.commands[spec.ID]; exists {
		panic(fmt.Sprintf("duplicate command ID '%s' registered", spec.ID))
	}
	r.commands[spec.ID] = spec
}

func MustRegister(spec *CommandSpec) {
	globalRegistry.MustRegister(spec)
}

func (r *Registry) Get(id string) (*CommandSpec, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	spec, exists := r.commands[id]
	if !exists {
		return nil, false
	}
	copy := *spec
	return &copy, true
}

func (r *Registry) List() []*CommandSpec {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*CommandSpec
	for _, spec := range r.commands {
		copy := *spec
		result = append(result, &copy)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	return result
}

func (r *Registry) ListByCategory(category string) []*CommandSpec {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*CommandSpec
	for _, spec := range r.commands {
		if strings.EqualFold(spec.Category, category) {
			copy := *spec
			result = append(result, &copy)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	return result
}

func (r *Registry) ListByRiskLevel(risk RiskLevel) []*CommandSpec {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*CommandSpec
	for _, spec := range r.commands {
		if spec.RiskLevel == risk {
			copy := *spec
			result = append(result, &copy)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	return result
}

func Register(spec *CommandSpec) {
	globalRegistry.Register(spec)
}

func Get(id string) (*CommandSpec, bool) {
	return globalRegistry.Get(id)
}

func List() []*CommandSpec {
	return globalRegistry.List()
}

func ListByCategory(category string) []*CommandSpec {
	return globalRegistry.ListByCategory(category)
}

func ListByRiskLevel(risk RiskLevel) []*CommandSpec {
	return globalRegistry.ListByRiskLevel(risk)
}

func ResolveAliasGroup(commandID string) string {
	spec, exists := globalRegistry.Get(commandID)
	if !exists || spec.AliasGroup == "" {
		return commandID
	}
	return spec.AliasGroup
}
