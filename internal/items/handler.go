package items

import (
	"github.com/ehsundar/dopamine/pkg/middleware/permission"
	"github.com/ehsundar/dopamine/pkg/storage"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
)

type handler struct {
	s storage.Storage
}

func RegisterHandlers(router *mux.Router, s storage.Storage) {
	hnd := &handler{
		s: s,
	}

	router.HandleFunc("/{table}/",
		permission.Middleware(hnd.HandleList, permission.List)).Methods("GET")
	router.HandleFunc("/{table}/",
		permission.Middleware(hnd.HandleInsertOne, permission.Create)).Methods("POST")
	router.HandleFunc("/{table}/{id:[0-9]+}/",
		permission.Middleware(hnd.HandleRetrieveOne, permission.Retrieve)).Methods("GET")
	router.HandleFunc("/{table}/{id:[0-9]+}/",
		permission.Middleware(hnd.HandleUpdateOne, permission.Update)).Methods("PUT")
	router.HandleFunc("/{table}/{id:[0-9]+}/",
		permission.Middleware(hnd.HandleDeleteOne, permission.Delete)).Methods("DELETE")
}

func (h *handler) HandleList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	table := vars["table"]

	items, err := h.s.GetAll(r.Context(), table)
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

func (h *handler) HandleInsertOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	table := vars["table"]

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

	i, err = h.s.InsertOne(r.Context(), table, i)
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

func (h *handler) HandleRetrieveOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	table := vars["table"]
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.WithContext(r.Context()).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	i, err := h.s.GetOne(r.Context(), table, id)
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

func (h *handler) HandleUpdateOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	table := vars["table"]
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

	i, err = h.s.UpdateOne(r.Context(), table, i)
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

func (h *handler) HandleDeleteOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	table := vars["table"]
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.WithContext(r.Context()).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.s.DeleteOne(r.Context(), table, id)
	if err != nil {
		log.WithContext(r.Context()).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
