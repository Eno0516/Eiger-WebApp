package dto

type LineAuthRequest struct {
	IDToken string `json:"id_token"` // LINE が発行した ID トークン
}

// RegisterRequest は POST /auth/signup のリクエストボディ。
// 個人情報（RealName・StudentID）は Sheets に書き込み、app DB には保存しない。
type RegisterRequest struct {
	DisplayName string `json:"display_name"` // 部内あだ名（app DB に保存）
	University  string `json:"university"`   // 大学（app DB に保存）
	Faculty     string `json:"faculty"`      // 学部（app DB に保存）
	Grade       int    `json:"grade"`        // 学年（app DB に保存）
	RealName    string `json:"real_name"`    // 本名（Sheets のみ）
	StudentID   string `json:"student_id"`   // 学籍番号（Sheets のみ）
}

type UpdateDisplayNameRequest struct {
	DisplayName string `json:"display_name"`
}

type CreateCircleRequest struct {
	Name               string `json:"name"`
	Description        string `json:"description"`
	ApplicationFormURL string `json:"application_form_url"`
	SpreadsheetID      string `json:"spreadsheet_id"`
}

type UpdateCircleMemberRoleRequest struct {
	Role string `json:"role"` // "visitor" / "member" / "admin"
}

type CreateEventRequest struct {
	Title        string  `json:"title"`
	Date         string  `json:"date"` // "YYYY-MM-DD"
	Description  string  `json:"description"`
	FormID       *string `json:"form_id"`
	LiffURL      string  `json:"liff_url"`
	LineGroupURL string  `json:"line_group_url"`
}

type UpdateParticipantStatusRequest struct {
	Status string `json:"status"` // "applied" / "confirmed" / "rejected"
}

type BulkUpdateParticipantsRequest struct {
	Participants []ParticipantStatusItem `json:"participants"`
}

type ParticipantStatusItem struct {
	UserID string `json:"user_id"`
	Status string `json:"status"`
}
