package shared

type Item map[string]interface{}

type Dom interface {
	CssSelect(cssSelector string) Dom
	CssSelectAll(cssSelector string) []Dom
	Attr(string) string
	Text() string
}

type Provider interface {
	IsRelevant(i *Item) bool
	GetUrl(i *Item) string
	Apply(d Dom, i *Item) ([]*Item, error)
}

type Schematizer[T any] interface {
	IsRelevant(i *Item) bool
	Schematize(i *Item) T
}
