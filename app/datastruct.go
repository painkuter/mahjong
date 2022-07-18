package app

import (
	"sync"

	"mahjong/app/ds"
)

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
	lock      sync.RWMutex
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

// pass - состояние хода каждого игрока по номеру: true = пасанул, false = делает ход, nil = ждем ответа
type pass map[int]bool
