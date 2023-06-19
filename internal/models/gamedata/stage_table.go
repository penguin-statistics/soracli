package gamedata

type StageTable struct {
	Stages map[string]*Stage `json:"stages"`
}

type RetroTable struct {
	StageList map[string]*Stage `json:"stageList"`
}

type Stage struct {
	StageID       string         `json:"stageId"`
	StageType     string         `json:"stageType"`
	ApCost        int            `json:"apCost"`
	Code          string         `json:"code"`
	ZoneID        string         `json:"zoneId"`
	StageDropInfo *StageDropInfo `json:"stageDropInfo"`
	DiffGroup     string         `json:"diffGroup"`
}

type StageDropInfo struct {
	DisplayDetailRewards []*DisplayDetailReward `json:"displayDetailRewards"`
}

type DisplayDetailReward struct {
	Id       string `json:"id"`
	DropType string `json:"dropType"`
	Type     string `json:"type"`
}
