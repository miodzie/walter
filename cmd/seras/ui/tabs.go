package main

import "github.com/charmbracelet/lipgloss"

type Tab struct {
	Title   string
	Content string
}

type TabGroup struct {
	tabs      []*Tab
	activeTab int
}

func (tg *TabGroup) ActiveTab() *Tab {
	return tg.tabs[tg.activeTab]
}

func (tg *TabGroup) NextTab() {
	tg.activeTab = min(tg.activeTab+1, len(tg.tabs)-1)
}

func (tg *TabGroup) PreviousTab() {
	tg.activeTab = max(tg.activeTab-1, 0)
}

func (tg *TabGroup) Render() string {
	var renderedTabs []string

	for i, t := range tg.tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(tg.tabs)-1, i == tg.activeTab
		if isActive {
			style = activeTabStyle.Copy()
		} else {
			style = inactiveTabStyle.Copy()
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t.Title))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}
