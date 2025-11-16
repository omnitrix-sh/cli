package styles

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

func HuhTheme() *huh.Theme {
	t := huh.ThemeBase()

	t.Focused.Base = t.Focused.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Focused.Title = t.Focused.Title.Foreground(lipgloss.Color(Text))
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(lipgloss.Color(Text))
	t.Focused.Directory = t.Focused.Directory.Foreground(lipgloss.Color(Text))
	t.Focused.Description = t.Focused.Description.Foreground(lipgloss.Color(SubText0))
	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(lipgloss.Color(Red))
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(lipgloss.Color(Red))
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(lipgloss.Color(Blue))
	t.Focused.NextIndicator = t.Focused.NextIndicator.Foreground(lipgloss.Color(Blue))
	t.Focused.PrevIndicator = t.Focused.PrevIndicator.Foreground(lipgloss.Color(Blue))
	t.Focused.Option = t.Focused.Option.Foreground(lipgloss.Color(Text))
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(lipgloss.Color(Blue))
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(lipgloss.Color(Green))
	t.Focused.SelectedPrefix = t.Focused.SelectedPrefix.Foreground(lipgloss.Color(Green))
	t.Focused.UnselectedPrefix = t.Focused.UnselectedPrefix.Foreground(lipgloss.Color(Text))
	t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(lipgloss.Color(Text))
	t.Focused.FocusedButton = t.Focused.FocusedButton.Foreground(lipgloss.Color(Base)).Background(lipgloss.Color(Blue))
	t.Focused.BlurredButton = t.Focused.BlurredButton.Foreground(lipgloss.Color(Text)).Background(lipgloss.Color(Base))

	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(lipgloss.Color(Teal))
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(lipgloss.Color(Overlay0))
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(lipgloss.Color(Blue))

	t.Blurred = t.Focused
	t.Blurred.Base = t.Blurred.Base.BorderStyle(lipgloss.HiddenBorder())

	t.Help.Ellipsis = t.Help.Ellipsis.Foreground(lipgloss.Color(SubText0))
	t.Help.ShortKey = t.Help.ShortKey.Foreground(lipgloss.Color(SubText0))
	t.Help.ShortDesc = t.Help.ShortDesc.Foreground(lipgloss.Color(Ovelay1))
	t.Help.ShortSeparator = t.Help.ShortSeparator.Foreground(lipgloss.Color(SubText0))
	t.Help.FullKey = t.Help.FullKey.Foreground(lipgloss.Color(SubText0))
	t.Help.FullDesc = t.Help.FullDesc.Foreground(lipgloss.Color(Ovelay1))
	t.Help.FullSeparator = t.Help.FullSeparator.Foreground(lipgloss.Color(SubText0))

	return t
}
