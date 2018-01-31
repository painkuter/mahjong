package app

const (
	// Server
	ADDR = ":8079"
	logFile = "logs/log_"
	timeout = 60*60*1000
	turntime = 60*1000 // 1 minute for the turn
	charset   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	urlLength = 8 //url length
	// Game
	wallSize  = 136
	//wsMessage types:
	gameType  = "game"
	messageType = "message"
	playersType = "players"
	startType   = "start"
	stopType = "stop"
)
