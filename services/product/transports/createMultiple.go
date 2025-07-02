// package transports

// import (
// 	"net/http"
// 	"sync"

// 	"gitag.ir/armogroup/armo/services/reality/models"
// 	"gitag.ir/armogroup/armo/services/reality/services/product/endpoints"
// 	"github.com/ARmo-BigBang/kit/response"
// 	"github.com/labstack/echo/v4"
// )

// func (r *resource) createMultiple(ctx echo.Context) error {
// 	var inputs []endpoints.CreateProductRequest

// 	if err := ctx.Bind(&inputs); err != nil {
// 		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest("Invalid request format"))
// 	}

// 	var products []models.Product
// 	var errs []response.ErrorResponse
// 	var wg sync.WaitGroup
// 	var mu sync.Mutex

// 	for _, input := range inputs {
// 		wg.Add(1)
// 		go func(input endpoints.CreateProductRequest) {
// 			defer wg.Done()
// 			product, err := r.service.Create(ctx.Request().Context(), input)
// 			mu.Lock()
// 			defer mu.Unlock()
// 			if err.StatusCode != 0 {
// 				errs = append(errs, err)
// 			} else {
// 				products = append(products, product)
// 			}
// 		}(input)
// 	}

// 	wg.Wait()

// 	if len(errs) > 0 {
// 		return ctx.JSON(http.StatusBadRequest, errs)
// 	}

// 	return ctx.JSON(http.StatusCreated, products)
// }

package transports

import (
	"net/http"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/services/product/endpoints"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/labstack/echo/v4"
)

func (r *resource) createMultiple(ctx echo.Context) error {
	var inputs []endpoints.CreateProductRequest

	if err := ctx.Bind(&inputs); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest("Invalid request format"))
	}

	var products []models.Product
	var errResp response.ErrorResponse

	for _, input := range inputs {
		product, err := r.service.Create(ctx.Request().Context(), input)
		if err.StatusCode != 0 {
			errResp = err
			break
		}
		products = append(products, product)
	}

	if errResp.StatusCode != 0 {
		return ctx.JSON(errResp.StatusCode, errResp)
	}

	return ctx.JSON(http.StatusCreated, products)
}
