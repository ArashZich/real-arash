package endpoints

import (
	"context"
	"errors"
	"fmt"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/response"
	"gorm.io/gorm"
)

func (s *service) CheckIsUniqueField(ctx context.Context, field, value string) (
	bool, models.User, response.ErrorResponse,
) {
	var exists bool
	var user models.User
	var count int64
	err := s.db.WithContext(ctx).Where(fmt.Sprintf("%s = ?", field), value).
		First(&user).Count(&count).Error
	exists = count > 0
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.With(ctx).Error(err)
		return exists, user, response.GormErrorResponse(err, "خطایی در یافتن کاربر رخ داده اس")
	}
	return exists, user, response.ErrorResponse{}
}
