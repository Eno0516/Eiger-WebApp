# 活動参加申し込み機能 設計ドキュメント

## 概要

管理者が作成した活動ページ（`event/{uuid}`）を通じて、ユーザーが活動への参加を申し込む機能。
LINEでURLを送付し、受け取ったユーザーが任意のブラウザでページを開いてそのまま申し込みまで完結できる設計。

---

## 全体フロー

```
[管理者] 活動参加確認用LIFFフォームを作成
      ↓
[管理者] 管理画面で event/{uuid} と LIFFフォームのURLを紐付ける
      ↓
[管理者] 活動ページのURLを参加希望者にLINEで送信
      ↓
[ユーザー] URLを開く（LINEのアプリ内ブラウザ・外部ブラウザどちらでも可）
      ↓
LINE認証セッションがない → LINE Login OAuth ページへ
      ↓
認証完了 → event/{uuid} に戻る
      ↓
circle_members のロールを確認
  ├─ visitor（未登録 or 参加申請未承認）
  │    └─ 「個人情報を登録してください」メッセージ表示
  │       「参加希望」ボタンは非表示
  │       「登録する」ボタン → 個人情報登録LIFFを別タブで開く
  │           ↓ 登録完了
  │       LIFFページ:「タブを閉じてください」表示
  │           ↓ ユーザーがタブを閉じて event/{uuid} のタブに戻る
  │           ↓ ウィンドウフォーカスイベントでロールを再取得
  │       member になっていれば「参加希望」ボタンが表示される
  │
  └─ member / admin
       └─ 活動詳細表示 + 「参加希望」ボタン
             ↓ クリック
       申し込みLIFFフォームを開く（同タブ or 別タブ）
             ↓ フォーム送信完了
       LIFFページ: 完了メッセージ（または自動で event/{uuid} に戻る）
```

---

## 設計判断

### LINE認証方式

LIFF（LINE内ブラウザ限定）ではなく、**LINE Login Web版（OAuth 2.0）** をこのアプリに組み込む。

- LINEのアプリ内ブラウザでも外部ブラウザでも動作する
- LIFF は申し込みフォーム・個人情報登録フォームのみに使用する

### LIFF→アプリへの戻り方

LIFFはウェブアプリへの任意URLリダイレクトが難しいため、以下の方式を採用：

- 個人情報登録LIFFは **別タブで開く**
- 登録完了後は「タブを閉じてください」メッセージを表示
- ユーザーが元のタブに戻ったとき、`window focus` イベントでロールを再取得して UI を更新

### ロールの取得タイミング

JWT にはサークルロールを含めない（`login.md` の方針通り）。
event/{uuid} ページ表示時および `window focus` 時に `GET /circles/:id/members/me` を叩いてロールを確認する。

---

## UI仕様

### event/{uuid} ページ

| ユーザー状態 | 表示内容 |
|---|---|
| 未認証 | LINE Loginボタンのみ表示 |
| visitor（登録済みだが未承認） | 「申し込みには承認が必要です」メッセージ |
| 未登録（circle_membersレコードなし） | 「個人情報を登録してください」＋「登録する」ボタン |
| member / admin | 活動詳細＋「参加希望」ボタン |

### 個人情報登録ページ（別タブ・LIFF）

- 登録完了後: 「登録が完了しました。このタブを閉じて元のページに戻ってください」

---

## フロントエンド実装方針

### window focus によるロール再取得

```typescript
// packages/shared/src/hooks/useWindowFocus.ts
import { useEffect } from 'react'

export function useWindowFocus(callback: () => void) {
  useEffect(() => {
    window.addEventListener('focus', callback)
    return () => window.removeEventListener('focus', callback)
  }, [callback])
}
```

event/{uuid} ページでの使用例：

```tsx
const { data: memberRole, refetch } = useQuery({
  queryKey: ['circle-member-role', circleId],
  queryFn: () => fetchMyCircleRole(circleId),
})

useWindowFocus(refetch)
```

---

## API エンドポイント

| Method | Path | 必要ロール | 説明 |
|--------|------|-----------|------|
| GET | `/events/:uuid` | ログイン済み | 活動ページ情報取得 |
| GET | `/circles/:id/members/me` | ログイン済み | 自分のサークルロール取得 |
| POST | `/circles/:id/apply` | ログイン済み | 参加申請（visitorレコード作成） |

---

## 未決事項

- 申し込みLIFFフォームを別タブで開くか同タブで開くか（UXの確認）
- 個人情報登録LIFFのLIFF IDをどこで管理するか（DBのcirclesテーブルに持つか、環境変数か）
- 活動ページ（event/{uuid}）の管理画面での作成フローの詳細設計

---

## 実装タスク（骨格）

### Backend

- [ ] `events` テーブルのマイグレーション（circle_id, uuid, liff_form_url, 活動詳細等）
- [ ] `GET /events/:uuid` ハンドラ
- [ ] `GET /circles/:id/members/me` ハンドラ（自分のロール返却）
- [ ] 管理画面用: `POST /events`・`PATCH /events/:uuid` ハンドラ

### Frontend（user アプリ）

- [ ] `event/{uuid}` ページ骨格
- [ ] ロール別UI出し分けコンポーネント
- [ ] `useWindowFocus` フック（`packages/shared/`）
- [ ] 個人情報登録LIFFを別タブで開くボタン

### Frontend（admin アプリ）

- [ ] 活動作成・編集ページ（LIFFフォームURL入力欄含む）
- [ ] event/{uuid} とLIFFフォームURLの紐付けUI
