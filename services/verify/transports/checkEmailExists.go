package transports

import (
	"net/http"

	"gitag.ir/armogroup/armo/services/reality/services/verify/endpoints"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/labstack/echo/v4"
)

func (r resource) checkEmailExists(ctx echo.Context) error {
	var input = &endpoints.CheckEmailExistsRequest{}
	errors := input.Validate(ctx.Request())

	if errors != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(errors))
	}

	exists, err := r.service.CheckEmailExists(ctx.Request().Context(), *input)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}

	return ctx.JSON(http.StatusOK, response.Success(exists))
}
