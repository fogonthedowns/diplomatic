package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
	"os"

	"github.com/fogonthedowns/diplomatic/driver"
	gameHandler "github.com/fogonthedowns/diplomatic/handler/http"
)

func main() {
	dbName := os.Getenv("DB_NAME")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	fmt.Printf("dbHost%v\n", dbHost)

	connection, err := driver.ConnectSQL(dbHost, dbPort, "root", dbPass, dbName)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	h := gameHandler.NewGameHandler(connection)
	mh := gameHandler.NewMovesHandler(connection)
	r.Route("/", func(rt chi.Router) {
		rt.Mount("/games", postRouter(h))
		rt.Mount("/moves", movesRouter(mh))
	})

	fmt.Println("Server listen at :8005")
	http.ListenAndServe(":8005", r)
}

// A completely separate router for posts routes
func postRouter(game *gameHandler.GameHandler) http.Handler {
	r := chi.NewRouter()
	r.Get("/", game.Fetch)
	r.Get("/{[0-9]+}", game.GetByID)
	r.Post("/", game.Create)
	r.Put("/{[0-9]+}", game.Update)
	// r.Delete("/{id:[0-9]+}", pHandler.Delete)

	return r
}

// A completely separate router for game mvoes
func movesRouter(move *gameHandler.MovesHandler) http.Handler {
	r := chi.NewRouter()
	r.Post("/", move.CreateOrUpdate)

	return r
}
