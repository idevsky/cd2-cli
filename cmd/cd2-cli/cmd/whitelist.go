package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/clouddrive/cd2-cli/internal/registry"
	"github.com/clouddrive/cd2-cli/internal/whitelist"
)

var whitelistCmd = &cobra.Command{
	Use:   "whitelist",
	Short: "Whitelist management commands",
}

var whitelistListCmd = &cobra.Command{
	Use:   "list",
	Short: "List whitelisted commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		showAll, _ := cmd.Flags().GetBool("all")
		riskLevel, _ := cmd.Flags().GetString("risk")
		category, _ := cmd.Flags().GetString("category")

		if whitelistMgr == nil {
			return fmt.Errorf("whitelist not initialized")
		}

		var commands []*whitelistCommandOutput
		if riskLevel != "" {
			commands = listByRiskLevel(riskLevel, showAll)
		} else if category != "" {
			commands = listByCategory(category, showAll)
		} else {
			commands = listAllCommands(showAll)
		}

		return outputResult(commands)
	},
}

var whitelistStatusCmd = &cobra.Command{
	Use:   "status [command-id]",
	Short: "Get status of a specific command",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		commandID := args[0]

		if whitelistMgr == nil {
			return fmt.Errorf("whitelist not initialized")
		}

		info, exists := whitelistMgr.GetCommand(commandID)
		if !exists {
			return outputResult(map[string]interface{}{
				"id":           commandID,
				"exists":       false,
				"error":        fmt.Sprintf("command '%s' not found in whitelist", commandID),
				"canonical_id": registry.ResolveAliasGroup(commandID),
			})
		}

		canonicalID := registry.ResolveAliasGroup(commandID)
		return outputResult(map[string]interface{}{
			"requested_id": commandID,
			"canonical_id": canonicalID,
			"id":           info.Name,
			"category":     info.Category,
			"description":  info.Description,
			"risk_level":   info.RiskLevel,
			"enabled":      info.Enabled,
			"exists":       true,
		})
	},
}

var whitelistEnableCmd = &cobra.Command{
	Use:   "enable [command-id]",
	Short: "Enable a command in the whitelist",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		commandID := args[0]

		if whitelistMgr == nil {
			return fmt.Errorf("whitelist not initialized")
		}

		canonicalID := registry.ResolveAliasGroup(commandID)

		err := whitelistMgr.EnableCommand(commandID)
		if err != nil {
			return err
		}

		return outputResult(map[string]interface{}{
			"requested_id": commandID,
			"canonical_id": canonicalID,
			"id":           canonicalID,
			"enabled":      true,
			"success":      true,
		})
	},
}

var whitelistDisableCmd = &cobra.Command{
	Use:   "disable [command-id]",
	Short: "Disable a command in the whitelist",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		commandID := args[0]

		if whitelistMgr == nil {
			return fmt.Errorf("whitelist not initialized")
		}

		canonicalID := registry.ResolveAliasGroup(commandID)

		err := whitelistMgr.DisableCommand(commandID)
		if err != nil {
			return err
		}

		return outputResult(map[string]interface{}{
			"requested_id": commandID,
			"canonical_id": canonicalID,
			"id":           canonicalID,
			"enabled":      false,
			"success":      true,
		})
	},
}

var whitelistResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset whitelist to registry defaults",
	RunE: func(cmd *cobra.Command, args []string) error {
		if whitelistMgr == nil {
			return fmt.Errorf("whitelist not initialized")
		}

		err := whitelistMgr.Reset()
		if err != nil {
			return err
		}

		return outputResult(map[string]interface{}{
			"success":        true,
			"whitelist_path": whitelistMgr.GetConfigPath(),
			"message":        "whitelist reset to registry defaults",
		})
	},
}

var whitelistPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show whitelist configuration file path",
	RunE: func(cmd *cobra.Command, args []string) error {
		if whitelistMgr == nil {
			return fmt.Errorf("whitelist not initialized")
		}

		return outputResult(map[string]string{
			"path":              whitelistMgr.GetConfigPath(),
			"whitelist_enabled": fmt.Sprintf("%v", whitelistMgr.IsEnabled()),
		})
	},
}

type whitelistCommandOutput struct {
	ID          string `json:"id"`
	Category    string `json:"category"`
	Description string `json:"description"`
	RiskLevel   string `json:"risk_level"`
	Enabled     bool   `json:"enabled"`
}

func listAllCommands(showAll bool) []*whitelistCommandOutput {
	var result []*whitelistCommandOutput
	for _, info := range whitelistMgr.ListCommands() {
		if showAll || info.Enabled {
			result = append(result, &whitelistCommandOutput{
				ID:          info.Name,
				Category:    info.Category,
				Description: info.Description,
				RiskLevel:   string(info.RiskLevel),
				Enabled:     info.Enabled,
			})
		}
	}
	return result
}

func listByCategory(category string, showAll bool) []*whitelistCommandOutput {
	var result []*whitelistCommandOutput
	for _, info := range whitelistMgr.ListCommandsByCategory(category) {
		if showAll || info.Enabled {
			result = append(result, &whitelistCommandOutput{
				ID:          info.Name,
				Category:    info.Category,
				Description: info.Description,
				RiskLevel:   string(info.RiskLevel),
				Enabled:     info.Enabled,
			})
		}
	}
	return result
}

func listByRiskLevel(riskLevel string, showAll bool) []*whitelistCommandOutput {
	var result []*whitelistCommandOutput
	for _, info := range whitelistMgr.ListCommandsByRisk(whitelist.RiskLevel(riskLevel)) {
		if showAll || info.Enabled {
			result = append(result, &whitelistCommandOutput{
				ID:          info.Name,
				Category:    info.Category,
				Description: info.Description,
				RiskLevel:   string(info.RiskLevel),
				Enabled:     info.Enabled,
			})
		}
	}
	return result
}

func init() {
	rootCmd.AddCommand(whitelistCmd)
	whitelistCmd.AddCommand(whitelistListCmd)
	whitelistCmd.AddCommand(whitelistStatusCmd)
	whitelistCmd.AddCommand(whitelistEnableCmd)
	whitelistCmd.AddCommand(whitelistDisableCmd)
	whitelistCmd.AddCommand(whitelistResetCmd)
	whitelistCmd.AddCommand(whitelistPathCmd)

	whitelistListCmd.Flags().Bool("all", false, "Show all commands including disabled ones")
	whitelistListCmd.Flags().String("risk", "", "Filter by risk level (low, medium, high, critical)")
	whitelistListCmd.Flags().String("category", "", "Filter by category")

	markAsLocal(whitelistCmd)
	markAsLocal(whitelistListCmd)
	markAsLocal(whitelistStatusCmd)
	markAsLocal(whitelistEnableCmd)
	markAsLocal(whitelistDisableCmd)
	markAsLocal(whitelistResetCmd)
	markAsLocal(whitelistPathCmd)

	markNeedsWhitelist(whitelistCmd)
	markNeedsWhitelist(whitelistListCmd)
	markNeedsWhitelist(whitelistStatusCmd)
	markNeedsWhitelist(whitelistEnableCmd)
	markNeedsWhitelist(whitelistDisableCmd)
	markNeedsWhitelist(whitelistResetCmd)
	markNeedsWhitelist(whitelistPathCmd)

	setCommandID(whitelistListCmd, "whitelist.list")
	setCommandID(whitelistStatusCmd, "whitelist.status")
	setCommandID(whitelistEnableCmd, "whitelist.enable")
	setCommandID(whitelistDisableCmd, "whitelist.disable")
	setCommandID(whitelistResetCmd, "whitelist.reset")
	setCommandID(whitelistPathCmd, "whitelist.path")
}
