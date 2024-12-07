package consts

import "errors"

var (
	EmptySearchResultErr = errors.New("empty search result")
	UnsupportSearchErr   = errors.New("unsupport search")
)
