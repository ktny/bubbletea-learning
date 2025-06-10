package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TODOアイテムの構造体
type todoItem struct {
	title     string
	completed bool
}

// TODOリストモデル
type todoModel struct {
	items    []todoItem // TODOアイテムのリスト
	cursor   int        // 現在選択している項目のインデックス
	viewport int        // ビューポートの開始位置
	height   int        // 表示可能な行数
}

// コンストラクタ
func NewTodoModel() todoModel {
	// 初期データ
	items := []todoItem{
		{title: "Bubble Teaの基本を学ぶ", completed: true},
		{title: "キーボードイベントを理解する", completed: true},
		{title: "スタイリングを適用する", completed: true},
		{title: "非同期処理を実装する", completed: true},
		{title: "リスト表示を作成する", completed: false},
		{title: "テキスト入力を実装する", completed: false},
		{title: "API連携を追加する", completed: false},
		{title: "複雑なレイアウトを構築する", completed: false},
	}

	return todoModel{
		items:    items,
		cursor:   0,
		viewport: 0,
		height:   10, // デフォルトの表示行数
	}
}

// Init - 初期化
func (m todoModel) Init() tea.Cmd {
	return nil
}

// Update - メッセージ処理
func (m todoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			m = m.moveCursorUp()
		case tea.KeyDown:
			m = m.moveCursorDown()
		case tea.KeyEnter, tea.KeySpace:
			m = m.toggleItem()
		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "q":
				return m, tea.Quit
			case "j":
				m = m.moveCursorDown()
			case "k":
				m = m.moveCursorUp()
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		// ウィンドウサイズの変更に対応
		m.height = msg.Height - 10 // ヘッダーやフッター分を引く
		if m.height < 1 {
			m.height = 1
		}
	}

	return m, nil
}

// カーソルを上に移動
func (m todoModel) moveCursorUp() todoModel {
	if m.cursor > 0 {
		m.cursor--
		// ビューポートの調整
		if m.cursor < m.viewport {
			m.viewport = m.cursor
		}
	}
	return m
}

// カーソルを下に移動
func (m todoModel) moveCursorDown() todoModel {
	if m.cursor < len(m.items)-1 {
		m.cursor++
		// ビューポートの調整
		if m.cursor >= m.viewport+m.height {
			m.viewport = m.cursor - m.height + 1
		}
	}
	return m
}

// 項目の完了状態を切り替え
func (m todoModel) toggleItem() todoModel {
	if m.cursor < len(m.items) {
		m.items[m.cursor].completed = !m.items[m.cursor].completed
	}
	return m
}

// View - UIの描画
func (m todoModel) View() string {
	// スタイル定義
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		MarginBottom(1)

	cursorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("14")).
		Bold(true)

	completedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Strikethrough(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7"))

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Italic(true).
		MarginTop(1)

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("12")).
		Padding(1, 2)

	// コンテンツの構築
	var content strings.Builder
	content.WriteString(titleStyle.Render("📝 TODOリスト"))
	content.WriteString("\n\n")

	// アイテムの表示（ビューポート内のみ）
	end := m.viewport + m.height
	if end > len(m.items) {
		end = len(m.items)
	}

	for i := m.viewport; i < end; i++ {
		item := m.items[i]
		checkbox := "[ ]"
		if item.completed {
			checkbox = "[✓]"
		}

		line := fmt.Sprintf("%s %s", checkbox, item.title)

		// スタイルの適用
		if i == m.cursor {
			content.WriteString(cursorStyle.Render("> " + line))
		} else if item.completed {
			content.WriteString(completedStyle.Render("  " + line))
		} else {
			content.WriteString(normalStyle.Render("  " + line))
		}

		if i < end-1 {
			content.WriteString("\n")
		}
	}

	// スクロールインジケーター
	if len(m.items) > m.height {
		scrollInfo := fmt.Sprintf("\n\n[%d/%d]", m.cursor+1, len(m.items))
		content.WriteString(normalStyle.Render(scrollInfo))
	}

	// ヘルプテキスト
	help := helpStyle.Render("\n↑/k: 上へ  ↓/j: 下へ  Enter/Space: 選択  q: 終了")
	content.WriteString(help)

	return borderStyle.Render(content.String())
}
