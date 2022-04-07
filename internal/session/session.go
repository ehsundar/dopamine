package session

import (
	"context"
	"fmt"
	"github.com/ehsundar/dopamine/pkg/storage"
)

const sessionsTableName = "sessions"

type SessionManager struct {
	s storage.Storage
}

func NewSessionManager(s storage.Storage) *SessionManager {
	return &SessionManager{s: s}
}

func (m *SessionManager) Create(ctx context.Context, username string) (sessionID string, err error) {
	s, err := m.s.InsertOne(ctx, sessionsTableName, &storage.Item{
		Contents: map[string]any{
			"username": username,
		},
	})
	if err != nil {
		return
	}
	sessionID = fmt.Sprintf("%d", s.ID)
	return
}

//
//func (m *SessionManager) RevokeUserSessions(ctx context.Context, username string) (err error) {
//	s, err := m.s.InsertOne(ctx, sessionsTableName, &storage.Item{
//		Contents: map[string]any{
//			"username": username,
//		},
//	})
//	if err != nil {
//		return
//	}
//	sessionID = fmt.Sprintf("%d", s.ID)
//	return
//}
