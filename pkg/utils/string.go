package utils

import (
	"bufio"
	"strings"

	"github.com/liuzl/gocc"
	"github.com/rs/zerolog/log"
)

func NovelContentFormat(input string) string {
	var result strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(input))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		result.WriteString("　　" + line + "\n")
	}

	return result.String()
}

func TrimRowSpaceInMultiParagraph(input string) string {
	var result strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(input))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		result.WriteString(line + "\n")
	}

	return result.String()
}

func ConvertTraditionalToSimplified(input string) string {
	converter, _ := gocc.New("t2s")
	output, err := converter.Convert(input)
	if err != nil {
		log.Warn().AnErr("err", err).Msg("ConvertTraditionalToSimplified err")
	}

	return output
}
