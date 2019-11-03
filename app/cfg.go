package app

const (
	// Server
	//ADDR = "0.0.0.0:8079"
	ADDR         = "0.0.0.0:8080"
	LogDir       = "logs"
	LogPrefix    = "log_"
	timeout      = 60 * 60 * 1000
	announceTime = 10 * 1000 // 10 sec to announce meld
	turnTime     = 60 * 1000 // 1 minute for the turn
	Charset      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	UrlLength    = 12 //room-url length

	// Game
	wallSize    = 136
	reserveSize = 16
	handSize    = 13

	//wsMessage types:
	gameType    = "game"    // игровые команды
	actionType  = "action"  // action
	messageType = "message" // чат
	playersType = "players" // добавление игрока в список
	startType   = "start"   // начало игры
	stopType    = "stop"    // окончание игры

	//game commands:
	skipCommand      = "skip"
	discardCommand   = "discard"
	announceCommand  = "announce"
	readyHandCommand = "ready_hand"
	getTileCommand   = "get_tile"
	// meld types
	ChowType    = "chow"
	PongType    = "pong"
	KongType    = "kong"
	MahjongType = "mahjong"
)
