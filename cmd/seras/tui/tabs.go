package main

import (
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type Tab struct {
	Title   string
	Content string
}

type TabGroup struct {
	tabs      []*Tab
	activeTab int
}

func (tg *TabGroup) ActiveContent() string {
	return tg.ActiveTab().Content
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
		isFirst, isActive := i == 0, i == tg.activeTab
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
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t.Title))
	}

	rendered := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	gap := tabGap.Render(strings.Repeat(" ", max(0, getTermWidth()-lipgloss.Width(rendered)-1)))
	gap += lipgloss.NewStyle().Foreground(highlightColor).Render("╮")

	return lipgloss.JoinHorizontal(lipgloss.Bottom, rendered, gap)
}

// Styling

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right

	return border
}

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	highlightColor    = lipgloss.AdaptiveColor{Light: "#E53935", Dark: "#F07178"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Copy().Border(activeTabBorder, true)
	tabGap            = inactiveTabStyle.Copy().
				BorderTop(false).
				BorderLeft(false).
				BorderRight(false)
)
