package auth

import (
	"context"
	"github.com/ehsundar/dopamine/internal/auth/token"
)

func GetSubject(ctx context.Context) *token.Subject {
	s, ok := ctx.Value(UserSubjectContextKey).(*token.Subject)
	if !ok {
		return nil
	}
	return s
}
