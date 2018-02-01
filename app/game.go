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
	case announceCommand:
	case discardCommand:
	default:
		//TODO: handle error
	}
	// TODO: handle
}
