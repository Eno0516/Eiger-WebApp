# フロントエンド作業ルール

## 構成

```
frontend/
├── apps/
│   ├── user/    # ユーザー向けアプリ (@eiger/user)
│   └── admin/   # 管理画面アプリ (@eiger/admin)
├── packages/
│   └── shared/  # 共通コンポーネント・ロジック
└── gen/         # OpenAPI → TypeScript 型定義（手を加えない）
```

## 技術スタック

- React 18 / TypeScript / Vite
- TanStack Query（サーバー状態管理）
- React Router v6
- Vitest + Testing Library（テスト）

## ルール

- `gen/` 配下のファイルは編集しない
- 単体テストはコンポーネントのそばに `xxx.test.tsx` として置く
- 共通コンポーネントは `packages/shared/` に置く
- 関数の中身の詳細実装は Copilot に委ねる
