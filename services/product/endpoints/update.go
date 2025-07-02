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

type UpdateProductRequest struct {
	Name         string `json:"name"`
	ThumbnailURI string `json:"thumbnail_uri"`
	Disabled     bool   `json:"disabled"`
}

func (c *UpdateProductRequest) Validate(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"name":          govalidity.New("name").Required().MinMaxLength(2, 200),
		"thumbnail_uri": govalidity.New("thumbnail_uri").Required(),
		"disabled":      govalidity.New("disabled"),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"name":          "نام محصول",
			"thumbnail_uri": "آیکون",
			"disabled":      "غیر فعال",
		},
	)

	errr := govalidity.ValidateBody(r, schema, c)

	if len(errr) > 0 {
		dumpedErrors := govalidity.DumpErrors(errr)
		return dumpedErrors
	}

	return nil
}

func (s *service) Update(ctx context.Context, id string, input UpdateProductRequest) (
	models.Product, response.ErrorResponse,
) {
	var product models.Product

	var user models.User
	err := s.db.WithContext(ctx).
		Preload("Invite").
		Preload("Roles").
		Preload("Organizations").
		First(&user, "id", id).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Product{}, response.GormErrorResponse(err, "خطا در یافتن کاربر")
	}

	if !policy.CanUpdateOrganization(ctx, user) {
		s.logger.With(ctx).Error("شما دسترسی ویرایش محصول ندارید")
		return models.Product{}, response.ErrorForbidden(nil, "شما دسترسی ویرایش محصول ندارید")
	}

	err = s.db.WithContext(ctx).First(&product, "id", id).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Product{}, response.GormErrorResponse(err, "خطایی در یافتن محصول رخ داده است")
	}

	product.DisabledAt = dtp.NullTime{
		Time:  time.Now(),
		Valid: input.Disabled,
	}
	product.Name = input.Name
	product.ThumbnailURI = input.ThumbnailURI

	err = s.db.WithContext(ctx).Save(&product).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Product{}, response.GormErrorResponse(err, "خطایی در بروزرسانی محصول رخ داده است")
	}
	return product, response.ErrorResponse{}
}
