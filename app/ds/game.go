package ds

type GameAction struct {
	Player int    `json:"player"`
	Action string `json:"action"` // skip / discard / announce
	Meld   string `json:"meld"`   // chow / pong / kong
	Value  Hand   `json:"value"`
}
