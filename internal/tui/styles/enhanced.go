package styles

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/omnitrix-sh/cli/internal/tui/theme"
)

// Enhanced styling utilities for better UI components

// Card returns a styled card component with rounded borders and shadows
func Card() lipgloss.Style {
	t := theme.CurrentTheme()
	return lipgloss.NewStyle().
		Background(t.Background()).
		Foreground(t.Text()).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.BorderNormal()).
		Padding(1, 2)
}

// Badge returns a styled badge component
func Badge() lipgloss.Style {
	t := theme.CurrentTheme()
	return lipgloss.NewStyle().
		Background(t.Primary()).
		Foreground(t.Background()).
		Padding(0, 1).
		Bold(true).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Primary())
}

// StatusBadge returns a styled status badge
func StatusBadge(statusType string) lipgloss.Style {
	t := theme.CurrentTheme()
	base := lipgloss.NewStyle().
		Padding(0, 1).
		Bold(true).
		Border(lipgloss.RoundedBorder())
	
	switch statusType {
	case "success":
		return base.
			Background(t.Success()).
			Foreground(t.Background()).
			BorderForeground(t.Success())
	case "warning":
		return base.
			Background(t.Warning()).
			Foreground(t.Background()).
			BorderForeground(t.Warning())
	case "error":
		return base.
			Background(t.Error()).
			Foreground(t.Background()).
			BorderForeground(t.Error())
	case "info":
		return base.
			Background(t.Info()).
			Foreground(t.Background()).
			BorderForeground(t.Info())
	default:
		return base.
			Background(t.TextMuted()).
			Foreground(t.Background()).
			BorderForeground(t.TextMuted())
	}
}

// ProgressBar creates a simple text-based progress bar
func ProgressBar(width int, progress float64) string {
	filled := int(progress * float64(width))
	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	return bar
}

// Tooltip returns a styled tooltip component
func Tooltip() lipgloss.Style {
	t := theme.CurrentTheme()
	return lipgloss.NewStyle().
		Background(t.BackgroundDarker()).
		Foreground(t.TextEmphasized()).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.BorderDim())
}

// InputField returns a styled input field
func InputField(focused bool) lipgloss.Style {
	t := theme.CurrentTheme()
	base := lipgloss.NewStyle().
		Background(t.Background()).
		Foreground(t.Text()).
		Padding(0, 1).
		Border(lipgloss.NormalBorder()).
		BorderForeground(t.BorderNormal())
	
	if focused {
		base = base.
			BorderForeground(t.BorderFocused()).
			Background(t.BackgroundSecondary())
	}
	
	return base
}

// MessageBubble returns a styled message bubble
func MessageBubble(isUser bool) lipgloss.Style {
	t := theme.CurrentTheme()
	base := lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		MarginBottom(1)
	
	if isUser {
		return base.
			Background(t.BackgroundSecondary()).
			Foreground(t.Text()).
			BorderForeground(t.Secondary()).
			Align(lipgloss.Right)
	}
	
	return base.
		Background(t.Background()).
		Foreground(t.Text()).
		BorderForeground(t.Primary()).
		Align(lipgloss.Left)
}

// Separator returns a styled separator line
func Separator(width int) string {
	t := theme.CurrentTheme()
	return lipgloss.NewStyle().
		Foreground(t.BorderDim()).
		Render(lipgloss.Place(width, 1, lipgloss.Center, lipgloss.Center, "─"))
}

// LoadingSpinner returns animated loading spinner states
func LoadingSpinner(frame int) string {
	spinners := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	return spinners[frame%len(spinners)]
}

// GradientText simulates gradient text using different shades
func GradientText(text string, startColor, endColor lipgloss.AdaptiveColor) []string {
	// Simple gradient simulation - in a real implementation you'd interpolate colors
	words := strings.Fields(text)
	result := make([]string, len(words))
	
	for i, word := range words {
		if i%2 == 0 {
			result[i] = lipgloss.NewStyle().Foreground(startColor).Render(word)
		} else {
			result[i] = lipgloss.NewStyle().Foreground(endColor).Render(word)
		}
	}
	
	return result
}
