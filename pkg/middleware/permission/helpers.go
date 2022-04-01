package permission

import (
	"context"
	"github.com/ehsundar/dopamine/internal/auth/token"
	"github.com/ehsundar/dopamine/pkg/middleware/auth"
)

func getSubject(ctx context.Context) *token.Subject {
	s, ok := ctx.Value(auth.UserSubjectContextKey).(*token.Subject)
	if !ok {
		return nil
	}
	return s
}
