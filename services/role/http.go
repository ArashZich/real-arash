package role

import (
	"net/http"
	"path/filepath"

	"gitag.ir/armogroup/armo/services/reality/services/mid"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(r *echo.Echo, service Service, logger log.Logger, prefix string) {
	res := resource{service, logger}

	g := r.Group(filepath.Join("", prefix))
	g.GET("/roles/list", res.list)

	rg := g.Group("")
	rg.Use(mid.EchoJWTHandler(), mid.BindUserToContext)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) list(ctx echo.Context) error {
	c := ctx.Request().Context()
	roles, er := r.service.Query(c)
	if er.StatusCode != 0 {
		return ctx.JSON(er.StatusCode, er)
	}
	return ctx.JSON(http.StatusOK, response.Success(roles))
}
