package main

import (
	"mahjong/app"
	"mahjong/app/common/log"
)

// just starter
func main() {
	l := log.InitLogging()
	defer l.Close()

	app.Main()
}
