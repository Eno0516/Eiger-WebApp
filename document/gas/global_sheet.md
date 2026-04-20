# Global Sheet — 初期登録データ書き込み

## 概要

初期登録時に本名・学籍番号・app UUID を Global Sheet に書き込む GAS web app。
閲覧権限は system_admin のみ（Google Sheets の共有設定で管理）。

## スプレッドシートスキーマ

| 列 | カラム名 | 型 | 説明 |
|---|---|---|---|
| A | app_uuid | string | `users.id`（Circle Sheet との JOIN キー） |
| B | real_name | string | 本名 |
| C | student_id | string | 学籍番号 |
| D | created_at | datetime | 書き込み日時（GAS が自動付与） |

## GAS web app エンドポイント

### POST（doPost）

フロントエンドから `UrlFetchApp` ではなく `fetch` で呼び出す。

**リクエスト（JSON body）**

```json
{
  "app_uuid": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "real_name": "山田太郎",
  "student_id": "2023001"
}
```

**レスポンス（成功）**

```json
{ "status": "ok" }
```

**レスポンス（失敗）**

```json
{ "status": "error", "message": "..." }
```

## 認証・セキュリティ

- GAS web app は「自分のアカウントで実行」「全員がアクセス可能」で公開する
- `app_uuid` は推測困難な UUID v4 のため、知っていないと書き込めない
- 本名・学籍番号はログに出力しない（GAS 側でも同様）
- スプレッドシート自体の閲覧権限を system_admin に限定する

## フロントエンド実装メモ

- `VITE_GAS_GLOBAL_SHEET_URL` 環境変数に GAS web app の URL を設定する
- `POST /auth/signup`（バックエンド）と並行して呼び出す
- 両方成功した場合のみ完了とみなす
- GAS が失敗した場合はバックエンド側もロールバックが必要（TODO: 検討）
