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
- Issue #2: Hello Worldアプリケーション（main.go）

### 次のステップ（Issue #3）
カウンターアプリケーションの実装：
- 上下キーでカウントの増減
- スペースキーでリセット
- TDDアプローチ：まずテストを書いてから実装

## 開発方針

1. **TDD（テスト駆動開発）を原則とする**
   - 新機能実装前にテストを作成
   - テストが失敗することを確認してから実装

2. **段階的な学習**
   - 各issueは前のissueの知識を基に構築
   - 基礎→中級→上級の順序を守る

## 注意事項

- 現在テストファイルは存在しないが、Issue #3からTDDを開始予定
- Bubble Teaアプリケーションは実際のターミナルで実行する必要がある（TTYが必要）
- 各issueブランチは`feature/`プレフィックスを使用