package app

import (
	"github.com/google/logger"
	ms "github.com/mitchellh/mapstructure"
)

//TODO: process game here
// processStatement processes players commands
func (s *statement) processStatement(playerNumber int, command interface{}, timer chan struct{}) *gameAction {
	var comm gameAction
	err := ms.Decode(command, &comm)
	if err != nil {
		logger.Warning(err)
		return nil
	}

	logger.Infof("Processing command %v", command)

	switch comm.Action {
	case skipCommand:
		// skip timer after 4 skips
		s.Pass[playerNumber] = true
		if s.Pass.checkSkip() {
			// pass timer
			timer <- struct{}{}
		}
	case announceCommand:
		s.lock.Lock()
		defer s.lock.Unlock()

		if len(s.Players[s.prevTurn()].Discard) == 0 {
			logger.Warning("Empty discard")
			return nil
		}
		lastTile := s.Players[s.prevTurn()].Discard[:1][0]
		//remove last tile from discard
		s.Players[s.prevTurn()].Discard.cutTile(lastTile)

		//append last tile to the hand
		h := s.Players[playerNumber].Hand
		comm.Value = append(comm.Value, lastTile)
		ok := false
		switch comm.Meld {
		case chowType:
			if playerNumber != (s.prevTurn()%4)+1 {
				// TODO: return error
				logger.Warning("Wrong turn for chow")
				return nil
			}
			ok = h.checkChow(comm.Value)
		case pongType:
			ok = h.checkPong(comm.Value)
		case kongType:
			ok = h.checkKong(comm.Value)
		case mahjongType:
			if !s.Players[playerNumber].IsReady {
				// TODO: return error: hand isn't ready
				return nil
			}
			//TODO: finish game, sand full statement
			logger.Infoln("MAHJONG!!!")
		}
		if !ok {
			return nil
		}

		s.Players[playerNumber].Open = append(s.Players[playerNumber].Open, comm.Value)
		for _, tile := range comm.Value {
			s.Players[playerNumber].Hand.cutTile(tile)
		}
		s.Step = playerNumber
	case discardCommand:
		s.lock.Lock()
		defer s.lock.Unlock()

		logger.Info(s.Step)
		if s.Step != playerNumber {
			logger.Warning("Wrong player number")
			return nil
		}
		if len(comm.Value) > 1 {
			logger.Warning("Wrong tiles number in the command")
			return nil
		}
		p := s.Players[playerNumber]
		//TODO: remove this:
		if p.CurrentTile == "" {
			logger.Warning("Empty current tile")
			return nil
		}
		//
		p.Hand = append(p.Hand, p.CurrentTile)
		p.CurrentTile = ""
		p.Hand.cutTile(comm.Value[0])

		p.Discard = append(p.Discard, comm.Value[0])
		// timer for announce
		s.nextTurn()
	case readyHandCommand:
		s.Players[playerNumber].IsReady = true
		//TODO: announce to all players
		logger.Infoln("Ready hand!")
	default:
		logger.Error("Wrong client command: " + comm.Action)
		// TODO: finish with player's error
	}
	return &comm
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
			*h = append((*h)[:i], (*h)[i+1:]...) //we use pointer to slice to avoid problems with capacity
			return
		}
	}
	// TODO: handle error
	logger.Warning("Tile not found")
}

// move turn to the next player
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
