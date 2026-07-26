package cmd

import (
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync operations",
}

var fileChangesCmd = &cobra.Command{
	Use:   "file-changes [path]",
	Short: "Sync file changes from cloud",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		path := args[0]

		result, err := cd2Client.Sync().SyncFileChangesFromCloud(ctx, path)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var startListenerCmd = &cobra.Command{
	Use:   "start-listener [path]",
	Short: "Start cloud event listener",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		path := args[0]

		err := cd2Client.Sync().StartCloudEventListener(ctx, path)
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "started"})
	},
}

var stopListenerCmd = &cobra.Command{
	Use:   "stop-listener [path]",
	Short: "Stop cloud event listener",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		path := args[0]

		err := cd2Client.Sync().StopCloudEventListener(ctx, path)
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "stopped"})
	},
}

var walkTestCmd = &cobra.Command{
	Use:   "walk-test <path>",
	Short: "Walk through folder test",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		path := args[0]

		result, err := cd2Client.Sync().WalkThroughFolderTest(ctx, path)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var cd1UserDataCmd = &cobra.Command{
	Use:   "cd1-user-data",
	Short: "Get CloudDrive1 user data",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()

		result, err := cd2Client.Sync().GetCloudDrive1UserData(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	syncCmd.AddCommand(fileChangesCmd)
	syncCmd.AddCommand(startListenerCmd)
	syncCmd.AddCommand(stopListenerCmd)
	syncCmd.AddCommand(walkTestCmd)
	syncCmd.AddCommand(cd1UserDataCmd)

	setCommandID(fileChangesCmd, "sync.file-changes")
	setCommandID(startListenerCmd, "sync.start-listener")
	setCommandID(stopListenerCmd, "sync.stop-listener")
	setCommandID(walkTestCmd, "sync.walk-test")
	setCommandID(cd1UserDataCmd, "sync.cd1-user-data")
}
