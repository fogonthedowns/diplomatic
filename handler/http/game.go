package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	// "strconv"

	// "github.com/go-chi/chi"

	db "diplomacy/db"
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
	db db.Crud
}

// Create a new post
func (g *GameHandler) Create(w http.ResponseWriter, r *http.Request) {
	game := model.Game{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&game)
	if err != nil {
		panic(err)
	}
	t := model.Territory("NOS")
	b := model.Territory("SYR")
	fmt.Printf("****%+v \n", t.ValidShipMovement(b))

	id, err := g.db.Create(r.Context(), &game)

	fmt.Printf("1 ****%+v \n", id)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Error")
	}

	respondwithJSON(w, http.StatusCreated, map[string]string{"message": "Successfully Created"})
}

// respondwithJSON write json response format
func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// respondwithError return error message
func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondwithJSON(w, code, map[string]string{"message": msg})
}

// // Fetch all post data
// func (p *GameDb) Fetch(w http.ResponseWriter, r *http.Request) {
// 	payload, _ := p.repo.Fetch(r.Context(), 5)

// 	respondwithJSON(w, http.StatusOK, payload)
// }

// // Update a post by id
// func (p *Post) Update(w http.ResponseWriter, r *http.Request) {
// 	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
// 	data := models.Post{ID: int64(id)}
// 	json.NewDecoder(r.Body).Decode(&data)
// 	payload, err := p.repo.Update(r.Context(), &data)

// 	if err != nil {
// 		respondWithError(w, http.StatusInternalServerError, "Server Error")
// 	}

// 	respondwithJSON(w, http.StatusOK, payload)
// }

// // GetByID returns a post details
// func (p *Post) GetByID(w http.ResponseWriter, r *http.Request) {
// 	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
// 	payload, err := p.repo.GetByID(r.Context(), int64(id))

// 	if err != nil {
// 		respondWithError(w, http.StatusNoContent, "Content not found")
// 	}

// 	respondwithJSON(w, http.StatusOK, payload)
// }

// // Delete a post
// func (p *Post) Delete(w http.ResponseWriter, r *http.Request) {
// 	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
// 	_, err := p.repo.Delete(r.Context(), int64(id))

// 	if err != nil {
// 		respondWithError(w, http.StatusInternalServerError, "Server Error")
// 	}

// 	respondwithJSON(w, http.StatusMovedPermanently, map[string]string{"message": "Delete Successfully"})
// }
