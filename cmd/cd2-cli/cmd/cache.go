package cmd

import (
	"strings"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Disk cache operations",
}

func init() {
	rootCmd.AddCommand(cacheCmd)

	cacheCmd.AddCommand(cacheStatsCmd)
	cacheCmd.AddCommand(cachePurgeCmd)
	cacheCmd.AddCommand(cacheFoldersCmd)
	cacheCmd.AddCommand(setEvictionCmd)
	cacheCmd.AddCommand(setFolderCacheCmd)
	cacheCmd.AddCommand(removeFolderCacheCmd)

	setFolderCacheCmd.Flags().Bool("enabled", true, "Enable disk cache")
	setFolderCacheCmd.Flags().Uint64("max-file-size", 0, "Max file size in bytes (0 = no limit)")
	setFolderCacheCmd.Flags().Uint64("min-file-size", 0, "Min file size in bytes (0 = no minimum)")
	setFolderCacheCmd.Flags().StringSlice("extensions", []string{}, "Extensions to cache (e.g. mp4,mkv)")

	setCommandID(cacheStatsCmd, "cache.stats")
	setCommandID(cachePurgeCmd, "cache.purge")
	setCommandID(cacheFoldersCmd, "cache.list-disk")
	setCommandID(setEvictionCmd, "cache.set-eviction")
	setCommandID(setFolderCacheCmd, "cache.set-folder")
	setCommandID(removeFolderCacheCmd, "cache.remove-folder")
}

var cacheStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Get disk cache stats",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.System().GetFileBufferDiskCacheStats(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var cachePurgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Purge disk cache",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		err := cd2Client.System().PurgeFileBufferDiskCache(ctx)
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "purged"})
	},
}

var cacheFoldersCmd = &cobra.Command{
	Use:   "folders",
	Short: "List folders with disk cache",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.System().ListDiskCacheFolders(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var setEvictionCmd = &cobra.Command{
	Use:   "set-eviction <strategy>",
	Short: "Set disk cache eviction strategy (LRU, FIFO)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		strategy := args[0]

		var strategyEnum pb.EvictionStrategy
		switch strings.ToUpper(strategy) {
		case "LARGEST_FIRST":
			strategyEnum = pb.EvictionStrategy_LARGEST_FIRST
		case "SMALLEST_FIRST":
			strategyEnum = pb.EvictionStrategy_SMALLEST_FIRST
		case "LRU":
			strategyEnum = pb.EvictionStrategy_LRU
		default:
			strategyEnum = pb.EvictionStrategy_LRU
		}

		err := cd2Client.System().SetDiskCacheEvictionStrategy(ctx, &pb.SetDiskCacheEvictionStrategyRequest{
			Strategy: strategyEnum,
		})
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "updated"})
	},
}

var setFolderCacheCmd = &cobra.Command{
	Use:   "set-folder <path>",
	Short: "Set folder disk cache",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		path := args[0]
		enabled, _ := cmd.Flags().GetBool("enabled")
		maxFileSize, _ := cmd.Flags().GetUint64("max-file-size")
		minFileSize, _ := cmd.Flags().GetUint64("min-file-size")
		extensions, _ := cmd.Flags().GetStringSlice("extensions")

		req := &pb.SetFolderDiskCacheRequest{
			Path:        path,
			Enabled:     enabled,
			MaxFileSize: maxFileSize,
			MinFileSize: minFileSize,
			Extensions:  extensions,
		}

		err := cd2Client.System().SetFolderDiskCache(ctx, req)
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "updated"})
	},
}

var removeFolderCacheCmd = &cobra.Command{
	Use:   "remove-folder <path>",
	Short: "Remove folder disk cache",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		path := args[0]

		err := cd2Client.System().RemoveFolderDiskCache(ctx, path)
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "removed"})
	},
}
