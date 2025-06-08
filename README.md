# Bubble Tea Learning 🫧

Bubble TeaフレームワークとGo言語を段階的に学習するためのプロジェクトです。

## 概要

このプロジェクトは、[Bubble Tea](https://github.com/charmbracelet/bubbletea)を使ったターミナルUIアプリケーション開発を基礎から学ぶことを目的としています。段階的なissueを通じて、シンプルなHello Worldから複雑なレイアウトまで、実践的に学習できます。

## 必要環境

- Go 1.23.0以上
- ターミナル環境

## セットアップ

```bash
# リポジトリのクローン
git clone https://github.com/ktny/bubbletea-learning.git
cd bubbletea-learning

# 依存関係のインストール
go mod download
```

## 実行方法

```bash
# Hello Worldアプリケーションの実行
go run main.go
```

## 学習計画

### 🟢 基礎編
- [x] **#2** [Hello World - 最初のBubble Teaアプリケーション](https://github.com/ktny/bubbletea-learning/issues/2)
  - Model-View-Updateパターンの理解
  - 基本的なアプリケーション構造
- [ ] **#3** [ユーザー入力 - キーボードイベントの処理](https://github.com/ktny/bubbletea-learning/issues/3)
  - キーボード入力の処理
  - 状態管理の基礎
- [ ] **#4** [スタイリング - Lip Glossでターミナルを彩る](https://github.com/ktny/bubbletea-learning/issues/4)
  - Lip Glossを使ったスタイリング
  - 色、枠線、レイアウトの適用

### 🟡 中級編
- [ ] **#5** [コマンドとメッセージ - 非同期処理の実装](https://github.com/ktny/bubbletea-learning/issues/5)
  - tea.Cmdを使った非同期処理
  - カスタムメッセージの定義
- [ ] **#6** [リスト表示 - 選択可能なメニューの作成](https://github.com/ktny/bubbletea-learning/issues/6)
  - 選択可能なリストUI
  - スクロールとカーソル管理
- [ ] **#7** [テキスト入力 - フォームの実装](https://github.com/ktny/bubbletea-learning/issues/7)
  - Bubblesライブラリの活用
  - フォーム入力とバリデーション

### 🔴 上級編
- [ ] **#8** [API連携 - HTTPリクエストとプログレス表示](https://github.com/ktny/bubbletea-learning/issues/8)
  - 外部APIとの連携
  - プログレス表示の実装
- [ ] **#9** [複雑なレイアウト - 分割画面とタブシステム](https://github.com/ktny/bubbletea-learning/issues/9)
  - 複雑なレイアウトの構築
  - 複数ビューの管理

## プロジェクト構造

```
.
├── README.md        # このファイル
├── go.mod          # Goモジュール定義
├── go.sum          # 依存関係のチェックサム
└── main.go         # メインアプリケーション
```

## 開発方針

- **TDD（テスト駆動開発）**: 各機能の実装前にテストを作成
- **段階的学習**: 基礎から上級まで順番に学習
- **実践的**: 実際に動くアプリケーションを作成

## 参考リソース

- [Bubble Tea 公式ドキュメント](https://github.com/charmbracelet/bubbletea)
- [Bubble Tea チュートリアル](https://github.com/charmbracelet/bubbletea/tree/master/tutorials)
- [Lip Gloss スタイリング](https://github.com/charmbracelet/lipgloss)
- [Bubbles コンポーネント](https://github.com/charmbracelet/bubbles)

## ライセンス

このプロジェクトは学習目的で作成されています。

---

💡 **改善提案**: 各issueを完了するごとに、学んだことをこのREADMEに追記していくと、学習の振り返りに役立ちます。