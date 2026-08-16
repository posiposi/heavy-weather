# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## リポジトリの現状

**ローカル開発環境（Docker）と Go モジュールのみが存在する。** ドメインロジック以降は未実装で、`internal/domain` はパッケージ宣言だけの状態（設計書 §9 第1段階で実体を入れる）。`documents/heavy-weather-architecture.md`（設計書）が唯一の真実の源であり、以下の記述はそこからの要約である。判断に迷う場合は設計書の該当セクションを参照すること。

**`documents/` は Git 管理外である。** 手元に存在しない場合は設計書を参照できないため、設計上の判断が必要な場面ではその旨を伝えて指示を仰ぐこと。要約であるこのファイルの記述だけで判断を進めない。

コードを追加した時点で、このファイルの「コマンド」セクションを実際に動くコマンドへ更新すること。

## 開発環境

`docker compose` で Go（`golang:1.26.6-bookworm`）と MySQL（`mysql:8.4.10`）を立ち上げる。ホストに Go を入れずに済む構成であり、テスト・lint・ビルドはコンテナ内で実行する。起動手順は `README.md`「ローカル開発環境」を参照。

```sh
cp .env.example .env   # 値は空なので記入が必要（README 参照）
docker compose up -d
docker compose exec app go test ./...
docker compose exec app go vet ./...
```

Go の alpine イメージは使わない。musl libc であり、本番の Lambda ランタイム `provided.al2023`（glibc）と揃わないため。

## コメント方針

**コメントは原則として書かない。** Go・設定ファイル（`compose.yaml` / `Dockerfile` / `.env.example` など）を問わず適用する。

書いてよいのは以下のみ。

- **分岐やパラメータ選択の理由**。「なぜこの条件なのか」「なぜこの値なのか」がコードから読み取れない箇所に限る
- **Go のエクスポートされたトップレベル名の doc コメント**（godoc として機能するため）

書かないもの。

- **技術・バージョン・イメージの採用理由。** 選定の背景は `documents/heavy-weather-architecture.md`（設計書）に、運用上踏まえるべき事項は `CLAUDE.md` に、使い方は `README.md` に書く。コメントに置くと同じ根拠が複数箇所に散り、更新漏れで矛盾する
- **設定値が何であるかの言い換え。** `character-set-server = utf8mb4` に「文字セットの設定」と付ける類
- **セクション見出し代わりのコメント**

判断に迷う場合は書かない側に倒す。読み手が根拠を必要とする情報であれば、それはコメントではなくドキュメントに書くべき分量である。

## プロジェクト概要

天気通知システム。EventBridge Scheduler で毎朝起動し、RDS からユーザーの購読都市を取得 → 外部天気 API から予報を取得 → EventBridge カスタムバス経由で SQS → 送信 Lambda → SES でメール配信する、Go / AWS サーバレス構成。

## 想定リポジトリ構成

```
cmd/publisher/main.go   発行側 Lambda のエントリポイント
cmd/sender/main.go      送信側 Lambda のエントリポイント
internal/domain/        値オブジェクト、降水判定ロジック
internal/usecase/       処理の組み立て
internal/adapter/weather/  天気 API クライアント
internal/adapter/rds/      UserRepository の実装
internal/adapter/notifier/ SES 送信
infra/                  IaC（SAM / CDK / Terraform — 第2段階で選定）
```

## ビルド

Lambda 向けバイナリのビルド（`provided.al2023` ランタイム、バイナリ名は `bootstrap` 固定）:

```sh
docker compose exec app env GOOS=linux GOARCH=arm64 CGO_ENABLED=0 \
  go build -tags lambda.norpc -o bootstrap ./cmd/publisher
```

**`cmd/` はまだ存在しないため、このコマンドは第2段階でエントリポイントを追加するまで通らない。** 現時点で通るのは `go build ./...` まで。

`arm64` を使う（x86 より安価、Go では移植性の問題がほぼない）。`-tags lambda.norpc` で不要な RPC 依存を除外する。

## アーキテクチャ上の不変条件

これらは設計判断の結果であり、破ると設計意図が失われる。

- **`internal/domain` は `aws-lambda-go` を import しない。** import していたら設計が誤っているという判定基準。分割されるのは `cmd/` 配下のみで、`internal/` は両 Lambda から参照される。
- **イベントに整形済み本文を載せない。** イベントの `payload` は意味的なデータ（都市、日付、気温、降雨時間帯）のみ。本文整形はチャネルアダプタ（送信側）の責務。ここを崩すと push 通知追加時にイベントスキーマの作り直しが必要になる。
- **チャネル選択は発行側の責務。** 発行側が設定を読み `(user × 有効チャネル)` の数だけイベントを発行する。消費側で判定させない。
- **`notificationId` は入力から決定的に生成する**（例: `hash(userId, date, type)`）。ランダム生成すると再実行時に別 ID となり下流の冪等性チェックをすり抜ける。冪等キーは `notificationId` + `channel` の複合。
- **処理順序は「ユーザー取得 → 都市の重複排除 → 予報取得」。** どの都市を取るかがユーザーの購読情報に依存するため。重複排除により外部 API 呼び出し回数がユーザー数でなく都市数に比例する。

## 実装上の落とし穴

### 発行側 Lambda

- `http.DefaultClient` はタイムアウトが無限。必ず `http.Client` にタイムアウトを設定する。
- `PutEvents` は1回10件まで、かつ部分失敗があり得る。`FailedEntryCount` を確認し**失敗したエントリのみ**再送する（全件リトライは重複配信）。
- API キー・DB 認証情報は初期化フェーズで一度だけ取得しキャッシュする（Secrets Manager: DB、SSM Parameter Store: 天気 API キー）。
- 外部 API 障害時は部分成功を許容し、取得できた都市だけ配信する。

### 送信側 Lambda

イベントソースマッピング（AWS 管理のポーラー）が SQS をポーリングして Lambda を同期呼び出しする構造に由来する注意点:

- **ハンドラは `error` を返さず `nil` を返す。** エラーを返すと部分バッチレスポンスが無視され、バッチ全体が再配信されて二重送信になる。`ReportBatchItemFailures` を有効にし、失敗したメッセージ ID のみ `BatchItemFailures` に入れる。
- 可視性タイムアウトは関数タイムアウトの6倍を目安に長く取る。
- DLQ は**キュー側**に `maxReceiveCount` 付きで設定する（Lambda の `onFailure` とは別の仕組み）。
- 1件あたりの処理順: 冪等性チェック（`INSERT ... ON DUPLICATE KEY UPDATE` で「送信中」レコード、既存なら即 return）→ 本文整形 → SES 送信 → `MessageId` 保存 → 送信済みに更新。
- 冪等性の挿入は `INSERT ... ON DUPLICATE KEY UPDATE notification_id = notification_id` とし、`RowsAffected() == 1` を挿入成功の判定に使う。`INSERT IGNORE` は NOT NULL 違反やデータ切り詰めまで警告に格下げして黙って通すため使わない。

### DB（MySQL 8.4 LTS / RDS Proxy は現時点で不採用）

**バージョンは MySQL 8.4 LTS の 8.4.10 を使用する。** ローカルの Docker イメージ、IaC のエンジンバージョン指定、CI のサービスコンテナすべてで揃えること。

- **8.0 を選ばない。** RDS for MySQL 8.0 は標準サポートが 2026年7月31日に終了済みで、新規構築しても Extended Support（vCPU 時間あたりの追加課金）の対象になる。コミュニティ側も 2026年4月30日に終息しており新しいマイナーバージョンは出ない。
- **9系も選ばない。** 次期 LTS の 9.7（2026年4月 GA）は RDS では Database Preview 環境のみの提供で、本番利用が禁じられている（インスタンスが作成から60日で自動削除される、RDS Proxy が使えない、サポート対象外）。RDS で一般提供され次第、Blue/Green デプロイでの移行を検討する。
- 8.4 の RDS 標準サポート終了は 2029年7月31日、Extended Support 終了は 2032年7月31日。
- 文字セットは `utf8mb4`、照合順序は `utf8mb4_0900_ai_ci`。ストレージエンジンは InnoDB（外部キー制約を使うため MyISAM は不可）。
- テーブル定義は `documents/heavy-weather-db-schema.md` を参照する。

接続の扱い:

- `sql.DB` はパッケージレベル変数として保持する（`sql.DB` 自体がコネクションプール）。
- `SetMaxOpenConns` を小さく設定する。
- `SetConnMaxIdleTime` を DB の `wait_timeout` より短くする（Lambda 凍結中に切断された接続を掴む問題の回避）。
- Lambda に予約された同時実行数を設定する。SQS の `maximumConcurrency` で SES のレート制限に合わせた流量調整も行う。

## 実装の進め方

設計書 §9 の段階順に進める。段階を飛ばすと切り分けが困難になる。

1. **第0段階（最優先・並行）**: SES 本番アクセス申請、送信ドメイン認証（DKIM/SPF/DMARC）、天気 API の実地検証。AWS 側の審査待ちが発生するため先に着手する。
2. **第1段階**: ドメインロジック（AWS 非依存の Go のみ）。天気 API レスポンス → `Forecast` 値オブジェクト、降水確率が閾値以上の連続時間帯を1区間にまとめる判定、通知イベント組み立て。ユーザーはハードコードでよいが `[]User` を受け取るシグネチャにしておく。境界値テスト（日をまたぐ降雨、確率が閾値ちょうど、データ欠損）を書く領域。
3. **第2段階**: `EventBridge Scheduler → Lambda → SES` を縦に1本通す。RDS も SQS も入れない。**このタイミングで IaC ツールを決定する**（後からの乗り換えは手戻りが大きい）。
4. **第3段階**: RDS と VPC。ネットワーク問題はタイムアウトするだけでエラーが出ないため、第2段階で Lambda 自体は動くと確認済みにしておく。
5. **第4段階**: EventBridge バス + SQS で非同期化。ここで分解が困難なら層分けが不十分だったというフィードバック。
6. **第5段階**: 異常系（冪等性テーブル、部分バッチ、DLQ とアラーム、SES バウンス/苦情によるチャネル無効化）。バウンス処理なしで送信を続けると SES の評判指標悪化でアカウント単位の送信停止リスクがある。

## 未確認事項

着手前に一次情報での確認が必要（設計書 §11）。これらの前提が崩れると設計から見直しになるものがある。

- 天気 API（Open-Meteo 想定）のレスポンス仕様、降水確率が実際に1時間粒度で返るか、レート制限、商用利用の可否、タイムゾーンの扱い
- RDS Proxy / Aurora 系の最新料金、Aurora DSQL の現在のサポート機能
- SES の送信レート上限（本番アクセス承認時に決定）
- ユーザー管理が社内ツール型（管理者が他人の設定を管理）かセルフサービス型か。Cognito と通知対象の関係、データモデル、push 通知のエンドポイント登録に影響する
