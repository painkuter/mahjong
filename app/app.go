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
	"mahjong/app/ds"

	"github.com/codemodus/parth"
	"github.com/gorilla/websocket"
)

var (
	Room        *room // active room for new players TODO use mutex for room
	application *App
)

func parseRoomName(w http.ResponseWriter, r *http.Request) string {
	roomUrl, err := parth.SegmentToString(r.URL.Path, -1)
	apperr.Check(err)
	if roomUrl == config.DefaultRoomURL {
		return application.ActiveRoomURL()
	} else {
		if len(roomUrl) != UrlLength {
			log.Info("Wrong room-Url")
			http.Error(w, "Room not found", 404)
			return ""
		}
	}
	return roomUrl
}

const roomPage = "view/index.html"
const roomPageOld = "view/index_old.html"

// roomHandler подключает нового игрока к комнате
func roomHandler(w http.ResponseWriter, r *http.Request) {
	roomUrl := parseRoomName(w, r)
	playerName, err := parth.SegmentToString(r.URL.Path, 0)
	apperr.Check(err)
	log.Infof(playerName)

	//var homeTempl = template.Must(template.ParseFiles(roomPageOld))
	var homeTempl = template.Must(template.ParseFiles(roomPage))
	data := ds.RoomResponse{r.Host, roomUrl, len(application.ActiveRoom().players) + 1}
	err = homeTempl.Execute(w, data)
	apperr.Check(err)
}

func appRoomHandler(w http.ResponseWriter, r *http.Request) {
	roomUrl := parseRoomName(w, r)
	data := ds.RoomResponse{r.Host, roomUrl, len(application.ActiveRoom().players) + 1}
	response, err := json.Marshal(data)
	apperr.Check(err)
	_, err = w.Write(response)
	apperr.Check(err)
	w.Header().Set("Content-Type", "application/json")
}

func roomsListHandler(w http.ResponseWriter, r *http.Request) {
	rooms := application.roomList()
	var roomsTempl = template.Must(template.ParseFiles("view/rooms.html"))
	data := struct {
		Rooms []string
	}{rooms}
	roomsTempl.Execute(w, data)
}

// newRoomHandler редиректит в новую комнату
func newRoomHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/room/"+NewRoom().Url, 302)
}

func ActiveRoom(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(application.ActiveRoomURL()))
}

func (a *App) ActiveRoomURL() string {
	return a.lastRoomUrl
}

func (a *App) ActiveRoom() *room {
	return a.rooms[a.lastRoomUrl]
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := NewWSConnection(w, r)
	if err != nil {
		return
	}

	roomURL, err := getRoomURL(r)
	if err != nil {
		http.Error(w, "Wrong room URL", 400)
	}
	playersCount := len(application.rooms[roomURL].players)

	playerName := getPlayerName(r, playersCount)
	log.Infof("Player %s has joined to room %s", playerName, roomURL)

	application.rooms[roomURL].AddPlayer(playerName, ws)
}

func NewWSConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{
		HandshakeTimeout: time.Second * 15,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		CheckOrigin:      websocket.IsWebSocketUpgrade,
	}
	ws, err := upgrader.Upgrade(w, r, http.Header{"Set-Cookie": {"sessionID=1234"}}) // fixme
	if e, ok := err.(websocket.HandshakeError); ok {
		log.Infof("Websocket handshake error: %s", e.Error())
		http.Error(w, "Not a websocket handshake", 400)
		return nil, e
	} else if err != nil {
		log.Error(err)
		http.Error(w, "Runtime error", 500)
		return nil, e
	}

	return ws, nil
}

func LiveHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("response"))
	apperr.Check(err)
	w.Header().Set("Content-Type", "application/json")
}

func Main() {
	NewApp()
	rh := http.RedirectHandler("/room", 301)
	http.Handle("/", rh)                             // Path to redirect to connect default room
	http.HandleFunc("/room/", roomHandler)           // Path to connect existed room
	http.HandleFunc("/room", roomHandler)            // Path to connect default room
	http.HandleFunc("/app/room", appRoomHandler)     // Path to connect default room
	http.HandleFunc("/rooms-list", roomsListHandler) // Rooms list
	http.HandleFunc("/new-room", newRoomHandler)     // Path to create new room -> redirecting to /room/[URL]
	http.HandleFunc("/ws", WsHandler)                // WebSocket handler
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/view/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	log.Infof("Handlers initialized. Server listening on:  http://%s", config.ADDR)
	if err := http.ListenAndServe(config.ADDR, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

type App struct {
	rooms       map[string]*room
	lastRoomUrl string
}

func NewApp() *App {

	application = &App{
		rooms: make(map[string]*room),
	}
	r := NewRoom()
	application.setRoom(r)
	return application
}

func (a *App) Room(url string) *room {
	return a.rooms[url]
}

func (a *App) setRoom(r *room) {
	log.Info("New room " + r.Url)
	if a.rooms == nil {
		a.rooms = map[string]*room{}
	}
	a.lastRoomUrl = r.Url
	a.rooms[r.Url] = r
}

func (a *App) roomList() []string {
	rooms := make([]string, len(a.rooms))
	for _, room := range application.rooms {
		rooms = append(rooms, room.Url)
	}
	return rooms
}

// getPlayerName возвращает установленное игроком имя или возвращает заглушку для анонима с учетом общего числа игроков
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
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return "", err
	}
	if len(params[config.DefaultRoomURL]) > 0 {
		if _, ok := application.rooms[params[config.DefaultRoomURL][0]]; ok { // looking for room by request param
			//room found
			return params[config.DefaultRoomURL][0], nil
		}
	}
	// this way is error
	log.Error("Error getting room-parameter")
	// Need to return 400 to client
	return application.ActiveRoomURL(), nil
}
