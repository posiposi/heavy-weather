---
paths:
  - "**/*.go"
---

# Go DDD Architecture Rules

本プロジェクトはリポジトリ直下を単一の Go モジュールとし、ドメイン駆動設計（DDD）+ CQRS パターンに従う。「Lambda のプロジェクト」ではなく「Lambda にもデプロイできる Go アプリケーション」として扱う。

## Layer Structure (Go)

| Layer | パッケージ | 責務 |
|-------|-----------|------|
| ドメイン層 | `internal/domain` | ビジネスロジックの核。値オブジェクト・エンティティ・リポジトリ抽象・降水判定などのドメインサービス |
| アプリケーション層 | `internal/usecase` | ユースケース（制御フロー）。取得したデータの組み立て、通知イベントの生成 |
| インフラ層（アダプタ） | `internal/adapter/*` | 技術的実装詳細。外部天気 API クライアント・リポジトリ実装・SES 送信 |
| エントリポイント | `cmd/*` | Lambda ハンドラの登録と依存の配線のみ |

インターフェース層（HTTP ハンドラ・ルーティング）は現時点では存在しない。将来管理画面を追加する場合に `internal/interfaces/` として切り出す。

## Directory Structure (Go)

```
├── cmd/
│   ├── publisher/main.go       # 発行側 Lambda のエントリポイント
│   └── sender/main.go          # 送信側 Lambda のエントリポイント
├── internal/
│   ├── domain/                 # 値オブジェクト、エンティティ、リポジトリ抽象、降水判定ロジック
│   ├── usecase/                # 処理の組み立て
│   └── adapter/
│       ├── weather/            # 天気 API クライアント（ForecastRepository の実装）
│       ├── rds/                # UserRepository 等の実装
│       └── notifier/           # SES 送信
├── infra/                      # IaC（SAM / CDK / Terraform）
├── Makefile
└── go.mod
```

- **ドメイン層を Lambda に依存させない。** `internal/domain` が `aws-lambda-go` を import していたら設計が誤っているという判定基準になる。AWS SDK への依存も同様にアダプタ層に閉じる
- **Lambda を分割してもロジックは共有する。** 分割されるのは `cmd/` 配下のみで、`internal/` は発行側・送信側の両方から参照される
- 現時点の規模では境界づけられたコンテキストごとのモジュール分割は行わない。ドメインが肥大化した場合に `internal/domain/{context}/` へ分割する
- 汎用パッケージ（`util` / `common` / `helper`）は作らない

## File Naming Convention (Go)

| Layer | Category | Naming Pattern | Example |
|-------|----------|----------------|---------|
| domain | Entity | `{entity}.go` | `user.go`, `notification.go` |
| domain | Value Object | `{name}.go` | `city_id.go`, `email.go`, `temperature.go`, `precipitation_probability.go` |
| domain | Enum/定数 | `{name}.go` | `channel.go`, `notification_type.go` |
| domain | ドメインサービス | `{name}.go` | `rain_hours.go` |
| domain | 共通基盤 | `{name}.go` | `error.go`, `id.go` |
| domain | Repository Interface (Query) | `{entity}_query_repository.go` | `user_query_repository.go` |
| domain | Repository Interface (Command) | `{entity}_command_repository.go` | `notification_command_repository.go` |
| adapter/{tech} | Repository Implementation | `{entity}_{query\|command}_repository.go` | `user_query_repository.go` |
| adapter/rds | DB モデル | `{entity}_model.go` | `user_model.go` |
| adapter/{tech} | テスト | `{対象ファイル名}_test.go` | `user_query_repository_test.go` |

インターフェースと実装はファイル名・型名とも同一にし、パッケージ名（`domain` / `rds`）で区別する。

```go
// internal/domain/user_query_repository.go
package domain
type UserQueryRepository interface { ... }

// internal/adapter/rds/user_query_repository.go
package rds
type UserQueryRepository struct { db *sql.DB }

func NewUserQueryRepository(db *sql.DB) domain.UserQueryRepository { ... }
```

### DTO変換の方針

- DTO（外部入出力の変換）はアプリケーション層（ユースケース）で行う。独自のmapperメソッドは不要
- 外部天気 API のレスポンス → ドメインの `Forecast` への変換はアダプタ層（`adapter/weather`）で行う。API 固有の JSON 構造をドメインに漏らさない
- 通知イベントの JSON へのシリアライズはアダプタ層の責務。ドメインは意味的なデータ（都市・日付・気温・降雨時間帯）のみを持つ
- Update操作はドメインモデル内にCommand構造体を定義して責務を寄せる

```go
type CommandUpdateSubscription struct {
    Cities   []CityID
    Channels []Channel
}
```

## Query/Command Separation (CQRS)

リポジトリインターフェースは原則として Query と Command に分離する。

### Query Repository

- 読み取り専用の操作を定義する
- データを変更しない（副作用なし）
- `Find`, `Get`, `List` 等の動詞を使用する

```go
type UserQueryRepository interface {
    FindByID(ctx context.Context, id UserID) (*User, error)
    ListSubscribers(ctx context.Context, t NotificationType) ([]*User, error)
}
```

### Command Repository

- 書き込み・状態変更の操作を定義する
- `Save`, `Delete`, `Update` 等の動詞を使用する

```go
type NotificationCommandRepository interface {
    // 冪等性レコードの登録。既存なら inserted=false を返す（INSERT ... ON CONFLICT DO NOTHING）
    TryInsertSending(ctx context.Context, id NotificationID, ch Channel) (inserted bool, err error)
    MarkSent(ctx context.Context, id NotificationID, ch Channel, messageID string) error
}
```

### Mixed Repository (例外的)

- 集約が小さく Query/Command 分離が過剰な場合のみ許容する
- 例: 予報のように「取得」と「キャッシュへの保存」が密結合な場合。外部 API を直接叩く実装から、事前取得したキャッシュを読む実装へ差し替えられるよう `ForecastRepository` として抽象化しておく

```go
type ForecastRepository interface {
    FindByCity(ctx context.Context, city CityID, date Date) (*Forecast, error)
    Save(ctx context.Context, f *Forecast) error
}
```

## Repository

- リポジトリインターフェースはドメイン層（`internal/domain`）に配置する
- 実装はアダプタ層（`internal/adapter/{rds,weather,notifier}`）に配置する
- DB行・API レスポンス → ドメインエンティティの変換はリポジトリ実装内のprivate関数で行う（独立したmapperは不要）
- 独立したmapper層は設けない（現在のサービス規模では不要）
- DB タグ付き構造体や API レスポンス構造体はアダプタ層の関心事であり、ドメインエンティティとは別型にする

## Dependency Rule

- Domain 層は他の層に依存しない（リポジトリはインターフェースのみ）。AWS SDK・`aws-lambda-go`・HTTP クライアントを import しない
- Usecase 層は Domain 層に依存する
- Adapter 層は Domain 層に依存する（リポジトリインターフェースを実装）
- Adapter 層の詳細が Domain/Usecase に漏れてはならない
- `cmd/` は全層を import してよいが、配線以外のロジックを持たない
