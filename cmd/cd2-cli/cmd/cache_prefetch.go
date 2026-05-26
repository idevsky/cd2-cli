package cmd

import (
	"encoding/json"
	"strconv"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

func init() {
	cacheCmd.AddCommand(prefetchCmd)
	cacheCmd.AddCommand(cancelPrefetchCmd)
	cacheCmd.AddCommand(closeReaderCmd)
	cacheCmd.AddCommand(prefetchHintsCmd)

	prefetchCmd.Flags().String("ranges", "", "JSON array of ByteRange objects")
	cancelPrefetchCmd.Flags().StringSlice("hint-ids", []string{}, "Hint IDs to cancel (as strings)")

	setCommandID(prefetchCmd, "cache.prefetch")
	setCommandID(cancelPrefetchCmd, "cache.cancel-prefetch")
	setCommandID(closeReaderCmd, "cache.close-reader")
	setCommandID(prefetchHintsCmd, "cache.prefetch-hints")
}

var prefetchCmd = &cobra.Command{
	Use:   "prefetch <path>",
	Short: "Prefetch file ranges (accept JSON ranges)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		path := args[0]
		rangesJson, _ := cmd.Flags().GetString("ranges")

		var ranges []*pb.ByteRange
		if rangesJson != "" {
			if err := json.Unmarshal([]byte(rangesJson), &ranges); err != nil {
				return err
			}
		}

		result, err := cd2Client.System().PrefetchFileRanges(ctx, &pb.PrefetchFileRangesRequest{
			Path:   path,
			Ranges: ranges,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var cancelPrefetchCmd = &cobra.Command{
	Use:   "cancel-prefetch <path>",
	Short: "Cancel file prefetch",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		path := args[0]
		hintIdsStr, _ := cmd.Flags().GetStringSlice("hint-ids")

		hintIds := make([]uint64, 0)
		for _, idStr := range hintIdsStr {
			id, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				return err
			}
			hintIds = append(hintIds, id)
		}

		err := cd2Client.System().CancelFilePrefetch(ctx, &pb.CancelFilePrefetchRequest{
			Path:    path,
			HintIds: hintIds,
		})
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "cancelled"})
	},
}

var closeReaderCmd = &cobra.Command{
	Use:   "close-reader <path>",
	Short: "Close file reader",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		path := args[0]

		err := cd2Client.System().CloseFileReader(ctx, path)
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "closed"})
	},
}

var prefetchHintsCmd = &cobra.Command{
	Use:   "prefetch-hints",
	Short: "Get active prefetch hints",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		result, err := cd2Client.System().GetActivePrefetchHints(ctx)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}
