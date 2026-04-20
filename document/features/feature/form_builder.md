# フォームビルダー機能 設計ドキュメント

## 概要

サークル admin がアプリ内でフォームを作成し、メンバーが LIFF 上で回答する機能。
Google Forms を使わず、アプリ内で完結させることで以下を実現する。

- LIFF SDK により LINE 認証済みユーザーとして自動識別（app UUID の手動入力不要）
- バックエンドが API 送信を制御できる
- サークルごとにフォームを管理できる

---

## フィールドタイプ

| タイプ | 説明 |
|---|---|
| `text` | テキスト入力欄（1行） |
| `single_select` | 単一選択（ラジオボタン） |
| `multi_select` | 複数選択（チェックボックス） |

---

## データフロー

### フォーム作成（admin）

```
admin がアプリ内でフォームを作成
  → POST /circles/{id}/forms
  → DB: forms + form_fields に保存
```

### フォーム回答（member）

```
member が活動ページからフォームリンクをタップ
  → LIFF で開く
  → LIFF SDK で LINE ID 取得 → JWT 発行済みのため app UUID が確定済み
  → GET /forms/{formId}/responses/me で過去の回答を取得し、あれば入力欄に pre-fill
  → 締め切りを過ぎている場合は回答フォームを非表示にして「締め切り済み」を表示
  → 回答送信: PUT /forms/{formId}/responses/me（JWT 付き）
  → 初回: form_responses を INSERT → form_response_fields を INSERT
  → 再回答: form_response_fields を UPDATE（form_responses の submitted_at も更新）
```

### 集計（admin）

```
admin が回答一覧を閲覧
  → GET /circles/{id}/forms/{formId}/responses
  → アプリ内で表示 or Circle Sheet にエクスポート（GAS 経由）
```

---

## DB スキーマ

```sql
-- フォーム定義
CREATE TABLE forms (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    circle_id   UUID NOT NULL REFERENCES circles(id) ON DELETE CASCADE,
    title       VARCHAR(255) NOT NULL,
    created_by  UUID NOT NULL REFERENCES users(id),
    closes_at   TIMESTAMPTZ,     -- NULL = 締め切りなし
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- フォームフィールド定義
CREATE TABLE form_fields (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    form_id     UUID NOT NULL REFERENCES forms(id) ON DELETE CASCADE,
    label       VARCHAR(255) NOT NULL,
    type        VARCHAR(20) NOT NULL CHECK (type IN ('text', 'single_select', 'multi_select')),
    options     JSONB,           -- single_select / multi_select の選択肢 ["A", "B", "C"]
    is_required BOOLEAN NOT NULL DEFAULT false,
    sort_order  SMALLINT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 回答セッション（1ユーザー1フォームにつき1レコード・上書き可）
CREATE TABLE form_responses (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    form_id      UUID NOT NULL REFERENCES forms(id) ON DELETE CASCADE,
    user_id      UUID NOT NULL REFERENCES users(id),
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (form_id, user_id)   -- 1フォームにつき1レコードのみ（再回答は UPDATE）
);

-- 各フィールドの回答値
CREATE TABLE form_response_fields (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    response_id UUID NOT NULL REFERENCES form_responses(id) ON DELETE CASCADE,
    field_id    UUID NOT NULL REFERENCES form_fields(id),
    value       TEXT            -- text: そのまま / select: 選択肢の値をカンマ区切り
);
```

---

## API エンドポイント（予定）

| メソッド | パス | 権限 | 説明 |
|---|---|---|---|
| POST | `/circles/{id}/forms` | circle_admin | フォーム作成 |
| GET | `/circles/{id}/forms` | member | フォーム一覧取得 |
| GET | `/forms/{formId}` | member | フォーム詳細取得（フィールド・締め切り含む） |
| GET | `/forms/{formId}/responses/me` | member | 自分の回答取得（re-fill 用） |
| PUT | `/forms/{formId}/responses/me` | member | 回答送信・上書き（締め切り後は 403） |
| GET | `/circles/{id}/forms/{formId}/responses` | circle_admin | 回答一覧取得 |

---

## GAS との連携

回答データは DB に保存するが、admin が Circle Sheet にエクスポートしたい場合は GAS 経由で書き出す。

- `GET /circles/{id}/forms/{formId}/responses` の結果をフロントが GAS に渡す
- GAS が Circle Sheet の所定の列に書き込む
- 個人情報（本名）が必要な場合は GAS が Global Sheet から JOIN して付与する

詳細は [circle_sheet.md](../../gas/circle_sheet.md) を参照。

---

## 決定事項まとめ

| 項目 | 決定内容 |
|---|---|
| 回答回数 | 1フォームにつき1レコード。再回答は上書き（`updated_at` を更新） |
| 再回答時の UX | フォームを開いたとき過去の回答を pre-fill して表示する |
| 締め切り | `forms.closes_at` で管理。NULL = 締め切りなし |
| 締め切り後の挙動 | フロントはフォームを非表示、バックエンドは PUT 時に 403 を返す |
