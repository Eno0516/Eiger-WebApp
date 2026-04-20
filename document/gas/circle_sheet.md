# Circle Sheet — サークル参加時の個人情報コピー

## 概要

サークルの活動フォーム回答を Circle Sheet にエクスポートする GAS web app。
フォームの回答はまずアプリ DB に保存され、admin がエクスポートを実行したタイミングで GAS が Circle Sheet に書き込む。
このとき GAS が Global Sheet から本名を検索して付与することで、DB に個人情報を持たずに集計できる。

Circle Sheet の閲覧権限はそのサークルの admin が管理（Google Sheets の共有設定）。

## Circle Sheet スキーマ

| 列 | カラム名 | 型 | 説明 |
|---|---|---|---|
| A | app_uuid | string | `users.id`（Global Sheet との JOIN キー） |
| B | real_name | string | Global Sheet から取得した本名 |
| C | display_name | string | バックエンド DB から取得（フロントが渡す） |
| D | university | string | バックエンド DB から取得（フロントが渡す） |
| E | faculty | string | バックエンド DB から取得（フロントが渡す） |
| F | grade | number | バックエンド DB から取得（フロントが渡す） |
| G | joined_at | datetime | 書き込み日時（GAS が自動付与） |

## GAS web app エンドポイント

### POST（doPost）

**リクエスト（JSON body）**

```json
{
  "circle_sheets_id": "Google Sheets の spreadsheetId",
  "app_uuid": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "display_name": "テスト太郎",
  "university": "○○大学",
  "faculty": "工学部",
  "grade": 2
}
```

**処理の流れ（GAS 内部）**

1. Global Sheet を `app_uuid` で検索して `real_name` を取得する
2. `circle_sheets_id` で指定された Circle Sheet の末尾行に書き込む
3. 完了レスポンスを返す

**レスポンス（成功）**

```json
{ "status": "ok" }
```

**レスポンス（失敗）**

```json
{ "status": "error", "message": "..." }
```

## 認証・セキュリティ

- GAS web app は「自分のアカウントで実行」で公開する
- `circle_sheets_id` はバックエンドの `circles.sheets_id` から取得し、フロントが渡す
- GAS はサービスアカウントで Global Sheet への読み取りアクセスを持つ
- Circle Sheet には GAS のサービスアカウントに編集権限を付与する必要がある

## 呼び出しタイミング

- `POST /circles/{id}/apply` がバックエンドで成功した後にフロントエンドが呼び出す
- すでに Circle Sheet に同じ `app_uuid` が存在する場合は書き込みをスキップする（二重参加防止）

## DB との連携

- `circles.sheets_id`：バックエンドの circles テーブルが持つ Circle Sheet の spreadsheetId
- system_admin がサークル作成時に設定する（`POST /admin/circles` のリクエストに含める）
