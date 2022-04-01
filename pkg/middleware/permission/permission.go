package permission

import (
	"github.com/gorilla/mux"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

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
