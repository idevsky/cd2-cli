//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"
)

func TestIntegrationSession_Revoke(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sessions, err := c.Session().GetSessions(ctx)
	if err != nil {
		t.Fatalf("GetSessions failed: %v", err)
	}

	if len(sessions.Sessions) == 0 {
		t.Skip("No sessions to revoke")
	}

	sessionToRevoke := sessions.Sessions[0].Id
	err = c.Session().RevokeSession(ctx, sessionToRevoke)
	if err != nil {
		t.Fatalf("RevokeSession failed: %v", err)
	}

	t.Logf("Successfully revoked session: %s", sessionToRevoke)
}

func TestIntegrationSession_RevokeOthers(t *testing.T) {
	c := getAuthClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sessionsBefore, err := c.Session().GetSessions(ctx)
	if err != nil {
		t.Fatalf("GetSessions failed: %v", err)
	}

	countBefore := len(sessionsBefore.Sessions)

	err = c.Session().RevokeOtherSessions(ctx)
	if err != nil {
		t.Fatalf("RevokeOtherSessions failed: %v", err)
	}

	sessionsAfter, err := c.Session().GetSessions(ctx)
	if err != nil {
		t.Fatalf("GetSessions after revoke failed: %v", err)
	}

	countAfter := len(sessionsAfter.Sessions)

	if countAfter >= countBefore {
		t.Errorf("Expected fewer sessions after revoke-others, got %d before and %d after", countBefore, countAfter)
	}

	t.Logf("Revoked %d session(s), %d remaining", countBefore-countAfter, countAfter)
}
