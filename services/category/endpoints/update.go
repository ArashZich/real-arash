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

type UpdateCategoryRequest struct {
	ParentID         int    `json:"parent_id"`
	AcceptedFileType string `json:"accepted_file_type"`
	Title            string `json:"title"`
	IconUrl          string `json:"icon_url"`
	ARPlacement      string `json:"ar_placement"`
	URL              string `json:"url"`
}

func (c *UpdateCategoryRequest) Validate(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"parent_id":          govalidity.New("parent_id").Optional(),
		"accepted_file_type": govalidity.New("accepted_file_type").Required(),
		"title":              govalidity.New("title").Required().MinMaxLength(2, 200),
		"icon_url":           govalidity.New("icon_url").Required(),
		"ar_placement":       govalidity.New("ar_placement").Optional(),
		"url":                govalidity.New("url").Optional(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"parent_id":          "دسته بندی مرتبط",
			"accepted_file_type": "نوع فایل های مجاز",
			"title":              "عنوان",
			"icon_url":           "آیکون",
			"ar_placement":       "محل AR",
			"url":                "آدرس URL",
		},
	)

	errr := govalidity.ValidateBody(r, schema, c)

	if len(errr) > 0 {
		dumpedErrors := govalidity.DumpErrors(errr)
		return dumpedErrors
	}

	return nil
}

func (s *service) Update(ctx context.Context, id string, input UpdateCategoryRequest) (
	models.Category, response.ErrorResponse,
) {
	var category models.Category

	if !policy.CanUpdateCategory(ctx) {
		s.logger.With(ctx).Error("شما دسترسی ویرایش دسته بندی ندارید")
		return models.Category{}, response.ErrorForbidden("شما دسترسی ویرایش دسته بندی ندارید")
	}

	err := s.db.WithContext(ctx).First(&category, "id = ?", id).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Category{}, response.GormErrorResponse(err, "خطایی در یافتن دسته رخ داده است")
	}

	category.ParentID = dtp.NullInt64{
		Int64: int64(input.ParentID),
		Valid: input.ParentID != 0,
	}
	category.AcceptedFileType = input.AcceptedFileType
	category.IconUrl = input.IconUrl
	category.Title = input.Title
	category.ARPlacement = dtp.NullString{
		String: input.ARPlacement,
		Valid:  input.ARPlacement != "",
	}
	category.URL = dtp.NullString{
		String: input.URL,
		Valid:  input.URL != "",
	}

	err = s.db.WithContext(ctx).Save(&category).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Category{}, response.GormErrorResponse(err, "خطایی در بروزرسانی دسته رخ داده است")
	}
	return category, response.ErrorResponse{}
}
