package app

import (
	"fmt"
	"strconv"

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

	wall    []int
	wind    int //player number
	step    int // current player number
	reserve []int

	// update statement for all players
	updateAll chan struct{}
	// game chat message
	message chan string
	// stop chanel
	stop chan int
}

func newRoom() *room {
	fmt.Println("New room")
	return &room{
		url:       "",
		updateAll: make(chan struct{}),
		stop:      make(chan int),
		message:   make(chan string),
		wall:      []int{1, 2, 3, 4, 5}, //fill bones list ... (1-136)
		wind:      1,
	}
}

// AddPlayer adds new player to the room
func (r *room) AddPlayer(name string, ws *websocket.Conn) {
	if len(r.players) < 4 {
		p := player{name, len(r.players) + 1, ws, r, []int{3, 5, 6}, []int{}, [][]int{}} //TODO: randomize hand here
		r.players = append(r.players, p)
		//start game after 4th player connected
		if len(r.players) == 4 {
			go r.run()
			// TODO: create new room
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
	for _, p := range r.players {
		p.wsMessage("start","")
	}
	// waiting for some changes
	for {
		select {
		case <-r.updateAll:
			r.updateAllPlayers()
		case pNumber := <-r.stop:
			fmt.Println("Player #" + strconv.Itoa(pNumber) + " stopped the game")
			r.stopAllPlayers()
		case pMessage := <-r.message:
			r.sendMessageToAllPlayers(pMessage)
		}
	}
}

// need factory?
func (r *room) updateAllPlayers() {
	for _, p := range r.players {
		p.sendStatement()
	}
}

func (r *room) stopAllPlayers() {
	for _, p := range r.players {
		p.stop()
	}
}

func (r *room) sendMessageToAllPlayers(message string) {
	for _, p := range r.players {
		p.sendMessage(message)
	}
}
