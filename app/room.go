package app

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type room struct {
	url string

	players []player
	//players map[int]player
	//player1 *player
	//player2 *player
	//player3 *player
	//player4 *player

	wall []int
	wind int //player number
	step int // current player number
	reserve []int

	// update statement for all players
	updateAll chan struct{}
	// stop chanel
	stop chan struct{}
}

func newRoom() *room {
	return &room{
		url:       "",
		updateAll: make(chan struct{}),
		stop: make(chan struct{}),
		wall: []int{1, 2, 3, 4, 5}, //fill bones list ... (1-136)
		wind: 1,
	}
}

// AddPlayer adds new player to the room
func (r *room) AddPlayer(name string, ws *websocket.Conn) {
	if len(r.players) < 4 {
		p := player{name, ws, r, []int{3, 5, 6}, []int{},[][]int{}} //TODO: randomize hand here
		r.players = append(r.players, p)

		if len(r.players) == 4 {
			go r.run()
		}
	} else {
		panic("Players count already equals four")
	}
}

func (r *room) run() {
	// creating receivers for all players
	for _, p := range r.players {
		p_ := p
		go p_.receiver()
	}
	//start the game
	fmt.Println("Starting the game")
	// waiting for some changes
	for {
		select {
		case <-r.updateAll:
			r.updateAllPlayers()
		case <- r.stop:
			r.stopAllPlayers()
		}

	}
}
func (r *room) updateAllPlayers() {
	for _, p := range r.players {
		p.sentStatement()
	}
}

func (r *room) stopAllPlayers() {
	for _, p := range r.players {
		p.stop()
	}
}