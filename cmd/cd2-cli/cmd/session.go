package cmd

import (
	"github.com/spf13/cobra"
)

var sessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Session management",
}

var listSessionsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all active sessions",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.Session().GetSessions(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var revokeSessionCmd = &cobra.Command{
	Use:   "revoke [session-id]",
	Short: "Revoke a session",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		sessionId := args[0]

		err := cd2Client.Session().RevokeSession(ctx, sessionId)
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "revoked"})
	},
}

var revokeOtherSessionsCmd = &cobra.Command{
	Use:   "revoke-others",
	Short: "Revoke all other sessions",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		err := cd2Client.Session().RevokeOtherSessions(ctx)
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "revoked"})
	},
}

func init() {
	rootCmd.AddCommand(sessionCmd)

	sessionCmd.AddCommand(listSessionsCmd)
	sessionCmd.AddCommand(revokeSessionCmd)
	sessionCmd.AddCommand(revokeOtherSessionsCmd)

	setCommandID(listSessionsCmd, "session.list")
	setCommandID(revokeSessionCmd, "session.revoke")
	setCommandID(revokeOtherSessionsCmd, "session.revoke-others")
}
