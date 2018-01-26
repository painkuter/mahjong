package main

import (
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

const ADDR = ":8079"

var Room = newRoom();

func homeHandler(c http.ResponseWriter, r *http.Request) {
	var homeTempl = template.Must(template.ParseFiles("view/index.html"))
	data := struct {
		Host       string
		Players int
	}{r.Host, len(Room.players)}
	homeTempl.Execute(c, data)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		return
	}

	playerName := "Player"
	params, _ := url.ParseQuery(r.URL.RawQuery)
	if len(params["name"]) > 0 {
		playerName = params["name"][0]
	}

	// TODO: Get or create a room
	//var room *room
	//if len(freeRooms) > 0 {
	//	for _, r := range freeRooms {
	//		room = r
	//		break
	//	}
	//} else {
	//	room = NewRoom("")
	//}

	Room.AddPlayer(playerName, ws)

	log.Printf("Player %s has joined to room", playerName)
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ws", wsHandler)

	http.HandleFunc("/view/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	if err := http.ListenAndServe(ADDR, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
