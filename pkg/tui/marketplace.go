package tui



type MarketplaceComponent struct{}

func (mp MarketplaceComponent) View(m model) string {
	preview := `Browse & install reshell terminal profiles.
Specify repositories to automatically import setup configurations (aliases, variables, functions, and package lists).

Example repositories:
- github.com/aaryansinhaa/reshell-java (Setup JDK, Maven, workspace env-vars, aliases)
- github.com/aaryansinhaa/reshell-react (Setup TypeScript environment, formatters, and templates)

Press 'i' key to install a package.`

	return CardStyle.Width(m.width - 30).Render(preview)
}
