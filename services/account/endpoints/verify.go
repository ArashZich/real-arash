package endpoints

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/response"
)

// rethink about getting the phone from the user of from the verification field it's self
func (s *service) checkAndDeleteVerificationBySessionCodeAndPhone(
	ctx context.Context, sessionCode string, phone string,
) response.ErrorResponse {
	var verification models.Verification
	err := s.db.WithContext(ctx).First(&verification, "session_code", sessionCode).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return response.GormErrorResponse(err, "خطای در اعتبار سنجی")
	}

	if verification.NotValidPhone(phone) || verification.Expired() {
		s.logger.Info(verification.NotValidPhone(phone), verification.Expired())
		return response.ErrorBadRequest(err, "زمان کد شما به پایان رسیده است لطفا دوباره امتحان کنید")
	}

	err = s.db.WithContext(ctx).Delete(&models.Verification{}, verification.ID).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return response.GormErrorResponse(err, "خطا در اعتبارسنجی")
	}

	return response.ErrorResponse{}
}

func (s *service) checkAndDeleteVerificationByCode(ctx context.Context, code string) (string, response.ErrorResponse) {
	var phone string
	var verification models.Verification
	err := s.db.WithContext(ctx).First(&verification, "code", code).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return phone, response.GormErrorResponse(err, "خطا در اعتبارسنجی. کد شما یافت نشد")
	}
	if verification.Expired() {
		s.logger.Info(verification.Phone, verification.Expired())
		return phone, response.ErrorBadRequest(err, "زمان کد شما به پایان رسیده است لطفا دوباره امتحان کنید")
	}

	phone = verification.Phone

	err = s.db.WithContext(ctx).Delete(&models.Verification{}, verification.ID).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return phone, response.GormErrorResponse(err, "خطا در اعتبارسنجی")
	}
	return phone, response.ErrorResponse{}
}
