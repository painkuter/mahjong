package app

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/url"
	"time"

	"mahjong/app/config"

	"github.com/codemodus/parth"
	"github.com/google/logger"
	"github.com/gorilla/websocket"
)

var (
	activeRooms = make(map[string]*room)
	Room        *room // active room for new players TODO use mutex for room
)

func roomHandler(w http.ResponseWriter, r *http.Request) {
	roomUrl, err := parth.SegmentToString(r.URL.Path, -1)
	check(err)
	if roomUrl == "room" {
		roomUrl = Room.url
	} else {
		if len(roomUrl) != urlLength {
			logger.Error("Wrong room-url")
			http.Error(w, "Room not found", 404)
			return
		}
	}

	playerName, err := parth.SegmentToString(r.URL.Path, 0)
	check(err)
	logger.Info(playerName)

	var homeTempl = template.Must(template.ParseFiles("view/index_old.html"))
	data := roomResponse{r.Host, roomUrl, len(Room.players) + 1}
	homeTempl.Execute(w, data)
}

func appRoomHandler(w http.ResponseWriter, r *http.Request) {
	roomUrl, err := parth.SegmentToString(r.URL.Path, -1)
	check(err)
	if roomUrl == "room" {
		roomUrl = Room.url
	} else {
		if len(roomUrl) != urlLength {
			logger.Error("Wrong room-url")
			http.Error(w, "Room not found", 404)
			return
		}
	}
	data := roomResponse{r.Host, roomUrl, len(Room.players) + 1}
	response, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	w.Write([]byte(response))
	w.Header().Set("Content-Type", "application/json")
}

func roomsListHandler(w http.ResponseWriter, r *http.Request) {
	var rooms []string
	for _, room := range activeRooms {
		rooms = append(rooms, room.url)
	}
	var roomsTempl = template.Must(template.ParseFiles("view/rooms.html"))
	data := struct {
		Rooms []string
	}{rooms}
	roomsTempl.Execute(w, data)
}

func newRoomHandler(w http.ResponseWriter, r *http.Request) {
	room := newRoom()
	http.Redirect(w, r, "/room/"+room.url, 302)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	//logger.Info("ws handler")

	upgrader := websocket.Upgrader{
		HandshakeTimeout: time.Second,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		logger.Error(err)
		http.Error(w, "Runtime error", 500)
		return
	}

	playerName := getPlayerName(r)
	roomURL := getRoomURL(r)
	logger.Infof("Player %s has joined to room %s", playerName, roomURL)

	activeRooms[roomURL].AddPlayer(playerName, ws)
}

func Main() {
	Room = newRoom()
	rh := http.RedirectHandler("/room", 301)
	http.Handle("/", rh)                             // Path to redirect to connect default room
	http.HandleFunc("/room", roomHandler)            // Path to connect default room
	http.HandleFunc("/app/room", appRoomHandler)     // Path to connect default room
	http.HandleFunc("/rooms-list", roomsListHandler) // Rooms list
	http.HandleFunc("/room/", roomHandler)           // Path to connect existed room
	http.HandleFunc("/new-room", newRoomHandler)     // Path to create new room -> redirecting to /room/[URL]
	http.HandleFunc("/ws", wsHandler)                // WebSocket handler
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/view/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	logger.Infof("Handlers initialized. Serve listening on: %s", config.ADDR)
	if err := http.ListenAndServe(config.ADDR, nil); err != nil {
		logger.Fatal("ListenAndServe:", err)
	}
}

func getPlayerName(r *http.Request) string {
	playerName := "Anonymous Player" //TODO: add number to anonymous name
	params, _ := url.ParseQuery(r.URL.RawQuery)
	if len(params["name"]) > 0 {
		playerName = params["name"][0]
	}
	return playerName
}

func getRoomURL(r *http.Request) string {
	params, _ := url.ParseQuery(r.URL.RawQuery)
	if len(params["room"]) > 0 {
		if _, ok := activeRooms[params["room"][0]]; ok { // looking for room by request param
			//room found
			return params["room"][0]
		}
		return Room.url
	}
	// this way is error
	logger.Error("Error getting room-parameter")
	// Need to return 400 to client
	return Room.url
}
