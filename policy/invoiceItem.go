package policy

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/services/role"
)

func CanAccessInvoiceItem(ctx context.Context) bool {
	return CanIssueInvoice(ctx)
}

func CanCreateInvoiceItem(ctx context.Context) bool {
	roles := ExtractRolesClaim(ctx)
	for _, r := range roles {
		if r == role.SuperAdmin {
			return true
		}
	}
	return false
}

func CanUpdateInvoiceItem(ctx context.Context) bool {
	return CanCreateInvoiceItem(ctx)
}

func CanDeleteInvoiceItem(ctx context.Context) bool {
	return CanCreateInvoiceItem(ctx)
}
