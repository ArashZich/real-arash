package policy

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/services/role"
)

func CanAccessDocument(ctx context.Context) bool {
	return CanCreateDocument(ctx)
}

func CanCreateDocument(ctx context.Context) bool {
	roles := ExtractRolesClaim(ctx)
	for _, r := range roles {
		if r == role.SuperAdmin || r == role.Admin || r == role.User {
			return true
		}
	}
	return false
}

func CanUpdateDocument(ctx context.Context) bool {
	return CanCreateDocument(ctx)
}

func CanDeleteDocument(ctx context.Context) bool {
	return CanCreateDocument(ctx)
}
