package app

import (
	"encoding/json"
	"sort"
	"sync"

	"github.com/google/logger"
	"github.com/gorilla/websocket"
	"fmt"
)

type playerConn struct {
	name   string
	number int //1-4
	m      sync.Mutex
	ws     *websocket.Conn
	r      *room
}

func (p *playerConn) sendStatement(s *statement) {
	p.wsMessage(gameType, s) //TODO: send full game statement
}

// Send message to game chat
func (p *playerConn) sendAction(action gameAction) {
	p.wsMessage(actionType, action)
}

func (p *playerConn) sendMessage(message string) {
	p.wsMessage(messageType, message)
}

func (p *playerConn) start() {
	p.wsMessage(startType, "start")
}

func (p *playerConn) stop(pNumber int) {
	p.wsMessage(stopType, pNumber)
	p.m.Lock()
	p.ws.WriteMessage(websocket.CloseMessage, []byte{})
	p.m.Unlock()
}

func (p *playerConn) receiver() {
	logger.Info("Listening for playerConn " + p.name)
	defer p.close()
	for {
		_, message, err := p.ws.ReadMessage() // TODO parse message type
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
			if !ok {
				// TODO: handle error
				logger.Error("Error parsing message body")
				continue
			}
			p.r.message <- request
		case stopType:
			p.r.stop <- p.number
			logger.Infof("Player %v gave up", p.number)
			//return
		case gameType:
			//TODO: update statement
			p.r.statement.processStatement(p.number, buf.Body, p.r.timer)
			p.r.updateAll <- struct{}{}
		case actionType:
			fmt.Println("action type")
			action := p.r.statement.processStatement(p.number, buf.Body, p.r.timer)
			// todo: process action - validation, calculation, resulting action for another [players
			p.r.sendAction(action)
		default:
			p.r.updateAll <- struct{}{}
		}
	}
}

func (p *playerConn) wsMessage(s string, b interface{}) {
	p.m.Lock()
	//logger.Infof("Lock on %p\n", p)
	text, err := json.Marshal(wsMessage{Status: s, Body: b})
	if err != nil {
		logger.Error(err)
	}
	//logger.Infof("Sending message %s to %p\n", s, p.ws)
	err = p.ws.WriteMessage(websocket.TextMessage, text)
	if err != nil {
		panic(err)
	}
	p.m.Unlock()
	//logger.Infof("Unlock on %p\n", p)
}

// TODO: handle player error
func (p *playerConn) playerError() {
	// player send wrong data -> auto defeat, disconnect
}

func (h hand) sortHand() {
	sort.Strings(h)
}

func (p *playerConn) close() {
	p.m.Lock()
	p.ws.Close()
	p.m.Unlock()
}

func parseAction(action string) gameAction {
	return gameAction{}
}