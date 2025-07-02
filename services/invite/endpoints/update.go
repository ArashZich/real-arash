package endpoints

import (
	"context"
	"net/http"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"
	"github.com/ARmo-BigBang/kit/dtp"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
)

type UpdateInviteRequest struct {
	Code  string `json:"code,omitempty"`
	Limit int    `json:"limit,omitempty"`
}

func (r *UpdateInviteRequest) Validate(req *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"code":  govalidity.New("code").Required(),
		"limit": govalidity.New("limit").Required(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"code":  "کد",
			"limit": "محدودیت",
		},
	)

	errr := govalidity.ValidateBody(req, schema, r)

	if len(errr) > 0 {
		dumpedErrors := govalidity.DumpErrors(errr)
		return dumpedErrors
	}

	return nil
}

func (s *service) Update(ctx context.Context, id string, input UpdateInviteRequest) (models.Invite, response.ErrorResponse) {
	var invite models.Invite
	if !policy.CanUpdateInvite(ctx) {
		s.logger.With(ctx).Error("شما اجازه ویرایش کد معرف ندارید.")
		return invite, response.ErrorForbidden("شما اجازه ویرایش کد معرف ندارید.")
	}
	err := s.db.WithContext(ctx).First(&invite, "id", id).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return invite, response.GormErrorResponse(err, "کد معرف یافت نشد.")
	}
	invite.Limit = input.Limit
	invite.Code = dtp.NullString{String: input.Code, Valid: input.Code != ""}

	err = s.db.WithContext(ctx).Save(&invite).Error

	return invite, response.ErrorResponse{}
}
