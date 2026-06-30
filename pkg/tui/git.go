package tui

import (
	"fmt"
)

type GitComponent struct{}

func (g GitComponent) View(m model) string {
	gitContent := "No global Git configurations read."
	if m.gitData != nil {
		gitContent = fmt.Sprintf("Name: %s\nEmail: %s\nGPG Signing: %t\nSigning Key: %s\n\nGlobal Aliases:\n",
			m.gitData.UserName, m.gitData.UserEmail, m.gitData.GpgSign, m.gitData.SigningKey,
		)
		for alias, value := range m.gitData.Aliases {
			gitContent += fmt.Sprintf("  %s = %s\n", alias, value)
		}
	}

	preview := fmt.Sprintf("%s\n%s", TitleStyle.Render("Git global config"), gitContent)
	return CardStyle.Width(m.width - 30).Render(preview)
}
