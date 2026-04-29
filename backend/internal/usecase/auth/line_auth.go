package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

type LineAuthInput struct {
	IDToken string // LINE が発行した ID トークン
}

type LineAuthOutput struct {
	Token     string // アプリが発行する JWT
	IsNewUser bool   // 初回ログインかどうか（true なら display_name 未設定）
}

// LineAuthUseCase は LINE ID トークンを検証し、ユーザーを取得または新規作成して JWT を返す。
type LineAuthUseCase struct {
	userRepo domainrepo.UserRepository
	verifier LINEIDTokenVerifier
	jwtGen   JWTGenerator
	log      zerolog.Logger
}

func NewLineAuthUseCase(
	userRepo domainrepo.UserRepository,
	verifier LINEIDTokenVerifier,
	jwtGen JWTGenerator,
	log zerolog.Logger,
) *LineAuthUseCase {
	return &LineAuthUseCase{userRepo: userRepo, verifier: verifier, jwtGen: jwtGen, log: log}
}

// Execute は LINE ID トークンを検証し、ユーザーを取得または新規作成して JWT を返す。
// 初回ログイン時は display_name が空のユーザーを作成する（フロント側で初回設定画面に誘導する）。
func (uc *LineAuthUseCase) Execute(ctx context.Context, input LineAuthInput) (*LineAuthOutput, error) {
	lineUserID, err := uc.verifier.Verify(ctx, input.IDToken)
	if err != nil {
		return nil, apperrors.NewUnauthorized("invalid id token", err)
	}

	var isNewUser bool
	u, err := uc.userRepo.FindByLineUserID(ctx, lineUserID)
	if err != nil {
		if !errors.Is(err, apperrors.ErrNotFound) {
			return nil, err
		}
		// 初回ログイン: ユーザーを新規作成する
		uc.log.Info().Str("line_user_id", lineUserID).Msg("初回ログイン: ユーザーを新規作成します")
		u, err = uc.userRepo.Create(ctx, &entity.User{LineUserID: lineUserID})
		if err != nil {
			return nil, fmt.Errorf("create user: %w", err)
		}
		isNewUser = true
	}

	token, err := uc.jwtGen.Generate(u.ID, u.IsSystemAdmin)
	if err != nil {
		return nil, fmt.Errorf("generate jwt: %w", err)
	}
	return &LineAuthOutput{Token: token, IsNewUser: isNewUser}, nil
}
