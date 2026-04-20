# テスト作業ルール

テストを書く際は必ず [tests.md](./tests.md) を参照すること。

## 役割分担

| 担当 | 作業内容 |
|------|----------|
| **Claude** | テストファイルの骨格作成・テスト構成の設計・モック方針の決定 |
| **GitHub Copilot** | テストケースの詳細実装・アサーションの補完・リファクタリング |

## テストの書き方

- **AAA 形式**（Arrange / Act / Assert）で統一する
- 各ブロックを `// Arrange` `// Act` `// Assert` のコメントで明示する
- 詳細・コード例は [tests.md](./tests.md) を確認

## 重要ポイント（抜粋）

- **APIテストは実 DB に繋がない** → `repository/memory` を DI して使う
- **Integrationテストのみ** 実 PostgreSQL を使う
- build tag を必ず付ける（`api` / `integration`）

詳細・コード例は [tests.md](./tests.md) を確認。
