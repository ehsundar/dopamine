package items

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/ehsundar/dopamine/pkg/middleware/auth"
	"github.com/ehsundar/dopamine/pkg/storage"
)

type handler struct {
	s storage.Storage
}

func getPermissionForTable(table string, apiType string) string {
	tablesConfig := viper.Sub("tables")
	tb := tablesConfig.GetStringMapString(table)
	perm, ok := tb[apiType]

	if ok {
		return perm
	} else {
		if table == "default" {
			return "public"
		} else {
			return getPermissionForTable("default", apiType)
		}
	}
}

func PermissionsMiddleware(next http.HandlerFunc, apiType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		table := vars["table"]

		subject := auth.GetSubject(r.Context())
		permissionForTable := getPermissionForTable(table, apiType)

		switch permissionForTable {
		case "public":
			break
		case "superuser":
			if subject == nil || !subject.Superuser {
				w.WriteHeader(http.StatusUnauthorized)
				log.Infof("unauthorized request on superuser api: %s -> %s", apiType, table)
				return
			}
			break
		default:
			if !lo.Contains(subject.Permissions, permissionForTable) {
				w.WriteHeader(http.StatusUnauthorized)
				log.Infof("unauthorized request: not enough permission: needed: %s, having: %s",
					permissionForTable, subject.Permissions)
				return
			}
			break
		}

		next(w, r)
	}
}

func RegisterHandlers(router *mux.Router, s storage.Storage) {
	hnd := &handler{
		s: s,
	}

	router.HandleFunc("/{table}/", PermissionsMiddleware(hnd.HandleList, auth.APIList)).Methods("GET")
	router.HandleFunc("/{table}/", hnd.HandleInsertOne).Methods("POST")

	router.HandleFunc("/{table}/{id:[0-9]+}/", hnd.HandleRetrieveOne).Methods("GET")
	router.HandleFunc("/{table}/{id:[0-9]+}/", hnd.HandleUpdateOne).Methods("PUT")
	router.HandleFunc("/{table}/{id:[0-9]+}/", hnd.HandleDeleteOne).Methods("DELETE")
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
