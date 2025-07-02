package transports

import (
	"net/http"

	"github.com/ARmo-BigBang/kit/response"
	"github.com/ARmo-BigBang/kit/restypes"
	"github.com/labstack/echo/v4"
)

type DeleteNotificationsRequest struct {
	IDs []int `json:"ids"`
}

func (r resource) deleteNotification(ctx echo.Context) error {
	var input DeleteNotificationsRequest
	if err := ctx.Bind(&input); err != nil {
		r.logger.With(ctx.Request().Context()).Error(err)
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(nil, "Invalid input"))
	}

	if len(input.IDs) == 0 {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(nil, "IDs are required"))
	}

	deletedNotifications, err := r.service.Delete(ctx.Request().Context(), input.IDs)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}

	result := restypes.DeleteResponse{
		IDs: deletedNotifications,
	}

	return ctx.JSON(http.StatusOK, response.Success(result))
}
