package transports

import (
	"net/http"

	"gitag.ir/armogroup/armo/services/reality/services/user/endpoints"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/labstack/echo/v4"
)

func (r resource) updateAvatar(ctx echo.Context) error {
	var input = &endpoints.UpdateUserAvatarRequest{}
	errors := input.Validate(ctx.Request())

	if errors != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(errors))
	}

	user, err := r.service.UpdateAvatar(ctx.Request().Context(), *input)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}
	return ctx.JSON(http.StatusOK, response.Success(user))
}
