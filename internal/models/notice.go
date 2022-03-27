package models

import (
	"encoding/json"

	"gopkg.in/guregu/null.v3"
)

type Notice struct {
	NoticeID  int             `bun:",pk,autoincrement" json:"id"`
	Existence json.RawMessage `json:"existence" swaggertype:"object"`
	Severity  null.Int        `json:"severity" swaggertype:"integer"`
	Content   json.RawMessage `json:"content_i18n"`
}
