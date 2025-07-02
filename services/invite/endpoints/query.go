package endpoints

import (
	"context"
	"fmt"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"
	"github.com/ARmo-BigBang/kit/response"
)

func (s *service) Query(ctx context.Context, offset, limit int, orderBy, order, query string) (
	[]models.Invite, response.ErrorResponse,
) {
	var invites []models.Invite
	if !policy.CanQueryInvite(ctx) {
		s.logger.With(ctx).Error("شما اجازه دسترسی به این کد معرف را ندارید.")
		return invites, response.ErrorForbidden("شما دسترسی لازم برای انجام این کار را ندارید")
	}

	query = fmt.Sprintf("%%%s%%", query)
	err := s.db.WithContext(ctx).Offset(offset).
		Where("title LIKE ?", query).
		Order(fmt.Sprintf("%s %s", orderBy, order)).
		Limit(limit).Find(&invites).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return invites, response.GormErrorResponse(err, "کد معرف یافت نشد.")
	}

	return invites, response.ErrorResponse{}
}
