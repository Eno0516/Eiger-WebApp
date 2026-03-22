# Claude 作業ルール

## 役割分担の方針

このプロジェクトでは **トークン消費を最適化** するために、以下の役割分担で運用する。

| 担当 | 作業内容 |
|------|----------|
| **Claude（あなた）** | アーキテクチャ設計・ファイル/関数の骨格作成・実装方針の決定 |
| **GitHub Copilot** | 関数の中身の詳細実装・補完・リファクタリング |

---

## Claude が行うこと（粒度：大）

- ディレクトリ構造・ファイルの新規作成
- interface・struct・関数シグネチャの定義
- 処理の流れ（TODO コメント）を骨格として記述
- 複雑なロジックの設計方針をコメントで明記
- 依存関係の整理・パッケージ構成の決定

## Claude が行わないこと

- 関数の中身の詳細な実装（Copilot に委ねる）
- 自明な getter/setter・単純な CRUD の実装
- テストコードの全量記述（方針とサンプル1件のみ書く）

---

## コード骨格の書き方

Claude がファイルを作成する際は、以下のフォーマットで骨格を残す。

```go
// GetUser はユーザーIDからユーザー情報を取得する。
// エラー時は domain.ErrNotFound を返す。
func (c *UserController) GetUser(ctx context.Context, req gen.GetUserRequestObject) (gen.GetUserResponseObject, error) {
    // TODO: usecase を呼び出してユーザーを取得する
    // TODO: 404 / 200 のレスポンスを返す
    panic("not implemented")
}
```

---

## プロジェクト概要

- **言語**: Go
- **アーキテクチャ**: モノレポ / Clean Architecture + DDD
- **API定義**: OpenAPI (oapi-codegen で自動生成)
- **構成**:
  - `api/` … openapi.yaml
  - `gen/` … 自動生成コード（手を加えない）
  - `internal/controller/` … gen の interface 実装
  - `internal/usecase/` … ビジネスロジック
  - `internal/domain/` … Entity / ValueObject / Repository interface
  - `internal/infrastructure/` … DB等の具体実装
  - `cmd/server/` … エントリーポイント

---

## 一般ルール

- 返答は **日本語** で行う
- ファイルの削除・リネームは必ず確認を取る
- `gen/` 配下のファイルは編集しない
- 1回の指示で変更するファイルは最小限にとどめる
