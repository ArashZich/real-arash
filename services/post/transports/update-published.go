package transports

import (
	"net/http"

	"gitag.ir/armogroup/armo/services/reality/services/post/endpoints"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/labstack/echo/v4"
)

func (r *resource) updatePublished(ctx echo.Context) error {
	id := ctx.Param("id")
	var input endpoints.UpdatePublishedRequest

	errs := ctx.Bind(&input)
	if errs != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(errs.Error()))
	}

	errors := input.Validate(ctx.Request())
	if len(errors) > 0 {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(errors))
	}

	post, err := r.service.UpdatePublished(ctx.Request().Context(), id, input)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}

	return ctx.JSON(http.StatusOK, response.Success(post))
}
