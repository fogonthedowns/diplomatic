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

	assert.Equal(t, false, moves[0].Dislodged)
	assert.Equal(t, false, moves[1].Dislodged)
	assert.Equal(t, false, moves[2].Dislodged)
	assert.Equal(t, false, moves[3].Dislodged)
	assert.Equal(t, false, moves[4].Dislodged)
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

	assert.Equal(t, false, moves[0].Dislodged)
	assert.Equal(t, false, moves[1].Dislodged)
	assert.Equal(t, false, moves[2].Dislodged)
	assert.Equal(t, false, moves[3].Dislodged)
	assert.Equal(t, false, moves[4].Dislodged)
}

func TestProcessMovesDislodge(t *testing.T) {
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
			OrderType:               SUPPORT,
			LocationStart:           ROME,
			LocationSubmitted:       VENICE,
			SecondLocationSubmitted: APULIA,
			UnitType:                ARMY,
		},
		{
			OrderType:         HOLD,
			LocationStart:     APULIA,
			LocationSubmitted: APULIA,
			UnitType:          NAVY,
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
	assert.Equal(t, ROME, moves[1].LocationResolved)
	assert.Equal(t, APULIA, moves[2].LocationResolved)
	assert.Equal(t, IONIAN_SEA, moves[3].LocationResolved)

	assert.Equal(t, false, moves[0].Dislodged)
	assert.Equal(t, false, moves[1].Dislodged)
	assert.Equal(t, true, moves[2].Dislodged)
	assert.Equal(t, false, moves[3].Dislodged)
}

func TestProcessMovesDislodgeReorderMoves(t *testing.T) {
	moves := make(Moves, 0)
	moves = []*Move{
		{
			OrderType:         HOLD,
			LocationStart:     APULIA,
			LocationSubmitted: APULIA,
			UnitType:          NAVY,
		},
		{
			OrderType:               MOVE,
			LocationStart:           VENICE,
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
			OrderType:         MOVE,
			LocationStart:     ADRIATIC_SEA,
			LocationSubmitted: IONIAN_SEA,
			UnitType:          NAVY,
		},
	}

	moves.ProcessMoves()
	assert.Equal(t, APULIA, moves[0].LocationResolved)
	assert.Equal(t, APULIA, moves[1].LocationResolved)
	assert.Equal(t, ROME, moves[2].LocationResolved)
	assert.Equal(t, IONIAN_SEA, moves[3].LocationResolved)

	assert.Equal(t, true, moves[0].Dislodged)
	assert.Equal(t, false, moves[1].Dislodged)
	assert.Equal(t, false, moves[2].Dislodged)
	assert.Equal(t, false, moves[3].Dislodged)
}

func TestProcessSupportHold(t *testing.T) {
	moves := make(Moves, 0)
	moves = []*Move{
		{
			OrderType:         HOLD,
			LocationStart:     VIENNA,
			LocationSubmitted: VIENNA,
			UnitType:          ARMY,
		},
		{
			OrderType:               SUPPORT,
			LocationStart:           TYROLIA,
			LocationSubmitted:       VIENNA,
			SecondLocationSubmitted: VIENNA,
			UnitType:                ARMY,
		},
		{
			OrderType:         MOVE,
			LocationStart:     BOHEMIA,
			LocationSubmitted: VIENNA,
			UnitType:          ARMY,
		},

		{
			OrderType:               SUPPORT,
			LocationStart:           GALICIA,
			LocationSubmitted:       BOHEMIA,
			SecondLocationSubmitted: VIENNA,
			UnitType:                ARMY,
		},
	}

	moves.ProcessMoves()
	assert.Equal(t, VIENNA, moves[0].LocationResolved)
	assert.Equal(t, TYROLIA, moves[1].LocationResolved)
	assert.Equal(t, BOHEMIA, moves[2].LocationResolved)
	assert.Equal(t, GALICIA, moves[3].LocationResolved)

	assert.Equal(t, false, moves[0].Dislodged)
	assert.Equal(t, false, moves[1].Dislodged)
	assert.Equal(t, false, moves[2].Dislodged)
	assert.Equal(t, false, moves[3].Dislodged)
}

func TestProcessMovesMoveViaConvoy(t *testing.T) {
	moves := make(Moves, 0)
	moves = []*Move{
		{
			OrderType:         MOVEVIACONVOY,
			LocationStart:     LONDON,
			LocationSubmitted: PICARDY,
			UnitType:          ARMY,
		},
		{
			OrderType:               CONVOY,
			LocationStart:           ENGLISH_CHANNEL,
			LocationSubmitted:       LONDON,
			SecondLocationSubmitted: PICARDY,
			UnitType:                NAVY,
		},
	}

	moves.ProcessMoves()
	assert.Equal(t, PICARDY, moves[0].LocationResolved)
	assert.Equal(t, ENGLISH_CHANNEL, moves[1].LocationResolved)

	assert.Equal(t, false, moves[0].Dislodged)
	assert.Equal(t, false, moves[1].Dislodged)
}

func TestInvalidProcessMovesMoveViaConvoy(t *testing.T) {
	moves := make(Moves, 0)
	moves = []*Move{
		{
			OrderType:         MOVEVIACONVOY,
			LocationStart:     LONDON,
			LocationSubmitted: TUNIS,
			UnitType:          ARMY,
		},
		{
			OrderType:               CONVOY,
			LocationStart:           ENGLISH_CHANNEL,
			LocationSubmitted:       LONDON,
			SecondLocationSubmitted: TUNIS,
			UnitType:                NAVY,
		},
		{
			OrderType:               CONVOY,
			LocationStart:           WESTERN_MEDITERRANEAN,
			LocationSubmitted:       LONDON,
			SecondLocationSubmitted: TUNIS,
			UnitType:                NAVY,
		},
	}

	moves.ProcessMoves()
	assert.Equal(t, LONDON, moves[0].LocationResolved)
	assert.Equal(t, ENGLISH_CHANNEL, moves[1].LocationResolved)
	assert.Equal(t, WESTERN_MEDITERRANEAN, moves[2].LocationResolved)

	assert.Equal(t, false, moves[0].Dislodged)
	assert.Equal(t, false, moves[1].Dislodged)
	assert.Equal(t, false, moves[2].Dislodged)
}

// TODO write test from LON to TUN, with support only in Western Med and Eng channel missing north atlantic.
