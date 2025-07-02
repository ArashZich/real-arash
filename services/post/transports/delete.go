// services/post/transports/delete.go
package transports

import (
	"net/http"

	"github.com/ARmo-BigBang/kit/response"
	"github.com/ARmo-BigBang/kit/restypes"
	"github.com/labstack/echo/v4"
)

type deleteRequest struct {
	IDs []uint `json:"ids"`
}

func (r *resource) delete(ctx echo.Context) error {
	var req deleteRequest
	if err := ctx.Bind(&req); err != nil {
		r.logger.With(ctx.Request().Context()).Info(err)
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest("Invalid request format"))
	}

	if len(req.IDs) == 0 {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest("No IDs provided"))
	}

	deletedPosts, err := r.service.Delete(ctx.Request().Context(), req.IDs)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}

	// Convert []uint to []int
	deletedPostIDs := make([]int, len(deletedPosts))
	for i, v := range deletedPosts {
		deletedPostIDs[i] = int(v)
	}

	result := restypes.DeleteResponse{
		IDs: deletedPostIDs,
	}

	return ctx.JSON(http.StatusOK, response.Success(result))
}
