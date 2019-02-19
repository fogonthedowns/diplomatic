package main

import (
    "github.com/go-chi/chi"
    "github.com/go-chi/chi/middleware"
    _ "github.com/go-sql-driver/mysql"
)

var router *chi.Mux
var db *sql.DB

const (
    dbName = "go-mysql-crud"
    dbPass = "12345"
    dbHost = "localhost"
    dbPort = "33066"
)

func routers() *chi.Mux {
    router.Get("/gamess", AllGames)
    router.Get("/games/{id}", DetailGame)
    router.Post("/games", CreateGame)
    router.Put("/games/{id}", UpdateGame)
    router.Delete("/games/{id}", DeleteGame)
    
    return router
}

func init() { 
    router = chi.NewRouter() 
    router.Use(middleware.Recoverer)  
    
    dbSource := fmt.Sprintf("root:%s@tcp(%s:%s)/%s?charset=utf8",  dbPass, dbHost, dbPort, dbName)
    
    var err error
    db, err = sql.Open("mysql", dbSource) 
    
    catch(err)
}

