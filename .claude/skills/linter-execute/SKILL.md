---
name: linter-execute
description: テスト通過後のlint確認・修正手順を定義するスキル。docker compose の app コンテナ内で Go 標準ツールチェーン（gofmt / go vet）を実行する。TDD実装フロー内でテストPASS後に使用する。
user-invocable: false
allowed-tools: Bash
---

# Linter実行スキル

## 実行タイミング

テストがPASSした後、コミット前にlint確認を行う。

## 前提

- 本プロジェクトはリポジトリ直下が単一の Go モジュール（フロントエンドは無い）
- lint は `docker compose` の `app` コンテナ内で実行する。ホストに Go を入れない構成のため、ホストで直接叩かない
- `Makefile` によるラッパーコマンドは**未整備**。整備された場合は本スキルの手順をそれに置き換える

## lint実行

リポジトリルートで下記コマンドを順に実行する。

```bash
# 1. フォーマット確認・修正
docker compose exec app gofmt -l -w .

# 2. vet（静的解析）
docker compose exec app go vet ./...
```

`gofmt -l` が出力したファイルはフォーマットが変更されたものなので、差分を確認すること。

## 追加のlinter（任意・未導入）

`golangci-lint` は**まだ導入されていない**。導入済みの環境で実行する場合のみ、上記の後に以下を実行する。

```bash
golangci-lint run ./...
```

未導入の状態で無理に実行しようとしないこと。導入の要否は設計段階で決める。

## ビルド確認

まず `docker compose exec app go build ./...` で全パッケージがビルドできることを確認する。

`cmd/` が追加される第2段階以降は、Lambda 向けバイナリがビルドできることも併せて確認する（`provided.al2023` ランタイム、バイナリ名は `bootstrap` 固定）。

```bash
docker compose exec app env GOOS=linux GOARCH=arm64 CGO_ENABLED=0 \
  go build -tags lambda.norpc -o bootstrap ./cmd/publisher
```

`cmd/sender` が実装済みであれば同様にビルド確認する。

## 制約事項

- `gofmt -w` → `go vet` の順で実行する
- lintエラーの場合はlintを無効化するアノテーション（`//nolint` 等）を使用せず、根本的な問題解決を行う
- ビルド成果物（`bootstrap`）をコミットに含めない
