package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestDashboardModel(t *testing.T) {
	t.Run("初期状態", func(t *testing.T) {
		m := NewDashboardModel()

		if len(m.panels) != 4 {
			t.Errorf("パネル数は4つであるべき、実際: %d", len(m.panels))
		}

		if m.activePanel != 0 {
			t.Errorf("初期アクティブパネルは0であるべき、実際: %d", m.activePanel)
		}

		if !m.panels[0].active {
			t.Error("最初のパネルがアクティブであるべき")
		}

		for i := 1; i < len(m.panels); i++ {
			if m.panels[i].active {
				t.Errorf("パネル%dはアクティブでないべき", i)
			}
		}

		if !m.showHelp {
			t.Error("初期状態でヘルプは表示されるべき")
		}
	})

	t.Run("Initメソッド", func(t *testing.T) {
		m := NewDashboardModel()
		cmd := m.Init()

		if cmd == nil {
			t.Error("Initは各パネルの初期化コマンドを返すべき")
		}
	})

	t.Run("Tabキーでパネル切り替え", func(t *testing.T) {
		m := NewDashboardModel()
		msg := tea.KeyMsg{Type: tea.KeyTab}

		// 0 → 1
		newModel, _ := m.Update(msg)
		updatedModel := newModel.(dashboardModel)

		if updatedModel.activePanel != 1 {
			t.Errorf("Tabキー後のアクティブパネルは1であるべき、実際: %d", updatedModel.activePanel)
		}

		if updatedModel.panels[0].active {
			t.Error("パネル0は非アクティブになるべき")
		}

		if !updatedModel.panels[1].active {
			t.Error("パネル1はアクティブになるべき")
		}

		// さらに3回Tabを押して一周
		for i := 0; i < 3; i++ {
			newModel, _ = updatedModel.Update(msg)
			updatedModel = newModel.(dashboardModel)
		}

		if updatedModel.activePanel != 0 {
			t.Errorf("4回Tab後は最初のパネルに戻るべき、実際: %d", updatedModel.activePanel)
		}
	})

	t.Run("Shift+Tabキーで逆方向パネル切り替え", func(t *testing.T) {
		m := NewDashboardModel()
		msg := tea.KeyMsg{Type: tea.KeyShiftTab}

		// 0 → 3 (逆方向)
		newModel, _ := m.Update(msg)
		updatedModel := newModel.(dashboardModel)

		if updatedModel.activePanel != 3 {
			t.Errorf("Shift+Tab後のアクティブパネルは3であるべき、実際: %d", updatedModel.activePanel)
		}

		if !updatedModel.panels[3].active {
			t.Error("パネル3はアクティブになるべき")
		}
	})

	t.Run("数字キーで直接パネル選択", func(t *testing.T) {
		m := NewDashboardModel()

		// '2'キーでパネル1を選択
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}}
		newModel, _ := m.Update(msg)
		updatedModel := newModel.(dashboardModel)

		if updatedModel.activePanel != 1 {
			t.Errorf("'2'キー後のアクティブパネルは1であるべき、実際: %d", updatedModel.activePanel)
		}

		// '4'キーでパネル3を選択
		msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'4'}}
		newModel, _ = updatedModel.Update(msg)
		updatedModel = newModel.(dashboardModel)

		if updatedModel.activePanel != 3 {
			t.Errorf("'4'キー後のアクティブパネルは3であるべき、実際: %d", updatedModel.activePanel)
		}
	})

	t.Run("F1キーでヘルプ表示切り替え", func(t *testing.T) {
		m := NewDashboardModel()
		originalHelp := m.showHelp
		msg := tea.KeyMsg{Type: tea.KeyF1}

		newModel, _ := m.Update(msg)
		updatedModel := newModel.(dashboardModel)

		if updatedModel.showHelp == originalHelp {
			t.Error("F1キーでヘルプ表示が切り替わるべき")
		}

		// もう一度押して元に戻る
		newModel, _ = updatedModel.Update(msg)
		updatedModel = newModel.(dashboardModel)

		if updatedModel.showHelp != originalHelp {
			t.Error("F1キーを2回押すとヘルプ表示が元に戻るべき")
		}
	})

	t.Run("F2キーでグローバルヘルプ切り替え", func(t *testing.T) {
		m := NewDashboardModel()
		originalGlobalHelp := m.globalHelp
		msg := tea.KeyMsg{Type: tea.KeyF2}

		newModel, _ := m.Update(msg)
		updatedModel := newModel.(dashboardModel)

		if updatedModel.globalHelp == originalGlobalHelp {
			t.Error("F2キーでグローバルヘルプが切り替わるべき")
		}
	})

	t.Run("Ctrl+Cで終了", func(t *testing.T) {
		m := NewDashboardModel()
		msg := tea.KeyMsg{Type: tea.KeyCtrlC}

		_, cmd := m.Update(msg)
		if cmd == nil {
			t.Error("Ctrl+CはQuitコマンドを返すべき")
		}
	})

	t.Run("ウィンドウサイズ変更", func(t *testing.T) {
		m := NewDashboardModel()
		msg := tea.WindowSizeMsg{Width: 120, Height: 40}

		newModel, _ := m.Update(msg)
		updatedModel := newModel.(dashboardModel)

		if updatedModel.width != 120 {
			t.Errorf("幅が更新されるべき、実際: %d", updatedModel.width)
		}

		if updatedModel.height != 40 {
			t.Errorf("高さが更新されるべき、実際: %d", updatedModel.height)
		}
	})

	t.Run("グローバルキー判定", func(t *testing.T) {
		m := NewDashboardModel()

		// グローバルキーのテスト
		globalKeys := []tea.KeyMsg{
			{Type: tea.KeyTab},
			{Type: tea.KeyShiftTab},
			{Type: tea.KeyF1},
			{Type: tea.KeyF2},
			{Type: tea.KeyCtrlC},
			{Type: tea.KeyRunes, Runes: []rune{'1'}},
			{Type: tea.KeyRunes, Runes: []rune{'2'}},
			{Type: tea.KeyRunes, Runes: []rune{'3'}},
			{Type: tea.KeyRunes, Runes: []rune{'4'}},
		}

		for _, key := range globalKeys {
			if !m.isGlobalKey(key) {
				t.Errorf("キー %v はグローバルキーとして認識されるべき", key)
			}
		}

		// 非グローバルキーのテスト
		nonGlobalKeys := []tea.KeyMsg{
			{Type: tea.KeyEnter},
			{Type: tea.KeyEsc},
			{Type: tea.KeyUp},
			{Type: tea.KeyDown},
			{Type: tea.KeyRunes, Runes: []rune{'a'}},
			{Type: tea.KeyRunes, Runes: []rune{'5'}},
		}

		for _, key := range nonGlobalKeys {
			if m.isGlobalKey(key) {
				t.Errorf("キー %v は非グローバルキーとして認識されるべき", key)
			}
		}
	})

	t.Run("アクティブパネルへのメッセージ転送", func(t *testing.T) {
		m := NewDashboardModel()
		m.activePanel = 0 // カウンターパネル

		// カウンターの初期値を確認
		counter := m.panels[0].model.(counterModel)
		initialCount := counter.count

		// ↑キーを送信（カウンター増加）
		msg := tea.KeyMsg{Type: tea.KeyUp}
		newModel, _ := m.Update(msg)
		updatedModel := newModel.(dashboardModel)

		// カウンターが増加していることを確認
		updatedCounter := updatedModel.panels[0].model.(counterModel)
		if updatedCounter.count != initialCount+1 {
			t.Errorf("アクティブパネル（カウンター）でカウントが増加すべき、実際: %d → %d", 
				initialCount, updatedCounter.count)
		}
	})
}

func TestDashboardModelView(t *testing.T) {
	t.Run("View表示確認", func(t *testing.T) {
		m := NewDashboardModel()
		view := m.View()

		if view == "" {
			t.Error("ビューは空であってはいけない")
		}

		// 必要な要素の確認
		requiredElements := []string{
			"ダッシュボード",
			"カウンター",
			"タイマー",
			"TODO",
			"GitHub",
			"Tab",
			"Ctrl+C",
		}

		for _, element := range requiredElements {
			if !strings.Contains(view, element) {
				t.Errorf("ビューに「%s」が含まれているべき", element)
			}
		}
	})

	t.Run("アクティブパネルの表示", func(t *testing.T) {
		m := NewDashboardModel()
		m.activePanel = 1 // タイマーパネル
		m.panels[0].active = false
		m.panels[1].active = true

		view := m.View()

		// アクティブパネルのマーカー確認
		if !strings.Contains(view, "タイマー ★") {
			t.Error("アクティブパネルには★マーカーが表示されるべき")
		}
	})

	t.Run("ヘルプ表示の切り替え", func(t *testing.T) {
		m := NewDashboardModel()

		// ヘルプ表示時
		m.showHelp = true
		view := m.View()
		if !strings.Contains(view, "Tab") {
			t.Error("ヘルプ表示時はキーバインディング情報が含まれるべき")
		}

		// ヘルプ非表示時
		m.showHelp = false
		view = m.View()
		if !strings.Contains(view, "F1でヘルプを表示") {
			t.Error("ヘルプ非表示時はF1キーの案内が表示されるべき")
		}
	})

	t.Run("グローバルヘルプの表示", func(t *testing.T) {
		m := NewDashboardModel()
		m.showHelp = true
		m.globalHelp = true

		view := m.View()

		// 詳細ヘルプの要素確認
		detailedHelpElements := []string{
			"グローバルキー",
			"パネル切替",
			"直接選択",
			"詳細ヘルプ",
		}

		for _, element := range detailedHelpElements {
			if !strings.Contains(view, element) {
				t.Errorf("グローバルヘルプに「%s」が含まれているべき", element)
			}
		}
	})
}

func TestTruncateContent(t *testing.T) {
	t.Run("コンテンツの切り詰め", func(t *testing.T) {
		m := NewDashboardModel()

		// 長いコンテンツ
		longContent := strings.Repeat("あいうえおかきくけこ", 10) + "\n" +
			strings.Repeat("さしすせそたちつてと", 10) + "\n" +
			strings.Repeat("なにぬねのはひふへほ", 10)

		// 幅と高さを制限
		truncated := m.truncateContent(longContent, 20, 2)
		lines := strings.Split(truncated, "\n")

		// 行数の確認
		if len(lines) != 2 {
			t.Errorf("切り詰め後の行数は2であるべき、実際: %d", len(lines))
		}

		// 各行の幅確認
		for i, line := range lines {
			if len(line) > 20 {
				t.Errorf("行%dの幅は20以下であるべき、実際: %d", i, len(line))
			}
		}
	})

	t.Run("短いコンテンツの補完", func(t *testing.T) {
		m := NewDashboardModel()

		// 短いコンテンツ
		shortContent := "短いテキスト"

		// 高さを指定
		truncated := m.truncateContent(shortContent, 50, 5)
		lines := strings.Split(truncated, "\n")

		// 行数の確認（空行で補完される）
		if len(lines) != 5 {
			t.Errorf("補完後の行数は5であるべき、実際: %d", len(lines))
		}

		// 最初の行は元のコンテンツ
		if lines[0] != shortContent {
			t.Errorf("最初の行は元のコンテンツであるべき、実際: %s", lines[0])
		}

		// 残りの行は空行
		for i := 1; i < len(lines); i++ {
			if lines[i] != "" {
				t.Errorf("行%dは空行であるべき、実際: '%s'", i, lines[i])
			}
		}
	})
}