package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWall(t *testing.T) {
	r := newRoom()
	assert.Equal(t, len(r.statement.Wall) == wallSize, true)
}
