package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestFormModel(t *testing.T) {
	t.Run("初期状態", func(t *testing.T) {
		m := NewFormModel()
		if m.focusIndex != nameInput {
			t.Errorf("初期フォーカスは名前フィールドであるべき、実際: %d", m.focusIndex)
		}
		if m.state != formInput {
			t.Error("初期状態はformInputであるべき")
		}
		if len(m.inputs) != 2 {
			t.Errorf("入力フィールドは2つであるべき、実際: %d", len(m.inputs))
		}
		if !m.inputs[nameInput].Focused() {
			t.Error("名前フィールドにフォーカスがあるべき")
		}
		if m.inputs[emailInput].Focused() {
			t.Error("メールフィールドにフォーカスがないべき")
		}
	})

	t.Run("Initメソッド", func(t *testing.T) {
		m := NewFormModel()
		cmd := m.Init()
		if cmd == nil {
			t.Error("Initはtextinput.Blinkコマンドを返すべき")
		}
	})

	t.Run("Tabキーで次のフィールドへ移動", func(t *testing.T) {
		m := NewFormModel()
		msg := tea.KeyMsg{Type: tea.KeyTab}

		// 名前 → メール
		newModel, _ := m.Update(msg)
		updatedModel := newModel.(formModel)
		if updatedModel.focusIndex != emailInput {
			t.Errorf("Tabキー後のフォーカスはメールフィールドであるべき、実際: %d", updatedModel.focusIndex)
		}
		if !updatedModel.inputs[emailInput].Focused() {
			t.Error("メールフィールドにフォーカスがあるべき")
		}

		// メール → 送信ボタン
		newModel, _ = updatedModel.Update(msg)
		updatedModel = newModel.(formModel)
		if updatedModel.focusIndex != submitButton {
			t.Errorf("Tabキー後のフォーカスは送信ボタンであるべき、実際: %d", updatedModel.focusIndex)
		}

		// 送信ボタン → 名前（ループ）
		newModel, _ = updatedModel.Update(msg)
		updatedModel = newModel.(formModel)
		if updatedModel.focusIndex != nameInput {
			t.Errorf("Tabキー後のフォーカスは名前フィールドに戻るべき、実際: %d", updatedModel.focusIndex)
		}
	})

	t.Run("Shift+Tabキーで前のフィールドへ移動", func(t *testing.T) {
		m := NewFormModel()
		m.focusIndex = emailInput
		msg := tea.KeyMsg{Type: tea.KeyShiftTab}

		// メール → 名前
		newModel, _ := m.Update(msg)
		updatedModel := newModel.(formModel)
		if updatedModel.focusIndex != nameInput {
			t.Errorf("Shift+Tab後のフォーカスは名前フィールドであるべき、実際: %d", updatedModel.focusIndex)
		}

		// 名前 → 送信ボタン（ループ）
		newModel, _ = updatedModel.Update(msg)
		updatedModel = newModel.(formModel)
		if updatedModel.focusIndex != submitButton {
			t.Errorf("Shift+Tab後のフォーカスは送信ボタンであるべき、実際: %d", updatedModel.focusIndex)
		}
	})

	t.Run("↓キーで次のフィールドへ移動", func(t *testing.T) {
		m := NewFormModel()
		msg := tea.KeyMsg{Type: tea.KeyDown}

		newModel, _ := m.Update(msg)
		updatedModel := newModel.(formModel)
		if updatedModel.focusIndex != emailInput {
			t.Errorf("↓キー後のフォーカスはメールフィールドであるべき、実際: %d", updatedModel.focusIndex)
		}
	})

	t.Run("↑キーで前のフィールドへ移動", func(t *testing.T) {
		m := NewFormModel()
		m.focusIndex = emailInput
		msg := tea.KeyMsg{Type: tea.KeyUp}

		newModel, _ := m.Update(msg)
		updatedModel := newModel.(formModel)
		if updatedModel.focusIndex != nameInput {
			t.Errorf("↑キー後のフォーカスは名前フィールドであるべき、実際: %d", updatedModel.focusIndex)
		}
	})

	t.Run("バリデーション - 名前が空", func(t *testing.T) {
		m := NewFormModel()
		// 名前フィールドは空のまま
		m.inputs[emailInput].SetValue("test@example.com")

		err := m.validate()
		if err == nil {
			t.Error("名前が空の場合はエラーが返されるべき")
		}
		if !strings.Contains(err.Error(), "名前") {
			t.Error("エラーメッセージに「名前」が含まれるべき")
		}
	})

	t.Run("バリデーション - メールが空", func(t *testing.T) {
		m := NewFormModel()
		m.inputs[nameInput].SetValue("山田太郎")
		// メールフィールドは空のまま

		err := m.validate()
		if err == nil {
			t.Error("メールが空の場合はエラーが返されるべき")
		}
		if !strings.Contains(err.Error(), "メールアドレス") {
			t.Error("エラーメッセージに「メールアドレス」が含まれるべき")
		}
	})

	t.Run("バリデーション - 不正なメール形式", func(t *testing.T) {
		m := NewFormModel()
		m.inputs[nameInput].SetValue("山田太郎")
		
		// 不正なメール形式のテストケース
		invalidEmails := []string{
			"invalid",
			"@example.com",
			"test@",
			"test@example",
			"test.example.com",
		}

		for _, email := range invalidEmails {
			m.inputs[emailInput].SetValue(email)
			err := m.validate()
			if err == nil {
				t.Errorf("メール形式が不正な場合はエラーが返されるべき: %s", email)
			}
			if !strings.Contains(err.Error(), "形式") {
				t.Error("エラーメッセージに「形式」が含まれるべき")
			}
		}
	})

	t.Run("バリデーション - 正常な入力", func(t *testing.T) {
		m := NewFormModel()
		m.inputs[nameInput].SetValue("山田太郎")
		
		// 正常なメール形式のテストケース
		validEmails := []string{
			"test@example.com",
			"user.name@example.com",
			"user+tag@example.co.jp",
			"test123@test-domain.com",
		}

		for _, email := range validEmails {
			m.inputs[emailInput].SetValue(email)
			err := m.validate()
			if err != nil {
				t.Errorf("正常な入力でエラーが返されるべきでない: %s, エラー: %v", email, err)
			}
		}
	})

	t.Run("送信処理 - バリデーションエラー", func(t *testing.T) {
		m := NewFormModel()
		m.focusIndex = submitButton
		// 入力なしで送信
		msg := tea.KeyMsg{Type: tea.KeyEnter}

		newModel, _ := m.Update(msg)
		updatedModel := newModel.(formModel)

		if updatedModel.state != formInput {
			t.Error("バリデーションエラー時は入力状態のままであるべき")
		}
		if updatedModel.errorMessage == "" {
			t.Error("エラーメッセージが設定されるべき")
		}
		if updatedModel.submitted {
			t.Error("送信済みフラグは立たないべき")
		}
	})

	t.Run("送信処理 - 正常", func(t *testing.T) {
		m := NewFormModel()
		m.inputs[nameInput].SetValue("山田太郎")
		m.inputs[emailInput].SetValue("taro@example.com")
		m.focusIndex = submitButton
		msg := tea.KeyMsg{Type: tea.KeyEnter}

		newModel, _ := m.Update(msg)
		updatedModel := newModel.(formModel)

		if updatedModel.state != formSubmitted {
			t.Error("送信後はformSubmitted状態になるべき")
		}
		if !updatedModel.submitted {
			t.Error("送信済みフラグが立つべき")
		}
		if updatedModel.errorMessage != "" {
			t.Error("エラーメッセージは空であるべき")
		}
	})

	t.Run("フォームデータ取得", func(t *testing.T) {
		m := NewFormModel()
		m.inputs[nameInput].SetValue("  山田太郎  ")  // 前後に空白
		m.inputs[emailInput].SetValue("  taro@example.com  ")

		data := m.getFormData()
		if data.name != "山田太郎" {
			t.Errorf("名前の空白がトリムされるべき、実際: '%s'", data.name)
		}
		if data.email != "taro@example.com" {
			t.Errorf("メールの空白がトリムされるべき、実際: '%s'", data.email)
		}
	})

	t.Run("Enterキーでフィールド移動", func(t *testing.T) {
		m := NewFormModel()
		msg := tea.KeyMsg{Type: tea.KeyEnter}

		// 名前フィールドでEnter
		newModel, _ := m.Update(msg)
		updatedModel := newModel.(formModel)
		if updatedModel.focusIndex != emailInput {
			t.Error("名前フィールドでEnterを押すとメールフィールドに移動すべき")
		}
	})

	t.Run("Escキーで終了", func(t *testing.T) {
		m := NewFormModel()
		msg := tea.KeyMsg{Type: tea.KeyEsc}

		_, cmd := m.Update(msg)
		if cmd == nil {
			t.Error("EscキーはQuitコマンドを返すべき")
		}
	})

	t.Run("Ctrl+Cで終了", func(t *testing.T) {
		m := NewFormModel()
		msg := tea.KeyMsg{Type: tea.KeyCtrlC}

		_, cmd := m.Update(msg)
		if cmd == nil {
			t.Error("Ctrl+CはQuitコマンドを返すべき")
		}
	})

	t.Run("送信完了後のqキーで終了", func(t *testing.T) {
		m := NewFormModel()
		m.state = formSubmitted
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}

		_, cmd := m.Update(msg)
		if cmd == nil {
			t.Error("送信完了後のqキーはQuitコマンドを返すべき")
		}
	})
}

func TestFormModelView(t *testing.T) {
	t.Run("入力画面の表示", func(t *testing.T) {
		m := NewFormModel()
		view := m.View()

		if view == "" {
			t.Error("ビューは空であってはいけない")
		}

		// 基本要素の確認
		requiredElements := []string{
			"ユーザー登録フォーム",
			"名前:",
			"メールアドレス:",
			"送信",
			"Tab:",
		}

		for _, element := range requiredElements {
			if !strings.Contains(view, element) {
				t.Errorf("ビューに「%s」が含まれているべき", element)
			}
		}
	})

	t.Run("エラーメッセージの表示", func(t *testing.T) {
		m := NewFormModel()
		m.errorMessage = "テストエラー"
		view := m.View()

		if !strings.Contains(view, "❌") {
			t.Error("エラーアイコンが表示されるべき")
		}
		if !strings.Contains(view, "テストエラー") {
			t.Error("エラーメッセージが表示されるべき")
		}
	})

	t.Run("送信完了画面の表示", func(t *testing.T) {
		m := NewFormModel()
		m.inputs[nameInput].SetValue("山田太郎")
		m.inputs[emailInput].SetValue("taro@example.com")
		m.state = formSubmitted
		view := m.View()

		requiredElements := []string{
			"フォーム送信完了",
			"✅",
			"正常に送信されました",
			"山田太郎",
			"taro@example.com",
			"q: 終了",
		}

		for _, element := range requiredElements {
			if !strings.Contains(view, element) {
				t.Errorf("送信完了画面に「%s」が含まれているべき", element)
			}
		}
	})
}