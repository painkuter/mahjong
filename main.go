package main

import (
	"mahjong/app"
	"mahjong/app/common"
)

// just starter
func main() {
	l := common.InitLogging()
	defer l.Close()

	app.Main()
}
