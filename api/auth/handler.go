package auth

import (
	"github.com/ehsundar/dopamine/api/auth/token"
	"github.com/ehsundar/dopamine/pkg/middleware/permission"
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

	grant := permission.Middleware(hnd.HandleGrantPermission, permission.StaticExtractorFactory("permissions.create"))
	drop := permission.Middleware(hnd.HandleDropPermission, permission.StaticExtractorFactory("permissions.delete"))

	router.HandleFunc(urlPrefix+"/authenticate/", hnd.HandleAuthenticate).Methods("POST")
	router.HandleFunc(urlPrefix+"/permissions/", grant).Methods("POST")
	router.HandleFunc(urlPrefix+"/permissions/", drop).Methods("DELETE")
}

func (h *handler) HandleAuthenticate(w http.ResponseWriter, r *http.Request) {
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
		log.WithError(err).Error("can not get user")
		w.WriteHeader(http.StatusNotFound)
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

func (h *handler) HandleGrantPermission(w http.ResponseWriter, r *http.Request) {
	req := &GrantPermissionRequest{}
	if err := req.Parse(r.Body); err != nil {
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

	user, err := h.s.GetOne(r.Context(), "users", userID)
	if err != nil {
		log.WithError(err).Error("can not get user")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	permissionsAny := user.Contents["permissions"].([]interface{})
	permissions := lo.Map(permissionsAny, func(p any, _ int) string {
		return p.(string)
	})

	user.Contents["permissions"] = lo.Uniq(append(permissions, req.Permission))
	_, err = h.s.UpdateOne(r.Context(), "users", user)
	if err != nil {
		log.WithError(err).Error("can not update user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	return
}

func (h *handler) HandleDropPermission(w http.ResponseWriter, r *http.Request) {
	req := &DropPermissionRequest{}
	if err := req.Parse(r.Body); err != nil {
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

	user, err := h.s.GetOne(r.Context(), "users", userID)
	if err != nil {
		log.WithError(err).Error("can not get user")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	permissionsAny := user.Contents["permissions"].([]interface{})
	permissions := lo.Map(permissionsAny, func(p any, _ int) string {
		return p.(string)
	})

	user.Contents["permissions"] = lo.Filter(permissions, func(p string, _ int) bool {
		return p != req.Permission
	})
	_, err = h.s.UpdateOne(r.Context(), "users", user)
	if err != nil {
		log.WithError(err).Error("can not update user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	return
}
