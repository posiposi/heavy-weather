---
name: tdd-workflow
description: テスト駆動開発（TDD）のワークフロー定義。テスト実行は docker compose の app コンテナ内（go test）で行う。Red→Green→Refactorのサイクル手順、テストの書き方、層別テスト方針との統合方法を定義する。テスト実装・実行時に使用する。
user-invocable: false
allowed-tools: Read, Write, Edit, Glob, Grep, Bash, TaskUpdate, TaskGet, TaskList
---

# TDDワークフロースキル

## TDDサイクル

### 1. Red フェーズ（テストを先に書く）

#### 手順

1. テストファイルを作成する（`*_test.go`）
2. テスト対象の期待する振る舞いをテストケースとして記述する
3. テストを実行し、**失敗すること**を確認する

#### テストの書き方

```go
func TestNewPrecipitationProbability(t *testing.T) {
    tests := []struct {
        name    string
        input   int
        wantErr bool
    }{
        {"0%は有効", 0, false},
        {"100%は有効", 100, false},
        {"範囲内の値で生成できる", 60, false},
        {"負の値は不正", -1, true},
        {"100を超える値は不正", 101, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := NewPrecipitationProbability(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("NewPrecipitationProbability(%d) error = %v, wantErr %v", tt.input, err, tt.wantErr)
                return
            }
            if !tt.wantErr && got.Value() != tt.input {
                t.Errorf("NewPrecipitationProbability(%d).Value() = %v, want %v", tt.input, got.Value(), tt.input)
            }
        })
    }
}
```

#### テスト名

- `t.Run` のテスト名は**日本語**で振る舞いを簡潔に表現する
- Table-driven tests を活用しテストケースの追加を容易にする

#### 層別のテスト方針

| 層 | パッケージ | テスト対象 | モック対象 |
|---|---|---|---|
| ドメイン層 | `internal/domain` | 値オブジェクト（`Forecast` / `PrecipitationProbability` / `RainHours` 等）、降水判定ロジック | なし（純粋なロジック。AWS SDK に依存しない） |
| ユースケース層 | `internal/usecase` | 処理の組み立て（ユーザー取得 → 都市の重複排除 → 予報取得 → イベント組み立て） | リポジトリ・外部APIクライアントのインターフェース |
| アダプタ層 | `internal/adapter/{weather,rds,notifier}` | 天気APIクライアント、`UserRepository` 実装、SES送信 | 外部依存（HTTPは固定レスポンス、DBは実接続） |
| エントリポイント | `cmd/{publisher,sender}` | ハンドラの配線 | ユースケース |

**最優先で厚くテストする領域は `internal/domain`。** AWS に一切依存せず、テストが最も書きやすい。特に以下の境界値を必ずカバーする。

- 日をまたぐ降雨（連続時間帯が日付境界を越えるケース）
- 降水確率が閾値**ちょうど**のケース（閾値以上を降雨とみなす仕様の境界）
- データ欠損（天気APIレスポンスに時間帯や確率が含まれないケース）

また、`notificationId` は入力から決定的に生成される仕様のため、**同じ入力から同じIDが生成されること**をテストで固定する。

#### 外部依存を伴うテストの方針

- 天気APIクライアントのテストは実APIを叩かず、**固定のJSONレスポンス**（テストデータ）に対して変換ロジックを検証する
- RDS に実接続するテストを書く場合はモックを使用せず、テストデータの投入・クリーンアップを各テストコードで必ず行い、テスト間の独立性を保つ
- 外部依存を伴うテストは `//go:build integration` タグで区別し、通常の `go test ./...` から外す
- **テスト用DB環境は現時点で未整備**（RDS 導入は実装ステップ第3段階）。それまでは統合テストを前提にしたタスクを立てない

### 2. Green フェーズ（最小限の実装）

#### 手順

1. テストが通る**最小限**のコードを実装する
2. 「動くコード」を最優先する
3. テストを実行し、**すべて通ること**を確認する

#### 原則

- 完璧なコードを書こうとしない
- テストに書かれていない機能は実装しない
- 既存のテストが壊れていないことも確認する

### 3. Refactor フェーズ（品質改善）

#### 手順

1. コードの重複を排除する
2. 命名を改善する（ドメイン語彙に準拠: `Forecast` / `City` / `RainHours` / `NotificationEvent` 等）
3. 不変性、ファクトリ関数等の値オブジェクト規約への適合を確認する
4. テストを再実行し、**すべて通ること**を確認する

#### リファクタリング観点

- 単一責任の原則に違反していないか
- 層の責務が正しく分離されているか
- 命名がドメインの言葉を使っているか
- 不要なコメントや死んだコードがないか

#### アーキテクチャ上の不変条件（CLAUDE.md より。破っていないか毎回確認する）

- `internal/domain` が `aws-lambda-go` を import していないこと
- 通知イベントに整形済みの本文を載せていないこと（本文整形は送信側アダプタの責務）
- `notificationId` をランダム生成していないこと（入力から決定的に生成する）

## テスト実行コマンド

`docker compose` の `app` コンテナ内で実行する。ホストに Go を入れない構成のため、ホストで直接叩かない。

```bash
# 特定パッケージのテスト実行
docker compose exec app go test ./internal/domain/...

# 全テスト実行
docker compose exec app go test ./...

# 詳細出力
docker compose exec app go test -v ./internal/domain/...

# integrationテスト実行（外部依存を伴うテストを書いた場合のみ）
docker compose exec app go test -tags=integration ./...
```

`Makefile` によるラッパーコマンド（`make test` 等）は**未整備**。整備された場合は本スキルの手順をそれに置き換える。

## テスト失敗時の対応

1. エラーメッセージを正確に読む
2. 期待値と実際の値の差分を確認する
3. Greenフェーズの実装を修正する（テストは変更しない）
4. テストを再実行する
5. 3回修正しても通らない場合は、テストの前提条件を見直す

## テストケースの原則

- テストケースは**そのコード固有の振る舞い（ドメインロジック・バリデーション・状態遷移等）のみ**を検証する
- 言語機能そのもの（基本的な演算、型変換等）はテストしない
- コード固有のロジックが絡むケースのみテストする
  - 例: 正規化（トリム、大文字小文字統一等）が等価性比較に影響するケース
  - 例: バリデーションルール（文字数制限、フォーマット検証等）
  - 例: ファクトリ関数での入力変換ロジック

## テストのアンチパターン

- テストが実装の詳細に依存している（unexportedな関数のテスト等）
- テスト間に順序依存がある
- テストデータが共有されている（各テストは独立すべき）
- モックが多すぎる（設計の問題のサイン）
- テスト名が意味のない名前になっている
