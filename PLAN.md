# `gr` - Git Worktree Review CLI

## Context

PRレビュー時に毎回手動でworktreeを切ってエディタを開く作業を自動化するCLIツール。
`git-wt`のラッパーとして、ブランチ名を渡すだけでレビュー環境を即座にセットアップする。

## 概要

- **コマンド名**: `gr`
- **モジュール**: `github.com/nu0ma/gr`
- **入力**: ブランチ名
- **動作**: worktree作成 → Cursorで自動オープン
- **クリーンアップ**: `gr finish` で worktree 削除
- **依存**: `git-wt`, `git`

## CLI設計

```
gr <branch>               # worktree作成 + Cursor起動
gr finish [branch]         # worktree削除（引数なし=現在のworktreeを削除）
gr list                    # アクティブなレビューworktree一覧
```

### `gr <branch>` の動作フロー

1. `git-wt` の存在確認
2. `git fetch origin <branch>` でリモートブランチを取得
3. 既存worktreeの確認 → 既にあればエラーで終了
4. `git-wt <branch>` でworktree作成（パスを取得）
5. `cursor <worktree-path>` でCursorを起動（シェルコマンドとして実行）
6. レビューセッション情報を状態ファイルに記録

### `gr finish [branch]` の動作フロー

1. 引数なし: cwd からworktreeを特定
2. `git-wt -d <branch>` でworktree削除
3. state ファイルからセッション削除

### `gr list` の動作フロー

1. 状態ファイルを読み込み
2. アクティブなレビュー一覧を表示（ブランチ名、パス、開始日時）

## 備考: git-wt のworktree命名

- ベースディレクトリ: デフォルト `.wt/`（`git config wt.basedir` で変更可能）
- ディレクトリ名: ブランチ名がそのまま使われる
- 例: `git wt feature-branch` → `.wt/feature-branch/`

## プロジェクト構造

```
gr/
├── go.mod
├── main.go                  # エントリポイント → cmd.Execute()
├── cmd/
│   ├── root.go              # ルートコマンド（引数あり→review, なし→list）
│   ├── review.go            # reviewロジック（デフォルトアクション）
│   ├── finish.go            # finishサブコマンド
│   └── list.go              # listサブコマンド
├── internal/
│   ├── gitwt/
│   │   └── gitwt.go         # git-wt外部コマンドラッパー
│   ├── gitutil/
│   │   └── gitutil.go       # git操作ユーティリティ（fetch, rev-parse等）
│   └── state/
│       ├── state.go         # レビューセッション状態管理
│       └── state_test.go    # テスト
├── Makefile
└── .gitignore
```

## 実装の詳細

### 1. `main.go`
- `cmd.Execute()` を呼ぶだけ

### 2. `cmd/root.go`
- cobra でルートコマンド定義
- `PersistentPreRunE` で `git-wt` の存在確認（`exec.LookPath`）
- 引数ありの場合は `review` に委譲
- `finish`, `list` サブコマンドを登録

### 3. `cmd/review.go`
- 位置引数: `<branch>`
- ロジック:
  - `gitwt.Exists(branch)` で既存worktree確認 → あればエラー終了
  - `gitutil.Fetch(origin, branch)` で fetch
  - `gitwt.Create(branch)` でworktree作成、パス取得
  - `exec.Command("cursor", worktreePath).Run()` でCursor起動
  - `state.Add(session)` でセッション記録

### 4. `cmd/finish.go`
- オプション位置引数: `[branch]`
- `--force` / `-D` フラグ
- 引数なし: `state.FindByCwd(cwd)` で現在のworktreeを特定
- `gitwt.Remove(branch, force)` で削除
- `state.Remove(branch)` でセッション削除

### 5. `cmd/list.go`
- `state.Load()` → 一覧表示
- 存在しないworktreeのエントリは自動削除（stale cleanup）

### 6. `internal/gitwt/gitwt.go`
- `Create(ctx, branch) (string, error)` - `git-wt <branch>` を実行、stdout からパス取得
- `Remove(ctx, branch, force) error` - `git-wt -d/-D <branch>`
- `Exists(ctx, branch) (bool, error)` - `git-wt --json` で既存worktree確認

### 7. `internal/gitutil/gitutil.go`
- `Fetch(ctx, remote, branch) error` - `git fetch <remote> <branch>`
- `GitCommonDir(ctx) (string, error)` - `git rev-parse --git-common-dir`
  - worktree内から実行しても正しい共通 `.git` ディレクトリを返す
  - 状態ファイルのパス解決に使用

### 8. `internal/state/state.go`
- ファイルパス: `<git-common-dir>/gr-state.json`
  - `git rev-parse --git-common-dir` で取得したパスに配置
  - worktree内からでもメインリポジトリの `.git` を参照できる
- `ReviewSession`: Branch, WorktreePath, StartedAt
- `Load(commonDir) (*State, error)`
- `Save(commonDir) error`
- `Add(session)`, `Remove(branch)`, `FindByCwd(cwd)`

## 主な依存パッケージ

- `github.com/spf13/cobra` - CLI フレームワーク
- `os/exec` - 外部コマンド実行（git-wt, git, cursor）
- `encoding/json` - state管理
- 標準ライブラリのみ（cobra以外は外部依存なし）

## 実装順序

1. プロジェクト初期化（go.mod, main.go, .gitignore）
2. `internal/gitutil` - git操作ユーティリティ
3. `internal/gitwt` - git-wtラッパー
4. `internal/state` - 状態管理 + テスト
5. `cmd/root.go` + `cmd/review.go` - メインフロー
6. `cmd/finish.go` - クリーンアップ
7. `cmd/list.go` - 一覧表示
8. Makefile

## 検証方法

1. `go build -o gr .` でビルド確認
2. `go test ./...` でユニットテスト実行
3. 実際のgitリポジトリ内で以下を手動テスト:
   - `gr <ブランチ名>` → worktree作成 + Cursor起動を確認
   - `gr <ブランチ名>` 2回目 → エラー終了を確認
   - `gr list` → アクティブレビュー一覧表示を確認
   - `gr finish <ブランチ名>` → worktree削除を確認
