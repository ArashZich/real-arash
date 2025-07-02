package policy

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/services/role"
)

func CanAccessPost(ctx context.Context) bool {
	return CanCreatePost(ctx)
}

func CanCreatePost(ctx context.Context) bool {
	roles := ExtractRolesClaim(ctx)
	for _, r := range roles {
		if r == role.SuperAdmin {
			return true
		}
	}
	return false
}

func CanUpdatePost(ctx context.Context) bool {
	return CanCreatePost(ctx)
}

func CanDeletePost(ctx context.Context) bool {
	return CanCreatePost(ctx)
}
