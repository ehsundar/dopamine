package auth

import (
	"context"
	"github.com/ehsundar/dopamine/internal/auth/token"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type ContextKey string

const UserSubjectContextKey ContextKey = "authUserSubject"

func Middleware(next http.HandlerFunc, manager *token.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			next(w, r)
			return
		}

		parts := strings.Split(header, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte("invalid authorization header"))
			if err != nil {
				log.WithError(err).Error("unable to write response")
			}
			return
		}
		tk := parts[1]

		subject, err := manager.Validate(tk)
		if err != nil {
			log.WithError(err).Error("invalid jwt")
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte("invalid jwt"))
			if err != nil {
				log.WithError(err).Error("unable to write response")
			}
			return
		}

		newCtx := context.WithValue(r.Context(), UserSubjectContextKey, subject)
		next(w, r.WithContext(newCtx))
	}
}
