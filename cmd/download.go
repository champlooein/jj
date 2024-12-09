package cmd

import (
	"fmt"
	"time"

	"github.com/champlooein/jj/internal/crawler"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	output string
	limit  int

	downloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download novel",
		Long:  `Crawler novel from novel repo and save it to disk`,
		Run: func(cmd *cobra.Command, args []string) {
			if verbose {
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			} else {
				zerolog.SetGlobalLevel(zerolog.InfoLevel)
			}

			crawler := crawler.NewCrawlerFromRepo(repo)
			info, err := crawler.Info(novelNo)
			if err != nil {
				log.Error().Err(err).Msg("get novel info error")
				return
			}

			fmt.Println(info.String())
			fmt.Print("Continue download?(yes[y] or no[n]) ")

			var s string
			fmt.Scanln(&s)
			switch s {
			case "yes", "y", "Y":
				log.Info().Msg("Download Start...")
				start := time.Now()

				chapterTitleToContentArr, err := crawler.Crawl(novelNo, limit)
				if err != nil {
					log.Err(err).Msg("Crawl novel error")
					return
				}

				err = crawler.Save(info.Title, info.Intro, chapterTitleToContentArr, output)
				if err != nil {
					log.Err(err).Msg("Save novel error")
					return
				}

				fmt.Printf("Download finish, enjoy yourself! (cost:%vs)", time.Since(start).Seconds())
			default:
				fmt.Println("Download terminated!")
			}

		},
	}
)

func init() {
	downloadCmd.Flags().StringVarP(&output, "output", "o", "./", "output folder")
	downloadCmd.Flags().StringVarP(&repo, "repo", "r", "banxia", "novel repo")
	downloadCmd.Flags().IntVarP(&limit, "limit", "l", 3, "concurrent crawling limit")
	downloadCmd.Flags().StringVarP(&novelNo, "no", "n", "", "novel number")
	if err := downloadCmd.MarkFlagRequired("no"); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(downloadCmd)
}
