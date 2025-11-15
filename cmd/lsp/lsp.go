package main

import (
	"context"
	"fmt"

	"github.com/omnitrix-sh/cli/internals/llm/tools"
)

func main() {
	t := tools.NewSourcegraphTool()
	r, _ := t.Run(context.Background(), tools.ToolCall{
		Input: `{"query": "context.WithCancel lang:go"}`,
	})
	fmt.Println(r.Content)
}
