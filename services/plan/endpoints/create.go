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

type CreatePlanRequest struct {
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

func (c *CreatePlanRequest) Validate(r *http.Request) govalidity.ValidityResponseErrors {
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
			"discount_price":   "قیمت تخفیف خورده",
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

func (s *service) Create(ctx context.Context, input CreatePlanRequest) (models.Plan, response.ErrorResponse) {

	if !policy.CanCreatePlan(ctx) {
		s.logger.With(ctx).Error("شما دسترسی ایجاد طرح را ندارید")
		return models.Plan{}, response.ErrorForbidden("شما دسترسی ایجاد طرح را ندارید")
	}

	Id := policy.ExtractIdClaim(ctx)
	id, _ := strconv.Atoi(Id)

	var categories []*models.Category
	if len(input.Categories) > 0 {
		err := s.db.Find(&categories, input.Categories).Error
		if err != nil {
			s.logger.With(ctx).Error(err)
			return models.Plan{}, response.GormErrorResponse(err, "ابزار مورد نظر یافت شند")
		}
	}
	plan := models.Plan{

		UserID:          id,
		Title:           input.Title,
		Description:     input.Description,
		Price:           input.Price,
		Categories:      categories,
		DiscountedPrice: input.DiscountedPrice,
		DayLength:       input.DayLength,
		ProductLimit:    input.ProductLimit,
		StorageLimitMB:  input.StorageLimitMB,
		IconUrl:         input.IconUrl,
	}

	err := s.db.WithContext(ctx).Preload("Categories").Create(&plan).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return plan, response.GormErrorResponse(err, "خطایی در ایجاد طرح رخ داد")
	}
	return plan, response.ErrorResponse{}
}
