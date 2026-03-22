# Eiger WebApp

## 概要

モノレポ構成の Web アプリケーション。  
バックエンドは Go、フロントエンドは React（npm workspaces）で構成される。

---

## 技術スタック

| カテゴリ | 技術 |
|----------|------|
| **バックエンド** | Go 1.23 / chi / PostgreSQL (pgx) |
| **フロントエンド** | React 18 / TypeScript / Vite / TanStack Query |
| **インフラ** | Docker / docker-compose |
| **マイグレーション** | golang-migrate |
| **ログ** | zerolog |
| **設定管理** | viper |

---

## ディレクトリ構成

```
/
├── backend/                    # Go バックエンド
│   ├── cmd/                    # エントリーポイント
│   ├── internal/               # アプリケーションコード（Clean Architecture）
│   │   ├── domain/             # Entity / ValueObject / Repository interface
│   │   ├── usecase/            # ビジネスロジック
│   │   ├── handler/            # HTTPハンドラ
│   │   ├── repository/         # Repository 実装
│   │   └── infrastructure/     # DB・外部サービスの具体実装
│   └── pkg/                    # 外部公開可能なユーティリティ
│
├── frontend/                   # フロントエンド（npm workspaces）
│   ├── apps/
│   │   ├── user/               # ユーザー向けアプリ
│   │   └── admin/              # 管理画面アプリ
│   └── packages/
│       └── shared/             # 共通コンポーネント・ロジック
│
├── tests/                      # クロスサービステスト
│   ├── api/                    # APIコントラクトテスト（build tag: api）
│   ├── integration/            # 統合テスト・DB+API（build tag: integration）
│   ├── e2e/                    # E2Eテスト・全スタック起動（build tag: e2e）
│   └── testdata/               # テスト用フィクスチャ
│
├── document/                   # 設計ドキュメント
├── docker-compose.yml          # 本番用
├── docker-compose.dev.yml      # 開発用
└── Makefile                    # 各種コマンド
```

---

## アーキテクチャ方針

- **Clean Architecture + DDD** をバックエンドに採用
- 依存の方向: `handler → usecase → domain ← repository`
- フロントエンドは **モノレポ（npm workspaces）** で user / admin を管理

### テスト戦略

| 種別 | 場所 | 内容 |
|------|------|------|
| 単体テスト | 各パッケージそば (`*_test.go`) | モック使用・DB不要 |
| APIテスト | `tests/api/` | APIサーバー単体・スキーマ検証 |
| 統合テスト | `tests/integration/` | DB + API・複数層またぎ |
| E2Eテスト | `tests/e2e/` | DB + API + フロント・操作シナリオ再現 |

---

## 開発環境のセットアップ

```bash
# 環境変数のセットアップ
cp .env.example .env

# Docker で起動
docker compose -f docker-compose.dev.yml up
```

詳細は各ディレクトリの README を参照。
