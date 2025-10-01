package auth

import "golang.org/x/crypto/bcrypt"

// hash pw
func HashPW(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hashed), err
}

//check hash

func CheckHashedPW(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil

}
