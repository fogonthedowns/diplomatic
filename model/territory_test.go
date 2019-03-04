package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolveConflictsBounce(t *testing.T) {
	tm := make(TerritoryMoves, 0)
	tm = TerritoryMoves{
		APULIA: []*Move{
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
		},
	}
	tm.ResolveConflicts()
	assert.Equal(t, VENICE, tm[APULIA][0].LocationResolved)
	assert.Equal(t, NAPLES, tm[APULIA][1].LocationResolved)
}

func TestResolveConflictsSupport(t *testing.T) {
	tm := make(TerritoryMoves, 0)
	tm = TerritoryMoves{
		APULIA: []*Move{
			{
				OrderType:               MOVE,
				LocationStart:           VENICE,
				LocationSubmitted:       APULIA,
				SecondLocationSubmitted: BLANK,
				MovePower:               1,
			},
			{
				OrderType:               MOVE,
				LocationStart:           NAPLES,
				LocationSubmitted:       APULIA,
				SecondLocationSubmitted: BLANK,
				MovePower:               0,
			},
		},
		ROME: []*Move{
			{
				OrderType:               SUPPORT,
				LocationStart:           ROME,
				LocationSubmitted:       VENICE,
				SecondLocationSubmitted: APULIA,
			},
		},
	}
	tm.ResolveConflicts()
	assert.Equal(t, APULIA, tm[APULIA][0].LocationResolved)
	assert.Equal(t, NAPLES, tm[APULIA][1].LocationResolved)
}
