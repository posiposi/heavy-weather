---
name: implement
description: GitHub IssueからPR作成までの開発ワークフローを実行する。Issue番号を引数として受け取り、仕様取得→調査→タスク分解→TDD実装→レビュー→コミット→PR作成を一気通貫で行う。「Issue #N を実装して」「#N をやって」のように指示された場合に使用する。
allowed-tools: Read, Write, Edit, Glob, Grep, Bash, Skill, Agent, WebFetch, TaskCreate, TaskUpdate, TaskList, TaskGet
user-invocable: true
---

# 概要

引数で受け取ったGitHub Issue番号（以下 `<Issue番号>`）の仕様に基づき、以下のフェーズを順番に実行する。
各フェーズ間の情報連携はClaude CodeのTasks機能を使用する。

## Phase 0: 前提確認とブランチ作成

> **注意**: 本リポジトリは現時点で **git 初期化・リモートリポジトリの作成が未実施**の可能性がある。まず前提が成立するか確認し、成立しない場合は Phase 0 / Phase 4 をスキップしてユーザーに状況を報告すること（勝手に `git init` やリポジトリ作成は行わない）。

```bash
git rev-parse --is-inside-work-tree   # git リポジトリか
git remote -v                          # リモートがあるか
gh repo view --json name               # GitHub 上にリポジトリがあるか
```

git リポジトリとして初期化済みで、リモートが存在する場合のみ以下を実行する。

1. 現在のブランチが `main` であることを確認し、最新の状態に更新する
   ```bash
   git checkout main && git pull origin main
   ```
2. Issue番号に基づいた作業ブランチを作成・切り替える
   ```bash
   git checkout -b feature/<Issue番号>-<Issueタイトルの要約(kebab-case)>
   ```

> **重要**: `main` ブランチでは直接作業しない。すべての変更は作業ブランチ上で行うこと。

Issue 管理がまだ運用されていない場合は、Phase 1 の代わりにユーザーから直接受け取った要件を仕様として扱い、Phase 2 以降を実行する。

## Phase 1: 仕様取得

1. `gh` CLIを使用してIssue #`<Issue番号>` の内容を取得する
   ```bash
   gh issue view <Issue番号> --json title,body,labels,comments
   ```
2. Issueのタイトル、本文、ラベル、コメントを収集する
3. TaskCreateで「仕様取得」タスクを作成する
   - subject: `Issue #<Issue番号> の仕様取得`
   - metadataにIssue仕様を格納する

## Phase 2: 調査・タスク分解

### Phase 2-1: 並列調査

以下のサブエージェントを**並列実行**する：

| サブエージェント | 目的 |
|---|---|
| `code-investigator` | Issue仕様に関連する既存コード・パターン・影響範囲を調査 |
| `log-investigator` | エラーログ・テスト結果・環境情報を調査（バグ修正Issue時） |

各サブエージェントには仕様取得タスクIDと自身のタスクIDを渡す。

### Phase 2-2: タスク分解

調査完了後、`task-decomposer`エージェントを起動し、1コミット粒度の実装タスクを作成する。

## Phase 3: 実装

### Phase 3-1: 実装ループ

タスク分解で作成された各実装タスクについて、以下のサイクルを実行する：

1. **テストコード作成**: `implementer`エージェントにテストコード作成を指示
2. **RED確認**: `implementer`エージェントにテスト実行（失敗確認）を指示
3. **実装**: `implementer`エージェントにコード実装を指示
4. **GREEN確認**: `implementer`エージェントにテスト実行（成功確認）を指示
5. **レビュー**: `code-reviewer`エージェントでレビュー実行
6. **修正**（指摘がある場合）: `review-fixer`エージェントで修正、再レビュー

### Phase 3-2: コミット

各タスク完了後、`/commit-commands:commit` でコミットする。

## Phase 4: PR作成

全タスク完了後（Phase 0 で git リポジトリ・リモートの存在を確認できている場合のみ）：

1. `pr-creator`エージェントでPRタイトル・本文を生成する
2. `/commit-commands:commit-push-pr` でPRを作成する

リモートリポジトリが存在しない場合は PR 作成を行わず、ローカルコミットまでで停止してユーザーに報告する。

## 注意事項

- 各フェーズの結果はTasks機能のmetadataで連携する
- テスト・lint・ビルドはホストの Go ツールチェーンで実行する（本プロジェクトにコンテナ環境は無い）
  - テスト: `go test ./...` — 詳細は `tdd-workflow` スキル
  - lint: `gofmt -w .` → `go vet ./...` — 詳細は `linter-execute` スキル
  - ビルド: `GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap ./cmd/publisher`
- `Makefile` / `go.mod` はまだ整備されていない可能性がある。存在を確認せずにラッパーコマンドを実行しない
- 設計上の判断に迷う場合は `heavy-weather-architecture.md`（唯一の真実の源）と `CLAUDE.md` を参照する
- 実装ステップは設計書 §9 の段階順（第0〜第5段階）に従う。段階を飛ばさない
- GitHub操作には `gh` CLIを使用する（MCP サーバーは使用しない）
