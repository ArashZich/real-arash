package endpoints

import (
	"context"
	"net/http"
	"strconv"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/modules/encrypt"
	"gitag.ir/armogroup/armo/services/reality/policy"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
)

type ChangePasswordRequest struct {
	Password    string `json:"password,omitempty"`
	NewPassword string `json:"new_password,omitempty"`
}

func (r *ChangePasswordRequest) Validate(req *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"password":     govalidity.New("password").Required(),
		"new_password": PasswordValidity,
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"password":     "رمز عبور",
			"new_password": "رمز عبور جدید",
		},
	)

	errr := govalidity.ValidateBody(req, schema, r)
	if len(errr) > 0 {
		dumpedErrors := govalidity.DumpErrors(errr)
		return dumpedErrors
	}

	return nil
}

func (s *service) ChangePassword(
	ctx context.Context, input ChangePasswordRequest,
) (
	LoginResponse,
	response.ErrorResponse,
) {
	var user models.User

	var loginResponse LoginResponse

	Id := policy.ExtractIdClaim(ctx)
	err := s.db.WithContext(ctx).First(&user, "id", Id).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return loginResponse, response.GormErrorResponse(err, "شما در سیستم ثبت نام نکرده‌اید")
	}

	if !user.ValidatePassword(input.Password) {
		return LoginResponse{}, response.ErrorBadRequest(err, "رمز عبور اشتباه می‌باشد")
	}

	var hashedPassword string
	hashedPassword, err = encrypt.HashPassword(input.NewPassword)
	if err != nil {
		s.logger.With(ctx).Error(err)
		return loginResponse, response.ErrorInternalServerError(err, "خطایی در سرور رخ داده است")
	}

	user.Password = hashedPassword
	err = s.db.WithContext(ctx).Save(&user).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		err := response.GormErrorResponse(err, "خطایی در ذخیره‌سازی رخ داده است")
		return loginResponse, err
	}

	userIDString := strconv.Itoa(int(user.ID))
	_, err = s.DeleteTokens(ctx, []string{userIDString})
	if err != nil {
		s.logger.Error(err)
		return LoginResponse{}, response.ErrorInternalServerError(nil, "خطایی در حذف توکن‌ها رخ داده است")
	}

	accessTkn, refreshTkn, responseError := s.generateTokens(ctx, user, "")
	if responseError.StatusCode != 0 {
		s.logger.With(ctx).Error(responseError)
		return loginResponse, responseError
	}
	loginResponse = LoginResponse{
		User:         user,
		AccessToken:  accessTkn,
		RefreshToken: refreshTkn,
	}
	return loginResponse, response.ErrorResponse{}
}
