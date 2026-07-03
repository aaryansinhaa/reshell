package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type ProfilesComponent struct{}

func (p ProfilesComponent) View(m model) string {
	availWidth := m.width - 30
	if availWidth < 45 {
		availWidth = 45
	}

	title := TitleStyle.Render("Configuration Profiles")
	desc := "Isolate workspaces into independent profiles (aliases, variables, functions, scripts)."

	start, end := m.getVisibleSlice(len(m.profilesData))
	var list []string

	if start > 0 {
		list = append(list, TextMuted.Render("  ▲ ..."))
	}
	for i := start; i < end; i++ {
		profileName := m.profilesData[i]
		isActive := profileName == m.activeProfile

		status := ""
		if isActive {
			status = SuccessLabel.Render(" (active)")
		}

		rawLine := fmt.Sprintf("%s%s", profileName, status)
		displayName := truncateString(rawLine, availWidth-4)
		if i == m.selectedIdx && m.activeTab == TabProfiles {
			list = append(list, SelectedStyle.Render("> "+displayName))
		} else {
			list = append(list, UnselectedStyle.Render("  "+displayName))
		}
	}
	if end < len(m.profilesData) {
		list = append(list, TextMuted.Render("  ▼ ..."))
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		TextMuted.Render(desc),
		"",
		lipgloss.JoinVertical(lipgloss.Left, list...),
	)

	return CardStyle.Width(availWidth).Render(content)
}
