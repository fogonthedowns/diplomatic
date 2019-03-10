package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProcessMovesBounce(t *testing.T) {
	moves := make(Moves, 0)
	moves = []*Move{
		{
			OrderType:               MOVE,
			LocationStart:           VENICE,
			LocationSubmitted:       APULIA,
			SecondLocationSubmitted: BLANK,
			UnitType:                ARMY,
		},
		{
			OrderType:               MOVE,
			LocationStart:           NAPLES,
			LocationSubmitted:       APULIA,
			SecondLocationSubmitted: BLANK,
			UnitType:                ARMY,
		},
		{
			OrderType:               SUPPORT,
			LocationStart:           ROME,
			LocationSubmitted:       VENICE,
			SecondLocationSubmitted: APULIA,
			UnitType:                ARMY,
		},
		{
			OrderType:               SUPPORT,
			LocationStart:           IONIAN_SEA,
			LocationSubmitted:       NAPLES,
			SecondLocationSubmitted: APULIA,
			UnitType:                NAVY,
		},
		{
			OrderType:         MOVE,
			LocationStart:     NORTH_SEA,
			LocationSubmitted: IONIAN_SEA,
			UnitType:          NAVY,
		},
	}

	moves.ProcessMoves()
	assert.Equal(t, VENICE, moves[0].LocationResolved)
	assert.Equal(t, NAPLES, moves[1].LocationResolved)
	assert.Equal(t, ROME, moves[2].LocationResolved)
	assert.Equal(t, IONIAN_SEA, moves[3].LocationResolved)
	assert.Equal(t, NORTH_SEA, moves[4].LocationResolved)
}

func TestProcessMovesCutSupport(t *testing.T) {
	moves := make(Moves, 0)
	moves = []*Move{
		{
			OrderType:               MOVE,
			LocationStart:           VENICE,
			LocationSubmitted:       APULIA,
			SecondLocationSubmitted: BLANK,
			UnitType:                ARMY,
		},
		{
			OrderType:               MOVE,
			LocationStart:           NAPLES,
			LocationSubmitted:       APULIA,
			SecondLocationSubmitted: BLANK,
			UnitType:                ARMY,
		},
		{
			OrderType:               SUPPORT,
			LocationStart:           ROME,
			LocationSubmitted:       VENICE,
			SecondLocationSubmitted: APULIA,
			UnitType:                ARMY,
		},
		{
			OrderType:               SUPPORT,
			LocationStart:           IONIAN_SEA,
			LocationSubmitted:       NAPLES,
			SecondLocationSubmitted: APULIA,
			UnitType:                NAVY,
		},
		{
			OrderType:         MOVE,
			LocationStart:     ADRIATIC_SEA,
			LocationSubmitted: IONIAN_SEA,
			UnitType:          NAVY,
		},
	}

	moves.ProcessMoves()
	assert.Equal(t, APULIA, moves[0].LocationResolved)
	assert.Equal(t, NAPLES, moves[1].LocationResolved)
	assert.Equal(t, ROME, moves[2].LocationResolved)
	assert.Equal(t, IONIAN_SEA, moves[3].LocationResolved)
	assert.Equal(t, ADRIATIC_SEA, moves[4].LocationResolved)
}
