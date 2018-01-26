package main

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

type player struct {
	name string
	ws   *websocket.Conn
	r    *room
	hand []int   // hand
	dump []int   // dump
	open [][]int //open
}

func (p *player) sentStatement() {
	t := time.Now().String()
	p.ws.WriteMessage(websocket.TextMessage, []byte("time = "+t+"\n")) //TODO: send full game statement
}

func (p *player) receiver() {
	fmt.Println("Listening for player " + p.name)
	for {
		_, message, err := p.ws.ReadMessage()
		if err != nil {
			panic("Error getting message from client")
		}
		p.r.updateAll <- struct{}{}
		fmt.Println(string(message))
	}
}
