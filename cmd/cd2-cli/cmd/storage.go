package cmd

import (
	"github.com/spf13/cobra"
)

var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Storage provider management",
	Long:  "Manage cloud storage providers (add, remove, configure, and list storage)",
}

func init() {
	rootCmd.AddCommand(storageCmd)
}
