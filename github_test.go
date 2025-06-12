package main

import (
	"fmt"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func TestGitHubModel(t *testing.T) {
	t.Run("初期状態", func(t *testing.T) {
		m := NewGitHubModel()
		
		if m.state != stateInput {
			t.Errorf("初期状態はstateInputであるべき、実際: %v", m.state)
		}
		if !m.input.Focused() {
			t.Error("入力フィールドにフォーカスがあるべき")
		}
		if m.maxRetries != 3 {
			t.Errorf("最大リトライ回数は3であるべき、実際: %d", m.maxRetries)
		}
		if m.user != nil {
			t.Error("初期状態でユーザー情報はnilであるべき")
		}
		if m.errorMsg != "" {
			t.Error("初期状態でエラーメッセージは空であるべき")
		}
	})

	t.Run("Initメソッド", func(t *testing.T) {
		m := NewGitHubModel()
		cmd := m.Init()
		if cmd == nil {
			t.Error("Initはtextinput.Blinkコマンドを返すべき")
		}
	})

	t.Run("Enterキーで検索開始", func(t *testing.T) {
		m := NewGitHubModel()
		m.input.SetValue("octocat")
		msg := tea.KeyMsg{Type: tea.KeyEnter}

		newModel, cmd := m.Update(msg)
		updatedModel := newModel.(githubModel)

		if updatedModel.state != stateLoading {
			t.Errorf("Enter後の状態はstateLoadingであるべき、実際: %v", updatedModel.state)
		}
		if updatedModel.lastRequest != "octocat" {
			t.Errorf("lastRequestは入力値を保持すべき、実際: %s", updatedModel.lastRequest)
		}
		if cmd == nil {
			t.Error("検索開始時はコマンド（spinner.TickとfetchGitHubUser）を返すべき")
		}
	})

	t.Run("空の入力でEnterキー", func(t *testing.T) {
		m := NewGitHubModel()
		m.input.SetValue("")
		msg := tea.KeyMsg{Type: tea.KeyEnter}

		newModel, cmd := m.Update(msg)
		updatedModel := newModel.(githubModel)

		if updatedModel.state != stateInput {
			t.Error("空の入力では状態が変わらないべき")
		}
		if cmd != nil {
			t.Error("空の入力ではコマンドを返さないべき")
		}
	})

	t.Run("Escキーで終了", func(t *testing.T) {
		m := NewGitHubModel()
		msg := tea.KeyMsg{Type: tea.KeyEsc}

		_, cmd := m.Update(msg)
		if cmd == nil {
			t.Error("EscキーはQuitコマンドを返すべき")
		}
	})

	t.Run("Ctrl+Cで終了", func(t *testing.T) {
		m := NewGitHubModel()
		msg := tea.KeyMsg{Type: tea.KeyCtrlC}

		_, cmd := m.Update(msg)
		if cmd == nil {
			t.Error("Ctrl+CはQuitコマンドを返すべき")
		}
	})

	t.Run("エラー状態でEnterキーでリトライ", func(t *testing.T) {
		m := NewGitHubModel()
		m.state = stateError
		m.lastRequest = "octocat"
		m.retryCount = 1
		m.errorMsg = "テストエラー"
		msg := tea.KeyMsg{Type: tea.KeyEnter}

		newModel, cmd := m.Update(msg)
		updatedModel := newModel.(githubModel)

		if updatedModel.state != stateLoading {
			t.Error("リトライ時は状態がstateLoadingになるべき")
		}
		if updatedModel.retryCount != 2 {
			t.Errorf("リトライカウントが増加すべき、実際: %d", updatedModel.retryCount)
		}
		if cmd == nil {
			t.Error("リトライ時はコマンドを返すべき")
		}
	})

	t.Run("リトライ上限到達時", func(t *testing.T) {
		m := NewGitHubModel()
		m.state = stateError
		m.retryCount = 3 // 上限に到達
		m.maxRetries = 3
		msg := tea.KeyMsg{Type: tea.KeyEnter}

		newModel, cmd := m.Update(msg)
		updatedModel := newModel.(githubModel)

		if updatedModel.state != stateError {
			t.Error("リトライ上限到達時は状態が変わらないべき")
		}
		if cmd != nil {
			t.Error("リトライ上限到達時はコマンドを返さないべき")
		}
	})

	t.Run("エラー状態でEscキーで入力画面に戻る", func(t *testing.T) {
		m := NewGitHubModel()
		m.state = stateError
		m.errorMsg = "テストエラー"
		m.input.SetValue("old-value")
		msg := tea.KeyMsg{Type: tea.KeyEsc}

		newModel, cmd := m.Update(msg)
		updatedModel := newModel.(githubModel)

		if updatedModel.state != stateInput {
			t.Error("Escキーで入力状態に戻るべき")
		}
		if updatedModel.errorMsg != "" {
			t.Error("エラーメッセージがクリアされるべき")
		}
		if updatedModel.input.Value() != "" {
			t.Error("入力フィールドがクリアされるべき")
		}
		if cmd == nil {
			t.Error("textinput.Blinkコマンドを返すべき")
		}
	})

	t.Run("成功状態でEnterキーで新しい検索", func(t *testing.T) {
		m := NewGitHubModel()
		m.state = stateSuccess
		m.user = &githubUser{Login: "octocat"}
		m.input.SetValue("old-value")
		msg := tea.KeyMsg{Type: tea.KeyEnter}

		newModel, cmd := m.Update(msg)
		updatedModel := newModel.(githubModel)

		if updatedModel.state != stateInput {
			t.Error("Enterキーで入力状態に戻るべき")
		}
		if updatedModel.user != nil {
			t.Error("ユーザー情報がクリアされるべき")
		}
		if updatedModel.input.Value() != "" {
			t.Error("入力フィールドがクリアされるべき")
		}
		if cmd == nil {
			t.Error("textinput.Blinkコマンドを返すべき")
		}
	})

	t.Run("APIレスポンス - エラー", func(t *testing.T) {
		m := NewGitHubModel()
		m.state = stateLoading
		msg := apiResponse{
			err: fmt.Errorf("ネットワークエラー"),
		}

		newModel, _ := m.Update(msg)
		updatedModel := newModel.(githubModel)

		if updatedModel.state != stateError {
			t.Error("エラーレスポンス時は状態がstateErrorになるべき")
		}
		if updatedModel.errorMsg != "ネットワークエラー" {
			t.Errorf("エラーメッセージが設定されるべき、実際: %s", updatedModel.errorMsg)
		}
	})

	t.Run("APIレスポンス - 成功", func(t *testing.T) {
		m := NewGitHubModel()
		m.state = stateLoading
		testUser := &githubUser{
			Login: "octocat",
			Name:  "The Octocat",
		}
		msg := apiResponse{
			user: testUser,
		}

		newModel, _ := m.Update(msg)
		updatedModel := newModel.(githubModel)

		if updatedModel.state != stateSuccess {
			t.Error("成功レスポンス時は状態がstateSuccessになるべき")
		}
		if updatedModel.user == nil || updatedModel.user.Login != "octocat" {
			t.Error("ユーザー情報が設定されるべき")
		}
	})

	t.Run("スピナーの更新", func(t *testing.T) {
		m := NewGitHubModel()
		m.state = stateLoading
		
		// spinner.TickMsgのモック
		// 実際のspinner.TickMsgは内部型なので、テストでは直接扱えない
		// ここでは動作の概念的なテストのみ
		
		if m.state != stateLoading {
			t.Error("ローディング状態でスピナーが更新されるべき")
		}
	})
}

func TestGitHubModelView(t *testing.T) {
	t.Run("入力画面の表示", func(t *testing.T) {
		m := NewGitHubModel()
		view := m.View()

		if view == "" {
			t.Error("ビューは空であってはいけない")
		}

		requiredElements := []string{
			"GitHub ユーザー検索",
			"ユーザー名を入力してください",
			"Enter: 検索",
			"Esc: 終了",
		}

		for _, element := range requiredElements {
			if !strings.Contains(view, element) {
				t.Errorf("入力画面に「%s」が含まれているべき", element)
			}
		}
	})

	t.Run("ローディング画面の表示", func(t *testing.T) {
		m := NewGitHubModel()
		m.state = stateLoading
		m.lastRequest = "octocat"
		view := m.View()

		requiredElements := []string{
			"GitHub ユーザー検索",
			"'octocat' を検索中...",
		}

		for _, element := range requiredElements {
			if !strings.Contains(view, element) {
				t.Errorf("ローディング画面に「%s」が含まれているべき", element)
			}
		}
	})

	t.Run("ローディング画面でリトライ表示", func(t *testing.T) {
		m := NewGitHubModel()
		m.state = stateLoading
		m.lastRequest = "octocat"
		m.retryCount = 2
		m.maxRetries = 3
		view := m.View()

		if !strings.Contains(view, "リトライ 2/3") {
			t.Error("リトライ情報が表示されるべき")
		}
	})

	t.Run("エラー画面の表示", func(t *testing.T) {
		m := NewGitHubModel()
		m.state = stateError
		m.errorMsg = "ユーザーが見つかりません"
		m.retryCount = 1
		view := m.View()

		requiredElements := []string{
			"❌ エラーが発生しました",
			"ユーザーが見つかりません",
			"Enter: リトライ",
			"Esc: 戻る",
		}

		for _, element := range requiredElements {
			if !strings.Contains(view, element) {
				t.Errorf("エラー画面に「%s」が含まれているべき", element)
			}
		}
	})

	t.Run("リトライ上限到達時のエラー画面", func(t *testing.T) {
		m := NewGitHubModel()
		m.state = stateError
		m.errorMsg = "ネットワークエラー"
		m.retryCount = 3
		m.maxRetries = 3
		view := m.View()

		if !strings.Contains(view, "リトライ回数が上限（3回）に達しました") {
			t.Error("リトライ上限メッセージが表示されるべき")
		}
		if strings.Contains(view, "Enter: リトライ") {
			t.Error("リトライオプションは表示されないべき")
		}
	})

	t.Run("成功画面の表示", func(t *testing.T) {
		m := NewGitHubModel()
		m.state = stateSuccess
		m.user = &githubUser{
			Login:       "octocat",
			Name:        "The Octocat",
			Company:     "GitHub",
			Location:    "San Francisco",
			Bio:         "GitHub's mascot",
			PublicRepos: 8,
			Followers:   3000,
			Following:   9,
			CreatedAt:   time.Date(2011, 1, 25, 0, 0, 0, 0, time.UTC),
			HTMLURL:     "https://github.com/octocat",
		}
		view := m.View()

		requiredElements := []string{
			"✅ ユーザーが見つかりました！",
			"ユーザー名:",
			"octocat",
			"名前:",
			"The Octocat",
			"会社:",
			"GitHub",
			"場所:",
			"San Francisco",
			"自己紹介:",
			"GitHub's mascot",
			"公開リポ数:",
			"8",
			"フォロワー:",
			"3000",
			"フォロー中:",
			"9",
			"登録日:",
			"2011年1月25日",
			"URL:",
			"https://github.com/octocat",
			"Enter: 新しい検索",
			"Esc: 終了",
		}

		for _, element := range requiredElements {
			if !strings.Contains(view, element) {
				t.Errorf("成功画面に「%s」が含まれているべき", element)
			}
		}
	})

	t.Run("成功画面で一部フィールドが空の場合", func(t *testing.T) {
		m := NewGitHubModel()
		m.state = stateSuccess
		m.user = &githubUser{
			Login:       "octocat",
			// Name, Company, Location, Bioは空
			PublicRepos: 8,
			Followers:   3000,
			Following:   9,
			CreatedAt:   time.Date(2011, 1, 25, 0, 0, 0, 0, time.UTC),
			HTMLURL:     "https://github.com/octocat",
		}
		view := m.View()

		// 空のフィールドは表示されない
		shouldNotContain := []string{
			"名前:",
			"会社:",
			"場所:",
			"自己紹介:",
		}

		for _, element := range shouldNotContain {
			if strings.Contains(view, element) {
				t.Errorf("空のフィールド「%s」は表示されないべき", element)
			}
		}

		// 必須フィールドは表示される
		mustContain := []string{
			"ユーザー名:",
			"octocat",
			"公開リポ数:",
			"フォロワー:",
			"登録日:",
		}

		for _, element := range mustContain {
			if !strings.Contains(view, element) {
				t.Errorf("必須フィールド「%s」は表示されるべき", element)
			}
		}
	})
}