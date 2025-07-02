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

type UpdatePlanRequest struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	Price           int    `json:"price"`
	DiscountedPrice int    `json:"discounted_price"`
	Categories      []int  `json:"categories"`
	DayLength       int    `json:"day_length"`
	ProductLimit    int    `json:"product_limit"`
	StorageLimitMB  int    `json:"storage_limit_mb"`
	IconUrl         string `json:"icon_url"`
}

func (c *UpdatePlanRequest) Validate(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"title":            govalidity.New("title").Required().MinMaxLength(2, 200),
		"description":      govalidity.New("description").Required().MinMaxLength(2, 200),
		"price":            govalidity.New("price").Required(),
		"discounted_price": govalidity.New("discounted_price").Optional(),
		"categories":       govalidity.New("categories").Required(),
		"day_length":       govalidity.New("day_length").Required(),
		"product_limit":    govalidity.New("product_limit").Required(),
		"storage_limit_mb": govalidity.New("storage_limit_mb").Required(),
		"icon_url":         govalidity.New("icon_url").Required(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"title":            "عنوان",
			"description":      "توضیحات",
			"price":            "قیمت تومان",
			"discounted_price": "قیمت تخفیف خورده",
			"categories":       "دسته بندی ها",
			"day_length":       "زمان دوره",
			"product_limit":    "محدودیت محصولات",
			"storage_limit_mb": "محدودیت فضای ذخیره سازی",
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

func (s *service) Update(ctx context.Context, id string, input UpdatePlanRequest) (
	models.Plan, response.ErrorResponse,
) {
	var plan models.Plan

	if !policy.CanUpdatePlan(ctx) {
		s.logger.With(ctx).Error("شما دسترسی ویرایش طرح ندارید")
		return models.Plan{}, response.ErrorForbidden(nil, "شما دسترسی ویرایش طرح ندارید")
	}

	err := s.db.WithContext(ctx).First(&plan, "id", id).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Plan{}, response.GormErrorResponse(err, "خطایی در یافتن طرح رخ داده است")
	}

	var categories []*models.Category
	err = s.db.WithContext(ctx).Find(&categories, input.Categories).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Plan{}, response.GormErrorResponse(err, "دسته مورد نظر یافت نشد")
	}

	plan.Title = input.Title
	plan.Description = input.Description
	plan.Categories = categories
	plan.Price = input.Price
	plan.DiscountedPrice = input.DiscountedPrice
	plan.DayLength = input.DayLength
	plan.ProductLimit = input.ProductLimit
	plan.StorageLimitMB = input.StorageLimitMB
	plan.IconUrl = input.IconUrl

	err = s.db.WithContext(ctx).Save(&plan).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Plan{}, response.GormErrorResponse(err, "خطایی در بروزرسانی طرح رخ داده است")
	}
	return plan, response.ErrorResponse{}
}
