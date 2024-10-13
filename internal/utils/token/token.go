package token

type TokenValidator interface {
	IsValid() error
}

type TokenMaker interface {
	CreateToken(payload TokenValidator) (string, error)
	VerifyToken(tokenString string, payload TokenValidator) (TokenValidator, error)
}
