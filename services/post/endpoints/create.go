// endpoints/create.go
package endpoints

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"
	"github.com/ARmo-BigBang/kit/response"
)

type CreatePostRequest struct {
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

func (c *CreatePostRequest) Validate(r *http.Request) map[string][]string {
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

	fmt.Println("Validation errors:", errors)
	return errors
}

func (s *service) Create(ctx context.Context, input CreatePostRequest) (models.Post, response.ErrorResponse) {
	Id := policy.ExtractIdClaim(ctx)
	id, _ := strconv.Atoi(Id)

	var user models.User
	err := s.db.First(&user, id).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Post{}, response.GormErrorResponse(err, "کاربر مورد نظر یافت نشد")
	}

	fmt.Println("User found:", user)

	post := models.Post{
		Title:           input.Title,
		Description:     input.Description,
		Content:         input.Content,
		CoverURL:        input.CoverURL,
		Tags:            input.Tags,
		MetaKeywords:    input.MetaKeywords,
		MetaTitle:       input.MetaTitle,
		MetaDescription: input.MetaDescription,
		Published:       input.Published,
		Views:           0, // Set initial value of Views to 0
		UserID:          user.ID,
		UserName:        formatUserName(user),
		UserAvatar:      user.AvatarUrl,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if post.Published {
		now := time.Now()
		post.PublishedAt = &now
	}

	fmt.Println("Creating post:", post)

	err = s.db.WithContext(ctx).Create(&post).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return post, response.GormErrorResponse(err, "خطا در ایجاد پست")
	}

	fmt.Println("Post created in DB:", post)
	return post, response.ErrorResponse{}
}

func formatUserName(user models.User) string {
	if user.Name == "" && user.LastName == "" {
		return ""
	}
	if user.Name == "" {
		return user.LastName
	}
	if user.LastName == "" {
		return user.Name
	}
	return user.Name + " " + user.LastName
}
