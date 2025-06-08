# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

これはBubble TeaフレームワークとGo言語を段階的に学習するためのプロジェクトです。基礎から上級まで9つのissueを通じて実践的に学習します。

## 開発コマンド

```bash
# アプリケーションの実行
go run main.go

# 依存関係のダウンロード
go mod download

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

### 次のステップ
Issue #5以降の中級編に進む準備完了

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

## 注意事項

- Bubble Teaアプリケーションは実際のターミナルで実行する必要がある（TTYが必要）
- 複数ファイル実行時は `go run .` または `go run main.go counter.go`
- テストは `go test -v ./...` で実行