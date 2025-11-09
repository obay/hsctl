package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	// These will be set by the build process using ldflags
	version   = "dev"
	commit    = "unknown"
	buildDate = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of hsctl",
	Long:  `Print the version number and build information of hsctl.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("hsctl version %s\n", version)
		fmt.Printf("Commit: %s\n", commit)
		fmt.Printf("Build Date: %s\n", buildDate)
		fmt.Printf("Go Version: %s\n", runtime.Version())
		fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
