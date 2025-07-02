package endpoints

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"
	"github.com/ARmo-BigBang/kit/response"
)

func (s *service) Get(ctx context.Context, uid string) (models.User, response.ErrorResponse) {
	var user models.User
	err := s.db.WithContext(ctx).
		Preload("Invite").
		Preload("Roles").
		Preload("Organizations").
		Preload("Organizations.Category").
		Preload("Organizations.Packages").
		Preload("Organizations.Products").
		Preload("Organizations.Packages.Plan").
		Preload("Organizations.Packages.Category").
		First(&user, "uid = ?", uid).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.User{}, response.GormErrorResponse(err, "خطا در یافتن کاربر")
	}

	if !policy.CanQueryUsers(ctx) {
		s.logger.With(ctx).Error("شما اجازه دسترسی به این کاربر را ندارید")
		return models.User{}, response.ErrorForbidden("شما اجازه دسترسی به این کاربر را ندارید")
	}
	return user, response.ErrorResponse{}
}
