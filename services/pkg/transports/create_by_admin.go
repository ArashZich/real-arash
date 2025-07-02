// create_by_admin-transports.go
package transports

import (
	"net/http"

	"gitag.ir/armogroup/armo/services/reality/services/pkg/endpoints"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/labstack/echo/v4"
)

func (r *resource) createByAdmin(ctx echo.Context) error {
	var input = &endpoints.CreatePackageByAdminRequest{}

	errors := input.Validate(ctx.Request())
	if errors != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(errors))
	}

	pkg, err := r.service.CreateByAdmin(ctx.Request().Context(), *input)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}
	return ctx.JSON(http.StatusCreated, response.Created(pkg))
}
