package reimaging

import (
	"github.com/forgoty/go-reimaging/reimaging/uploadcommand"
	vkw "github.com/forgoty/go-reimaging/reimaging/vkwrapper"
	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload PATH",
	Short: "Upload photo directory to vk album",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		upload(args)
	},
}

var (
	title      string
	uploadPath string
)

func init() {
	rootCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().IntVarP(&AlbumID, "album-id", "", 0, "Using an existing album to upload")
	uploadCmd.Flags().StringVarP(&title, "title", "t", "", "Create new vk album with title and uploads")
	uploadCmd.Flags().StringVarP(&uploadPath, "path", "p", "", "Set Upload Folder")
}

func upload(args []string) {
	vkWrapper := vkw.NewVKWrapper()
	albumUploader := uploadcommand.NewAlbumUploader(vkWrapper)
	if AlbumID != 0 {
		ids := []int{AlbumID}
		albums := vkWrapper.GetAlbumsByAlbumIds(ids)
		albumUploader.Upload(AlbumID, uploadPath, albums[0].Title)
		return
	}
	album := vkWrapper.CreateAlbum(title)
	albumUploader.Upload(album.ID, uploadPath, title)
}
