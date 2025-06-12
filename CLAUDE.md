# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

これはBubble TeaフレームワークとGo言語を段階的に学習するためのプロジェクトです。基礎から上級まで9つのissueを通じて実践的に学習します。

## 開発コマンド

```bash
# アプリケーションの実行
go run . <app_name>

# 利用可能なアプリ
go run . counter  # カウンターアプリ
go run . timer    # タイマーアプリ  
go run . todo     # TODOリストアプリ
go run . form     # フォームアプリ
go run . github   # GitHub APIアプリ

# 依存関係のダウンロード
go mod download

# テストの実行
go test -v ./...

# Go標準のフォーマット
go fmt ./...

# ビルド
go build -o bubbletea-learning
```

## アーキテクチャ

現在は基本的なBubble TeaのModel-View-Update (MVU)パターンで実装されています：

- **Model**: アプリケーションの状態を保持する構造体
- **Init()**: 初期化時のコマンドを返す（現在はnil）
- **Update()**: メッセージ（キーボード入力など）を処理し、新しいモデルとコマンドを返す
- **View()**: モデルに基づいてUIを文字列として描画

## 学習進捗と次のステップ

### 完了済み
- **Issue #2**: Hello Worldアプリケーション（main.go）
  - Bubble TeaのMVUアーキテクチャ理解
  - Model, View, Update の基本実装
  
- **Issue #3**: ユーザー入力 - キーボードイベントの処理（counter.go, counter_test.go）
  - tea.KeyMsgでのキーボードイベント処理
  - TDD（テスト駆動開発）の実践
  - MVUパターンでの状態管理
  
- **Issue #4**: スタイリング - Lip Glossでターミナルを彩る
  - Lip Glossによる条件付きスタイリング
  - 枠線、色分け、レイアウト制御
  - ターミナルUI設計原則の理解

- **Issue #5**: 非同期処理 - タイマーアプリケーション（timer.go, timer_test.go）
  - tea.Tickによる定期実行
  - 時間計測とミリ秒精度表示
  - 開始/停止/リセット機能の実装

- **Issue #6**: リスト管理 - TODOアプリケーション（todo.go, todo_test.go）
  - ビューポートによるスクロール制御
  - カーソル移動と項目選択
  - 大量データの効率的な表示

- **Issue #7**: コンポーネント活用 - フォームアプリケーション（form.go, form_test.go）
  - Bubblesライブラリのtextinput使用
  - フォーカス管理とフィールド間移動
  - バリデーションとエラーハンドリング

- **Issue #8**: HTTP API統合 - GitHub検索アプリケーション（github.go, github_test.go）
  - Bubblesライブラリのspinner使用
  - 非同期HTTP通信とAPI統合
  - エラーハンドリングとリトライ機能
  - JSON パースと構造化データ表示

### 次のステップ
Issue #9の上級編に進む準備完了

## 開発方針

1. **TDD（テスト駆動開発）を原則とする**
   - 新機能実装前にテストを作成
   - テストが失敗することを確認してから実装

2. **段階的な学習**
   - 各issueは前のissueの知識を基に構築
   - 基礎→中級→上級の順序を守る

## 学習内容の詳細

### Issue #2: Bubble Tea基礎
**MVUアーキテクチャの理解**
- **Model**: アプリケーション状態を保持する構造体
- **View**: モデルに基づいてUI文字列を生成する純粋関数
- **Update**: イベント処理とモデル更新のロジック

**インターフェースの実装**
```go
type tea.Model interface {
    Init() Cmd
    Update(Msg) (Model, Cmd)
    View() string
}
```

### Issue #3: キーボードイベント処理
**tea.KeyMsgの活用**
```go
case tea.KeyMsg:
    switch msg.Type {
    case tea.KeyUp:     // 特殊キー
    case tea.KeyRunes:  // 文字キー
```

**TDD実践**
- テスト作成 → 失敗確認 → 実装 → テスト通過
- `t.Run()`によるサブテスト構造
- 型アサーション: `newModel.(counterModel)`

**状態管理パターン**
- イミュータブルな更新
- `return m, nil` で新しいモデルとコマンドを返す

### Issue #4: Lip Glossスタイリング
**スタイル定義**
```go
style := lipgloss.NewStyle().
    Foreground(lipgloss.Color("10")).
    Border(lipgloss.RoundedBorder()).
    Padding(1, 2)
```

**条件付きスタイリング**
- 数値の正負に応じた色分け
- 動的なスタイル切り替え

**レイアウト制御**
- Padding（内側余白）とMargin（外側余白）
- Border（境界線）による視覚的階層

**設計原則**
- 内容と見た目の分離
- コンポーネント別スタイル定義
- 一貫性のあるカラーパレット

### Issue #5: 非同期処理とタイマー
**tea.Tickの活用**
```go
// 定期実行の設定
tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
    return tickMsg{time: t}
})
```

**時間の計測と表示**
- time.Durationの操作
- ミリ秒精度での時間表示
- formatDuration関数による可読性向上

**状態機械の実装**
- 停止→実行中→一時停止のサイクル
- 状態に応じたコマンド制御

### Issue #6: ビューポートとリスト管理
**仮想スクロールの概念**
```go
// ビューポート調整
if m.cursor >= m.viewport+m.height {
    m.viewport = m.cursor - m.height + 1
}
```

**カーソル管理**
- 境界チェックによる安全な移動
- キーボードナビゲーション（↑↓、j/k）
- 選択状態の視覚的フィードバック

### Issue #7: Bubblesコンポーネント統合
**textinputの使用**
```go
ti := textinput.New()
ti.Placeholder = "例: 山田太郎"
ti.Focus()
ti.CharLimit = 50
```

**フォーカス管理**
- 複数フィールド間の移動
- Tab/Shift+Tab での順次移動
- フィールド状態に応じたスタイル変更

**バリデーション設計**
- 入力チェックの分離
- エラーメッセージの表示制御
- 正規表現によるメール検証

### Issue #8: HTTP API統合と非同期処理
**tea.Cmd での HTTP通信**
```go
func fetchGitHubUser(username string) tea.Cmd {
    return func() tea.Msg {
        // HTTP通信処理
        return apiResponse{user: &user, err: nil}
    }
}
```

**状態管理の拡張**
- 入力→ローディング→成功/エラーの状態遷移
- リトライ機能の実装
- 状態に応じたUI切り替え

**spinnerコンポーネント**
```go
sp := spinner.New()
sp.Spinner = spinner.Dot
sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
```

**エラーハンドリング戦略**
- ネットワークエラーの分類
- HTTPステータスコード別処理
- ユーザーフレンドリーなエラーメッセージ

**JSON構造体設計**
- GitHub API レスポンスのマッピング
- オプショナルフィールドの条件付き表示
- time.Time のフォーマット処理

## 注意事項

- Bubble Teaアプリケーションは実際のターミナルで実行する必要がある（TTYが必要）
- 複数ファイル実行時は `go run .` または `go run main.go counter.go`
- テストは `go test -v ./...` で実行