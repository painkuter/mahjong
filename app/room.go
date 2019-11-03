package app

import (
	"mahjong/app/ds"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"log"

	"github.com/gorilla/websocket"
)

type room struct {
	Url string

	players []playerConn
	//players map[int]playerConn
	statement *statement

	// update statement for all players
	updateAll chan struct{}
	// timer
	timer chan struct{}
	// game chat message
	message chan string
	// stop chanel
	stop chan int

	lock sync.RWMutex
}

func NewRoom() *room {
	url := generateUrl()
	wall := generateWall()
	wall = randomizeWall(wall)
	wall, reserve := generateReserve(wall)
	east := randomEast()
	statement := &statement{
		Players: make(map[int]*PlayerStatement, 4),
		Reserve: reserve,
		East:    east,
		Wind:    1, // East = 1
		//Step:    east, // starting from east = 1
		Step: 0,
	}
	// Fill players statements
	for i := 1; i <= 4; i++ {
		var h ds.Hand
		wall, h = generateHand(wall) // fixme
		pStatement := PlayerStatement{Hand: h}
		statement.Players[i] = &pStatement
		//TODO: add wind
	}

	statement.Wall = wall

	r := &room{
		Url:       url,
		updateAll: make(chan struct{}),
		stop:      make(chan int),
		message:   make(chan string),
		statement: statement,
	}
	log.Printf("New room " + url)
	activeRooms[r.Url] = r
	Room = r
	return r
}

func (r *room) Statement() *statement {
	return r.statement
}

// AddPlayer adds new playerConn to the room
func (r *room) AddPlayer(name string, ws *websocket.Conn) {
	// add locks to room
	r.lock.Lock()
	defer r.lock.Unlock()
	log.Printf("New websocet %p\n", ws)
	if len(r.players) < 4 {
		p := playerConn{
			name:   name,
			number: len(r.players) + 1,
			lock:   sync.Mutex{},
			ws:     ws,
			room:   r,
		}
		r.players = append(r.players, p)
		r.statement.Players[len(r.players)].Name = name

		var players []string
		for _, p_ := range r.players { // kill me for this naming
			players = append(players, p_.name)
		}

		// push player to players list:
		for _, p_ := range r.players {
			p_.wsMessage(playersType, players)
		}

		//start game after 4th player connected
		if len(r.players) == 4 {
			go r.run()
			//TODO: check rooms list
			//Room = NewRoom() todo uncomment
		}
	} else {
		log.Fatal("Players count already equals four")
		//TODO: return some error
	}
}

func (r *room) run() {
	// creating receivers for all players
	for i := range r.players {
		go r.players[i].receiver()
	}

	// start the game
	log.Printf("Starting the game")
	for _, p := range r.players {
		//TODO: check the players
		p.start()
	}

	// first turn
	r.sendStartStatement()
	//room.statement.getFromWall()
	//room.updateAllPlayers()
	r.statement.nextTurn()

	// waiting for some changes
	for {
		select {
		case <-r.timer:
			r.statement.nextTurn()
		case <-r.updateAll:
			r.updateAllPlayers()
		case pNumber := <-r.stop:
			log.Printf("Player #%v stopped the game", pNumber)
			r.stopAllPlayers(pNumber)
		case pMessage := <-r.message:
			r.sendMessageToAllPlayers(pMessage)
		}
	}
}

// need factory?
func (r *room) updateAllPlayers() {
	for i, p := range r.players {
		p.sendStatement(r.statement.statementByPlayerNumber(i + 1))
	}
}

func (r *room) sendStartStatement() {
	for i, p := range r.players {
		p.sendStatement(r.statement.statementByPlayerNumber(i + 1))
	}
}

func (r *room) sendAction(action *gameAction) {
	if action == nil { // error actions
		return
	}
	for i, p := range r.players {
		p.sendAction(r.statement.actionByPlayer(i+1, action))
	}
}

func (r *room) stopAllPlayers(pNumber int) {
	for _, p := range r.players {
		p.stop(pNumber)
	}
}

func (r *room) sendMessageToAllPlayers(message string) {
	for _, p := range r.players {
		p.sendMessage(message)
	}
}

// commented for development
func randomizeWall(wall []string) []string {
	/*	list := rand.Perm(wallSize)
		w := make([]string, wallSize)
		for i, _ := range wall {
			w[i] = wall[list[i]]
		}
		return w*/
	return wall
}

func generateWall() []string {
	var wall []string
	for i := 1; i <= 4; i++ { // loop to multiply each tile by 4
		// suites
		for j := 1; j <= 9; j++ {
			for k := 1; k <= 3; k++ {
				wall = append(wall, strconv.Itoa(k)+"_"+strconv.Itoa(j)+"_"+strconv.Itoa(i))
			}
		}

		// winds
		for j := 1; j <= 4; j++ {
			wall = append(wall, strconv.Itoa(4)+"_"+strconv.Itoa(j)+"_"+strconv.Itoa(i))
		}

		// dragons
		for j := 1; j <= 3; j++ {
			wall = append(wall, strconv.Itoa(5)+"_"+strconv.Itoa(j)+"_"+strconv.Itoa(i))
		}
	}
	return wall
}

func generateReserve(w []string) (wall, reserve []string) {
	return w[reserveSize:], w[:reserveSize]
}

func generateHand(w []string) (wall, h ds.Hand) {
	h = make(ds.Hand, handSize)
	copy(h, w[:handSize])
	h.SortHand()
	return w[handSize:], h
}

func randomEast() int {
	//return rand.Intn(4) + 1 // 1-4, not 0-3
	return 1
}

func (s statement) statementByPlayerNumber(playerNumber int) *statement {
	// filter statement for selected playerConn (remove foreign hands, the wall and the reserve)
	privateStatement := &statement{
		Players: make(map[int]*PlayerStatement, 4),
		Step:    s.Step,
		Wind:    s.Wind,
		East:    s.East,
	}

	for j, player := range s.Players {
		if j == playerNumber {
			privateStatement.Players[100] = player
		} else {
			privateStatement.Players[j] = &PlayerStatement{
				Open:    player.Open,
				Discard: player.Discard,
				Name:    player.Name, // TODO add name
			}
			//privateStatement.Players[j] = player
		}
	}
	return privateStatement
}

func generateUrl() string {
	if true {
		return "AAA"
	}
	var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	url := make([]byte, urlLength)
	for i := range url {
		url[i] = charset[rnd.Intn(len(charset))]
	}
	return string(url)
}

func (s statement) actionByPlayer(player int, action *gameAction) gameAction {
	player = action.Player
	switch action.Action {
	case skipCommand:
		return gameAction{
			Player: player,
			Action: skipCommand,
		}
	case announceCommand:
		action.Player = player
	case discardCommand:
		action.Player = player
	case readyHandCommand:
		action.Player = player
	case getTileCommand:
		action.Player = player
	}
	return *action
}
