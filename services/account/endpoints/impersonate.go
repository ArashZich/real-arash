package endpoints

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"
	"github.com/ARmo-BigBang/kit/response"
)

func (s *service) Impersonate(ctx context.Context, id string, currentAccessToken string) (LoginResponse, response.ErrorResponse,
) {
	var loginResponse LoginResponse

	var user models.User
	err := s.db.WithContext(ctx).
		Preload("Roles").
		Preload("Invite").
		First(&user, "id", id).Error

	if !policy.CanImpersonate(ctx) {
		s.logger.With(ctx).Error("شما اجازه دسترسی به این بخش را ندارید")
		return loginResponse, response.ErrorForbidden(err, "شما اجازه دسترسی به این بخش را ندارید")
	}
	if err != nil {
		s.logger.With(ctx).Error(err)
		return loginResponse, response.GormErrorResponse(err, "خطا در یافتن کاربر")
	}

	accessToken, refreshToken, responseError := s.generateTokens(ctx, user, currentAccessToken)
	if responseError.StatusCode != 0 {
		s.logger.With(ctx).Error(responseError)
		return loginResponse, responseError
	}
	loginResponse = LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return loginResponse, response.ErrorResponse{}
}
