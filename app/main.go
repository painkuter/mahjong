package app

import (
	"html/template"
	"net/http"
	"net/url"

	"github.com/codemodus/parth"
	"github.com/google/logger"
	"github.com/gorilla/websocket"
)

var (
	activeRooms  []room
	activeRooms_ = make(map[string]*room)
	Room         *room
	// Init templates
)

func roomHandler(c http.ResponseWriter, r *http.Request) {
	roomUrl, err := parth.SegmentToString(r.URL.Path, -1)
	check(err)
	if roomUrl == "room" {
		roomUrl = ""
	} else {
		if len(roomUrl) != urlLength {
			logger.Error("Wrong room-url")
			return // 400
		}
	}
	var homeTempl = template.Must(template.ParseFiles("view/index.html"))
	data := struct {
		Host     string
		RoomName string
		Players  int
	}{r.Host, roomUrl, len(Room.players)}
	homeTempl.Execute(c, data)
}

func roomsListHandler(w http.ResponseWriter, r *http.Request) {
	rooms := []string{}
	for _, room := range activeRooms_ {
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
	http.Redirect(w, r, "/room/"+room.url, 301)
}

func connectRoomHandler(w http.ResponseWriter, r *http.Request) {
	roomUrl := r.URL.Path[1:]
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	check(err)
	playerName := getPlayerName(r)
	activeRooms_[roomUrl].AddPlayer(playerName, ws)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		// TODO: logging
		// check(err)
		return
	}

	playerName := getPlayerName(r)

	Room.AddPlayer(playerName, ws)

	logger.Infof("Player %s has joined to room %s", playerName, Room.url)
}

func Main() {
	Room  = newRoom()
	rh := http.RedirectHandler("/room", 301)
	http.Handle("/", rh)                             // Path to redirect to connect default room
	http.HandleFunc("/room", roomHandler)            // Path to connect default room
	http.HandleFunc("/rooms-list", roomsListHandler) // Rooms list
	http.HandleFunc("/room/", roomHandler)           // Path to connect existed room
	http.HandleFunc("/new-room", newRoomHandler)     // Path to create new room -> redirectin to /room/[URL]
	http.HandleFunc("/ws", wsHandler)
	//http.HandleFunc("/new-room", newRoomWsHandler)

	http.HandleFunc("/view/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	logger.Info("Handlers initialized")
	if err := http.ListenAndServe(ADDR, nil); err != nil {
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
