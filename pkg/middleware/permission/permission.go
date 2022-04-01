package permission

import (
	"github.com/gorilla/mux"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Middleware(next http.HandlerFunc, apiType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		table := vars["table"]

		subject := getSubject(r.Context())
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
