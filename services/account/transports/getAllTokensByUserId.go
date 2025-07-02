package transports

import (
	"net/http"

	"github.com/ARmo-BigBang/kit/response"
	"github.com/labstack/echo/v4"
)

func (r resource) getAllTokensByUserId(ctx echo.Context) error {
	userID := ctx.Param("user_id")
	tokens, err := r.service.GetAllTokensByUserId(ctx.Request().Context(), userID)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}
	return ctx.JSON(http.StatusOK, response.Success(tokens))
}
