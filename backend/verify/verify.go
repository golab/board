package verify

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"os"
)

func VerifyWithKey(message, sig string, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(message))
	expectedMAC := mac.Sum(nil)
	messageMAC, _ := hex.DecodeString(sig)
	return hmac.Equal(messageMAC, expectedMAC)
}

func Verify(message, sig string) bool {
	key := GetKey()
	if len(key) == 0 {
		return true
	}
	return VerifyWithKey(message, sig, key)
}

func GetKey() []byte {
	hexString := os.Getenv("LOCALKEY")
	if len(hexString) == 0 {
		return []byte{}
	}
	byteArray, _ := hex.DecodeString(hexString)
	return byteArray
}

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
