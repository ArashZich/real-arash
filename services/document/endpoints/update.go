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

type UpdateDocumentRequest struct {
	Title    *string `json:"title,omitempty"`     // pointer برای اختیاری بودن
	ShopLink *string `json:"shop_link,omitempty"` // pointer برای اختیاری بودن
}

func (c *UpdateDocumentRequest) Validate(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"title":     govalidity.New("title").Optional().MinMaxLength(2, 200),
		"shop_link": govalidity.New("shop_link").Optional().Url(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"title":     "عنوان",
			"shop_link": "لینک فروشگاه",
		},
	)

	errr := govalidity.ValidateBody(r, schema, c)

	if len(errr) > 0 {
		dumpedErrors := govalidity.DumpErrors(errr)
		return dumpedErrors
	}

	return nil
}

func (s *service) Update(ctx context.Context, id string, input UpdateDocumentRequest) (
	models.Document, response.ErrorResponse,
) {
	var document models.Document

	if !policy.CanUpdateDocument(ctx) {
		s.logger.With(ctx).Error("شما دسترسی ویرایش سند ندارید")
		return models.Document{}, response.ErrorForbidden(nil, "شما دسترسی ویرایش سند ندارید")
	}

	err := s.db.WithContext(ctx).First(&document, "id", id).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Document{}, response.GormErrorResponse(err, "خطایی در یافتن سند رخ داده است")
	}

	// فقط فیلدهایی که ارسال شدن رو آپدیت کن
	if input.Title != nil {
		document.Title = *input.Title
	}

	if input.ShopLink != nil {
		document.ShopLink = dtp.NullString{
			String: *input.ShopLink,
			Valid:  *input.ShopLink != "",
		}
	}

	err = s.db.WithContext(ctx).Save(&document).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Document{}, response.GormErrorResponse(err, "خطایی در بروزرسانی سند رخ داده است")
	}

	// آپدیت title محصول اگر title document تغییر کرد
	if input.Title != nil {
		var product models.Product
		err = s.db.WithContext(ctx).Where("product_uid = ?", document.ProductUID).First(&product).Error
		if err == nil { // اگر محصول پیدا شد
			product.Name = *input.Title
			err = s.db.WithContext(ctx).Save(&product).Error
			if err != nil {
				s.logger.With(ctx).Error("خطا در آپدیت نام محصول: ", err)
				// در صورت خطا در آپدیت محصول، لاگ می‌کنیم ولی خطا برنمی‌گردانیم
				// چون document با موفقیت آپدیت شده
			}
		} else {
			s.logger.With(ctx).Error("خطا در یافتن محصول برای آپدیت نام: ", err)
		}
	}

	return document, response.ErrorResponse{}
}
