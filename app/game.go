package app

import (
	"github.com/google/logger"
	ms "github.com/mitchellh/mapstructure"
)

//TODO: process game here

func processStatement(command interface{}) {
	var c playerCommand
	err := ms.Decode(command, &c)
	if err != nil {
		logger.Warning(err)
		return
	}

	switch c.Status {
	case skipCommand:
		// skip timer after 3 skips
	case announceCommand:
	case discardCommand:
		// timer
	default:
		//TODO: handle error
	}
	// TODO: handle
}

func (s *statement) getFromWall() {
	currentPlayer := s.Players[s.Step]
	currentPlayer.CurrentTile = s.Wall[1]
	s.Wall = s.Wall[1:]
}
