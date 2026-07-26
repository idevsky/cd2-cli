//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/clouddrive/cd2-cli/internal/client"
	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

type TestConfig struct {
	Host    string
	Port    string
	User    string
	Pass    string
	Address string
}

func loadTestConfig() *TestConfig {
	cfg := &TestConfig{
		Host: getEnvOrDefault("CD2_HOST", "localhost"),
		Port: getEnvOrDefault("CD2_PORT", "19798"),
		User: getEnvOrDefault("CD2_USER", "admin"),
		Pass: getEnvOrDefault("CD2_PASS", "admin123"),
	}
	cfg.Address = fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	return cfg
}

func getTestClient(t *testing.T) *client.Client {
	cfg := loadTestConfig()

	c, err := client.NewClient(client.Config{
		Address: cfg.Address,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	return c
}

func getAuthClient(t *testing.T) *client.Client {
	cfg := loadTestConfig()
	c := getTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := c.Public().GetToken(ctx, &pb.GetTokenRequest{
		UserName: cfg.User,
		Password: cfg.Pass,
	})
	if err != nil {
		c.Close()
		t.Fatalf("Login failed: %v", err)
	}

	c.SetToken(resp.Token)
	return c
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
