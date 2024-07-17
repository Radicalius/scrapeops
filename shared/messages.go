package shared

type Query_ struct {
	Selector  string
	Attribute string
}

type HttpAsyncMessage struct {
	Url      string
	Callback string
	Queries  []Query_
}

type HttpAsyncResponse struct {
	Results [][]string
}

type EmptyRequest struct{}
