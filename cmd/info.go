package cmd

import (
	"fmt"

	"github.com/champlooein/jj/internal/crawler"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	repo    string
	novelNo string

	infoCmd = &cobra.Command{
		Use:   "info",
		Short: "Show novel info",
		Long:  `Show novel info like title、author and intro`,
		Run: func(cmd *cobra.Command, args []string) {
			crawler := crawler.NewCrawlerFromRepo(repo)
			info, err := crawler.Info(novelNo)
			if err != nil {
				log.Err(err).Msg("get novel info error")
				return
			}

			fmt.Println(info.String())
			return
		},
	}
)

func init() {
	infoCmd.Flags().StringVarP(&repo, "repo", "r", "banxia", "novel repo")
	infoCmd.Flags().StringVarP(&novelNo, "no", "n", "", "novel number")
	if err := infoCmd.MarkFlagRequired("no"); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(infoCmd)
}
