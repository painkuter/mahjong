package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"

	"mahjong/app"
	"mahjong/app/common/log"
	"mahjong/app/config"
)

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
	time.Sleep(1000 * time.Millisecond)
	if err := http.ListenAndServe(config.ADDR, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func room(w http.ResponseWriter, r *http.Request) {
	var homeTempl = template.Must(template.ParseFiles("./view/index.html"))
	data := app.RoomResponse{Host: r.Host, RoomName: "AAA", Players: 4}
	homeTempl.Execute(w, data)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := app.NewWSConnection(w, r)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open("./examples/game.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	err = ws.WriteMessage(websocket.TextMessage, buf)
	//_, err = w.Write([]byte(buf))
	if err != nil {
		log.Fatal(err)
	}
}
