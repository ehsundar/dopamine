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
	Superuser bool
}

func NewManager(signingKey []byte) *Manager {
	return &Manager{signingKey: signingKey}
}

func (tg *Manager) Generate(userID string, subject *Subject) (string, error) {
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
		Id:        userID,
		IssuedAt:  time.Now().Unix(),
		Issuer:    "dopamine",
		NotBefore: time.Now().Unix(),
		Subject:   subjectStr,
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(tg.signingKey)
}

func (tg *Manager) Validate(token string) (*Subject, error) {
	subject := Subject{}

	tk, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return tg.signingKey, nil
	})

	if claims, ok := tk.Claims.(jwt.StandardClaims); ok && tk.Valid {
		err = json.Unmarshal([]byte(claims.Subject), &subject)
		log.WithError(err).Error(ErrInvalidClaimSubject)
		return &subject, ErrInvalidClaimSubject
	} else {
		return nil, ErrInvalidStandardClaim
	}
}
