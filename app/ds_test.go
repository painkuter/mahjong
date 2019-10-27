package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHand_CheckChow_false(t *testing.T) {
	h := hand([]string{"1_1_1", "2_1_2", "2_1_3", "3_7_1", "2_6_2", "2_7_3", "3_4_3",
		"4_2_1", "5_3_2", "5_2_3", "3_1_2", "1_8_2", "2_2_4", "2_9_3", "3_8_3"})
	assert.True(t, !h.CheckChow())
}

func TestHand_CheckChow_True(t *testing.T) {
	h := hand([]string{"1_1_1", "1_2_2", "2_1_3", "3_7_1", "2_6_2", "2_7_3", "3_4_3",
		"4_2_1", "5_3_2", "5_2_3", "3_1_2", "1_8_2", "2_2_4", "1_3_3", "3_8_3"})
	assert.True(t, h.CheckChow())
}
