package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// カスタムメッセージ型
type tickMsg time.Time

// タイマーの状態
type timerState int

const (
	stopped timerState = iota
	running
	paused
)

// タイマーモデル
type timerModel struct {
	duration   time.Duration // 経過時間
	state      timerState    // タイマーの状態
	startTime  time.Time     // 開始時刻
	pausedTime time.Duration // 一時停止時の累積時間
}

// タイマーモデルのコンストラクタ
func NewTimerModel() timerModel {
	return timerModel{}.handleReset()
}

// Init - 初期化時のコマンド
func (m timerModel) Init() tea.Cmd {
	return nil
}

// tick コマンドを生成する関数
func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// handleStartStop - Start/Stop トグル処理
func (m timerModel) handleStartStop() (timerModel, tea.Cmd) {
	switch m.state {
	case stopped:
		m.state = running
		m.startTime = time.Now()
		m.pausedTime = 0
		return m, tickCmd()
	case running:
		m.state = paused
		m.pausedTime = m.duration
		return m, nil
	case paused:
		m.state = running
		m.startTime = time.Now()
		return m, tickCmd()
	}
	return m, nil
}

// handleReset - リセット処理
func (m timerModel) handleReset() timerModel {
	m.state = stopped
	m.duration = 0
	m.pausedTime = 0
	m.startTime = time.Time{}
	return m
}

// handleTick - tick メッセージ処理
func (m timerModel) handleTick() (tea.Model, tea.Cmd) {
	if m.state == running {
		// 現在時刻から開始時刻を引いて、一時停止時の累積時間を加算
		m.duration = time.Since(m.startTime) + m.pausedTime
		return m, tickCmd() // 次のtickを予約
	}
	return m, nil
}

// Update - メッセージ処理
func (m timerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Start/Stop処理の判定
		isStartStop := msg.Type == tea.KeySpace ||
			(msg.Type == tea.KeyRunes && string(msg.Runes) == "s")

		if isStartStop {
			return m.handleStartStop()
		}

		// その他のキー処理
		switch msg.Type {
		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "r":
				return m.handleReset(), nil
			case "q":
				return m, tea.Quit
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		}

	case tickMsg:
		return m.handleTick()
	}

	return m, nil
}

// formatDuration - MM:SS.ms形式で時間をフォーマット
func formatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	milliseconds := int(d.Milliseconds()) % 1000 / 100 // 1桁のミリ秒
	return fmt.Sprintf("%02d:%02d.%d", minutes, seconds, milliseconds)
}

// View - UIの描画
func (m timerModel) View() string {
	// スタイル定義
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")). // 青
		Align(lipgloss.Center)

	timeStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("10")). // 緑
		Align(lipgloss.Center).
		Width(12)

	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")). // グレー
		Italic(true)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")). // グレー
		Italic(true)

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("12")). // 青
		Padding(1, 2)

	// 状態テキスト
	var stateText string
	switch m.state {
	case stopped:
		stateText = "停止中"
	case running:
		stateText = "実行中"
	case paused:
		stateText = "一時停止中"
	}

	content := fmt.Sprintf(
		"%s\n\n%s\n\n%s\n\n%s",
		titleStyle.Render("⏱️  タイマー"),
		timeStyle.Render(formatDuration(m.duration)),
		statusStyle.Render(fmt.Sprintf("状態: %s", stateText)),
		helpStyle.Render(
			"s/スペース: スタート/ストップ\n"+
				"r: リセット\n"+
				"q: 終了",
		),
	)

	return borderStyle.Render(content)
}
