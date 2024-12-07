package cmd

import (
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jj",
	Short: "A novel crawler",
	Long:  "A novel crawler that supports downloading novels from different book repo and archiving it into txt format according to unified typesetting rules",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		glog.Fatal(err.Error())
	}
}
