package main

import (
	"github.com/ehsundar/dopamine/pkg/storage"
	"log"
	"net/http"

	"github.com/ehsundar/dopamine/internal/items"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	s := storage.NewSqliteStorage("./dopamine.db")

	router := mux.NewRouter()

	_ = items.NewHandler(router, s)

	err := http.ListenAndServe("0.0.0.0:8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
