package endpoints

import (
	"context"
	"net/http"
	"time"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"
	"github.com/ARmo-BigBang/kit/dtp"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
)

type UpdatePackageRequest struct {
	ProductLimit   int    `json:"product_limit"`
	StorageLimitMB int    `json:"storage_limit_mb"`
	ExpiredAt      string `json:"expired_at"`
}

func (c *UpdatePackageRequest) Validate(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"product_limit":    govalidity.New("product_limit").Optional(),
		"storage_limit_mb": govalidity.New("storage_limit_mb").Optional(),
		"expired_at":       govalidity.New("expired_at").Optional(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"product_limit":    "محدودیت محصولات",
			"storage_limit_mb": "محدودیت فضای ذخیره سازی",
			"expired_at":       "تاریخ انقضا",
			"icon_url":         "آدرس آیکون",
		},
	)

	errr := govalidity.ValidateBody(r, schema, c)

	if len(errr) > 0 {
		dumpedErrors := govalidity.DumpErrors(errr)
		return dumpedErrors
	}

	return nil
}

func (s *service) Update(ctx context.Context, id string, input UpdatePackageRequest) (
	models.Package, response.ErrorResponse,
) {
	var pkg models.Package

	if !policy.CanUpdatePackage(ctx) {
		s.logger.With(ctx).Error("شما دسترسی ویرایش سازمان ندارید")
		return models.Package{}, response.ErrorForbidden(nil, "شما دسترسی ویرایش بسته ندارید")
	}

	err := s.db.WithContext(ctx).First(&pkg, "id", id).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Package{}, response.GormErrorResponse(err, "خطایی در یافتن بسته رخ داده است")
	}

	t, err := time.Parse(time.RFC3339, input.ExpiredAt)
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Package{}, response.ErrorInternalServerError("تاریخ انقضا معتبر نیست")
	}

	pkg.ProductLimit = input.ProductLimit
	pkg.StorageLimitMB = input.StorageLimitMB
	pkg.ExpiredAt = dtp.NullTime{
		Time:  t,
		Valid: input.ExpiredAt != "",
	}

	err = s.db.WithContext(ctx).Save(&pkg).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Package{}, response.GormErrorResponse(err, "خطایی در بروزرسانی بسته رخ داده است")
	}
	return pkg, response.ErrorResponse{}
}
