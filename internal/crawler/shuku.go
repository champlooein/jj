package crawler

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/champlooein/jj/pkg/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

var (
	chapterTitleMatchRegexp = `(?m)^第[ \f\r\t\v0123456789一二三四五六七八九十零〇百千两]+[章卷节](([ \f\r\t\v]+.*$)|($))`

	shukuNovelDetailUrl = shukuRepo.url + "/yanqing/%s.html"

	defaultShukuCrawler = shukuCrawler{saver: saver{}}
)

type shukuCrawler struct {
	saver
}

func (c shukuCrawler) Info(novelNo string) (info NovelMetaInfo, err error) {
	var title, author, intro string

	httpResp, err := utils.HttpGet(fmt.Sprintf(shukuNovelDetailUrl, novelNo), false)
	if err != nil {
		return info, err
	}
	defer httpResp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(httpResp.Body)
	if err != nil {
		return info, errors.Wrap(err, "can't parse html")
	}

	s := doc.Find(".article-title").Text()
	title, author = strings.Split(s, "_")[0], s[strings.Index(s, "_")+len("_"):strings.Index(s, "【")]

	var firstPageUrl string
	doc.Find(".list").Find("a").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		firstPageUrl, _ = selection.Attr("href")
		return false
	})

	pageContent, err := c.crawlPage(firstPageUrl)
	if err != nil {
		return info, err
	}

	matches := regexp.MustCompile(chapterTitleMatchRegexp).FindAllStringIndex(pageContent, 1)
	if len(matches) == 0 {
		return NovelMetaInfo{Title: title, Author: author}, nil
	}

	intro = c.extractIntro(pageContent[:matches[0][0]])
	return NovelMetaInfo{Title: title, Author: author, Intro: intro}, nil
}

func (c shukuCrawler) Crawl(novelNo string, n int) (chapterTitleToContentArr []*lo.Entry[string, string], err error) {
	httpResp, err := utils.HttpGet(fmt.Sprintf(shukuNovelDetailUrl, novelNo), false)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(httpResp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse html")
	}

	var pageTitleToUrlArr []*lo.Entry[string, string]
	doc.Find(".list").Find("a").Each(func(i int, s *goquery.Selection) {
		pageUrl, _ := s.Attr("href")
		pageTitleToUrlArr = append(pageTitleToUrlArr, &lo.Entry[string, string]{Key: s.Text(), Value: pageUrl})
	})

	var (
		eg errgroup.Group
		m  sync.Map
	)
	eg.SetLimit(n)

	for _, pageTitleToUrl := range pageTitleToUrlArr {
		eg.Go(func() error {
			pageContent, subErr := c.crawlPage(pageTitleToUrl.Value)
			if subErr != nil {
				return errors.WithMessagef(subErr, "crawl page err, Title: %s", pageTitleToUrl.Key)
			}
			log.Debug().Str("pageTitle", pageTitleToUrl.Key).Msg("crawl page ok")

			m.Store(pageTitleToUrl.Key, pageContent)
			return nil
		})
	}
	if err = eg.Wait(); err != nil {
		return nil, err
	}

	var sb strings.Builder
	for _, pageTitleToUrl := range pageTitleToUrlArr {
		v, _ := m.Load(pageTitleToUrl.Key)
		pageContent := v.(string)

		_, err = sb.WriteString(fmt.Sprintf("%s\n", pageContent))
		if err != nil {
			return nil, errors.WithMessagef(err, "write page content err, Title: %s", pageTitleToUrl.Key)
		}
	}

	return c.pageToChapterFormat(sb.String()), nil
}

func (c shukuCrawler) crawlPage(chapterUrl string) (string, error) {
	httpResp, err := utils.HttpGet(chapterUrl, false)
	if err != nil {
		return "", err
	}
	defer httpResp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(httpResp.Body)
	if err != nil {
		return "", errors.Wrap(err, "can't parse html")
	}

	return utils.TrimRowSpaceInMultiParagraph(doc.Find(".book_con").Text()), nil
}

func (c shukuCrawler) pageToChapterFormat(input string) (chapterTitleToContentArr []*lo.Entry[string, string]) {
	matches := regexp.MustCompile(chapterTitleMatchRegexp).FindAllStringIndex(input, -1)
	if len(matches) == 0 {
		chapterTitleToContentArr = []*lo.Entry[string, string]{{Key: "正文", Value: input}}
	}

	for i, match := range matches {
		contentStart, contentEnd := match[1], 0
		if i < len(matches)-1 {
			contentEnd = matches[i+1][0]
		} else {
			contentEnd = len(input)
		}
		chapterTitleToContentArr = append(chapterTitleToContentArr, &lo.Entry[string, string]{
			Key:   input[match[0]:match[1]],
			Value: input[contentStart:contentEnd],
		})
	}

	return chapterTitleToContentArr
}

func (c shukuCrawler) extractIntro(input string) string {
	headers := []string{"简介：", "简介", "文案：", "文案"}

	var result strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(input))

	for scanner.Scan() {
		for _, header := range headers {
			line, cnt := scanner.Text(), 0
			if line != header {
				cnt++
			}
			if cnt == len(headers) {
				result.WriteString(line)
			}
		}
	}

	return utils.NovelContentFormat(result.String())
}
