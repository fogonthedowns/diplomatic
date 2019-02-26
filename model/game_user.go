package model

import "errors"

type GameUser struct {
	UserId  int     `json:"user_id"`
	GameId  int     `json:"game_id"`
	Country Country `json:"country"`
}

func Validate(gameusers []GameUser, c Country) error {
	if len(gameusers) == 7 {
		return errors.New("The Selected Game is Full")
	}

	for _, game := range gameusers {
		if game.Country == c {
			return errors.New("The Country is already Selected")
		}
	}
	return nil
}
