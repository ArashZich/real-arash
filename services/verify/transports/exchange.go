package transports

import (
	"net/http"

	"gitag.ir/armogroup/armo/services/reality/services/verify/endpoints"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/labstack/echo/v4"
)

type ExchangeResponse struct {
	SessionCode string `json:"session_code"`
}

func (r resource) exchange(ctx echo.Context) error {
	var input = &endpoints.ExchangeRequest{}
	errors := input.Validate(ctx.Request())

	if errors != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(errors))
	}

	sessionCode, err := r.service.Exchange(ctx.Request().Context(), *input)
	if err.StatusCode != 0 {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(nil))
	}

	result := &ExchangeResponse{
		SessionCode: sessionCode,
	}

	return ctx.JSON(http.StatusOK, response.Success(result))
}
