package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var authTokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Manage CLI authentication token",
}

var authTokenSetCmd = &cobra.Command{
	Use:   "set [TOKEN]",
	Short: "Set authentication token in config",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		token := args[0]

		configPath := getConfigPath()
		if err := saveTokenToConfig(configPath, token); err != nil {
			return err
		}

		viper.Set("token", token)

		return outputResult(map[string]interface{}{
			"success":      true,
			"config_path":  configPath,
			"message":      "token saved to config",
			"token_length": len(token),
		})
	},
}

var authTokenClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear authentication token from config",
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath := getConfigPath()
		if err := clearTokenFromConfig(configPath); err != nil {
			return err
		}

		return outputResult(map[string]interface{}{
			"success":     true,
			"config_path": configPath,
			"message":     "token cleared from config",
		})
	},
}

var authTokenShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current authentication token",
	RunE: func(cmd *cobra.Command, args []string) error {
		redact, _ := cmd.Flags().GetBool("redact")

		token := viper.GetString("token")
		if token == "" {
			return outputResult(map[string]interface{}{
				"token":   "",
				"source":  "none",
				"message": "no token configured",
			})
		}

		source := getTokenSource(cmd)

		if redact {
			if len(token) > 8 {
				token = token[:4] + "..." + token[len(token)-4:]
			} else {
				token = "***"
			}
		}

		return outputResult(map[string]interface{}{
			"token":    token,
			"source":   source,
			"redacted": redact,
		})
	},
}

func getConfigPath() string {
	if cfgFile != "" {
		return cfgFile
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(os.Getenv("HOME"), ".config", "cd2-cli.yaml")
	}
	return filepath.Join(home, ".config", "cd2-cli.yaml")
}

func getTokenSource(cmd *cobra.Command) string {
	root := cmd
	for root.HasParent() {
		root = root.Parent()
	}

	if root.PersistentFlags().Changed("token") {
		return "flag"
	}
	if os.Getenv("CD2_CLI_TOKEN") != "" {
		return "env"
	}
	return "config"
}

func saveTokenToConfig(configPath string, token string) error {
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	existingData := make(map[string]interface{})
	if data, err := os.ReadFile(configPath); err == nil {
		yaml.Unmarshal(data, &existingData)
	}

	existingData["token"] = token

	data, err := yaml.Marshal(existingData)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return atomicWriteFile(configPath, data, 0600)
}

func atomicWriteFile(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	tmpFile, err := os.CreateTemp(dir, ".cd2-cli-config-*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	cleanup := func() {
		tmpFile.Close()
		os.Remove(tmpPath)
	}

	if _, err := tmpFile.Write(data); err != nil {
		cleanup()
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	if err := tmpFile.Chmod(perm); err != nil {
		cleanup()
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		cleanup()
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		cleanup()
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

func clearTokenFromConfig(configPath string) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg map[string]interface{}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	delete(cfg, "token")

	newData, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return atomicWriteFile(configPath, newData, 0600)
}

func init() {
	authCmd.AddCommand(authTokenCmd)
	authTokenCmd.AddCommand(authTokenSetCmd)
	authTokenCmd.AddCommand(authTokenClearCmd)
	authTokenCmd.AddCommand(authTokenShowCmd)

	authTokenShowCmd.Flags().Bool("redact", false, "Redact token in output")

	markAsLocal(authTokenCmd)
	markAsLocal(authTokenSetCmd)
	markAsLocal(authTokenClearCmd)
	markAsLocal(authTokenShowCmd)

	setCommandID(authTokenSetCmd, "auth.token-set")
	setCommandID(authTokenClearCmd, "auth.token-clear")
	setCommandID(authTokenShowCmd, "auth.token-show")
}
