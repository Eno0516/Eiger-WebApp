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

## features と pages の責務分担

### features（機能モジュール）

各機能は `src/features/<feature-name>/` 以下に配置する。

```
features/<feature-name>/
├── index.tsx        # メインコンポーネント（エントリポイント）
├── components/      # サブコンポーネント（機能内でのみ使用）
└── types/           # 型定義
```

- その機能に**特化したUI・ロジックのみ**を持つ
  - 例: `features/login/` はログインフォームだけを担当する
- ページ遷移・ボタンのような**ページ固有の関心事は含めない**
- `api/`・`hooks/` は必要になった時点で追加する
- `index.tsx` からのみ外部に公開する（内部実装は外から直接 import しない）

### pages（ページコンポーネント）

- ルートに 1:1 で対応するコンポーネントを置く
- features のコンポーネントを **組み合わせる・配置する**責務を持つ
- ページ固有のUI（戻るボタン・ヘッダー等）はここに書く

```tsx
// pages/LoginPage.tsx の例
import { LoginForm } from '../features/login'

export function LoginPage() {
  return (
    <>
      <button onClick={() => navigate(-1)}>戻る</button>
      <LoginForm />
    </>
  )
}
```

## TDD フロー（フロントエンド）

### 1. コンポーネント・関数設計（Claude）

実装前に以下を決める。

- Props の型定義
- 表示すべき状態（ローディング・エラー・成功）
- ユーザー操作とその副作用

### 2. テスト骨格を先に作成（Claude）

```tsx
// features/login/components/LoginForm.test.tsx
import { render, screen } from '@testing-library/react'
import { LoginForm } from './LoginForm'

describe('LoginForm', () => {
  it('フォームが正しく表示される', () => {
    // Arrange
    // TODO: render

    // Act / Assert
    // TODO: 要素の存在を検証
  })

  it('送信時に onSubmit が呼ばれる', () => {
    // Arrange / Act / Assert
  })
})
```

### 3. 実装（Copilot）

テストが Red の状態でコンポーネントを実装し、Green にする。

## ルール

- `gen/` 配下のファイルは編集しない
- **実装ファイルより先にテストファイルを作成する**
- 単体テストはコンポーネントのそばに `xxx.test.tsx` として置く
- 共通コンポーネントは `packages/shared/` に置く
- 関数の中身の詳細実装は Copilot に委ねる
