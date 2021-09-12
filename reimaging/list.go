package reimaging

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	vkw "github.com/forgoty/go-reimaging/reimaging/vkwrapper"
)

var listCmd = &cobra.Command{
	Use:   "list USERID",
	Short: "Get list os user photo albums",
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
		list(args)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&System, "system", "s", false, "Enable system albums")
}

func list(args []string) {
	userID, _ := strconv.Atoi(args[0])
	albums := vkw.NewVKWrapper().GetAlbums(userID, System)

	for _, album := range albums {
		fmt.Printf("%s(%d) - id:%d\n", album.Title, album.Size, album.ID)
	}
}
