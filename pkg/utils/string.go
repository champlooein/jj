package utils

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"github.com/liuzl/gocc"
	"github.com/rs/zerolog/log"
)

// NovelContentFormat 段缩进
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

// TrimRowSpaceInMultiParagraph 去掉每行前的空白及空白行
func TrimRowSpaceInMultiParagraph(input string) string {
	var result strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(input))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

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

func RemoveStringSpaces(input string) string {
	sb := strings.Builder{}
	for _, char := range input {
		if char == ' ' || char == '\t' || char == '\n' || char == '\r' {
			continue
		}

		sb.WriteRune(char)
	}

	return sb.String()
}

func FormatChapterTitle(input string, n int) string {
	if ok, _ := regexp.Match(`^第\d+章`, []byte(input)); ok {
		return input
	}

	return fmt.Sprintf("第%d章 %s", n, input)
}
