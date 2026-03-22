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

## ルール

- `gen/` 配下のファイルは編集しない（`go generate` で再生成する）
- 関数の中身の詳細実装は Copilot に委ねる
- 単体テストは各ファイルのそばに `xxx_test.go` として置く
- テストコードは方針とサンプル1件のみ書く（全量は Copilot に委ねる）
