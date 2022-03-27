package main

import (
	"github.com/ehsundar/dopamine/internal/auth"
	"github.com/ehsundar/dopamine/internal/auth/token"
	"github.com/spf13/viper"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"

	"github.com/ehsundar/dopamine/internal/items"
	authMw "github.com/ehsundar/dopamine/pkg/middleware/auth"
	"github.com/ehsundar/dopamine/pkg/storage"
)

func main() {
	LoadConfig()

	s := storage.NewSqliteStorage(viper.GetString("dsn"))

	signingKey := token.LoadSigningKey()
	tokenManager := token.NewManager(signingKey)

	router := mux.NewRouter()

	items.RegisterHandlers(router, s)
	auth.RegisterHandlers(router, s, tokenManager)

	handlerFunc := http.Handler(router).ServeHTTP
	httpServer := authMw.NewAuthMiddleware(handlerFunc, tokenManager)

	serverAddr := "0.0.0.0:8080"
	log.Infof("serving http server at %s", serverAddr)
	err := http.ListenAndServe(serverAddr, httpServer)
	if err != nil {
		log.Error(err)
		panic(err)
	}
}
