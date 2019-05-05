package model

import (
	"errors"
	"strconv"
	"time"
)

type GamePhase int

const (
	Waiting       = GamePhase(0)
	Spring        = GamePhase(1)
	SpringRetreat = GamePhase(2)
	Fall          = GamePhase(3)
	FallRetreat   = GamePhase(4)
	FallBuild     = GamePhase(5)
)

func NewPhase(phase GamePhase) GamePhase {
	switch phase {
	case Waiting:
		return Spring
	case Spring:
		return SpringRetreat
	case SpringRetreat:
		return Fall
	case Fall:
		return FallRetreat
	case FallRetreat:
		return FallBuild
	case FallBuild:
		return Spring
	default:
		panic("phase not present int changePhase()")
	}
}

func (g GamePhase) ValidatePhaseUpdate(updateToPhase GamePhase) error {
	if g >= updateToPhase {
		return errors.New("The Phase has already been updated")
	}
	return nil
}

// HasPhaseEnded() returns true if the current phase is over
func (g GamePhase) HasPhaseEnded(endTime string) (hasEnded bool) {
	// if the game has not started yet return false
	if g < Waiting {
		return false
	}
	now := time.Now()
	timestamp, err := strconv.ParseInt(endTime, 10, 64)
	if err != nil {
		panic(err)
	}
	phaseEndTime := time.Unix(timestamp, 0)
	hasEnded = phaseEndTime.Before(now)

	return hasEnded
}
