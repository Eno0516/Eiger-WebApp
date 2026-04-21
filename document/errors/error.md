# エラーハンドリング設計 概要

## 基本方針

- エラーは**構造化**して持ち回り、HTTP ステータスコードとクライアント向けメッセージを一元管理する
- 内部エラー文字列はクライアントへ露出させない
- 詳細は各レイヤーのドキュメントを参照

## エラーの全体フロー

```
[DB / 外部API]
     │  生のエラー（pgx.ErrNoRows 等）
     ▼
[backend: repository]
     │  apperrors センチネルに変換
     ▼
[backend: usecase]
     │  apperrors.NewXxx() でステータスを付与
     ▼
[backend: handler]
     │  HTTPStatusCode / HTTPMessage でレスポンス化
     ▼
HTTP レスポンス  { "error": "..." }  status: 4xx / 5xx
     ▼
[frontend: httpClient interceptor]
     │  AppError に変換
     ▼
[frontend: TanStack Query / notifyError]
     │  status → display を自動決定
     ▼
[frontend: UI]
     Snackbar / Dialog で通知
```

## 詳細ドキュメント

| 対象 | ファイル |
|---|---|
| バックエンド（Go） | [backend/error.md](backend/error.md) |
| フロントエンド（TypeScript） | [frontend/error.md](frontend/error.md) |
