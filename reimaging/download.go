package reimaging

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/forgoty/go-reimaging/reimaging/downloadcommand"
	"github.com/forgoty/go-reimaging/reimaging/validator"
	vkw "github.com/forgoty/go-reimaging/reimaging/vkwrapper"
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

	downloadCmd.Flags().BoolVarP(&System, "system", "s", false, "Enable system albums for download")
	downloadCmd.Flags().IntVarP(&AlbumID, "album-id", "", 0, "Use specific album ID to download")
	downloadCmd.Flags().StringVarP(&path, "path", "p", "", "Set Download Folder")
}

func download(args []string) {
	validateDownloadDirictory()
	albumDownloader := getAlbumDownloader(parseUserID(args[0]))

	if AlbumID != 0 {
		albumDownloader.DownloadAlbumByID(AlbumID)
		return
	}
	albumDownloader.DownloadAll()
}

func validateDownloadDirictory() {
	_, err := validator.ValidateDownloadDir(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func parseUserID(rawUserID string) int {
	userID, err := strconv.Atoi(rawUserID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return userID
}

func getAlbumDownloader(userID int) *downloadcommand.AlbumDownloader {
	vkWrapper := vkw.NewVKWrapper()
	options := downloadcommand.NewDownloadOptions(userID, path, System)
	return downloadcommand.NewAlbumDownloader(vkWrapper, options)
}
