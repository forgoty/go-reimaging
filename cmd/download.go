package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/forgoty/go-reimaging/cmd/auth"
	"github.com/forgoty/go-reimaging/cmd/validators"
)

var downloadCmd = &cobra.Command{
	Use:   "download USERID",
	Short: "Download photo albums of specific user",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("More then one argument provided")
		}

		_, error := strconv.Atoi(args[0])
		if error != nil {
			return errors.New("Unvalid USERID has been provided. Need Integer")
		}

		return nil
	},
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
	userId := args[0]

	_, error := validators.ValidateDownloadDir(path)
	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	}
	vk := auth.GetClient(Auth)
	response := GetAlbums(vk, userId)
	for _, album := range response.Items {
		fmt.Printf("%s(%d) - id:%d\n", album.Title, album.Size, album.ID)
	}
}
