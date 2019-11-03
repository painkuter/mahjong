package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"

	"mahjong/app"
	"mahjong/app/apperr"
	"mahjong/app/common"
	"mahjong/app/common/log"
	"mahjong/app/config"
)

func die(w http.ResponseWriter, r *http.Request) {
	os.Exit(1)
}

func toRoom(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/ws?"+app.Room.Url, 302) // fixme
}

type testCon struct {
	conn []*websocket.Conn
}

func main() {

	l := log.InitLogging()
	defer l.Close()

	r := app.NewRoom()
	go func() {
		http.HandleFunc("/ws", app.WsHandler)
		http.HandleFunc("/room", app.ActiveRoom)
		http.HandleFunc("/to-room", toRoom)
		http.HandleFunc("/die", die)
		http.HandleFunc("/live", app.LiveHandler)
		err := http.ListenAndServe(config.ADDR, nil)
		if err != nil {
			log.Error(err)
			return
		}
	}()
	time.Sleep(1000 * time.Millisecond)

	//fmt.Println("Creating clients")
	// Convert http://127.0.0.1 to ws://127.0.0.1
	//u := "ws" + strings.TrimPrefix("http://0.0.0.0:8079", "http")
	//u := "ws" + strings.TrimPrefix(s.URL, "http")
	u := "ws://" + config.ADDR

	messageCh := make([]chan string, 4)
	var testPlayers testCon
	for i := 0; i < 3; i++ {
		messageCh[i] = make(chan string, 10)
		// Connect to the server
		url := u + "/ws?room=" + r.Url + "&name=player_" + strconv.Itoa(i)
		log.Info(url)
		ws, _, err := websocket.DefaultDialer.Dial(url, nil)
		apperr.Check(err)
		log.Infof("Adding player conn %p\n", ws)
		if err != nil {
			log.Fatalf("%v", err)
		}
		testPlayers.conn = append(testPlayers.conn, ws)
		go common.Receive(ws, messageCh[i])
		time.Sleep(100 * time.Millisecond)
	}
	time.Sleep(time.Millisecond * 100) // waiting for all players
	log.Info("Single game ready. Waiting for 4th played connection")

	for {
		testPlayers.makeTurn(r.Statement().Step, r.Statement().Players[r.Statement().Step])
		time.Sleep(100 * time.Millisecond)
	}
}

func (c testCon) makeTurn(step int, statement *app.PlayerStatement) {
	switch step {
	case 1, 2, 3:
		kong := statement.Hand.FindKong()
		if kong != nil {
			turn := fmt.Sprintf(`{"status":"action","body":{"action":"announce", "value":%s, "meld":"kong"}}`, kong.Print())
			err := c.conn[step-1].WriteMessage(websocket.TextMessage, []byte(turn))
			apperr.Check(err)
			return
		}

		pong := statement.Hand.FindKong()
		if pong != nil {
			turn := fmt.Sprintf(`{"status":"action","body":{"action":"announce", "value":%s, "meld":"kong"}}`, pong.Print())
			err := c.conn[step-1].WriteMessage(websocket.TextMessage, []byte(turn))
			apperr.Check(err)
			return
		}

		chow := statement.Hand.FindChow() // TODO add WithTile()
		if chow != nil {
			turn := fmt.Sprintf(`{"status":"action","body":{"action":"announce", "value":%s, "meld":"chow"}}`, chow.Print())
			err := c.conn[step-1].WriteMessage(websocket.TextMessage, []byte(turn))
			apperr.Check(err)
			return
		}

		turn := fmt.Sprintf(`{"status":"action","body":{"action":"discard", "value":["` + statement.Hand[0] + `"]}}`)
		err := c.conn[step-1].WriteMessage(websocket.TextMessage, []byte(turn))
		apperr.Check(err)

		// Поиск комбинации в руке + последний тайл из дискарда
	}
}
