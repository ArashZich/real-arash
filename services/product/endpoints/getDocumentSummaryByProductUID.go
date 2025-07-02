// File: endpoints/getDocumentsSummaryByProductUID.go

package endpoints

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/google/uuid"
)

type DocumentSummary struct {
	Title      string    `json:"title"`
	FileURI    string    `json:"file_uri"`
	AssetURI   string    `json:"asset_uri"`
	PreviewURI string    `json:"preview_uri"`
	ShopLink   string    `json:"shop_link"` // فیلد جدید اضافه شده
	ProductUID uuid.UUID `json:"product_uid"`
}

func (s *service) GetDocumentSummaryByProductUID(ctx context.Context, productUID string) (DocumentSummary, response.ErrorResponse) {
	var document models.Document
	err := s.db.WithContext(ctx).Where("product_uid = ?", productUID).First(&document).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return DocumentSummary{}, response.GormErrorResponse(err, "خطا در یافتن داکیومنت")
	}

	summary := DocumentSummary{
		Title:      document.Title,
		FileURI:    document.FileURI,
		AssetURI:   document.AssetURI.String,
		PreviewURI: document.PreviewURI,
		ShopLink:   document.ShopLink.String, // فیلد جدید اضافه شده
		ProductUID: document.ProductUID,
	}

	return summary, response.ErrorResponse{}
}
