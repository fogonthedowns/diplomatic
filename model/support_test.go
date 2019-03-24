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
			PieceOwner:              ITALY,
		},
		{
			OrderType:               MOVE,
			LocationStart:           NAPLES,
			LocationSubmitted:       APULIA,
			SecondLocationSubmitted: BLANK,
			PieceOwner:              ENGLAND,
		},
		{
			OrderType:               SUPPORT,
			LocationStart:           ROME,
			LocationSubmitted:       VENICE,
			SecondLocationSubmitted: APULIA,
			PieceOwner:              ITALY,
		},
	}

	supportMove := Move{
		OrderType:               SUPPORT,
		LocationStart:           ROME,
		LocationSubmitted:       VENICE,
		SecondLocationSubmitted: APULIA,
		PieceOwner:              ITALY,
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
			PieceOwner:              RUSSIA,
		},
		{
			OrderType:               MOVE,
			LocationStart:           NAPLES,
			LocationSubmitted:       APULIA,
			SecondLocationSubmitted: BLANK,
			PieceOwner:              ITALY,
		},
		{
			OrderType:               SUPPORT,
			LocationStart:           ROME,
			LocationSubmitted:       VENICE,
			SecondLocationSubmitted: APULIA,
			PieceOwner:              RUSSIA,
		},
		{
			OrderType:               SUPPORT,
			LocationStart:           ADRIATIC_SEA,
			LocationSubmitted:       VENICE,
			SecondLocationSubmitted: APULIA,
			PieceOwner:              ITALY,
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
			PieceOwner:              ITALY,
		},
		{
			OrderType:               MOVE,
			LocationStart:           NAPLES,
			LocationSubmitted:       APULIA,
			SecondLocationSubmitted: BLANK,
			PieceOwner:              TURKEY,
		},
		{
			OrderType:               SUPPORT,
			LocationStart:           ROME,
			LocationSubmitted:       VENICE,
			SecondLocationSubmitted: APULIA,
			PieceOwner:              GERMANY,
		},
		{
			OrderType:               SUPPORT,
			LocationStart:           ADRIATIC_SEA,
			LocationSubmitted:       VENICE,
			SecondLocationSubmitted: APULIA,
			UnitType:                NAVY,
			PieceOwner:              TURKEY,
		},
		{
			OrderType:         MOVE,
			LocationStart:     IONIAN_SEA,
			LocationSubmitted: ADRIATIC_SEA,
			UnitType:          NAVY,
			PieceOwner:        ITALY,
		},
	}

	moves.CalculateSupport()
	assert.Equal(t, 1, moves[0].MovePower)
	assert.Equal(t, 0, moves[1].MovePower)
}

func TestCanNotCutSupportOfSelf(t *testing.T) {
	moves := make(Moves, 0)
	moves = []*Move{
		{
			OrderType:               MOVE,
			LocationStart:           VENICE,
			LocationSubmitted:       APULIA,
			SecondLocationSubmitted: BLANK,
			PieceOwner:              ITALY,
		},
		{
			OrderType:               MOVE,
			LocationStart:           NAPLES,
			LocationSubmitted:       APULIA,
			SecondLocationSubmitted: BLANK,
			PieceOwner:              TURKEY,
		},
		{
			OrderType:               SUPPORT,
			LocationStart:           ROME,
			LocationSubmitted:       VENICE,
			SecondLocationSubmitted: APULIA,
			PieceOwner:              GERMANY,
		},
		{
			OrderType:               SUPPORT,
			LocationStart:           ADRIATIC_SEA,
			LocationSubmitted:       VENICE,
			SecondLocationSubmitted: APULIA,
			UnitType:                NAVY,
			PieceOwner:              TURKEY,
		},
		{
			OrderType:         MOVE,
			LocationStart:     IONIAN_SEA,
			LocationSubmitted: ADRIATIC_SEA,
			UnitType:          NAVY,
			PieceOwner:        TURKEY,
		},
	}

	moves.CalculateSupport()
	assert.Equal(t, 2, moves[0].MovePower)
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
			UnitType:                ARMY,
			PieceOwner:              TURKEY,
		},
		{
			OrderType:               MOVE,
			LocationStart:           NAPLES,
			LocationSubmitted:       APULIA,
			SecondLocationSubmitted: BLANK,
			UnitType:                ARMY,
			PieceOwner:              ITALY,
		},
		{
			OrderType:               SUPPORT,
			LocationStart:           ROME,
			LocationSubmitted:       VENICE,
			SecondLocationSubmitted: APULIA,
			UnitType:                ARMY,
			PieceOwner:              RUSSIA,
		},
		{
			OrderType:               SUPPORT,
			LocationStart:           ADRIATIC_SEA,
			LocationSubmitted:       VENICE,
			SecondLocationSubmitted: APULIA,
			UnitType:                NAVY,
			PieceOwner:              ITALY,
		},
		{
			OrderType:         MOVE,
			LocationStart:     NORTH_SEA,
			LocationSubmitted: ADRIATIC_SEA,
			UnitType:          NAVY,
			PieceOwner:        ENGLAND,
		},
	}

	moves.CalculateSupport()
	assert.Equal(t, 2, moves[0].MovePower)
	assert.Equal(t, 0, moves[1].MovePower)
}
