package ds

import (
	"fmt"
	"strconv"

	"mahjong/app/apperr"
)

func (h Hand) CheckPong(pong []string) bool {
	if !(len(pong) == 3) {
		return false
	}
	if !((pong[0][:4] == pong[1][:4]) && (pong[1][:4] == pong[2][:4])) {
		return false
	}
	count := 0
	for _, el := range h {
		if (el == pong[0]) || (el == pong[1]) || (el == pong[2]) {
			count++
		}
	}
	return count == 3
}

func (h Hand) CheckKong(kong []string) bool {
	if !(len(kong) == 4) {
		return false
	}
	if !((kong[0][:4] == kong[1][:4]) && (kong[1][:4] == kong[2][:4]) && (kong[2][:4] == kong[3][:4])) {
		return false
	}
	count := 0
	for _, el := range h {
		if (el == kong[0]) || (el == kong[1]) || (el == kong[2]) || (el == kong[3]) {
			count++
		}
	}
	return count == 4
}

func (h Hand) CheckChow(chow []string) bool {
	if !(len(chow) == 3) {
		return false
	}
	if !((chow[0][0] == chow[1][0]) && (chow[1][0] == chow[2][0])) {
		return false
	}

	a, err := strconv.ParseInt(chow[0][2:3], 10, 64)
	apperr.Check(err)
	b, err := strconv.ParseInt(chow[1][2:3], 10, 64)
	apperr.Check(err)
	c, err := strconv.ParseInt(chow[2][2:3], 10, 64)
	apperr.Check(err)

	if a > b {
		a, b = b, a
	}
	if a > c {
		a, c = c, a
	}

	// a is min
	expectSum := a*2 + 1 + 2
	sum := b + c
	if expectSum != sum {
		return false
	}

	//math.Abs(a - b) +math.Abs(a-c)+math.Abs(b-c)==4 && a!=b && b!=c && a!=c

	count := 0
	for _, el := range h {
		if (el == chow[0]) || (el == chow[1]) || (el == chow[2]) {
			count++
		}
	}
	fmt.Println("Chow")
	return count == 3
}

func (h Hand) CheckMahjong(mahjong [][]string) bool {

	return true
}
