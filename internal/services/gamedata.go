package services

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gopkg.in/guregu/null.v3"

	"github.com/penguin-statistics/soracli/internal/consts"
	"github.com/penguin-statistics/soracli/internal/models"
	"github.com/penguin-statistics/soracli/internal/models/gamedata"
	"github.com/penguin-statistics/soracli/internal/pkg/client"
	"github.com/penguin-statistics/soracli/internal/pkg/gdutils"
)

var ErrCannotGetFromRemote = errors.New("cannot get from remote")

type GameDataService struct {
	ItemService *ItemService

	http     *http.Client
	pgclient *client.Penguin
}

func NewGameDataService(itemService *ItemService, client *client.Penguin) *GameDataService {
	return &GameDataService{
		ItemService: itemService,
		http: &http.Client{
			Timeout: time.Second * 10,
		},
		pgclient: client,
	}
}

var dropTypeOrderMapping = map[string]int{
	"REGULAR":          0,
	"SPECIAL":          1,
	"EXTRA":            2,
	"FURNITURE":        3,
	"RECOGNITION_ONLY": 4,
}

func (s *GameDataService) RenderNewEvent(ctx context.Context, info *gamedata.NewEventBasicInfo) (*gamedata.RenderedObjects, error) {
	log.Info().Interface("info", info).Msg("rendering new event")
	zone, err := s.renderNewZone(info)
	if err != nil {
		return nil, err
	}

	timeRange := s.renderNewTimeRange(info)

	activity, err := s.renderNewActivity(info)
	if err != nil {
		return nil, err
	}

	importStages, err := s.fetchLatestStages(ctx, []string{info.ArkZoneId})
	if err != nil {
		return nil, err
	}

	stages := make([]*models.Stage, 0)
	dropInfosMap := make(map[string][]*models.DropInfo)
	for _, gamedataStage := range importStages {
		log.Debug().Interface("stage", gamedataStage).Msg("rendering stage")
		stage, dropInfosForOneStage, err := s.genStageAndDropInfosFromGameData(ctx, info.Server, gamedataStage, 0, timeRange)
		if err != nil {
			return nil, err
		}
		stages = append(stages, stage)
		dropInfosMap[stage.ArkStageID] = dropInfosForOneStage
	}

	return &gamedata.RenderedObjects{
		Zone:         zone,
		Stages:       stages,
		DropInfosMap: dropInfosMap,
		TimeRange:    timeRange,
		Activity:     activity,
	}, nil
}

func (s *GameDataService) UpdateNewEvent(ctx context.Context, renderedObjects *gamedata.RenderedObjects) error {
	log.Info().Interface("renderedObjects", renderedObjects).Msg("updating new event")

	return s.pgclient.PostJSON("/save", renderedObjects)
}

func (s *GameDataService) renderNewZone(info *gamedata.NewEventBasicInfo) (*models.Zone, error) {
	nameMap := make(map[string]string)
	for _, lang := range consts.Languages {
		nameMap[lang] = info.ZoneName
	}
	name, err := json.Marshal(nameMap)
	if err != nil {
		return nil, err
	}

	existenceMap := make(map[string]map[string]any)
	for _, s := range consts.Servers {
		existenceMap[s] = map[string]any{
			"exist": false,
		}
		if s == info.Server {
			existenceMap[s]["exist"] = true
			existenceMap[s]["openTime"] = info.StartTime.UnixMilli()
			fakeEndTime := time.UnixMilli(consts.FakeEndTimeMilli)
			endTime := &fakeEndTime
			if info.EndTime != nil {
				endTime = info.EndTime
			}
			existenceMap[s]["closeTime"] = endTime.UnixMilli()
		}
	}
	existence, err := json.Marshal(existenceMap)
	if err != nil {
		return nil, err
	}

	backgroundStr := consts.ZoneBackgroundPath + info.ArkZoneId + consts.ZoneBackgroundExtension
	background := null.StringFrom(backgroundStr)

	return &models.Zone{
		ArkZoneID:  info.ArkZoneId,
		Index:      0,
		Category:   info.ZoneCategory,
		Type:       info.ZoneType,
		Name:       name,
		Existence:  existence,
		Background: background,
	}, nil
}

func (s *GameDataService) renderNewTimeRange(info *gamedata.NewEventBasicInfo) *models.TimeRange {
	log.Info().Interface("info", info).Msg("rendering time range")
	fakeEndTime := time.UnixMilli(consts.FakeEndTimeMilli)
	endTime := &fakeEndTime
	if info.EndTime != nil {
		endTime = info.EndTime
	}

	name := null.StringFrom(gdutils.GetZonePrefixFromArkZoneID(info.ArkZoneId))
	startTimeInComment := info.StartTime.In(consts.LocMap[info.Server]).Format("2006/1/2 15:04")
	endTimeInComment := "?"
	if info.EndTime != nil {
		endTimeInComment = info.EndTime.In(consts.LocMap[info.Server]).Format("2006/1/2 15:04")
	}
	comment := null.StringFrom(consts.ServerNameMapping[info.Server] + info.ZoneName + " " + startTimeInComment + " - " + endTimeInComment)
	return &models.TimeRange{
		StartTime: info.StartTime,
		EndTime:   endTime,
		Server:    info.Server,
		Name:      name,
		Comment:   comment,
	}
}

func (s *GameDataService) renderNewActivity(info *gamedata.NewEventBasicInfo) (*models.Activity, error) {
	fakeEndTime := time.UnixMilli(consts.FakeEndTimeMilli)
	endTime := &fakeEndTime
	if info.EndTime != nil {
		endTime = info.EndTime
	}

	nameMap := make(map[string]string)
	for _, lang := range consts.Languages {
		nameMap[lang] = info.ZoneName
	}
	name, err := json.Marshal(nameMap)
	if err != nil {
		return nil, err
	}

	existenceMap := make(map[string]map[string]any)
	for _, s := range consts.Servers {
		existenceMap[s] = map[string]any{
			"exist": false,
		}
		if s == info.Server {
			existenceMap[s]["exist"] = true
		}
	}
	existence, err := json.Marshal(existenceMap)
	if err != nil {
		return nil, err
	}

	return &models.Activity{
		StartTime: info.StartTime,
		EndTime:   endTime,
		Name:      name,
		Existence: existence,
	}, nil
}

func (s *GameDataService) fetchLatestStages(ctx context.Context, arkZoneIds []string) ([]*gamedata.Stage, error) {
	u := "https://raw.githubusercontent.com/Kengxxiao/ArknightsGameData/master/zh_CN/gamedata/excel/stage_table.json"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, err
	}

	res, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, ErrCannotGetFromRemote
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	stageTable := gamedata.StageTable{}

	json.Unmarshal([]byte(body), &stageTable)

	importStages := make([]*gamedata.Stage, 0)
	for _, stage := range stageTable.Stages {
		if len(arkZoneIds) > 0 && !linq.From(arkZoneIds).Contains(stage.ZoneID) {
			continue
		}

		switch {
		case gdutils.IsCampaignStage(stage):
		case gdutils.IsGuideStage(stage):
		case gdutils.IsDailyStage(stage):
		case gdutils.IsChallengeModeStage(stage):
		case gdutils.IsTrainingStage(stage):
		case gdutils.IsStoryStage(stage):
		case gdutils.IsNormalModeExStage(stage):
		default:
			importStages = append(importStages, stage)
		}
	}
	linq.From(importStages).
		DistinctByT(func(stage *gamedata.Stage) string { return stage.StageID }).
		SortT(func(a, b *gamedata.Stage) bool { return gdutils.CompareStageCode(a.Code, b.Code) }).
		ToSlice(&importStages)

	log.Debug().Interface("stages", importStages).Msg("fetched latest stages")
	return importStages, nil
}

func (s *GameDataService) genStageAndDropInfosFromGameData(ctx context.Context, server string, gamedataStage *gamedata.Stage, zoneId int, timeRange *models.TimeRange) (*models.Stage, []*models.DropInfo, error) {
	codeMap := make(map[string]string)
	for _, lang := range consts.Languages {
		codeMap[lang] = gamedataStage.Code
	}
	code, err := json.Marshal(codeMap)
	if err != nil {
		return nil, nil, err
	}

	existenceMap := make(map[string]map[string]any)
	for _, s := range consts.Servers {
		existenceMap[s] = map[string]any{
			"exist": false,
		}
		if s == server {
			existenceMap[s]["exist"] = true
			existenceMap[s]["openTime"] = timeRange.StartTime.UnixMilli()
			fakeEndTime := time.UnixMilli(consts.FakeEndTimeMilli)
			endTime := &fakeEndTime
			if timeRange.EndTime != nil {
				endTime = timeRange.EndTime
			}
			existenceMap[s]["closeTime"] = endTime.UnixMilli()
		}
	}
	existence, err := json.Marshal(existenceMap)
	if err != nil {
		return nil, nil, err
	}

	stage := &models.Stage{
		ArkStageID: gamedataStage.StageID,
		ZoneID:     zoneId,
		StageType:  gamedataStage.StageType,
		Sanity:     null.IntFrom(int64(gamedataStage.ApCost)),
		Code:       code,
		Existence:  existence,
	}

	itemsMap, err := s.ItemService.GetItemsMapByArkId(ctx)
	if err != nil {
		return nil, nil, err
	}
	var activityToken string
	for _, reward := range gamedataStage.StageDropInfo.DisplayDetailRewards {
		if reward.Type == consts.ItemTypeActivity && activityToken == "" {
			activityToken = reward.Id
			break
		}
	}

	extrasMap := map[string]string{}
	if activityToken != "" {
		extrasMap["arkItemId"] = activityToken
	}

	groupedRewards := make(map[int][]*gamedata.DisplayDetailReward)
	groupedRewards[2] = make([]*gamedata.DisplayDetailReward, 0)
	groupedRewards[3] = make([]*gamedata.DisplayDetailReward, 0)
	groupedRewards[4] = make([]*gamedata.DisplayDetailReward, 0)
	for _, reward := range gamedataStage.StageDropInfo.DisplayDetailRewards {
		if reward.DropType > 4 || reward.DropType < 2 {
			continue
		}
		groupedRewards[reward.DropType] = append(groupedRewards[reward.DropType], reward)
	}

	dropInfos := make([]*models.DropInfo, 0)
	for dropType, rewards := range groupedRewards {
		items := make([]*models.Item, 0)
		for _, reward := range rewards {
			if reward.Type == consts.ItemTypeMaterial {
				item := itemsMap[reward.Id]
				items = append(items, item)
				bounds := s.decideItemBounds(item)
				dropInfos = append(dropInfos, &models.DropInfo{
					Server:      server,
					ItemID:      null.IntFrom(int64(item.ItemID)),
					DropType:    gdutils.RewardTypeMap[dropType],
					Accumulable: true,
					Bounds:      bounds,
				})
			}
		}

		// add dropinfo for dropType
		dropTypeBounds := s.decideDropTypeBounds(dropType, items)
		dropInfos = append(dropInfos, &models.DropInfo{
			Server:      server,
			DropType:    gdutils.RewardTypeMap[dropType],
			Accumulable: true,
			Bounds:      dropTypeBounds,
		})
	}

	// add dropinfo for furniture
	if gamedataStage.ApCost != 0 {
		item := itemsMap[consts.FurnitureArkItemID]
		dropInfos = append(dropInfos, &models.DropInfo{
			Server:      server,
			ItemID:      null.IntFrom(int64(item.ItemID)),
			DropType:    consts.DropTypeFurniture,
			Accumulable: true,
			Bounds:      &models.Bounds{Upper: 1, Lower: 0},
		})
	}

	// add dropinfo for recognition only item
	if activityToken != "" {
		extras, err := json.Marshal(extrasMap)
		if err != nil {
			return nil, nil, err
		}
		dropInfos = append(dropInfos, &models.DropInfo{
			Server:      server,
			DropType:    consts.DropTypeRecognitionOnly,
			Accumulable: false,
			Extras:      extras,
		})
	}

	linq.From(dropInfos).SortT(func(a, b *models.DropInfo) bool {
		if a.ItemID.Valid && b.ItemID.Valid || !a.ItemID.Valid && !b.ItemID.Valid {
			if a.DropType == b.DropType {
				return a.ItemID.Int64 < b.ItemID.Int64
			} else {
				return dropTypeOrderMapping[a.DropType] < dropTypeOrderMapping[b.DropType]
			}
		} else {
			return a.ItemID.Valid
		}
	}).ToSlice(&dropInfos)

	return stage, dropInfos, nil
}

func (s *GameDataService) decideItemBounds(item *models.Item) *models.Bounds {
	var upper int
	var lower int

	switch {
	case item.Rarity >= 2:
		upper = 1
		lower = 0
	case item.Rarity == 1:
		upper = 3
		lower = 0
	case item.Rarity == 0:
		upper = 5
		lower = 0
	}

	bounds := &models.Bounds{
		Upper: upper,
		Lower: lower,
	}
	return bounds
}

func (s *GameDataService) decideDropTypeBounds(dropType int, items []*models.Item) *models.Bounds {
	if dropType == 2 || dropType == 3 {
		return &models.Bounds{Upper: len(items), Lower: 0}
	}
	if dropType == 4 {
		if len(items) == 0 {
			return &models.Bounds{Upper: 0, Lower: 0}
		} else {
			return &models.Bounds{Upper: 1, Lower: 0}
		}
	}
	return nil
}
