package cache

import (
	"sync"

	"gopkg.in/guregu/null.v3"

	"github.com/penguin-statistics/soracli/internal/models"
	"github.com/penguin-statistics/soracli/internal/models/shims"
	"github.com/penguin-statistics/soracli/internal/models/types"
	"github.com/penguin-statistics/soracli/internal/pkg/cache"
)

type Flusher func() error

var (
	ItemDropSetByStageIDAndRangeID   *cache.Set[[]int]
	ItemDropSetByStageIdAndTimeRange *cache.Set[[]int]

	CliGameDataSeed *cache.Singular[types.CliGameDataSeedResponse]
	ItemByArkID     *cache.Set[models.Item]
	ShimItems       *cache.Singular[[]*shims.Item]
	ShimItemByArkID *cache.Set[shims.Item]
	ItemsMapById    *cache.Singular[map[int]*models.Item]
	ItemsMapByArkID *cache.Singular[map[string]*models.Item]

	Notices *cache.Singular[[]*models.Notice]

	Activities     *cache.Singular[[]*models.Activity]
	ShimActivities *cache.Singular[[]*shims.Activity]

	Stages           *cache.Singular[[]*models.Stage]
	StageByArkID     *cache.Set[models.Stage]
	ShimStages       *cache.Set[[]*shims.Stage]
	ShimStageByArkID *cache.Set[shims.Stage]
	StagesMapByID    *cache.Singular[map[int]*models.Stage]
	StagesMapByArkID *cache.Singular[map[string]*models.Stage]

	TimeRanges               *cache.Set[[]*models.TimeRange]
	TimeRangeByID            *cache.Set[models.TimeRange]
	TimeRangesMap            *cache.Set[map[int]*models.TimeRange]
	MaxAccumulableTimeRanges *cache.Set[map[int]map[int][]*models.TimeRange]

	Zones           *cache.Singular[[]*models.Zone]
	ZoneByArkID     *cache.Set[models.Zone]
	ShimZones       *cache.Singular[[]*shims.Zone]
	ShimZoneByArkID *cache.Set[shims.Zone]

	Properties map[string]string

	once sync.Once

	CacheSetMap             map[string]Flusher
	CacheSingularFlusherMap map[string]Flusher
)

func Initialize() {
	once.Do(func() {
		initializeCaches()
	})
}

func Delete(name string, key null.String) error {
	if key.Valid {
		if _, ok := CacheSetMap[name]; ok {
			if err := CacheSetMap[name](); err != nil {
				return err
			}
		}
	} else {
		if _, ok := CacheSingularFlusherMap[name]; ok {
			if err := CacheSingularFlusherMap[name](); err != nil {
				return err
			}
		} else if _, ok := CacheSetMap[name]; ok {
			if err := CacheSetMap[name](); err != nil {
				return err
			}
		}
	}
	return nil
}

func initializeCaches() {
	CacheSetMap = make(map[string]Flusher)
	CacheSingularFlusherMap = make(map[string]Flusher)

	// drop_info
	ItemDropSetByStageIDAndRangeID = cache.NewSet[[]int]("itemDropSet#server|stageId|rangeId")
	ItemDropSetByStageIdAndTimeRange = cache.NewSet[[]int]("itemDropSet#server|stageId|startTime|endTime")

	CacheSetMap["itemDropSet#server|stageId|rangeId"] = ItemDropSetByStageIDAndRangeID.Flush
	CacheSetMap["itemDropSet#server|stageId|startTime|endTime"] = ItemDropSetByStageIdAndTimeRange.Flush

	// item
	CliGameDataSeed = cache.NewSingular[types.CliGameDataSeedResponse]("cliGameDataSeed")
	ItemByArkID = cache.NewSet[models.Item]("item#arkItemId")
	ShimItems = cache.NewSingular[[]*shims.Item]("shimItems")
	ShimItemByArkID = cache.NewSet[shims.Item]("shimItem#arkItemId")
	ItemsMapById = cache.NewSingular[map[int]*models.Item]("itemsMapById")
	ItemsMapByArkID = cache.NewSingular[map[string]*models.Item]("itemsMapByArkId")

	CacheSingularFlusherMap["items"] = CliGameDataSeed.Delete
	CacheSetMap["item#arkItemId"] = ItemByArkID.Flush
	CacheSingularFlusherMap["shimItems"] = ShimItems.Delete
	CacheSetMap["shimItem#arkItemId"] = ShimItemByArkID.Flush
	CacheSingularFlusherMap["itemsMapById"] = ItemsMapById.Delete
	CacheSingularFlusherMap["itemsMapByArkId"] = ItemsMapByArkID.Delete

	// notice
	Notices = cache.NewSingular[[]*models.Notice]("notices")

	CacheSingularFlusherMap["notices"] = Notices.Delete

	// activity
	Activities = cache.NewSingular[[]*models.Activity]("activities")
	ShimActivities = cache.NewSingular[[]*shims.Activity]("shimActivities")

	CacheSingularFlusherMap["activities"] = Activities.Delete
	CacheSingularFlusherMap["shimActivities"] = ShimActivities.Delete

	// stage
	Stages = cache.NewSingular[[]*models.Stage]("stages")
	StageByArkID = cache.NewSet[models.Stage]("stage#arkStageId")
	ShimStages = cache.NewSet[[]*shims.Stage]("shimStages#server")
	ShimStageByArkID = cache.NewSet[shims.Stage]("shimStage#server|arkStageId")
	StagesMapByID = cache.NewSingular[map[int]*models.Stage]("stagesMapById")
	StagesMapByArkID = cache.NewSingular[map[string]*models.Stage]("stagesMapByArkId")

	CacheSingularFlusherMap["stages"] = Stages.Delete
	CacheSetMap["stage#arkStageId"] = StageByArkID.Flush
	CacheSetMap["shimStages#server"] = ShimStages.Flush
	CacheSetMap["shimStage#server|arkStageId"] = ShimStageByArkID.Flush
	CacheSingularFlusherMap["stagesMapById"] = StagesMapByID.Delete
	CacheSingularFlusherMap["stagesMapByArkId"] = StagesMapByArkID.Delete

	// time_range
	TimeRanges = cache.NewSet[[]*models.TimeRange]("timeRanges#server")
	TimeRangeByID = cache.NewSet[models.TimeRange]("timeRange#rangeId")
	TimeRangesMap = cache.NewSet[map[int]*models.TimeRange]("timeRangesMap#server")
	MaxAccumulableTimeRanges = cache.NewSet[map[int]map[int][]*models.TimeRange]("maxAccumulableTimeRanges#server")

	CacheSetMap["timeRanges#server"] = TimeRanges.Flush
	CacheSetMap["timeRange#rangeId"] = TimeRangeByID.Flush
	CacheSetMap["timeRangesMap#server"] = TimeRangesMap.Flush
	CacheSetMap["maxAccumulableTimeRanges#server"] = MaxAccumulableTimeRanges.Flush

	// zone
	Zones = cache.NewSingular[[]*models.Zone]("zones")
	ZoneByArkID = cache.NewSet[models.Zone]("zone#arkZoneId")
	ShimZones = cache.NewSingular[[]*shims.Zone]("shimZones")
	ShimZoneByArkID = cache.NewSet[shims.Zone]("shimZone#arkZoneId")

	CacheSingularFlusherMap["zones"] = Zones.Delete
	CacheSetMap["zone#arkZoneId"] = ZoneByArkID.Flush
	CacheSingularFlusherMap["shimZones"] = ShimZones.Delete
	CacheSetMap["shimZone#arkZoneId"] = ShimZoneByArkID.Flush
}
