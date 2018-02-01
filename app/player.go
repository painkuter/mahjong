package app

import (
	"encoding/json"
	"fmt"

	"github.com/google/logger"
	"github.com/gorilla/websocket"
)

type player struct {
	name   string
	number int //1-4
	ws     *websocket.Conn
	r      *room
}

func (p *player) sendStatement(s statement) {
	p.wsMessage(gameType, s) //TODO: send full game statement
}

// Send message to game chat
func (p *player) sendMessage(message string) {
	p.wsMessage(messageType, message)
}

func (p *player) start() {
	p.wsMessage(startType, "start")
}

func (p *player) stop() {
	p.wsMessage(stopType, "stop")
	p.ws.WriteMessage(websocket.CloseMessage, []byte{})
}

func (p *player) receiver() {
	logger.Info("Listening for player " + p.name)
	for {
		_, message, err := p.ws.ReadMessage()
		if err != nil {
			p.r.stop <- p.number
			logger.Error(err)
			break
		}
		// TODO: parse message here
		var buf wsMessage
		err = json.Unmarshal(message, &buf)
		check(err)

		switch buf.Status {
		case messageType:
			request, ok := buf.Body.(string)
			if !ok{
				// TODO: handle error
				logger.Error("Error parsing message body")
				continue
			}
			p.r.message <- request
		case gameType:
			//TODO: update statement
			processStatement(buf.Body)
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
		logger.Error(err)
	}
	p.ws.WriteMessage(websocket.TextMessage, text)
}
