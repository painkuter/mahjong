package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"

	"mahjong/app"
	"mahjong/app/common"
	"mahjong/app/common/log"
	"mahjong/app/config"
	"mahjong/app/ds"
)

var ws *websocket.Conn

func main() {

	l := log.InitLogging()
	defer l.Close()

	//r := app.NewRoom()
	go func() {
		http.HandleFunc("/ws", wsHandler)
		http.HandleFunc("/room", room)
		//http.HandleFunc("/to-room", toRoom)
		//http.HandleFunc("/die", die)
		fs := http.FileServer(http.Dir("static"))
		http.Handle("/static/", http.StripPrefix("/static/", fs))

		http.HandleFunc("/live", app.LiveHandler)
		http.HandleFunc("/view/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, r.URL.Path[1:])
		})
	}()

	time.Sleep(100 * time.Millisecond)
	if err := http.ListenAndServe(config.ADDR, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func room(w http.ResponseWriter, r *http.Request) {
	var homeTempl = template.Must(template.ParseFiles("./view/index.html"))
	data := ds.RoomResponse{Host: r.Host, RoomName: "AAA", Players: 4}
	err := homeTempl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(string(request))
	ws, err = app.NewWSConnection(w, r)
	if err != nil {
		log.Fatal(err)
	}
	go common.Receive(ws)

	writeWSMessage(getDataFromJSON("./examples/e2e_data/game_data.json"))
	time.Sleep(time.Second)

	writeWSMessage(getDataFromJSON("./examples/e2e_data/announce_player_100.json"))

	writeWSMessage(getDataFromJSON("./examples/e2e_data/announce_player_1.json"))
	writeWSMessage(getDataFromJSON("./examples/e2e_data/announce_player_2.json"))
	writeWSMessage(getDataFromJSON("./examples/e2e_data/announce_player_3.json"))

	writeWSMessage(getDataFromJSON("./examples/e2e_data/discard.json"))
	writeWSMessage(getDataFromJSON("./examples/e2e_data/discard_2.json"))
	writeWSMessage(getDataFromJSON("./examples/e2e_data/discard_3.json"))
	writeWSMessage(getDataFromJSON("./examples/e2e_data/discard_4.json"))
	writeWSMessage(getDataFromJSON("./examples/e2e_data/discard_5.json"))
}

func getDataFromJSON(fileName string) []byte {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	return buf
}

func writeWSMessage(message []byte) {
	time.Sleep(time.Millisecond * 200)
	err := ws.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Fatal(err)
	}
}
