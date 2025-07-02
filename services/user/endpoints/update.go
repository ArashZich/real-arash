package endpoints

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/modules/encrypt"
	"gitag.ir/armogroup/armo/services/reality/policy"
	accountentpoints "gitag.ir/armogroup/armo/services/reality/services/account/endpoints"
	"gitag.ir/armogroup/armo/services/reality/services/role"
	"github.com/ARmo-BigBang/kit/dtp"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
	"gorm.io/gorm"
)

type UpdateUserRequest struct {
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
	Gender              string     `json:"gender,omitempty"`
	DateOfBirth         time.Time  `json:"date_of_birth,omitempty"`
	SuspendedAt         bool       `json:"suspended_at,omitempty"`
	MadeOfficialAt      bool       `json:"made_official_at,omitempty"`
	EmailVerifiedAt     bool       `json:"email_verified_at,omitempty"`
	PhoneVerifiedAt     bool       `json:"phone_verified_at,omitempty"`
	Roles               []int      `json:"roles,omitempty"`
	MadeProfilePublicAt bool       `json:"made_profile_public_at,omitempty"`
}

func (r *UpdateUserRequest) Validate(req *http.Request) govalidity.ValidityResponseErrors {
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
		"password":               govalidity.New("password").Optional(),
		"avatar_url":             govalidity.New("avatar_url").Optional(),
		"invite":                 govalidity.New("invite").Optional().MinMaxLength(2, 200),
		"country_code":           govalidity.New("country_code").Required(),
		"city":                   govalidity.New("city").Optional(),
		"roles":                  govalidity.New("roles").Optional(),
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
			"roles":                  "نقش ها",
			"phone":                  "تلفن",
			"email":                  "ایمیل",
			"title":                  "عنوان",
			"grade":                  "درجه",
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

func (s *service) isEnoughUpdateData(input UpdateUserRequest) bool {
	if input.Name != "" && input.LastName != "" && input.Username != "" && input.Phone != "" && input.Email != "" && input.IDCode != "" &&
		input.CountryCode != "" &&
		input.Gender != "" && (input.DateOfBirth != time.Time{}) && input.PhoneVerifiedAt {
		return true
	}
	return false
}

func (s *service) Update(ctx context.Context, id string, input UpdateUserRequest) (models.User, response.ErrorResponse) {
	var user models.User
	tx := s.db.WithContext(ctx).Begin() // Begin a new transaction
	if tx.Error != nil {
		return models.User{}, response.GormErrorResponse(tx.Error, "خطایی در سرور رخ داده است")
	}

	if !policy.CanUpdateUser(ctx) {
		s.logger.With(ctx).Error(" شما اجازه دسترسی به این کاربر را ندارید")
		return models.User{}, response.ErrorForbidden("شما اجازه دسترسی به این کاربر را ندارید")
	}
	err := tx.Preload("Roles").First(&user, "id", id).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.User{}, response.GormErrorResponse(err, "خطا در یافتن کاربر")
	}

	profileCompletedAt := dtp.NullTime{Valid: s.isEnoughUpdateData(input), Time: time.Now()}

	if input.Password != "" {
		var hashedPassword string
		hashedPassword, err = encrypt.HashPassword(input.Password)
		if err != nil {
			return models.User{}, response.ErrorInternalServerError("خطایی در سرور رخ داده است")
		}
		user.Password = hashedPassword
	}

	//check the unique fields not to be duplicated and if user object found and the user id is not the same with the input id return/proper/error
	if input.Email != "" && input.Email != user.Email.String {
		var tmpUser models.User
		err = tx.First(&tmpUser, "email = ? AND id != ?", input.Email, id).Error
		if err != gorm.ErrRecordNotFound && err != nil {
			s.logger.With(ctx).Error(err)
			return models.User{}, response.GormErrorResponse(err, "خطا در یافتن کاربر")
		}
		if err == nil && tmpUser.ID != user.ID {
			return models.User{}, response.GormErrorResponse(err, "ایمیل وارد شده تکراری است")
		}
	}

	if input.Username != "" {
		var tmpUser models.User
		err = tx.First(
			&tmpUser, "username = ? AND id != ?", input.Username, id,
		).Error
		if err != gorm.ErrRecordNotFound && err != nil {
			s.logger.With(ctx).Error(err)
			return models.User{}, response.GormErrorResponse(err, "خطا در یافتن کاربر")
		}
		if err == nil && tmpUser.ID != user.ID {
			return models.User{}, response.GormErrorResponse(err, "نام کاربری وارد شده تکراری است")
		}
	}

	if input.Phone != "" {
		var tmpUser models.User
		err = tx.First(&tmpUser, "phone = ? AND id != ?", input.Phone, id).Error
		if err != gorm.ErrRecordNotFound && err != nil {
			s.logger.With(ctx).Error(err)
			return models.User{}, response.GormErrorResponse(err, "خطا در یافتن کاربر")
		}
		if err == nil && tmpUser.ID != user.ID {
			return models.User{}, response.GormErrorResponse(err, "شماره تلفن وارد شده تکراری است")
		}
	}

	if input.IDCode != "" {
		var tmpUser models.User
		err = tx.First(&tmpUser, "id_code = ? AND id != ?", input.IDCode, id).Error
		if err != gorm.ErrRecordNotFound && err != nil {
			s.logger.With(ctx).Error(err)
			return models.User{}, response.GormErrorResponse(err, "خطا در یافتن کاربر")
		}
		if err == nil && tmpUser.ID != user.ID {
			return models.User{}, response.GormErrorResponse(err, "کد ملی وارد شده تکراری است")
		}
	}

	user.Bio = input.Bio
	user.Name = input.Name
	user.Phone = input.Phone
	user.Grade = input.Grade
	user.Title = input.Title
	if input.Email != user.Email.String {
		user.Email = dtp.NullString{
			String: input.Email,
			Valid:  input.Email != "",
		}
	}
	user.IDCode = dtp.NullString{
		String: input.IDCode,
		Valid:  input.IDCode != "",
	}
	user.LastName = input.LastName
	user.Nickname = input.Nickname
	user.Username = dtp.NullString{
		String: input.Username,
		Valid:  input.Username != "",
	}
	user.AvatarUrl = input.AvatarUrl
	user.CountryCode = input.CountryCode
	user.City = input.City
	user.Gender = input.Gender
	user.DateOfBirth = dtp.NullTime{
		Time:  input.DateOfBirth,
		Valid: input.DateOfBirth != time.Time{},
	}
	user.SuspendedAt = dtp.NullTime{
		Valid: input.SuspendedAt,
		Time:  time.Now(),
	}
	user.MadeOfficialAt = dtp.NullTime{
		Valid: input.MadeOfficialAt,
		Time:  time.Now(),
	}
	user.EmailVerifiedAt = dtp.NullTime{
		Valid: input.EmailVerifiedAt,
		Time:  time.Now(),
	}
	user.PhoneVerifiedAt = dtp.NullTime{
		Valid: input.PhoneVerifiedAt,
		Time:  time.Now(),
	}
	user.MadeProfilePublicAt = dtp.NullTime{
		Valid: input.MadeProfilePublicAt,
		Time:  time.Now(),
	}
	user.ProfileCompletedAt = profileCompletedAt

	var roles []*models.Role
	if len(input.Roles) != 0 {
		err = s.db.Find(&roles, input.Roles).Error
		if err != nil {
			s.logger.With(ctx).Error(err)
			return user, response.GormErrorResponse(err, "خطا در یافتن نقش ها")
		}
	}

	Id := policy.ExtractIdClaim(ctx)
	userID := fmt.Sprintf("%v", user.ID)
	// if user is admin and wants to delete admin role, return error
	// first find admin role in roles
	var adminRole *models.Role
	for _, r := range roles {
		if r.Title == role.SuperAdmin {
			adminRole = r
			break
		}
	}

	for _, r := range user.Roles {
		if r.Title == role.SuperAdmin && adminRole == nil && userID == Id {
			err := response.ErrorForbidden("شما اجازه حذف نقش مدیر را ندارید")
			return models.User{}, err
		}
	}
	err = tx.Model(&user).Association("Roles").Replace(roles)
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.User{}, response.GormErrorResponse(err, "خطا در ذخیره کاربر")
	}

	invite := &models.Invite{
		Code:  dtp.NullString{String: input.Invite.Code, Valid: input.Invite.Code != ""},
		Limit: input.Invite.Limit,
	}

	err = tx.
		Where("user_id", id).
		Delete(&models.Invite{}).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.User{}, response.GormErrorResponse(err, "خطا در ذخیره کاربر")
	}

	// if there is any invite code received, create a new invite
	if invite.Code.Valid {
		err = tx.Model(&user).Association("Invite").
			Replace(invite)
		if err != nil {
			return models.User{}, response.GormErrorResponse(err, "خطا در ذخیره کاربر")
		}
	}
	err = tx.Save(&user).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.User{}, response.GormErrorResponse(err, "خطا در ذخیره کاربر")
	}

	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback() // Rollback the transaction if there's any error or panic
		} else {
			err = tx.Commit().Error // Commit the transaction if there are no errors
		}
	}()
	return user, response.ErrorResponse{}
}
