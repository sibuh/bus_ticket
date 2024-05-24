package pass

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pswd string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pswd), bcrypt.DefaultCost)
	return string(b), err
}
func CheckHashPassword(pswd, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pswd))
	return err == nil

}
