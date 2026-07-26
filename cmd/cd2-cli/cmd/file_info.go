package cmd

import (
	pb "github.com/clouddrive/cd2-cli/pkg/proto"
	"github.com/spf13/cobra"
)

func init() {
	fileCmd.AddCommand(searchFilesCmd)
	fileCmd.AddCommand(downloadUrlCmd)
	fileCmd.AddCommand(filePropertiesCmd)
	fileCmd.AddCommand(spaceInfoCmd)
	fileCmd.AddCommand(membershipsCmd)
	fileCmd.AddCommand(metadataCmd)
	fileCmd.AddCommand(originalPathCmd)
	fileCmd.AddCommand(addSharedLinkCmd)

	searchFilesCmd.Flags().String("path", "/", "Search path")
	searchFilesCmd.Flags().Bool("fuzzy", false, "Use fuzzy matching")
	searchFilesCmd.Flags().Bool("content", false, "Search file content")
	downloadUrlCmd.Flags().Bool("preview", false, "Get preview URL")
	downloadUrlCmd.Flags().Bool("direct", false, "Get direct URL if available")
	addSharedLinkCmd.Flags().String("password", "", "Shared link password")

	setCommandID(searchFilesCmd, "file.search")
	setCommandID(downloadUrlCmd, "file.download-url")
	setCommandID(filePropertiesCmd, "file.properties")
	setCommandID(spaceInfoCmd, "file.space")
	setCommandID(membershipsCmd, "file.memberships")
	setCommandID(metadataCmd, "file.metadata")
	setCommandID(originalPathCmd, "file.original-path")
	setCommandID(addSharedLinkCmd, "file.add-shared-link")
}

var searchFilesCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for files",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		query := args[0]
		path, _ := cmd.Flags().GetString("path")
		fuzzyMatch, _ := cmd.Flags().GetBool("fuzzy")
		contentSearch, _ := cmd.Flags().GetBool("content")

		files, err := cd2Client.File().GetSearchResults(ctx, &pb.SearchRequest{
			SearchFor:     query,
			Path:          path,
			FuzzyMatch:    fuzzyMatch,
			ForceRefresh:  false,
			ContentSearch: &contentSearch,
		})
		if err != nil {
			return err
		}
		return outputResult(files)
	},
}

var downloadUrlCmd = &cobra.Command{
	Use:   "download-url [path]",
	Short: "Get download URL for a file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		path := args[0]
		preview, _ := cmd.Flags().GetBool("preview")
		direct, _ := cmd.Flags().GetBool("direct")

		url, err := cd2Client.File().GetDownloadUrl(ctx, &pb.GetDownloadUrlPathRequest{
			Path:         path,
			Preview:      preview,
			LazyRead:     false,
			GetDirectUrl: direct,
		})
		if err != nil {
			return err
		}
		return outputResult(url)
	},
}

var filePropertiesCmd = &cobra.Command{
	Use:   "properties [path]",
	Short: "Get file detail properties",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		path := args[0]

		props, err := cd2Client.File().GetFileDetailProperties(ctx, path)
		if err != nil {
			return err
		}
		return outputResult(props)
	},
}

var spaceInfoCmd = &cobra.Command{
	Use:   "space [path]",
	Short: "Get space info for a cloud path",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		path := args[0]

		info, err := cd2Client.File().GetSpaceInfo(ctx, path)
		if err != nil {
			return err
		}
		return outputResult(info)
	},
}

var membershipsCmd = &cobra.Command{
	Use:   "memberships [path]",
	Short: "Get cloud account memberships",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		path := args[0]

		memberships, err := cd2Client.File().GetCloudMemberships(ctx, path)
		if err != nil {
			return err
		}
		return outputResult(memberships)
	},
}

var metadataCmd = &cobra.Command{
	Use:   "metadata [path]",
	Short: "Get file metadata",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		path := args[0]

		meta, err := cd2Client.File().GetMetaData(ctx, path)
		if err != nil {
			return err
		}
		return outputResult(meta)
	},
}

var originalPathCmd = &cobra.Command{
	Use:   "original-path [path]",
	Short: "Get original path from search result",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		path := args[0]

		result, err := cd2Client.File().GetOriginalPath(ctx, path)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var addSharedLinkCmd = &cobra.Command{
	Use:   "add-shared-link [url] [folder]",
	Short: "Add shared link to folder",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := getTimeoutContext()
		defer cancel()
		url := args[0]
		folder := args[1]
		password, _ := cmd.Flags().GetString("password")

		err := cd2Client.File().AddSharedLink(ctx, &pb.AddSharedLinkRequest{
			SharedLinkUrl:  url,
			SharedPassword: &password,
			ToFolder:       folder,
		})
		if err != nil {
			return err
		}
		return outputResult(map[string]string{"status": "success"})
	},
}
