package crawler

import (
	"fmt"
	"strings"
	"sync"

	"github.com/champlooein/jj/internal/consts"
	"github.com/rs/zerolog/log"

	"github.com/PuerkitoBio/goquery"
	"github.com/champlooein/jj/pkg/utils"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"golang.org/x/net/html"
	"golang.org/x/sync/errgroup"
)

var (
	banxiaNovelDetailUrl = banxiaRepo.url + "/books/%s.html"

	defaultBanxiaCrawler = banxiaCrawler{saver: saver{}}
)

type banxiaCrawler struct {
	saver
}

func (c banxiaCrawler) Info(novelNo string) (info NovelMetaInfo, err error) {
	var title, author, intro string

	httpResp, err := utils.HttpGet(fmt.Sprintf(banxiaNovelDetailUrl, novelNo), false)
	if err != nil {
		return info, err
	}
	defer httpResp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(httpResp.Body)
	if err != nil {
		return info, errors.Wrapf(err, "can't parse html, url:%s", fmt.Sprintf(banxiaNovelDetailUrl, novelNo))
	}
	if strings.Contains(doc.Find("title").Text(), consts.BanxiaNovelNotExistTitle) {
		return info, errors.Wrapf(err, "novel not exist, novelNo:%s", novelNo)
	}

	title = utils.ConvertTraditionalToSimplified(doc.Find(".book-describe").Find("h1").Text())
	intro = utils.NovelContentFormat(strings.Replace(utils.ConvertTraditionalToSimplified(utils.ExtractNovelTextFromHtml(doc.Find(".describe-html").Nodes[0])), "文案:", "", 1))
	doc.Find(".book-describe").Find("p").EachWithBreak(func(i int, s *goquery.Selection) bool {
		author = utils.ConvertTraditionalToSimplified(strings.Replace(s.Text(), "作者︰", "", 1))
		return false
	})

	return NovelMetaInfo{Title: title, Author: author, Intro: intro}, nil
}

func (c banxiaCrawler) Crawl(novelNo string, n int) (chapterTitleToContentArr []*lo.Entry[string, string], err error) {
	httpResp, err := utils.HttpGet(fmt.Sprintf(banxiaNovelDetailUrl, novelNo), false)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(httpResp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse html")
	}

	chapterTitleToUrlArr := make([]*lo.Entry[string, string], 0)
	doc.Find(".book-list").Find("a").Each(func(i int, s *goquery.Selection) {
		if href, exist := s.Attr("href"); exist {
			chapterTitleToUrlArr = append(chapterTitleToUrlArr, &lo.Entry[string, string]{Key: utils.ConvertTraditionalToSimplified(s.Text()), Value: banxiaRepo.url + href})
		}
	})

	var (
		eg errgroup.Group
		m  sync.Map
	)
	eg.SetLimit(n)

	for _, chapterTitleToUrl := range chapterTitleToUrlArr {
		eg.Go(func() error {
			var chapterContent string

			chapterContent, err = c.crawlChapter(chapterTitleToUrl.Value)
			if err != nil {
				return errors.WithMessagef(err, "crawl chapter err, Title: %s", chapterTitleToUrl.Key)
			}
			log.Debug().Str("chapterTitle", chapterTitleToUrl.Key).Msg("crawl chapter ok")

			m.Store(chapterTitleToUrl.Key, chapterContent)
			return nil
		})
	}
	err = eg.Wait()
	if err != nil {
		return nil, err
	}

	for _, chapterTitleToUrl := range chapterTitleToUrlArr {
		chapterContent, _ := m.Load(chapterTitleToUrl.Key)
		chapterTitleToContentArr = append(chapterTitleToContentArr, &lo.Entry[string, string]{Key: chapterTitleToUrl.Key, Value: chapterContent.(string)})
	}

	return chapterTitleToContentArr, nil
}

func (c banxiaCrawler) crawlChapter(chapterUrl string) (string, error) {
	httpResp, err := utils.HttpGet(chapterUrl, false)
	if err != nil {
		return "", err
	}
	defer httpResp.Body.Close()

	// 解析html
	doc, err := goquery.NewDocumentFromReader(httpResp.Body)
	if err != nil {
		return "", errors.Wrap(err, "can't parse html")
	}

	// 删除第一行章节名和最后一行的广告
	text := doc.Find("#nr1").Contents().FilterFunction(func(i int, s *goquery.Selection) bool {
		if len(s.Nodes) == 0 || s.Nodes[0].Type != html.TextNode {
			return false
		}
		if len(strings.TrimSpace(s.Nodes[0].Data)) == 0 {
			return false
		}

		return true
	})
	text.First().Remove()
	doc.Find("#nr1").Find("span").Remove()

	return utils.ConvertTraditionalToSimplified(utils.ExtractNovelTextFromHtml(doc.Find("#nr1").Nodes[0])), nil
}
