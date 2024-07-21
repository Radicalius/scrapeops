package main

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/Radicalius/scrapeops/shared"
)

type Dom struct {
	Html *goquery.Selection
}

func (d Dom) CssSelect(selector string) shared.Dom {
	return Dom{
		Html: d.Html.Find(selector).First(),
	}
}

func (d Dom) CssSelectAll(selector string) []shared.Dom {
	res := make([]shared.Dom, 0)
	d.Html.Find(selector).Each(func(i int, s *goquery.Selection) {
		res = append(res, Dom{
			Html: s,
		})
	})

	return res
}

func (d Dom) Text() string {
	return d.Html.Text()
}

func (d Dom) Attr(attrName string) string {
	val, exists := d.Html.Attr(attrName)
	if !exists {
		return ""
	}

	return val
}

func HttpFetch(url string) (shared.Dom, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Dom{
		Html: doc.Selection,
	}, nil
}
