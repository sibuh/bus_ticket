package paseto

import (
	"event_ticket/internal/utils/token"
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type pasetoMaker struct {
	signingKey []byte
	paseto     *paseto.V2
}

func NewPaetoMaker(key string) (token.TokenMaker, error) {
	if len(key) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	return &pasetoMaker{
		signingKey: []byte(key),
		paseto:     paseto.NewV2(),
	}, nil
}

func (p *pasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload := token.NewPayload(username, duration)

	return p.paseto.Encrypt(p.signingKey, payload, nil)

}
func (p *pasetoMaker) VerifyToken(tokenString string) (*token.Payload, error) {
	var payload token.Payload
	err := p.paseto.Decrypt(tokenString, p.signingKey, &payload, nil)
	if err != nil {
		return nil, err
	}
	if !payload.Valid() {
		return nil, fmt.Errorf("token has been expired")
	}
	return &payload, err
}
