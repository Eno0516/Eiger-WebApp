// Package membership は、「サークルのメンバーか」の判定の入口。
// 判定そのものは本体(認証・サークル管理)が持つため、ミニアプリ側は interface だけを定義し、
// 本体が実装を渡す(Register の deps)。
package membership

import (
	"context"

	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain"
)

type Directory interface {
	IsMember(ctx context.Context, circleID, memberID string) (bool, error)
}

// Require は、memberID が circleID のメンバーでなければ domain.ErrNotMember を返す。
func Require(ctx context.Context, dir Directory, circleID, memberID string) error {
	ok, err := dir.IsMember(ctx, circleID, memberID)
	if err != nil {
		return err
	}
	if !ok {
		return domain.ErrNotMember
	}
	return nil
}
