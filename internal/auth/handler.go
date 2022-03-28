package auth

import (
	"github.com/ehsundar/dopamine/internal/auth/token"
	"github.com/ehsundar/dopamine/pkg/storage"
	"github.com/gorilla/mux"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	urlPrefix = "/auth"
)

type handler struct {
	s       storage.Storage
	manager *token.Manager
}

func RegisterHandlers(router *mux.Router, s storage.Storage, manager *token.Manager) {
	hnd := &handler{
		s:       s,
		manager: manager,
	}

	router.HandleFunc(urlPrefix+"/authenticate/", hnd.HandleAuthenticate).Methods("POST")
}

func (h handler) HandleAuthenticate(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Error("can not read request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req := AuthenticateRequest{}
	err = req.Parse(body)
	if err != nil {
		log.WithError(err).Error("can not read request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(req.Username)
	if err != nil {
		log.WithError(err).Error("username must be an int")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	item, err := h.s.GetOne(r.Context(), "users", userID)
	if err == storage.ErrTableNotExist {
		log.WithError(err).Error("no user registered")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		log.WithError(err).Error("can not get user item")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if item.Contents["password"] != req.Password {
		log.WithError(err).Error("invalid password")
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	superuser := item.Contents["superuser"].(bool)
	permissionsAny := item.Contents["permissions"].([]interface{})
	permissions := lo.Map(permissionsAny, func(p any, _ int) string {
		return p.(string)
	})

	subject := &token.Subject{
		UserID:      req.Username,
		Superuser:   superuser,
		Permissions: permissions,
	}
	tk, err := h.manager.Generate(subject)
	if err != nil {
		log.WithError(err).Error("error generating jwt")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := AuthenticateResponse{Token: tk}

	w.Header().Set("Content-Type", "application/json")
	err = resp.Render(w)
	if err != nil {
		log.WithContext(r.Context()).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
