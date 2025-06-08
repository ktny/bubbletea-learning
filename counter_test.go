package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestCounterModel(t *testing.T) {
	t.Run("初期状態", func(t *testing.T) {
		m := NewCounterModel()
		if m.count != 0 {
			t.Errorf("初期カウントは0であるべき、実際: %d", m.count)
		}
	})

	t.Run("Initメソッド", func(t *testing.T) {
		m := NewCounterModel()
		cmd := m.Init()
		if cmd != nil {
			t.Error("Initはnilを返すべき")
		}
	})

	t.Run("↑キーでカウンター増加", func(t *testing.T) {
		m := NewCounterModel()
		msg := tea.KeyMsg{Type: tea.KeyUp}

		newModel, cmd := m.Update(msg)
		updatedModel := newModel.(counterModel)

		if updatedModel.count != 1 {
			t.Errorf("↑キー後のカウントは1であるべき、実際: %d", updatedModel.count)
		}
		if cmd != nil {
			t.Error("コマンドはnilであるべき")
		}
	})

	t.Run("↓キーでカウンター減少", func(t *testing.T) {
		m := NewCounterModel()
		m.count = 5
		msg := tea.KeyMsg{Type: tea.KeyDown}

		newModel, cmd := m.Update(msg)
		updatedModel := newModel.(counterModel)

		if updatedModel.count != 4 {
			t.Errorf("↓キー後のカウントは4であるべき、実際: %d", updatedModel.count)
		}
		if cmd != nil {
			t.Error("コマンドはnilであるべき")
		}
	})

	t.Run("スペースキーでリセット", func(t *testing.T) {
		m := NewCounterModel()
		m.count = 10
		msg := tea.KeyMsg{Type: tea.KeySpace}

		newModel, cmd := m.Update(msg)
		updatedModel := newModel.(counterModel)

		if updatedModel.count != 0 {
			t.Errorf("スペースキー後のカウントは0であるべき、実際: %d", updatedModel.count)
		}
		if cmd != nil {
			t.Error("コマンドはnilであるべき")
		}
	})

	t.Run("qキーで終了", func(t *testing.T) {
		m := NewCounterModel()
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}

		_, cmd := m.Update(msg)

		if cmd == nil {
			t.Error("qキーはQuitコマンドを返すべき")
		}
	})

	t.Run("Ctrl+Cで終了", func(t *testing.T) {
		m := NewCounterModel()
		msg := tea.KeyMsg{Type: tea.KeyCtrlC}

		_, cmd := m.Update(msg)

		if cmd == nil {
			t.Error("Ctrl+CはQuitコマンドを返すべき")
		}
	})

	t.Run("その他のキーは無視", func(t *testing.T) {
		m := NewCounterModel()
		initialCount := m.count
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}

		newModel, cmd := m.Update(msg)
		updatedModel := newModel.(counterModel)

		if updatedModel.count != initialCount {
			t.Errorf("無関係なキーではカウントは変更されないべき")
		}
		if cmd != nil {
			t.Error("コマンドはnilであるべき")
		}
	})

	t.Run("View表示確認", func(t *testing.T) {
		m := NewCounterModel()
		m.count = 42
		view := m.View()

		if view == "" {
			t.Error("ビューは空であってはいけない")
		}
		// ビューにカウント値が含まれているか確認
		expectedContent := "42"
		if !contains(view, expectedContent) {
			t.Errorf("ビューにカウント値 %s が含まれているべき", expectedContent)
		}
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr || contains(s[1:], substr)
}
