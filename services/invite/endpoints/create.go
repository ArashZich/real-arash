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

type CreateInviteRequest struct {
	Code   string `json:"code,omitempty"`
	Limit  int    `json:"limit,omitempty"`
	UserID int    `json:"user_id,omitempty"`
}

func (r *CreateInviteRequest) Validate(req *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"code":    govalidity.New("code").Required(),
		"limit":   govalidity.New("limit").Required(),
		"user_id": govalidity.New("user_id").Required(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"code":    "کد",
			"limit":   "محدودیت",
			"user_id": "شناسه کاربر",
		},
	)

	errr := govalidity.ValidateBody(req, schema, r)

	if len(errr) > 0 {
		dumpedErrors := govalidity.DumpErrors(errr)
		return dumpedErrors
	}

	return nil
}

func (s *service) Create(ctx context.Context, input CreateInviteRequest) (models.Invite, response.ErrorResponse) {
	var invite models.Invite
	if !policy.CanCreateInvite(ctx) {
		s.logger.With(ctx).Error("شما اجازه ایجاد کد معرف ندارید.")
		return invite, response.ErrorForbidden("شما اجازه ایجاد کد معرف ندارید.")
	}

	invite = models.Invite{
		Code:   dtp.NullString{String: input.Code, Valid: input.Code != ""},
		Limit:  input.Limit,
		UserID: input.UserID,
	}

	err := s.db.WithContext(ctx).Create(&invite).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return invite, response.GormErrorResponse(err, "خطایی در ذخیره کد معرف رخ داده است")
	}
	return invite, response.ErrorResponse{}
}
