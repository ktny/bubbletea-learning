package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// パネルの種類
type panelType int

const (
	panelCounter panelType = iota
	panelTimer
	panelTodo
	panelGithub
)

// パネル情報
type panel struct {
	title    string
	model    tea.Model
	panelType panelType
	active   bool
}

// ダッシュボードの状態
type dashboardModel struct {
	panels       []panel
	activePanel  int
	width        int
	height       int
	showHelp     bool
	globalHelp   bool
}

// コンストラクタ
func NewDashboardModel() dashboardModel {
	panels := []panel{
		{
			title:     "カウンター",
			model:     NewCounterModel(),
			panelType: panelCounter,
			active:    true,
		},
		{
			title:     "タイマー",
			model:     NewTimerModel(),
			panelType: panelTimer,
			active:    false,
		},
		{
			title:     "TODO",
			model:     NewTodoModel(),
			panelType: panelTodo,
			active:    false,
		},
		{
			title:     "GitHub",
			model:     NewGitHubModel(),
			panelType: panelGithub,
			active:    false,
		},
	}

	return dashboardModel{
		panels:      panels,
		activePanel: 0,
		width:       80,
		height:      24,
		showHelp:    true,
		globalHelp:  false,
	}
}

// Init - 初期化
func (m dashboardModel) Init() tea.Cmd {
	// 全パネルの初期化コマンドを収集
	var cmds []tea.Cmd
	for _, panel := range m.panels {
		if cmd := panel.model.Init(); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return tea.Batch(cmds...)
}

// Update - メッセージ処理
func (m dashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// グローバルキーバインディング
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyTab:
			// 次のパネルに切り替え
			m.panels[m.activePanel].active = false
			m.activePanel = (m.activePanel + 1) % len(m.panels)
			m.panels[m.activePanel].active = true
			return m, nil

		case tea.KeyShiftTab:
			// 前のパネルに切り替え
			m.panels[m.activePanel].active = false
			m.activePanel = (m.activePanel - 1 + len(m.panels)) % len(m.panels)
			m.panels[m.activePanel].active = true
			return m, nil

		case tea.KeyF1:
			// ヘルプの表示切り替え
			m.showHelp = !m.showHelp
			return m, nil

		case tea.KeyF2:
			// グローバルヘルプの表示切り替え
			m.globalHelp = !m.globalHelp
			return m, nil

		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "1":
				// カウンターパネルに直接切り替え
				m.panels[m.activePanel].active = false
				m.activePanel = 0
				m.panels[m.activePanel].active = true
				return m, nil
			case "2":
				// タイマーパネルに直接切り替え
				m.panels[m.activePanel].active = false
				m.activePanel = 1
				m.panels[m.activePanel].active = true
				return m, nil
			case "3":
				// TODOパネルに直接切り替え
				m.panels[m.activePanel].active = false
				m.activePanel = 2
				m.panels[m.activePanel].active = true
				return m, nil
			case "4":
				// GitHubパネルに直接切り替え
				m.panels[m.activePanel].active = false
				m.activePanel = 3
				m.panels[m.activePanel].active = true
				return m, nil
			}
		}

		// アクティブパネルにメッセージを転送（グローバルキー以外）
		if !m.isGlobalKey(msg) {
			activeModel := m.panels[m.activePanel].model
			newModel, cmd := activeModel.Update(msg)
			m.panels[m.activePanel].model = newModel
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}

	case tea.WindowSizeMsg:
		// ウィンドウサイズ変更
		m.width = msg.Width
		m.height = msg.Height

		// 各パネルにサイズ変更を通知
		for i, panel := range m.panels {
			// パネルサイズを計算（2x2レイアウト）
			panelWidth := msg.Width / 2
			panelHeight := (msg.Height - 4) / 2 // ヘッダー・フッター分を除外

			// ウィンドウサイズメッセージを作成
			panelSizeMsg := tea.WindowSizeMsg{
				Width:  panelWidth - 2, // ボーダー分を除外
				Height: panelHeight - 2,
			}

			newModel, cmd := panel.model.Update(panelSizeMsg)
			m.panels[i].model = newModel
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}

	default:
		// その他のメッセージは全パネルに配信
		for i, panel := range m.panels {
			newModel, cmd := panel.model.Update(msg)
			m.panels[i].model = newModel
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

// グローバルキーかどうかを判定
func (m dashboardModel) isGlobalKey(msg tea.KeyMsg) bool {
	switch msg.Type {
	case tea.KeyTab, tea.KeyShiftTab, tea.KeyF1, tea.KeyF2, tea.KeyCtrlC:
		return true
	case tea.KeyRunes:
		switch string(msg.Runes) {
		case "1", "2", "3", "4":
			return true
		}
	}
	return false
}

// View - UIの描画
func (m dashboardModel) View() string {
	// スタイル定義
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		Background(lipgloss.Color("235")).
		Padding(0, 1).
		Width(m.width)

	activePanelStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("12")).
		Padding(1)

	inactivePanelStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).
		Padding(1)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Background(lipgloss.Color("235")).
		Padding(0, 1).
		Width(m.width)

	// タイトルバー
	title := titleStyle.Render("🎛️  Bubble Tea ダッシュボード - 統合アプリケーション")

	// パネルサイズ計算
	panelWidth := (m.width - 4) / 2 // 2列レイアウト
	panelHeight := (m.height - 6) / 2 // タイトル・ヘルプ分を除外

	// 各パネルの描画
	var panels [4]string
	for i, panel := range m.panels {
		// パネルのスタイルを選択
		var style lipgloss.Style
		if i == m.activePanel {
			style = activePanelStyle
		} else {
			style = inactivePanelStyle
		}

		// パネルタイトル
		panelTitle := fmt.Sprintf("[%d] %s", i+1, panel.title)
		if i == m.activePanel {
			panelTitle += " ★"
		}

		// パネル内容を取得
		content := panel.model.View()

		// サイズに合わせて調整
		content = m.truncateContent(content, panelWidth-4, panelHeight-4)

		// パネル全体
		panelContent := panelTitle + "\n" + strings.Repeat("─", panelWidth-4) + "\n" + content
		panels[i] = style.Width(panelWidth).Height(panelHeight).Render(panelContent)
	}

	// 2x2レイアウトで配置
	topRow := lipgloss.JoinHorizontal(lipgloss.Top, panels[0], panels[1])
	bottomRow := lipgloss.JoinHorizontal(lipgloss.Top, panels[2], panels[3])
	mainContent := lipgloss.JoinVertical(lipgloss.Left, topRow, bottomRow)

	// ヘルプテキスト
	var helpText string
	if m.showHelp {
		if m.globalHelp {
			helpText = "グローバルキー: Tab/Shift+Tab:パネル切替 | 1-4:直接選択 | F1:ヘルプ | F2:詳細ヘルプ | Ctrl+C:終了"
		} else {
			helpText = "Tab:次 | Shift+Tab:前 | 1-4:選択 | F1:ヘルプ切替 | F2:詳細 | Ctrl+C:終了"
		}
	} else {
		helpText = "F1でヘルプを表示"
	}

	help := helpStyle.Render(helpText)

	// 最終的なレイアウト
	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		mainContent,
		help,
	)
}

// コンテンツを指定されたサイズに切り詰める
func (m dashboardModel) truncateContent(content string, width, height int) string {
	lines := strings.Split(content, "\n")

	// 高さの調整
	if len(lines) > height {
		lines = lines[:height]
	}

	// 幅の調整
	for i, line := range lines {
		if len(line) > width {
			lines[i] = line[:width-3] + "..."
		}
	}

	// 不足行を空行で埋める
	for len(lines) < height {
		lines = append(lines, "")
	}

	return strings.Join(lines, "\n")
}