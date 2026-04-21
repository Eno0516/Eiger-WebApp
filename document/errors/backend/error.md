# バックエンド エラーハンドリング設計

## 基本方針

- HTTP ステータスコードの決定責務は **usecase 層** が持つ
- **handler 層** は usecase が返したエラーからステータスとメッセージを読み取るだけ
- クライアントへは `apperrors.HTTPMessage()` 経由の**整形済みメッセージ**のみを返し、内部エラー文字列を露出させない

---

## 各層の責務

| 層 | エラーに対してやること | やってはいけないこと |
|---|---|---|
| `domain/entity` | `errors.New()` でドメインエラーを定義する | HTTP ステータスコードを持たせる |
| `repository` | DB エラーを `apperrors` センチネルに変換して返す | ビジネスロジックを書く |
| `usecase` | `apperrors.NewXxx()` でエラーを生成・ラップしてステータスを付与する | raw な数値リテラル (`400` 等) や `http.StatusXxx` を直書きする |
| `handler / controller` | `apperrors.HTTPStatusCode(err)` と `apperrors.HTTPMessage(err)` でレスポンスを組み立てる | エラーの種類を自前で判定する / `err.Error()` をそのままクライアントへ返す |

---

## エラー伝播フロー

```
[DB / 外部API]
     │  pgx.ErrNoRows など生のエラー
     ▼
[repository]
     │  apperrors.ErrNotFound / apperrors.NewInternal(...) に変換
     ▼
[usecase]
     │  apperrors.NewBadRequest / NewConflict / NewNotFound ...
     │  またはそのまま fmt.Errorf("...: %w", err) でラップして伝播
     ▼
[handler / controller]
     │  apperrors.HTTPStatusCode(err) → HTTP ステータス
     │  apperrors.HTTPMessage(err)    → クライアント向けメッセージ
     ▼
[フロントエンド]
     { "error": "email already in use" }  HTTP 409
```

> `fmt.Errorf("create user: %w", appErr)` でラップしても `errors.As()` がチェーンを辿るため、`HTTPStatusCode()` / `HTTPMessage()` は正しく動作する。

---

## コード例

### domain/entity — 純粋な Go エラー

```go
var (
    ErrUserDisplayNameRequired  = errors.New("display name is required")
    ErrUserLineUserIDRequired   = errors.New("line user id is required")
)
```

### repository — DB エラーをセンチネルに変換

```go
func (r *UserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
    row, err := r.q.GetUserByID(ctx, id)
    if errors.Is(err, pgx.ErrNoRows) {
        return nil, apperrors.ErrNotFound             // センチネルをそのまま返す
    }
    if err != nil {
        return nil, apperrors.NewInternal("get user", err)
    }
    return toEntity(row), nil
}
```

### usecase — 型付きコンストラクタでステータスを付与

```go
// バリデーション失敗 → 400
if err := user.Validate(); err != nil {
    return nil, apperrors.NewBadRequest("validation failed", err)
}

// 重複 → 409
if existing != nil {
    return nil, apperrors.NewConflict("email already in use", nil)
}

// repository エラーはコンテキストを追加してそのまま伝播
if err := uc.repo.Create(ctx, user); err != nil {
    return nil, fmt.Errorf("create user: %w", err)
}
```

### handler — 読み取るだけ

```go
out, err := h.createUC.Execute(r.Context(), input)
if err != nil {
    if apperrors.HTTPStatusCode(err) >= 500 {
        h.log.Error().Err(err).Msg("internal error") // 5xx はサーバーログに残す
    }
    h.writeError(w, apperrors.HTTPStatusCode(err), apperrors.HTTPMessage(err))
    return
}
```

---

## apperrors パッケージ API（`pkg/apperrors`）

### センチネルエラー

repository 層が DB エラーをこれらに変換する。`errors.Is()` で比較可能。

| 変数 | ステータス |
|---|---|
| `apperrors.ErrNotFound` | 404 |
| `apperrors.ErrBadRequest` | 400 |
| `apperrors.ErrConflict` | 409 |
| `apperrors.ErrUnauthorized` | 401 |
| `apperrors.ErrForbidden` | 403 |
| `apperrors.ErrInternal` | 500 |

### 型付きコンストラクタ

usecase 層でエラーを生成するときに使う。

```go
apperrors.NewBadRequest(message string, err error)   // 400
apperrors.NewUnauthorized(message string, err error) // 401
apperrors.NewForbidden(message string, err error)    // 403
apperrors.NewNotFound(message string, err error)     // 404
apperrors.NewConflict(message string, err error)     // 409
apperrors.NewInternal(message string, err error)     // 500
```

### handler 向けユーティリティ

```go
apperrors.HTTPStatusCode(err error) int    // AppError でなければ 500 を返す
apperrors.HTTPMessage(err error) string    // AppError でなければ "internal server error" を返す
```

---

## ルール（チェックリスト）

- [ ] `usecase` に raw な数値リテラル (`400`, `409` ...) や `http.StatusXxx` を書いていない
- [ ] `repository` は pgx/DB エラーを必ず `apperrors` に変換してから返している
- [ ] `handler` は `err.Error()` をそのままクライアントへ返していない
- [ ] `handler` は 5xx エラーをサーバーログに記録している
