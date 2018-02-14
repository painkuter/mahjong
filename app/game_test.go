package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindPong(t *testing.T) {
	a := "1_1_1"
	b := "1_1_2"
	c := "1_1_3"
	h := hand{"1_1_1", "1_2_3", "1_1_2", "1_1_3"}
	assert.Equal(t, h.findPong([]string{a, b, c}), true)
}

func TestFindPong_NotFound(t *testing.T) {
	a := "1_1_1"
	b := "1_1_2"
	c := "1_1_3"
	h := hand{"1_1_1", "1_2_3", "1_1_3"}
	assert.Equal(t, h.findPong([]string{a, b, c}), false)
}

func TestFindKong(t *testing.T) {
	a := "1_1_1"
	b := "1_1_2"
	c := "1_1_3"
	d := "1_1_4"
	h := hand{"1_1_1", "1_2_3", "1_1_2", "1_1_3", "1_1_4"}
	assert.Equal(t, h.findKong([]string{a, b, c, d}), true)
}

func TestFindKong_NotFound(t *testing.T) {
	a := "1_1_1"
	b := "1_1_2"
	c := "1_1_3"
	d := "1_1_4"
	h := hand{"1_1_1", "1_2_3", "1_1_3", "1_1_4"}
	assert.Equal(t, h.findKong([]string{a, b, c, d}), false)
}
