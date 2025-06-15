// Package constants provides common constants used across the Bubble Tea applications.
package constants

import "time"

// Application dimensions
const (
	DefaultWidth        = 50
	DefaultHeight       = 20
	DefaultPanelWidth   = 40
	DefaultPanelHeight  = 15
	MinWidth            = 20
	MinHeight           = 10
)

// Timer constants
const (
	TickInterval = 100 * time.Millisecond
)

// Form constants
const (
	FormFieldMaxLength = 50
	EmailMaxLength     = 100
)

// GitHub API constants
const (
	GitHubAPITimeout = 10 * time.Second
	GitHubAPIBaseURL = "https://api.github.com"
	MaxRetries       = 3
	DemoDelay        = 500 * time.Millisecond // For API simulation
)

// Dashboard constants
const (
	PanelRows    = 2
	PanelColumns = 2
	PanelCount   = PanelRows * PanelColumns
)

// Key binding help texts
const (
	QuitHelp         = "q/Ctrl+C: 終了"
	NavigationHelp   = "↑/↓/j/k: 移動"
	SelectHelp       = "Enter/Space: 選択"
	TabHelp          = "Tab/Shift+Tab: 切り替え"
	EscHelp          = "Esc: 戻る"
	F1Help           = "F1: ヘルプ切り替え"
)

// Application titles
const (
	CounterTitle   = "🔢 カウンターアプリ"
	TimerTitle     = "⏱️  タイマーアプリ"
	TodoTitle      = "📝 TODOリスト"
	FormTitle      = "📋 フォーム入力"
	GitHubTitle    = "🐙 GitHub ユーザー検索"
	DashboardTitle = "🎛️  Bubble Tea ダッシュボード - 統合アプリケーション"
)