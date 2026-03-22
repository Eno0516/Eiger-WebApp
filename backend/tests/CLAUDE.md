# テスト作業ルール

## テスト種別と配置

| 種別 | 場所 | build tag | 内容 |
|------|------|-----------|------|
| 単体テスト | `internal/xxx/xxx_test.go` | なし | モック使用・DB不要 |
| APIテスト | `tests/api/` | `api` | APIサーバー単体・スキーマ検証 |
| 統合テスト | `tests/integration/` | `integration` | DB + 複数層またぎ |

## build tag の使い方

```go
//go:build integration

package integration_test
```

```bash
# 単体テストのみ（普段の開発）
go test ./...

# APIテスト
go test -tags=api ./tests/api/...

# 統合テスト
go test -tags=integration ./tests/integration/...
```

## ルール

- テストコードの全量記述は Copilot に委ねる（Claude は方針とサンプル1件のみ）
- `testdata/` にフィクスチャデータを置く
- 統合テストは必ず build tag を付ける
