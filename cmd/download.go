package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/forgoty/go-reimaging/cmd/validators"
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
	userId, error := strconv.Atoi(args[0])
	if error != nil {
		panic(errors.New("Unvalid USERID has been provided. Need Integer"))
	}
	fmt.Println("download called with", userId)

	validPath, error := validators.ValidateDownloadDir(path)
	if error != nil {
		panic(error)
	}
	fmt.Println("Path:", validPath)
}
