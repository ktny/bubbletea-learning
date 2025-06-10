package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestTodoModel(t *testing.T) {
	t.Run("初期状態", func(t *testing.T) {
		m := NewTodoModel()
		if m.cursor != 0 {
			t.Errorf("初期カーソル位置は0であるべき、実際: %d", m.cursor)
		}
		if len(m.items) == 0 {
			t.Error("初期アイテムが存在するべき")
		}
		if m.viewport != 0 {
			t.Errorf("初期ビューポートは0であるべき、実際: %d", m.viewport)
		}
	})

	t.Run("Initメソッド", func(t *testing.T) {
		m := NewTodoModel()
		cmd := m.Init()
		if cmd != nil {
			t.Error("Initはnilを返すべき")
		}
	})

	t.Run("↓キーでカーソル移動", func(t *testing.T) {
		m := NewTodoModel()
		msg := tea.KeyMsg{Type: tea.KeyDown}

		newModel, _ := m.Update(msg)
		updatedModel := newModel.(todoModel)

		if updatedModel.cursor != 1 {
			t.Errorf("↓キー後のカーソルは1であるべき、実際: %d", updatedModel.cursor)
		}
	})

	t.Run("↑キーでカーソル移動", func(t *testing.T) {
		m := NewTodoModel()
		m.cursor = 2
		msg := tea.KeyMsg{Type: tea.KeyUp}

		newModel, _ := m.Update(msg)
		updatedModel := newModel.(todoModel)

		if updatedModel.cursor != 1 {
			t.Errorf("↑キー後のカーソルは1であるべき、実際: %d", updatedModel.cursor)
		}
	})

	t.Run("jキーで下へ移動", func(t *testing.T) {
		m := NewTodoModel()
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}

		newModel, _ := m.Update(msg)
		updatedModel := newModel.(todoModel)

		if updatedModel.cursor != 1 {
			t.Errorf("jキー後のカーソルは1であるべき、実際: %d", updatedModel.cursor)
		}
	})

	t.Run("kキーで上へ移動", func(t *testing.T) {
		m := NewTodoModel()
		m.cursor = 2
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}

		newModel, _ := m.Update(msg)
		updatedModel := newModel.(todoModel)

		if updatedModel.cursor != 1 {
			t.Errorf("kキー後のカーソルは1であるべき、実際: %d", updatedModel.cursor)
		}
	})

	t.Run("カーソル移動の境界チェック（上限）", func(t *testing.T) {
		m := NewTodoModel()
		m.cursor = 0
		msg := tea.KeyMsg{Type: tea.KeyUp}

		newModel, _ := m.Update(msg)
		updatedModel := newModel.(todoModel)

		if updatedModel.cursor != 0 {
			t.Errorf("最上部でのカーソルは0のままであるべき、実際: %d", updatedModel.cursor)
		}
	})

	t.Run("カーソル移動の境界チェック（下限）", func(t *testing.T) {
		m := NewTodoModel()
		m.cursor = len(m.items) - 1
		msg := tea.KeyMsg{Type: tea.KeyDown}

		newModel, _ := m.Update(msg)
		updatedModel := newModel.(todoModel)

		if updatedModel.cursor != len(m.items)-1 {
			t.Errorf("最下部でのカーソルは変わらないべき、実際: %d", updatedModel.cursor)
		}
	})

	t.Run("Enterキーで項目の完了状態を切り替え", func(t *testing.T) {
		m := NewTodoModel()
		initialState := m.items[0].completed
		msg := tea.KeyMsg{Type: tea.KeyEnter}

		newModel, _ := m.Update(msg)
		updatedModel := newModel.(todoModel)

		if updatedModel.items[0].completed == initialState {
			t.Error("Enterキーで完了状態が切り替わるべき")
		}
	})

	t.Run("スペースキーで項目の完了状態を切り替え", func(t *testing.T) {
		m := NewTodoModel()
		initialState := m.items[0].completed
		msg := tea.KeyMsg{Type: tea.KeySpace}

		newModel, _ := m.Update(msg)
		updatedModel := newModel.(todoModel)

		if updatedModel.items[0].completed == initialState {
			t.Error("スペースキーで完了状態が切り替わるべき")
		}
	})

	t.Run("qキーで終了", func(t *testing.T) {
		m := NewTodoModel()
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}

		_, cmd := m.Update(msg)

		if cmd == nil {
			t.Error("qキーはQuitコマンドを返すべき")
		}
	})

	t.Run("ビューポートの調整（下方向）", func(t *testing.T) {
		m := NewTodoModel()
		m.height = 3 // 小さいビューポート
		m.cursor = 0
		m.viewport = 0

		// カーソルをビューポート外に移動
		for i := 0; i < 4; i++ {
			msg := tea.KeyMsg{Type: tea.KeyDown}
			newModel, _ := m.Update(msg)
			m = newModel.(todoModel)
		}

		if m.viewport == 0 {
			t.Error("カーソルがビューポート外に移動したらビューポートも調整されるべき")
		}
		if m.cursor < m.viewport || m.cursor >= m.viewport+m.height {
			t.Error("カーソルは常にビューポート内にあるべき")
		}
	})

	t.Run("ビューポートの調整（上方向）", func(t *testing.T) {
		m := NewTodoModel()
		m.height = 3
		m.cursor = 4
		m.viewport = 2

		// カーソルを上に移動
		for i := 0; i < 3; i++ {
			msg := tea.KeyMsg{Type: tea.KeyUp}
			newModel, _ := m.Update(msg)
			m = newModel.(todoModel)
		}

		if m.viewport != m.cursor {
			t.Error("カーソルがビューポート上部より上に移動したらビューポートも調整されるべき")
		}
	})

	t.Run("ウィンドウサイズ変更の処理", func(t *testing.T) {
		m := NewTodoModel()
		msg := tea.WindowSizeMsg{Width: 80, Height: 24}

		newModel, _ := m.Update(msg)
		updatedModel := newModel.(todoModel)

		expectedHeight := 24 - 10
		if updatedModel.height != expectedHeight {
			t.Errorf("ウィンドウサイズ変更後の高さは%dであるべき、実際: %d", expectedHeight, updatedModel.height)
		}
	})
}

func TestTodoModelView(t *testing.T) {
	t.Run("View表示確認", func(t *testing.T) {
		m := NewTodoModel()
		view := m.View()

		if view == "" {
			t.Error("ビューは空であってはいけない")
		}

		// 基本的な要素が含まれているか確認
		if !contains(view, "TODO") {
			t.Error("ビューにTODOという文字が含まれているべき")
		}
		if !contains(view, "[") && !contains(view, "]") {
			t.Error("ビューにチェックボックスが含まれているべき")
		}
	})

	t.Run("完了項目の表示", func(t *testing.T) {
		m := NewTodoModel()
		// 最初の項目を完了にする
		m.items[0].completed = true
		view := m.View()

		if !contains(view, "✓") {
			t.Error("完了項目にはチェックマークが表示されるべき")
		}
	})

	t.Run("カーソル表示", func(t *testing.T) {
		m := NewTodoModel()
		view := m.View()

		if !contains(view, ">") {
			t.Error("カーソル位置に > が表示されるべき")
		}
	})
}
