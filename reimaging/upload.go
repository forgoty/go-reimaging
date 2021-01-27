package reimaging

import (
	"fmt"

	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload PATH",
	Short: "Upload photo directory to vk album",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		upload(args)
	},
}

var title string

func init() {
	rootCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().IntVarP(&AlbumID, "album-id", "", 0, "Using an existing album to upload")
	uploadCmd.Flags().StringVarP(&title, "title", "t", "", "Create new vk album with title and uploads")
}

func upload(args []string) {
	fmt.Println("upload called with", args[0])
	fmt.Println(AlbumID, title)
}
