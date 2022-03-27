package main

import (
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/penguin-statistics/soracli/cmd"
)

func main() {
	app := &cli.App{
		Name:  "sora",
		Usage: "Penguin Statistics Admin CLI",
		Commands: []*cli.Command{
			{
				Name:    "render",
				Aliases: []string{"r"},
				Usage:   "renders a new game data",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "arkZoneId",
						Aliases:  []string{"zi"},
						Usage:    "arkZoneId",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "zoneName",
						Aliases:  []string{"zn"},
						Usage:    "zoneName",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "zoneCategory",
						Aliases:  []string{"zc"},
						Usage:    "zoneCategory",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "zoneType",
						Aliases: []string{"zt"},
						Usage:   "zoneType",
					},
					&cli.StringFlag{
						Name:     "server",
						Aliases:  []string{"s"},
						Usage:    "server",
						Required: true,
					},
					&cli.TimestampFlag{
						Name:     "startTime",
						Aliases:  []string{"st"},
						Usage:    "startTime",
						Required: true,
						Layout:   time.RFC3339,
					},
					&cli.TimestampFlag{
						Name:     "endTime",
						Aliases:  []string{"et"},
						Usage:    "endTime",
						Required: true,
						Layout:   time.RFC3339,
					},
					&cli.StringFlag{
						Name:     "output",
						Aliases:  []string{"o"},
						Usage:    "output file",
						Required: true,
						Value:    "rendered.json",
					},
				},
				Action: func(c *cli.Context) error {
					return cmd.Render(c)
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "baseUrl",
				Aliases:  []string{"u"},
				Usage:    "base url of the admin api, without trailing slash",
				Required: false,
				Value:    "https://penguin-stats.io/api/admin",
			},
			&cli.StringFlag{
				Name:     "token",
				Usage:    "bearer token for authentication to the admin api; required",
				Required: true,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
