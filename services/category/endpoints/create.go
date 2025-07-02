package endpoints

import (
	"context"
	"net/http"
	"strconv"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"
	"github.com/ARmo-BigBang/kit/dtp"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
)

type CreateCategoryRequest struct {
	ParentID         int    `json:"parent_id"`
	AcceptedFileType string `json:"accepted_file_type"`
	Title            string `json:"title"`
	IconUrl          string `json:"icon_url"`
	ARPlacement      string `json:"ar_placement"`
	URL              string `json:"url"`
}

func (c *CreateCategoryRequest) Validate(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"parent_id":          govalidity.New("parent_id").Optional(),
		"accepted_file_type": govalidity.New("accepted_file_type").Required(),
		"title":              govalidity.New("title").Required().MinMaxLength(2, 200),
		"icon_url":           govalidity.New("icon_url").Required(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"parent_id":          "دسته بندی مرتبط",
			"accepted_file_type": "نوع فایل های مجاز",
			"title":              "عنوان",
			"icon_url":           "آیکون",
		},
	)

	errr := govalidity.ValidateBody(r, schema, c)

	if len(errr) > 0 {
		dumpedErrors := govalidity.DumpErrors(errr)
		return dumpedErrors
	}

	return nil
}

func (s *service) Create(ctx context.Context, input CreateCategoryRequest) (models.Category, response.ErrorResponse) {

	if !policy.CanCreateCategory(ctx) {
		s.logger.With(ctx).Error("شما دسترسی ایجاد دسته بندی ندارید")
		return models.Category{}, response.ErrorForbidden("شما دسترسی ایجاد دسته بندی ندارید")
	}

	Id := policy.ExtractIdClaim(ctx)
	id, _ := strconv.Atoi(Id)

	category := models.Category{
		ParentID: dtp.NullInt64{
			Int64: int64(input.ParentID),
			Valid: input.ParentID != 0,
		},
		UserID:           id,
		AcceptedFileType: input.AcceptedFileType,
		IconUrl:          input.IconUrl,
		Title:            input.Title,
		ARPlacement: dtp.NullString{
			String: input.ARPlacement,
			Valid:  input.ARPlacement != "",
		},
		URL: dtp.NullString{
			String: input.URL,
			Valid:  input.URL != "",
		},
	}

	err := s.db.WithContext(ctx).Create(&category).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return category, response.GormErrorResponse(err, "خطایی در ایجاد دسته رخ داد")
	}
	return category, response.ErrorResponse{}
}
