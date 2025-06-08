package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
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
	return fmt.Sprintf(
		"カウンター: %d\n\n"+
			"↑: 増加\n"+
			"↓: 減少\n"+
			"スペース: リセット\n"+
			"q または Ctrl+C: 終了\n",
		m.count,
	)
}