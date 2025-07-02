package endpoints

import (
	"context"
	"net/http"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
)

type UpdateUserAvatarRequest struct {
	AvatarUrl string `json:"avatar_url,omitempty"`
}

func (req *UpdateUserAvatarRequest) Validate(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"avatar_url": govalidity.New("avatar_url").Required(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"avatar_url": "آدرس آواتار",
		},
	)

	errr := govalidity.ValidateBody(r, schema, req)
	if len(errr) > 0 {
		return govalidity.DumpErrors(errr)
	}

	return nil
}

func (s *service) UpdateAvatar(ctx context.Context, input UpdateUserAvatarRequest) (models.User, response.ErrorResponse) {
	var user models.User
	Id := policy.ExtractIdClaim(ctx)
	err := s.db.WithContext(ctx).First(&user, "id", Id).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return user, response.GormErrorResponse(err, "خطا در یافتن کاربر")
	}

	if !policy.CanUpdateAvatar(ctx, user) {
		s.logger.With(ctx).Error("شما اجازه دسترسی به این کاربر را ندارید")
		return user, response.ErrorForbidden("شما اجازه دسترسی به این کاربر را ندارید")
	}

	user.AvatarUrl = input.AvatarUrl

	err = s.db.WithContext(ctx).Save(&user).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return user, response.GormErrorResponse(err, "خطا در ذخیره کاربر")
	}

	return user, response.ErrorResponse{}
}
