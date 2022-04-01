package appentry

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"

	"github.com/penguin-statistics/soracli/internal/cmd"
	"github.com/penguin-statistics/soracli/internal/models/cache"
	"github.com/penguin-statistics/soracli/internal/pkg/client"
	"github.com/penguin-statistics/soracli/internal/pkg/filepath"
	"github.com/penguin-statistics/soracli/internal/services"
)

func CliApp(c *cli.Context) (*cmd.CliApp, error) {
	logFileName := fmt.Sprintf("logs/soracli-%s.log", time.Now().Format("20060102-150405"))
	logFile, err := os.OpenFile(filepath.UnderDataDir(logFileName), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	logWriters := zerolog.MultiLevelWriter(
		zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"},
		zerolog.ConsoleWriter{Out: logFile, NoColor: true, TimeFormat: "2006-01-02 15:04:05 Z07:00"},
	)

	level := zerolog.DebugLevel
	if c.Bool("verbose") {
		level = zerolog.TraceLevel
	}
	log.Logger = zerolog.New(logWriters).With().Timestamp().Logger().Level(level)

	var app *cmd.CliApp

	opts := []fx.Option{
		fx.Supply(client.NewHTTPFromCliContext(c)),
		fx.Provide(services.NewItemService),
		fx.Provide(services.NewGameDataService),
		fx.Provide(cmd.NewCliApp),
		fx.Invoke(cache.Initialize),
		fx.Populate(&app),
	}

	err = fx.New(opts...).Start(c.Context)
	if err != nil {
		return nil, err
	}

	return app, nil
}
