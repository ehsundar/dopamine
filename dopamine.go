package main

import (
	"github.com/ehsundar/dopamine/internal/auth"
	"github.com/ehsundar/dopamine/internal/auth/token"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"

	"github.com/ehsundar/dopamine/internal/items"
	"github.com/ehsundar/dopamine/pkg/storage"
)

func main() {
	s := storage.NewSqliteStorage("./dopamine.db")

	signingKey := token.LoadSigningKey()
	tokenManager := token.NewManager(signingKey)

	router := mux.NewRouter()

	_ = items.NewHandler(router, s, tokenManager)
	auth.NewHandler(router, s, tokenManager)

	serverAddr := "0.0.0.0:8080"
	log.Infof("serving http server at %s", serverAddr)
	err := http.ListenAndServe(serverAddr, router)
	if err != nil {
		log.Error(err)
		panic(err)
	}
}
