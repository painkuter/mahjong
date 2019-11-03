package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"

	"mahjong/app/apperr"
	"mahjong/app/common"
	log2 "mahjong/app/common/log"
)

func TestWall(t *testing.T) {
	r := NewRoom()
	fmt.Println(len(r.statement.Wall))
	assert.Equal(t, len(r.statement.Wall) == wallSize-4*handSize-reserveSize, true)
}

func TestRoom_Run(t *testing.T) {
	l := log2.InitLogging()
	defer l.Close()

	r := NewRoom()
	//fmt.Println("Creating server")
	s := httptest.NewServer(http.HandlerFunc(WsHandler))
	defer s.Close()
	//fmt.Println("Creating clients")
	// Convert http://127.0.0.1 to ws://127.0.0.1
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	messageCh := make([]chan string, 4)
	var testPlayers []*websocket.Conn
	for i := 0; i < 4; i++ {
		messageCh[i] = make(chan string, 10)
		// Connect to the server
		url := u + "?room=" + r.Url + "&name=player_" + strconv.Itoa(i)
		log.Printf(url)
		ws, _, err := websocket.DefaultDialer.Dial(url, nil)
		log.Printf("Adding player conn %p\n", ws)
		if err != nil {
			t.Fatalf("%v", err)
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
	fmt.Println(<-messageCh[0]) //statement: map[...]

	time.Sleep(time.Millisecond * 100) // waiting for all players
	testPlayers[0].WriteMessage(websocket.TextMessage, []byte(`{"status":"action","body":{"action":"discard", "value":["1_7_1"]}}`))
	fmt.Println(<-messageCh[0]) //
	testPlayers[1].WriteMessage(websocket.TextMessage, []byte(`{"status":"action","body":{"action":"discard", "value":["1_1_2"]}}`))
	fmt.Println(<-messageCh[0]) //
	testPlayers[2].WriteMessage(websocket.TextMessage, []byte(`{"status":"action","body":{"action":"discard", "value":["1_4_2"]}}`))
	fmt.Println(<-messageCh[0]) //
	fmt.Println(<-messageCh[0]) //
	fmt.Println("done")
}

//func TestMain(lock *testing.M) {
//	l := common.InitLogging()
//	defer l.Close()
//	//os.Exit(lock.Run())
//}

func HelpReceiver(ws *websocket.Conn, messageCh chan string) {
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			panic(err)
		}
		var buf WsMessage
		err = json.Unmarshal(message, &buf)
		apperr.Check(err)
		var str string
		switch buf.Body.(type) {
		case []interface{}:
			str = fmt.Sprintf("%v", buf.Body)
		default:
			str = fmt.Sprintf("%v", buf.Body)
		}

		messageCh <- str
	}
}
