package policy

import (
	"context"
	"strconv"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/services/role"
)

func CanImpersonate(ctx context.Context) bool {
	roles := ExtractRolesClaim(ctx)
	for _, r := range roles {
		if r == role.SuperAdmin {
			return true
		}
	}
	return false
}

func ownerOfAllTokens(tokens []models.Token, theID int) bool {
	var allBelongToUser bool
	for _, token := range tokens {
		allBelongToUser = token.UserID == theID
		if !allBelongToUser {
			break
		}
	}
	return allBelongToUser
}

func CanGetAllTokensByUserId(ctx context.Context, tokens []models.Token) bool {
	roles := ExtractRolesClaim(ctx)
	Id := ExtractIdClaim(ctx)
	theID, _ := strconv.Atoi(Id)
	for _, r := range roles {
		if r == role.SuperAdmin || ownerOfAllTokens(tokens, theID) {
			return true
		}
	}
	return false
}
