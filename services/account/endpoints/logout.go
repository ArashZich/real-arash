package endpoints

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/response"
)

func (s *service) Logout(ctx context.Context, accessTokens string) (string, response.ErrorResponse) {
	var token models.Token
	tx := s.db.WithContext(ctx).
		Where("access_token = ?", accessTokens).First(&token)
	if tx.Error != nil {
		s.logger.With(ctx).Error(tx.Error)
		return "", response.GormErrorResponse(tx.Error, "خطا در خروج")
	}

	err := tx.Delete(&models.Token{}).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return "", response.GormErrorResponse(err, "خطا در خروج")
	}

	return token.AccessToken, response.ErrorResponse{}
}
