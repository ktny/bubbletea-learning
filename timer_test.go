package main

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func TestTimerModel(t *testing.T) {
	t.Run("初期状態", func(t *testing.T) {
		m := NewTimerModel()
		if m.state != stopped {
			t.Errorf("初期状態はstopped であるべき、実際: %v", m.state)
		}
		if m.duration != 0 {
			t.Errorf("初期時間は0であるべき、実際: %v", m.duration)
		}
	})

	t.Run("Initメソッド", func(t *testing.T) {
		m := NewTimerModel()
		cmd := m.Init()
		if cmd != nil {
			t.Error("Initはnilを返すべき")
		}
	})

	t.Run("スペースキーでスタート", func(t *testing.T) {
		m := NewTimerModel()
		msg := tea.KeyMsg{Type: tea.KeySpace}

		newModel, cmd := m.Update(msg)
		updatedModel := newModel.(timerModel)

		if updatedModel.state != running {
			t.Errorf("スペースキー後の状態はrunningであるべき、実際: %v", updatedModel.state)
		}
		if cmd == nil {
			t.Error("tickコマンドが返されるべき")
		}
	})

	t.Run("sキーでスタート", func(t *testing.T) {
		m := NewTimerModel()
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}

		newModel, cmd := m.Update(msg)
		updatedModel := newModel.(timerModel)

		if updatedModel.state != running {
			t.Errorf("sキー後の状態はrunningであるべき、実際: %v", updatedModel.state)
		}
		if cmd == nil {
			t.Error("tickコマンドが返されるべき")
		}
	})

	t.Run("実行中にスペースキーで一時停止", func(t *testing.T) {
		m := NewTimerModel()
		m.state = running
		m.startTime = time.Now()
		msg := tea.KeyMsg{Type: tea.KeySpace}

		newModel, cmd := m.Update(msg)
		updatedModel := newModel.(timerModel)

		if updatedModel.state != paused {
			t.Errorf("実行中にスペースキー後の状態はpausedであるべき、実際: %v", updatedModel.state)
		}
		if cmd != nil {
			t.Error("一時停止時はコマンドはnilであるべき")
		}
	})

	t.Run("一時停止中にスペースキーで再開", func(t *testing.T) {
		m := NewTimerModel()
		m.state = paused
		m.pausedTime = time.Second * 5
		msg := tea.KeyMsg{Type: tea.KeySpace}

		newModel, cmd := m.Update(msg)
		updatedModel := newModel.(timerModel)

		if updatedModel.state != running {
			t.Errorf("一時停止中にスペースキー後の状態はrunningであるべき、実際: %v", updatedModel.state)
		}
		if cmd == nil {
			t.Error("tickコマンドが返されるべき")
		}
	})

	t.Run("rキーでリセット", func(t *testing.T) {
		m := NewTimerModel()
		m.state = running
		m.duration = time.Minute * 2
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}

		newModel, cmd := m.Update(msg)
		updatedModel := newModel.(timerModel)

		if updatedModel.state != stopped {
			t.Errorf("rキー後の状態はstoppedであるべき、実際: %v", updatedModel.state)
		}
		if updatedModel.duration != 0 {
			t.Errorf("rキー後の時間は0であるべき、実際: %v", updatedModel.duration)
		}
		if cmd != nil {
			t.Error("リセット時はコマンドはnilであるべき")
		}
	})

	t.Run("qキーで終了", func(t *testing.T) {
		m := NewTimerModel()
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}

		_, cmd := m.Update(msg)

		if cmd == nil {
			t.Error("qキーはQuitコマンドを返すべき")
		}
	})

	t.Run("Ctrl+Cで終了", func(t *testing.T) {
		m := NewTimerModel()
		msg := tea.KeyMsg{Type: tea.KeyCtrlC}

		_, cmd := m.Update(msg)

		if cmd == nil {
			t.Error("Ctrl+CはQuitコマンドを返すべき")
		}
	})

	t.Run("tickMsg処理（実行中）", func(t *testing.T) {
		m := NewTimerModel()
		m.state = running
		m.startTime = time.Now().Add(-time.Second * 3) // 3秒前に開始
		msg := tickMsg(time.Now())

		newModel, cmd := m.Update(msg)
		updatedModel := newModel.(timerModel)

		if updatedModel.duration < time.Second*2 {
			t.Errorf("3秒前に開始したので、duration は少なくとも2秒以上であるべき、実際: %v", updatedModel.duration)
		}
		if cmd == nil {
			t.Error("次のtickコマンドが返されるべき")
		}
	})

	t.Run("tickMsg処理（停止中）", func(t *testing.T) {
		m := NewTimerModel()
		m.state = stopped
		initialDuration := m.duration
		msg := tickMsg(time.Now())

		newModel, cmd := m.Update(msg)
		updatedModel := newModel.(timerModel)

		if updatedModel.duration != initialDuration {
			t.Errorf("停止中はdurationは変更されないべき、初期: %v, 実際: %v", initialDuration, updatedModel.duration)
		}
		if cmd != nil {
			t.Error("停止中はコマンドはnilであるべき")
		}
	})

	t.Run("tickMsg処理（一時停止中）", func(t *testing.T) {
		m := NewTimerModel()
		m.state = paused
		initialDuration := m.duration
		msg := tickMsg(time.Now())

		newModel, cmd := m.Update(msg)
		updatedModel := newModel.(timerModel)

		if updatedModel.duration != initialDuration {
			t.Errorf("一時停止中はdurationは変更されないべき、初期: %v, 実際: %v", initialDuration, updatedModel.duration)
		}
		if cmd != nil {
			t.Error("一時停止中はコマンドはnilであるべき")
		}
	})

	t.Run("その他のキーは無視", func(t *testing.T) {
		m := NewTimerModel()
		initialState := m.state
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}

		newModel, cmd := m.Update(msg)
		updatedModel := newModel.(timerModel)

		if updatedModel.state != initialState {
			t.Errorf("無関係なキーでは状態は変更されないべき")
		}
		if cmd != nil {
			t.Error("コマンドはnilであるべき")
		}
	})
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{"0秒", 0, "00:00"},
		{"30秒", time.Second * 30, "00:30"},
		{"1分", time.Minute, "01:00"},
		{"1分30秒", time.Minute + time.Second*30, "01:30"},
		{"10分5秒", time.Minute*10 + time.Second*5, "10:05"},
		{"59分59秒", time.Minute*59 + time.Second*59, "59:59"},
		{"1時間", time.Hour, "60:00"},
		{"1時間30分", time.Hour + time.Minute*30, "90:00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDuration(tt.duration)
			if result != tt.expected {
				t.Errorf("formatDuration(%v) = %s, 期待値: %s", tt.duration, result, tt.expected)
			}
		})
	}
}

func TestTimerModelView(t *testing.T) {
	t.Run("View表示確認", func(t *testing.T) {
		m := NewTimerModel()
		m.duration = time.Minute*5 + time.Second*30
		view := m.View()

		if view == "" {
			t.Error("ビューは空であってはいけない")
		}

		// ビューに時間が正しくフォーマットされて含まれているか確認
		expectedTime := "05:30"
		if !contains(view, expectedTime) {
			t.Errorf("ビューに時間 %s が含まれているべき", expectedTime)
		}
	})
}