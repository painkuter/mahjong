package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrevTurn(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		s := statement{
			Step: 1,
		}
		assert.Equal(t, 4, s.prevTurn())
	})

	t.Run("2", func(t *testing.T) {
		s := statement{
			Step: 2,
		}
		assert.Equal(t, 1, s.prevTurn())
	})

	t.Run("3", func(t *testing.T) {
		s := statement{
			Step: 3,
		}
		assert.Equal(t, 2, s.prevTurn())
	})

	t.Run("4", func(t *testing.T) {
		s := statement{
			Step: 4,
		}
		assert.Equal(t, 3, s.prevTurn())
	})
}
