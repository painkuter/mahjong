package common

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"

	"mahjong/app/ds"
)

func Receive(ws *websocket.Conn, messageCh chan string) {
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			panic(err)
		}
		var buf ds.WsMessage
		err = json.Unmarshal(message, &buf)
		Check(err)
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
