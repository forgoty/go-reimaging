package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list USERID",
	Short: "Get list os user photo albums",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args[0])
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&Auth, "auth", "a", false, "Enable authorization")
	listCmd.Flags().BoolVarP(&System, "system", "s", false, "Enable system albums for download")
}
