package transports

import (
	"net/http"

	"github.com/ARmo-BigBang/kit/response"
	"github.com/labstack/echo/v4"
)

func (r resource) logout(ctx echo.Context) error {
	// get access token from header
	accessToken := ctx.Request().Header.Get("Authorization")
	accessToken = accessToken[len("Bearer "):]
	logoutResponse, err := r.service.Logout(ctx.Request().Context(), accessToken)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}

	return ctx.JSON(http.StatusOK, response.Success(logoutResponse))
}
