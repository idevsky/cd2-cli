//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

func TestIntegrationWebhook_Template(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	template, err := c.Webhook().GetWebhookConfigTemplate(ctx)
	if err != nil {
		t.Fatalf("GetWebhookConfigTemplate failed: %v", err)
	}

	if template.Result == "" {
		t.Error("Expected non-empty template")
	}

	t.Logf("Template: %s", template.Result)
}

func TestIntegrationWebhook_AddAndRemove(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	hooksBefore, err := c.Webhook().GetWebhookConfigs(ctx)
	if err != nil {
		t.Fatalf("GetWebhookConfigs failed: %v", err)
	}
	countBefore := len(hooksBefore.Webhooks)

	template, err := c.Webhook().GetWebhookConfigTemplate(ctx)
	if err != nil {
		t.Fatalf("GetWebhookConfigTemplate failed: %v", err)
	}

	filename := "test-webhook-" + time.Now().Format("20060102150405") + ".json"
	err = c.Webhook().AddWebhookConfig(ctx, &pb.WebhookRequest{
		FileName: filename,
		Content:  template.Result,
	})
	if err != nil {
		t.Fatalf("AddWebhookConfig failed: %v", err)
	}

	hooksAfterAdd, err := c.Webhook().GetWebhookConfigs(ctx)
	if err != nil {
		t.Fatalf("GetWebhookConfigs after add failed: %v", err)
	}
	countAfterAdd := len(hooksAfterAdd.Webhooks)

	if countAfterAdd != countBefore+1 {
		t.Errorf("Expected %d webhooks after add, got %d", countBefore+1, countAfterAdd)
	}

	err = c.Webhook().RemoveWebhookConfig(ctx, filename)
	if err != nil {
		t.Fatalf("RemoveWebhookConfig failed: %v", err)
	}

	hooksAfterRemove, err := c.Webhook().GetWebhookConfigs(ctx)
	if err != nil {
		t.Fatalf("GetWebhookConfigs after remove failed: %v", err)
	}
	countAfterRemove := len(hooksAfterRemove.Webhooks)

	if countAfterRemove != countBefore {
		t.Errorf("Expected %d webhooks after remove, got %d", countBefore, countAfterRemove)
	}

	t.Logf("Added and removed webhook: %s", filename)
}

func TestIntegrationWebhook_Change(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	template, err := c.Webhook().GetWebhookConfigTemplate(ctx)
	if err != nil {
		t.Fatalf("GetWebhookConfigTemplate failed: %v", err)
	}

	filename := "test-webhook-change-" + time.Now().Format("20060102150405") + ".json"
	err = c.Webhook().AddWebhookConfig(ctx, &pb.WebhookRequest{
		FileName: filename,
		Content:  template.Result,
	})
	if err != nil {
		t.Fatalf("AddWebhookConfig failed: %v", err)
	}
	defer c.Webhook().RemoveWebhookConfig(ctx, filename)

	modifiedContent := template.Result + " // modified"
	err = c.Webhook().ChangeWebhookConfig(ctx, &pb.WebhookRequest{
		FileName: filename,
		Content:  modifiedContent,
	})
	if err != nil {
		t.Fatalf("ChangeWebhookConfig failed: %v", err)
	}

	hooks, err := c.Webhook().GetWebhookConfigs(ctx)
	if err != nil {
		t.Fatalf("GetWebhookConfigs failed: %v", err)
	}

	var found bool
	for _, h := range hooks.Webhooks {
		if h.FileName == filename {
			found = true
			t.Logf("Changed webhook: %s, content length: %d", h.FileName, len(h.Content))
			break
		}
	}

	if !found {
		t.Error("Expected to find changed webhook in list")
	}
}
