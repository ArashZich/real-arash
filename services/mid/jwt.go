package mid

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"gitag.ir/armogroup/armo/services/reality/config"
	"github.com/ARmo-BigBang/kit/jwd"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func EchoJWTHandler() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(
		middleware.JWTConfig{
			ContextKey:    "user",
			SigningKey:    []byte(config.AppConfig.AccessTokenSigningKey),
			SigningMethod: jwt.SigningMethodHS256.Name,
			Claims:        jwt.MapClaims{},
			TokenLookup:   "header:Authorization",
			AuthScheme:    "Bearer",
			ParseTokenFunc: func(auth string, c echo.Context) (interface{}, error) {
				keyFunc := func(t *jwt.Token) (interface{}, error) {
					if t.Method.Alg() != "HS256" {
						return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
					}
					return []byte(config.AppConfig.AccessTokenSigningKey), nil
				}

				// claims are of type `jwt.MapClaims` when token is created with `jwt.Parse`
				token, err := jwt.Parse(auth, keyFunc)
				if err != nil {
					return nil, err
				}
				if !token.Valid {
					return nil, errors.New("invalid token")
				}

				// Check token expiration
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					if exp, ook := claims["ExpiresAt"].(float64); ook {
						currentTime := time.Now().Unix()
						if int64(exp) < currentTime {
							return nil, errors.New("token has expired")
						}
					}
				}

				// FIXME: check token in redis or database and if not found return token not found error

				return token, nil
			},
		},
	)
}

func BindUserToContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			log.Fatal("Converting user key in context to jwt.Token failed In BindUserToContext() middleware ")
		}

		ctx = context.WithValue(ctx, jwd.UserContextKey, token)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}
