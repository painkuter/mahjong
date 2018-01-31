package main

import "mahjong/app"

// just starter
func main() {
	l := app.InitLogging()
	defer l.Close()

	app.Main()
}