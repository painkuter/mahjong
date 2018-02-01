package app

type wsMessage struct {
	Status string      `json:"status"` //
	Body   interface{} `json:"body"`
	//type
	//code
}

type statement struct {
	Players map[int]playerStatement `json:"players"`
	Wall    []string                `json:"wall"`
	Wind    int                     `json:"wind"` //wind of game
	Step    int                     `json:"step"` //current player number
	Reserve []string                `json:"reserve"`
}

type playerStatement struct {
	Hand        []string   `json:"hand"`         // [private]
	Discard     []string   `json:"discard"`      // [public]
	Open        [][]string `json:"open"`         // [public]
	CurrentTail string     `json:"current_tail"` // [private]
	Wind        int        `json:"wind"`         // [public]
}
