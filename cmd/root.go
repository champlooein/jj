package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	verbose bool

	rootCmd = &cobra.Command{
		Use:   "jj",
		Short: "A novel crawler",
		Long:  "A novel crawler that supports downloading novels from different book repo and archiving it into txt format according to unified typesetting rules",
	}
)

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}
