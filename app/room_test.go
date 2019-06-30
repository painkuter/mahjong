package app

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"mahjong/app/common"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
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

	for i := 0; i < 4; i++ {
		// Connect to the server
		ws, _, err := websocket.DefaultDialer.Dial(u+"?room="+r.url, nil)
		fmt.Printf("Adding player conn %p\n", ws)
		if err != nil {
			t.Fatalf("%v", err)
		}
		r.AddPlayer(fmt.Sprintf("test_player_%d", i), ws)
		//defer ws.Close()
	}
	os.Exit(1)
	fmt.Println("Running")
	r.run()

	fmt.Println("Testing")
	messageType, buf, err := r.players[0].ws.ReadMessage()
	fmt.Println(messageType)
	fmt.Println(string(buf))
	fmt.Println(err)
}

//func TestMain(m *testing.M) {
//	l := common.InitLogging()
//	defer l.Close()
//	//os.Exit(m.Run())
//}
