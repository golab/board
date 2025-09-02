package verify

import (
	"golang.org/x/crypto/bcrypt"
)

func Hash(input string) string {
	hashedBytes, _ := bcrypt.GenerateFromPassword(
		[]byte(input),
		bcrypt.DefaultCost)
	return string(hashedBytes)
}

func CorrectPassword(input, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(input))
	return err == nil
}
