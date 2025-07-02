package endpoints

import (
	"context"
	"net/http"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
)

type CheckEmailExistsRequest struct {
	Email string `json:"email"`
}

func (r *CheckEmailExistsRequest) Validate(req *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"email": govalidity.New("email").Required(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"email": "ایمیل",
		},
	)

	errr := govalidity.ValidateBody(req, schema, r)
	if len(errr) > 0 {
		return govalidity.DumpErrors(errr)
	}

	return nil
}

func (s *service) CheckEmailExists(ctx context.Context, input CheckEmailExistsRequest) (bool, response.ErrorResponse) {
	var user models.User
	var count int64
	err := s.db.WithContext(ctx).First(&user, "email", input.Email).Count(&count).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return false, response.ErrorBadRequest(err, "خطایی در یافتن ایمیل رخ داده است")
	}

	exists := count > 0
	return exists, response.ErrorResponse{}
}
