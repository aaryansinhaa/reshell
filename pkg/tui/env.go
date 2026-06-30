package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type EnvComponent struct{}

func (e EnvComponent) View(m model) string {
	availWidth := m.width - 30
	if availWidth < 45 {
		availWidth = 45
	}

	start, end := m.getVisibleSlice(len(m.envData))
	var envList []string
	if start > 0 {
		envList = append(envList, TextMuted.Render("  ▲ ..."))
	}
	for i := start; i < end; i++ {
		v := m.envData[i]
		status := SuccessLabel.Render("on")
		if !v.Enabled {
			status = TextMuted.Render("off")
		}
		rawLine := fmt.Sprintf("export %s=%s (%s)", v.Name, v.Value, status)
		displayName := truncateString(rawLine, availWidth-4)
		if i == m.selectedIdx && m.activeTab == TabEnv {
			envList = append(envList, SelectedStyle.Render("> "+displayName))
		} else {
			envList = append(envList, UnselectedStyle.Render("  "+displayName))
		}
	}
	if end < len(m.envData) {
		envList = append(envList, TextMuted.Render("  ▼ ..."))
	}

	return lipgloss.JoinVertical(lipgloss.Left, envList...)
}
