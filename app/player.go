package app

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"

	"mahjong/app/apperr"
	"mahjong/app/common/log"
	"mahjong/app/ds"
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
func (p *playerConn) sendAction(action GameAction) {
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
	err := p.ws.WriteMessage(websocket.CloseMessage, []byte{})
	apperr.Check(err)
	p.lock.Unlock()
}

func (p *playerConn) receiver() {
	log.Info("Listening for playerConn " + p.name)
	defer p.close()
	for {
		//p.lock.Lock()
		//fmt.Println("locked")
		var buf ds.WsMessage
		err := p.ws.ReadJSON(&buf) // TODO parse message type
		if err != nil {
			p.room.stop <- p.number
			log.Error(err)
			//p.lock.Unlock()
			break
		}
		//spew.Dump(buf)
		// TODO: parse message here
		//err = json.Unmarshal(message, &buf)
		//apperr.Check(err)
		p.handleMessage(buf)
		//p.lock.Unlock()
		//fmt.Println("unlock")
	}
}

func (p *playerConn) handleMessage(buf ds.WsMessage) {
	switch buf.Status {
	case messageType:
		request, ok := buf.Body.(string)
		if !ok {
			// TODO: handle error
			log.Error("Error parsing message body")
			return
		}
		p.room.message <- request
	case stopType:
		p.room.stop <- p.number
		log.Infof("Player %v gave up", p.number)
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

func (p *playerConn) wsMessage(s string, b interface{}) {
	p.lock.Lock()
	//log.Printf("Lock on %p\n", p)
	text, err := json.Marshal(WsMessage{Status: s, Body: b})
	if err != nil {
		log.Error(err)
	}
	log.Infof("Sending message %s:%s to %d\n", s, text, p.number)
	err = p.ws.WriteMessage(websocket.TextMessage, text) // TODO process error - disconnect player
	apperr.Check(err)
	p.lock.Unlock()
	//log.Printf("Unlock on %p\n", p)
}

// TODO: handle player error
func (p *playerConn) playerError() {
	// player send wrong data -> auto defeat, disconnect
}

/*func (h ds.Hand) sortHand() {
	sort.Strings(h)
}*/

func (p *playerConn) close() {
	p.lock.Lock()
	apperr.Check(p.ws.Close())
	p.lock.Unlock()
}

func parseAction(action string) GameAction {
	return GameAction{}
}
