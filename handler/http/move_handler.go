package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	// "strconv"
	// "strings"

	// "github.com/go-chi/chi"

	db "diplomacy/db"
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
	db db.Crud
}

// A Player can send moves
// Validate a player is part of the game
// Validate they own the country
// Validate time/phase
// Count the pieces_moves table, when the records are complete or when time expires update the pieces_moves.location_resolved
// Update the game phase, year and phase_end based on the orders_interval
func (g *MovesHandler) CreateOrUpdate(w http.ResponseWriter, r *http.Request) {
	gameInput := model.GameInput{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&gameInput)
	if err != nil {
		panic(err)
	}
	// t := model.Territory("NOS")
	// b := model.Territory("SYR")
	// game.AddGameSquares()
	// fmt.Printf("****%+v \n", t.ValidSeaMovement(b))

	id, err := g.db.Create(r.Context(), &gameInput)

	fmt.Printf("1 ****%+v \n", id)

	if err != nil {
		respondwithJSON(w, http.StatusInternalServerError, model.ErrorMessage{Message: fmt.Sprintf("%v\n", err)})
	} else {
		respondwithJSON(w, http.StatusCreated, map[string]string{"message": "Successfully Created"})
	}
}
