package endpoints

import (
	"context"
	"strconv"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"
	"github.com/ARmo-BigBang/kit/response"
)

func (s *service) Delete(ctx context.Context, ids []int) ([]int, response.ErrorResponse) {
	Id := policy.ExtractIdClaim(ctx)
	id, _ := strconv.Atoi(Id)

	var user models.User
	err := s.db.WithContext(ctx).
		Preload("Invite").
		Preload("Roles").
		Preload("Organizations").
		First(&user, "id", id).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return []int{}, response.GormErrorResponse(err, "خطا در یافتن کاربر")
	}
	if !policy.CanCreateOrganization(ctx, user) {
		s.logger.With(ctx).Error("شما دسترسی حذف سازمان را ندارید")
		return []int{}, response.ErrorForbidden("شما دسترسی حذف سازمان را ندارید")
	}

	err = s.db.WithContext(ctx).
		Where("id IN ?", ids).
		Delete(&models.Organization{}).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return []int{}, response.GormErrorResponse(err, "خطا در حذف سازمان")
	}

	return ids, response.ErrorResponse{}
}
