package models

import (
	"encoding/json"
	"time"
)

type Activity struct {
	ActivityID int             `bun:",pk,autoincrement" json:"id"`
	StartTime  *time.Time      `json:"startTime"`
	EndTime    *time.Time      `json:"endTime"`
	Name       json.RawMessage `json:"name"`
	Existence  json.RawMessage `json:"existence"`
}
