package utils

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

func ExtractNovelTextFromHtml(node *html.Node) string {
	s, _ := goquery.NewDocumentFromNode(node).Html()
	s = strings.ReplaceAll(s, "<br>", "\n")
	s = strings.ReplaceAll(s, "<br/>", "\n")

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(s))
	return doc.Text()
}
