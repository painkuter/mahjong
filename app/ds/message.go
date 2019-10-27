package ds

type WsMessage struct {
	Status string      `json:"status"` //
	Body   interface{} `json:"body"`
}
