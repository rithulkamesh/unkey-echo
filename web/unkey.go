package web

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	unkeygo "github.com/unkeyed/unkey-go"
	"github.com/unkeyed/unkey-go/models/components"
)

func unkeyMiddleware(state *AppState) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.String(http.StatusUnauthorized, "Unauthorized")
			}

			request := components.V1KeysVerifyKeyRequest{
				APIID: unkeygo.String(os.Getenv("UNKEY_API_ID")),
				Key:   authHeader,
			}
			ctx := context.Background()
			res, err := state.UnkeyClient.Keys.VerifyKey(ctx, request)
			if err != nil {
				log.Printf("Error verifying key: %v", err)
				return c.String(http.StatusInternalServerError, "Internal Server Error")
			}

			if res.V1KeysVerifyKeyResponse != nil && res.V1KeysVerifyKeyResponse.Valid {
				return next(c)
			}

			return c.String(http.StatusUnauthorized, "Unauthorized: Invalid API key")
		}
	}
}
