package transports

import (
	"net/http"

	"gitag.ir/armogroup/armo/services/reality/services/notification/endpoints"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/labstack/echo/v4"
)

func (r resource) createNotification(ctx echo.Context) error {
	var input endpoints.CreateNotificationRequest
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(nil, "invalid input"))
	}

	notifications, err := r.service.Create(ctx.Request().Context(), input)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}

	return ctx.JSON(http.StatusCreated, response.Created(notifications))
}
