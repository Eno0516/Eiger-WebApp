# テスト戦略

## 全体構成

```
E2E テスト（tests/e2e/）
    ↑ 実ブラウザ + 実サーバー + 実 DB
─────────────────────────────────────────
Backend                      Frontend
  Integration テスト            Page テスト
  API テスト                    Component テスト
  Unit テスト                   Unit テスト
```

---

## テストの書き方

### AAA（Arrange-Act-Assert）形式

全テストケースは **AAA 形式** で記述する。

```go
func TestXxx_Scenario(t *testing.T) {
    // Arrange — テストに必要なデータ・モック・依存を準備する
    repo := &mockUserRepo{...}
    uc   := NewXxxUseCase(repo, ...)

    // Act — テスト対象の処理を1回だけ実行する
    out, err := uc.Execute(ctx, input)

    // Assert — 結果を検証する
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if out.Field != want {
        t.Errorf("Field = %v, want %v", out.Field, want)
    }
}
```

- 各ブロックは `// Arrange` `// Act` `// Assert` のコメントで明示する
- Act は原則1行（複数の副作用を一度に検証しない）
- Assert は期待値を先に書く（`want`, `got` の順）

---

## Backend

### テスト種別一覧

| 種別 | 配置 | build tag | DB | 使用する Repository |
|------|------|-----------|-----|---------------------|
| **Unit テスト** | `internal/usecase/*/` | なし | ❌ 不要 | `repository/memory` |
| **API テスト** | `tests/api/` | `api` | ❌ 不要 | `repository/memory` |
| **Integration テスト** | `tests/integration/` | `integration` | ✅ 実 PostgreSQL | `repository/postgres` |

### Unit テスト

- UseCase のビジネスロジックを検証する
- `repository/memory` を DI するため DB 不要
- バリデーション・重複チェック・エラーハンドリング等を網羅する

### API テスト

- HTTP のインターフェース（スキーマ・ステータスコード・エラー形式）を検証する
- **実 DB には繋がない**。`repository/memory` を DI してサーバーを起動する
- フロントエンドとの契約が崩れていないことを確認するのが目的

```go
//go:build api

func TestCreateUser_API(t *testing.T) {
    repo := memory.NewUserRepository()
    uc := user.NewCreateUserUseCase(repo)
    handler := httphandler.NewUserHandler(uc, logger)
    router := httphandler.NewRouter(handler, logger)

    srv := httptest.NewServer(router)
    defer srv.Close()
    // ...
}
```

**DB をモックする理由：**
- DB の状態に依存しないのでテストが安定する
- 外部サービス（PostgreSQL）なしで CI が完結する
- 検証したいのは「HTTP の契約」であり DB の永続化ではない

### Integration テスト

- 複数層をまたいで実際の PostgreSQL と通信する
- testcontainers などで DB コンテナを起動して使う
- DB スキーマ・マイグレーション・トランザクションの整合性を確認するのが目的
- **実 DB を使う唯一の Backend テスト**

### build tag の使い方

```bash
# Unit テストのみ（普段の開発）
go test ./...

# API テスト
go test -tags=api ./tests/api/...

# Integration テスト
go test -tags=integration ./tests/integration/...
```

---

## Frontend

### テスト種別一覧

| 種別 | 対象 | DB・API | モック方法 |
|------|------|---------|-----------|
| **Unit テスト** | 関数・hooks・ユーティリティ | ❌ 不要 | — |
| **Component テスト** | 個別コンポーネント | ❌ 不要 | API レスポンスをモック |
| **Page テスト** | ページ全体の画面・導線 | ❌ 不要 | API レスポンスをモック（MSW 等） |

### Unit テスト

- 純粋関数・カスタム hooks・バリデーションロジック等を検証する
- DOM・API・外部依存は持たない

### Component テスト

- 単一コンポーネントのレンダリング・インタラクションを検証する
- API 呼び出しはモック（MSW や vi.fn()）で差し替える
- 実サーバーには繋がない

### Page テスト

- ページ単位の画面遷移・表示内容・ユーザー操作フローを検証する
- MSW（Mock Service Worker）等で API レスポンスをモックし、実サーバー不要
- 複数コンポーネントの組み合わせが正しく機能するかを確認するのが目的

---

## E2E テスト

| 配置 | ツール | DB | サーバー |
|------|--------|----|---------|
| `tests/e2e/` | Playwright 等 | ✅ 実 PostgreSQL | ✅ 実サーバー |

- 実ブラウザ + 実バックエンド + 実 DB で本番に近い環境で検証する
- ユーザーの主要導線（ログイン・登録・操作）を網羅する
- 実行コストが高いため、重要なフローのみに絞る
