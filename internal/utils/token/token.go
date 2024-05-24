package token

type TokenMaker interface {
	CreateToken(username string) (string, error)
	VerifyToken(token string) (*Payload, error)
}
