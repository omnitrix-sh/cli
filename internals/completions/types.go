package completions

type CompletionItem struct {
	Title string
	Value string
}

type CompletionProvider interface {
	GetId() string
	GetEntry() CompletionItem
	GetChildEntries(query string) ([]CompletionItem, error)
}
