package healthcheck

import "github.com/labstack/echo/v4"

// RegisterHandlers registers the handlers that perform healthchecks.
func RegisterHandlers(e *echo.Echo, version string) {
	e.Match([]string{"GET", "HEAD"}, "/healthcheck", healthcheck(version))
}

// healthcheck responds to a healthcheck request.
func healthcheck(version string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(200, "OK V-"+version)
	}
}
