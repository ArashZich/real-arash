package transports

import (
	"net/http"

	"gitag.ir/armogroup/armo/services/reality/services/organization/endpoints"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/labstack/echo/v4"
)

func (r *resource) setupShowroomUrl(ctx echo.Context) error {
	// Parse the request body into SetupShowroomUrlRequest struct
	var input endpoints.SetupShowroomUrlRequest
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(err.Error()))
	}

	// Validate the input
	errors := input.Validate(ctx.Request())
	if errors != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(errors))
	}

	// Call the service method to handle the request
	err := r.service.SetupShowroomUrl(ctx.Request().Context(), input)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}
	return ctx.JSON(http.StatusOK, map[string]string{"message": "ShowroomUrl setup completed"})
}
