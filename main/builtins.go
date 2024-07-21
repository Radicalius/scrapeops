package main

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/Radicalius/scrapeops/shared"
)

func HttpAsyncHandler(message shared.HttpAsyncMessage, ctx shared.Context) error {
	res := shared.HttpAsyncResponse{
		JoinKey: message.JoinKey,
		Results: make([][]string, 0),
	}

	resp, err := http.Get(message.Url)
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	for _, query := range message.Queries {
		selection := doc.Find(query.Selector)
		res.Results = append(res.Results, selection.Map(func(i int, s *goquery.Selection) string {
			if query.Attribute == "text" {
				return s.Text()
			} else {
				return s.AttrOr(query.Attribute, "")
			}
		}))
	}

	return shared.Emit(ctx, message.Callback, res)
}
