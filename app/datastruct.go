package app

import (
	"sync"

	"mahjong/app/ds"
)

type WsMessage struct {
	Status string      `json:"status"` //
	Body   interface{} `json:"body"`
	//type
	//code
}

type statement struct {
	Players   map[int]*PlayerStatement `json:"players"`
	LastTile  string                   `json:"last_tile,omitempty"`
	Wall      ds.Hand                  `json:"wall"`
	East      int                      `json:"east"`       //east-player number (1-4)
	Wind      int                      `json:"wind"`       //wind of round (changing after 4 rounds)
	Step      int                      `json:"step"`       //current player number (1-4)
	StepCount int                      `json:"step_count"` // порядковый номер хода
	Reserve   ds.Hand                  `json:"reserve,omitempty"`
	Pass      pass                     `json:"-"`
	lock      sync.RWMutex             `json:"-"`
}

type PlayerStatement struct {
	Name        string  `json:"name"`                   // [public]
	CurrentTile string  `json:"current_tile,omitempty"` // [private]
	Hand        ds.Hand `json:"hand,omitempty"`         // [private]
	//Available   []hand `json:"available"`    // [private]
	Discard ds.Hand   `json:"discard,omitempty"` // [public] у каждого свой дискард
	Open    []ds.Hand `json:"open,omitempty"`    // [public]
	Wind    int       `json:"wind"`              // [public]
	IsReady bool      `json:"is_ready,omitempty"`
}

func (ps *PlayerStatement) GetDiscard() *ds.Hand {
	if ps.Discard == nil {
		return &ds.Hand{}
	}
	return &ps.Discard
}

type gameAction struct {
	Player int     `json:"player"`
	Action string  `json:"action"` // skip / discard / announce
	Meld   string  `json:"meld"`   // chow / pong / kong
	Value  ds.Hand `json:"value"`
}

type pass map[int]bool

type RoomResponse struct {
	Host     string `json:"host"` // зачем здесь хост?
	RoomName string `json:"room_name"`
	Players  int    `json:"players"`
}

/*type gameActionOld struct {
	Player int      `json:"player"`
	Action string   `json:"action"`
	Meld   string   `json:"meld"`
	Value  []string `json:"value"`
}*/
