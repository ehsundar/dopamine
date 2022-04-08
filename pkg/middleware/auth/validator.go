package auth

import (
	"net/http"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"

	"github.com/ehsundar/dopamine/api/auth/token"
)

func WithLoginRequired(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subject, ok := r.Context().Value(UserSubjectContextKey).(*token.Subject)
		if subject == nil || !ok {
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte("login required"))
			if err != nil {
				log.WithError(err).Error("can not write response")
			}
			return
		}

		next(w, r)
	}
}

func WithPermissions(next http.HandlerFunc, perms ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subject, ok := r.Context().Value(UserSubjectContextKey).(*token.Subject)
		if !ok {
			log.WithContext(r.Context()).Error("subject is not of type *token.Subject")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if subject.Superuser {
			next(w, r)
			return
		}

		authorized := lo.Some(perms, subject.Permissions)
		if !authorized {
			log.Warning("unauthorized request")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
