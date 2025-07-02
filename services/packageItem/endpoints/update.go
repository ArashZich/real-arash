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

type UpdatePackageItemRequest struct {
	Title string `json:"title"`
	Price int    `json:"price"`
}

func (c *UpdatePackageItemRequest) Validate(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"title": govalidity.New("title").Required().MinMaxLength(2, 200),
		"price": govalidity.New("price").Required(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"title": "عنوان",
			"price": "مبلغ",
		},
	)

	errr := govalidity.ValidateBody(r, schema, c)

	if len(errr) > 0 {
		dumpedErrors := govalidity.DumpErrors(errr)
		return dumpedErrors
	}

	return nil
}

func (s *packageItem) Update(ctx context.Context, id string, input UpdatePackageItemRequest) (
	models.PackageItem, response.ErrorResponse,
) {
	var packageItem models.PackageItem

	if !policy.CanUpdatePackageItem(ctx) {
		s.logger.With(ctx).Error("شما دسترسی ویرایش آیتم را ندارید")
		return models.PackageItem{}, response.ErrorForbidden(nil, "شما دسترسی ویرایش آیتم را ندارید")
	}

	err := s.db.WithContext(ctx).First(&packageItem, "id", id).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.PackageItem{}, response.GormErrorResponse(err, "خطایی در یافتن آیتم رخ داده است")
	}

	packageItem.Title = input.Title
	packageItem.Price = input.Price

	err = s.db.WithContext(ctx).Save(&packageItem).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.PackageItem{}, response.GormErrorResponse(err, "خطایی در بروزرسانی آیتم رخ داده است")
	}
	return packageItem, response.ErrorResponse{}
}
