package policy

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/services/role"
)

func CanAccessInvite(ctx context.Context) bool {
	return CanGetInvite(ctx)
}

func CanGetInvite(ctx context.Context) bool {
	roles := ExtractRolesClaim(ctx)
	for _, r := range roles {
		if r == role.SuperAdmin {
			return true
		}
	}
	return false
}

func CanQueryInvite(ctx context.Context) bool {
	return CanGetInvite(ctx)
}

func CanCreateInvite(ctx context.Context) bool {
	return CanGetInvite(ctx)
}

func CanUpdateInvite(ctx context.Context) bool {
	return CanGetInvite(ctx)
}

func CanDeleteInvite(ctx context.Context) bool {
	return CanGetInvite(ctx)
}
