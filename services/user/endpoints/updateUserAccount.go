package endpoints

import (
	"context"
	"net/http"
	"time"

	accountentpoints "gitag.ir/armogroup/armo/services/reality/services/account/endpoints"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"
	"github.com/ARmo-BigBang/kit/dtp"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
)

type UpdateUserAccountRequest struct {
	Name                string `json:"name,omitempty"`
	LastName            string `json:"last_name,omitempty"`
	Email               string `json:"email,omitempty"`
	Phone               string `json:"phone,omitempty"`
	AvatarUrl           string `json:"avatar_url,omitempty"`
	CompanyName         string `json:"company_name"`
	MadeProfilePublicAt bool   `json:"made_profile_public_at,omitempty"`
	// Username            string `json:"username,omitempty"`
	// Bio                 string    `json:"biography,omitempty"`
	// IDCode              string    `json:"id_code,omitempty"`
	// Socials             []int     `json:"socials,omitempty"`
	// Nickname            string    `json:"nickname,omitempty"`
	// CountryCode         string    `json:"country_code,omitempty"`
	// City                string    `json:"city,omitempty"`
	// Gender              string    `json:"gender,omitempty"`
	// DateOfBirth         time.Time `json:"date_of_birth,omitempty"`
}

func (req *UpdateUserAccountRequest) Validate(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"phone":                  accountentpoints.PhoneNumberValidity,
		"name":                   govalidity.New("name").Required().MinMaxLength(2, 200),
		"email":                  govalidity.New("email").Required(),
		"last_name":              govalidity.New("last_name").Required().MinMaxLength(2, 200),
		"avatar_url":             govalidity.New("avatar_url").Optional(),
		"company_name":           govalidity.New("company_name").Optional().MinMaxLength(2, 200),
		"made_profile_public_at": govalidity.New("made_profile_public_at").Optional(),
		// "username":               govalidity.New("username").Required().MinMaxLength(2, 200),
		// "biography":              govalidity.New("biography").Optional().MinMaxLength(2, 200),
		// "id_code":                govalidity.New("id_code").Optional(),
		// "socials":                govalidity.New("socials").Optional(),
		// "nickname":               govalidity.New("nickname").MinMaxLength(2, 200).Optional(),
		// "country_code":           govalidity.New("country_code").Required(),
		// "city":                   govalidity.New("city").Optional(),
		// "gender":                 govalidity.New("gender").Required(),
		// "date_of_birth":          govalidity.New("date_of_birth").Required(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"name":                   "نام",
			"email":                  "ایمیل",
			"phone":                  "تلفن",
			"last_name":              "نام خانوادگی",
			"avatar_url":             "آدرس آواتار",
			"company_name":           "نام شرکت",
			"made_profile_public_at": "ساخت پروفایل عمومی",
			// "username":               "نام کاربری",
			// "biography":              "بیوگرافی",
			// "id_code":                "کد ملی",
			// "socials":                "شبکه‌های اجتماعی",
			// "nickname":               "نام مستعار",
			// "country_code":           "کد کشور",
			// "city":                   "شهر",
			// "gender":                 "جنسیت",
			// "date_of_birth":          "تاریخ تولد",
		},
	)

	errr := govalidity.ValidateBody(r, schema, req)
	if len(errr) > 0 {
		return govalidity.DumpErrors(errr)
	}

	return nil
}

func (s *service) isEnoughUpdateAccountData(user models.User, input UpdateUserAccountRequest) bool {
	/*&& input.Gender != "" && input.CountryCode != "" && input.IDCode != "" && input.Username != "" && (input.DateOfBirth != time.Time{})*/
	if input.Name != "" && input.LastName != "" && input.Phone != "" && input.Email != "" && user.PhoneVerifiedAt.Valid {
		return true
	}
	return false
}

func (s *service) UpdateAccount(ctx context.Context, input UpdateUserAccountRequest) (models.User, response.ErrorResponse) {
	var user models.User
	Id := policy.ExtractIdClaim(ctx)
	err := s.db.WithContext(ctx).First(&user, "id", Id).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return user, response.GormErrorResponse(err, "خطا در یافتن کاربر")
	}

	if !policy.CanUpdateAccount(ctx, user) {
		return user, response.ErrorForbidden("شما اجازه دسترسی به این کاربر را ندارید")
	}

	phoneVerifiedAt := user.PhoneVerifiedAt
	if input.Phone != user.Phone {
		phoneVerifiedAt = dtp.NullTime{Valid: false, Time: time.Now()}
	}

	emailVerifiedAt := user.EmailVerifiedAt
	if input.Email != user.Email.String {
		emailVerifiedAt = dtp.NullTime{Valid: false, Time: time.Now()}
	}

	profileCompletedAt := dtp.NullTime{Valid: s.isEnoughUpdateAccountData(user, input), Time: time.Now()}

	user.Name = input.Name
	if input.Email != user.Email.String {
		user.Email = dtp.NullString{
			String: input.Email,
			Valid:  input.Email != "",
		}
	}
	user.Phone = input.Phone
	user.LastName = input.LastName
	user.AvatarUrl = input.AvatarUrl
	user.CompanyName = dtp.NullString{String: input.CompanyName, Valid: input.CompanyName != ""}
	user.PhoneVerifiedAt = phoneVerifiedAt
	user.MadeProfilePublicAt = dtp.NullTime{
		Valid: input.MadeProfilePublicAt,
		Time:  time.Now(),
	}
	user.ProfileCompletedAt = profileCompletedAt
	user.EmailVerifiedAt = emailVerifiedAt
	// user.Username = dtp.NullString{String: input.Username, Valid: input.Username != ""}
	// user.Bio = input.Bio
	// user.IDCode = dtp.NullString{String: input.IDCode, Valid: input.IDCode != ""}
	// user.Nickname = input.Nickname
	// user.CountryCode = input.CountryCode
	// user.City = input.City
	// user.Gender = input.Gender
	// user.DateOfBirth = dtp.NullTime{
	// 	Time:  input.DateOfBirth,
	// 	Valid: input.DateOfBirth != time.Time{},
	// }

	err = s.db.WithContext(ctx).Save(&user).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return user, response.GormErrorResponse(err, "خطا در ذخیره کاربر")
	}

	return user, response.ErrorResponse{}
}
