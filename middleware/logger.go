// middleware/logger.go
package middleware

import (
	"bytes"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var skipPaths = []string{
	"/api/v1/views",
	"/api/v1/views/duration",
	"/health",
	"/metrics",
}

func NewCustomLogger() echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper:       skipLogger,
		Format:        logFormat,
		CustomTagFunc: customTagFunc,
	})
}

func skipLogger(c echo.Context) bool {
	path := c.Path()
	for _, skipPath := range skipPaths {
		if path == skipPath {
			return true
		}
	}
	return false
}

const logFormat = `{
    "timestamp": "${time_rfc3339}",
    "duration": "${latency_human}",
    "status": ${status},
    "method": "${method}",
    "path": "${uri}",
    "ip": "${remote_ip}"
    ${error}
}` + "\n"

func customTagFunc(c echo.Context, buf *bytes.Buffer) (int, error) {
	if err := c.Get("error"); err != nil {
		buf.WriteString(fmt.Sprintf(`,"error":"%v"`, err))
	}
	return 0, nil
}
