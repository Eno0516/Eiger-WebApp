//go:build api

package api_test

// APIテストの方針:
// - DBは使わない（モック or testcontainers で差し替え）
// - APIサーバーだけ起動してリクエスト/レスポンスを検証する
// - フロントとのI/F（スキーマ・ステータスコード・エラー形式）が正しいかを確認する
