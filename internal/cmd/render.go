package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/manifoldco/promptui"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"gopkg.in/guregu/null.v3"

	"github.com/penguin-statistics/soracli/internal/models/gamedata"
	"github.com/penguin-statistics/soracli/internal/pkg/filepath"
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
	rendered, err := a.GameDataService.RenderNewEvent(c.Context, c.String("sourceUrl"), &gamedata.NewEventBasicInfo{
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

	filename := filepath.UnderDataDir(fmt.Sprintf("rendered-%s.json", rendered.Zone.ArkZoneID))
	writeToFile(filename, rendered)

	// open rendered file in editor
	editor := c.String("editor")
	log.Info().Msgf("opening rendered file in editor: %s", editor)
	if err := exec.Command(editor, filename).Run(); err != nil {
		log.Error().Err(err).Msg("failed to open rendered file in editor. you may want to open it manually")
	}

	prompt := promptui.Prompt{
		Label:     "Continue with edited file content? (you can edit the file before continuing, soracli will read the file back before continuing)",
		IsConfirm: true,
	}

	if _, err := prompt.Run(); err != nil {
		return err
	}

	rendered, err = readFromFile(filename)
	if err != nil {
		return err
	}

	if err = a.GameDataService.UpdateNewEvent(c.Context, rendered); err != nil {
		return err
	}

	log.Info().Msg("successfully updated game data")
	return nil
}

func readFromFile(filename string) (*gamedata.RenderedObjects, error) {
	log.Info().Msgf("reading rendered game data back from %s", filename)
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
	log.Info().Msgf("writing rendered game data to %s", filename)
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
