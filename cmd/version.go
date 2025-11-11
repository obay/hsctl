package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// This will be set by the build process using ldflags
	version = "dev"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of hscli",
	Long:  `Print the version number and build information of hscli.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("hscli version %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
