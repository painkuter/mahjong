package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/google/logger"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"

	"mahjong/app/common"
	"time"
)

func TestWall(t *testing.T) {
	r := newRoom()
	fmt.Println(len(r.statement.Wall))
	assert.Equal(t, len(r.statement.Wall) == wallSize-4*handSize-reserveSize, true)
}

func TestRoom_Run(t *testing.T) {
	l := common.InitLogging()
	defer l.Close()

	r := newRoom()
	fmt.Println("Creating server")
	s := httptest.NewServer(http.HandlerFunc(wsHandler))
	defer s.Close()
	fmt.Println("Creating clients")
	// Convert http://127.0.0.1 to ws://127.0.0.1
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	messageCh := make([]chan string, 4)
	var testPlayers []*websocket.Conn
	for i := 0; i < 4; i++ {
		messageCh[i] = make(chan string, 10)
		// Connect to the server
		url := u + "?room=" + r.url + "&name=player_" + strconv.Itoa(i)
		logger.Info(url)
		ws, _, err := websocket.DefaultDialer.Dial(url, nil)
		logger.Infof("Adding player conn %p\n", ws)
		if err != nil {
			t.Fatalf("%v", err)
		}
		testPlayers = append(testPlayers, ws)
		go testReceiver(ws, messageCh[i])
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("Testing")
	<-messageCh[0]              //[player_0]
	<-messageCh[0]              //[player_0 player_1]
	<-messageCh[0]              //[player_0 player_1 player_2]
	<-messageCh[0]              //[player_0 player_1 player_2 player_3]
	<-messageCh[0]              //start
	fmt.Println(<-messageCh[0]) //map

	time.Sleep(time.Millisecond * 100)
	testPlayers[0].WriteMessage(websocket.TextMessage, []byte(`{"status":"action","body":{"player":0, "action":"discard", "value":["1_7_1"]}}`))
	fmt.Println(<-messageCh[0]) //
	fmt.Println(<-messageCh[0]) //
	fmt.Println(<-messageCh[0]) //
	fmt.Println(<-messageCh[0]) //
	fmt.Println("done")
}

//func TestMain(m *testing.M) {
//	l := common.InitLogging()
//	defer l.Close()
//	//os.Exit(m.Run())
//}

func testReceiver(ws *websocket.Conn, messageCh chan string) {
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			panic(err)
		}
		var buf wsMessage
		err = json.Unmarshal(message, &buf)
		check(err)
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
