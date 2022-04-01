package permission

import (
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Middleware(next http.HandlerFunc, extractor Extractor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subject := getSubject(r.Context())
		permissionForTable := extractor(r)

		switch permissionForTable {
		case "public":
			break
		case "superuser":
			if subject == nil || !subject.Superuser {
				w.WriteHeader(http.StatusUnauthorized)
				log.Infof("unauthorized request on superuser api")
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
