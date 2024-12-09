package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/champlooein/jj/cmd"
)

func main() {
	cmd.Execute()
}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, FormatTimestamp: func(i interface{}) string {
		return time.Now().Format("2006-01-02 15:04:05") // 自定义时间格式
	}})
}
