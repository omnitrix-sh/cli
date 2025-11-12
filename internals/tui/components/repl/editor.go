package repl

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/omnitrix-sh/cli/editor"
	"github.com/omnitrix-sh/cli/internals/app"
	agent "github.com/omnitrix-sh/cli/internals/llm/agents"
	"github.com/omnitrix-sh/cli/internals/tui/layout"
	"github.com/omnitrix-sh/cli/internals/tui/styles"
	util "github.com/omnitrix-sh/cli/internals/utils"
)

type EditorCmp interface {
	tea.Model
	layout.Focusable
	layout.Sizeable
	layout.Bordered
	layout.Bindings
}

type editorCmp struct {
	app        *app.App
	editor     editor.Editor
	editorMode editor.EditorMode
	sessionID  string
	focused    bool
	width      int
	height     int
}

type editorKeyMap struct {
	SendMessage    key.Binding
	SendMessageI   key.Binding
	InsertMode     key.Binding
	NormaMode      key.Binding
	VisualMode     key.Binding
	VisualLineMode key.Binding
}

var editorKeyMapValue = editorKeyMap{
	SendMessage: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "send message normal mode"),
	),
	SendMessageI: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "send message insert mode"),
	),
	InsertMode: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "insert mode"),
	),
	NormaMode: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "normal mode"),
	),
	VisualMode: key.NewBinding(
		key.WithKeys("v"),
		key.WithHelp("v", "visual mode"),
	),
	VisualLineMode: key.NewBinding(
		key.WithKeys("V"),
		key.WithHelp("V", "visual line mode"),
	),
}

func (m *editorCmp) Init() tea.Cmd {
	return m.editor.Init()
}

func (m *editorCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case editor.EditorModeMsg:
		m.editorMode = msg.Mode
	case SelectedSessionMsg:
		if msg.SessionID != m.sessionID {
			m.sessionID = msg.SessionID
		}
	}
	if m.IsFocused() {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, editorKeyMapValue.SendMessage):
				if m.editorMode == editor.ModeNormal {
					return m, m.Send()
				}
			case key.Matches(msg, editorKeyMapValue.SendMessageI):
				if m.editorMode == editor.ModeInsert {
					return m, m.Send()
				}
			}
		}
		u, cmd := m.editor.Update(msg)
		m.editor = u.(editor.Editor)
		return m, cmd
	}
	return m, nil
}

func (m *editorCmp) Blur() tea.Cmd {
	m.focused = false
	return nil
}

func (m *editorCmp) BorderText() map[layout.BorderPosition]string {
	title := "New Message"
	if m.focused {
		title = lipgloss.NewStyle().Foreground(styles.Primary).Render(title)
	}
	return map[layout.BorderPosition]string{
		layout.BottomLeftBorder: title,
	}
}

func (m *editorCmp) Focus() tea.Cmd {
	m.focused = true
	return m.editor.Tick()
}

func (m *editorCmp) GetSize() (int, int) {
	return m.width, m.height
}

func (m *editorCmp) IsFocused() bool {
	return m.focused
}

func (m *editorCmp) SetSize(width int, height int) {
	m.width = width
	m.height = height
	m.editor.SetSize(width, height)
}

func (m *editorCmp) Send() tea.Cmd {
	return func() tea.Msg {
		messages, _ := m.app.Messages.List(m.sessionID)
		if hasUnfinishedMessages(messages) {
			return util.InfoMsg("Assistant is still working on the previous message")
		}
		a, _ := agent.NewCoderAgent(m.app)

		content := strings.Join(m.editor.GetBuffer().Lines(), "\n")
		go a.Generate(m.sessionID, content)

		return m.editor.Reset()
	}
}

func (m *editorCmp) View() string {
	return m.editor.View()
}

func (m *editorCmp) BindingKeys() []key.Binding {
	return layout.KeyMapToSlice(editorKeyMapValue)
}

func NewEditorCmp(app *app.App) EditorCmp {
	editor := editor.NewEditor(
		editor.WithFileName("message.md"),
	)
	return &editorCmp{
		app:    app,
		editor: editor,
	}
}
