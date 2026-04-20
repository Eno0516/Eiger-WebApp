# フロントエンド エラーハンドリング設計

## 基本方針

- API エラーの**表示方法はバックエンドの HTTP ステータスコードから自動決定**する
- 開発者はエラーを意識して書く必要がない。TanStack Query 経由のエラーは何もしなくても自動で通知される
- どうしても表示方法を変えたい場合だけ `AppError` の `display` オプションで override する

---

## エラーの流れ（自動）

```
バックエンド
  └── HTTP レスポンス { "error": "email already in use" }  status: 409
          ↓
httpClient interceptor（packages/shared/src/api/httpClient.ts）
  └── axios エラー → AppError に変換
      new AppError('email already in use', { status: 409 })
          ↓
TanStack Query（QueryCache / MutationCache の onError）
  └── notifyError(err) を自動呼び出し
          ↓
ErrorNotificationProvider
  └── status → display を自動決定 → Snackbar または Dialog を表示
```

**TanStack Query を使った API 呼び出しでは、エラーハンドリングのコードを一切書かなくてよい。**

```tsx
// これだけで、エラー時は自動でスナックバー/ダイアログが出る
const { mutate } = useMutation({ mutationFn: circleApi.apply })
```

---

## status → display の自動マッピング

バックエンドが返す HTTP ステータスコードに基づいて表示方法が自動決定される。

| status | display | 表示コンポーネント | 理由 |
|---|---|---|---|
| 401 | `none` | なし | interceptor がログインページへリダイレクト済み |
| 403 | `dialog` | Dialog | 権限エラーは目立たせる |
| 5xx | `dialog` | Dialog | サーバー障害はユーザーに明示する |
| 0（ネットワーク断） | `dialog` | Dialog | 通信できていないことを明示する |
| その他（404, 409 等） | `snackbar` | Snackbar | 軽微なエラーは控えめに通知 |

---

## イベントハンドラ・async 処理でのエラー通知

TanStack Query を使わない処理（フォームのバリデーション等）では、
`useErrorNotification()` フックの `notifyError` を手動で呼ぶ。

```tsx
import { useErrorNotification } from '@eiger/shared/errors'

function SomeComponent() {
  const { notifyError } = useErrorNotification()

  const handleClick = async () => {
    try {
      await someAsyncOperation()
    } catch (err) {
      notifyError(err)  // AppError でも普通の Error でも何でも渡してよい
    }
  }
}
```

`notifyError` に渡した値は内部で `AppError.from(err)` に通される。

| 渡す値 | 変換結果 |
|---|---|
| `AppError` | そのまま使用（status / display を引き継ぐ） |
| `Error` | message を引き継ぎ、status=500 → dialog として扱う |
| それ以外 | `'予期しないエラーが発生しました'` status=500 → dialog |

---

## 表示方法を override する

デフォルトのマッピングと異なる表示をしたい場合のみ、`display` を明示する。

```tsx
import { AppError } from '@eiger/shared/errors'

// 409 は本来 snackbar だが、この操作ではダイアログで伝えたい場合
throw new AppError('このメールアドレスは既に使用されています', {
  status: 409,
  display: 'dialog',  // ← override
})
```

```tsx
// エラーを表示したくない場合（自前でハンドリングするなど）
throw new AppError('...', { status: 422, display: 'none' })
```

---

## AppError クラス API（`@eiger/shared/errors`）

```ts
new AppError(message: string, options: AppErrorOptions)

interface AppErrorOptions {
  status: number       // HTTP ステータスコード（ネットワークエラーは 0）
  display?: ErrorDisplay  // 省略時は status から自動決定
}

type ErrorDisplay = 'snackbar' | 'dialog' | 'none'

// ユーティリティ
AppError.from(err: unknown): AppError  // 任意のエラーを AppError に変換
AppError.is(err: unknown): err is AppError  // 型ガード
```

---

## ルール（チェックリスト）

- [ ] TanStack Query を使っている API 呼び出しでは `catch` を書いていない
- [ ] イベントハンドラ・async 処理の `catch` では `notifyError(err)` を呼んでいる
- [ ] `display` の override は本当に必要な場合だけ行っている
- [ ] `err.message` を直接 UI に表示していない（`notifyError` に委ねる）
