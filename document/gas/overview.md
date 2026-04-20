# GAS（Google Apps Script）概要

## 役割

GAS はバックエンドを経由させたくない個人情報（本名・学籍番号）の Sheets 書き込みを担当する。
アプリバックエンド（Go）は個人情報に一切触れない。

## 責務の分担

| 処理 | 担当 | 理由 |
|---|---|---|
| 本名・学籍番号 → Global Sheet | GAS | 個人情報をバックエンドに通さない |
| display_name・大学・学部・学年 → DB | バックエンド（Go） | 個人情報ではない、API で管理 |
| Global Sheet → Circle Sheet コピー | GAS | 個人情報の連鎖、バックエンド不要 |

## フロントエンドからの呼び出し方式

フロントエンドが GAS と バックエンドの両方を並行して呼び出す（JWT は GAS に渡さない）。

```
フロント（/sign_up ページ）
  ├─ GAS web app を呼ぶ（本名・学籍番号・app UUID）
  │     └─ GAS が Global Sheet に書き込む
  └─ POST /auth/signup（display_name・大学・学部・学年 + JWT）
        └─ バックエンドが user_profiles を作成・is_registered = true
```

サークル参加時：

```
フロント（サークル申し込み後）
  ├─ POST /circles/{id}/apply → バックエンドで circle_members 作成
  └─ GAS web app を呼ぶ（app UUID + circle_sheets_id）
        └─ GAS が Global Sheet から本名を検索 → Circle Sheet に書き込む
```

## GAS web app 一覧

| web app | エンドポイント定義 |
|---|---|
| 初期登録データ書き込み | [global_sheet.md](global_sheet.md) |
| サークルシートへのコピー | [circle_sheet.md](circle_sheet.md) |
