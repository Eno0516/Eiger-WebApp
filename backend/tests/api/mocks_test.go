//go:build api

package api_test

import (
	"context"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
	authuc "github.com/enoro/eiger-webapp/backend/internal/usecase/auth"
)

// =============================================================================
// モック定義
// =============================================================================

type mockVerifier struct {
	lineUserID string
	err        error
}

func (m *mockVerifier) Verify(_ context.Context, _ string) (string, error) {
	return m.lineUserID, m.err
}

type mockJWTGen struct {
	token string
	err   error
}

func (m *mockJWTGen) Generate(_ string, _ bool) (string, error) {
	return m.token, m.err
}

// mockLineAuthUserRepo は LineAuth 用モック。FindByID・Update は呼ばれない想定。
type mockLineAuthUserRepo struct {
	findByLineUserIDUser *entity.User
	findByLineUserIDErr  error
	createdUser          *entity.User
	createErr            error
}

func (m *mockLineAuthUserRepo) FindByID(_ context.Context, _ string) (*entity.User, error) {
	panic("not called in these tests")
}
func (m *mockLineAuthUserRepo) FindByLineUserID(_ context.Context, _ string) (*entity.User, error) {
	return m.findByLineUserIDUser, m.findByLineUserIDErr
}
func (m *mockLineAuthUserRepo) FindAll(_ context.Context) ([]*entity.User, error) {
	panic("not called in these tests")
}
func (m *mockLineAuthUserRepo) Create(_ context.Context, user *entity.User) (*entity.User, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	if m.createdUser != nil {
		return m.createdUser, nil
	}
	return user, nil
}
func (m *mockLineAuthUserRepo) Update(_ context.Context, user *entity.User) (*entity.User, error) {
	panic("not called in these tests")
}
func (m *mockLineAuthUserRepo) Delete(_ context.Context, _ string) error {
	panic("not called in these tests")
}

var _ domainrepo.UserRepository = (*mockLineAuthUserRepo)(nil)

// mockFullUserRepo は全メソッドを設定可能な汎用 UserRepository モック。
type mockFullUserRepo struct {
	findByIDUser  *entity.User
	findByIDErr   error
	findByLineUID *entity.User
	findByLineErr error
	createUser    *entity.User
	createErr     error
	updateUser    *entity.User
	updateErr     error
}

func (m *mockFullUserRepo) FindByID(_ context.Context, _ string) (*entity.User, error) {
	return m.findByIDUser, m.findByIDErr
}
func (m *mockFullUserRepo) FindByLineUserID(_ context.Context, _ string) (*entity.User, error) {
	return m.findByLineUID, m.findByLineErr
}
func (m *mockFullUserRepo) FindAll(_ context.Context) ([]*entity.User, error) {
	panic("not called in these tests")
}
func (m *mockFullUserRepo) Create(_ context.Context, user *entity.User) (*entity.User, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	if m.createUser != nil {
		return m.createUser, nil
	}
	return user, nil
}
func (m *mockFullUserRepo) Update(_ context.Context, user *entity.User) (*entity.User, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	if m.updateUser != nil {
		return m.updateUser, nil
	}
	return user, nil
}
func (m *mockFullUserRepo) Delete(_ context.Context, _ string) error {
	panic("not called in these tests")
}

var _ domainrepo.UserRepository = (*mockFullUserRepo)(nil)

// mockUserProfileRepoAPI は UserProfileRepository の最小モック。
type mockUserProfileRepoAPI struct {
	createProfile *entity.UserProfile
	createErr     error
}

func (m *mockUserProfileRepoAPI) FindByUserID(_ context.Context, _ string) (*entity.UserProfile, error) {
	panic("not called in these tests")
}
func (m *mockUserProfileRepoAPI) Create(_ context.Context, p *entity.UserProfile) (*entity.UserProfile, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	if m.createProfile != nil {
		return m.createProfile, nil
	}
	return p, nil
}
func (m *mockUserProfileRepoAPI) Update(_ context.Context, p *entity.UserProfile) (*entity.UserProfile, error) {
	panic("not called in these tests")
}

var _ domainrepo.UserProfileRepository = (*mockUserProfileRepoAPI)(nil)

// =============================================================================
// テストサーバー設定
// =============================================================================

// apiTestConfig はテストサーバーを構築するためのリポジトリ設定。
// nil のフィールドは対応する UseCase を nil で構成する（呼び出されなければ panic しない）。
type apiTestConfig struct {
	userRepo        domainrepo.UserRepository
	userProfileRepo domainrepo.UserProfileRepository
	lineVerifier    authuc.LINEIDTokenVerifier
	jwtGenUC        authuc.JWTGenerator // LineAuth UseCase 用モック JWT 生成器

	circleRepo       domainrepo.CircleRepository
	circleMemberRepo domainrepo.CircleMemberRepository

	eventRepo            domainrepo.EventRepository
	eventParticipantRepo domainrepo.EventParticipantRepository
}
