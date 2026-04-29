# frontend

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

## アプリ内ディレクトリ構造

```
src/
├── app/
│   ├── layouts/     # レイアウトコンポーネント
│   ├── providers.tsx
│   └── router.tsx
├── features/        # 機能単位のモジュール
│   └── <feature-name>/
│       ├── index.tsx        # メインコンポーネント（エントリポイント）
│       ├── components/      # サブコンポーネント
│       └── types/           # 型定義
│       # ※ api/ hooks/ は必要になった時点で追加
└── pages/           # ルートに対応するページ（features を組み合わせる・ページ固有UIを持つ）
```

### features と pages の責務

| 層 | 責務 | 例 |
|----|------|----|
| `features/` | 機能に特化したUI・ロジック | ログインフォーム本体 |
| `pages/` | features の組み合わせ＋ページ固有UI | 戻るボタン、ページタイトル配置 |

## 技術スタック

- React 18 / TypeScript / Vite
- TanStack Query（サーバー状態管理）
- React Router v6
- Vitest + Testing Library（テスト）

## 開発

```bash
# 依存インストール
pnpm install

# 開発サーバー起動（user アプリ）
pnpm --filter @eiger/user dev

# 開発サーバー起動（admin アプリ）
pnpm --filter @eiger/admin dev
```
