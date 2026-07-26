package cmd

import (
	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

var offlineCmd = &cobra.Command{
	Use:   "offline",
	Short: "Offline download operations",
}

var addOfflineCmd = &cobra.Command{
	Use:   "add [urls] [folder]",
	Short: "Add offline download files",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		urls := args[0]
		folder := args[1]

		result, err := cd2Client.Offline().AddOfflineFiles(ctx, &pb.AddOfflineFileRequest{
			Urls:     urls,
			ToFolder: folder,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var removeOfflineCmd = &cobra.Command{
	Use:   "remove <cloud-name> <account-id> <info-hash...>",
	Short: "Remove offline files",
	Args:  cobra.MinimumNArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		cloudName := args[0]
		accountId := args[1]
		infoHashes := args[2:]

		result, err := cd2Client.Offline().RemoveOfflineFiles(ctx, &pb.RemoveOfflineFilesRequest{
			CloudName:      cloudName,
			CloudAccountId: accountId,
			InfoHashes:     infoHashes,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var listOfflineCmd = &cobra.Command{
	Use:   "list [path]",
	Short: "List offline files by path",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		path := args[0]

		result, err := cd2Client.Offline().ListOfflineFilesByPath(ctx, path)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var listAllOfflineCmd = &cobra.Command{
	Use:   "list-all <cloud-name> <account-id>",
	Short: "List all offline files",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		cloudName := args[0]
		accountId := args[1]
		page, _ := cmd.Flags().GetUint32("page")

		result, err := cd2Client.Offline().ListAllOfflineFiles(ctx, &pb.OfflineFileListAllRequest{
			CloudName:      cloudName,
			CloudAccountId: accountId,
			Page:           page,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var quotaOfflineCmd = &cobra.Command{
	Use:   "quota <cloud-name> <account-id>",
	Short: "Get offline quota info",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		cloudName := args[0]
		accountId := args[1]

		result, err := cd2Client.Offline().GetOfflineQuotaInfo(ctx, &pb.OfflineQuotaRequest{
			CloudName:      cloudName,
			CloudAccountId: accountId,
		})
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var restartOfflineCmd = &cobra.Command{
	Use:   "restart <cloud-name> <account-id> <info-hash>",
	Short: "Restart offline task",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		cloudName := args[0]
		accountId := args[1]
		infoHash := args[2]

		err := cd2Client.Offline().RestartOfflineTask(ctx, &pb.RestartOfflineFileRequest{
			CloudName:      cloudName,
			CloudAccountId: accountId,
			InfoHash:       infoHash,
		})
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "restarted"})
	},
}

var clearOfflineCmd = &cobra.Command{
	Use:   "clear [cloud-name] [account-id]",
	Short: "Clear offline files",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		cloudName := args[0]
		accountId := args[1]
		filter, _ := cmd.Flags().GetString("filter")
		deleteFiles, _ := cmd.Flags().GetBool("delete")

		var filterEnum pb.ClearOfflineFileRequest_Filter
		switch filter {
		case "finished":
			filterEnum = pb.ClearOfflineFileRequest_Finished
		case "error":
			filterEnum = pb.ClearOfflineFileRequest_Error
		case "downloading":
			filterEnum = pb.ClearOfflineFileRequest_Downloading
		default:
			filterEnum = pb.ClearOfflineFileRequest_All
		}

		err := cd2Client.Offline().ClearOfflineFiles(ctx, &pb.ClearOfflineFileRequest{
			CloudName:      cloudName,
			CloudAccountId: accountId,
			Filter:         filterEnum,
			DeleteFiles:    deleteFiles,
		})
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "cleared"})
	},
}

func init() {
	rootCmd.AddCommand(offlineCmd)

	offlineCmd.AddCommand(addOfflineCmd)
	offlineCmd.AddCommand(removeOfflineCmd)
	offlineCmd.AddCommand(listOfflineCmd)
	offlineCmd.AddCommand(listAllOfflineCmd)
	offlineCmd.AddCommand(quotaOfflineCmd)
	offlineCmd.AddCommand(restartOfflineCmd)
	offlineCmd.AddCommand(clearOfflineCmd)

	clearOfflineCmd.Flags().String("filter", "all", "Filter type: all, finished, error, downloading")
	clearOfflineCmd.Flags().Bool("delete", false, "Delete files")
	listAllOfflineCmd.Flags().Uint32("page", 1, "Page number")

	setCommandID(addOfflineCmd, "offline.add")
	setCommandID(removeOfflineCmd, "offline.remove")
	setCommandID(listOfflineCmd, "offline.list")
	setCommandID(listAllOfflineCmd, "offline.list-all")
	setCommandID(quotaOfflineCmd, "offline.quota")
	setCommandID(restartOfflineCmd, "offline.restart")
	setCommandID(clearOfflineCmd, "offline.clear")
}
