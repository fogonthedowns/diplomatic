package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	// "github.com/go-chi/chi"

	// db "diplomacy/db/engine"
	gamecrud "diplomacy/db/gamecrud"
	"diplomacy/driver"
	model "diplomacy/model"
)

// NewGameHandler() creates a new HTTP handler
// To handle web requests for Game
func NewGameHandler(db *driver.DB) *GameHandler {
	return &GameHandler{
		db: gamecrud.NewEngine(db.SQL),
	}
}

type GameHandler struct {
	db gamecrud.Engine
}

// Creates a new game
// default password game is false
// writes a game row
// TODO respond with ID!!
func (g *GameHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	if err != nil {
		respondwithJSON(w, http.StatusInternalServerError, model.ErrorMessage{Message: fmt.Sprintf("%v", err)})
	} else {
		respondwithJSON(w, http.StatusCreated, map[string]int64{"id": id})
	}
}

// Used to Join Game
// Update a post by game id
// Read UserGames Table
// Counts the number of Users in a game
// when the Game is full, return an error
// does not allow duplicate countries

// Create Piece records, setting the user.id
// Create Territory records, setting the user.id
// TODO(:3/1) Go over every controller and ensure validations are set up correctly
func (g *GameHandler) Update(w http.ResponseWriter, r *http.Request) {
	p := strings.Split(r.URL.Path, "/")

	var id int
	if len(p) == 3 {
		id, _ = strconv.Atoi(p[2])
	}

	// id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	data := model.GameInput{Id: int64(id)}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		respondwithJSON(w, http.StatusInternalServerError, model.ErrorMessage{Message: fmt.Sprintf("%v\n", err)})
		return
	}

	payload, code, err := g.db.Update(r.Context(), &data)

	if err != nil {
		switch code {
		case 409:
			respondwithJSON(w, http.StatusConflict, map[string]string{"message": fmt.Sprintf("%v\n", err)})
		default:
			respondwithJSON(w, http.StatusInternalServerError, model.ErrorMessage{Message: fmt.Sprintf("%v\n", err)})
		}
	} else {
		respondwithJSON(w, http.StatusOK, payload)
	}
}

type ErrorMessage struct {
	message string `json:"message"`
}

// respondwithJSON write json response format
func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("error: %v \n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// respondwithError return error message
func respondWithError(w http.ResponseWriter, code int, msg string) {
	em := ErrorMessage{message: msg}
	respondwithJSON(w, code, em)
}

// // Fetch all post data
func (g *GameHandler) Fetch(w http.ResponseWriter, r *http.Request) {
	payload, _ := g.db.Fetch(r.Context(), 5)
	respondwithJSON(w, http.StatusOK, payload)
}

// GetByID returns a post details
func (g *GameHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	var id int
	if len(path) == 3 {
		id, _ = strconv.Atoi(path[2])
	}

	payload, err := g.db.GetByID(r.Context(), int64(id))

	if err != nil {
		respondWithError(w, http.StatusNoContent, "Content not found")
	}

	respondwithJSON(w, http.StatusOK, payload)
}
