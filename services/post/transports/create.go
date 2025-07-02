package transports

import (
	"fmt"
	"net/http"

	"gitag.ir/armogroup/armo/services/reality/services/post/endpoints"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/labstack/echo/v4"
)

func (r *resource) create(ctx echo.Context) error {
	var input endpoints.CreatePostRequest

	errs := ctx.Bind(&input)
	if errs != nil {
		fmt.Println("Error binding request:", errs)
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(errs.Error()))
	}

	fmt.Println("Received input:", input)

	errors := input.Validate(ctx.Request())
	if len(errors) > 0 {
		fmt.Println("Validation errors:", errors)
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(errors))
	}

	fmt.Println("Validation passed. Creating post...")

	post, err := r.service.Create(ctx.Request().Context(), input)
	if err.StatusCode != 0 {
		fmt.Println("Error creating post:", err)
		return ctx.JSON(err.StatusCode, err)
	}

	fmt.Println("Post created successfully:", post)
	return ctx.JSON(http.StatusCreated, response.Created(post))
}
