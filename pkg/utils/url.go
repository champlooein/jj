package utils

import (
	"net/url"
	"path"
	"strings"

	"github.com/pkg/errors"
)

func GetUrlLastSegment(inputURL string) (string, error) {
	u, err := url.Parse(inputURL)
	if err != nil {
		return "", errors.Wrap(err, "parse url error")
	}

	segs := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(segs) == 0 {
		return "", errors.Errorf("no path segments in the url, url : %s", inputURL)
	}

	lastSeg := strings.TrimSuffix(segs[len(segs)-1], path.Ext(segs[len(segs)-1]))

	return lastSeg, nil
}
