package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Payload contains the payload data of the token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayLoad creates a new token payload with specific duration and username
func NewPayLoad(username string, duration time.Duration) (*Payload, error) {

	//creating random uuid
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	//creates the payload
	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

// Valid checks if the token payload is laid or not
func (payload *Payload) Valid() error {
	if time.Now().After((payload.ExpiredAt)) {
		return errors.New("token has expired")
	}
	return nil
}
