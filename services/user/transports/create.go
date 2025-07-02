package transports

import (
	"net/http"

	"gitag.ir/armogroup/armo/services/reality/services/user/endpoints"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/labstack/echo/v4"
)

func (r resource) create(ctx echo.Context) error {
	var input = &endpoints.CreateUserRequest{}
	errors := input.Validate(ctx.Request())

	if errors != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(errors, "اطلاعات ورودی نامعتبر است"))
	}

	user, err := r.service.Create(ctx.Request().Context(), *input)

	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}

	return ctx.JSON(http.StatusCreated, response.Created(user))
}
