package transports

import (
	"fmt"
	"net/http"

	"gitag.ir/armogroup/armo/services/reality/services/notification/endpoints"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/labstack/echo/v4"
)

func (r resource) updateNotification(ctx echo.Context) error {
	id := ctx.Param("id")
	var nid uint
	if _, err := fmt.Sscan(id, &nid); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(nil, "invalid notification id"))
	}

	var input endpoints.UpdateNotificationRequest
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(nil, "invalid input"))
	}

	notification, err := r.service.Update(ctx.Request().Context(), nid, input)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}

	return ctx.JSON(http.StatusOK, response.Success(notification))
}
