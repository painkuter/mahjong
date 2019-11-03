package ds

import (
	"encoding/json"

	"mahjong/app/apperr"
)

type WsMessage struct {
	Status string      `json:"status"` //
	Body   interface{} `json:"body"`
}

func (ws WsMessage) Print() string {
	result, err := json.Marshal(ws)
	apperr.Check(err)
	return string(result)
}
