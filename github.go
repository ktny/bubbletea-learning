package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// GitHubアプリの状態
type githubState int

const (
	stateInput githubState = iota
	stateLoading
	stateSuccess
	stateError
)

// GitHubユーザー情報の構造体
type githubUser struct {
	Login       string    `json:"login"`
	Name        string    `json:"name"`
	Company     string    `json:"company"`
	Blog        string    `json:"blog"`
	Location    string    `json:"location"`
	Email       string    `json:"email"`
	Bio         string    `json:"bio"`
	PublicRepos int       `json:"public_repos"`
	PublicGists int       `json:"public_gists"`
	Followers   int       `json:"followers"`
	Following   int       `json:"following"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	AvatarURL   string    `json:"avatar_url"`
	HTMLURL     string    `json:"html_url"`
}

// API応答メッセージ
type apiResponse struct {
	user *githubUser
	err  error
}

// GitHubモデル
type githubModel struct {
	input       textinput.Model
	spinner     spinner.Model
	state       githubState
	user        *githubUser
	errorMsg    string
	retryCount  int
	maxRetries  int
	lastRequest string
}

// コンストラクタ
func NewGitHubModel() githubModel {
	// テキスト入力の設定
	ti := textinput.New()
	ti.Placeholder = "例: octocat"
	ti.Focus()
	ti.CharLimit = 39 // GitHubのユーザー名制限
	ti.Width = 30

	// スピナーの設定
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))

	return githubModel{
		input:      ti,
		spinner:    sp,
		state:      stateInput,
		maxRetries: 3,
	}
}

// Init - 初期化
func (m githubModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update - メッセージ処理
func (m githubModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case stateInput:
			switch msg.Type {
			case tea.KeyEnter:
				username := strings.TrimSpace(m.input.Value())
				if username != "" {
					m.state = stateLoading
					m.lastRequest = username
					m.retryCount = 0
					return m, tea.Batch(
						m.spinner.Tick,
						fetchGitHubUser(username),
					)
				}
			case tea.KeyEsc, tea.KeyCtrlC:
				return m, tea.Quit
			}

		case stateError:
			switch msg.Type {
			case tea.KeyEnter:
				// リトライ
				if m.retryCount < m.maxRetries {
					m.state = stateLoading
					m.retryCount++
					return m, tea.Batch(
						m.spinner.Tick,
						fetchGitHubUser(m.lastRequest),
					)
				}
			case tea.KeyEsc:
				// 入力画面に戻る
				m.state = stateInput
				m.errorMsg = ""
				m.input.SetValue("")
				m.input.Focus()
				return m, textinput.Blink
			case tea.KeyCtrlC:
				return m, tea.Quit
			}

		case stateSuccess:
			switch msg.Type {
			case tea.KeyEnter:
				// 新しい検索
				m.state = stateInput
				m.user = nil
				m.input.SetValue("")
				m.input.Focus()
				return m, textinput.Blink
			case tea.KeyEsc, tea.KeyCtrlC:
				return m, tea.Quit
			}
		}

	case apiResponse:
		if msg.err != nil {
			m.state = stateError
			m.errorMsg = msg.err.Error()
			return m, nil
		}
		m.state = stateSuccess
		m.user = msg.user
		return m, nil

	case spinner.TickMsg:
		if m.state == stateLoading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	// テキスト入力の更新
	if m.state == stateInput {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}

	return m, nil
}

// GitHubユーザー情報を取得
func fetchGitHubUser(username string) tea.Cmd {
	return func() tea.Msg {
		// API呼び出し
		url := fmt.Sprintf("https://api.github.com/users/%s", username)
		
		client := &http.Client{
			Timeout: 10 * time.Second,
		}
		
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return apiResponse{err: err}
		}
		
		// User-Agentヘッダーを設定（GitHub API要件）
		req.Header.Set("User-Agent", "bubbletea-learning")
		
		resp, err := client.Do(req)
		if err != nil {
			return apiResponse{err: fmt.Errorf("ネットワークエラー: %w", err)}
		}
		defer resp.Body.Close()

		// ステータスコードのチェック
		if resp.StatusCode == 404 {
			return apiResponse{err: fmt.Errorf("ユーザー '%s' が見つかりません", username)}
		}
		if resp.StatusCode != 200 {
			return apiResponse{err: fmt.Errorf("APIエラー: ステータスコード %d", resp.StatusCode)}
		}

		// JSONパース
		var user githubUser
		if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
			return apiResponse{err: fmt.Errorf("JSONパースエラー: %w", err)}
		}

		// 少し遅延を入れて読み込み画面を見えやすくする（デモ用）
		time.Sleep(500 * time.Millisecond)

		return apiResponse{user: &user}
	}
}

// View - UIの描画
func (m githubModel) View() string {
	// スタイル定義
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		MarginBottom(1)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7")).
		Width(12)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("14"))

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")).
		Bold(true)

	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("10"))

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Italic(true).
		MarginTop(1)

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("12")).
		Padding(1, 2).
		Width(50)

	var content string

	switch m.state {
	case stateInput:
		content = titleStyle.Render("🐙 GitHub ユーザー検索") + "\n\n"
		content += "ユーザー名を入力してください:\n"
		content += m.input.View() + "\n\n"
		content += helpStyle.Render("Enter: 検索  Esc: 終了")

	case stateLoading:
		content = titleStyle.Render("🐙 GitHub ユーザー検索") + "\n\n"
		content += m.spinner.View() + " " + 
			fmt.Sprintf("'%s' を検索中...", m.lastRequest) + "\n\n"
		if m.retryCount > 0 {
			content += fmt.Sprintf("リトライ %d/%d\n", m.retryCount, m.maxRetries)
		}

	case stateError:
		content = titleStyle.Render("🐙 GitHub ユーザー検索") + "\n\n"
		content += errorStyle.Render("❌ エラーが発生しました") + "\n\n"
		content += m.errorMsg + "\n\n"
		if m.retryCount < m.maxRetries {
			content += helpStyle.Render("Enter: リトライ  Esc: 戻る  Ctrl+C: 終了")
		} else {
			content += errorStyle.Render(fmt.Sprintf("リトライ回数が上限（%d回）に達しました", m.maxRetries)) + "\n"
			content += helpStyle.Render("Esc: 戻る  Ctrl+C: 終了")
		}

	case stateSuccess:
		if m.user != nil {
			content = titleStyle.Render("🐙 GitHub ユーザー情報") + "\n\n"
			content += successStyle.Render("✅ ユーザーが見つかりました！") + "\n\n"
			
			// ユーザー情報の表示
			content += labelStyle.Render("ユーザー名:") + " " + valueStyle.Render(m.user.Login) + "\n"
			
			if m.user.Name != "" {
				content += labelStyle.Render("名前:") + " " + valueStyle.Render(m.user.Name) + "\n"
			}
			
			if m.user.Company != "" {
				content += labelStyle.Render("会社:") + " " + valueStyle.Render(m.user.Company) + "\n"
			}
			
			if m.user.Location != "" {
				content += labelStyle.Render("場所:") + " " + valueStyle.Render(m.user.Location) + "\n"
			}
			
			if m.user.Bio != "" {
				content += labelStyle.Render("自己紹介:") + " " + valueStyle.Render(m.user.Bio) + "\n"
			}
			
			content += labelStyle.Render("公開リポ数:") + " " + valueStyle.Render(fmt.Sprintf("%d", m.user.PublicRepos)) + "\n"
			content += labelStyle.Render("フォロワー:") + " " + valueStyle.Render(fmt.Sprintf("%d", m.user.Followers)) + "\n"
			content += labelStyle.Render("フォロー中:") + " " + valueStyle.Render(fmt.Sprintf("%d", m.user.Following)) + "\n"
			content += labelStyle.Render("登録日:") + " " + valueStyle.Render(m.user.CreatedAt.Format("2006年1月2日")) + "\n"
			content += labelStyle.Render("URL:") + " " + valueStyle.Render(m.user.HTMLURL) + "\n\n"
			
			content += helpStyle.Render("Enter: 新しい検索  Esc: 終了")
		}
	}

	return borderStyle.Render(content)
}