package app

const (
	// Server
	ADDR = ":8079"
	logDir = "logs"
	logPrefix = "log_"
	timeout = 60*60*1000
	announceTime = 10*1000 // 10 sec to announce meld
	turnTime = 60*1000 // 1 minute for the turn
	charset   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	urlLength = 8 //room-url length

	// Game
	wallSize  = 136
	reserveSize = 16
	handSize = 13

	//wsMessage types:
	gameType  = "game"
	messageType = "message"
	playersType = "players"
	startType   = "start"
	stopType = "stop"

	//game commands:
	skipCommand = "skip"
	discardCommand = "discard"
	announceCommand = "announce"
	pigHandCommand = "pig_hand"
	// meld types
	chowType = "chow"
	pongType = "pong"
	kongType = "kong"
	mahjongType = "mahjong"
)
