package page

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/omnitrix-sh/cli/internals/app"
	"github.com/omnitrix-sh/cli/internals/tui/components/repl"
	"github.com/omnitrix-sh/cli/internals/tui/layout"
)

var ReplPage PageID = "repl"

func NewReplPage(app *app.App) tea.Model {
	return layout.NewBentoLayout(
		layout.BentoPanes{
			layout.BentoLeftPane:        repl.NewSessionsCmp(app),
			layout.BentoRightTopPane:    repl.NewMessagesCmp(app),
			layout.BentoRightBottomPane: repl.NewEditorCmp(app),
		},
		layout.WithBentoLayoutCurrentPane(layout.BentoRightBottomPane),
	)
}
