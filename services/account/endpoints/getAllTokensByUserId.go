package endpoints

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"
	"github.com/ARmo-BigBang/kit/response"
)

func (s *service) GetAllTokensByUserId(ctx context.Context, userId string) ([]models.Token, response.ErrorResponse) {
	var tokens []models.Token
	err := s.db.WithContext(ctx).Where("user_id = ?", userId).Find(&tokens).Error

	if !policy.CanGetAllTokensByUserId(ctx, tokens) {
		s.logger.With(ctx).Error("شما اجازه دسترسی به این بخش را ندارید")
		return tokens, response.ErrorForbidden(err, "شما اجازه دسترسی به این بخش را ندارید")
	}

	return tokens, response.ErrorResponse{}
}
