package app

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/google/logger"
	"github.com/gorilla/websocket"
)

type room struct {
	url string

	players []player
	//players map[int]player
	statement *statement

	// update statement for all players
	updateAll chan struct{}
	// game chat message
	message chan string
	// stop chanel
	stop chan int
}

func newRoom() *room {
	url := generateUrl()
	wall := generateWall()
	wall = randomizeWall(wall)
	wall, reserve := generateReserve(wall)

	statement := &statement{
		Players: make(map[string]playerStatement, 4),
		Reserve: reserve,
		Wind:    randomWind(),
	}
	// Fill players statements
	for i := 1; i <= 4; i++ {
		var hand []string
		wall, hand = generateHand(wall)
		pStatement := playerStatement{Hand: hand}
		statement.Players[strconv.Itoa(i)] = pStatement
		//TODO: add wind
	}

	statement.Wall = wall

	r := &room{
		url:       url,
		updateAll: make(chan struct{}),
		stop:      make(chan int),
		message:   make(chan string),
		statement: statement,
	}
	logger.Info("New room " + url)
	activeRooms_[r.url] = r
	return r
}

// AddPlayer adds new player to the room
func (r *room) AddPlayer(name string, ws *websocket.Conn) {
	if len(r.players) < 4 {
		p := player{name, len(r.players) + 1, ws, r}
		r.players = append(r.players, p)

		// push player to players list:
		for _, p_ := range r.players {
			var players []string
			for _, p__ := range r.players { // kill me for this naming
				players = append(players, p__.name)
			}
			p_.wsMessage(playersType, players)
		}

		//start game after 4th player connected
		if len(r.players) == 4 {
			go r.run()
			//TODO: check rooms list
			Room = newRoom()
		}
	} else {
		logger.Fatal("Players count already equals four")
		//TODO: return some error
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
		p.start()
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
	for i, p := range r.players {
		p.sendStatement(r.statement.statementByPlayerNumber(i))
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

func randomizeWall(wall []string) []string {
	list := rand.Perm(wallSize)
	w := make([]string, wallSize)
	for i, _ := range wall {
		w[i] = wall[list[i]]
	}
	return w
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

func generateHand(w []string) (wall, hand []string) {
	return w[handSize:], w[:handSize]
}

func randomWind() int {
	return rand.Intn(4)
}

func (s statement) statementByPlayerNumber(playerNumber int) statement {
	// filter statement for selected player (remove foreign hands, the wall and thr reserve)
	privateStatement := statement{
		Players: make(map[string]playerStatement, 4),
		Step: s.Step,
		Wind: s.Wind,
	}

	for j, player := range s.Players {
		if j == strconv.Itoa(playerNumber) {
			privateStatement.Players["me"] = player
		} else {
			privateStatement.Players[j] = playerStatement{
				Open:    player.Open,
				Discard: player.Discard,
			}
		}
	}

	return privateStatement
}

func generateUrl() string {
	var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	url := make([]byte, urlLength)
	for i := range url {
		url[i] = charset[rnd.Intn(len(charset))]
	}
	return string(url)
}
