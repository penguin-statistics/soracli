package cmd

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"gopkg.in/guregu/null.v3"

	"github.com/penguin-statistics/soracli/internal/models/gamedata"
	"github.com/penguin-statistics/soracli/internal/services"
)

type CliApp struct {
	GameDataService *services.GameDataService
}

func NewCliApp(gameDataService *services.GameDataService) *CliApp {
	return &CliApp{
		GameDataService: gameDataService,
	}
}

func (a *CliApp) RenderGameData(c *cli.Context) error {
	rendered, err := a.GameDataService.RenderNewEvent(c.Context, &gamedata.NewEventBasicInfo{
		ArkZoneId:    c.String("arkZoneId"),
		ZoneName:     c.String("zoneName"),
		ZoneCategory: c.String("zoneCategory"),
		ZoneType:     null.NewString(c.String("zoneType"), c.IsSet("zoneType")),
		Server:       c.String("server"),
		StartTime:    c.Timestamp("startTime"),
		EndTime:      c.Timestamp("endTime"),
	})
	if err != nil {
		return err
	}
	fileName := c.String("output")
	if fileName == "" {
		fileName = "game-data.json"
	}
	log.Info().Msgf("Writing rendered game data to %s", fileName)
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	// marshal rendered to JSON
	b, err := json.Marshal(rendered)
	if err != nil {
		return err
	}

	// write to file
	_, err = f.Write(b)
	if err != nil {
		return err
	}

	return nil
}
