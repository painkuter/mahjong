package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindPong(t *testing.T) {
	a := "1_1_1"
	b := "1_1_2"
	c := "1_1_3"
	h := ds.Hand{"1_1_1", "1_2_3", "1_1_2", "1_1_3"}
	assert.Equal(t, true, h.checkPong([]string{a, b, c}))
}

func TestFindPong_NotFound(t *testing.T) {
	a := "1_1_1"
	b := "1_1_2"
	c := "1_1_3"
	h := ds.Hand{"1_1_1", "1_2_3", "1_1_3"}
	assert.Equal(t, false, h.checkPong([]string{a, b, c}))
}

func TestFindKong(t *testing.T) {
	a := "1_1_1"
	b := "1_1_2"
	c := "1_1_3"
	d := "1_1_4"
	h := ds.Hand{"1_1_1", "1_2_3", "1_1_2", "1_1_3", "1_1_4"}
	assert.Equal(t, true, h.checkKong([]string{a, b, c, d}))
}

func TestFindKong_NotFound(t *testing.T) {
	a := "1_1_1"
	b := "1_1_2"
	c := "1_1_3"
	d := "1_1_4"
	h := hand{"1_1_1", "1_2_3", "1_1_3", "1_1_4"}
	assert.Equal(t, false, h.checkKong([]string{a, b, c, d}))
}

func TestCutTile(t *testing.T) {
	h := hand{"1_1_1", "1_9_3", "1_2_3", "1_1_3", "1_1_4"}
	h.cutTile("1_9_3")
	assert.Equal(t, hand{"1_1_1", "1_2_3", "1_1_3", "1_1_4"}, h)
}
