package crawler

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/champlooein/jj/internal/consts"
	"github.com/champlooein/jj/pkg/utils"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type Crawler interface {
	Search(novelTitle, novelAuthor string) (novelNo string, meta NovelMetaInfo, err error)
	Info(novelNo string) (info NovelMetaInfo, err error)
	Crawl(novelNo string, n int) (chapterTitleToContentArr []*lo.Entry[string, string], err error)
	Save(novelTitle, novelIntro string, chapterTitleToContentArr []*lo.Entry[string, string], path string) (err error)
}

type repo struct {
	name  string
	url   string
	intro string
}

type NovelMetaInfo struct {
	Title  string
	Author string
	Intro  string
}

type saver struct{}

var (
	banxiaRepo = repo{
		name: "banxia",
		url:  "https://www.xbanxia.com",
	}
	shukuRepo = repo{
		name: "shuku",
		url:  "https://www.52shuku.vip",
	}

	DefaultCrawler = defaultBanxiaCrawler
)

func NewCrawlerFromRepo(r string) Crawler {
	switch r {
	case banxiaRepo.name:
		return defaultBanxiaCrawler
	case shukuRepo.name:
		return defaultShukuCrawler
	default:
		glog.Warning("unknown repo, using default")
		return DefaultCrawler
	}
}

func (s saver) Save(novelTitle, novelIntro string, chapterTitleToContentArr []*lo.Entry[string, string], path string) error {
	if ext := filepath.Ext(path); ext != "" {
		if ext != ".txt" {
			return errors.Errorf("path invalid, not a txt file, path: %v", path)
		}
	} else {
		path += string(filepath.Separator) + novelTitle + ".txt"
	}

	var sb strings.Builder
	if len(novelIntro) > 0 {
		sb.WriteString(fmt.Sprintf("%s\n%s\n", consts.NovelIntroKey, novelIntro))
	}
	for _, novelTitleToContent := range chapterTitleToContentArr {
		sb.WriteString(fmt.Sprintf("%s\n%s\n", novelTitleToContent.Key, utils.NovelContentFormat(novelTitleToContent.Value)))
	}

	return utils.WriteToFile(path, sb.String())
}

func (r repo) String() string {
	s := fmt.Sprintf("Repo: %s\nWebsite: %s\n", r.name, r.url)
	if strings.TrimSpace(r.intro) != "" {
		s += fmt.Sprintf("Intro:\n  %s\n", r.intro)
	}

	return s
}

func (r NovelMetaInfo) String() string {
	return fmt.Sprintf("书名：%s\n作者：%s\n简介：\n%s", r.Title, r.Author, r.Intro)
}
