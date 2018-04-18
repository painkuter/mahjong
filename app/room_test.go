package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestWall(t *testing.T) {
	r := newRoom()
	fmt.Println(len(r.statement.Wall))
	assert.Equal(t, len(r.statement.Wall) == wallSize - 4*handSize - reserveSize, true)
}
