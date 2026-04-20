# ログイン機能 設計ドキュメント

## 概要

LINE Login（LIFF）による認証を採用。OAuth 2.0 + OpenID Connect ベースで、Google Login と同等のセキュリティ。
パスワードやメールアドレスはアプリ内に保持しない。
複数サークルへの対応を前提とした設計で、ロールはサークルごとに管理する。

---

## 認証方式

| 項目 | 内容 |
|------|------|
| プロトコル | OAuth 2.0 + OpenID Connect（LINE Login） |
| フロント側 | LIFF SDK（`@line/liff`）でIDトークンを取得 |
| バックエンド側 | LINEの公開鍵でIDトークンの署名を検証し、`line_user_id` を取得 |
| セキュリティ原則 | フロントから `line_user_id` を直接送らず、IDトークンをバックエンドで検証する |

---

## ロール定義

### システムロール（users テーブル）

| ロール | 説明 |
|--------|------|
| `is_system_admin = true` | サークルの新規作成が可能なスーパー管理者。DBを直接操作して付与する |
| `is_system_admin = false` | 一般ユーザー（デフォルト） |

### サークル内ロール（circle_members テーブル）

| ロール | 説明 |
|--------|------|
| `admin` | サークル管理者。全機能使用可能。メンバーのロール変更が可能 |
| `member` | Adminに承認されたメンバー。サークル内部データを閲覧可能 |
| `visitor` | 参加申請済みだが未承認の状態。公開情報のみ閲覧可能 |

**ロール昇格フロー:**
```
LINE ログイン → ログイン済み（circle_members レコードなし）
        ↓ 参加申請（LIFFフォーム）
visitor（申請中）
        ↓ Admin が承認
member or admin
```

circle_members にレコードがない状態 = visitor と同等の権限（全サークルの概要のみ閲覧可能）

---

## DBスキーマ

```sql
-- users テーブル
CREATE TABLE users (
  id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  line_user_id      VARCHAR(255) NOT NULL UNIQUE,  -- LINE の sub クレーム
  display_name      VARCHAR(255) NOT NULL,          -- 部内あだ名（個人を特定できる本名はNG）
  line_id           VARCHAR(255),                   -- LINE連携後に保存（メッセージ送信用）
  is_system_admin   BOOLEAN NOT NULL DEFAULT false,
  created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- circles テーブル
CREATE TABLE circles (
  id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name                 VARCHAR(255) NOT NULL,
  description          TEXT,
  spreadsheet_id       TEXT,  -- 個人情報・回答が蓄積される Spreadsheet の ID
  created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at           TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- circle_members テーブル（中間テーブル）
CREATE TABLE circle_members (
  circle_id   UUID NOT NULL REFERENCES circles(id) ON DELETE CASCADE,
  user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role        VARCHAR(50) NOT NULL DEFAULT 'visitor',  -- admin | member | visitor
  joined_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (circle_id, user_id)
);
```

---

## ユーザーフロー

### ログインフロー
```
1. ユーザーが「LINEでログイン」ボタンをクリック
2. LIFF SDK が LINE Login 画面を起動
3. 認証成功後、LIFF SDK が IDトークンを取得
4. フロントがIDトークンをバックエンドに送信
5. バックエンドがLINEの公開鍵でIDトークンを検証し line_user_id を取得
6. line_user_id でDBを検索
   - 初回: users レコードを新規作成（display_name は後から設定）
   - 既存: ユーザー情報を取得
7. JWTを発行して返却
```

### 参加申請フロー
```
1. ユーザーがサークル一覧からサークルを選択
2. 「参加申請」ボタンをクリック
3. LIFFフォームが起動
4. LIFF SDK が自動で line_user_id を取得
5. ユーザーが本名など必要事項を記入して送信
6. バックエンドが line_user_id + 回答内容を受け取り:
   - app DB: circle_members に visitor レコード作成 + line_user_id を保存
   - Sheets: サービスアカウント経由で個人情報・回答内容を書き込み（永続保存しない）
7. サークルの admin 全員に通知（アプリ内通知）
```

### Admin によるメンバー承認フロー
```
1. Admin が「メンバー管理」ページを開く
2. visitor 一覧が表示される（app DB から取得）
3. Admin が対象ユーザーを選択
4. バックエンドが Sheets API（サービスアカウント）で該当ユーザーの回答内容を取得
   → フロントに返すのみ。DBには保存しない（個人情報非保持の原則）
5. Admin が回答内容を確認後、「承認」ボタンをクリック
6. circle_members の role を member に更新
7. 承認されたユーザーにアプリ内通知を送信
```

---

## Google Sheets 連携

### 用途
- 参加申請時の個人情報・回答内容の書き込み
- 承認フローでの回答内容の読み取り（pass-through 表示、DBに保存しない）

### アクセス方式
サービスアカウントを使用。個人のGoogleアカウントに依存せず、JSONキーで認証。

```
Go backend → Sheets API（サービスアカウント）→ Spreadsheet
```

### 個人情報保護原則
- **保存しない**: Sheetsの個人情報はapp DBに書き込まない
- **pass-through**: 表示のためにSheetsから読んだデータはレスポンスで返すのみ
- **アクセス制御**: Sheets読み取りエンドポイントはadminロール限定
- **ログ非出力**: 個人情報をGoのログに出力しない
- **キャッシュ無効**: `Cache-Control: no-store` をレスポンスヘッダーに付与

---

## API エンドポイント

| Method | Path | 必要ロール | 説明 |
|--------|------|-----------|------|
| POST | `/auth/line` | 不要 | LINE IDトークンを受け取りJWTを発行 |
| POST | `/auth/logout` | ログイン済み | ログアウト |
| GET | `/circles` | ログイン済み | 全サークル一覧取得 |
| POST | `/circles` | system_admin | サークル新規作成 |
| POST | `/circles/:id/apply` | ログイン済み | 参加申請（LIFFフォーム送信・visitorレコード作成） |
| GET | `/circles/:id/members` | admin（対象サークル） | メンバー一覧取得（visitor含む） |
| GET | `/circles/:id/members/:userId/form-response` | admin（対象サークル） | 対象ユーザーのSpreadsheet回答を取得（pass-through） |
| PATCH | `/circles/:id/members/:userId/role` | admin（対象サークル） | ロール変更（visitor→member 等） |

---

## 認証トークン方式

**JWT（アクセストークンのみ）** を採用。

| 項目 | 値 |
|------|----|
| 署名アルゴリズム | HS256 |
| 有効期限 | 1時間 |
| Payload | `user_id`, `is_system_admin`, `exp` |

サークル内ロールはリクエストのたびにDBから取得する（JWTには含めない）。
サークルをまたぐロールをJWTに含めると肥大化するため。

---

## ミドルウェア

```
認証不要:         POST /auth/line
ログイン済み:      GET /circles, POST /circles/:id/apply
サークルadmin:    GET /circles/:id/members, PATCH /circles/:id/members/:userId/role
system_admin:    POST /circles
```

---

## 通知

アプリ内通知として管理する。

| イベント | 通知先 |
|---|---|
| 参加申請があった | 対象サークルの admin 全員 |
| 承認された（visitor→member） | 承認されたユーザー本人 |

---

## 環境変数

```
LINE_CHANNEL_ID=            # LINE Login チャネルID（IDトークン検証に使用）
LINE_CHANNEL_SECRET=        # LINE Login チャネルシークレット
LIFF_ID=                    # LIFFアプリのID
JWT_SECRET_KEY=             # JWT署名鍵
GOOGLE_SERVICE_ACCOUNT=     # Sheets API 用サービスアカウントのJSON（base64エンコード推奨）
```

---

## 未決事項

- display_name の初回設定タイミング（ログイン直後に入力画面を出すか）
- Sheetsのどの列を表示するか（サークルごとにLIFFフォームの構成が違うため）

---

## 実装タスク（骨格）

### Backend
- [ ] `001_create_users.sql` マイグレーション（line_user_id, is_system_admin）
- [ ] `002_create_circles.sql` マイグレーション
- [ ] `003_create_circle_members.sql` マイグレーション
- [ ] sqlc query 追加（users / circles / circle_members）
- [ ] `User`, `Circle`, `CircleMember` エンティティ定義
- [ ] `LineAuthUsecase` 骨格（IDトークン検証 → JWT発行）
- [ ] `ApplyCircleUsecase` 骨格（visitor レコード作成 + Sheets書き込み + 通知）
- [ ] `ApproveCircleMemberUsecase` 骨格（role更新 + 通知）
- [ ] `GetFormResponseUsecase` 骨格（Sheets API で回答取得・pass-through）
- [ ] `AuthHandler` 追加（`/auth/line`）
- [ ] `CircleHandler` 追加
- [ ] JWT 生成・検証パッケージ（`pkg/jwt/`）
- [ ] LINE IDトークン検証パッケージ（`pkg/lineauth/`）
- [ ] Google Sheets API クライアント（`internal/infrastructure/sheets/`、サービスアカウント認証）
- [ ] 認証ミドルウェア・サークルロール検証ミドルウェア
- [ ] 環境変数追加（`.env.example`）

### Frontend（admin アプリ）
- [ ] メンバー管理ページ（visitor一覧・Sheets回答表示・承認ボタン）
- [ ] アプリ内通知コンポーネント

### Frontend（user アプリ）
- [ ] LINEログインページ（LIFF SDK 初期化・IDトークン取得・バックエンド送信）
- [ ] display_name 初回設定画面
- [ ] サークル一覧ページ
- [ ] 参加申請ボタン（LIFFフォーム起動）
- [ ] 認証状態管理（Context or Zustand）
- [ ] JWTをAuthorizationヘッダに付与するインターセプター
- [ ] 未認証時のリダイレクト処理
- [ ] サークル内ロールに応じたルートガード
- [ ] アプリ内通知コンポーネント

### Infrastructure
- [ ] `.env.example` に認証関連の環境変数追加
- [ ] Google サービスアカウント作成・Spreadsheet への共有設定
