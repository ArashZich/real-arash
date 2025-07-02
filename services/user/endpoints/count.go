package endpoints

import (
	"context"
	"fmt"
	"strings"

	"gitag.ir/armogroup/armo/services/reality/database"
	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"
	"github.com/ARmo-BigBang/kit/response"
)

func (s *service) Count(ctx context.Context, query, suspendedAt, isOfficial, isProfileComplete, hasOrganization, hasPackages string, affiliateCodes []string) (int64, response.ErrorResponse) {
	if !policy.CanQueryUsers(ctx) {
		return 0, response.ErrorForbidden("Access denied for counting users")
	}

	var count int64
	query = fmt.Sprintf("%%%s%%", strings.ToLower(query))

	tx := s.db.WithContext(ctx).
		Model(&models.User{}).
		Where("LOWER(name) LIKE ? OR LOWER(last_name) LIKE ? OR LOWER(phone) LIKE ? OR LOWER(email) LIKE ? OR LOWER(id_code) LIKE ? OR LOWER(username) LIKE ?", query, query, query, query, query, query)

	if suspendedAt != "" && suspendedAt != "-1" {
		tx = tx.Where("suspended_at" + database.NullClause(suspendedAt))
	}
	if isOfficial != "" && isOfficial != "-1" {
		tx = tx.Where("made_official_at" + database.NullClause(isOfficial))
	}
	if isProfileComplete != "" && isProfileComplete != "-1" {
		tx = tx.Where("profile_completed_at" + database.NullClause(isProfileComplete))
	}

	// فیلتر کردن کاربرانی که سازمان دارند
	if hasOrganization != "" && hasOrganization != "-1" {
		if hasOrganization == "1" {
			tx = tx.Where("EXISTS (SELECT 1 FROM organizations WHERE organizations.user_id = users.id AND organizations.deleted_at = 0)")
		} else if hasOrganization == "0" {
			tx = tx.Where("NOT EXISTS (SELECT 1 FROM organizations WHERE organizations.user_id = users.id AND organizations.deleted_at = 0)")
		}
	}

	// فیلتر کردن کاربرانی که در سازمان‌های خود پکیج دارند
	if hasPackages != "" && hasPackages != "-1" {
		if hasPackages == "1" {
			tx = tx.Where("EXISTS (SELECT 1 FROM organizations o JOIN packages p ON p.organization_id = o.id WHERE o.user_id = users.id AND o.deleted_at = 0 AND p.deleted_at = 0)")
		} else if hasPackages == "0" {
			tx = tx.Where("NOT EXISTS (SELECT 1 FROM organizations o JOIN packages p ON p.organization_id = o.id WHERE o.user_id = users.id AND o.deleted_at = 0 AND p.deleted_at = 0)")
		}
	}

	if len(affiliateCodes) > 0 && policy.CanUserGod(ctx) {
		affiliateCondition := "ARRAY['" + strings.Join(affiliateCodes, "', '") + "'] && STRING_TO_ARRAY(users.affiliate_codes, ', ')"
		tx = tx.Where(affiliateCondition)
	}

	err := tx.Count(&count).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return 0, response.GormErrorResponse(err, "Error in counting users")
	}

	return count, response.ErrorResponse{}
}
