// endpoints/update.go
package endpoints

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"
	"github.com/ARmo-BigBang/kit/response"
)

type UpdatePostRequest struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	Content         string `json:"content"`
	CoverURL        string `json:"cover_url"`
	Tags            string `json:"tags"`
	MetaKeywords    string `json:"meta_keywords"`
	MetaTitle       string `json:"meta_title"`
	MetaDescription string `json:"meta_description"`
	Published       bool   `json:"published"`
}

func (c *UpdatePostRequest) Validate(r *http.Request) map[string][]string {
	errors := make(map[string][]string)
	if c.Title == "" {
		errors["title"] = append(errors["title"], "عنوان باید مقدار داشته باشد")
	}
	if len(c.Title) < 2 || len(c.Title) > 200 {
		errors["title"] = append(errors["title"], "عنوان باید حداقل 2 و حداکثر 200 کاراکتر باشد")
	}

	// Check if description is provided and validate its length
	if c.Description != "" {
		if len(c.Description) < 2 || len(c.Description) > 500 {
			errors["description"] = append(errors["description"], "توضیحات باید حداقل 2 و حداکثر 500 کاراکتر باشد")
		}
	}

	if c.Content == "" {
		errors["content"] = append(errors["content"], "محتوا باید مقدار داشته باشد")
	}

	if c.CoverURL == "" {
		errors["cover_url"] = append(errors["cover_url"], "آدرس کاور باید مقدار داشته باشد")
	}

	if c.MetaKeywords == "" {
		errors["meta_keywords"] = append(errors["meta_keywords"], "کلمات کلیدی متا باید مقدار داشته باشد")
	}

	if c.MetaTitle == "" {
		errors["meta_title"] = append(errors["meta_title"], "عنوان متا باید مقدار داشته باشد")
	}

	if c.MetaDescription == "" {
		errors["meta_description"] = append(errors["meta_description"], "توضیحات متا باید مقدار داشته باشد")
	}

	if c.Tags == "" {
		errors["tags"] = append(errors["tags"], "برچسب‌ها باید مقدار داشته باشد")
	}

	return errors
}

func (s *service) Update(ctx context.Context, id string, input UpdatePostRequest) (models.Post, response.ErrorResponse) {
	var post models.Post

	if !policy.CanUpdatePost(ctx) {
		s.logger.With(ctx).Error("شما دسترسی ویرایش پست را ندارید")
		return models.Post{}, response.ErrorForbidden(nil, "شما دسترسی ویرایش پست را ندارید")
	}

	err := s.db.WithContext(ctx).First(&post, "id = ?", id).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Post{}, response.GormErrorResponse(err, "پست مورد نظر یافت نشد")
	}

	Id := policy.ExtractIdClaim(ctx)
	userID, _ := strconv.Atoi(Id)

	var user models.User
	err = s.db.First(&user, userID).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Post{}, response.GormErrorResponse(err, "کاربر مورد نظر یافت نشد")
	}

	err = s.db.WithContext(ctx).Model(&post).Updates(models.Post{
		Title:           input.Title,
		Description:     input.Description,
		Content:         input.Content,
		CoverURL:        input.CoverURL,
		Tags:            input.Tags,
		MetaKeywords:    input.MetaKeywords,
		MetaTitle:       input.MetaTitle,
		MetaDescription: input.MetaDescription,
		Published:       input.Published,
		UserID:          user.ID,
		UserName:        user.Name + " " + user.LastName,
		UserAvatar:      user.AvatarUrl,
		UpdatedAt:       time.Now(),
	}).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Post{}, response.GormErrorResponse(err, "خطا در بروزرسانی پست")
	}
	return post, response.ErrorResponse{}
}
