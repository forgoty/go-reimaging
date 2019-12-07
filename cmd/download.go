package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download USERID",
	Short: "Download photo albums of specific user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		download(args)
	},
}

var path string

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().BoolVarP(&Auth, "auth", "a", false, "Enable authorization")
	downloadCmd.Flags().BoolVarP(&System, "system", "s", false, "Enable system albums for download")
	downloadCmd.Flags().IntVarP(&AlbumId, "album-id", "", 0, "Use specific album ID to download")
	downloadCmd.Flags().StringVarP(&path, "path", "p", "", "Set Download Folder")
}

func download(args []string) {
	fmt.Println("download called with", args[0])
	fmt.Println(System, Auth, AlbumId, path)
}
