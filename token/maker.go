package token

import "time"

// Maker is an interface for managing tokens
type Maker interface {
	//CreateToken creates a new token for a specific username and the duration
	CreateToken(username string, duration time.Duration) (string, error)
	//VerifyToken is to validate the token
	VerifyToken(token string) (*Payload, error)
}
