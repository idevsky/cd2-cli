package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	systemCmd.AddCommand(dirCacheCmd)

	dirCacheCmd.AddCommand(dirCacheSizeCmd)
	dirCacheCmd.AddCommand(dirCacheVacuumCmd)
	dirCacheCmd.AddCommand(dirCacheProgressCmd)
	dirCacheCmd.AddCommand(dirCacheTableCmd)

	setCommandID(dirCacheSizeCmd, "system.dir-cache-size")
	setCommandID(dirCacheVacuumCmd, "system.dir-cache-vacuum")
	setCommandID(dirCacheProgressCmd, "system.dir-cache-progress")
	setCommandID(dirCacheTableCmd, "system.dir-cache-table")
}

var dirCacheCmd = &cobra.Command{
	Use:   "dir-cache",
	Short: "Directory cache commands",
}

var dirCacheSizeCmd = &cobra.Command{
	Use:   "size",
	Short: "Get directory cache DB size",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.System().GetDirCacheDbSize(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var dirCacheVacuumCmd = &cobra.Command{
	Use:   "vacuum",
	Short: "Vacuum directory cache",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		err := cd2Client.System().VacuumDirCache(ctx)
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "vacuumed"})
	},
}

var dirCacheProgressCmd = &cobra.Command{
	Use:   "progress",
	Short: "Get vacuum progress",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.System().GetVacuumProgress(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var dirCacheTableCmd = &cobra.Command{
	Use:   "table",
	Short: "Get directory cache table",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.System().GetDirCacheTable(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}
