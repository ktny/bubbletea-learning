package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ktny/bubbletea-learning/pkg/common"
	"github.com/ktny/bubbletea-learning/pkg/constants"
	"github.com/ktny/bubbletea-learning/pkg/styles"
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
		// Handle common quit keys
		if cmd := common.HandleQuitKeys(msg); cmd != nil {
			return m, cmd
		}
		
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
		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "s": // Alternative space for reset
				m.count = 0
				return m, nil
			}
		}
	}
	return m, nil
}

func (m counterModel) View() string {
	// カウンター数値のスタイル（条件付き）
	var countStyle = styles.CounterZeroStyle
	if m.count > 0 {
		countStyle = styles.CounterPositiveStyle
	} else if m.count < 0 {
		countStyle = styles.CounterNegativeStyle
	}

	content := fmt.Sprintf(
		"%s\n\nカウンター: %s\n\n%s",
		styles.TitleStyle.Render(constants.CounterTitle),
		countStyle.Render(fmt.Sprintf("%d", m.count)),
		styles.HelpStyle.Render(
			"↑: 増加\n"+
			"↓: 減少\n"+
			"スペース/s: リセット\n"+
			constants.QuitHelp,
		),
	)

	return styles.BorderStyle.Render(content)
}