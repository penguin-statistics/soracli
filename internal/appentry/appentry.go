package appentry

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"

	"github.com/penguin-statistics/soracli/internal/cmd"
	"github.com/penguin-statistics/soracli/internal/models/cache"
	"github.com/penguin-statistics/soracli/internal/pkg/client"
	"github.com/penguin-statistics/soracli/internal/services"
)

func CliApp(c *cli.Context) (*cmd.CliApp, error) {
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()

	var app *cmd.CliApp

	opts := []fx.Option{
		fx.Supply(client.NewHTTPFromCliContext(c)),
		fx.Provide(services.NewItemService),
		fx.Provide(services.NewGameDataService),
		fx.Provide(cmd.NewCliApp),
		fx.Invoke(cache.Initialize),
		fx.Populate(&app),
	}

	err := fx.New(opts...).Start(c.Context)
	if err != nil {
		return nil, err
	}

	return app, nil
}
