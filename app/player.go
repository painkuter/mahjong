package app

import (
	"encoding/json"
	"mahjong/app/common"
	"sort"
	"sync"

	"github.com/google/logger"
	"github.com/gorilla/websocket"
)

type playerConn struct {
	name   string
	number int //1-4
	lock   sync.Mutex
	ws     *websocket.Conn
	room   *room
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
	p.lock.Lock()
	p.ws.WriteMessage(websocket.CloseMessage, []byte{})
	p.lock.Unlock()
}

func (p *playerConn) receiver() {
	logger.Info("Listening for playerConn " + p.name)
	defer p.close()
	for {
		var buf WsMessage
		/*_, message,*/ err := p.ws.ReadJSON(&buf) // TODO parse message type
		if err != nil {
			p.room.stop <- p.number
			logger.Error(err)
			break
		}
		// TODO: parse message here
		//err = json.Unmarshal(message, &buf)
		//check(err)

		switch buf.Status {
		case messageType:
			request, ok := buf.Body.(string)
			if !ok {
				// TODO: handle error
				logger.Error("Error parsing message body")
				continue
			}
			p.room.message <- request
		case stopType:
			p.room.stop <- p.number
			logger.Infof("Player %v gave up", p.number)
			//return
		case gameType:
			//TODO: update statement
			p.room.statement.processStatement(p.number, buf.Body, p.room.timer)
			p.room.updateAll <- struct{}{}
		case actionType:
			//fmt.Println("action type")
			action := p.room.statement.processStatement(p.number, buf.Body, p.room.timer)
			// todo: process action - validation, calculation, resulting action for another [players
			p.room.sendAction(action)
		default:
			p.room.updateAll <- struct{}{}
		}
	}
}

func (p *playerConn) wsMessage(s string, b interface{}) {
	p.lock.Lock()
	//logger.Infof("Lock on %p\n", p)
	text, err := json.Marshal(WsMessage{Status: s, Body: b})
	if err != nil {
		logger.Error(err)
	}
	//logger.Infof("Sending message %s to %p\n", s, p.ws)
	err = p.ws.WriteMessage(websocket.TextMessage, text)
	common.Check(err)
	p.lock.Unlock()
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
	p.lock.Lock()
	common.Check(p.ws.Close())
	p.lock.Unlock()
}

func parseAction(action string) gameAction {
	return gameAction{}
}
