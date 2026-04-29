package lineauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

// Verifier は LINE ID トークンの検証を行う。
// usecase/auth の LINEIDTokenVerifier interface を実装する。
type Verifier struct {
	channelID string // LINE_CHANNEL_ID 環境変数から渡す
}

func NewVerifier(channelID string) *Verifier {
	return &Verifier{channelID: channelID}
}

// Verify は ID トークンを検証して LineUserID（sub クレーム）を返す。
// トークンが無効・期限切れ・channelID 不一致の場合はエラーを返す。
func (v *Verifier) Verify(ctx context.Context, idToken string) (string, error) {
	resp, err := http.PostForm("https://api.line.me/oauth2/v2.1/verify", url.Values{
		"id_token":  {idToken},
		"client_id": {v.channelID},
	})
	if err != nil {
		return "", apperrors.NewUnauthorized("LINE ID token verification request failed", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errBody struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&errBody)
		return "", apperrors.NewUnauthorized(
			fmt.Sprintf("LINE token verification failed: %s", errBody.ErrorDescription), nil)
	}

	var result struct {
		Sub string `json:"sub"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", apperrors.NewUnauthorized("failed to decode LINE verify response", err)
	}
	if result.Sub == "" {
		return "", apperrors.NewUnauthorized("LINE token missing sub claim", nil)
	}
	return result.Sub, nil
}
