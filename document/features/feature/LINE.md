# LINE連携機能 設計ドキュメント

## 概要

LINE Login（LIFF）をアプリの認証基盤として使用し、かつ公式LINE（Messaging API）でメッセージ送信を行う機能。

---

## 使用するLINEサービス

| サービス | 用途 |
|---|---|
| LINE Login | アプリへのログイン認証（OAuth 2.0 + OpenID Connect） |
| LIFF（LINE Front-end Framework） | LINE Login の実行・userId取得・LIFFフォーム |
| LINE Messaging API | BotからユーザーへのDM送信 |

すべて無料枠で利用可能。Messaging APIは月1000通まで無料。

---

## ログイン認証の仕組み（LIFF方式）

```
① ユーザーが「LINEでログイン」ボタンをクリック
        ↓
② LIFF SDK が LINE Login 画面を起動
        ↓
③ 認証成功後、LIFF SDK が IDトークンを取得
        ↓
④ フロントが IDトークン をバックエンドに送信（userId を直接送らない）
        ↓
⑤ バックエンドが LINE の公開鍵で署名を検証し、line_user_id を取得
        ↓
⑥ JWT を発行して返却
```

**セキュリティ原則**: userId をフロントからそのまま送るのではなく、LINE が署名した IDトークンをバックエンドで検証する。

## Messaging API 連携（メッセージ送信）

```
① メンバーが公式LINEを友だち追加（事前に必要）
        ↓
② ログイン時に取得した line_user_id を app DB に保存済み
        ↓
③ Admin が送信操作 or 自動トリガー
        ↓
④ バックエンドが app DB の line_user_id を使って push 送信
```

---

## 前提条件

- BotからのDMは、**ユーザーが先に公式LINEを友だち追加している場合のみ送信可能**
- 友だち追加が未完了のユーザーにはDMを送れない
- LINE連携は任意か必須かを要検討（フリマ機能・送迎機能で招待リンクが必要なため、実質必須になりうる）

---

## ユーザーフロー

### LINE連携フロー
```
1. メンバーがアプリ内「LINEと連携する」ボタンをクリック
2. LIFFページが起動
3. liff.init() → liff.getProfile() で userId を取得
4. userId を API に POST
5. バックエンドが users テーブルの line_user_id を更新
6. 連携完了表示
```

### 招待リンク送信フロー（Adminが実行）
```
1. Admin が対象メンバーに招待リンクを送信
2. バックエンドが対象ユーザーの line_user_id を取得
3. LINE Messaging API の pushMessage で個人チャットに送信
4. 送信完了
```

---

## API エンドポイント

| Method | Path | 必要ロール | 説明 |
|--------|------|-----------|------|
| POST | `/line/connect` | member以上 | LIFFから取得したuserIdを紐付け |
| DELETE | `/line/connect` | member以上 | LINE連携を解除 |
| POST | `/line/send/:userId` | admin | 指定ユーザーに招待リンクを送信 |
| POST | `/line/webhook` | 不要（LINE署名検証） | LINEからのWebhookを受信 |

---

## DBスキーマ変更

`users` テーブルに以下カラムを追加。

```sql
-- xxx_add_line_user_id.sql

-- +migrate Up
ALTER TABLE users
  ADD COLUMN line_user_id VARCHAR(255);

-- +migrate Down
ALTER TABLE users
  DROP COLUMN IF EXISTS line_user_id;
```

`line_user_id` は NULL 許容（未連携の場合はNULL）。

---

## フロントエンド実装

```typescript
// LIFF初期化・userId取得
await liff.init({ liffId: process.env.LIFF_ID })

if (!liff.isLoggedIn()) {
  liff.login()
  return
}

const profile = await liff.getProfile()
// profile.userId をバックエンドに送る
await api.post('/line/connect', { lineUserId: profile.userId })
```

LIFFのSDKは `@line/liff` パッケージで導入。

---

## 環境変数

```
LINE_CHANNEL_SECRET=       # Messaging API チャネルシークレット（Webhook署名検証に使用）
LINE_CHANNEL_ACCESS_TOKEN= # Messaging API チャネルアクセストークン（メッセージ送信に使用）
LIFF_ID=                   # LIFFアプリのID
```

---

## LINE Developersコンソール設定

| 設定 | 内容 |
|---|---|
| Messaging APIチャネル | Webhookを有効化し、バックエンドの `/line/webhook` を設定 |
| LINE Loginチャネル | LIFFアプリを追加し、エンドポイントURLをフロントのLIFFページに設定 |
| LIFF | Endpoint URLにアプリ内のLINE連携ページのURLを設定 |

---

## セキュリティ要件

| 要件 | 対応方針 |
|---|---|
| Webhook署名検証 | `X-Line-Signature` ヘッダを `LINE_CHANNEL_SECRET` で検証 |
| userId の詐称防止 | LIFFはLINEが署名するため、フロントから送られたuserIdは信頼可能 |

---

## 未決事項

- LINE連携を「任意」「必須」どちらにするか
- 連携していないユーザーへのフォールバック（招待リンクを送れない場合の代替手段）

---

## LIFFフォームによるイベント参加管理

### 背景・Googleフォームでは要件を満たせない理由

Googleフォームでは LINE userId を取得できないため、回答者と LINE アカウントを回答時点で紐付けることが不可能。
結果として「参加確定者だけに LINE で URL を送る」が技術的に成立しない。

→ 回答時点で LINE userId を取得できる **LIFFフォーム** に置き換える。

---

### LIFFフォームのメリット

| 項目 | 内容 |
|---|---|
| userId 自動取得 | 回答時点で LINE userId を取得し、回答内容と 1 行で Spreadsheet に保存 |
| 動的フォーム | アプリ登録済みユーザーには「名前・学部」などの入力項目を非表示にできる |
| 参加確定者への push | userId を持つ参加確定者だけに絞って push 送信が可能 |

---

### スプレッドシートへの保存構成

直接書き込みは APIキー漏洩・個人情報漏洩の危険があるため **絶対NG**。

```
LIFFフォーム
   ↓ HTTPS
自前API（バックエンド）  ← 保持せず中継するだけ
   ↓
Googleスプレッドシート
```

- 通信は HTTPS で暗号化されるため個人情報保護法上も問題なし
- 保存先を Google Sheets に一本化し、権限管理は Google に委譲

---

### 本人性モデル

#### アプリ登録済みユーザー

```
Googleアカウント → アプリUUID → LINE userId
```

- 初回 Google ログインで本人性を確立
- LIFF フォーム回答時は LINE userId で照合するため再ログイン不要
- 入力済みの「名前・学部」などの項目を非表示にできる

#### アプリ未登録ユーザー（イベント参加者のみ）

```
LINE userId のみ
```

- Google ログイン不要
- 必要最小限の情報だけ取得

---

### 参加確定者への LINE 送信フロー

```
1. LIFFフォームで回答（LINE userId + 回答内容を Spreadsheet に保存）
2. Spreadsheet で参加確定フラグを付ける（Admin が手動 or 自動化）
3. 参加確定かつ userId がある行だけを対象に push 送信
4. 公式 LINE の友だちが何千人いても userId リストにだけ送るため問題なし
```

---

### 個人情報保護法との整合性

| 観点 | 対応 |
|---|---|
| 保存先の一本化 | Google Sheets のみ。自前APIは中継のみで保持しない |
| 権限管理 | Google アカウントに委譲 |
| 取得情報の最小化 | アプリ未登録者からは必要最小限のみ取得 |
| 利用目的の明示 | フォーム内に明記することで適法 |

---

## 実装タスク（骨格）

### Backend
- [ ] `xxx_add_line_user_id.sql` マイグレーション作成
- [ ] `users` テーブルの `line_user_id` カラム対応（sqlc query追加）
- [ ] `LINE連携 usecase` 骨格（Connect / Disconnect）
- [ ] `LINE送信 usecase` 骨格（pushMessage）
- [ ] `LineHandler` 追加（`/line/connect`, `/line/send/:userId`）
- [ ] `WebhookHandler` 追加（`/line/webhook`・署名検証込み）
- [ ] LINE Messaging API クライアント（`internal/infrastructure/line/`）
- [ ] 環境変数追加（`.env.example`）

### Frontend（user アプリ）
- [ ] `@line/liff` パッケージ導入
- [ ] LINE連携ページ（LIFF起動・userId取得・API送信）
- [ ] 連携済み/未連携の表示切り替え
