package app

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"mahjong/app/apperr"
	"mahjong/app/common/log"
	"mahjong/app/config"

	"github.com/codemodus/parth"
	"github.com/gorilla/websocket"
)

var (
	activeRooms = make(map[string]*room)
	Room        *room // active room for new players TODO use mutex for room
)

func roomHandler(w http.ResponseWriter, r *http.Request) {
	roomUrl, err := parth.SegmentToString(r.URL.Path, -1)
	apperr.Check(err)
	if roomUrl == "room" {
		roomUrl = Room.Url
	} else {
		if len(roomUrl) != UrlLength {
			log.Info("Wrong room-Url")
			http.Error(w, "Room not found", 404)
			return
		}
	}

	playerName, err := parth.SegmentToString(r.URL.Path, 0)
	apperr.Check(err)
	log.Infof(playerName)

	//var homeTempl = template.Must(template.ParseFiles("view/index_old.html"))
	var homeTempl = template.Must(template.ParseFiles("view/index.html"))
	data := roomResponse{r.Host, roomUrl, len(Room.players) + 1}
	homeTempl.Execute(w, data)
}

func appRoomHandler(w http.ResponseWriter, r *http.Request) {
	roomUrl, err := parth.SegmentToString(r.URL.Path, -1)
	apperr.Check(err)
	if roomUrl == "room" {
		roomUrl = Room.Url
	} else {
		if len(roomUrl) != UrlLength {
			log.Info("Wrong room-Url")
			http.Error(w, "Room not found", 404)
			return
		}
	}
	data := roomResponse{r.Host, roomUrl, len(Room.players) + 1}
	response, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	w.Write(response)
	w.Header().Set("Content-Type", "application/json")
}

func roomsListHandler(w http.ResponseWriter, r *http.Request) {
	var rooms []string
	for _, room := range activeRooms {
		rooms = append(rooms, room.Url)
	}
	var roomsTempl = template.Must(template.ParseFiles("view/rooms.html"))
	data := struct {
		Rooms []string
	}{rooms}
	roomsTempl.Execute(w, data)
}

func newRoomHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/room/"+NewRoom().Url, 302)
}

func ActiveRoom(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(Room.Url))
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		HandshakeTimeout: time.Second,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		CheckOrigin:      websocket.IsWebSocketUpgrade,
	}
	ws, err := upgrader.Upgrade(w, r, http.Header{"Set-Cookie": {"sessionID=1234"}}) // fixme
	if e, ok := err.(websocket.HandshakeError); ok {
		log.Infof("Websocket handshake error: %s", e.Error())
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Error(err)
		http.Error(w, "Runtime error", 500)
		return
	}

	roomURL, err := getRoomURL(r)
	if err != nil {
		http.Error(w, "Wrong room URL", 400)
	}
	playersCount := len(activeRooms[roomURL].players)

	playerName := getPlayerName(r, playersCount)
	log.Infof("Player %s has joined to room %s", playerName, roomURL)

	activeRooms[roomURL].AddPlayer(playerName, ws)
}

func LiveHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("response"))
	apperr.Check(err)
	w.Header().Set("Content-Type", "application/json")
}

func Main() {
	NewRoom()
	rh := http.RedirectHandler("/room", 301)
	http.Handle("/", rh)                             // Path to redirect to connect default room
	http.HandleFunc("/room", roomHandler)            // Path to connect default room
	http.HandleFunc("/app/room", appRoomHandler)     // Path to connect default room
	http.HandleFunc("/rooms-list", roomsListHandler) // Rooms list
	http.HandleFunc("/room/", roomHandler)           // Path to connect existed room
	http.HandleFunc("/new-room", newRoomHandler)     // Path to create new room -> redirecting to /room/[URL]
	http.HandleFunc("/ws", WsHandler)                // WebSocket handler
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/view/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	log.Infof("Handlers initialized. Serve listening on: %s", config.ADDR)
	if err := http.ListenAndServe(config.ADDR, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func getPlayerName(r *http.Request, playersCount int) string {
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Error("Error getting player name: ", err)
	}
	if len(params["name"]) > 0 {
		return params["name"][0]
	}
	return "Anonymous Player #" + strconv.Itoa(playersCount+1)
}

func getRoomURL(r *http.Request) (string, error) {
	params, _ := url.ParseQuery(r.URL.RawQuery)
	if len(params["room"]) > 0 {
		if _, ok := activeRooms[params["room"][0]]; ok { // looking for room by request param
			//room found
			return params["room"][0], nil
		}
		return Room.Url, nil
	}
	// this way is error
	log.Error("Error getting room-parameter")
	// Need to return 400 to client
	return Room.Url, nil
}
