package app

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

type player struct {
	name   string
	number int //1-4
	ws     *websocket.Conn
	r      *room
	hand   []int   // hand
	dump   []int   // dump
	open   [][]int //open
}

func (p *player) sendStatement(s statement) {
	p.wsMessage("game", s) //TODO: send full game statement
}

// Send message to game chat
func (p *player) sendMessage(message string) {
	p.wsMessage("message", message)
}

func (p *player) start() {
	p.wsMessage("start", "start")
}

func (p *player) stop() {
	p.wsMessage("stop", "stop")
	p.ws.WriteMessage(websocket.CloseMessage, []byte{})
}

func (p *player) receiver() {
	fmt.Println("Listening for player " + p.name)
	for {
		_, message, err := p.ws.ReadMessage()
		if err != nil {
			p.r.stop <- p.number
			fmt.Println(err)
			break
			//panic("Error getting message from client")
		}
		// TODO: parse message here
		var buf wsMessage
		err = json.Unmarshal(message, &buf)
		check(err)

		switch buf.Status {
		case "message":
			p.r.message <- string(buf.Body.(string))
		case "game":
			p.r.updateAll <- struct{}{}
		default:
			p.r.updateAll <- struct{}{}
		}
	}
	p.ws.Close()
}

func (p *player) wsMessage(s string, b interface{}) {
	text, err := json.Marshal(wsMessage{Status: s, Body: b})
	if err != nil {
		panic(err)
	}
	p.ws.WriteMessage(websocket.TextMessage, text)
}
