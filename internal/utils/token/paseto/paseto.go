package paseto

import (
	"event_ticket/internal/model"
	"event_ticket/internal/utils/token"
	"fmt"
	"log"
	"net/http"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type pasetoMaker struct {
	signingKey []byte
	paseto     *paseto.V2
}

func NewPasetoMaker(key string) token.TokenMaker {
	if len(key) != chacha20poly1305.KeySize {
		log.Default().Println("Wrong size signing key")
	}

	return &pasetoMaker{
		signingKey: []byte(key),
		paseto:     paseto.NewV2(),
	}
}

func (p *pasetoMaker) CreateToken(payload token.TokenValidator) (string, error) {
	tokenString, err := p.paseto.Encrypt(p.signingKey, payload, nil)
	if err != nil {
		newErr := model.Error{
			ErrCode:   http.StatusInternalServerError,
			Message:   "failed to create token by paseto encryption",
			RootError: err,
		}
		return "", &newErr
	}
	return tokenString, nil
}
func (p *pasetoMaker) VerifyToken(tokenString string, payload token.TokenValidator) (token.TokenValidator, error) {
	err := p.paseto.Decrypt(tokenString, p.signingKey, &payload, nil)
	if err != nil {
		return nil, err
	}
	fmt.Println("decrypted payload:-->", payload)
	er := payload.IsValid()

	if er != nil {
		return nil, er
	}

	return payload, err
}
