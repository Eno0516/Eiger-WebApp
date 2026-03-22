# CI/CD 作業ルール

## ワークフロー構成

| ファイル | トリガー | 内容 |
|----------|----------|------|
| `backend.yml` | `backend/**` の変更 | Go テスト・ビルド |
| `frontend.yml` | `frontend/**` の変更 | Lint・テスト・ビルド |
| `e2e.yml` | `tests/e2e/**` / `backend/**` / `frontend/**` の変更 | E2E テスト |

## ルール

- `paths` フィルターを必ず設定し、関係ないCIを走らせない
- `gen/` の再生成は CI に組み込む（yaml とコードの乖離を防ぐ）
- E2E は CI の最終ゲートとして位置づける
