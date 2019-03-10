package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddSupportPointsToMove(t *testing.T) {
	moves := make(Moves, 0)
	moves = []*Move{
		{
			OrderType:               MOVE,
			LocationStart:           VENICE,
			LocationSubmitted:       APULIA,
			SecondLocationSubmitted: BLANK,
		},
		{
			OrderType:               MOVE,
			LocationStart:           NAPLES,
			LocationSubmitted:       APULIA,
			SecondLocationSubmitted: BLANK,
		},
		{
			OrderType:               SUPPORT,
			LocationStart:           ROME,
			LocationSubmitted:       VENICE,
			SecondLocationSubmitted: APULIA,
		},
	}

	supportMove := Move{
		OrderType:               SUPPORT,
		LocationStart:           ROME,
		LocationSubmitted:       VENICE,
		SecondLocationSubmitted: APULIA,
	}

	moves.AddSupportPointsToMove(supportMove)
	assert.Equal(t, 1, moves[0].MovePower)
	assert.Equal(t, 0, moves[1].MovePower)
}

func TestAddSupportPointsToMoveTwoX(t *testing.T) {
	moves := make(Moves, 0)
	moves = []*Move{
		{
			OrderType:               MOVE,
			LocationStart:           VENICE,
			LocationSubmitted:       APULIA,
			SecondLocationSubmitted: BLANK,
		},
		{
			OrderType:               MOVE,
			LocationStart:           NAPLES,
			LocationSubmitted:       APULIA,
			SecondLocationSubmitted: BLANK,
		},
		{
			OrderType:               SUPPORT,
			LocationStart:           ROME,
			LocationSubmitted:       VENICE,
			SecondLocationSubmitted: APULIA,
		},
		{
			OrderType:               SUPPORT,
			LocationStart:           ADRIATIC_SEA,
			LocationSubmitted:       VENICE,
			SecondLocationSubmitted: APULIA,
		},
	}

	moves.CalculateSupport()
	assert.Equal(t, 2, moves[0].MovePower)
	assert.Equal(t, 0, moves[1].MovePower)
}

func TestAddSupportPointsToMoveCutSupport(t *testing.T) {
	moves := make(Moves, 0)
	moves = []*Move{
		{
			OrderType:               MOVE,
			LocationStart:           VENICE,
			LocationSubmitted:       APULIA,
			SecondLocationSubmitted: BLANK,
		},
		{
			OrderType:               MOVE,
			LocationStart:           NAPLES,
			LocationSubmitted:       APULIA,
			SecondLocationSubmitted: BLANK,
		},
		{
			OrderType:               SUPPORT,
			LocationStart:           ROME,
			LocationSubmitted:       VENICE,
			SecondLocationSubmitted: APULIA,
		},
		{
			OrderType:               SUPPORT,
			LocationStart:           ADRIATIC_SEA,
			LocationSubmitted:       VENICE,
			SecondLocationSubmitted: APULIA,
			UnitType:                NAVY,
		},
		{
			OrderType:         MOVE,
			LocationStart:     IONIAN_SEA,
			LocationSubmitted: ADRIATIC_SEA,
			UnitType:          NAVY,
		},
	}

	moves.CalculateSupport()
	assert.Equal(t, 1, moves[0].MovePower)
	assert.Equal(t, 0, moves[1].MovePower)
}

func TestAddSupportPointsToMoveDoubleSupport(t *testing.T) {
	moves := make(Moves, 0)
	moves = []*Move{
		{
			OrderType:               MOVE,
			LocationStart:           VENICE,
			LocationSubmitted:       APULIA,
			SecondLocationSubmitted: BLANK,
		},
		{
			OrderType:               MOVE,
			LocationStart:           NAPLES,
			LocationSubmitted:       APULIA,
			SecondLocationSubmitted: BLANK,
		},
		{
			OrderType:               SUPPORT,
			LocationStart:           ROME,
			LocationSubmitted:       VENICE,
			SecondLocationSubmitted: APULIA,
		},
		{
			OrderType:               SUPPORT,
			LocationStart:           ADRIATIC_SEA,
			LocationSubmitted:       VENICE,
			SecondLocationSubmitted: APULIA,
			UnitType:                NAVY,
		},
		{
			OrderType:         MOVE,
			LocationStart:     NORTH_SEA,
			LocationSubmitted: ADRIATIC_SEA,
			UnitType:          NAVY,
		},
	}

	moves.CalculateSupport()
	assert.Equal(t, 2, moves[0].MovePower)
	assert.Equal(t, 0, moves[1].MovePower)
}
