---
name: task-analysis
description: タスク分解・調査フェーズで使用する分析スキル。Issue仕様の構造化、コードベース調査の手法、タスク粒度の判断基準を定義する。Issue分析・コード調査・タスク設計時に使用する。
user-invocable: false
allowed-tools: Read, Glob, Grep, TaskCreate, TaskUpdate, TaskGet, TaskList
---

# タスク分析スキル

## Issue仕様の構造化

Issue仕様を以下の観点で構造化して把握する。

### 機能要件の抽出

- **目的**: この変更で何を達成するか
- **ユーザーストーリー**: 誰が、何を、なぜ行うか
- **受け入れ条件**: 完了の判定基準
- **制約条件**: 技術的・ビジネス的な制約

### 技術要件の抽出

- **変更種別**: 新規追加 / 修正 / 削除 / リファクタリング
- **影響範囲**: ドメイン（`internal/domain`） / ユースケース（`internal/usecase`） / アダプタ（`internal/adapter/{weather,rds,notifier}`） / エントリポイント（`cmd/{publisher,sender}`） / インフラ（`infra/`）
- **対象 Lambda**: 発行側 / 送信側 / 両方（`internal/` は両方から参照される点に注意）
- **データモデル変更**: RDS スキーマ変更の有無
- **イベントスキーマ変更**: EventBridge カスタムバスに流すイベントの `payload` 変更の有無
- **外部依存**: 天気 API / SES / SQS / Secrets Manager / SSM Parameter Store のいずれに触れるか
- **実装ステップ**: 設計書 §9 の第何段階に相当するか（段階を飛ばすタスクは立てない）

## コードベース調査の手法

### ドメイン概念の検索

1. Issue内のキーワードからドメイン概念を抽出する
2. `internal/domain` 配下で関連する値オブジェクト・判定ロジックを検索する
3. リポジトリインターフェース（`UserRepository` 等）を確認する
4. `internal/usecase` のユースケースを確認する

主要なドメイン語彙: `Forecast`（予報）/ `City`・`CityID` / `User`・`UserID` / `RainHours`（降雨時間帯）/ `NotificationEvent` / `NotificationID` / `Channel`（email、将来 push）/ `PrecipitationProbability`。設計書と異なる語を新たに導入しない。

> **注意**: 本リポジトリはまだ実装コードが存在しない（設計書と `CLAUDE.md` のみ）。既存コードが見つからないのは異常ではないので、その場合は設計書 `heavy-weather-architecture.md` の該当セクションを一次情報として調査する。

### 既存パターンの調査

1. 類似機能の実装を見つける（同じ種類のエンティティやユースケース）
2. ディレクトリ構造とファイル命名のパターンを把握する
3. テストの書き方のパターンを把握する
4. DI・モジュール構成のパターンを把握する

### 影響範囲の調査

1. 変更対象ファイルのimport/依存関係を追跡する
2. RDS スキーマの関連テーブルを確認する
3. パッケージ間の依存関係を確認する（特に `internal/domain` が AWS SDK / `aws-lambda-go` に依存していないこと）
4. `internal/` の変更は発行側・送信側の**両 Lambda に影響する**ため、両方の呼び出し元を確認する

## タスク粒度の判断基準

### 適切な粒度

- 1つのテストファイルで検証できる範囲
- 1つの層（domain / usecase / adapter / cmd）内で完結する変更
- diffが概ね200行以内に収まる変更
- レビュアーが短時間で理解できる変更

### 粒度が大きすぎるサイン

- 複数の層にまたがる変更が1タスクに含まれている
- テストケースが10個以上必要になる
- 変更ファイルが5個以上ある

### 粒度が小さすぎるサイン

- 単体ではテストできない（依存する実装が別タスク）
- 変更がインターフェース定義のみで振る舞いがない
- 他のタスクと必ずセットでないと意味がない
