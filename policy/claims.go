package policy

import (
	"context"
	"fmt"

	"github.com/ARmo-BigBang/kit/jwd"
)

func ExtractRolesClaim(ctx context.Context) []string {
	claims, _ := jwd.Claims(ctx)
	roles := claims["Roles"].([]interface{})

	stringRoles := make([]string, len(roles))

	for _, r := range roles {
		rr := fmt.Sprintf("%v", r)
		stringRoles = append(stringRoles, rr)
	}
	return stringRoles
}

func ExtractIdClaim(ctx context.Context) string {
	claims, _ := jwd.Claims(ctx)
	Id := fmt.Sprintf("%v", claims["Id"])
	return Id
}
