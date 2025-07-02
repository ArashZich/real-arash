package endpoints

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/faker"
	"github.com/ARmo-BigBang/kit/response"
)

func (s *service) ApprovePhone(ctx context.Context, code string) (string, response.ErrorResponse) {
	var user models.User
	var phone string
	var ok string

	phone, responseError := s.checkAndDeleteVerificationByCode(ctx, code)
	if responseError.StatusCode != 0 {
		s.logger.With(ctx).Error(responseError)
		return ok, responseError
	}

	_, user, err := s.findUser(ctx, phone)
	if err != nil {
		s.logger.With(ctx).Error(err)
		return ok, response.GormErrorResponse(err, "خطایی در یافتن شماره موبایل رخ داده است")
	}

	user.PhoneVerifiedAt = faker.SQLNow()
	err = s.db.WithContext(ctx).Save(&user).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return ok, response.GormErrorResponse(err, "خطایی در ذخیره شماره موبایل رخ داده است")
	}

	ok = "شماره موبایل با موفقیت تایید شد"

	return ok, response.ErrorResponse{}
}
