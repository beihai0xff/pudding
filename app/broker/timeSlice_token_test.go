package broker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchedule_formatTokenName(t *testing.T) {
	assert.Equal(t, "pudding_token:1", s.formatTokenName(1))
	assert.Equal(t, "pudding_token:50", s.formatTokenName(50))
	assert.Equal(t, "pudding_token:100", s.formatTokenName(100))
	assert.Equal(t, "pudding_token:10000000", s.formatTokenName(10000000))
}

func TestSchedule_parseNowFromToken(t *testing.T) {

	assert.Equal(t, uint64(10000000), s.parseNowFromToken("pudding_token:10000000"))
	assert.Equal(t, uint64(100), s.parseNowFromToken("pudding_token:100"))
	assert.Equal(t, uint64(50), s.parseNowFromToken("pudding_token:50"))
	assert.Equal(t, uint64(0), s.parseNowFromToken("pudding_token:-2"))
	assert.Equal(t, uint64(0), s.parseNowFromToken("pudding_token:wewq"))
	assert.Equal(t, uint64(0), s.parseNowFromToken("broken_token:100"))

}
