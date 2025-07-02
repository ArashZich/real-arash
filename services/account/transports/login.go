package transports

import (
	"fmt"
	"net/http"

	"gitag.ir/armogroup/armo/services/reality/services/account/endpoints"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/labstack/echo/v4"
)

func (r resource) login(ctx echo.Context) error {
	var input = &endpoints.LoginRequest{}
	errors := input.Validate(ctx.Request())

	fmt.Println(errors)

	if errors != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(errors))
	}
	loginResponse, err := r.service.Login(ctx.Request().Context(), *input)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}

	return ctx.JSON(http.StatusOK, response.Success(loginResponse))
}
