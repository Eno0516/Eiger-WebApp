# 初期登録機能 設計ドキュメント

## 概要

LINE Login 後の初回ユーザー向けに、外部 LIFF フォームへ誘導する登録ページ。
個人情報の入力・管理はすべて LIFF フォーム（外部）で行い、アプリ内には保持しない。
フォーム送信後はタブを閉じ、元ページが `window.focus` イベントで完了を検知してリダイレクトする。

---

## 登録状態の管理

`users` テーブルに `is_registered` カラムを追加して初期登録の完了を管理する。
`circle_members` のロール（visitor / member / admin）はサークル内役割の管理に専念させる。

```sql
ALTER TABLE users ADD COLUMN is_registered BOOLEAN NOT NULL DEFAULT false;
```

| 状態 | users.is_registered | circle_members.role |
|------|---------------------|---------------------|
| LINE ログインのみ | false | レコードなし |
| 初期登録完了 | true | レコードなし（サークル未申請） |
| サークル申請中 | true | visitor |
| 承認済みメンバー | true | member / admin |

---

## 個人情報の保存方針

初期登録フォームで収集する個人情報は **Google Sheets にのみ保存** し、app DB には書き込まない。
app DB には `is_registered = true` と `display_name`（部内あだ名）のみ記録する。

### 収集する個人情報

| 項目 | 保存先 |
|------|--------|
| 大学・学部・学科 | Sheets のみ |
| 本名 | Sheets のみ |
| 学籍番号 | Sheets のみ |
| 部内あだ名（display_name） | app DB（users テーブル） |

Sheets の行は `line_user_id` をキーとして管理する。

### Sheets データのアクセス権限

| アクセス者 | 閲覧範囲 |
|-----------|---------|
| `system_admin` | 全ユーザーの初期登録データ |
| サークル `admin` | 自分のサークルの `circle_members` に存在するユーザーのみ |
| それ以外 | 閲覧不可 |

サークル admin が閲覧する際、バックエンドは `circle_members` を参照してフィルタリングを行う（pass-through 表示）。

---

## 活動フォームとの自動紐付け

ユーザーが活動フォームを回答する際、**毎回名前・学籍番号を入力させない**ことを原則とする。
LIFF SDK が自動で `line_user_id` を付与するため、バックエンドが集計時に初期登録 Sheets と JOIN して名前などを補完する。

```
[初期登録 Sheets]
line_user_id | 本名 | 学籍番号 | 大学 | 学部 | 学科

[活動フォーム Sheets]（LIFF 経由で line_user_id を自動付与）
line_user_id | 回答内容 | タイムスタンプ

↓ バックエンドが集計 API 呼び出し時に JOIN

[admin への集計結果]
本名 | 学籍番号 | 回答内容
```

### 制約
- 活動フォームの LIFF は `line_user_id` を自動送信する設計であること
- JOIN は集計 API のレスポンス生成時のみ行い、結果は app DB に保存しない（pass-through）

---

## ページ仕様

### `/sign_up`（初期登録ページ）

| 項目 | 内容 |
|------|------|
| 認証要否 | LINE 認証済みであること（未認証 → LINE ログイン画面へ） |
| 個人情報登録 | 不要（is_registered = false でもアクセス可能） |
| クエリパラメータ | `redirect_url`（登録完了後の遷移先。省略時は `/`） |
| 登録完了済みでのアクセス | `/` にリダイレクト |

#### 表示内容

- 登録手順の説明テキスト
- LIFF フォームを新規タブで開くボタン
- フォーム送信後の案内（タブを閉じるよう促す）

### `/`（ホーム）

| 項目 | 内容 |
|------|------|
| 認証要否 | LINE 認証済みであること（未認証 → LINE ログイン画面へ） |
| 個人情報登録 | 不要（is_registered = false でも閲覧可能） |

---

## ユーザーフロー

### 初期登録フロー

```
1. ユーザーが `/sign_up?redirect_url=<遷移先>` にアクセス（LINE 認証済み）
2. ページ上の「登録フォームを開く」ボタンをクリック
3. 外部 LIFF フォームが新規タブで開く
4. ユーザーが本名・学籍番号・大学・学部・学科・あだ名を記入して送信
5. バックエンドが受け取る:
   - app DB: users.is_registered = true, display_name を更新
   - Sheets: 個人情報を line_user_id をキーに書き込み（app DB には書かない）
6. LIFF フォームのタブを閉じる（LIFF 側で自動 or ユーザーが手動）
7. 元のタブ（/sign_up）が window の focus イベントを検知
8. GET /users/me でユーザー状態を再取得
   - is_registered = true → `redirect_url`（省略時は `/`）へ遷移
   - is_registered = false → そのまま待機（再度フォームへのリンクを表示）
```

### 活動フォーム集計フロー

```
1. ユーザーが LIFF 活動フォームを回答（line_user_id が自動付与）
2. 活動フォーム Sheets に line_user_id + 回答内容が記録される
3. circle admin が集計 API を呼び出す
4. バックエンドが:
   a. 活動フォーム Sheets から回答行を取得
   b. 各行の line_user_id で初期登録 Sheets を参照（JOIN）
   c. 本名・学籍番号等を補完してレスポンスを返す（app DB には保存しない）
```

### ホーム（`/`）アクセスフロー

```
1. ユーザーが `/` にアクセス
2. 認証状態を確認
   - 未認証 → LINE 認証ページへリダイレクト
   - 認証済み（is_registered = false でも可） → ホームコンテンツを表示
```

---

## LIFF フォーム連携

| 項目 | 内容 |
|------|------|
| フォーム URL | 環境変数 `VITE_LIFF_SIGNUP_FORM_URL` で設定 |
| 開き方 | `target="_blank"` で新規タブ |
| 完了検知 | 親ウィンドウの `window.addEventListener('focus', ...)` |
| 完了判定 | focus 後に `GET /users/me` を呼び `is_registered` を確認 |

---

## ルーティング一覧

| パス | 認証 | is_registered | 説明 |
|------|------|--------------|------|
| `/` | 必須（LINE 認証） | 不要 | ホーム。未認証は LINE ログインへリダイレクト |
| `/sign_up` | 必須（LINE 認証） | 不要 | 初期登録ページ。完了済みなら `/` へリダイレクト |

---

## API エンドポイント

| Method | Path | 必要ロール | 説明 |
|--------|------|-----------|------|
| GET | `/users/me` | ログイン済み | 自分のユーザー情報取得（is_registered 確認に使用） |
| POST | `/auth/signup` | ログイン済み | 初期登録フォーム受付（is_registered 更新・Sheets 書き込み） |
| GET | `/circles/:id/members/:userId/signup-form-response` | admin（対象サークル）または system_admin | 対象ユーザーの初期登録 Sheets 回答を取得（pass-through） |
| GET | `/circles/:id/activities/:activityId/responses` | admin（対象サークル） | 活動フォーム回答一覧（初期登録 Sheets と JOIN して名前等を補完） |

---

## 環境変数

```
VITE_LIFF_SIGNUP_FORM_URL=   # 初期登録 LIFF フォームの URL
```

---

## 未決事項

- LIFF フォームタブを閉じるタイミング（フォーム側で自動 or ユーザーが手動）
- `redirect_url` のホワイトリスト検証（オープンリダイレクト対策）

---

## 実装タスク（骨格）

### Backend

- [ ] `users` テーブルに `is_registered` カラム追加（マイグレーション）
- [ ] `POST /auth/signup` エンドポイント骨格（is_registered 更新・Sheets 書き込み）
- [ ] `GET /circles/:id/members/:userId/signup-form-response` 骨格（pass-through）
  - circle admin: 自サークルメンバーのみ閲覧可
  - system_admin: 全ユーザー閲覧可
- [ ] `GET /circles/:id/activities/:activityId/responses` 骨格
  - 活動フォーム Sheets + 初期登録 Sheets を line_user_id で JOIN
  - 結果は pass-through（app DB に保存しない）
- [ ] `User` エンティティに `IsRegistered` フィールド追加
- [ ] Sheets JOIN ロジック（`internal/infrastructure/sheets/`）

### Frontend（user アプリ）

- [ ] `SignUpPage`（`/sign_up`）骨格
  - LIFF フォーム URL をボタンとして表示（`target="_blank"`）
  - `window.addEventListener('focus', ...)` で完了を検知
  - `GET /users/me` で `is_registered` を確認して遷移
  - 登録完了済みでアクセスした場合は `/` へリダイレクト
- [ ] `features/signup/` に `LiffFormLink` コンポーネント骨格
- [ ] `RequireAuth` ガードを `/` にも適用
- [ ] `router.tsx` に `/sign_up` ルート追加（RequireAuth 配下）
- [ ] `.env.example` に `VITE_LIFF_SIGNUP_FORM_URL` 追加
