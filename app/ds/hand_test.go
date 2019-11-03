package ds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHand_CheckChow_false(t *testing.T) {
	h := Hand([]string{"1_1_1", "2_1_2", "2_1_3", "3_7_1", "2_6_2", "2_7_3", "3_4_3",
		"4_2_1", "5_3_2", "5_2_3", "3_1_2", "1_8_2", "2_2_4", "2_9_3", "3_8_3"})
	assert.Nil(t, h.FindChow())
}

func TestHand_CheckChow_True(t *testing.T) {
	h := Hand([]string{"1_1_1", "1_2_2", "2_1_3", "3_7_1", "2_6_2", "2_7_3", "3_4_3",
		"4_2_1", "5_3_2", "5_2_3", "3_1_2", "1_8_2", "2_2_4", "1_3_3", "3_8_3"})
	assert.NotNil(t, h.FindChow())
}

func TestFindPong(t *testing.T) {
	a := "1_1_1"
	b := "1_1_2"
	c := "1_1_3"
	h := Hand{"1_1_1", "1_2_3", "1_1_2", "1_1_3"}
	assert.Equal(t, true, h.CheckPong([]string{a, b, c}))
}

func TestFindPong_NotFound(t *testing.T) {
	a := "1_1_1"
	b := "1_1_2"
	c := "1_1_3"
	h := Hand{"1_1_1", "1_2_3", "1_1_3"}
	assert.Equal(t, false, h.CheckPong([]string{a, b, c}))
}

func TestFindKong(t *testing.T) {
	a := "1_1_1"
	b := "1_1_2"
	c := "1_1_3"
	d := "1_1_4"
	h := Hand{"1_1_1", "1_2_3", "1_1_2", "1_1_3", "1_1_4"}
	assert.Equal(t, true, h.CheckKong([]string{a, b, c, d}))
}

func TestFindKong_NotFound(t *testing.T) {
	a := "1_1_1"
	b := "1_1_2"
	c := "1_1_3"
	d := "1_1_4"
	h := Hand{"1_1_1", "1_2_3", "1_1_3", "1_1_4"}
	assert.Equal(t, false, h.CheckKong([]string{a, b, c, d}))
}

func TestCutTile(t *testing.T) {
	h := Hand{"1_1_1", "1_9_3", "1_2_3", "1_1_3", "1_1_4"}
	h.CutTile("1_9_3")
	assert.Equal(t, Hand{"1_1_1", "1_2_3", "1_1_3", "1_1_4"}, h)
}
