package cmd

import (
	"encoding/json"
	"os"

	"github.com/manifoldco/promptui"
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
		ArkZoneId:    c.String("ark-zone-id"),
		ZoneName:     c.String("zone-name"),
		ZoneCategory: c.String("zone-category"),
		ZoneType:     null.NewString(c.String("zone-type"), c.IsSet("zone-type")),
		Server:       c.String("server"),
		StartTime:    c.Timestamp("start-time"),
		EndTime:      c.Timestamp("end-time"),
	})
	if err != nil {
		return err
	}

	filename := c.String("output")
	writeToFile(filename, rendered)

	prompt := promptui.Prompt{
		Label:     "File written. Continue with file content? (hint: you can edit the file before continuing)",
		IsConfirm: true,
	}

	if _, err := prompt.Run(); err != nil {
		return err
	}

	rendered, err = readFromFile(filename)
	if err != nil {
		return err
	}

	return a.GameDataService.UpdateNewEvent(c.Context, rendered)
}

func readFromFile(filename string) (*gamedata.RenderedObjects, error) {
	log.Info().Msgf("Reading rendered game data back from %s", filename)
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var rendered gamedata.RenderedObjects
	err = json.NewDecoder(f).Decode(&rendered)
	if err != nil {
		return nil, err
	}

	return &rendered, nil
}

func writeToFile(filename string, data interface{}) error {
	log.Info().Msgf("Writing rendered game data to %s", filename)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	// marshal rendered to JSON
	b, err := json.MarshalIndent(data, "", "  ")
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
