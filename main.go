package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	unkeygo "github.com/unkeyed/unkey-go"
	"github.com/unkeyed/unkey-go/models/components"
)

type AppState struct {
	UnkeyClient *unkeygo.Unkey
}

func InitializeAppState() (*AppState, error) {
	unkeyClient := unkeygo.New(
		unkeygo.WithSecurity(os.Getenv("UNKEY_ROOT_KEY")),
	)

	return &AppState{
		UnkeyClient: unkeyClient,
	}, nil
}

func unkeyMiddleware(state *AppState) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.String(http.StatusUnauthorized, "Unauthorized: No API key provided")
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

func main() {
	appState, err := InitializeAppState()
	if err != nil {
		log.Fatalf("Failed to initialize app state: %v", err)
	}

	e := echo.New()
	e.Use(unkeyMiddleware(appState))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
