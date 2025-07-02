// services/organization/transports/getPublicByUid.go

package transports

import (
	"net/http"

	"github.com/ARmo-BigBang/kit/response"
	"github.com/labstack/echo/v4"
)

func (r *resource) getPublicByUid(ctx echo.Context) error {
	uid := ctx.Param("uid")
	if uid == "" {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest("شناسه سازمان الزامی است"))
	}

	organization, err := r.service.GetPublicByUID(ctx.Request().Context(), uid)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}

	return ctx.JSON(http.StatusOK, response.Success(organization))
}
