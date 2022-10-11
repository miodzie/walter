package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
	"os"
	"strings"
)

type model struct {
	tabs     *TabGroup
	textarea textarea.Model
	viewport viewport.Model
	messages []string
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, tea.EnterAltScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)
	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyRight, tea.KeyCtrlL, tea.KeyTab:
			m.tabs.NextTab()
			return m, nil
		case tea.KeyLeft, tea.KeyCtrlH, tea.KeyShiftTab:
			m.tabs.PreviousTab()
			return m, nil
		case tea.KeyEnter:
			m.messages = append(m.messages, lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Render("You: ")+m.textarea.Value())
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

var (
	windowStyle = lipgloss.NewStyle().BorderForeground(highlightColor).
			Padding(0, 0, 0, 1).
		//Align(lipgloss.Left).
		Border(lipgloss.NormalBorder()).UnsetBorderTop()
	//docStyle    = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	docStyle = lipgloss.NewStyle()
)

func (m model) View() string {
	doc := strings.Builder{}
	// Tabs
	tabs := m.tabs.Render()
	doc.WriteString(tabs)

	// Body
	doc.WriteString("\n")
	doc.WriteString(windowStyle.Width(getTermWidth()).Render(m.viewport.View()))

	// Input
	doc.WriteString("\n")
	doc.WriteString(m.textarea.View())

	return docStyle.Render(doc.String())
}

func getTermWidth() int {
	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	width := physicalWidth - windowStyle.GetHorizontalFrameSize() - 2
	return width
}

func main() {
	tabs := []*Tab{
		{"Logs", "A bunch of log content I guess"},
		{"Discord", "Discord sucks"},
		{"IRC", "IRC sucks too"},
	}
	tabGroup := &TabGroup{tabs: tabs}
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()
	ta.Prompt = "┃ "
	ta.CharLimit = 280
	ta.SetWidth(getTermWidth())
	ta.SetHeight(1)
	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	vp := viewport.New(100, 5)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	m := model{tabs: tabGroup, textarea: ta, viewport: vp}

	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
