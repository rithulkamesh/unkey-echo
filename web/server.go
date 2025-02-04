package web

import (
	"log"
	"os"

	"github.com/kordlab/marketplace/config"
	"github.com/kordlab/marketplace/data"
	"github.com/labstack/echo/v4"
	unkeygo "github.com/unkeyed/unkey-go"
)

type AppState struct {
	UnkeyClient *unkeygo.Unkey
	MongoDB     *data.MongoDB
	RedisDB     *data.RedisDB
	AuthHandler *AuthHandler
}

func initializeAppState() (*AppState, error) {
	unkeyClient := unkeygo.New(
		unkeygo.WithSecurity(os.Getenv("UNKEY_ROOT_KEY")),
	)

	cfg := config.LoadConfig()
	mongodb, err := data.NewMongoDB(cfg)
	if err != nil {
		return nil, err
	}

	redis, err := data.NewRedisDB(cfg)
	if err != nil {
		return nil, err
	}
	appState := &AppState{
		UnkeyClient: unkeyClient,
		MongoDB:     mongodb,
		RedisDB:     redis,
	}
	appState.AuthHandler = NewAuthHandler(mongodb, redis)
	return appState, nil
}

func Serve() {
	appState, err := initializeAppState()
	if err != nil {
		log.Fatalf("Failed to initialize app state: %v", err)
	}

	e := echo.New()
	registerAuthRoutes(e, appState.AuthHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
