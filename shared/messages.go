package shared

type Query_ struct {
	Selector  string
	Attribute string
}

type HttpAsyncMessage struct {
	Url      string
	Callback string
	JoinKey  string
	Queries  []Query_
}

type HttpAsyncResponse struct {
	JoinKey string
	Results [][]string
}

type EmptyRequest struct{}
