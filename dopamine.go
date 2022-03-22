package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/ehsundar/dopamine/internal/items"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./dopamine.db")
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	_ = items.NewHandler(router, db)

	err = http.ListenAndServe("0.0.0.0:8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
