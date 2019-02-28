package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	// "strconv"
	// "strings"

	// "github.com/go-chi/chi"

	// db "diplomacy/db"
	gamecrud "diplomacy/db/gamecrud"
	"diplomacy/driver"
	model "diplomacy/model"
)

// NewGameHandler() creates a new HTTP handler
// To handle web requests for Game
func NewMovesHandler(db *driver.DB) *MovesHandler {
	return &MovesHandler{
		db: gamecrud.NewMovesEngine(db.SQL),
	}
}

type MovesHandler struct {
	db gamecrud.MovesEngine
}

// A Player can send moves
// Validate a player is part of the game
// Validate they own the country
// Validate time/phase
// Count the pieces_moves table, when the records are complete or when time expires update the pieces_moves.location_resolved
// Update the game phase, year and phase_end based on the orders_interval
func (g *MovesHandler) CreateOrUpdate(w http.ResponseWriter, r *http.Request) {
	moveInput := model.MoveInput{}

	decoder := json.NewDecoder(r.Body)

	// UnmarshalJSON() is called on GameInput
	// will do some data validation
	err := decoder.Decode(&moveInput)
	if err != nil {
		respondwithJSON(w, http.StatusInternalServerError, model.ErrorMessage{Message: fmt.Sprintf("%v\n", err)})
		return
	}
	// t := model.Territory("NOS")
	// b := model.Territory("SYR")
	// game.AddGameSquares()
	// fmt.Printf("****%+v \n", t.ValidSeaMovement(b))

	_, err = g.db.Create(r.Context(), &moveInput)

	// fmt.Printf("1 ****%+v \n", id)

	if err != nil {
		respondwithJSON(w, http.StatusInternalServerError, model.ErrorMessage{Message: fmt.Sprintf("%v\n", err)})
		return
	}

	respondwithJSON(w, http.StatusCreated, map[string]string{"message": "Successfully Created"})

}
