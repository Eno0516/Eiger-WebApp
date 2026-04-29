package auth

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

// RegisterInput は初期登録フォームから受け取るデータ。
// 本名・学籍番号はフロントエンドから GAS 経由で直接 Sheets に書き込むため含まない。
type RegisterInput struct {
	UserID      string // JWT から取得した users.id
	DisplayName string // 部内あだ名（app DB の user_profiles に保存）
	University  string // 大学（app DB の user_profiles に保存）
	Faculty     string // 学部（app DB の user_profiles に保存）
	Grade       int    // 学年（app DB の user_profiles に保存）
}

type RegisterOutput struct {
	Profile *entity.UserProfile
}

// RegisterUseCase は初期登録フォームの送信を受け付けるユースケース。
// user_profiles を作成し、users.is_registered を true に更新する。
// 本名・学籍番号はフロントエンドから GAS 経由で直接 Sheets に書き込むため、このユースケースは関与しない。
type RegisterUseCase struct {
	userRepo        domainrepo.UserRepository
	userProfileRepo domainrepo.UserProfileRepository
	log             zerolog.Logger
}

func NewRegisterUseCase(
	userRepo domainrepo.UserRepository,
	userProfileRepo domainrepo.UserProfileRepository,
	log zerolog.Logger,
) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo:        userRepo,
		userProfileRepo: userProfileRepo,
		log:             log,
	}
}

// Execute は初期登録を完了させる。
// 1. UserProfile を作成する（display_name・大学・学部・学年）
// 2. users.is_registered を true に更新する
// すでに登録済みの場合は apperrors.NewConflict を返す。
func (uc *RegisterUseCase) Execute(ctx context.Context, input RegisterInput) (*RegisterOutput, error) {
	u, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}
	if u.IsRegistered {
		return nil, apperrors.NewConflict("already registered", nil)
	}
	profile := buildProfile(input)
	if err := profile.Validate(); err != nil {
		return nil, apperrors.NewBadRequest("validation failed", err)
	}
	created, err := uc.userProfileRepo.Create(ctx, profile)
	if err != nil {
		return nil, apperrors.NewInternal("create user profile", err)
	}
	u.IsRegistered = true
	if _, err := uc.userRepo.Update(ctx, u); err != nil {
		return nil, apperrors.NewInternal("update user", err)
	}
	return &RegisterOutput{Profile: created}, nil
}

// buildProfile は RegisterInput から UserProfile エンティティを生成する。
func buildProfile(input RegisterInput) *entity.UserProfile {
	return &entity.UserProfile{
		UserID:      input.UserID,
		DisplayName: input.DisplayName,
		University:  input.University,
		Faculty:     input.Faculty,
		Grade:       input.Grade,
	}
}
