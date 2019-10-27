package app

import (
	"sort"
	"strconv"
	"sync"

	"mahjong/app/common"
)

type WsMessage struct {
	Status string      `json:"status"` //
	Body   interface{} `json:"body"`
	//type
	//code
}

type statement struct {
	Players map[int]*PlayerStatement `json:"players"`
	Wall    hand                     `json:"wall"`
	East    int                      `json:"east"` //east-player number (1-4)
	Wind    int                      `json:"wind"` //wind of round (changing after 4 rounds)
	Step    int                      `json:"step"` //current player number (1-4)
	Reserve hand                     `json:"reserve"`
	Pass    pass                     `json:"-"`
	lock    sync.RWMutex             `json:"-"`
}

type PlayerStatement struct {
	Name        string `json:"name"`         // [public]
	CurrentTile string `json:"current_tile"` // [private]
	Hand        hand   `json:"hand"`         // [private]
	//Available   []hand `json:"available"`    // [private]
	Discard hand   `json:"discard"` // [public]
	Open    []hand `json:"open"`    // [public]
	Wind    int    `json:"wind"`    // [public]
	IsReady bool   `json:"is_ready"`
}

type gameAction struct {
	Player int    `json:"player"`
	Action string `json:"action"` // skip / discard / announce
	Meld   string `json:"meld"`   // chow / pong / kong
	Value  hand   `json:"value"`
}

type pass map[int]bool

type hand []string

func (h hand) Int(i int) int {
	if i >= len(h) {
		return 0
	}
	v, err := strconv.ParseInt(string(h[i][0])+string(h[i][2]), 10, 64)
	common.Check(err)
	return int(v)
}

// implement sort.Interface
func (h hand) Len() int {
	return len(h)
}

func (h hand) Less(i, j int) bool {
	return h.Int(i) < h.Int(j)
}

func (h hand) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h hand) WithTile(tile string) hand {
	return append(h, tile)
}

func (h hand) CheckChow() bool {
	h2 := make(hand, len(h))
	copy(h2, h)
	sort.Sort(h2)
	var t1, t2 int
	for i := range h2 {
		if t1+1 == t2 && t2+1 == h2.Int(i) {
			return true
		}
		t1 = t2
		t2 = h2.Int(i)
	}
	// sort
	return false
}

func (h hand) CheckPong() bool {
	// map [string(1_2)] => count
	return false
}

func (h hand) CheckKong() bool {
	// map [string(1_2)] => count
	return false
}

type roomResponse struct {
	Host     string `json:"host"`
	RoomName string `json:"room_name"`
	Players  int    `json:"players"`
}

type gameActionOld struct {
	Player int      `json:"player"`
	Action string   `json:"action"`
	Meld   string   `json:"meld"`
	Value  []string `json:"value"`
}
