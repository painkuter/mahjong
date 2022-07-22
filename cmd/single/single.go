package main

import (
	"encoding/json"
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
	"mahjong/app/ds"
)

// die нужен для удаленного перезапуска приложения
func die(w http.ResponseWriter, _ *http.Request) { // TODO fix me
	w.Write([]byte("restarting"))
	os.Exit(1)
}

func newRoomHandler(a *app.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ws?"+a.ActiveRoomURL(), 302) // fixme
	}
}

type testCon struct {
	conn []*websocket.Conn
}

func main() {

	l := log.InitLogging()
	defer l.Close()

	application := app.NewApp()
	newRoom := application.ActiveRoomURL()
	go func() {
		http.HandleFunc("/ws", app.WsHandler)
		http.HandleFunc("/room", app.ActiveRoom)
		http.HandleFunc("/to-room", newRoomHandler(application))
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
	//host := "ws" + strings.TrimPrefix("http://0.0.0.0:8079", "http")
	//host := "ws" + strings.TrimPrefix(s.URL, "http")
	host := "ws://" + config.ADDR

	messageCh := make([]chan string, 4)
	var testPlayers testCon
	for i := 0; i < 3; i++ {
		messageCh[i] = make(chan string, 10)
		// Connect to the server
		url := host + "/ws?room=" + newRoom + "&name=player_" + strconv.Itoa(i)
		log.Info(url)
		ws, _, err := websocket.DefaultDialer.Dial(url, nil)
		apperr.Check(err)
		log.Infof("Adding player conn %p\n", ws)
		if err != nil {
			log.Fatalf("%v", err)
		}
		testPlayers.conn = append(testPlayers.conn, ws)
		go common.ReceiveTo(ws, messageCh[i])
		time.Sleep(100 * time.Millisecond)
	}
	time.Sleep(time.Millisecond * 100) // waiting for all players
	log.Info("Single game ready. Waiting for 4th played connection")

	for {
		testPlayers.makeTurn(application.Room(newRoom).Statement().Step, application.Room(newRoom).Statement().Players[application.Room(newRoom).Statement().Step])
		time.Sleep(100 * time.Millisecond)
	}
}

func (c testCon) makeTurn(step int, statement *app.PlayerStatement) {
	switch step {
	case 1, 2, 3:
		kong := statement.Hand.FindKong()
		if kong != nil {
			turn := c.printActionTurn("announce", kong)
			err := c.conn[step-1].WriteMessage(websocket.TextMessage, []byte(turn))
			apperr.Check(err)
			return
		}

		pong := statement.Hand.FindKong()
		if pong != nil {
			turn := c.printActionTurn("announce", pong)
			err := c.conn[step-1].WriteMessage(websocket.TextMessage, []byte(turn))
			apperr.Check(err)
			return
		}

		chow := statement.Hand.FindChow() // TODO add WithTile()
		if chow != nil {
			turn := c.printActionTurn("announce", chow)
			err := c.conn[step-1].WriteMessage(websocket.TextMessage, []byte(turn))
			apperr.Check(err)
			return
		}
		turn := c.printActionTurn("discard", statement.Hand)
		err := c.conn[step-1].WriteMessage(websocket.TextMessage, []byte(turn))
		apperr.Check(err)

		// Поиск комбинации в руке + последний тайл из дискарда
	}
}

func (c testCon) printActionTurn(action string, h ds.Hand) string {
	act := ds.GameAction{
		Action: action,
		Value:  h,
	}

	result := ds.WsMessage{
		Status: "action",
		Body:   act,
	}

	buf, err := json.Marshal(result)
	apperr.Check(err)
	return string(buf)
}
