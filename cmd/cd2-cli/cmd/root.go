package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/clouddrive/cd2-cli/internal/client"
	"github.com/clouddrive/cd2-cli/internal/whitelist"
)

const (
	annotationCommandID      = "cd2-cli:command-id"
	annotationIsLocal        = "cd2-cli:is-local"
	annotationSkipWhitelist  = "cd2-cli:skip-whitelist"
	annotationNeedsWhitelist = "cd2-cli:needs-whitelist"
)

type exitCodeError struct {
	data map[string]interface{}
}

func (e *exitCodeError) Error() string {
	if msg, ok := e.data["status"].(string); ok {
		return msg
	}
	return "exit code error"
}

func exitCode(data map[string]interface{}) error {
	return &exitCodeError{data: data}
}

var (
	cfgFile          string
	whitelistCfgFile string
	serverAddr       string
	authToken        string
	outputJSON       bool
	useTLS           bool
	skipTLSVerify    bool
	timeout          string
	whitelistMgr     *whitelist.Manager
	initialized      bool
)

var cd2Client *client.Client

var rootCmd = &cobra.Command{
	Use:   "cd2-cli",
	Short: "CloudDrive2 CLI - Manage CloudDrive2 instances via gRPC API",
	Long: `cd2-cli is a command-line interface tool for managing CloudDrive2 instances.
It communicates with CloudDrive2 via its gRPC API and outputs structured JSON data,
making it suitable for use by AI agents and automation tools.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return persistentPreRunE(cmd)
	},
}

func Execute() {
	markCompletionAsLocal()
	if err := rootCmd.Execute(); err != nil {
		outputError(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/cd2-cli.yaml)")
	rootCmd.PersistentFlags().StringVar(&whitelistCfgFile, "whitelist-config", "", "whitelist config file (default is $HOME/.cd2-cli-whitelist.yaml)")
	rootCmd.PersistentFlags().StringVarP(&serverAddr, "server", "s", "", "CloudDrive2 server address")
	rootCmd.PersistentFlags().StringVarP(&authToken, "token", "t", "", "Authentication token")
	rootCmd.PersistentFlags().BoolVarP(&outputJSON, "json", "j", false, "Output in JSON format")
	rootCmd.PersistentFlags().BoolVar(&useTLS, "tls", false, "Use TLS connection")
	rootCmd.PersistentFlags().BoolVar(&skipTLSVerify, "skip-tls-verify", false, "Skip TLS certificate verification")
	rootCmd.PersistentFlags().StringVar(&timeout, "timeout", "30s", "Timeout for remote API calls (e.g. 30s, 1m)")

	viper.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server"))
	viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("json", rootCmd.PersistentFlags().Lookup("json"))
	viper.BindPFlag("tls", rootCmd.PersistentFlags().Lookup("tls"))
	viper.BindPFlag("skip-tls-verify", rootCmd.PersistentFlags().Lookup("skip-tls-verify"))
	viper.BindPFlag("timeout", rootCmd.PersistentFlags().Lookup("timeout"))

	viper.SetEnvPrefix("CD2_CLI")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	registerCommands()
}

func persistentPreRunE(cmd *cobra.Command) error {
	if !initialized {
		if err := loadConfig(); err != nil {
			return err
		}
		initialized = true
	}

	shouldInitWhitelist := needsWhitelistInit(cmd)
	if shouldInitWhitelist && whitelistMgr == nil {
		if err := initWhitelist(); err != nil {
			return fmt.Errorf("failed to initialize whitelist: %w", err)
		}
	}

	if shouldSkipWhitelistCheck(cmd) {
		return nil
	}

	if err := checkWhitelist(cmd); err != nil {
		return err
	}

	if cd2Client == nil {
		return initClient()
	}
	return nil
}

func needsWhitelistInit(cmd *cobra.Command) bool {
	if cmd.Annotations != nil {
		if _, ok := cmd.Annotations[annotationIsLocal]; ok {
			if _, ok := cmd.Annotations[annotationNeedsWhitelist]; !ok {
				return false
			}
		}
	}
	for parent := cmd.Parent(); parent != nil; parent = parent.Parent() {
		if parent.Annotations != nil {
			if _, ok := parent.Annotations[annotationIsLocal]; ok {
				if _, ok := parent.Annotations[annotationNeedsWhitelist]; !ok {
					return false
				}
			}
		}
	}
	return true
}

func shouldSkipWhitelistCheck(cmd *cobra.Command) bool {
	if cmd.Name() == "completion" || cmd.Name() == "help" || cmd.Name() == "__complete" || cmd.Name() == "__completeScript" {
		return true
	}
	if cmd.Annotations != nil {
		if _, ok := cmd.Annotations[annotationSkipWhitelist]; ok {
			return true
		}
		if _, ok := cmd.Annotations[annotationIsLocal]; ok {
			return true
		}
	}
	for parent := cmd.Parent(); parent != nil; parent = parent.Parent() {
		if parent.Name() == "completion" || parent.Name() == "help" {
			return true
		}
		if parent.Annotations != nil {
			if _, ok := parent.Annotations[annotationSkipWhitelist]; ok {
				return true
			}
			if _, ok := parent.Annotations[annotationIsLocal]; ok {
				return true
			}
		}
	}
	return false
}

func checkWhitelist(cmd *cobra.Command) error {
	if whitelistMgr == nil || !whitelistMgr.IsEnabled() {
		return nil
	}

	commandID := getCommandID(cmd)
	if commandID == "" {
		return fmt.Errorf("command '%s' has no command ID registered and is blocked by whitelist (default deny)", cmd.CommandPath())
	}

	allowed, reason := whitelistMgr.IsAllowed(commandID)
	if !allowed {
		info, _ := whitelistMgr.GetCommand(commandID)
		riskLevel := "unknown"
		if info != nil {
			riskLevel = string(info.RiskLevel)
		}
		return fmt.Errorf("command '%s' is blocked by whitelist (risk level: %s): %s", commandID, riskLevel, reason)
	}
	return nil
}

func getCommandID(cmd *cobra.Command) string {
	if cmd.Annotations != nil {
		if id, ok := cmd.Annotations[annotationCommandID]; ok {
			return id
		}
	}
	return ""
}

func loadConfig() error {
	viper.SetDefault("server", "localhost:19798")
	viper.SetDefault("json", true)
	viper.SetDefault("tls", false)
	viper.SetDefault("skip-tls-verify", false)
	viper.SetDefault("timeout", "30s")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}

		configPath := filepath.Join(home, ".config", "cd2-cli.yaml")
		viper.SetConfigFile(configPath)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok && !os.IsNotExist(err) {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	return nil
}

func initWhitelist() error {
	whitelistPath := getWhitelistConfigPath()
	mgr, err := whitelist.NewManager(whitelistPath)
	if err != nil {
		return err
	}
	whitelistMgr = mgr
	return nil
}

func getWhitelistConfigPath() string {
	if whitelistCfgFile != "" {
		return whitelistCfgFile
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(os.Getenv("HOME"), ".cd2-cli-whitelist.yaml")
	}
	return filepath.Join(home, ".cd2-cli-whitelist.yaml")
}

func initClient() error {
	timeoutDur := parseTimeout()

	cfg := client.Config{
		Address:       viper.GetString("server"),
		Token:         viper.GetString("token"),
		Timeout:       timeoutDur,
		UseTLS:        viper.GetBool("tls"),
		SkipVerifyTLS: viper.GetBool("skip-tls-verify"),
	}

	c, err := client.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}

	cd2Client = c
	return nil
}

func parseTimeout() time.Duration {
	timeoutStr := viper.GetString("timeout")
	if timeoutStr == "" {
		return 30 * time.Second
	}
	dur, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return 30 * time.Second
	}
	return dur
}

func getTimeoutContext() (context.Context, context.CancelFunc) {
	dur := parseTimeout()
	return context.WithTimeout(context.Background(), dur)
}

func closeClient() {
	if cd2Client != nil {
		cd2Client.Close()
	}
}

func parseProtoJSON(data []byte, msg proto.Message) error {
	return protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}.Unmarshal(data, msg)
}

func outputResult(data interface{}) error {
	if viper.GetBool("json") {
		if msg, ok := data.(proto.Message); ok {
			result, err := protojson.MarshalOptions{
				EmitUnpopulated: true,
			}.Marshal(msg)
			if err != nil {
				return err
			}
			fmt.Println(string(result))
			return nil
		}
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetEscapeHTML(false)
		return encoder.Encode(data)
	}

	fmt.Println(data)
	return nil
}

func outputError(err error) {
	if exitErr, ok := err.(*exitCodeError); ok {
		if viper.GetBool("json") {
			json.NewEncoder(os.Stdout).Encode(exitErr.data)
		} else {
			fmt.Println(exitErr.data)
		}
		return
	}
	if viper.GetBool("json") {
		result := map[string]string{"error": err.Error()}
		json.NewEncoder(os.Stdout).Encode(result)
	} else {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
}

func setCommandID(cmd *cobra.Command, id string) {
	if cmd.Annotations == nil {
		cmd.Annotations = make(map[string]string)
	}
	cmd.Annotations[annotationCommandID] = id
}

func markAsLocal(cmd *cobra.Command) {
	if cmd.Annotations == nil {
		cmd.Annotations = make(map[string]string)
	}
	cmd.Annotations[annotationIsLocal] = "true"
}

func markSkipWhitelist(cmd *cobra.Command) {
	if cmd.Annotations == nil {
		cmd.Annotations = make(map[string]string)
	}
	cmd.Annotations[annotationSkipWhitelist] = "true"
}

func markNeedsWhitelist(cmd *cobra.Command) {
	if cmd.Annotations == nil {
		cmd.Annotations = make(map[string]string)
	}
	cmd.Annotations[annotationNeedsWhitelist] = "true"
}
