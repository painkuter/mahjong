package app

type wsMessage struct {
	Status string      `json:"status"` //
	Body   interface{} `json:"body"`
	//type
	//code
}

type statement struct {
	Players map[int]*playerStatement `json:"players"`
	Wall    hand                     `json:"wall"`
	East    int                      `json:"east"` //east-player number (1-4)
	Wind    int                      `json:"wind"` //wind of round (changing after 4 rounds)
	Step    int                      `json:"step"` //current player number (1-4)
	Reserve hand                     `json:"reserve"`
	Pass    pass                     `json:"-"`
}

type playerStatement struct {
	Name        string `json:"name"`         // [public]
	CurrentTile string `json:"current_tile"` // [private]
	Hand        hand   `json:"hand"`         // [private]
	//Available   []hand `json:"available"`    // [private]
	Discard hand   `json:"discard"` // [public]
	Open    []hand `json:"open"`    // [public]
	Wind    int    `json:"wind"`    // [public]
	IsReady bool   `json:"is_ready"`
}

type playerCommand struct {
	Status string `json:"status"` // skip / discard / announce
	Meld   string `json:"meld"`   // chow / pong / kong
	Tiles  hand   `json:"tiles"`
}

type pass map[int]bool

type hand []string

type roomResponse struct {
	Host     string `json:"host"`
	RoomName string `json:"room_name"`
	Players  int    `json:"players"`
}

type gameAction struct {
	Player int      `json:"player"`
	Action string   `json:"action"`
	Meld   string   `json:"meld"`
	Value  []string `json:"value"`
}
