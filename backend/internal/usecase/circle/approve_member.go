package circle

import (
	"context"
	"fmt"

	"github.com/enoro/eiger-webapp/backend/internal/domain/entity"
	domainrepo "github.com/enoro/eiger-webapp/backend/internal/domain/repository"
	"github.com/enoro/eiger-webapp/backend/pkg/apperrors"
)

type ApproveMemberInput struct {
	CircleID     string
	TargetUserID string
}

// ApproveMemberUseCase はサークル admin が visitor を member に昇格させる。
// 昇格後、対象ユーザーにアプリ内通知を送る。
type ApproveMemberUseCase struct {
	circleMemberRepo domainrepo.CircleMemberRepository
	// TODO: Spreadsheet から Form 回答を取得する SheetsClient を追加する
}

func NewApproveMemberUseCase(repo domainrepo.CircleMemberRepository) *ApproveMemberUseCase {
	return &ApproveMemberUseCase{circleMemberRepo: repo}
}

// Execute は対象ユーザーを visitor から member に昇格させ、本人に通知する。
func (uc *ApproveMemberUseCase) Execute(ctx context.Context, input ApproveMemberInput) error {
	member, err := uc.circleMemberRepo.FindByCircleAndUser(ctx, input.CircleID, input.TargetUserID)
	if err != nil {
		return err
	}
	if member.Role != entity.CircleRoleVisitor {
		return apperrors.NewBadRequest("target user is not a visitor", nil)
	}
	member.Role = entity.CircleRoleMember
	if err := uc.circleMemberRepo.Upsert(ctx, member); err != nil {
		return fmt.Errorf("upsert member: %w", err)
	}
	// NOTE: 対象ユーザーへの通知は別途実装予定
	return nil
}
