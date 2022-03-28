package token

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
	"time"
)

var (
	ErrInvalidStandardClaim = errors.New("invalid jwt standard claim")
	ErrInvalidClaimSubject  = errors.New("invalid jwt standard claim subject")
)

type Manager struct {
	signingKey []byte
}

type Subject struct {
	UserID      string
	Superuser   bool
	Permissions []string
}

func NewManager(signingKey []byte) *Manager {
	return &Manager{signingKey: signingKey}
}

func (tg *Manager) Generate(subject *Subject) (string, error) {
	subjectStr := ""
	if subject != nil {
		j, err := json.Marshal(subject)
		if err != nil {
			return "", err
		}
		subjectStr = string(j)
	}
	claims := jwt.StandardClaims{
		Audience:  "",
		ExpiresAt: time.Now().AddDate(0, 0, 3).Unix(),
		Id:        subject.UserID,
		IssuedAt:  time.Now().Unix(),
		Issuer:    "dopamine",
		NotBefore: time.Now().Unix(),
		Subject:   subjectStr,
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(tg.signingKey)
}

func (tg *Manager) Validate(token string) (*Subject, error) {
	tk, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return tg.signingKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := tk.Claims.(*jwt.StandardClaims); ok && tk.Valid {
		subject := Subject{}
		err = json.Unmarshal([]byte(claims.Subject), &subject)
		if err != nil {
			log.WithError(err).Error(ErrInvalidClaimSubject)
		}
		return &subject, err
	} else {
		return nil, ErrInvalidStandardClaim
	}
}
