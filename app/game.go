package app

import (
	"github.com/google/logger"
	ms "github.com/mitchellh/mapstructure"
)

//TODO: process game here

func (s *statement) processStatement(playerNumber int, command interface{}, timer chan struct{}) {
	var c playerCommand
	err := ms.Decode(command, &c)
	if err != nil {
		logger.Warning(err)
		return
	}

	logger.Infof("Processing command %v", command)

	switch c.Status {
	case skipCommand:
		// skip timer after 4 skips
		s.Pass[playerNumber] = true
		if s.Pass.checkSkip() {
			// pass timer
			timer <- struct{}{}
		}
	case announceCommand:
		if len(s.Players[s.prevTurn()].Discard) == 0 {
			logger.Warning("Empty discard")
			return
		}
		lastTile := s.Players[s.prevTurn()].Discard[:1][0]
		//remove last tile from discard
		s.Players[s.prevTurn()].Discard.cutTile(lastTile)

		//append last tile to the hand
		h := s.Players[playerNumber].Hand
		c.Tiles = append(c.Tiles, lastTile)
		ok := false
		switch c.Meld {
		case chowType:
			if playerNumber != (s.prevTurn()%4)+1 {
				// TODO: return error
				logger.Warning("Wrong turn for chow")
				return
			}
			ok = h.findChow(c.Tiles)
		case pongType:
			ok = h.findPong(c.Tiles)
		case kongType:
			ok = h.findKong(c.Tiles)
		case mahjongType:
			if !s.Players[playerNumber].IsReady {
				// TODO: return error: hand isn't ready
				return
			}
			//TODO: finish game, sand full statement
			logger.Infoln("MAHJONG!!!")
		}
		if !ok {
			return
		}

		s.Players[playerNumber].Open = append(s.Players[playerNumber].Open, c.Tiles)
		for _, tile := range c.Tiles {
			s.Players[playerNumber].Hand.cutTile(tile)
		}
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
	case readyHandCommand:
		s.Players[playerNumber].IsReady = true
		//TODO: announce to all players
		logger.Infoln("Ready hand!")
	default:
		logger.Error("Wrong client command")
		// TODO: finish with player's error
	}
}

func (s *statement) getFromWall() {
	currentPlayer := s.Players[s.Step]
	currentPlayer.CurrentTile = s.Wall[0]
	s.Wall = s.Wall[1:]
}

// removes tile from the hand
func (h *hand) cutTile(tile string) {
	for i, elem := range *h {
		if tile == elem {
			*h = append((*h)[:i], (*h)[i+1:]...) //TODO: fix pointers
			return
		}
	}
	// TODO: handle error
	logger.Warning("Tile not found")
}

// move turn to hte next player
func (s *statement) nextTurn() {
	s.Step = (s.Step % 4) + 1
	s.getFromWall()
}

// returns last player's number
func (s *statement) prevTurn() int {
	return (s.Step+4)%4 - 1
}

// returns last tile name
func (s *statement) lastTile() string {
	return s.Players[s.prevTurn()].Discard[:1][0]
}

func (p pass) checkSkip() bool {
	for _, el := range p {
		if !el {
			return false
		}
	}
	return true
}
