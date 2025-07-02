// endpoints/update-views.go
package endpoints

import (
	"context"
	"net/http"
	"time"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/response"
)

type UpdateViewsRequest struct {
	Increment int `json:"increment"`
}

func (c *UpdateViewsRequest) Validate(r *http.Request) map[string][]string {
	errors := make(map[string][]string)
	if c.Increment <= 0 {
		errors["increment"] = append(errors["increment"], "مقدار افزایش بازدید باید بزرگتر از صفر باشد")
	}
	return errors
}

func (s *service) UpdateViews(ctx context.Context, id string, input UpdateViewsRequest) (models.Post, response.ErrorResponse) {
	var post models.Post

	// if !policy.CanUpdatePost(ctx) {
	// 	s.logger.With(ctx).Error("شما دسترسی ویرایش پست را ندارید")
	// 	return models.Post{}, response.ErrorForbidden(nil, "شما دسترسی ویرایش پست را ندارید")
	// }

	err := s.db.WithContext(ctx).First(&post, "id = ?", id).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Post{}, response.GormErrorResponse(err, "پست مورد نظر یافت نشد")
	}

	post.Views += input.Increment
	post.UpdatedAt = time.Now()

	err = s.db.WithContext(ctx).Save(&post).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Post{}, response.GormErrorResponse(err, "خطا در بروزرسانی پست")
	}
	return post, response.ErrorResponse{}
}
