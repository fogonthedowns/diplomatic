package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	gamecrud "github.com/fogonthedowns/diplomatic/db/gamecrud"
	"github.com/fogonthedowns/diplomatic/driver"
	model "github.com/fogonthedowns/diplomatic/model"
)

// NewMovesHandler() creates a new HTTP handler
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
	moveInput := model.Move{}

	decoder := json.NewDecoder(r.Body)

	// UnmarshalJSON() is called on GameInput
	err := decoder.Decode(&moveInput)
	if err != nil {
		respondwithJSON(w, http.StatusInternalServerError, model.ErrorMessage{Message: fmt.Sprintf("%v\n", err)})
		return
	}

	code, err := g.db.CreateOrUpdate(r.Context(), &moveInput)

	switch code {
	case 403:
		respondwithJSON(w, http.StatusForbidden, model.ErrorMessage{Message: fmt.Sprintf("%v\n", err)})
	case 400:
		respondwithJSON(w, http.StatusForbidden, model.ErrorMessage{Message: fmt.Sprintf("%v\n", err)})
	case 200:
		respondwithJSON(w, http.StatusCreated, map[string]string{"message": "Successfully Created"})
	default:
		respondwithJSON(w, http.StatusInternalServerError, model.ErrorMessage{Message: fmt.Sprintf("%v\n", err)})

	}
}
