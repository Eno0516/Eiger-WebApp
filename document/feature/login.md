# ログイン機能 設計ドキュメント

## 概要

メールアドレス＋パスワードによる認証機能。初回登録時にメールアドレス認証を必須とする。
パスワードリセット機能も含む。

---

## ロール定義

| ロール | 説明 |
|--------|------|
| `admin` | 管理者。全操作が可能。他ユーザーのロールを変更できる |
| `member` | admin に承認されたユーザー。内部情報（個人情報含む）を閲覧可能 |
| `pending` | メール認証済みだが、まだ admin に承認されていないユーザー。公開情報のみ閲覧可能 |

**ロール昇格フロー:**
```
登録 → メール認証完了 → pending（デフォルト）→ admin が承認 → member or admin
```

DBの `role` カラムで管理する。初回登録時は `pending` として保存。

---

## ユーザーフロー

### 新規登録フロー
```
1. ユーザーが名前・メール・パスワードを入力して登録
2. パスワードをハッシュ化してDBに保存（role=pending, email_verified_at=NULL）
3. メール認証トークンを生成してDBに保存（有効期限: 24時間）
4. 認証メールを送信
5. ユーザーがメール内のリンクをクリック
6. トークン検証 → email_verified_at に現在時刻を記録
7. 登録完了（この時点で role=pending のままログイン可能）
```

### ロール昇格フロー（admin操作）
```
1. admin がユーザー一覧で pending を確認
2. admin が対象ユーザーを member または admin に昇格
3. 対象ユーザーの次回 JWT 発行時から新ロールが反映される
```

### ログインフロー
```
1. ユーザーがメールアドレス＋パスワードを入力
2. DBからユーザーを取得
3. パスワードハッシュを検証
4. メール認証済みかチェック（未認証の場合はエラー）
5. JWT（アクセストークン）を発行して返却（payloadにroleを含む）
```

### パスワードリセットフロー
```
1. ユーザーがメールアドレスを入力
2. DBからユーザーを検索（存在しなくても同じレスポンスを返す ※ユーザー存在確認を防ぐ）
3. リセットトークンを生成してDBに保存（有効期限: 24時間）
4. リセットメールを送信
5. ユーザーがメール内のリンクをクリック → 新パスワード入力画面
6. トークン検証 → パスワードを更新、トークンを無効化
```

---

## API エンドポイント

| Method | Path | 必要ロール | 説明 |
|--------|------|-----------|------|
| POST | `/auth/register` | 不要 | 新規ユーザー登録 |
| POST | `/auth/login` | 不要 | ログイン・JWT発行 |
| GET | `/auth/verify-email` | 不要 | メールアドレス認証 |
| POST | `/auth/resend-verification` | 不要 | 認証メール再送信 |
| POST | `/auth/forgot-password` | 不要 | パスワードリセットメール送信 |
| POST | `/auth/reset-password` | 不要 | パスワードリセット実行 |
| PATCH | `/admin/users/:id/role` | admin | ユーザーのロール変更（pending→member, member→admin 等） |

---

## DBスキーマ変更

既存の `users` テーブルに以下カラムを追加。既存レコードなし前提のため、マイグレーションのみ対応。

```sql
-- 002_add_auth_columns.sql

-- +migrate Up
ALTER TABLE users
  ADD COLUMN password_hash                  VARCHAR(255) NOT NULL,
  ADD COLUMN role                           VARCHAR(50)  NOT NULL DEFAULT 'pending',
  ADD COLUMN email_verified_at              TIMESTAMPTZ,
  ADD COLUMN email_verification_token       VARCHAR(255),
  ADD COLUMN email_verification_expires_at  TIMESTAMPTZ,
  ADD COLUMN password_reset_token           VARCHAR(255),
  ADD COLUMN password_reset_expires_at      TIMESTAMPTZ;

-- +migrate Down
ALTER TABLE users
  DROP COLUMN IF EXISTS password_hash,
  DROP COLUMN IF EXISTS role,
  DROP COLUMN IF EXISTS email_verified_at,
  DROP COLUMN IF EXISTS email_verification_token,
  DROP COLUMN IF EXISTS email_verification_expires_at,
  DROP COLUMN IF EXISTS password_reset_token,
  DROP COLUMN IF EXISTS password_reset_expires_at;
```

マイグレーションファイル: `backend/internal/infrastructure/database/migrations/002_add_auth_columns.sql`

---

## ドメインモデル変更

### `User` エンティティへの追加フィールド

```go
type UserRole string

const (
    RoleAdmin   UserRole = "admin"
    RoleMember  UserRole = "member"
    RoleVisitor UserRole = "pending"
)

type User struct {
    ID                          string
    Name                        string
    Email                       string
    PasswordHash                string     // bcrypt ハッシュ
    Role                        UserRole
    EmailVerifiedAt             *time.Time // nil = 未認証
    EmailVerificationToken      *string
    EmailVerificationExpiresAt  *time.Time
    PasswordResetToken          *string
    PasswordResetExpiresAt      *time.Time
    CreatedAt                   time.Time
    UpdatedAt                   time.Time
}
```

### バリデーション追加
- パスワード: 8文字以上（詳細実装は Copilot に委ねる）

---

## 認証トークン方式

**JWT（アクセストークンのみ）** を採用。

| 項目 | 値 |
|------|----|
| 署名アルゴリズム | HS256 |
| 有効期限 | 1時間 |
| Payload | `user_id`, `email`, `role`, `exp` |

Payload に `role` を含めることで、ミドルウェアレベルでロール制御が可能。

---

## メール送信

| 項目 | 内容 |
|------|------|
| 使用ライブラリ | `gopkg.in/gomail.v2`（要追加）|
| メール認証トークン有効期限 | 24時間 |
| パスワードリセットトークン有効期限 | 24時間 |
| 再送信の制限 | 1分に1回 |

### 環境変数
```
SMTP_HOST=
SMTP_PORT=
SMTP_USER=
SMTP_PASSWORD=
EMAIL_FROM=
APP_BASE_URL=   # メール内リンクのベースURL
```

### 開発環境
`docker-compose.dev.yml` に **mailhog** を追加。

```yaml
mailhog:
  image: mailhog/mailhog
  ports:
    - "1025:1025"   # SMTP
    - "8025:8025"   # Web UI
```

開発時は `SMTP_HOST=mailhog`, `SMTP_PORT=1025` で設定。

---

## セキュリティ要件

| 要件 | 対応方針 |
|------|----------|
| パスワードハッシュ | `bcrypt`（cost=12）。`golang.org/x/crypto` は go.mod に存在 |
| メール・リセットトークン | `crypto/rand` で128bit生成 |
| JWT秘密鍵 | 環境変数 `JWT_SECRET_KEY` で管理 |
| ブルートフォース対策 | ミドルウェアで同一IPの連続失敗を検知（インメモリ実装、Redisなし） |
| ユーザー存在確認防止 | forgot-password は存在しないメールでも同一レスポンスを返す |

---

## ミドルウェア

```
認証不要:    POST /auth/*
pending以上: ログイン済みであれば到達可能な公開寄りエンドポイント
member以上:  個人情報・内部情報を含むエンドポイント
adminのみ:   PATCH /admin/users/:id/role など管理系エンドポイント
```

JWT の `role` フィールドを見てルーティングレベルで制御する。

---

## 未決事項（残り）

なし。全て決定済み。

---

## 実装タスク（骨格）

### Backend
- [x] `002_add_auth_columns.sql` マイグレーション作成
- [x] `sqlc.yaml` + `db/queries/user.sql` 作成（`sqlc generate` で `internal/infrastructure/database/sqlcgen/` に生成）
- [x] `User` エンティティ更新（role=pending デフォルト・認証フィールド追加）
- [x] `UserRepository` interface に `FindByVerificationToken`, `FindByResetToken`, `UpdateRole` 追加
- [x] postgres `UserRepository` を sqlcgen ベースの骨格に書き換え
- [ ] `sqlc generate` を実行して `sqlcgen/` を生成する
- [ ] `RegisterUsecase` 骨格
- [ ] `LoginUsecase` 骨格
- [ ] `VerifyEmailUsecase` 骨格
- [ ] `ForgotPasswordUsecase` 骨格
- [ ] `ResetPasswordUsecase` 骨格
- [ ] `UpdateUserRoleUsecase` 骨格（admin 専用）
- [ ] `AuthHandler` 追加
- [ ] `AdminHandler` 追加（`PATCH /admin/users/:id/role`）
- [ ] JWT 生成・検証パッケージ追加（`pkg/jwt/`）
- [ ] メール送信 interface + 骨格（`internal/infrastructure/email/`）
- [ ] JWT 認証ミドルウェア・ロール検証ミドルウェア追加（pending / member / admin の3段階）
- [ ] `gomail` を go.mod に追加

### Frontend（admin / user 両アプリ）
- [ ] ログインページ
- [ ] 登録ページ
- [ ] メール認証完了ページ
- [ ] パスワードリセット申請ページ
- [ ] パスワードリセット実行ページ
- [ ] 認証状態管理（Context or Zustand）
- [ ] axios インターセプターで JWT を Authorization ヘッダに付与
- [ ] 未認証時のリダイレクト処理
- [ ] ロールに応じたルートガード

### Infrastructure
- [ ] `docker-compose.dev.yml` に mailhog 追加
- [ ] `.env.example` に認証・メール関連の環境変数追加
