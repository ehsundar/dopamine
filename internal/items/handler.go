package items

import (
	"github.com/ehsundar/dopamine/internal/auth/token"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/ehsundar/dopamine/pkg/storage"
)

type Handler struct {
	s storage.Storage
}

func NewHandler(router *mux.Router, s storage.Storage, manager *token.Manager) *Handler {
	hnd := &Handler{
		s: s,
	}

	router.HandleFunc("/{namespace}/", hnd.HandleList).Methods("GET")
	router.HandleFunc("/{namespace}/", hnd.HandleInsertOne).Methods("POST")

	router.HandleFunc("/{namespace}/{id}/", hnd.HandleRetrieveOne).Methods("GET")
	router.HandleFunc("/{namespace}/{id}/", hnd.HandleUpdateOne).Methods("PUT")
	router.HandleFunc("/{namespace}/{id}/", hnd.HandleDeleteOne).Methods("DELETE")

	return hnd
}

func (h *Handler) HandleList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]

	items, err := h.s.GetAll(r.Context(), namespace)
	if err != nil {
		log.WithContext(r.Context()).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := storage.ItemsToJSON(items, true)
	if err != nil {
		log.WithContext(r.Context()).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(result)
	if err != nil {
		log.WithContext(r.Context()).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandleInsertOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithContext(r.Context()).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	i, err := storage.ItemFromJSON(body)
	if err != nil {
		log.WithContext(r.Context()).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	i, err = h.s.InsertOne(r.Context(), namespace, i)
	if err != nil {
		log.WithContext(r.Context()).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := i.ToJSON(true)
	if err != nil {
		log.WithContext(r.Context()).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(result)
	if err != nil {
		log.WithContext(r.Context()).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandleRetrieveOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.WithContext(r.Context()).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	i, err := h.s.GetOne(r.Context(), namespace, id)
	if err != nil {
		log.WithContext(r.Context()).Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	result, err := i.ToJSON(true)
	if err != nil {
		log.WithContext(r.Context()).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(result)
	if err != nil {
		log.WithContext(r.Context()).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandleUpdateOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.WithContext(r.Context()).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithContext(r.Context()).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	i, err := storage.ItemFromJSON(body)
	if err != nil {
		log.WithContext(r.Context()).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	i.ID = id

	i, err = h.s.UpdateOne(r.Context(), namespace, i)
	if err != nil {
		log.WithContext(r.Context()).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result, err := i.ToJSON(true)
	if err != nil {
		log.WithContext(r.Context()).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(result)
	if err != nil {
		log.WithContext(r.Context()).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandleDeleteOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.WithContext(r.Context()).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.s.DeleteOne(r.Context(), namespace, id)
	if err != nil {
		log.WithContext(r.Context()).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
