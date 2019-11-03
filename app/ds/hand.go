package ds

import (
	"encoding/json"
	"log"
	"sort"
	"strconv"

	"mahjong/app/apperr"
)

type Hand []string

func (h Hand) Int(i int) int {
	if i >= len(h) {
		return 0
	}
	v, err := strconv.ParseInt(string(h[i][0])+string(h[i][2]), 10, 64)
	apperr.Check(err)
	return int(v)
}

// implement sort.Interface
/*func (h Hand) Len() int {
	return len(h)
}

func (h Hand) Less(i, j int) bool {
	return h.Int(i) < h.Int(j)
}

func (h Hand) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
*/
func (h Hand) WithTile(tile string) Hand {
	return append(h, tile)
}

/*func (h Hand) CheckChow() bool {
	h2 := make(Hand, len(h))
	copy(h2, h)
	sort.Sort(h2)
	var t1, t2 int
	for i := range h2 {
		if t1+1 == t2 && t2+1 == h2.Int(i) {
			return true
		}
		t1 = t2
		t2 = h2.Int(i)
	}
	// sort
	return false
}*/

func (h Hand) FindChow() Hand {
	h2 := make(Hand, len(h))
	copy(h2, h)
	h2.SortHand()
	//sort.Sort(h2)

	var (
		t1, t2       int
		tile1, tile2 string
	)
	for i, tile := range h2 {
		if t1+1 == t2 && t2+1 == h2.Int(i) {
			return Hand{tile1, tile2, tile}
		}
		t1 = t2
		tile1 = tile2

		t2 = h2.Int(i)
		tile2 = tile
	}
	return nil
}

func (h Hand) Print() string {
	result, err := json.Marshal(h)
	apperr.Check(err)
	return string(result)
}

/*
func (h Hand) CheckPong() bool {
	//m := make(map[int]int)

	// map [string(1_2)] => count
	return false
}

func (h Hand) CheckKong() bool {
	// map [string(1_2)] => count
	return false
}*/

func (h Hand) SortHand() {
	sort.Strings(h)
}

// removes tile from the hand
func (h *Hand) CutTile(tile string) {
	if tile == "" || len(*h) == 0 {
		return
	}
	for i, elem := range *h {
		if tile == elem {
			*h = append((*h)[:i], (*h)[i+1:]...) //we use pointer to slice to avoid problems with capacity
			return
		}
	}
	// TODO: handle error
	log.Print("Tile not found")
}
