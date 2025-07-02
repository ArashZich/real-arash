package endpoints

import (
	"context"
	"net/http"
	"time"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"
	"github.com/ARmo-BigBang/kit/response"
)

type UpdatePublishedRequest struct {
	Published bool `json:"published"`
}

func (c *UpdatePublishedRequest) Validate(r *http.Request) map[string][]string {
	errors := make(map[string][]string)
	if c.Published != true && c.Published != false {
		errors["published"] = append(errors["published"], "وضعیت انتشار باید مشخص شود")
	}
	return errors
}

func (s *service) UpdatePublished(ctx context.Context, id string, input UpdatePublishedRequest) (models.Post, response.ErrorResponse) {
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

	post.Published = input.Published
	if post.Published {
		now := time.Now()
		post.PublishedAt = &now
	} else {
		post.PublishedAt = nil
	}
	post.UpdatedAt = time.Now()

	err = s.db.WithContext(ctx).Save(&post).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Post{}, response.GormErrorResponse(err, "خطا در بروزرسانی پست")
	}
	return post, response.ErrorResponse{}
}
