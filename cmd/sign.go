package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"
)

func SignWithKey(message string, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

func Sign(message string) string {
	key := GetKey()
	if len(key) == 0 {
		return ""
	}
	return SignWithKey(message, key)
}

func GetKey() []byte {
	hexString := os.Getenv("LOCALKEY")
	if len(hexString) == 0 {
		return []byte{}
	}
	byteArray, _ := hex.DecodeString(hexString)
	return byteArray
}
