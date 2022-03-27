package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/penguin-statistics/soracli/internal/appentry"
)

func Render(c *cli.Context) error {
	app, err := appentry.CliApp(c)
	if err != nil {
		return err
	}

	return app.RenderGameData(c)
}
