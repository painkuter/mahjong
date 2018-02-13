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
	assert.Equal(t, h.findPong(a,b,c), true)
}
