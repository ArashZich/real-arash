// services/organization/endpoints/getPublicByUid.go

package endpoints

import (
	"context"
	"errors"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationPublicResponse struct {
	Name        string `json:"name"`
	CompanyName string `json:"company_name"`
	Website     string `json:"website"`
	CompanyLogo string `json:"company_logo"`
}

func (s *service) GetPublicByUID(ctx context.Context, uid string) (OrganizationPublicResponse, response.ErrorResponse) {
	organizationUID, err := uuid.Parse(uid)
	if err != nil {
		return OrganizationPublicResponse{}, response.ErrorBadRequest("شناسه سازمان نامعتبر است")
	}

	var result OrganizationPublicResponse
	err = s.db.WithContext(ctx).
		Model(&models.Organization{}).
		Select("name, company_name, website, company_logo").
		Where("organization_uid = ?", organizationUID).
		First(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return OrganizationPublicResponse{}, response.ErrorNotFound("سازمان یافت نشد")
		}
		s.logger.With(ctx).Error(err)
		return OrganizationPublicResponse{}, response.GormErrorResponse(err, "خطایی در یافتن سازمان رخ داده است")
	}

	return result, response.ErrorResponse{}
}
