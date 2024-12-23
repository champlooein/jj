package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of jj",
	Long:  `All software has versions. This is jj's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("jj version jj0.0.4")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
