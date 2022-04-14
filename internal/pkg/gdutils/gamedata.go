package gdutils

import (
	"strings"

	"github.com/penguin-statistics/soracli/internal/consts"
	"github.com/penguin-statistics/soracli/internal/models/gamedata"
)

const (
	ArkStageIDMarkTraining      = "tr"
	ArkStageIDMarkStory         = "st"
	ArkStageIDMarkEx            = "ex"
	ArkStageIDMarkChallengeMode = "#f#"
)

const (
	StageTypeMain         = "MAIN"
	StageTypeSub          = "SUB"
	StageTypeDaily        = "DAILY"
	StageTypeActivity     = "ACTIVITY"
	StageTypeCampaign     = "CAMPAIGN"
	StageTypeSpecialStory = "SPECIAL_STORY"
	StageTypeGuide        = "GUIDE"
)

var RewardTypeMap = map[int]string{
	2: consts.DropTypeRegular,
	3: consts.DropTypeSpecial,
	4: consts.DropTypeExtra,
}

func IsTrainingStage(stage *gamedata.Stage) bool {
	return strings.HasPrefix(GetArkStageIDSecondPart(stage), ArkStageIDMarkTraining) || strings.HasPrefix(stage.StageID, ArkStageIDMarkTraining+"_")
}

func IsStoryStage(stage *gamedata.Stage) bool {
	return strings.HasPrefix(GetArkStageIDSecondPart(stage), ArkStageIDMarkStory) || strings.HasPrefix(stage.StageID, ArkStageIDMarkStory+"_")
}

func IsChallengeModeStage(stage *gamedata.Stage) bool {
	return strings.HasSuffix(GetArkStageIDSecondPart(stage), ArkStageIDMarkChallengeMode)
}

func IsNormalModeExStage(stage *gamedata.Stage) bool {
	return strings.HasPrefix(GetArkStageIDSecondPart(stage), ArkStageIDMarkEx) && !IsChallengeModeStage(stage)
}

func IsEasyDiffGroupStage(stage *gamedata.Stage) bool {
	return stage.DiffGroup == consts.StageDiffGroupEasy
}

func IsCampaignStage(stage *gamedata.Stage) bool {
	return stage.StageType == StageTypeCampaign
}

func IsGuideStage(stage *gamedata.Stage) bool {
	return stage.StageType == StageTypeGuide
}

func IsDailyStage(stage *gamedata.Stage) bool {
	return stage.StageType == StageTypeDaily
}

func GetZonePrefixFromArkZoneID(arkZoneID string) string {
	index := strings.Index(arkZoneID, "_zone")
	if index == -1 {
		return ""
	}
	return arkZoneID[0:index]
}

func GetArkStageIDSecondPart(stage *gamedata.Stage) string {
	zonePrefix := GetZonePrefixFromArkZoneID(stage.ZoneID)
	if zonePrefix == "" || len(zonePrefix) >= len(stage.StageID) || !strings.HasPrefix(stage.StageID, zonePrefix) {
		index := strings.Index(stage.StageID, "_")
		if index == -1 {
			return ""
		}
		return stage.StageID[index+1:]
	}
	return stage.StageID[len(zonePrefix)+1:]
}
