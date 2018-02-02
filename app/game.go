package app

import (
	"github.com/google/logger"
	ms "github.com/mitchellh/mapstructure"
)

//TODO: process game here

func (s *statement) processStatement(playerNumber int, command interface{}) {
	var c playerCommand
	err := ms.Decode(command, &c)
	if err != nil {
		logger.Warning(err)
		return
	}

	logger.Infof("Processing command %v", command)

	switch c.Status {
	case skipCommand:
		// skip timer after 3 skips
	case announceCommand:
		lastTile := s.Players[s.prevTurn()].Discard[:1][0]
		s.Players[playerNumber].Open = append(s.Players[playerNumber].Open, []string{lastTile})
		s.Step = playerNumber
	case discardCommand:
		if s.Step != playerNumber {
			logger.Warning("Wrong player number")
			return
		}
		if len(c.Tiles) > 1 {
			logger.Warning("Wrong tiles number in the command")
			return
		}
		p := s.Players[playerNumber]
		//TODO: remove this:
		if p.CurrentTile == "" {
			logger.Warning("Empty current tile")
			return
		}
		//
		p.Hand = append(p.Hand, p.CurrentTile)
		p.CurrentTile = ""
		p.Hand.cutTile(c.Tiles[0])

		p.Discard = append(p.Discard, c.Tiles[0])
		// timer for announce
		s.nextTurn()
	default:
		//TODO: handle error
	}
}

func (s *statement) getFromWall() {
	currentPlayer := s.Players[s.Step]
	currentPlayer.CurrentTile = s.Wall[0]
	s.Wall = s.Wall[1:]
}

func (h hand) cutTile(tile string) {
	for i, elem := range h {
		if tile == elem {
			h = append(h[:i], h[i+1:]...)
			return
		}
	}
	// TODO: handle error
	logger.Warning("Tile not found")
}

func (s *statement) nextTurn() {
	s.Step = (s.Step % 4) + 1
	s.getFromWall()
}

func (s *statement) prevTurn() int {
	return (s.Step+4)%4 - 1
}

func (s *statement) lastTile() string {
	return s.Players[s.prevTurn()].Discard[:1][0]
}
