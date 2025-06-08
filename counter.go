package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type counterModel struct {
	count int
}

func NewCounterModel() counterModel {
	return counterModel{
		count: 0,
	}
}

func (m counterModel) Init() tea.Cmd {
	return nil
}

func (m counterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			m.count++
			return m, nil
		case tea.KeyDown:
			m.count--
			return m, nil
		case tea.KeySpace:
			m.count = 0
			return m, nil
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "q":
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m counterModel) View() string {
	// カウンター数値のスタイル（条件付き）
	var countStyle lipgloss.Style
	if m.count > 0 {
		countStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")) // 緑
	} else if m.count < 0 {
		countStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))  // 赤
	} else {
		countStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))  // 白
	}

	// ヘルプテキストのスタイル
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).  // グレー
		Italic(true)

	// アプリケーション全体の枠線スタイル
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("12")). // 青
		Padding(1, 2)

	content := fmt.Sprintf(
		"カウンター: %s\n\n%s",
		countStyle.Render(fmt.Sprintf("%d", m.count)),
		helpStyle.Render(
			"↑: 増加\n"+
			"↓: 減少\n"+
			"スペース: リセット\n"+
			"q または Ctrl+C: 終了",
		),
	)

	return borderStyle.Render(content)
}