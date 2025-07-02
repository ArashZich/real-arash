package endpoints

import (
	"context"
	"net/http"
	"time"

	"gitag.ir/armogroup/armo/services/reality/config"
	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/models/aggregation"
	"gitag.ir/armogroup/armo/services/reality/modules/encrypt"
	"github.com/ARmo-BigBang/kit/dtp"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govalidityl"
	"github.com/hoitek-go/govalidity/govaliditym"
)

type RegisterRequest struct {
	Name        string `json:"name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Password    string `json:"password,omitempty"`
	SessionCode string `json:"session_code,omitempty"`
	InviteCode  string `json:"invite_code,omitempty"`
}

func (r *RegisterRequest) Validate(req *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"phone":        PhoneNumberValidity,
		"password":     PasswordValidity,
		"name":         govalidity.New("name").MinLength(2).MaxLength(200).Required(),
		"last_name":    govalidity.New("last_name").MinLength(2).MaxLength(200).Required(),
		"session_code": govalidity.New("session_code").Required(),
		"invite_code":  govalidity.New("invite_code").AlphaNum(govalidityl.EnUS).Optional(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"name":         "نام",
			"last_name":    "نام خانوادگی",
			"phone":        "شماره موبایل",
			"password":     "رمز عبور",
			"session_code": "کد ارسالی",
			"invite_code":  "کد دعوت",
		},
	)

	errr := govalidity.ValidateBody(req, schema, r)
	if len(errr) > 0 {
		dumpedErrors := govalidity.DumpErrors(errr)
		return dumpedErrors
	}

	return nil
}

func (s *service) Register(ctx context.Context, input RegisterRequest) (LoginResponse, response.ErrorResponse) {
	var user models.User
	var exists bool
	var res LoginResponse

	exists, user, err := s.findUser(ctx, input.Phone)
	if exists {
		return res, response.ErrorBadRequest(nil, "کاربری با این نام کاربری یافت شد")
	}

	if err != nil {
		s.logger.With(ctx).Error(err)
		return res, response.GormErrorResponse(err, "خطایی رخ داده است")
	}

	if config.AppConfig.RequiredSignUpInviteCode {
		if input.InviteCode == "" {
			return res, response.ErrorBadRequest(nil, "کد دعوت الزامی است")
		}

		invite, er := s.useInviteCode(ctx, input.InviteCode)
		if er.StatusCode != 0 {
			s.logger.With(ctx).Error(err)
			return res, er
		}

		registerInvite := aggregation.RegisterInvite{
			HostID: invite.UserID,
			UserID: int(user.ID),
		}

		err := s.db.WithContext(ctx).Create(&registerInvite).Error
		if err != nil {
			s.logger.With(ctx).Error(err)
			return res, response.GormErrorResponse(err, "خطایی در ساخت ارتباط بین کاربر و میزبان رخ داده است")
		}
	}

	eerr := s.checkAndDeleteVerificationBySessionCodeAndPhone(ctx, input.SessionCode, input.Phone)
	if eerr.StatusCode != 0 {
		s.logger.With(ctx).Error(eerr)
		return res, eerr
	}

	var hashedPassword string
	hashedPassword, err = encrypt.HashPassword(input.Password)
	if err != nil {
		s.logger.With(ctx).Error(err)
		return res, response.ErrorInternalServerError(nil, "خطایی در ساخت رمز عبور رخ داده است")
	}

	var roles []models.Role
	err = s.db.WithContext(ctx).Find(&roles).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return res, response.GormErrorResponse(err, "خطایی در یافتن نقش‌ ها رخ داده است")
	}

	user = models.User{
		Name:            input.Name,
		LastName:        input.LastName,
		Phone:           input.Phone,
		Password:        hashedPassword,
		Username:        dtp.NullString{String: input.Phone, Valid: true},
		PhoneVerifiedAt: dtp.NullTime{Time: time.Now(), Valid: true},
		Roles:           []*models.Role{&roles[2]},
	}
	err = s.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return res, response.GormErrorResponse(err, "خطایی در ثبت کاربر رخ داده است")
	}

	accessToken, refreshToken, errr := s.generateTokens(ctx, user, "")
	if errr.StatusCode != 0 {
		return res, errr
	}

	res = LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return res, response.ErrorResponse{}
}
