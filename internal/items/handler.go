package items

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(router *mux.Router, db *sql.DB) *Handler {
	hnd := &Handler{
		db: db,
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

	err := createTable(r.Context(), h.db, namespace)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.WithContext(r.Context()).Error(err)
		return
	}

	items, err := listItems(r.Context(), h.db, namespace)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.WithContext(r.Context()).Error(err)
		return
	}

	mapped := lo.Map(items, func(i *Item, _ int) map[string]any {
		m := make(map[string]any)
		json.Unmarshal([]byte(i.Contents), &m)
		m["id"] = i.ID
		return m
	})

	result, err := json.Marshal(mapped)
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

	err := createTable(r.Context(), h.db, namespace)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.WithContext(r.Context()).Error(err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.WithContext(r.Context()).Error(err)
		return
	}

	bodyJson := make(map[string]any)
	err = json.Unmarshal(body, &bodyJson)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.WithContext(r.Context()).Error(err)
		return
	}

	bodyClean, err := json.Marshal(bodyJson)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.WithContext(r.Context()).Error(err)
		return
	}

	i, err := insertOneItem(r.Context(), h.db, namespace, string(bodyClean))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.WithContext(r.Context()).Error(err)
		return
	}

	m := make(map[string]any)
	err = json.Unmarshal([]byte(i.Contents), &m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.WithContext(r.Context()).Error(err)
		return
	}

	m["id"] = i.ID
	m["created_at"] = i.CreatedAt

	result, err := json.Marshal(m)
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

	err = createTable(r.Context(), h.db, namespace)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.WithContext(r.Context()).Error(err)
		return
	}

	i, err := getItem(r.Context(), h.db, namespace, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.WithContext(r.Context()).Error(err)
		return
	}

	m := make(map[string]any)
	err = json.Unmarshal([]byte(i.Contents), &m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.WithContext(r.Context()).Error(err)
		return
	}

	m["id"] = i.ID
	m["created_at"] = i.CreatedAt

	result, err := json.Marshal(m)
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

	bodyJson := make(map[string]any)
	err = json.Unmarshal(body, &bodyJson)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.WithContext(r.Context()).Error(err)
		return
	}

	bodyClean, err := json.Marshal(bodyJson)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.WithContext(r.Context()).Error(err)
		return
	}

	i, err := updateItem(r.Context(), h.db, namespace, id, string(bodyClean))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.WithContext(r.Context()).Error(err)
		return
	}

	m := make(map[string]any)
	err = json.Unmarshal([]byte(i.Contents), &m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.WithContext(r.Context()).Error(err)
		return
	}

	m["id"] = i.ID
	m["created_at"] = i.CreatedAt

	result, err := json.Marshal(m)
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

	err = deleteItem(r.Context(), h.db, namespace, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.WithContext(r.Context()).Error(err)
		return
	}
}