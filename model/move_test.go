package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProcessMovesSpecialBounce(t *testing.T) {
	moves := make(Moves, 0)
	moves = []*Move{
		{
			OrderType:               MOVE,
			LocationStart:           MARSEILLES,
			LocationSubmitted:       SPAIN_SOUTH_COAST,
			SecondLocationSubmitted: BLANK,
			UnitType:                NAVY,
		},
		{
			OrderType:               MOVE,
			LocationStart:           MID_ATLANTIC_OCEAN,
			LocationSubmitted:       SPAIN_NORTH_COAST,
			SecondLocationSubmitted: BLANK,
			UnitType:                NAVY,
		},
	}

	moves.ProcessMoves()
	assert.Equal(t, MARSEILLES, moves[0].LocationResolved)
	assert.Equal(t, MID_ATLANTIC_OCEAN, moves[1].LocationResolved)

	assert.Equal(t, false, moves[0].Dislodged)
	assert.Equal(t, false, moves[1].Dislodged)
}

func TestProcessRetreatDestroy(t *testing.T) {
	moves := make(Moves, 0)
	moves = []*Move{
		{
			OrderType:         MOVE,
			LocationStart:     MARSEILLES,
			LocationSubmitted: SPAIN_SOUTH_COAST,
			UnitType:          NAVY,
			Dislodged:         true,
			DislodgedFrom:     SPAIN_SOUTH_COAST,
			IsActive:          true,
		},
	}
	moves.ProcessMoves()
	assert.Equal(t, false, moves[0].IsActive)
	assert.Equal(t, MARSEILLES, moves[0].LocationResolved)
}

func TestRetreatSuccess(t *testing.T) {
	moves := make(Moves, 0)
	moves = []*Move{
		{
			OrderType:         MOVE,
			LocationStart:     MARSEILLES,
			LocationSubmitted: GULF_OF_LYON,
			UnitType:          NAVY,
			Dislodged:         true,
			IsActive:          true,
			DislodgedFrom:     SPAIN_SOUTH_COAST,
		},
	}
	moves.ProcessMoves()
	assert.Equal(t, true, moves[0].IsActive)
	assert.Equal(t, GULF_OF_LYON, moves[0].LocationResolved)
}

func TestProcessMovesRussiaBounce(t *testing.T) {
	moves := make(Moves, 0)
	moves = []*Move{
		{
			OrderType:               MOVE,
			LocationStart:           NORWAY,
			LocationSubmitted:       ST_PETERSBURG_NORTH_COAST,
			SecondLocationSubmitted: BLANK,
			UnitType:                NAVY,
		},
		{
			OrderType:               MOVE,
			LocationStart:           MOSCOW,
			LocationSubmitted:       ST_PETERSBURG,
			SecondLocationSubmitted: BLANK,
			UnitType:                ARMY,
		},
	}

	moves.ProcessMoves()
	assert.Equal(t, NORWAY, moves[0].LocationResolved)
	assert.Equal(t, MOSCOW, moves[1].LocationResolved)

	assert.Equal(t, false, moves[0].Dislodged)
	assert.Equal(t, false, moves[1].Dislodged)
}

func TestSimpleMove(t *testing.T) {
	moves := make(Moves, 0)
	moves = []*Move{
		{
			OrderType:               MOVE,
			LocationStart:           KIEL,
			LocationSubmitted:       HOLLAND,
			SecondLocationSubmitted: BLANK,
			UnitType:                NAVY,
		},
	}

	moves.ProcessMoves()
	assert.Equal(t, HOLLAND, moves[0].LocationResolved)
	assert.Equal(t, false, moves[0].Dislodged)
}

func TestProcessMovesSpecialSupport(t *testing.T) {
	moves := make(Moves, 0)
	moves = []*Move{
		{
			OrderType:               MOVE,
			LocationStart:           MARSEILLES,
			LocationSubmitted:       SPAIN_SOUTH_COAST,
			SecondLocationSubmitted: BLANK,
			UnitType:                NAVY,
		},
		{
			OrderType:               MOVE,
			LocationStart:           MID_ATLANTIC_OCEAN,
			LocationSubmitted:       SPAIN_NORTH_COAST,
			SecondLocationSubmitted: BLANK,
			UnitType:                NAVY,
		},
		{
			OrderType:               SUPPORT,
			LocationStart:           PORTUGAL,
			LocationSubmitted:       MID_ATLANTIC_OCEAN,
			SecondLocationSubmitted: SPAIN_NORTH_COAST,
			UnitType:                ARMY,
			SupportUnitType:         NAVY,
		},
	}

	moves.ProcessMoves()
	assert.Equal(t, MARSEILLES, moves[0].LocationResolved)
	assert.Equal(t, SPAIN_NORTH_COAST, moves[1].LocationResolved)

	assert.Equal(t, false, moves[0].Dislodged)
	assert.Equal(t, false, moves[1].Dislodged)
}

func TestProcessBulgariaSupport(t *testing.T) {
	moves := make(Moves, 0)
	moves = []*Move{
		{
			OrderType:               MOVE,
			LocationStart:           CONSTANTINOPLE,
			LocationSubmitted:       BULGARIA_EAST_COAST,
			SecondLocationSubmitted: BLANK,
			UnitType:                NAVY,
		},
		{
			OrderType:               MOVE,
			LocationStart:           GREECE,
			LocationSubmitted:       BULGARIA_SOUTH_COAST,
			SecondLocationSubmitted: BLANK,
			UnitType:                NAVY,
		},
		{
			OrderType:               SUPPORT,
			LocationStart:           BLACK_SEA,
			LocationSubmitted:       CONSTANTINOPLE,
			SecondLocationSubmitted: BULGARIA_EAST_COAST,
			UnitType:                NAVY,
			SupportUnitType:         NAVY,
		},
	}

	moves.ProcessMoves()
	assert.Equal(t, BULGARIA_EAST_COAST, moves[0].LocationResolved)
	assert.Equal(t, GREECE, moves[1].LocationResolved)

	assert.Equal(t, false, moves[0].Dislodged)
	assert.Equal(t, false, moves[1].Dislodged)
}

func TestProcessBulgariaBounce(t *testing.T) {
	moves := make(Moves, 0)
	moves = []*Move{
		{
			OrderType:               MOVE,
			LocationStart:           CONSTANTINOPLE,
			LocationSubmitted:       BULGARIA_EAST_COAST,
			SecondLocationSubmitted: BLANK,
			UnitType:                NAVY,
		},
		{
			OrderType:               MOVE,
			LocationStart:           GREECE,
			LocationSubmitted:       BULGARIA_SOUTH_COAST,
			SecondLocationSubmitted: BLANK,
			UnitType:                NAVY,
		},
	}

	moves.ProcessMoves()
	assert.Equal(t, CONSTANTINOPLE, moves[0].LocationResolved)
	assert.Equal(t, GREECE, moves[1].LocationResolved)

	assert.Equal(t, false, moves[0].Dislodged)
	assert.Equal(t, false, moves[1].Dislodged)
}

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
			SupportUnitType:         ARMY,
		},
		{
			OrderType:               SUPPORT,
			LocationStart:           IONIAN_SEA,
			LocationSubmitted:       NAPLES,
			SecondLocationSubmitted: APULIA,
			UnitType:                NAVY,
			SupportUnitType:         ARMY,
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
			PieceOwner:              ITALY,
		},
		{
			OrderType:               MOVE,
			LocationStart:           NAPLES,
			LocationSubmitted:       APULIA,
			SecondLocationSubmitted: BLANK,
			UnitType:                ARMY,
			PieceOwner:              TURKEY,
		},
		{
			OrderType:               SUPPORT,
			LocationStart:           ROME,
			LocationSubmitted:       VENICE,
			SecondLocationSubmitted: APULIA,
			UnitType:                ARMY,
			SupportUnitType:         ARMY,
			PieceOwner:              ITALY,
		},
		{
			OrderType:               SUPPORT,
			LocationStart:           IONIAN_SEA,
			LocationSubmitted:       NAPLES,
			SecondLocationSubmitted: APULIA,
			UnitType:                NAVY,
			SupportUnitType:         NAVY,
			PieceOwner:              TURKEY,
		},
		{
			OrderType:         MOVE,
			LocationStart:     ADRIATIC_SEA,
			LocationSubmitted: IONIAN_SEA,
			UnitType:          NAVY,
			PieceOwner:        ITALY,
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
			PieceOwner:              ITALY,
		},
		{
			OrderType:               SUPPORT,
			LocationStart:           ROME,
			LocationSubmitted:       VENICE,
			SecondLocationSubmitted: APULIA,
			UnitType:                ARMY,
			SupportUnitType:         ARMY,
			PieceOwner:              GERMANY,
		},
		{
			OrderType:         HOLD,
			LocationStart:     APULIA,
			LocationSubmitted: APULIA,
			UnitType:          NAVY,
			PieceOwner:        TURKEY,
		},
		{
			OrderType:         MOVE,
			LocationStart:     ADRIATIC_SEA,
			LocationSubmitted: IONIAN_SEA,
			UnitType:          NAVY,
			PieceOwner:        TURKEY,
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
			PieceOwner:        ITALY,
		},
		{
			OrderType:               MOVE,
			LocationStart:           VENICE,
			LocationSubmitted:       APULIA,
			SecondLocationSubmitted: BLANK,
			UnitType:                ARMY,
			PieceOwner:              RUSSIA,
		},
		{
			OrderType:               SUPPORT,
			LocationStart:           ROME,
			LocationSubmitted:       VENICE,
			SecondLocationSubmitted: APULIA,
			UnitType:                ARMY,
			SupportUnitType:         ARMY,
			PieceOwner:              RUSSIA,
		},

		{
			OrderType:         MOVE,
			LocationStart:     ADRIATIC_SEA,
			LocationSubmitted: IONIAN_SEA,
			UnitType:          NAVY,
			PieceOwner:        ITALY,
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

func TestProcessIllegalSupportHold(t *testing.T) {
	moves := make(Moves, 0)
	moves = []*Move{
		{
			OrderType:         HOLD,
			LocationStart:     VIENNA,
			LocationSubmitted: VIENNA,
			UnitType:          ARMY,
			PieceOwner:        ITALY,
		},
		{
			OrderType:               SUPPORT,
			LocationStart:           MOSCOW,
			LocationSubmitted:       VIENNA,
			SecondLocationSubmitted: VIENNA,
			UnitType:                ARMY,
			PieceOwner:              ITALY,
			SupportUnitType:         ARMY,
		},
		{
			OrderType:         MOVE,
			LocationStart:     BOHEMIA,
			LocationSubmitted: VIENNA,
			UnitType:          ARMY,
			PieceOwner:        RUSSIA,
		},

		{
			OrderType:               SUPPORT,
			LocationStart:           GALICIA,
			LocationSubmitted:       BOHEMIA,
			SecondLocationSubmitted: VIENNA,
			UnitType:                ARMY,
			PieceOwner:              RUSSIA,
			SupportUnitType:         ARMY,
		},
	}

	moves.ProcessMoves()
	assert.Equal(t, VIENNA, moves[0].LocationResolved)
	assert.Equal(t, MOSCOW, moves[1].LocationResolved)
	assert.Equal(t, VIENNA, moves[2].LocationResolved)
	assert.Equal(t, GALICIA, moves[3].LocationResolved)

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
			PieceOwner:        RUSSIA,
		},
		{
			OrderType:               SUPPORT,
			LocationStart:           TYROLIA,
			LocationSubmitted:       VIENNA,
			SecondLocationSubmitted: VIENNA,
			UnitType:                ARMY,
			SupportUnitType:         ARMY,
			PieceOwner:              RUSSIA,
		},
		{
			OrderType:         MOVE,
			LocationStart:     BOHEMIA,
			LocationSubmitted: VIENNA,
			UnitType:          ARMY,
			PieceOwner:        TURKEY,
		},

		{
			OrderType:               SUPPORT,
			LocationStart:           GALICIA,
			LocationSubmitted:       BOHEMIA,
			SecondLocationSubmitted: VIENNA,
			UnitType:                ARMY,
			SupportUnitType:         ARMY,
			PieceOwner:              TURKEY,
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

func TestInvalidProcessMovesInvalidMoveViaConvoy(t *testing.T) {
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

func TestInvalidProcessMovesMoveViaConvoy(t *testing.T) {
	moves := make(Moves, 0)
	moves = []*Move{
		{
			OrderType:               CONVOY,
			LocationStart:           MID_ATLANTIC_OCEAN,
			LocationSubmitted:       LONDON,
			SecondLocationSubmitted: TUNIS,
			UnitType:                NAVY,
		},
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
	assert.Equal(t, MID_ATLANTIC_OCEAN, moves[0].LocationResolved)
	assert.Equal(t, TUNIS, moves[1].LocationResolved)
	assert.Equal(t, ENGLISH_CHANNEL, moves[2].LocationResolved)
	assert.Equal(t, WESTERN_MEDITERRANEAN, moves[3].LocationResolved)

	assert.Equal(t, false, moves[0].Dislodged)
	assert.Equal(t, false, moves[1].Dislodged)
	assert.Equal(t, false, moves[2].Dislodged)
	assert.Equal(t, false, moves[3].Dislodged)
}

func TestLongPathProcessTwoConvoys(t *testing.T) {
	moves := make(Moves, 0)
	moves = []*Move{
		{
			OrderType:               CONVOY,
			LocationStart:           MID_ATLANTIC_OCEAN,
			LocationSubmitted:       LONDON,
			SecondLocationSubmitted: TUNIS,
			UnitType:                NAVY,
		},
		{
			OrderType:         MOVEVIACONVOY,
			LocationStart:     LONDON,
			LocationSubmitted: TUNIS,
			UnitType:          ARMY,
		},

		{
			OrderType:               CONVOY,
			LocationStart:           ADRIATIC_SEA,
			LocationSubmitted:       ALBANIA,
			SecondLocationSubmitted: APULIA,
			UnitType:                NAVY,
		},
		{
			OrderType:         MOVEVIACONVOY,
			LocationStart:     ALBANIA,
			LocationSubmitted: APULIA,
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
	assert.Equal(t, MID_ATLANTIC_OCEAN, moves[0].LocationResolved)
	assert.Equal(t, TUNIS, moves[1].LocationResolved)
	assert.Equal(t, ENGLISH_CHANNEL, moves[4].LocationResolved)
	assert.Equal(t, WESTERN_MEDITERRANEAN, moves[5].LocationResolved)

	assert.Equal(t, false, moves[0].Dislodged)
	assert.Equal(t, false, moves[1].Dislodged)
	assert.Equal(t, false, moves[4].Dislodged)
	assert.Equal(t, false, moves[5].Dislodged)

	assert.Equal(t, ADRIATIC_SEA, moves[2].LocationResolved)
	assert.Equal(t, APULIA, moves[3].LocationResolved)

	assert.Equal(t, false, moves[2].Dislodged)
	assert.Equal(t, false, moves[3].Dislodged)
}

func TestProcessDislodgeConvoy(t *testing.T) {
	moves := make(Moves, 0)
	moves = []*Move{
		{
			Id:                1,
			OrderType:         MOVEVIACONVOY,
			LocationStart:     LONDON,
			LocationSubmitted: PICARDY,
			UnitType:          ARMY,
			PieceOwner:        ENGLAND,
		},
		{
			Id:                      2,
			OrderType:               CONVOY,
			LocationStart:           ENGLISH_CHANNEL,
			LocationSubmitted:       LONDON,
			SecondLocationSubmitted: PICARDY,
			UnitType:                NAVY,
			PieceOwner:              GERMANY,
		},
		{
			Id:                3,
			OrderType:         MOVE,
			LocationStart:     NORTH_SEA,
			LocationSubmitted: ENGLISH_CHANNEL,
			UnitType:          NAVY,
			PieceOwner:        FRANCE,
		},
		{
			Id:                      4,
			OrderType:               SUPPORT,
			LocationStart:           IRISH_SEA,
			LocationSubmitted:       NORTH_SEA,
			SecondLocationSubmitted: ENGLISH_CHANNEL,
			UnitType:                NAVY,
			PieceOwner:              FRANCE,
		},
	}

	moves.ProcessMoves()
	assert.Equal(t, LONDON, moves[0].LocationResolved)
	assert.Equal(t, ENGLISH_CHANNEL, moves[1].LocationResolved)

	assert.Equal(t, false, moves[0].Dislodged)
	assert.Equal(t, true, moves[1].Dislodged)
	assert.Equal(t, false, moves[2].Dislodged)
	assert.Equal(t, false, moves[3].Dislodged)
	assert.Equal(t, NORTH_SEA, moves[1].DislodgedFrom)
}

func TestLongPathProcessTwoDislodgedConvoys(t *testing.T) {
	moves := make(Moves, 0)
	moves = []*Move{
		{
			Id:                      0,
			OrderType:               CONVOY,
			LocationStart:           MID_ATLANTIC_OCEAN,
			LocationSubmitted:       LONDON,
			SecondLocationSubmitted: TUNIS,
			UnitType:                NAVY,
			PieceOwner:              ENGLAND,
		},
		{
			Id:                1,
			OrderType:         MOVEVIACONVOY,
			LocationStart:     LONDON,
			LocationSubmitted: TUNIS,
			UnitType:          ARMY,
			PieceOwner:        ENGLAND,
		},

		{
			Id:                      2,
			OrderType:               CONVOY,
			LocationStart:           ADRIATIC_SEA,
			LocationSubmitted:       ALBANIA,
			SecondLocationSubmitted: APULIA,
			UnitType:                NAVY,
			PieceOwner:              TURKEY,
		},
		{
			Id:                3,
			OrderType:         MOVEVIACONVOY,
			LocationStart:     ALBANIA,
			LocationSubmitted: APULIA,
			UnitType:          ARMY,
			PieceOwner:        TURKEY,
		},

		{
			Id:                      4,
			OrderType:               CONVOY,
			LocationStart:           ENGLISH_CHANNEL,
			LocationSubmitted:       LONDON,
			SecondLocationSubmitted: TUNIS,
			UnitType:                NAVY,
			PieceOwner:              ENGLAND,
		},
		{
			Id:                      5,
			OrderType:               CONVOY,
			LocationStart:           WESTERN_MEDITERRANEAN,
			LocationSubmitted:       LONDON,
			SecondLocationSubmitted: TUNIS,
			UnitType:                NAVY,
			PieceOwner:              ENGLAND,
		},
		{
			Id:                6,
			OrderType:         MOVE,
			LocationStart:     NORTH_SEA,
			LocationSubmitted: ENGLISH_CHANNEL,
			UnitType:          NAVY,
			PieceOwner:        FRANCE,
		},
		{
			Id:                      7,
			OrderType:               SUPPORT,
			LocationStart:           IRISH_SEA,
			LocationSubmitted:       NORTH_SEA,
			SecondLocationSubmitted: ENGLISH_CHANNEL,
			UnitType:                NAVY,
			PieceOwner:              FRANCE,
			SupportUnitType:         NAVY,
		},
		{
			Id:                8,
			OrderType:         MOVE,
			LocationStart:     TYRRHENIAN_SEA,
			LocationSubmitted: WESTERN_MEDITERRANEAN,
			UnitType:          NAVY,
			PieceOwner:        FRANCE,
		},
		{
			Id:                      9,
			OrderType:               SUPPORT,
			LocationStart:           GULF_OF_LYON,
			LocationSubmitted:       TYRRHENIAN_SEA,
			SecondLocationSubmitted: WESTERN_MEDITERRANEAN,
			UnitType:                NAVY,
			PieceOwner:              FRANCE,
			SupportUnitType:         NAVY,
		},
		{
			Id:                10,
			OrderType:         MOVE,
			LocationStart:     IONIAN_SEA,
			LocationSubmitted: ADRIATIC_SEA,
			UnitType:          NAVY,
			PieceOwner:        ITALY,
		},
		{
			Id:                      11,
			OrderType:               SUPPORT,
			LocationStart:           VENICE,
			LocationSubmitted:       IONIAN_SEA,
			SecondLocationSubmitted: ADRIATIC_SEA,
			UnitType:                NAVY,
			PieceOwner:              ITALY,
			SupportUnitType:         NAVY,
		},
	}

	moves.ProcessMoves()
	assert.Equal(t, MID_ATLANTIC_OCEAN, moves[0].LocationResolved)
	assert.Equal(t, LONDON, moves[1].LocationResolved)
	assert.Equal(t, ENGLISH_CHANNEL, moves[4].LocationResolved)
	assert.Equal(t, WESTERN_MEDITERRANEAN, moves[5].LocationResolved)

	assert.Equal(t, false, moves[0].Dislodged)
	assert.Equal(t, false, moves[1].Dislodged)
	assert.Equal(t, true, moves[4].Dislodged)
	assert.Equal(t, true, moves[5].Dislodged)

	assert.Equal(t, ADRIATIC_SEA, moves[2].LocationResolved)
	assert.Equal(t, ALBANIA, moves[3].LocationResolved)

	assert.Equal(t, true, moves[2].Dislodged)
	assert.Equal(t, false, moves[3].Dislodged)
	assert.Equal(t, IONIAN_SEA, moves[2].DislodgedFrom)
}
