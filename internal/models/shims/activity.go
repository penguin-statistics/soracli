package shims

import (
	"encoding/json"

	"gopkg.in/guregu/null.v3"
)

type Activity struct {
	ActivityID int             `bun:",pk,autoincrement" json:"-"`
	Start      int64           `json:"start"`
	End        null.Int        `json:"end,omitempty" swaggertype:"integer"`
	LabelI18n  json.RawMessage `json:"label_i18n" swaggertype:"object"`
	Existence  json.RawMessage `json:"existence" swaggertype:"object"`
}
