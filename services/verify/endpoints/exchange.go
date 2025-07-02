package endpoints

import (
	"context"
	"net/http"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
)

type ExchangeRequest struct {
	Code string `json:"code,omitempty"`
}

func (r *ExchangeRequest) Validate(req *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"code": govalidity.New("code").Required(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"code": "کد",
		},
	)

	errr := govalidity.ValidateBody(req, schema, r)
	if len(errr) > 0 {
		return govalidity.DumpErrors(errr)
	}

	return nil
}

func (s *service) Exchange(ctx context.Context, input ExchangeRequest) (string, response.ErrorResponse) {
	var verification models.Verification
	err := s.db.WithContext(ctx).Find(&verification, "code", input.Code).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return "", response.GormErrorResponse(err, "خطایی در یافتن کد رخ داده است")
	}
	if verification.Expired() {
		return "", response.GormErrorResponse(err, "کد منقضی شده است.")
	}

	if err = s.db.WithContext(ctx).Save(&verification).Error; err != nil {
		s.logger.With(ctx).Error(err)
		return "", response.GormErrorResponse(err, "خطایی در ذخیره کد رخ داده است")
	}

	sessionCode := verification.SessionCode
	return sessionCode, response.ErrorResponse{}
}
