package app

import (
	ms "github.com/mitchellh/mapstructure"

	"mahjong/app/common/log"
)

//TODO: process game here
// processStatement processes players commands
func (s *statement) processStatement(playerNumber int, command interface{}, timer chan struct{}) *gameAction {
	log.Infof("Processing command %v from player #%d", command, playerNumber)

	var comm gameAction
	err := ms.Decode(command, &comm)
	if err != nil {
		log.Info(err)
		return nil
	}
	comm.Player = playerNumber

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

		//if len(s.Players[s.prevTurn()].Discard) == 0 {
		//	log.Info("Empty discard")
		//	return nil
		//}
		lastTile := s.lastTile()
		//remove last tile from discard
		s.Players[s.prevTurn()].GetDiscard().CutTile(lastTile)

		//append last tile to the hand
		h := s.Players[playerNumber].Hand
		/*		if lastTile != "" {
				comm.Value = append(comm.Value, lastTile)
			}*/
		var ok bool
		switch comm.Meld {
		case chowType:
			if playerNumber != s.Step {
				//if playerNumber != (s.prevTurn()%4)+1 {
				// TODO: return error
				log.Info("Wrong turn for chow")
				return nil
			}
			ok = h.CheckChow(comm.Value)
		case pongType:
			ok = h.CheckPong(comm.Value)
		case kongType:
			ok = h.CheckKong(comm.Value)
		case mahjongType:
			if !s.Players[playerNumber].IsReady {
				// TODO: return error: hand isn't ready
				return nil
			}
			//TODO: finish game, sand full statement
			log.Infof("MAHJONG!!!")
		default:
			log.Warning("Wrong command")
		}
		if !ok {
			return nil
		}

		s.Players[playerNumber].Open = append(s.Players[playerNumber].Open, comm.Value)
		for _, tile := range comm.Value {
			s.Players[playerNumber].Hand.CutTile(tile)
		}
		s.Step = playerNumber // в случае анонса ход переходит к игроку, который забрал тайл
	case discardCommand:
		s.lock.Lock()
		defer s.lock.Unlock()

		log.Infof("Player #%d step", s.Step)
		if s.Step != playerNumber {
			log.Info("Wrong player number")
			return nil
		}
		if len(comm.Value) > 1 {
			log.Info("Wrong tiles number in the command")
			return nil
		}
		p := s.Players[playerNumber]
		//TODO: remove this:
		if p.CurrentTile == "" {
			log.Info("Empty current tile")
			return nil
		}
		//
		p.Hand = append(p.Hand, p.CurrentTile)
		p.CurrentTile = ""
		p.Hand.CutTile(comm.Value[0])

		p.Discard = append(p.Discard, comm.Value[0])
		// timer for announce
		s.nextTurn()
	case readyHandCommand:
		s.Players[playerNumber].IsReady = true
		//TODO: announce to all players
		log.Infof("Ready hand!")
	default:
		log.Info("Wrong client command: " + comm.Action)
		// TODO: finish with player's error
	}
	return &comm
}

func (s *statement) getFromWall() {
	currentPlayer := s.Players[s.Step]
	currentPlayer.CurrentTile = s.Wall[0]
	s.Wall = s.Wall[1:]
}

/*
// removes tile from the hand
func (h *ds.Hand) cutTile(tile string) {
	for i, elem := range *h {
		if tile == elem {
			*h = append((*h)[:i], (*h)[i+1:]...) //we use pointer to slice to avoid problems with capacity
			return
		}
	}
	// TODO: handle error
	log.Info("Tile not found")
}*/

// move turn to the next player
func (s *statement) nextTurn() {
	s.Step = (s.Step % 4) + 1
	s.getFromWall()
}

// returns last player's number
func (s *statement) prevTurn() int {
	return (s.Step+2)%4 + 1
}

// returns last tile name
func (s *statement) lastTile() string {
	pt := s.prevTurn()
	discard := s.Players[pt].Discard // только боты могут объявлять на первом ходу
	if discard != nil {
		return discard[:1][0]
	}
	return ""
}

func (p pass) checkSkip() bool {
	for _, el := range p {
		if !el {
			return false
		}
	}
	return true
}
