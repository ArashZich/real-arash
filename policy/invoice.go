package policy

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/services/role"
)

func CanAccessInvoice(ctx context.Context) bool {
	return CanCreatePlan(ctx)
}

func CanIssueInvoice(ctx context.Context) bool {
	roles := ExtractRolesClaim(ctx)
	for _, r := range roles {
		if r == role.SuperAdmin {
			return true
		}
	}
	return false
}

func CanUpdateInvoice(ctx context.Context) bool {
	return CanIssueInvoice(ctx)
}

func CanDeleteInvoice(ctx context.Context) bool {
	return CanIssueInvoice(ctx)
}
