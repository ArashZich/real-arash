package transports

import (
	"net/http"

	"gitag.ir/armogroup/armo/services/reality/services/user/endpoints"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/labstack/echo/v4"
)

func (r resource) addAffiliateCodes(ctx echo.Context) error {
	id := ctx.Param("id")

	var input endpoints.AddAffiliateCodeRequest
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest("Invalid request body"))
	}

	user, err := r.service.AddAffiliateCode(ctx.Request().Context(), id, input)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}

	return ctx.JSON(http.StatusOK, response.Success(user))
}
