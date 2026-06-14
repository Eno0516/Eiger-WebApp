# ドキュメント作業ルール

## ルール

- ドキュメントは Markdown で記述する
- 実装の変更を行った場合、関連するドキュメントも合わせて更新する
- 図が必要な場合は Plant UML 記法を使う
- 各機能は `features/feature/login.md` を参考に同じ構成で書く

## ドキュメント構成

```
document/
├── CLAUDE.md                   # このファイル（ドキュメント作業ルール）
│
├── api/                        # OpenAPI 定義
│   ├── openapi.yaml            # ルート定義（エントリーポイント）
│   ├── redocly.yaml            # Redocly 設定
│   ├── components/
│   │   └── schemas/            # 共通スキーマ定義
│   │       ├── Error.yaml
│   │       ├── LoginResponse.yaml
│   │       ├── MessageResponse.yaml
│   │       └── User.yaml
│   └── paths/                  # エンドポイント定義（ルート別）
│       ├── admin/
│       │   └── users/
│       │       └── {id}/
│       │           └── role.yaml
│       ├── auth/
│       │   ├── login.yaml
│       │   ├── register.yaml
│       │   ├── forgot-password.yaml
│       │   ├── resend-verification.yaml
│       │   ├── reset-password.yaml
│       │   └── verify-email.yaml
│       └── users/
│           ├── index.yaml
│           └── {id}.yaml
│
├── db/                         # DB 設計
│   └── migration/              # マイグレーションファイル（予定）
│
├── errors/                     # エラー定義・方針
│   ├── error.md                # エラー共通方針
│   ├── backend/
│   │   └── error.md            # バックエンド側エラー詳細
│   └── frontend/
│       └── error.md            # フロントエンド側エラー詳細
│
├── gas/                        # GAS（Google Apps Script）設計
│   ├── overview.md             # GASの役割・全体像・フロントとの連携方式
│   ├── global_sheet.md         # Global Sheet スキーマ・書き込み仕様
│   └── circle_sheet.md         # Circle Sheet スキーマ・コピー仕様
│
├── features/                   # 機能設計
│   ├── features.md             # 機能一覧
│   └── feature/                # 機能別詳細設計
│       ├── login.md            # ログイン機能（他の機能の書き方の参考にする）
│       ├── LINE.md             # LINE 連携機能
│       └── advertisement.md    # 広告機能
│
├── pages/                      # 画面設計
│   ├── pages.md                # 全画面一覧
│   └── user/                   # ユーザー向けページ詳細
│       └── club_pages.md       # サークル画面の設計
│
├── deploy/                     # デプロイ構成
│   └── deploy.md               # Vercel / Railway / Neon の構成・手順・料金
│
└── tests/                      # テスト戦略
    ├── CLAUDE.md               # テスト作業ルール（テストを書く際は必ず参照）
    ├── tests.md                # テスト戦略詳細（種別・DB モック方針・構成）
    ├── backend/                # バックエンドテスト仕様（予定）
    ├── frontend/               # フロントエンドテスト仕様（予定）
    └── e2e/                    # E2E テスト仕様（予定）
```

## 各ディレクトリの用途早見表

| ディレクトリ | 用途 | 実装時の参照タイミング |
|---|---|---|
| `gas/` | GAS 設計・Sheets 連携仕様 | GAS 実装時・Sheets スキーマ変更時 |
| `api/` | OpenAPI 定義 | API エンドポイントの追加・変更時 |
| `db/` | DB スキーマ・マイグレーション | テーブル設計・変更時 |
| `errors/` | エラーコード・エラー方針 | エラーハンドリング実装時 |
| `features/` | 機能要件・仕様 | 機能実装の開始時 |
| `pages/` | 画面設計・導線 | フロントエンドのページ実装時 |
| `deploy/` | デプロイ構成・手順・料金 | デプロイ作業時・環境変数追加時 |
| `tests/` | テスト戦略・種別・モック方針 | テストを書く際は必ず参照 |

## ページ設計ドキュメント

- ページ設計は `pages/pages.md` に画面一覧を管理する
- 各ページの詳細設計は `pages/{ユーザー種別}/` 配下にページ単位で作成する
