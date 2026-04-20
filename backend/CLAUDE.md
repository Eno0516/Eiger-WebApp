# バックエンド作業ルール

## アーキテクチャ

- **Clean Architecture + DDD**
- 依存の方向: `handler → usecase → domain ← repository`

## 構成

```
backend/
├── cmd/server/          # エントリーポイント
├── gen/                 # 自動生成コード（手を加えない）
├── internal/
│   ├── controller/      # gen の interface を実装する（HTTP の関心事なし）
│   ├── usecase/         # ビジネスロジック
│   ├── domain/          # Entity / ValueObject / Repository interface
│   ├── handler/         # oapi-codegen が生成するルーティング層
│   ├── repository/      # Repository interface の実装
│   └── infrastructure/  # DB・外部サービスの具体実装
└── pkg/                 # 外部公開可能なユーティリティ
```

## 技術スタック

- Go 1.23 / chi / PostgreSQL (pgx)
- oapi-codegen（OpenAPI → Go コード自動生成）
- golang-migrate（マイグレーション）
- zerolog（ログ）/ viper（設定管理）

## エラーハンドリングパターン

HTTP ステータスコードの決定は **usecase 層** で行い、**handler 層** がそれを読み取ってフロントへ返す。

### 各層の責務

| 層 | やること | やらないこと |
|---|---|---|
| `domain/entity` | `errors.New()` でドメインエラーを定義する | HTTP コードを持たせる |
| `repository` | DB エラーを `apperrors` センチネルに変換する | ビジネスロジックを書く |
| `usecase` | `apperrors.NewXxx()` でエラーを生成・ラップする | raw な数値リテラル (`400` など) を書く |
| `handler/controller` | `apperrors.HTTPStatusCode(err)` と `apperrors.HTTPMessage(err)` でレスポンスを返す | エラーの種類を自前で判定する |

### コード例

**domain/entity** — 純粋な Go エラー（HTTP コード不要）

```go
var (
    ErrUserDisplayNameRequired = errors.New("display name is required")
    ErrUserGoogleIDRequired    = errors.New("google id is required")
)
```

**repository** — DB エラー → `apperrors` センチネルに変換

```go
func (r *UserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
    row, err := r.q.GetUserByID(ctx, id)
    if errors.Is(err, pgx.ErrNoRows) {
        return nil, apperrors.ErrNotFound   // センチネルをそのまま返す
    }
    if err != nil {
        return nil, apperrors.NewInternal("get user", err)
    }
    return toEntity(row), nil
}
```

**usecase** — 型付きコンストラクタ (`NewXxx`) を使う。raw な数値リテラルは書かない。

```go
// バリデーションエラー → 400
if err := user.Validate(); err != nil {
    return nil, apperrors.NewBadRequest("validation failed", err)
}

// 重複エラー → 409
if existing != nil {
    return nil, apperrors.NewConflict("email already in use", nil)
}

// repository のエラーは fmt.Errorf でコンテキストを追加して上流へ伝播
// （エラーチェーンが維持されるため HTTPStatusCode() は正しく動作する）
if err := r.repo.Create(ctx, user); err != nil {
    return nil, fmt.Errorf("create user: %w", err)
}
```

**handler/controller** — `apperrors` を通じてコードとメッセージを取得する

```go
out, err := h.createUC.Execute(r.Context(), input)
if err != nil {
    if apperrors.HTTPStatusCode(err) >= 500 {
        h.log.Error().Err(err).Msg("internal error")
    }
    h.writeError(w, apperrors.HTTPStatusCode(err), apperrors.HTTPMessage(err))
    return
}
```

### 利用可能なコンストラクタ（`pkg/apperrors`）

```go
apperrors.NewBadRequest(message string, err error)   // 400
apperrors.NewUnauthorized(message string, err error) // 401
apperrors.NewForbidden(message string, err error)    // 403
apperrors.NewNotFound(message string, err error)     // 404
apperrors.NewConflict(message string, err error)     // 409
apperrors.NewInternal(message string, err error)     // 500

// handler 層での読み取り
apperrors.HTTPStatusCode(err error) int    // AppError でなければ 500
apperrors.HTTPMessage(err error) string    // AppError でなければ "internal server error"
```

### ルール

1. `usecase` 層で raw な `http.StatusXxx` や数値リテラル（`400` など）を書かない
2. `repository` 層は必ず pgx/DB エラーを `apperrors` センチネルに変換して返す
3. `handler` 層は `err.Error()` を直接クライアントへ返さない（内部情報の漏洩防止）
4. `fmt.Errorf("...: %w", err)` でラップしてもエラーチェーンが維持されるため問題ない

---

## コード骨格の書き方

Claude がファイルを作成する際は以下のフォーマットで骨格を残す。

```go
// GetUser はユーザーIDからユーザー情報を取得する。
// エラー時は domain.ErrNotFound を返す。
func (c *UserController) GetUser(ctx context.Context, req gen.GetUserRequestObject) (gen.GetUserResponseObject, error) {
    // TODO: usecase を呼び出してユーザーを取得する
    // TODO: 404 / 200 のレスポンスを返す
    panic("not implemented")
}
```

## TDD フロー（バックエンド）

### 1. 関数設計（Claude）

実装前に以下を決める。

- 関数シグネチャ（引数・戻り値）
- 正常系・エラー系の一覧
- 依存する interface（usecase / repository）

### 2. テスト骨格を先に作成（Claude）

```go
// usecase/user/get_user_test.go
package user_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestGetUser_Success(t *testing.T) {
    // Arrange
    // TODO: mock repository をセットアップ

    // Act
    // TODO: usecase を実行

    // Assert
    // TODO: 返り値を検証
}

func TestGetUser_NotFound(t *testing.T) {
    // Arrange / Act / Assert
}
```

### 3. 実装（Copilot）

テストが Red の状態で実装ファイルを作成し、Green にする。

## ルール

- `gen/` 配下のファイルは編集しない（`go generate` で再生成する）
- **実装ファイルより先にテストファイルを作成する**
- 関数の中身の詳細実装は Copilot に委ねる
- 単体テストは各ファイルのそばに `xxx_test.go` として置く
- テストコードは方針とサンプル1件のみ書く（全量は Copilot に委ねる）
