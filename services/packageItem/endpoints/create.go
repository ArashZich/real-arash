package endpoints

import (
	"context"
	"net/http"
	"strconv"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"

	"github.com/ARmo-BigBang/kit/response"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
)

type CreatePackageItemRequest struct {
	Title string `json:"title"`
	Price int    `json:"price"`
}

func (c *CreatePackageItemRequest) Validate(r *http.Request) govalidity.ValidityResponseErrors {
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

func (s *packageItem) Create(ctx context.Context, input CreatePackageItemRequest) (models.PackageItem, response.ErrorResponse) {

	if !policy.CanCreatePackageItem(ctx) {
		s.logger.With(ctx).Error("شما دسترسی ایجاد آیتم را ندارید")
		return models.PackageItem{}, response.ErrorForbidden("شما دسترسی ایجاد آیتم را ندارید")
	}

	Id := policy.ExtractIdClaim(ctx)
	id, _ := strconv.Atoi(Id)

	packageItem := models.PackageItem{

		UserID: id,
		Title:  input.Title,
		Price:  input.Price,
	}

	err := s.db.WithContext(ctx).Create(&packageItem).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return packageItem, response.GormErrorResponse(err, "خطایی در ایجاد آیتم رخ داد")
	}
	return packageItem, response.ErrorResponse{}
}
