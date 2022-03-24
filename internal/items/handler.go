package items

import (
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

func NewHandler(router *mux.Router, s storage.Storage) *Handler {
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
		w.WriteHeader(http.StatusInternalServerError)
		log.WithContext(r.Context()).Error(err)
		return
	}

	result, err := storage.ItemsToJSON(items, true)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.WithContext(r.Context()).Error(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.WithContext(r.Context()).Error(err)
		return
	}
}

func (h *Handler) HandleInsertOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.WithContext(r.Context()).Error(err)
		return
	}

	i, err := storage.ItemFromJSON(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.WithContext(r.Context()).Error(err)
		return
	}

	i, err = h.s.InsertOne(r.Context(), namespace, i)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.WithContext(r.Context()).Error(err)
		return
	}

	result, err := i.ToJSON(true)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.WithContext(r.Context()).Error(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.WithContext(r.Context()).Error(err)
		return
	}
}

func (h *Handler) HandleRetrieveOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.WithContext(r.Context()).Error(err)
		return
	}

	i, err := h.s.GetOne(r.Context(), namespace, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.WithContext(r.Context()).Error(err)
		return
	}

	result, err := i.ToJSON(true)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.WithContext(r.Context()).Error(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.WithContext(r.Context()).Error(err)
		return
	}
}

func (h *Handler) HandleUpdateOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.WithContext(r.Context()).Error(err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.WithContext(r.Context()).Error(err)
		return
	}

	i, err := storage.ItemFromJSON(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.WithContext(r.Context()).Error(err)
		return
	}
	i.ID = id

	i, err = h.s.UpdateOne(r.Context(), namespace, i)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.WithContext(r.Context()).Error(err)
		return
	}
	result, err := i.ToJSON(true)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.WithContext(r.Context()).Error(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.WithContext(r.Context()).Error(err)
		return
	}
}

func (h *Handler) HandleDeleteOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.WithContext(r.Context()).Error(err)
		return
	}

	err = h.s.DeleteOne(r.Context(), namespace, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.WithContext(r.Context()).Error(err)
		return
	}
}
