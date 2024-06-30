package shared

type Query_ struct {
	Selector  string
	Attribute string
}

type HttpAsyncMessage struct {
	Url      string
	Callback string
	Queries  map[string]Query_
}

type HttpAsyncResponse struct {
	Results map[string][]string
}
