package endpoints

import (
	"context"
	"net/http"
	"time"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/modules/encrypt"
	"gitag.ir/armogroup/armo/services/reality/policy"
	accountentpoints "gitag.ir/armogroup/armo/services/reality/services/account/endpoints"
	"github.com/ARmo-BigBang/kit/dtp"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/google/uuid"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
)

type InviteData struct {
	Limit int    `json:"limit,omitempty"`
	Code  string `json:"code,omitempty"`
}

type CreateUserRequest struct {
	Bio                 string     `json:"biography,omitempty" `
	Name                string     `json:"name,omitempty"`
	Phone               string     `json:"phone,omitempty"`
	Email               string     `json:"email,omitempty"`
	Title               string     `json:"title,omitempty"`
	Grade               int        `json:"grade,omitempty"`
	IDCode              string     `json:"id_code,omitempty"`
	LastName            string     `json:"last_name,omitempty"`
	Username            string     `json:"username,omitempty"`
	Nickname            string     `json:"nickname,omitempty"`
	Password            string     `json:"password,omitempty"`
	AvatarUrl           string     `json:"avatar_url,omitempty"`
	Invite              InviteData `json:"invite,omitempty"`
	CountryCode         string     `json:"country_code,omitempty"`
	City                string     `json:"city,omitempty"`
	Roles               []int      `json:"roles,omitempty"`
	Organizations       []int      `json:"organizations,omitempty"`
	Gender              string     `json:"gender,omitempty"`
	DateOfBirth         time.Time  `json:"date_of_birth,omitempty"`
	SuspendedAt         bool       `json:"suspended_at,omitempty"`
	MadeOfficialAt      bool       `json:"made_official_at,omitempty"`
	EmailVerifiedAt     bool       `json:"email_verified_at,omitempty"`
	PhoneVerifiedAt     bool       `json:"phone_verified_at,omitempty"`
	MadeProfilePublicAt bool       `json:"made_profile_public_at,omitempty"`
}

func (r *CreateUserRequest) Validate(req *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"phone":     accountentpoints.PhoneNumberValidity,
		"biography": govalidity.New("biography").MinMaxLength(2, 200).Optional(),
		"name":      govalidity.New("name").Required().MinLength(2).MaxLength(200),
		"email":     govalidity.New("email").Email().Required(),
		"title":     govalidity.New("title").Required().MinMaxLength(2, 200),
		"grade":     govalidity.New("grade"),
		"id_code":   govalidity.New("id_code"),
		"last_name": govalidity.New("last_name").Required().MinLength(2).MaxLength(100),
		// TODO: username should include not only number to prevent reserved phone hack
		"username":               govalidity.New("username").MinMaxLength(2, 200).Required(),
		"nickname":               govalidity.New("nickname").MinMaxLength(2, 200).Optional(),
		"password":               govalidity.New("password").MinLength(6).Required(),
		"avatar_url":             govalidity.New("avatar_url").Optional(),
		"invite":                 govalidity.New("invite").Optional().MinMaxLength(2, 200),
		"country_code":           govalidity.New("country_code").Required(),
		"city":                   govalidity.New("city").Optional(),
		"roles":                  govalidity.New("roles").Optional(),
		"organizations":          govalidity.New("organizations").Optional(),
		"gender":                 govalidity.New("gender").Required(),
		"date_of_birth":          govalidity.New("date_of_birth").Optional(),
		"suspended_at":           govalidity.New("suspended_at").Optional(),
		"made_official_at":       govalidity.New("made_official_at").Optional(),
		"email_verified_at":      govalidity.New("email_verified_at").Optional(),
		"phone_verified_at":      govalidity.New("phone_verified_at").Optional(),
		"made_profile_public_at": govalidity.New("made_profile_public_at").Optional(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"biography":              "بیوگرافی",
			"name":                   "نام",
			"phone":                  "تلفن",
			"email":                  "ایمیل",
			"title":                  "عنوان",
			"grade":                  "درجه",
			"roles":                  "نقش ها",
			"organization":           "سازمان ها",
			"id_code":                "کد شناسایی",
			"last_name":              "نام خانوادگی",
			"username":               "نام کاربری",
			"nickname":               "نام مستعار",
			"password":               "رمز عبور",
			"avatar_url":             "آدرس تصویر آواتار",
			"invite":                 "دعوت",
			"country_code":           "کد کشور",
			"city":                   "شهر",
			"gender":                 "جنسیت",
			"date_of_birth":          "تاریخ تولد",
			"suspended_at":           "تعلیق شده در",
			"made_official_at":       "تاریخ رسمی شدن",
			"email_verified_at":      "تاریخ تأیید ایمیل",
			"phone_verified_at":      "تاریخ تأیید تلفن",
			"made_profile_public_at": "تاریخ عمومی شدن پروفایل",
		},
	)

	errr := govalidity.ValidateBody(req, schema, r)

	if len(errr) > 0 {
		return govalidity.DumpErrors(errr)
	}

	return nil
}

func (s *service) isEnoughCreateData(input CreateUserRequest) bool {
	if input.Name != "" && input.LastName != "" && input.Username != "" && input.Phone != "" && input.Email != "" && input.IDCode != "" &&
		input.CountryCode != "" &&
		input.Gender != "" && (input.DateOfBirth != time.Time{}) && input.PhoneVerifiedAt {
		return true
	}
	return false
}

func (s *service) Create(ctx context.Context, input CreateUserRequest) (models.User, response.ErrorResponse) {

	if !policy.CanCreateUser(ctx) {
		s.logger.With(ctx).Error("شما اجازه دسترسی به این کاربر را ندارید")
		return models.User{}, response.ErrorForbidden(nil, "شما اجازه دسترسی به این کاربر را ندارید")
	}

	profileCompletedAt := dtp.NullTime{Valid: s.isEnoughCreateData(input), Time: time.Now()}

	var hashedPassword string
	hashedPassword, err := encrypt.HashPassword(input.Password)
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.User{}, response.ErrorInternalServerError(nil, "خطایی در سرور رخ داده است")
	}

	var roles []*models.Role
	if len(input.Roles) != 0 {
		err := s.db.Find(&roles, input.Roles).Error
		if err != nil {
			s.logger.With(ctx).Error(err)
			return models.User{}, response.ErrorBadRequest(nil, "نقش‌های ارسالی معتبر نمی‌باشد")
		}
	}

	var organizations []*models.Organization
	if len(input.Organizations) != 0 {
		err := s.db.Find(&organizations, input.Organizations).Error
		if err != nil {
			s.logger.With(ctx).Error(err)
			return models.User{}, response.ErrorBadRequest(nil, "سازمان های ارسالی معتبر نمی‌باشد")
		}
	}

	user := models.User{
		UID:           uuid.New(),
		Name:          input.Name,
		Phone:         input.Phone,
		Title:         input.Title,
		LastName:      input.LastName,
		Password:      hashedPassword,
		CountryCode:   input.CountryCode,
		City:          input.City,
		Roles:         roles,
		Organizations: organizations,
		Gender:        input.Gender,
		DateOfBirth: dtp.NullTime{
			Time:  input.DateOfBirth,
			Valid: input.DateOfBirth != time.Time{},
		},
		Bio:       input.Bio,
		Grade:     input.Grade,
		Nickname:  input.Nickname,
		AvatarUrl: input.AvatarUrl,
		Username: dtp.NullString{
			String: input.Username,
			Valid:  input.Username != "",
		},
		Email: dtp.NullString{
			String: input.Email,
			Valid:  input.Email != "",
		},
		IDCode: dtp.NullString{
			String: input.IDCode,
			Valid:  input.IDCode != "",
		},
		Invite: &models.Invite{
			Code:  dtp.NullString{String: input.Invite.Code, Valid: input.Invite.Code != ""},
			Limit: input.Invite.Limit,
		},
		SuspendedAt: dtp.NullTime{
			Valid: input.SuspendedAt,
			Time:  time.Now(),
		},
		MadeOfficialAt: dtp.NullTime{
			Valid: input.MadeOfficialAt,
			Time:  time.Now(),
		},
		EmailVerifiedAt: dtp.NullTime{
			Valid: input.EmailVerifiedAt,
			Time:  time.Now(),
		},
		PhoneVerifiedAt: dtp.NullTime{
			Valid: input.PhoneVerifiedAt,
			Time:  time.Now(),
		},
		MadeProfilePublicAt: dtp.NullTime{
			Valid: input.MadeProfilePublicAt,
			Time:  time.Now(),
		},
		ProfileCompletedAt: profileCompletedAt,
	}
	err = s.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.User{}, response.GormErrorResponse(err, "خطایی در ثبت کاربر رخ داده است")
	}
	return user, response.ErrorResponse{}
}
