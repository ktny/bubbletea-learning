# 🫧 Bubble Tea Learning Journey

段階的にBubble TeaフレームワークとGo言語を学習するためのプロジェクトです。基礎から上級まで9つのissueを通じて実践的に学習します。

## 📚 学習進捗

### ✅ 完了済み
- **Issue #2**: Hello Worldアプリケーション - Bubble TeaのMVUアーキテクチャ基礎
- **Issue #3**: ユーザー入力 - キーボードイベント処理とTDD実践
- **Issue #4**: スタイリング - Lip Glossによる条件付きスタイリング
- **Issue #5**: コマンドとメッセージ - 非同期処理とタイマー実装
- **Issue #6**: リスト表示 - 選択可能なTODOリストとビューポート機能

### 🚧 進行中
次に取り組む予定のissue

### 📋 今後の学習
- **Issue #7**: テキスト入力 - Bubblesライブラリを使ったフォーム実装
- **Issue #8**: API連携 - HTTPリクエストとプログレス表示
- **Issue #9**: 複雑なレイアウト - 分割画面とタブシステム

## 🚀 アプリケーション実行

```bash
# カウンターアプリ（基本的なキーボード操作）
go run . counter

# タイマーアプリ（非同期処理とミリ秒精度）
go run . timer

# TODOリストアプリ（リスト操作とビューポート）
go run . todo
```

## 🧪 テスト実行

```bash
# 全てのテストを実行
go test -v ./...

# カバレッジ付きでテスト実行
go test -cover ./...
```

## 🏗️ アーキテクチャ

### Model-View-Update (MVU)パターン
```go
type Model interface {
    Init() Cmd           // 初期化処理
    Update(Msg) (Model, Cmd)  // イベント処理
    View() string        // UI描画
}
```

### 学習した技術要素
- **基礎**: MVUアーキテクチャ、インターフェース実装
- **イベント処理**: tea.KeyMsg、型アサーション
- **スタイリング**: Lip Gloss、条件付きスタイル
- **非同期処理**: tea.Cmd、tea.Tick、カスタムメッセージ
- **リスト管理**: ビューポート、スクロール、カーソル移動
- **TDD**: テスト駆動開発、サブテスト構造

## 📖 ドキュメント

各実装の詳細な学習内容は `CLAUDE.md` を参照してください。

## 🛠️ 開発環境

- Go 1.21+
- Bubble Tea v0.27+
- Lip Gloss v0.9+

---

**学習目標**: 実用的なターミナルUIアプリケーションの構築スキルを段階的に習得