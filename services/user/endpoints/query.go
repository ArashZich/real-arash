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

func (s *service) Query(
	ctx context.Context, offset, limit int, orderBy, order, query, suspendedAt, isOfficial, isProfileComplete, hasOrganization, hasPackages string, affiliateCodes []string,
) ([]models.User, response.ErrorResponse) {
	if !policy.CanQueryUsers(ctx) {
		return nil, response.ErrorForbidden("Access denied for querying users")
	}

	var users []models.User
	query = fmt.Sprintf("%%%s%%", strings.ToLower(query))

	tx := s.db.WithContext(ctx).
		Order(fmt.Sprintf("%s %s", orderBy, order)).
		Offset(offset).Limit(limit).
		Preload("Roles").
		Preload("Organizations").
		Preload("Organizations.Category").
		Preload("Organizations.Packages").
		Preload("Organizations.Packages.Plan").
		Preload("Organizations.Packages.Category")

	tx.Where(
		"LOWER(name) LIKE ? OR LOWER(last_name) LIKE ? OR LOWER(phone) LIKE ? OR LOWER(email) LIKE ? OR LOWER(id_code) LIKE ? OR LOWER(username) LIKE ?", query, query, query, query, query, query,
	)

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

	// Apply affiliateCodes condition based on the CanUserGod policy
	if policy.CanUserGod(ctx) {
		if len(affiliateCodes) > 0 {
			affiliateCondition := "ARRAY['" + strings.Join(affiliateCodes, "', '") + "'] && STRING_TO_ARRAY(users.affiliate_codes, ', ')"
			tx = tx.Where(affiliateCondition)
		}
	} else if len(affiliateCodes) > 0 {
		affiliateCondition := "ARRAY['" + strings.Join(affiliateCodes, "', '") + "'] && STRING_TO_ARRAY(users.affiliate_codes, ', ')"
		tx = tx.Where(affiliateCondition)
	} else {
		// If not a "god" user and no affiliate codes are provided, access is restricted
		return nil, response.ErrorForbidden("Access denied due to insufficient privileges or missing affiliate codes")
	}

	err := tx.Find(&users).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return nil, response.GormErrorResponse(err, "Error finding users")
	}

	return users, response.ErrorResponse{}
}
