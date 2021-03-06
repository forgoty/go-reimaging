package reimaging

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "reimaging",
	Version: "0.1",
	Short:   "Simple photo downloader/uploader for vk.com",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// Flag Global Values
var (
	System  bool
	AlbumID int
)

func init() {
	cobra.OnInitialize()
}
