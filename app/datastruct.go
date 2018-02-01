package app

type wsMessage struct {
	Status string      `json:"status"` //
	Body   interface{} `json:"body"`
	//type
	//code
}

type statement struct {
	Players map[string]playerStatement `json:"players"`
	Wall    []string                   `json:"wall"`
	Wind    int                        `json:"wind"` //wind of game
	Step    int                        `json:"step"` //current player number
	Reserve []string                   `json:"reserve"`
}

type playerStatement struct {
	CurrentTile string     `json:"current_tile"` // [private]
	Hand        []string   `json:"hand"`         // [private]
	Discard     []string   `json:"discard"`      // [public]
	Open        [][]string `json:"open"`         // [public]
	Wind        int        `json:"wind"`         // [public]
}
type playerCommand struct {
	Status string   `json:"status"` // skip / discard / announce
	Meld   string   `json:"meld"`   // chow / pong / kong
	Tiles  []string `json:"tiles"`
}
