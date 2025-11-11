package page

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/omnitrix-sh/cli/internals/tui/components/repl"
	"github.com/omnitrix-sh/cli/internals/tui/layout"
)

var ReplPage PageID = "repl"

func NewReplPage() tea.Model {
	return layout.NewBentoLayout(
		layout.BentoPanes{
			layout.BentoLeftPane:        repl.NewThreadsCmp(),
			layout.BentoRightTopPane:    repl.NewMessagesCmp(),
			layout.BentoRightBottomPane: repl.NewEditorCmp(),
		},
	)
}
