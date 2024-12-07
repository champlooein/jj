package utils

import (
	"net/http"

	"github.com/pkg/errors"
)

func HttpGet(url string, noRedirect bool) (*http.Response, error) {
	cli := http.DefaultClient
	if noRedirect {
		cli = &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse // 禁用自动重定向
			},
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating http request, url: %s", url)
	}

	resp, err := cli.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "error executing http request, url: %s", url)
	}

	if resp.StatusCode != http.StatusOK {
		if noRedirect && (resp.StatusCode == http.StatusMovedPermanently || resp.StatusCode == http.StatusFound) {
			return resp, nil
		}

		return nil, errors.Errorf("error executing http request, url: %s, http_code: %d", url, resp.StatusCode)
	}

	return resp, nil
}
