# DB スキーマ設計

## 保存場所の方針

| データ種別 | 保存先 | 理由 |
|-----------|--------|------|
| 認証・権限・広告用属性・表示名 | app DB | 認証・権限制御・広告ターゲティングに必要 |
| 本名（初期登録時） | Google Sheets のみ | 個人情報非保持の原則 |
| 緊急連絡先などサークル固有情報 | Google Sheets のみ | 個人情報非保持の原則 |

Sheets の行は `users.id`（app UUID）をキーとして管理する。

---

## テーブル一覧

```
users                # 認証・システム管理（最小限）
user_profiles        # アプリ内プロフィール
user_line_accounts   # 公式LINE連携（任意・後から追加）
circles              # サークル
circle_members       # ユーザー ↔ サークル（ロール管理）
```

---

## テーブル定義

### users

認証・システム管理のみを担う。アプリ内表示情報は `user_profiles` へ分離。

```sql
CREATE TABLE users (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  line_user_id    VARCHAR(255) NOT NULL UNIQUE,  -- LINE の sub クレーム（認証用）
  is_system_admin BOOLEAN NOT NULL DEFAULT false,
  is_registered   BOOLEAN NOT NULL DEFAULT false, -- 初期登録（LIFFフォーム）完了フラグ
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

### user_profiles

アプリ内プロフィール。`users` と 1:1 対応。初期登録完了時に作成。

```sql
CREATE TABLE user_profiles (
  user_id      UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  display_name VARCHAR(255) NOT NULL,  -- 部内あだ名（個人特定不可のもの）
  university   VARCHAR(255),           -- 広告ターゲティング用
  faculty      VARCHAR(255),           -- 広告ターゲティング用
  grade        SMALLINT,               -- 学年（1〜4）。留年時は本人が更新
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

**`grade` の運用ルール:**
- ユーザーが自己申告で入力し、プロフィール編集で更新可能
- 留年・休学などの場合は本人が手動更新する
- 広告用途のため、多少の不正確さは許容

### user_line_accounts

公式 LINE アカウント連携後に作成。未連携ユーザーはレコードなし。

```sql
CREATE TABLE user_line_accounts (
  user_id   UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  line_id   VARCHAR(255) NOT NULL,  -- 公式LINEでのメッセージ送信に使用
  linked_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

**`line_user_id` との違い:**

| カラム | テーブル | 用途 |
|--------|---------|------|
| `line_user_id` | users | LINE Login の認証識別子（UID） |
| `line_id` | user_line_accounts | 公式LINEアカウントでのメッセージ送信先 |

### circles

サークル情報。

```sql
CREATE TABLE circles (
  id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name                 VARCHAR(255) NOT NULL,
  description          TEXT NOT NULL DEFAULT '',
  application_form_url TEXT NOT NULL DEFAULT '',  -- 参加申請 LIFF フォーム URL
  spreadsheet_id       VARCHAR(255) NOT NULL DEFAULT '',  -- サークル固有 Sheets ID
  uuid_column          VARCHAR(10) NOT NULL DEFAULT 'A',  -- Sheets 内の UUID カラム位置
  created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at           TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

### circle_members

ユーザーとサークルの中間テーブル。ロールを管理。

```sql
CREATE TABLE circle_members (
  circle_id UUID NOT NULL REFERENCES circles(id) ON DELETE CASCADE,
  user_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role      VARCHAR(50) NOT NULL DEFAULT 'visitor',  -- admin | member | visitor
  joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (circle_id, user_id)
);
```

**ロール定義:**

| role | 説明 |
|------|------|
| `admin` | サークル管理者。全機能使用可能・メンバーのロール変更が可能 |
| `member` | admin に承認されたメンバー。サークル内部データを閲覧可能 |
| `visitor` | 参加申請済みだが未承認。公開情報のみ閲覧可能 |

`circle_members` にレコードがない状態 = 全サークルの概要のみ閲覧可能（visitor 相当）

---

## ユーザー状態遷移

```
LINE ログイン
  → users レコード作成（is_registered = false）
        ↓ 初期登録（LIFF フォーム送信）
  → is_registered = true、user_profiles レコード作成
        ↓ サークル参加申請
  → circle_members: visitor
        ↓ admin が承認
  → circle_members: member または admin
        ↓ 公式 LINE 連携
  → user_line_accounts レコード作成
```

---

## Google Sheets との役割分担

| データ | 保存先 | キー |
|--------|--------|------|
| 本名 | 初期登録 Sheets | `users.id`（app UUID） |
| 緊急連絡先・サークル固有情報 | サークル固有 Sheets | `users.id`（app UUID） |
| 活動フォーム回答 | 活動フォーム Sheets | `users.id`（app UUID） |

活動フォーム集計時、バックエンドが活動フォーム Sheets と初期登録 Sheets を `users.id` で JOIN して本名等を補完する（pass-through、app DB には保存しない）。

---

## migration ファイル一覧

| ファイル | 内容 |
|---------|------|
| `001_create_users.sql` | users テーブル作成 |
| `002_create_circles.sql` | circles テーブル作成 |
| `003_create_circle_members.sql` | circle_members テーブル作成 |
| `004_create_user_profiles.sql` | user_profiles テーブル作成 |
| `005_create_user_line_accounts.sql` | user_line_accounts テーブル作成 |
