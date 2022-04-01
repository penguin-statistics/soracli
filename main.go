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
		Name:  "soracli",
		Usage: "Penguin Statistics Admin CLI",
		Commands: []*cli.Command{
			{
				Name:    "render",
				Aliases: []string{"r"},
				Usage:   "renders a new game data",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "ark-zone-id",
						Aliases:  []string{"zi"},
						Usage:    "ark zone ID",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "zone-name",
						Aliases:  []string{"zn"},
						Usage:    "zone name",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "zone-category",
						Aliases:  []string{"zc"},
						Usage:    "zone category",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "zone-type",
						Aliases: []string{"zt"},
						Usage:   "zone type",
					},
					&cli.StringFlag{
						Name:     "server",
						Aliases:  []string{"s"},
						Usage:    "server",
						Required: true,
					},
					&cli.TimestampFlag{
						Name:     "start-time",
						Aliases:  []string{"st"},
						Usage:    "zone start time",
						Required: true,
						Layout:   time.RFC3339,
					},
					&cli.TimestampFlag{
						Name:     "end-time",
						Aliases:  []string{"et"},
						Usage:    "zone end time",
						Required: true,
						Layout:   time.RFC3339,
					},
					&cli.StringFlag{
						Name:    "editor",
						Aliases: []string{"e"},
						Usage:   "editor",
						Value:   os.Getenv("EDITOR"),
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
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "verbose output",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
