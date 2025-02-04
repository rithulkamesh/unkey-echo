package web

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
	unkeygo "github.com/unkeyed/unkey-go"
)

var Port = ":3423"

type AppState struct {
	UnkeyClient *unkeygo.Unkey
}

func initializeAppState() (*AppState, error) {
	unkeyClient := unkeygo.New(
		unkeygo.WithSecurity(os.Getenv("UNKEY_ROOT_KEY")),
	)

	return &AppState{
		UnkeyClient: unkeyClient,
	}, nil
}

func Serve() {
	appState, err := initializeAppState()
	if err != nil {
		log.Fatalf("Failed to initialize app state: %v", err)
	}

	e := echo.New()
	e.Use(unkeyMiddleware(appState))

	registerAuthRoutes(e)

	e.Logger.Fatal(e.Start(Port))
}
