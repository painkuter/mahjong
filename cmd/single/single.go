package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"

	"mahjong/app"
	"mahjong/app/common"
	"mahjong/app/config"
)

func main() {

	r := app.NewRoom()
	go func() {
		http.HandleFunc("/ws", app.WsHandler)
		http.HandleFunc("/live", app.LiveHandler)
		err := http.ListenAndServe(config.ADDR, nil)
		if err != nil {
			fmt.Println(err)
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
	var testPlayers []*websocket.Conn
	for i := 0; i < 3; i++ {
		messageCh[i] = make(chan string, 10)
		// Connect to the server
		url := u + "/ws?room=" + r.Url + "&name=player_" + strconv.Itoa(i)
		log.Printf(url)
		ws, _, err := websocket.DefaultDialer.Dial(url, nil)
		common.Check(err)
		log.Printf("Adding player conn %p\n", ws)
		if err != nil {
			log.Fatalf("%v", err)
		}
		testPlayers = append(testPlayers, ws)
		go common.Receive(ws, messageCh[i])
		time.Sleep(100 * time.Millisecond)
	}
	log.Printf("Game started")

	<-messageCh[0]              //[player_0]
	<-messageCh[0]              //[player_0 player_1]
	<-messageCh[0]              //[player_0 player_1 player_2]
	<-messageCh[0]              //[player_0 player_1 player_2 player_3]
	<-messageCh[0]              //start
	fmt.Println(<-messageCh[0]) //map

	time.Sleep(time.Millisecond * 100) // waiting for all players
	testPlayers[0].WriteMessage(websocket.TextMessage, []byte(`{"status":"action","body":{"action":"discard", "value":["1_7_1"]}}`))
	fmt.Println(<-messageCh[0]) //
	testPlayers[1].WriteMessage(websocket.TextMessage, []byte(`{"status":"action","body":{"action":"discard", "value":["1_1_2"]}}`))
	fmt.Println(<-messageCh[0]) //
	fmt.Println(<-messageCh[0]) //
	fmt.Println(<-messageCh[0]) //
	fmt.Println("done")
	for {
		makeTurn(r.Statement().Step, r.Statement().Players[r.Statement().Step])
	}
}

func makeTurn(step int, statement *app.PlayerStatement) {
	switch step {
	case 0, 1, 2:
		//statement.Hand
		// Поиск комбинации в руке + последний тайл из дискарда
	}
}
