package common

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"

	"mahjong/app/apperr"
	"mahjong/app/ds"
)

func ReceiveTo(ws *websocket.Conn, messageCh chan string) {
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			panic(err)
		}
		var buf ds.WsMessage
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

func Receive(ws *websocket.Conn) {
	for {
		_, message, err := ws.ReadMessage()
		if websocket.IsCloseError(err, 1001) {
			fmt.Println("connection closed")
			return
		}
		if err != nil {
			panic(err)
		}
		var buf ds.WsMessage
		err = json.Unmarshal(message, &buf)
		apperr.Check(err)
		var str string
		switch buf.Body.(type) {
		case []interface{}:
			str = fmt.Sprintf("%v", buf.Body)
		default:
			str = fmt.Sprintf("%v", buf.Body)
		}

		fmt.Println(str)
	}
}
