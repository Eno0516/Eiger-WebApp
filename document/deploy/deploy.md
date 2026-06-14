# デプロイ構成

## 構成概要

| 役割 | サービス | 料金 |
|---|---|---|
| フロントエンド（user / admin） | Vercel | 無料 |
| バックエンド（Go サーバー） | Railway | 実質無料（後述） |
| データベース（PostgreSQL） | Neon | 無料枠あり |

---

## Vercel（フロントエンド）

- `frontend/apps/user/` と `frontend/apps/admin/` をそれぞれ静的サイトとしてデプロイ
- GitHub リポジトリと連携し、`main` ブランチへのマージで自動デプロイ
- 無料枠で十分運用可能

### 設定が必要な環境変数（user アプリ）

```
VITE_LIFF_ID=              # LINE LIFF アプリの ID
VITE_LIFF_SIGNUP_FORM_URL= # 初期登録 LIFF フォームの URL
```

---

## Railway（バックエンド）

- Go の HTTP サーバーをコンテナとして常駐稼働させる
- Render の無料枠は15分無操作でスリープするため NG（LINE ログイン時にタイムアウトが発生する）
- Railway は Hobby プラン（$5/月）だが **$5 分のクレジットが毎月付与される**
- このアプリの規模（サークル管理・イベント管理）では実際の使用量は $0.10/月 程度のため **実質無料**

### 上限設定（勝手に課金されないために必須）

```
Railway ダッシュボード
→ Account Settings
→ Usage
→ Spending Limit を $5 に設定
```

上限に達するとサービスが停止し、それ以上は課金されない。
合わせて Usage の通知メール（Notification）もオンにしておくこと。

### 設定が必要な環境変数

```
DB_USER=
DB_PASSWORD=
DB_NAME=
JWT_SECRET_KEY=
LINE_CHANNEL_ID=
LINE_CHANNEL_SECRET=
LIFF_ID=
```

---

## Neon（PostgreSQL）

- サーバーレス PostgreSQL
- 無料枠: ストレージ 0.5GB、コンピュート 191時間/月
- このアプリの規模では無料枠で十分
- Railway からの接続には Neon が発行する接続文字列（`DATABASE_URL`）を使用する

---

## LIFF 登録の前提

LIFF の Endpoint URL にはデプロイ済みの **HTTPS URL** が必要。
そのため、**Vercel へのデプロイが完了してから** LINE Developers コンソールで LIFF 登録を行う。

### 登録手順

```
1. Vercel デプロイ完了 → user アプリの URL を確認
2. LINE Developers コンソール
   → LINE Login チャネルを作成
   → LIFF アプリを追加
   → Endpoint URL に Vercel の URL を設定
3. 発行された LIFF ID を Vercel の環境変数（VITE_LIFF_ID）に設定
4. Railway の環境変数（LINE_CHANNEL_ID, LINE_CHANNEL_SECRET）を設定
```

---

## デプロイ順序

```
1. Neon でデータベース作成 → 接続文字列を取得
2. Railway にバックエンドをデプロイ → API の URL を確認
3. Vercel にフロントエンドをデプロイ → user アプリの URL を確認
4. LINE Developers で LIFF 登録 → LIFF ID を取得
5. 各サービスの環境変数を設定して動作確認
```
